package sql

import (
	"context"
	"errors"
	"go.opencensus.io/trace"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	"github.com/muhammadisa/bareksanews/util/mocker"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type sqlNewsTagTestSuite struct {
	suite.Suite
}

func TestNewsTagTestSuite(t *testing.T) {
	suite.Run(t, new(sqlNewsTagTestSuite))
}

func (ts *sqlNewsTagTestSuite) TestReadNewsTagsTagIDAndTagByNewsID() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	// test case
	tests := []struct {
		Name      string
		NewsID    string
		TagID     string
		Tag       string
		All       bool
		WantError bool
	}{
		{
			Name:      "read news tag success all true",
			NewsID:    uuid.NewV4().String(),
			TagID:     uuid.NewV4().String(),
			Tag:       "game",
			All:       true,
			WantError: false,
		},
		{
			Name:      "read news tag success all false",
			NewsID:    uuid.NewV4().String(),
			TagID:     uuid.NewV4().String(),
			Tag:       "game",
			All:       false,
			WantError: false,
		},
		{
			Name:      "read news tag failed all true",
			NewsID:    uuid.NewV4().String(),
			TagID:     uuid.NewV4().String(),
			Tag:       "game",
			All:       true,
			WantError: true,
		},
		{
			Name:      "read news tag failed all false",
			NewsID:    uuid.NewV4().String(),
			TagID:     uuid.NewV4().String(),
			Tag:       "game",
			All:       false,
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
				mock.ExpectPrepare(queryReadNewsTags)
				mock.ExpectQuery(queryReadNewsTags).
					WillReturnRows(sqlmock.NewRows([]string{"tag_id", "tag"}).
						AddRow(test.TagID, test.Tag))

				newsTags := repository.ReadNewsTagsTagIDAndTagByNewsID(ctx, test.NewsID, test.All)
				if test.All {
					ts.Assert().Equal(newsTags[0], test.Tag)
				} else {
					ts.Assert().Equal(newsTags[0], test.TagID)
				}

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryReadNewsTags).
					WillReturnError(errorDummy)
				mock.ExpectQuery(queryReadNewsTags).
					WillReturnError(errorDummy)

				newsTags := repository.ReadNewsTagsTagIDAndTagByNewsID(ctx, test.NewsID, test.All)
				ts.Assert().Nil(newsTags)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlNewsTagTestSuite) TestRemoveNewsTagsByNewsID() {
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
			Name: "delete news tag success",
			Request: &pb.Select{
				Id: uuid.NewV4().String(),
			},
			WantError: false,
		},
		{
			Name: "delete news tag failed",
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
				mock.ExpectPrepare(queryRemoveNewsTagsByNewsID)
				mock.ExpectExec(queryRemoveNewsTagsByNewsID).
					WithArgs(test.Request.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.RemoveNewsTagsByNewsID(ctx, test.Request)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryRemoveNewsTagsByNewsID).
					WillReturnError(errorDummy)
				mock.ExpectExec(queryRemoveNewsTagsByNewsID).
					WithArgs(test.Request.Id).
					WillReturnError(errorDummy)

				err := repository.RemoveNewsTagsByNewsID(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}
