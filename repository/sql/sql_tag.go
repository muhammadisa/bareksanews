package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/muhammadisa/bareksanews/model"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
)

func (r *readWrite) WriteTag(ctx context.Context, req *pb.Tag) (res *pb.Tag, err error) {
	const funcName = `WriteTag`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	currentTime := time.Now()
	req.CreatedAt = currentTime.Unix()
	req.UpdatedAt = currentTime.Unix()
	stmt, err := r.db.Prepare(queryWriteTag)
	if err != nil {
		return res, err
	}
	result, err := stmt.ExecContext(
		ctx,
		req.Id,      // id
		req.Tag,     // tag
		currentTime, // created_at
		currentTime, // updated_at
	)
	if err != nil {
		return res, err
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return res, fmt.Errorf("failed to insert reason : %+v", err)
	}
	return req, nil
}

func (r *readWrite) ModifyTag(ctx context.Context, req *pb.Tag) (res *pb.Tag, err error) {
	const funcName = `ModifyTag`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	var oldTag model.Tag

	stmt, err := r.db.Prepare(queryLookupCreateAtTag)
	if err != nil {
		return res, err
	}
	row := stmt.QueryRow(req.Id)
	err = row.Scan(
		&oldTag.ID,      // id
		&oldTag.Created, // created_at
	)
	if err != nil || err == sql.ErrNoRows {
		return res, err
	}
	oldTag.UseUnixTimeStamp()

	currentTime := time.Now()
	req.CreatedAt = oldTag.CreatedAt
	req.UpdatedAt = currentTime.Unix()
	stmt, err = r.db.Prepare(queryUpdateTag)
	if err != nil {
		return res, err
	}
	result, err := stmt.ExecContext(
		ctx,
		req.Tag,        // tag
		oldTag.Created, // created_at
		currentTime,    // updated_at
		req.Id,         // id
	)
	if err != nil {
		return res, err
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return res, fmt.Errorf("failed to insert reason : %+v", err)
	}
	return req, nil
}

func (r *readWrite) RemoveTag(ctx context.Context, req *pb.Select) (err error) {
	const funcName = `RemoveTag`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	stmt, err := r.db.Prepare(queryRemoveTag)
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

func (r *readWrite) ReadTags(ctx context.Context) (res *pb.Tags, err error) {
	const funcName = `ReadTags`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	var tags pb.Tags
	var tag model.Tag

	stmt, err := r.db.Prepare(queryReadTags)
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
			&tag.ID,      // id
			&tag.Tag,     // tag
			&tag.Created, // created_at
			&tag.Updated, // updated_at
		)
		if err != nil {
			return res, err
		}
		tag.UseUnixTimeStamp()
		tags.Tags = append(tags.Tags, &pb.Tag{
			Id:        tag.ID,
			Tag:       tag.Tag,
			CreatedAt: tag.CreatedAt,
			UpdatedAt: tag.UpdatedAt,
		})
	}
	return &tags, nil
}
