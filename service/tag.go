package service

import (
	"context"
	"fmt"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s service) AddTag(ctx context.Context, tag *pb.Tag) (res *emptypb.Empty, err error) {
	const funcName = `AddTag`
	_, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	tag.Id = uuid.NewV4().String()
	newTag, err := s.repo.ReadWriter.WriteTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	return nil, s.repo.CacheReadWriter.SetTag(ctx, newTag)
}

func (s service) EditTag(ctx context.Context, tag *pb.Tag) (res *emptypb.Empty, err error) {
	const funcName = `EditTag`
	_, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	err = s.repo.CacheReadWriter.UnsetTag(ctx, tag.Id)
	if err != nil {
		return nil, err
	}
	updatedTag, err := s.repo.ReadWriter.ModifyTag(ctx, tag)
	if err != nil {
		return nil, err
	}
	return nil, s.repo.CacheReadWriter.SetTag(ctx, updatedTag)
}

func (s service) DeleteTag(ctx context.Context, selectTag *pb.Select) (res *emptypb.Empty, err error) {
	const funcName = `DeleteTag`
	_, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	err = s.repo.ReadWriter.RemoveTag(ctx, selectTag)
	if err != nil {
		return nil, err
	}
	err = s.repo.CacheReadWriter.UnsetTag(ctx, selectTag.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s service) GetTags(ctx context.Context, _ *emptypb.Empty) (res *pb.Tags, err error) {
	const funcName = `GetTags`
	_, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	res, err = s.repo.CacheReadWriter.GetTags(ctx)
	if err != nil {
		return nil, err
	}
	if len(res.Tags) == 0 {
		fmt.Println("from database")
		res, err = s.repo.ReadWriter.ReadTags(ctx)
		if err != nil {
			return nil, err
		}
		_ = s.repo.CacheReadWriter.ReloadTags(ctx, res)
		return res, nil
	}
	fmt.Println("from cache")
	return res, nil
}
