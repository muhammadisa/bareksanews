package sql

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/muhammadisa/bareksanews/protoc/api/v1"
	"github.com/muhammadisa/bareksanews/util/mocker"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type sqlNewsTestSuite struct {
	suite.Suite
}

func TestNewsTestSuite(t *testing.T) {
	suite.Run(t, new(sqlNewsTestSuite))
}

func (ts *sqlNewsTestSuite) TestReadNewsesByStatusAndTopicID() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	now := time.Now()

	// test case
	tests := []struct {
		Name      string
		Request   *pb.News
		WantError bool
	}{
		{
			Name: "read news success",
			Request: &pb.News{
				Id:        uuid.NewV4().String(),
				TopicId: uuid.NewV4().String(),
				Title:     "health",
				Content:  "Talk about health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
			},
			WantError: false,
		},
		{
			Name: "read news success",
			Request: &pb.News{
				Id:        uuid.NewV4().String(),
				TopicId: uuid.NewV4().String(),
				Title:     "health",
				Content:  "Talk about health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
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
				mock.ExpectPrepare(queryReadNewsesByStatusAndTopicID)
				mock.ExpectQuery(queryReadNewsesByStatusAndTopicID).
					WithArgs(test.Request.TopicId, test.Request.Status).
					WillReturnRows(sqlmock.NewRows([]string{"id", "topic_id", "title", "content", "status", "created_at", "updated_at"}).
						AddRow(test.Request.Id, test.Request.TopicId, test.Request.Title, test.Request.Content, test.Request.Status, now, now))

				newses, err := repository.ReadNewsesByStatusAndTopicID(ctx, test.Request.Status, test.Request.TopicId)
				ts.Assert().NoError(err)
				ts.Assert().NotNil(newses)
				ts.Assert().NotNil(len(newses.Newses))

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryReadNewsesByStatusAndTopicID).
					WillReturnError(errorDummy)
				mock.ExpectQuery(queryReadNewsesByStatusAndTopicID).
					WithArgs(test.Request.TopicId, test.Request.Status).
					WillReturnError(errorDummy)

				newses, err := repository.ReadNewsesByStatusAndTopicID(ctx, test.Request.Status, test.Request.TopicId)
				ts.Assert().Error(err)
				ts.Assert().Nil(newses)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlNewsTestSuite) TestReadNewsesByTopicID() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	now := time.Now()

	// test case
	tests := []struct {
		Name      string
		Request   *pb.News
		WantError bool
	}{
		{
			Name: "read news success",
			Request: &pb.News{
				Id:        uuid.NewV4().String(),
				TopicId: uuid.NewV4().String(),
				Title:     "health",
				Content:  "Talk about health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
			},
			WantError: false,
		},
		{
			Name: "read news success",
			Request: &pb.News{
				Id:        uuid.NewV4().String(),
				TopicId: uuid.NewV4().String(),
				Title:     "health",
				Content:  "Talk about health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
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
				mock.ExpectPrepare(queryReadNewsesByTopicID)
				mock.ExpectQuery(queryReadNewsesByTopicID).
					WithArgs(test.Request.TopicId).
					WillReturnRows(sqlmock.NewRows([]string{"id", "topic_id", "title", "content", "status", "created_at", "updated_at"}).
						AddRow(test.Request.Id, test.Request.TopicId, test.Request.Title, test.Request.Content, test.Request.Status, now, now))

				newses, err := repository.ReadNewsesByTopicID(ctx, test.Request.TopicId)
				ts.Assert().NoError(err)
				ts.Assert().NotNil(newses)
				ts.Assert().NotNil(len(newses.Newses))

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryReadNewsesByTopicID).
					WillReturnError(errorDummy)
				mock.ExpectQuery(queryReadNewsesByTopicID).
					WithArgs(test.Request.TopicId).
					WillReturnError(errorDummy)

				newses, err := repository.ReadNewsesByTopicID(ctx, test.Request.TopicId)
				ts.Assert().Error(err)
				ts.Assert().Nil(newses)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlNewsTestSuite) TestReadNewsesByStatus() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	now := time.Now()

	// test case
	tests := []struct {
		Name      string
		Request   *pb.News
		WantError bool
	}{
		{
			Name: "read news success",
			Request: &pb.News{
				Id:        uuid.NewV4().String(),
				TopicId: uuid.NewV4().String(),
				Title:     "health",
				Content:  "Talk about health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
			},
			WantError: false,
		},
		{
			Name: "read news success",
			Request: &pb.News{
				Id:        uuid.NewV4().String(),
				TopicId: uuid.NewV4().String(),
				Title:     "health",
				Content:  "Talk about health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
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
				mock.ExpectPrepare(queryReadNewsesByStatus)
				mock.ExpectQuery(queryReadNewsesByStatus).
					WithArgs(test.Request.Status).
					WillReturnRows(sqlmock.NewRows([]string{"id", "topic_id", "title", "content", "status", "created_at", "updated_at"}).
						AddRow(test.Request.Id, test.Request.TopicId, test.Request.Title, test.Request.Content, test.Request.Status, now, now))

				newses, err := repository.ReadNewsesByStatus(ctx, test.Request.Status)
				ts.Assert().NoError(err)
				ts.Assert().NotNil(newses)
				ts.Assert().NotNil(len(newses.Newses))

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryReadNewsesByStatus).
					WillReturnError(errorDummy)
				mock.ExpectQuery(queryReadNewsesByStatus).
					WithArgs(test.Request.Status).
					WillReturnError(errorDummy)

				newses, err := repository.ReadNewsesByStatus(ctx, test.Request.Status)
				ts.Assert().Error(err)
				ts.Assert().Nil(newses)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlNewsTestSuite) TestReadNewses() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	now := time.Now()

	// test case
	tests := []struct {
		Name      string
		Request   *pb.News
		WantError bool
	}{
		{
			Name: "read news success",
			Request: &pb.News{
				Id:        uuid.NewV4().String(),
				TopicId: uuid.NewV4().String(),
				Title:     "health",
				Content:  "Talk about health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
			},
			WantError: false,
		},
		{
			Name: "read news success",
			Request: &pb.News{
				Id:        uuid.NewV4().String(),
				TopicId: uuid.NewV4().String(),
				Title:     "health",
				Content:  "Talk about health",
				CreatedAt: now.Unix(),
				UpdatedAt: now.Unix(),
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
				mock.ExpectPrepare(queryReadNewses)
				mock.ExpectQuery(queryReadNewses).
					WillReturnRows(sqlmock.NewRows([]string{"id", "topic_id", "title", "content", "status", "created_at", "updated_at"}).
						AddRow(test.Request.Id, test.Request.TopicId, test.Request.Title, test.Request.Content, test.Request.Status, now, now))

				newses, err := repository.ReadNewses(ctx)
				ts.Assert().NoError(err)
				ts.Assert().NotNil(newses)
				ts.Assert().NotNil(len(newses.Newses))

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryReadNewses).
					WillReturnError(errorDummy)
				mock.ExpectQuery(queryReadNewses).
					WillReturnError(errorDummy)

				newses, err := repository.ReadNewses(ctx)
				ts.Assert().Error(err)
				ts.Assert().Nil(newses)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlNewsTestSuite) TestRemoveNews() {
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
			Name: "delete news success",
			Request: &pb.Select{
				Id: uuid.NewV4().String(),
			},
			WantError: false,
		},
		{
			Name: "delete news failed",
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
				mock.ExpectPrepare(queryRemoveNews)
				mock.ExpectExec(queryRemoveNews).
					WithArgs(test.Request.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				err := repository.RemoveNews(ctx, test.Request)
				ts.Assert().NoError(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryRemoveNews).
					WillReturnError(errorDummy)
				mock.ExpectExec(queryRemoveNews).
					WithArgs(test.Request.Id).
					WillReturnError(errorDummy)

				err := repository.RemoveNews(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlNewsTestSuite) TestModifyNews() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	now := time.Now()

	// test case
	tests := []struct {
		Name      string
		Request   *pb.News
		WantError bool
	}{
		{
			Name: "modify news success",
			Request: &pb.News{
				Id:       uuid.NewV4().String(),
				Title:    "health",
				Content: "Talk about health",
			},
			WantError: false,
		},
		{
			Name: "modify news failed",
			Request: &pb.News{
				Id:       uuid.NewV4().String(),
				Title:    "health",
				Content: "Talk about health",
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
				mock.ExpectPrepare(queryLookupCreateAtNews)
				mock.ExpectQuery(queryLookupCreateAtNews).
					WithArgs(test.Request.Id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(test.Request.Id, now))
				mock.ExpectPrepare(queryUpdateNews)
				mock.ExpectExec(queryUpdateNews).
					WithArgs(test.Request.TopicId, test.Request.Title, test.Request.Content, test.Request.Status, currentDate, currentDate, test.Request.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				updatedNews, err := repository.ModifyNews(ctx, test.Request)
				ts.Assert().NoError(err)
				ts.Assert().NotEqual(updatedNews, &pb.News{})
				ts.Assert().Equal(updatedNews.Id, test.Request.Id)
				ts.Assert().Equal(updatedNews.TopicId, test.Request.TopicId)
				ts.Assert().Equal(updatedNews.Title, test.Request.Title)
				ts.Assert().Equal(updatedNews.Content, test.Request.Content)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryLookupCreateAtNews)
				mock.ExpectQuery(queryLookupCreateAtNews).
					WillReturnError(errorDummy)
				mock.ExpectPrepare(queryUpdateNews)
				mock.ExpectExec(queryUpdateNews).
					WillReturnError(errorDummy)

				_, err := repository.ModifyNews(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}

func (ts *sqlNewsTestSuite) TestWriteNews() {
	// sql mock
	mockDB, mock, err := mocker.SQLMocker()
	ts.Require().NoError(err)
	ts.Require().NotNil(mockDB)
	ts.Require().NotNil(mock)

	// test case
	tests := []struct {
		Name      string
		Request   *pb.News
		WantError bool
	}{
		{
			Name: "write news success",
			Request: &pb.News{
				Id:       uuid.NewV4().String(),
				Title:    "health",
				Content: "Talk about health",
			},
			WantError: false,
		},
		{
			Name: "write news failed",
			Request: &pb.News{
				Id:       uuid.NewV4().String(),
				Title:    "health",
				Content: "Talk about health",
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
				mock.ExpectPrepare(queryWriteNews)
				mock.ExpectExec(queryWriteNews).
					WithArgs(test.Request.Id, test.Request.TopicId, test.Request.Title, test.Request.Content, test.Request.Status, currentDate, currentDate).
					WillReturnResult(sqlmock.NewResult(1, 1))

				newNews, err := repository.WriteNews(ctx, test.Request)
				ts.Assert().NoError(err)
				ts.Assert().NotEqual(newNews, &pb.News{})
				ts.Assert().Equal(newNews.Id, test.Request.Id)
				ts.Assert().Equal(newNews.TopicId, test.Request.TopicId)
				ts.Assert().Equal(newNews.Title, test.Request.Title)
				ts.Assert().Equal(newNews.Content, test.Request.Content)

				err = mock.ExpectationsWereMet()
				ts.Assert().NoError(err)
			} else {
				mock.ExpectPrepare(queryWriteNews).
					WillReturnError(errorDummy)
				mock.ExpectExec(queryWriteNews).
					WithArgs(test.Request.Id, test.Request.TopicId, test.Request.Title, test.Request.Content, test.Request.Status, currentDate, currentDate).
					WillReturnError(errorDummy)

				_, err := repository.WriteNews(ctx, test.Request)
				ts.Assert().Error(err)

				err = mock.ExpectationsWereMet()
				ts.Assert().Error(err)
			}
		})
	}
}
