package sql

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	uuid "github.com/satori/go.uuid"
)

func (r *readWrite) RemoveNewsTagsByNewsID(ctx context.Context, req *pb.Select) error {
	const funcName = `RemoveNewsTagsByNewsID`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	stmt, err := r.db.Prepare(queryRemoveNewsTagsByNewsID)
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
		return fmt.Errorf("failed to delete reason : %+v", err)
	}
	return nil
}

func (r *readWrite) ReadNewsTagsTagIDAndTagByNewsID(ctx context.Context, newsID string, all bool) (res []string) {
	const funcName = `ReadNewsTagsTagIDAndTagByNewsID`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	var tagID, tag string
	stmt, err := r.db.Prepare(queryReadNewsTags)
	if err != nil {
		return nil
	}
	mutex.Lock()
	row, err := stmt.QueryContext(ctx, newsID)
	if err != nil {
		return nil
	}
	mutex.Unlock()
	for row.Next() {
		err = row.Scan(
			&tagID,
			&tag,
		)
		if err != nil {
			return nil
		}
		if all {
			res = append(res, tag)
		} else {
			res = append(res, tagID)
		}
	}
	return res
}

func (r *readWrite) WriteNewsTags(ctx context.Context, newsID string, tagIDs []string, new bool) error {
	const funcName = `WriteNewsTags`
	_, span := r.tracer.StartSpan(ctx, funcName)
	defer span.End()

	timeNow := time.Now()
	length := len(tagIDs)

	if length < 0 {
		return nil
	}

	if !new {
		err := r.RemoveNewsTagsByNewsID(ctx, &pb.Select{Id: newsID})
		if err != nil {
			return err
		}
	}

	var valueStrings []string
	var valueArgs []interface{}
	for _, tagID := range tagIDs {
		id := uuid.NewV4().String()
		valueStrings = append(valueStrings, "(?,?,?,?,?)")
		valueArgs = append(valueArgs, id)      // id
		valueArgs = append(valueArgs, newsID)  // news_id
		valueArgs = append(valueArgs, tagID)   // tag_id
		valueArgs = append(valueArgs, timeNow) // created_at
		valueArgs = append(valueArgs, timeNow) // updated_at
	}

	query := fmt.Sprintf(queryWriteBulkNewsTags, strings.Join(valueStrings, ","))
	result, err := r.db.Exec(query, valueArgs...)
	if err != nil {
		return err
	}
	if affected, err := result.RowsAffected(); affected != int64(length) || err != nil {
		return fmt.Errorf("failed to insert reason : %+v", err)
	}
	return nil
}
