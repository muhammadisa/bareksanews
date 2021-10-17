package sql

import (
	"context"
	"errors"
	"go.opencensus.io/trace"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	"github.com/muhammadisa/bareksanews/util/mocker"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type sqlTagTestSuite struct {
	suite.Suite
}

func TestTagTestSuite(t *testing.T) {
	suite.Run(t, new(sqlTagTestSuite))
}

func (ts *sqlTagTestSuite) TestReadTags() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	now := time.Now()

	// test case
	tests := []struct {
		Name      string
		Request   *pb.Tag
		WantError bool
	}{
		{
			Name: "read tag success",
			Request: &pb.Tag{
				Id:        uuid.NewV4().String(),
				Tag:       "health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
			},
			WantError: false,
		},
	}

	repository := &readWrite{db: mockDB, tracer: trace.DefaultTracer}
	errorDummy := errors.New("sql error while executing query")
	ctx := context.Background()
	defer ctx.Done()

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectPrepare(queryReadTags)
				mock.ExpectQuery(queryReadTags).
					WillReturnRows(sqlmock.NewRows([]string{"id", "tag", "created_at", "updated_at"}).
						AddRow(test.Request.Id, test.Request.Tag, now, now))

				tags, err := repository.ReadTags(ctx)
				ts.Assert().NoError(err)
				ts.Assert().NotNil(tags)
				ts.Assert().NotNil(len(tags.Tags))

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryReadTags).
					WillReturnError(errorDummy)
				mock.ExpectQuery(queryReadTags).
					WillReturnError(errorDummy)

				tags, err := repository.ReadTags(ctx)
				ts.Assert().NoError(err)
				ts.Assert().Nil(tags)
				ts.Assert().Len(len(tags.Tags), 0)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			}
		})
	}
}

func (ts *sqlTagTestSuite) TestRemoveTag() {
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
			Name: "delete tag success",
			Request: &pb.Select{
				Id: uuid.NewV4().String(),
			},
			WantError: false,
		},
		{
			Name: "delete tag failed",
			Request: &pb.Select{
				Id: uuid.NewV4().String(),
			},
			WantError: true,
		},
	}

	repository := &readWrite{db: mockDB, tracer: trace.DefaultTracer}
	errorDummy := errors.New("sql error while executing query")
	ctx := context.Background()
	defer ctx.Done()

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectPrepare(queryRemoveTag)
				mock.ExpectExec(queryRemoveTag).
					WithArgs(test.Request.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.RemoveTag(ctx, test.Request)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryRemoveTag).
					WillReturnError(errorDummy)
				mock.ExpectExec(queryRemoveTag).
					WithArgs(test.Request.Id).
					WillReturnError(errorDummy)

				err := repository.RemoveTag(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlTagTestSuite) TestModifyTag() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	now := time.Now()

	// test case
	tests := []struct {
		Name      string
		Request   *pb.Tag
		WantError bool
	}{
		{
			Name: "modify tag success",
			Request: &pb.Tag{
				Id:  uuid.NewV4().String(),
				Tag: "health",
			},
			WantError: false,
		},
		{
			Name: "modify tag failed",
			Request: &pb.Tag{
				Id:  uuid.NewV4().String(),
				Tag: "health",
			},
			WantError: true,
		},
	}

	repository := &readWrite{db: mockDB, tracer: trace.DefaultTracer}
	errorDummy := errors.New("sql error while executing query")
	currentDate := mocker.AnyTime{}
	ctx := context.Background()
	defer ctx.Done()

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectPrepare(queryLookupCreateAtTag)
				mock.ExpectQuery(queryLookupCreateAtTag).
					WithArgs(test.Request.Id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(test.Request.Id, now))
				mock.ExpectPrepare(queryUpdateTag)
				mock.ExpectExec(queryUpdateTag).
					WithArgs(test.Request.Tag, currentDate, currentDate, test.Request.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				updatedTag, err := repository.ModifyTag(ctx, test.Request)
				ts.Assert().NoError(err)
				ts.Assert().NotEqual(updatedTag, &pb.Tag{})
				ts.Assert().Equal(updatedTag.Id, test.Request.Id)
				ts.Assert().Equal(updatedTag.Tag, test.Request.Tag)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryLookupCreateAtTag)
				mock.ExpectQuery(queryLookupCreateAtTag).
					WillReturnError(errorDummy)
				mock.ExpectPrepare(queryUpdateTag)
				mock.ExpectExec(queryUpdateTag).
					WillReturnError(errorDummy)

				_, err := repository.ModifyTag(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlTagTestSuite) TestWriteTag() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	// test case
	tests := []struct {
		Name      string
		Request   *pb.Tag
		WantError bool
	}{
		{
			Name: "write tag success",
			Request: &pb.Tag{
				Id:  uuid.NewV4().String(),
				Tag: "health",
			},
			WantError: false,
		},
		{
			Name: "write tag failed",
			Request: &pb.Tag{
				Id:  uuid.NewV4().String(),
				Tag: "health",
			},
			WantError: true,
		},
	}

	repository := &readWrite{db: mockDB, tracer: trace.DefaultTracer}
	errorDummy := errors.New("sql error while executing query")
	currentDate := mocker.AnyTime{}
	ctx := context.Background()
	defer ctx.Done()

	for _, test := range tests {
		ts.Run(test.Name, func() {
			if !test.WantError {
				mock.ExpectPrepare(queryWriteTag)
				mock.ExpectExec(queryWriteTag).
					WithArgs(test.Request.Id, test.Request.Tag, currentDate, currentDate).
					WillReturnResult(sqlmock.NewResult(1, 1))

				newTag, err := repository.WriteTag(ctx, test.Request)
				ts.Assert().NoError(err)
				ts.Assert().NotEqual(newTag, &pb.Tag{})
				ts.Assert().Equal(newTag.Id, test.Request.Id)
				ts.Assert().Equal(newTag.Tag, test.Request.Tag)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryWriteTag).
					WillReturnError(errorDummy)
				mock.ExpectExec(queryWriteTag).
					WithArgs(test.Request.Id, test.Request.Tag, currentDate, currentDate).
					WillReturnError(errorDummy)

				_, err := repository.WriteTag(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}
