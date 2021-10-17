package service

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"

	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s service) AddNews(ctx context.Context, news *pb.News) (res *emptypb.Empty, err error) {
	news.Id = uuid.NewV4().String()
	newNews, err := s.repo.ReadWriter.WriteNews(ctx, news)
	if err != nil {
		return nil, err
	}
	err = s.repo.ReadWriter.WriteNewsTags(ctx, newNews.Id, news.NewsTagIds, true)
	if err != nil {
		_ = s.repo.ReadWriter.RemoveNews(ctx, &pb.Select{Id: newNews.Id})
		return nil, err
	}
	return nil, s.repo.CacheReadWriter.InvalidateNewses(ctx)
}

func (s service) EditNews(ctx context.Context, news *pb.News) (res *emptypb.Empty, err error) {
	_ = s.repo.CacheReadWriter.UnsetNews(ctx, news.Id)
	oldNews, err := s.repo.ReadWriter.ModifyNews(ctx, news)
	if err != nil {
		return nil, err
	}
	err = s.repo.ReadWriter.WriteNewsTags(ctx, oldNews.Id, news.NewsTagIds, false)
	if err != nil {
		return nil, err
	}
	return nil, s.repo.CacheReadWriter.InvalidateNewses(ctx)
}

func (s service) DeleteNews(ctx context.Context, selectNews *pb.Select) (res *emptypb.Empty, err error) {
	_ = s.repo.CacheReadWriter.UnsetNews(ctx, selectNews.Id)
	err = s.repo.ReadWriter.RemoveNews(ctx, selectNews)
	if err != nil {
		return nil, err
	}
	return nil, s.repo.CacheReadWriter.InvalidateNewses(ctx)
}

func (s service) filterRedisKeyGenerator(filters *pb.Filters) (res string) {
	if filters.TopicId != "" && filters.Status != 0 {
		res = fmt.Sprintf("topic_id_status_%s_%d", filters.TopicId, filters.Status)
	} else if filters.Status != 0 {
		res = fmt.Sprintf("status_%d", filters.Status)
	} else if filters.TopicId != "" {
		res = fmt.Sprintf("topic_id_%s", filters.TopicId)
	} else if filters.TopicId == "" && filters.Status == 0 {
		res = "none"
	}
	return
}

func (s service) GetNewses(ctx context.Context, filters *pb.Filters) (res *pb.Newses, err error) {
	filterValue := s.filterRedisKeyGenerator(filters)
	if reload := s.repo.CacheReadWriter.ReloadRequired(ctx, filterValue); reload {
		_ = s.repo.CacheReadWriter.SetFilter(ctx, filterValue)
		if filters.TopicId != "" && filters.Status != 0 {
			res, err = s.repo.ReadWriter.ReadNewsesByStatusAndTopicID(ctx, filters.Status, filters.TopicId)
		} else if filters.Status != 0 {
			res, err = s.repo.ReadWriter.ReadNewsesByStatus(ctx, filters.Status)
		} else if filters.TopicId != "" {
			res, err = s.repo.ReadWriter.ReadNewsesByTopicID(ctx, filters.TopicId)
		} else if filters.TopicId == "" && filters.Status == 0 {
			res, err = s.repo.ReadWriter.ReadNewses(ctx)
		}
		err = s.repo.CacheReadWriter.ReloadNewses(ctx, res)
		if err != nil {
			return nil, err
		}
		fmt.Println("database")
		return res, nil
	} else {
		fmt.Println("redis")
		return s.repo.CacheReadWriter.GetNewses(ctx)
	}
}
