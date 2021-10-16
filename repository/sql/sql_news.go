package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/muhammadisa/bareksanews/model"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	uuid "github.com/satori/go.uuid"
)

func (r *readWrite) WriteNews(ctx context.Context, req *pb.News) (res *pb.News, err error) {
	currentTime := time.Now()
	req.Id = uuid.NewV4().String()
	req.CreatedAt = currentTime.Unix()
	req.UpdatedAt = currentTime.Unix()
	stmt, err := r.db.Prepare(queryWriteNews)
	if err != nil {
		return res, err
	}
	result, err := stmt.ExecContext(
		ctx,
		req.Id,      // id
		req.TopicId, // topic_id
		req.Title,   // title
		req.Content, // content
		req.Status,  // status
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

func (r *readWrite) ModifyNews(ctx context.Context, req *pb.News) (res *pb.News, err error) {
	var oldNews model.News
	stmt, err := r.db.Prepare(queryLookupCreateAtNews)
	if err != nil {
		return res, err
	}
	row := stmt.QueryRow(req.Id)
	err = row.Scan(
		&oldNews.ID,      // id
		&oldNews.Created, // created_at
	)
	if err != nil || err == sql.ErrNoRows {
		return res, err
	}
	oldNews.UseUnixTimeStamp()

	currentTime := time.Now()
	req.CreatedAt = oldNews.CreatedAt
	req.UpdatedAt = currentTime.Unix()
	stmt, err = r.db.Prepare(queryUpdateNews)
	if err != nil {
		return res, err
	}
	result, err := stmt.ExecContext(
		ctx,
		req.TopicId,     // topic_id
		req.Title,       // title
		req.Content,     // content
		req.Status,      // status
		oldNews.Created, // created_at
		currentTime,     // updated_at
		req.Id,          // id
	)
	if err != nil {
		return res, err
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return res, fmt.Errorf("failed to insert reason : %+v", err)
	}
	return req, nil
}

func (r *readWrite) RemoveNews(ctx context.Context, req *pb.Select) error {
	stmt, err := r.db.Prepare(queryRemoveNews)
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

func (r *readWrite) rowsNewsesNextAndScan(ctx context.Context, row *sql.Rows) (res *pb.Newses, err error) {
	var newses pb.Newses
	var news model.News
	for row.Next() {
		err = row.Scan(
			&news.ID,      // id
			&news.TopicID, // topic_id
			&news.Title,   // title
			&news.Content, // content
			&news.Status,  // status
			&news.Created, // created_at
			&news.Updated, // updated_at
		)
		if err != nil {
			return res, err
		}
		news.UseUnixTimeStamp()
		newses.Newses = append(newses.Newses, &pb.News{
			Id:           news.ID,
			TopicId:      news.TopicID,
			Title:        news.Title,
			Content:      news.Content,
			Status:       news.Status,
			NewsTagIds:   r.ReadNewsTagsTagIDAndTagByNewsID(ctx, news.ID, false),
			NewsTagNames: r.ReadNewsTagsTagIDAndTagByNewsID(ctx, news.ID, true),
			CreatedAt:    news.CreatedAt,
			UpdatedAt:    news.UpdatedAt,
		})
	}
	return &newses, nil
}

func (r *readWrite) ReadNewsesByStatusAndTopicID(ctx context.Context, status int32, topicID string) (res *pb.Newses, err error) {
	stmt, err := r.db.Prepare(queryReadNewsesByStatusAndTopicID)
	if err != nil {
		return res, err
	}
	mutex.Lock()
	row, err := stmt.QueryContext(ctx, topicID, status)
	if err != nil {
		return res, err
	}
	mutex.Unlock()
	return r.rowsNewsesNextAndScan(ctx, row)
}

func (r *readWrite) ReadNewsesByTopicID(ctx context.Context, topicID string) (res *pb.Newses, err error) {
	stmt, err := r.db.Prepare(queryReadNewsesByTopicID)
	if err != nil {
		return res, err
	}
	mutex.Lock()
	row, err := stmt.QueryContext(ctx, topicID)
	if err != nil {
		return res, err
	}
	mutex.Unlock()
	return r.rowsNewsesNextAndScan(ctx, row)
}

func (r *readWrite) ReadNewsesByStatus(ctx context.Context, status int32) (res *pb.Newses, err error) {
	stmt, err := r.db.Prepare(queryReadNewsesByStatus)
	if err != nil {
		return res, err
	}
	mutex.Lock()
	row, err := stmt.QueryContext(ctx, status)
	if err != nil {
		return res, err
	}
	mutex.Unlock()
	return r.rowsNewsesNextAndScan(ctx, row)
}

func (r *readWrite) ReadNewses(ctx context.Context) (res *pb.Newses, err error) {
	stmt, err := r.db.Prepare(queryReadNewses)
	if err != nil {
		return res, err
	}
	mutex.Lock()
	row, err := stmt.QueryContext(ctx)
	if err != nil {
		return res, err
	}
	mutex.Unlock()
	return r.rowsNewsesNextAndScan(ctx, row)
}
