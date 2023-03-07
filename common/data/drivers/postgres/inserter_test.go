// //go:build integration

package postgres

import (
	"context"
	"encoding/json"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"common/convert"
	"common/data/drivers"
	"common/data/model"
)

func createTestInserter[T model.Model](t *testing.T, log *logrus.Entry, db *sqlx.DB) *inserter[T] {
	t.Helper()

	var entity T

	return &inserter[T]{
		log: log.WithField("service", "[INSERTER-INTEGRATION-TEST]"),
		ext: db,
		sql: sq.Insert(entity.TableName()),
	}
}

func InsertUser(t *testing.T, log *logrus.Entry, db *sqlx.DB) {
	testInserter := createTestInserter[model.User](t, log, db)

	ctx := context.Background()

	expUser := &model.User{
		Username:  convert.ToPtr("username"),
		FirstName: convert.ToPtr("first_name"),
		LastName:  convert.ToPtr("last_name"),
	}

	testCases := []struct {
		name string
		in   model.User
		out  *model.User
		err  error
	}{
		{
			name: "ok user",
			in:   *expUser,
			out:  expUser,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := testInserter.Insert(ctx, tc.in)
			require.ErrorIs(t, err, tc.err)
			drivers.TestEqualWithoutFields(t, tc.out, u, "ID", "CreatedAt")
			require.NotEmpty(t, u.ID)
		})
	}
}

func InsertNews(t *testing.T, log *logrus.Entry, db *sqlx.DB) {
	testInserter := createTestInserter[model.News](t, log, db)

	ctx := context.Background()

	expNews := &model.News{
		Media: &model.NewsMedia{
			Title: convert.ToPtr("title"),
			Text:  convert.ToPtr("text"),
			Resources: []model.NewsMediaResource{
				{
					Type: convert.ToPtr("type"),
					URL:  convert.ToPtr("url"),
					Meta: json.RawMessage(`{}`),
				},
			},
		},
		Source: convert.ToPtr("source"),
	}

	testCases := []struct {
		name string
		in   model.News
		out  *model.News
		err  error
	}{
		{
			name: "ok news",
			in:   *expNews,
			out:  expNews,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := testInserter.Insert(ctx, tc.in)
			require.ErrorIs(t, err, tc.err)
			drivers.TestEqualWithoutFields(t, tc.out, u, "ID", "CreatedAt")
			require.NotEmpty(t, u.ID)
		})
	}
}

func InsertUsersBatch(t *testing.T, log *logrus.Entry, db *sqlx.DB) {
	testInserter := createTestInserter[model.User](t, log, db)

	ctx := context.Background()

	expUsers := []model.User{
		{
			Username:  convert.ToPtr("username1"),
			FirstName: convert.ToPtr("first_name1"),
			LastName:  convert.ToPtr("last_name1"),
		},
		{
			Username:  convert.ToPtr("username2"),
			FirstName: convert.ToPtr("first_name2"),
			LastName:  convert.ToPtr("last_name2"),
		},
		{
			Username:  convert.ToPtr("username3"),
			FirstName: convert.ToPtr("first_name3"),
			LastName:  convert.ToPtr("last_name3"),
		},
	}

	testCases := []struct {
		name string
		in   []model.User
		out  []model.User
		err  error
	}{
		{
			name: "ok user batch",
			in:   expUsers,
			out:  expUsers,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testInserter.InsertBatch(ctx, tc.in)
			require.ErrorIs(t, err, tc.err)
			for i, u := range tc.in {
				drivers.TestEqualWithoutFields(t, tc.out[i], tc.in[i], "ID", "CreatedAt")
				require.NotEmpty(t, u.ID)
			}

		})
	}
}

func InsertNewsBatch(t *testing.T, log *logrus.Entry, db *sqlx.DB) {
	testInserter := createTestInserter[model.News](t, log, db)

	ctx := context.Background()

	expNews := []model.News{
		{
			Media: &model.NewsMedia{
				Title: convert.ToPtr("title1"),
				Text:  convert.ToPtr("text1"),
				Resources: []model.NewsMediaResource{
					{
						Type: convert.ToPtr("type1"),
						URL:  convert.ToPtr("url1"),
						Meta: json.RawMessage(`{}`),
					},
				},
			},
			Source: convert.ToPtr("source1"),
		},
		{
			Media: &model.NewsMedia{
				Title: convert.ToPtr("title2"),
				Text:  convert.ToPtr("text2"),
				Resources: []model.NewsMediaResource{
					{
						Type: convert.ToPtr("type2"),
						URL:  convert.ToPtr("url2"),
						Meta: json.RawMessage(`{}`),
					},
				},
			},
			Source: convert.ToPtr("source2"),
		},
		{
			Media: &model.NewsMedia{
				Title: convert.ToPtr("title3"),
				Text:  convert.ToPtr("text3"),
				Resources: []model.NewsMediaResource{
					{
						Type: convert.ToPtr("type3"),
						URL:  convert.ToPtr("url3"),
						Meta: json.RawMessage(`{}`),
					},
				},
			},
			Source: convert.ToPtr("source3"),
		},
	}

	testCases := []struct {
		name string
		in   []model.News
		out  []model.News
		err  error
	}{
		{
			name: "ok news batch",
			in:   expNews,
			out:  expNews,
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testInserter.InsertBatch(ctx, tc.in)
			require.ErrorIs(t, err, tc.err)
			for i, u := range tc.in {
				drivers.TestEqualWithoutFields(t, tc.out[i], tc.in[i], "ID", "CreatedAt")
				require.NotEmpty(t, u.ID)
			}

		})
	}
}

func TestInserter(t *testing.T) {
	suite := drivers.NewSuite(t)
	suite.AddTests(InsertUser, InsertNews, InsertUsersBatch, InsertNewsBatch)

	suite.SetupSuite()
	defer suite.CleanupSuite()

	suite.TestRunIntegration()
}
