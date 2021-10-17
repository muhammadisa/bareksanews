package cache

import (
	"context"
	"github.com/go-redis/redismock/v8"
	"github.com/muhammadisa/bareksanews/constant"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	"github.com/stretchr/testify/suite"
	"go.opencensus.io/trace"
	"testing"
)

type cacheNewsTestSuite struct {
	suite.Suite
}

func TestCacheNewsTestSuite(t *testing.T) {
	suite.Run(t, new(cacheNewsTestSuite))
}

func (ts *cacheNewsTestSuite) TestGetFilter() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	mock.ExpectGet(constant.FilterNewses).SetVal("ok")

	data := redisCache.GetFilter(ctx)
	ts.Assert().Equal("ok", data)

	err := mock.ExpectationsWereMet()
	ts.Assert().NoError(err)
}

func (ts *cacheNewsTestSuite) TestSetFilter() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	mock.ExpectSet(constant.FilterNewses, "filter_1", 0).SetVal("filter_1")

	err := redisCache.SetFilter(ctx, "filter_1")
	ts.Assert().NoError(err)

	err = mock.ExpectationsWereMet()
	ts.Assert().NoError(err)
}

func (ts *cacheNewsTestSuite) TestUnsetNews() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	mock.ExpectHDel(constant.News, "field_1").SetVal(1)

	err := redisCache.UnsetNews(ctx, "field_1")
	ts.Assert().NoError(err)

	err = mock.ExpectationsWereMet()
	ts.Assert().NoError(err)
}

func (ts *cacheNewsTestSuite) TestInvalidateNewses() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	mock.ExpectSet(constant.ReloadNewses, true, 0).SetVal("true")

	err := redisCache.InvalidateNewses(ctx)
	ts.Assert().NoError(err)

	err = mock.ExpectationsWereMet()
	ts.Assert().NoError(err)
}

func (ts *cacheNewsTestSuite) TestReloadRequired() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	mock.ExpectGet(constant.FilterNewses).SetVal("filter_1")
	mock.ExpectGet(constant.ReloadNewses).SetVal("0")

	reload := redisCache.ReloadRequired(ctx, "filter_1_not_same")
	ts.Assert().True(reload)

	err := mock.ExpectationsWereMet()
	ts.Assert().NoError(err)
}

func (ts *cacheTagTestSuite) TestGetNewses() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	tests := []struct {
		Name      string
		Newses    map[string]string
		WantError bool
	}{
		{
			Name: "get newses success",
			Newses: map[string]string{
				"9366c83d-4c1e-40ab-93ca-30b9548aebf7": "{\n  \"id\": \"9366c83d-4c1e-40ab-93ca-30b9548aebf7\",\n  \"topic_id\": \"d95cb090-0906-471a-80ef-3714c6451920\",\n  \"title\": \"title news number 1\",\n  \"content\": \"content news 1\",\n  \"news_tag_ids\": [\n    \"0f4e1e74-9238-4afb-87c2-108e569ff866\",\n    \"4e358682-e2b7-4ecf-9e1f-4373bffd661a\"\n  ],\n  \"news_tag_names\": [\n    \"health\",\n    \"game\"\n  ],\n  \"status\": 1,\n  \"created_at\": 1634323641,\n  \"updated_at\": 1634338479\n}",
			},
			WantError: false,
		},
		{
			Name: "get newses failed",
			Newses: map[string]string{
				"9366c83d-4c1e-40ab-93ca-30b9548aebf7": "{\n  \"id\": \"9366c83d-4c1e-40ab-93ca-30b9548aebf7\",\n  \"topic_id\": \"d95cb090-0906-471a-80ef-3714c6451920\",\n  \"title\": \"title news number 1\",\n  \"content\": \"content news 1\",\n  \"news_tag_ids\": [\n    \"0f4e1e74-9238-4afb-87c2-108e569ff866\",\n    \"4e358682-e2b7-4ecf-9e1f-4373bffd661a\"\n  ],\n  \"news_tag_names\": [\n    \"health\",\n    \"game\"\n  ],\n  \"status\": 1,\n  \"created_at\": 1634323641,\n  \"updated_at\": 1634338479\n}",
			},
			WantError: false,
		},
	}

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectHGetAll(constant.News).SetVal(test.Newses)

				newsesData, err := redisCache.GetNewses(ctx)
				ts.Assert().NoError(err)

				for _, news := range newsesData.Newses {
					ts.Assert().Equal(news.Id, "9366c83d-4c1e-40ab-93ca-30b9548aebf7")
					ts.Assert().Equal(news.TopicId, "d95cb090-0906-471a-80ef-3714c6451920")
					ts.Assert().Equal(news.Title, "title news number 1")
					ts.Assert().Equal(news.Content, "content news 1")
					ts.Assert().Equal(news.Status, int32(1))
					ts.Assert().Equal(news.CreatedAt, int64(1634323641))
					ts.Assert().Equal(news.UpdatedAt, int64(1634338479))
				}

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectHGetAll(constant.News).SetVal(test.Newses)

				tagsData, err := redisCache.GetNewses(ctx)
				ts.Assert().Error(err)
				ts.Assert().Nil(tagsData)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *cacheTagTestSuite) TestReloadNewses() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db, tracer: trace.DefaultTracer}

	tests := []struct {
		Name      string
		Newses    *pb.Newses
		MapNewses map[string]string
		WantError bool
	}{
		{
			Name: "reload newses success",
			MapNewses: map[string]string{
				"9366c83d-4c1e-40ab-93ca-30b9548aebf7": "{\"id\":\"9366c83d-4c1e-40ab-93ca-30b9548aebf7\",\"topic_id\":\"d95cb090-0906-471a-80ef-3714c6451920\",\"title\":\"title news number 1\",\"content\":\"content news 1\",\"news_tag_ids\":[\"0f4e1e74-9238-4afb-87c2-108e569ff866\",\"4e358682-e2b7-4ecf-9e1f-4373bffd661a\"],\"news_tag_names\":[\"health\",\"game\"],\"status\":1,\"created_at\":1634323641,\"updated_at\":1634338479}",
			},
			Newses: &pb.Newses{
				Newses: []*pb.News{
					{
						Id:           "9366c83d-4c1e-40ab-93ca-30b9548aebf7",
						TopicId:      "d95cb090-0906-471a-80ef-3714c6451920",
						Title:        "title news number 1",
						Content:      "content news 1",
						Status:       1,
						NewsTagIds:   []string{"0f4e1e74-9238-4afb-87c2-108e569ff866", "4e358682-e2b7-4ecf-9e1f-4373bffd661a"},
						NewsTagNames: []string{"health", "game"},
						CreatedAt:    1634323641,
						UpdatedAt:    1634338479,
					},
				},
			},
			WantError: false,
		},
	}

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectHMSet(constant.News, test.MapNewses).SetVal(true)

				err := redisCache.ReloadNewses(ctx, test.Newses)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectHMSet(constant.News, test.MapNewses).RedisNil()

				err := redisCache.ReloadNewses(ctx, test.Newses)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}
