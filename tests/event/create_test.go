package event

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/hukz/service"
	"github.com/factly/hukz/tests"
	"github.com/factly/hukz/util"
	"github.com/gavv/httpexpect"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func TestEventCreate(t *testing.T) {

	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("Unprocessable event", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidData).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("Undecodable event", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("same name event already exist in db", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "events"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))

		e.POST(basePath).
			WithJSON(Data).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		tests.ExpectationsMet(t, mock)
	})

	t.Run("Invalid tags in body", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "events"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(0))

		Data["tags"] = postgres.Jsonb{
			RawMessage: []byte(`"test"`),
		}

		e.POST(basePath).
			WithJSON(Data).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		tests.ExpectationsMet(t, mock)
		Data["tags"] = postgres.Jsonb{
			RawMessage: []byte(`{"test":"tag"}`),
		}
	})

	t.Run("Add event in the database", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "events"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(0))

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "events"`).
			WithArgs(tests.AnyTime{}, tests.AnyTime{}, nil, 1, 1, Data["name"], Data["tags"]).
			WillReturnRows(sqlmock.
				NewRows([]string{"id"}).
				AddRow(1))
		mock.ExpectCommit()

		util.SubscribeEvents = func(events ...string) error { return nil }

		e.POST(basePath).
			WithJSON(Data).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusCreated)

		tests.ExpectationsMet(t, mock)
	})
}
