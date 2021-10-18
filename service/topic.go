package service

import (
	"context"
	"fmt"

	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s service) AddTopic(ctx context.Context, topic *pb.Topic) (*emptypb.Empty, error) {
	const funcName = `AddTopic`
	_, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	topic.Id = uuid.NewV4().String()
	newTopic, err := s.repo.ReadWriter.WriteTopic(ctx, topic)
	if err != nil {
		return nil, err
	}
	return nil, s.repo.CacheReadWriter.SetTopic(ctx, newTopic)
}

func (s service) EditTopic(ctx context.Context, topic *pb.Topic) (res *emptypb.Empty, err error) {
	const funcName = `EditTopic`
	_, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	err = s.repo.CacheReadWriter.UnsetTopic(ctx, topic.Id)
	if err != nil {
		return nil, err
	}
	updatedTopic, err := s.repo.ReadWriter.ModifyTopic(ctx, topic)
	if err != nil {
		return nil, err
	}
	return nil, s.repo.CacheReadWriter.SetTopic(ctx, updatedTopic)
}

func (s service) DeleteTopic(ctx context.Context, selectTopic *pb.Select) (res *emptypb.Empty, err error) {
	const funcName = `DeleteTopic`
	_, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	err = s.repo.ReadWriter.RemoveTopic(ctx, selectTopic)
	if err != nil {
		return nil, err
	}
	err = s.repo.CacheReadWriter.UnsetTopic(ctx, selectTopic.Id)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s service) GetTopics(ctx context.Context, _ *emptypb.Empty) (res *pb.Topics, err error) {
	const funcName = `GetTopics`
	_, span := s.tracer.StartSpan(ctx, funcName)
	defer span.End()

	res, err = s.repo.CacheReadWriter.GetTopics(ctx)
	if err != nil {
		return nil, nil
	}
	if len(res.Topics) == 0 {
		fmt.Println("from database")
		res, err = s.repo.ReadWriter.ReadTopics(ctx)
		if err != nil {
			return nil, err
		}
		_ = s.repo.CacheReadWriter.ReloadTopics(ctx, res)
		return res, nil
	}
	fmt.Println("from cache")
	return res, nil
}
