package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/muhammadisa/bareksanews/model"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
)

func (r *readWrite) WriteTopic(ctx context.Context, req *pb.Topic) (res *pb.Topic, err error) {
	const funcName = `WriteTopic`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	currentTime := time.Now()
	req.CreatedAt = currentTime.Unix()
	req.UpdatedAt = currentTime.Unix()
	stmt, err := r.db.Prepare(queryWriteTopic)
	if err != nil {
		return res, err
	}
	result, err := stmt.ExecContext(
		ctx,
		req.Id,       // id
		req.Title,    // title
		req.Headline, // headline
		currentTime,  // created_at
		currentTime,  // updated_at
	)
	if err != nil {
		return res, err
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return res, fmt.Errorf("failed to insert reason : %+v", err)
	}
	return req, nil
}

func (r *readWrite) ModifyTopic(ctx context.Context, req *pb.Topic) (res *pb.Topic, err error) {
	const funcName = `ModifyTopic`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	var oldTopic model.Topic

	stmt, err := r.db.Prepare(queryLookupCreateAtTopic)
	if err != nil {
		return res, err
	}
	row := stmt.QueryRow(req.Id)
	err = row.Scan(
		&oldTopic.ID,      // id
		&oldTopic.Created, // created_at
	)
	if err != nil || err == sql.ErrNoRows {
		return res, err
	}
	oldTopic.UseUnixTimeStamp()

	currentTime := time.Now()
	req.CreatedAt = oldTopic.CreatedAt
	req.UpdatedAt = currentTime.Unix()
	stmt, err = r.db.Prepare(queryUpdateTopic)
	if err != nil {
		return res, err
	}
	result, err := stmt.ExecContext(
		ctx,
		req.Title,        // title
		req.Headline,     // headline
		oldTopic.Created, // created_at
		currentTime,      // updated_at
		req.Id,           // id
	)
	if err != nil {
		return res, err
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return res, fmt.Errorf("failed to insert reason : %+v", err)
	}
	return req, nil
}

func (r *readWrite) RemoveTopic(ctx context.Context, req *pb.Select) error {
	const funcName = `RemoveTopic`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	stmt, err := r.db.Prepare(queryRemoveTopic)
	if err != nil {
		return err
	}
	result, err := stmt.ExecContext(
		ctx,
		req.Id, // id
	)
	if err != nil {
		return err
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return fmt.Errorf("failed to insert reason : %+v", err)
	}
	return nil
}

func (r *readWrite) ReadTopics(ctx context.Context) (res *pb.Topics, err error) {
	const funcName = `ReadTopics`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	var topics pb.Topics
	var topic model.Topic

	stmt, err := r.db.Prepare(queryReadTopics)
	if err != nil {
		return res, err
	}
	mutex.Lock()
	row, err := stmt.QueryContext(ctx)
	if err != nil {
		return res, err
	}
	mutex.Unlock()
	for row.Next() {
		err = row.Scan(
			&topic.ID,       // id
			&topic.Title,    // title
			&topic.Headline, // headline
			&topic.Created,  // created_at
			&topic.Updated,  // updated_at
		)
		if err != nil {
			return res, err
		}
		topic.UseUnixTimeStamp()
		topics.Topics = append(topics.Topics, &pb.Topic{
			Id:        topic.ID,
			Title:     topic.Title,
			Headline:  topic.Headline,
			CreatedAt: topic.CreatedAt,
			UpdatedAt: topic.UpdatedAt,
		})
	}
	return &topics, nil
}
