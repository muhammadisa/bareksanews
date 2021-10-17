package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redismock/v8"
	"github.com/muhammadisa/bareksanews/constant"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"go.opencensus.io/trace"
	"testing"
	"time"
)

type cacheTopicTestSuite struct {
	suite.Suite
}

func TestCacheTopicTestSuite(t *testing.T) {
	suite.Run(t, new(cacheTopicTestSuite))
}

func (ts *cacheTopicTestSuite) TestReloadTopics() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	tests := []struct {
		Name      string
		Topics    *pb.Topics
		MapTopics map[string]string
		WantError bool
	}{
		{
			Name: "reload topics success",
			MapTopics: map[string]string{
				"c63b17cc-e227-4947-a01f-74f429ce99be": "{\"id\":\"c63b17cc-e227-4947-a01f-74f429ce99be\",\"title\":\"health\",\"headline\":\"this is headline\",\"created_at\":1634304927,\"updated_at\":1634304948}",
			},
			Topics: &pb.Topics{
				Topics: []*pb.Topic{
					{
						Id:        "c63b17cc-e227-4947-a01f-74f429ce99be",
						Title:     "health",
						Headline:  "this is headline",
						CreatedAt: 1634304927,
						UpdatedAt: 1634304948,
					},
				},
			},
			WantError: false,
		},
	}

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectHMSet(constant.Topics, test.MapTopics).SetVal(true)

				err := redisCache.ReloadTopics(ctx, test.Topics)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectHMSet(constant.Topics, test.MapTopics).RedisNil()

				err := redisCache.ReloadTopics(ctx, test.Topics)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *cacheTopicTestSuite) TestGetTopics() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	tests := []struct {
		Name      string
		Topics    map[string]string
		WantError bool
	}{
		{
			Name: "get topics success",
			Topics: map[string]string{
				"c63b17cc-e227-4947-a01f-74f429ce99be": "{\n  \"id\": \"c63b17cc-e227-4947-a01f-74f429ce99be\",\n  \"title\": \"tech\",\n \"headline\": \"this is headline\",\n  \"created_at\": 1634304927,\n  \"updated_at\": 1634304948\n}",
			},
			WantError: false,
		},
		{
			Name: "get topics failed",
			Topics: map[string]string{
				"c63b17cc-e227-4947-a01f-74f429ce99be": "{\n  \"id\": \"c63b17cc-e227-4947-a01f-74f429ce99be\",\n  \"title\": \"tech\",\n \"headline\": \"this is headline\",\n  \"created_at\": 1634304927,\n  \"updated_at\": 1634304948\n}",
			},
			WantError: false,
		},
	}

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectHGetAll(constant.Topics).SetVal(test.Topics)

				topicsData, err := redisCache.GetTopics(ctx)
				ts.Assert().NoError(err)

				for _, topic := range topicsData.Topics {
					ts.Assert().Equal(topic.Id, "c63b17cc-e227-4947-a01f-74f429ce99be")
					ts.Assert().Equal(topic.Title, "tech")
					ts.Assert().Equal(topic.Headline, "this is headline")
				}

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectHGetAll(constant.Topics).SetVal(test.Topics)

				topicsData, err := redisCache.GetTopics(ctx)
				ts.Assert().Error(err)
				ts.Assert().Nil(topicsData)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *cacheTopicTestSuite) TestUnsetTopic() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	tests := []struct {
		Name      string
		Id        string
		WantError bool
	}{
		{
			Name:      "unset topic success",
			Id:        uuid.NewV4().String(),
			WantError: false,
		},
		{
			Name:      "unset topic failed",
			Id:        uuid.NewV4().String(),
			WantError: true,
		},
	}

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectHDel(constant.Topics, test.Id).SetVal(1)

				err := redisCache.UnsetTopic(ctx, test.Id)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectHDel(constant.Topics, test.Id).SetErr(errors.New("dummy"))

				err := redisCache.UnsetTopic(ctx, test.Id)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			}
		})
	}
}

func (ts *cacheTopicTestSuite) TestSetTopic() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}
	now := time.Now()

	tests := []struct {
		Name         string
		Topic        *pb.Topic
		MarshalError bool
		WantError    bool
	}{
		{
			Name: "set topic success no marshal error",
			Topic: &pb.Topic{
				Id:        uuid.NewV4().String(),
				Title:     "health",
				Headline:  "this is headline",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
			},
			MarshalError: false,
			WantError:    false,
		},
		{
			Name: "set topic failed marshal error",
			Topic: &pb.Topic{
				Id:        uuid.NewV4().String(),
				Title:     "health",
				Headline:  "this is headline",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
			},
			MarshalError: true,
			WantError:    true,
		},
	}

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				topicByte, err := json.Marshal(test.Topic)
				ts.Assert().NoError(err)

				mock.ExpectHSetNX(constant.Topics, test.Topic.Id, string(topicByte)).SetVal(true)

				err = redisCache.SetTopic(ctx, test.Topic)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				topicByte, err := json.Marshal(test.Topic)
				ts.Assert().NoError(err)

				if !test.MarshalError {
					mock.ExpectHSetNX(constant.Topics, test.Topic.Id, string(topicByte)).SetVal(false)
				} else {
					mock.ExpectHSetNX(constant.Topics, test.Topic.Id, "").SetVal(false)
				}

				err = redisCache.SetTopic(ctx, test.Topic)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}
