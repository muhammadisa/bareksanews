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
	"testing"
	"time"
)

type cacheTagTestSuite struct {
	suite.Suite
}

func TestCacheTagTestSuite(t *testing.T) {
	suite.Run(t, new(cacheTagTestSuite))
}

func (ts *cacheTagTestSuite) TestReloadTags() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db}

	tests := []struct {
		Name      string
		Tags      *pb.Tags
		MapTags   map[string]string
		WantError bool
	}{
		{
			Name: "reload tags success",
			MapTags: map[string]string{
				"c63b17cc-e227-4947-a01f-74f429ce99be": "{\"id\":\"c63b17cc-e227-4947-a01f-74f429ce99be\",\"tag\":\"tech\",\"created_at\":1634304927,\"updated_at\":1634304948}",
			},
			Tags: &pb.Tags{
				Tags: []*pb.Tag{
					{
						Id:        "c63b17cc-e227-4947-a01f-74f429ce99be",
						Tag:       "tech",
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
				mock.ExpectHMSet(constant.Tags, test.MapTags).SetVal(true)

				err := redisCache.ReloadTags(ctx, test.Tags)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectHMSet(constant.Tags, test.MapTags).RedisNil()

				err := redisCache.ReloadTags(ctx, test.Tags)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *cacheTagTestSuite) TestGetTags() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db}

	tests := []struct {
		Name      string
		Tags      map[string]string
		WantError bool
	}{
		{
			Name: "get tags success",
			Tags: map[string]string{
				"c63b17cc-e227-4947-a01f-74f429ce99be": "{\n  \"id\": \"c63b17cc-e227-4947-a01f-74f429ce99be\",\n  \"tag\": \"tech\",\n  \"created_at\": 1634304927,\n  \"updated_at\": 1634304948\n}",
			},
			WantError: false,
		},
		{
			Name: "get tags failed",
			Tags: map[string]string{
				"c63b17cc-e227-4947-a01f-74f429ce99be": "{\n  \"id\": \"c63b17cc-e227-4947-a01f-74f429ce99be\",\n  \"tag\": \"tech\",\n  \"created_at\": 1634304927,\n  \"updated_at\": 1634304948\n}",
			},
			WantError: false,
		},
	}

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectHGetAll(constant.Tags).SetVal(test.Tags)

				tagsData, err := redisCache.GetTags(ctx)
				ts.Assert().NoError(err)

				for _, tag := range tagsData.Tags {
					ts.Assert().Equal(tag.Id, "c63b17cc-e227-4947-a01f-74f429ce99be")
					ts.Assert().Equal(tag.Tag, "tech")
				}

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectHGetAll(constant.Tags).SetVal(test.Tags)

				tagsData, err := redisCache.GetTags(ctx)
				ts.Assert().Error(err)
				ts.Assert().Nil(tagsData)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *cacheTagTestSuite) TestUnsetTag() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db}

	tests := []struct {
		Name      string
		Id        string
		WantError bool
	}{
		{
			Name:      "unset tag success",
			Id:        uuid.NewV4().String(),
			WantError: false,
		},
		{
			Name:      "unset tag failed",
			Id:        uuid.NewV4().String(),
			WantError: true,
		},
	}

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectHDel(constant.Tags, test.Id).SetVal(1)

				err := redisCache.UnsetTag(ctx, test.Id)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectHDel(constant.Tags, test.Id).SetErr(errors.New("dummy"))

				err := redisCache.UnsetTag(ctx, test.Id)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			}
		})
	}
}

func (ts *cacheTagTestSuite) TestSetTag() {
	db, mock := redismock.NewClientMock()
	ctx := context.Background()
	redisCache := &cache{redis: db}
	now := time.Now()

	tests := []struct {
		Name         string
		Tag          *pb.Tag
		MarshalError bool
		WantError    bool
	}{
		{
			Name: "set tag success no marshal error",
			Tag: &pb.Tag{
				Id:        uuid.NewV4().String(),
				Tag:       "health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
			},
			MarshalError: false,
			WantError:    false,
		},
		{
			Name: "set tag failed marshal error",
			Tag: &pb.Tag{
				Id:        uuid.NewV4().String(),
				Tag:       "health",
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
				tagByte, err := json.Marshal(test.Tag)
				ts.Assert().NoError(err)

				mock.ExpectHSetNX(constant.Tags, test.Tag.Id, string(tagByte)).SetVal(true)

				err = redisCache.SetTag(ctx, test.Tag)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				tagByte, err := json.Marshal(test.Tag)
				ts.Assert().NoError(err)

				if !test.MarshalError {
					mock.ExpectHSetNX(constant.Tags, test.Tag.Id, string(tagByte)).SetVal(false)
				} else {
					mock.ExpectHSetNX(constant.Tags, test.Tag.Id, "").SetVal(false)
				}

				err = redisCache.SetTag(ctx, test.Tag)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}
