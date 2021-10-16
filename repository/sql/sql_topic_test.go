package sql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	"github.com/muhammadisa/bareksanews/util/mocker"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type sqlTopicTestSuite struct {
	suite.Suite
}

func TestTopicTestSuite(t *testing.T) {
	suite.Run(t, new(sqlTopicTestSuite))
}

func (ts *sqlTopicTestSuite) TestReadTopics() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	now := time.Now()

	// test case
	tests := []struct {
		Name      string
		Request   *pb.Topic
		WantError bool
	}{
		{
			Name: "read topic success",
			Request: &pb.Topic{
				Id:        uuid.NewV4().String(),
				Title:     "health",
				Headline:  "Talk about health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
			},
			WantError: false,
		},
	}

	repository := &readWrite{db: mockDB}
	errorDummy := errors.New("sql error while executing query")
	ctx := context.Background()
	defer ctx.Done()

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectPrepare(queryReadTopics)
				mock.ExpectQuery(queryReadTopics).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "headline", "created_at", "updated_at"}).
						AddRow(test.Request.Id, test.Request.Title, test.Request.Headline, now, now))

				topics, err := repository.ReadTopics(ctx)
				ts.Assert().NoError(err)
				ts.Assert().NotNil(topics)
				ts.Assert().NotNil(len(topics.Topics))

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryReadTopics).
					WillReturnError(errorDummy)
				mock.ExpectQuery(queryReadTopics).
					WillReturnError(errorDummy)

				topics, err := repository.ReadTopics(ctx)
				ts.Assert().NoError(err)
				ts.Assert().Nil(topics)
				ts.Assert().Len(len(topics.Topics), 0)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			}
		})
	}
}

func (ts *sqlTopicTestSuite) TestRemoveTopic() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	// test case
	tests := []struct {
		Name      string
		Request   *pb.Select
		WantError bool
	}{
		{
			Name: "delete topic success",
			Request: &pb.Select{
				Id: uuid.NewV4().String(),
			},
			WantError: false,
		},
		{
			Name: "delete topic failed",
			Request: &pb.Select{
				Id: uuid.NewV4().String(),
			},
			WantError: true,
		},
	}

	repository := &readWrite{db: mockDB}
	errorDummy := errors.New("sql error while executing query")
	ctx := context.Background()
	defer ctx.Done()

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectPrepare(queryRemoveTopic)
				mock.ExpectExec(queryRemoveTopic).
					WithArgs(test.Request.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.RemoveTopic(ctx, test.Request)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryRemoveTopic).
					WillReturnError(errorDummy)
				mock.ExpectExec(queryRemoveTopic).
					WithArgs(test.Request.Id).
					WillReturnError(errorDummy)

				err := repository.RemoveTopic(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlTopicTestSuite) TestModifyTopic() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	now := time.Now()

	// test case
	tests := []struct {
		Name      string
		Request   *pb.Topic
		WantError bool
	}{
		{
			Name: "modify topic success",
			Request: &pb.Topic{
				Id:       uuid.NewV4().String(),
				Title:    "health",
				Headline: "Talk about health",
			},
			WantError: false,
		},
		{
			Name: "modify topic failed",
			Request: &pb.Topic{
				Id:       uuid.NewV4().String(),
				Title:    "health",
				Headline: "Talk about health",
			},
			WantError: true,
		},
	}

	repository := &readWrite{db: mockDB}
	errorDummy := errors.New("sql error while executing query")
	currentDate := mocker.AnyTime{}
	ctx := context.Background()
	defer ctx.Done()

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectPrepare(queryLookupCreateAtTopic)
				mock.ExpectQuery(queryLookupCreateAtTopic).
					WithArgs(test.Request.Id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(test.Request.Id, now))
				mock.ExpectPrepare(queryUpdateTopic)
				mock.ExpectExec(queryUpdateTopic).
					WithArgs(test.Request.Title, test.Request.Headline, currentDate, currentDate, test.Request.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				updatedTopic, err := repository.ModifyTopic(ctx, test.Request)
				ts.Assert().NoError(err)
				ts.Assert().NotEqual(updatedTopic, &pb.Topic{})
				ts.Assert().Equal(updatedTopic.Id, test.Request.Id)
				ts.Assert().Equal(updatedTopic.Title, test.Request.Title)
				ts.Assert().Equal(updatedTopic.Headline, test.Request.Headline)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryLookupCreateAtTopic)
				mock.ExpectQuery(queryLookupCreateAtTopic).
					WillReturnError(errorDummy)
				mock.ExpectPrepare(queryUpdateTopic)
				mock.ExpectExec(queryUpdateTopic).
					WillReturnError(errorDummy)

				_, err := repository.ModifyTopic(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlTopicTestSuite) TestWriteTopic() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	// test case
	tests := []struct {
		Name      string
		Request   *pb.Topic
		WantError bool
	}{
		{
			Name: "write topic success",
			Request: &pb.Topic{
				Id:       uuid.NewV4().String(),
				Title:    "health",
				Headline: "Talk about health",
			},
			WantError: false,
		},
		{
			Name: "write topic failed",
			Request: &pb.Topic{
				Id:       uuid.NewV4().String(),
				Title:    "health",
				Headline: "Talk about health",
			},
			WantError: true,
		},
	}

	repository := &readWrite{db: mockDB}
	errorDummy := errors.New("sql error while executing query")
	currentDate := mocker.AnyTime{}
	ctx := context.Background()
	defer ctx.Done()

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectPrepare(queryWriteTopic)
				mock.ExpectExec(queryWriteTopic).
					WithArgs(test.Request.Id, test.Request.Title, test.Request.Headline, currentDate, currentDate).
					WillReturnResult(sqlmock.NewResult(1, 1))

				newTopic, err := repository.WriteTopic(ctx, test.Request)
				ts.Assert().NoError(err)
				ts.Assert().NotEqual(newTopic, &pb.Topic{})
				ts.Assert().Equal(newTopic.Id, test.Request.Id)
				ts.Assert().Equal(newTopic.Title, test.Request.Title)
				ts.Assert().Equal(newTopic.Headline, test.Request.Headline)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryWriteTopic).
					WillReturnError(errorDummy)
				mock.ExpectExec(queryWriteTopic).
					WithArgs(test.Request.Id, test.Request.Title, test.Request.Headline, currentDate, currentDate).
					WillReturnError(errorDummy)

				_, err := repository.WriteTopic(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}
