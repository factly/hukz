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

func TestEventUpdate(t *testing.T) {
	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("invalid event id", func(t *testing.T) {
		e.PUT(path).
			WithPath("event_id", "invalid_id").
			WithJSON(Data).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("event record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(Columns))

		e.PUT(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			WithJSON(Data).
			Expect().
			Status(http.StatusNotFound)
		tests.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable event body", func(t *testing.T) {
		e.PUT(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			WithJSON(invalidData).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable event body", func(t *testing.T) {
		e.PUT(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("event name exist in db", func(t *testing.T) {
		SelectMock(mock, 1)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "events"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).
				AddRow(1))

		Data["name"] = "event.new"
		e.PUT(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			WithJSON(Data).
			Expect().
			Status(http.StatusUnprocessableEntity)
		tests.ExpectationsMet(t, mock)
		Data["name"] = "event.done"
	})

	t.Run("tags invalid in event body", func(t *testing.T) {
		SelectMock(mock, 1)

		Data["tags"] = postgres.Jsonb{
			RawMessage: []byte(`"test"`),
		}
		e.PUT(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			WithJSON(Data).
			Expect().
			Status(http.StatusUnprocessableEntity)
		tests.ExpectationsMet(t, mock)
		Data["tags"] = postgres.Jsonb{
			RawMessage: []byte(`{"test":"tag"}`),
		}
	})

	util.SubscribeEvents = func(events ...string) error { return nil }
	t.Run("updated event", func(t *testing.T) {
		SelectMock(mock, 1)

		mock.ExpectExec(`UPDATE \"events\"`).
			WithArgs(tests.AnyTime{}, 1, Data["name"], Data["tags"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		e.PUT(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			WithJSON(Data).
			Expect().
			Status(http.StatusOK)
		tests.ExpectationsMet(t, mock)
	})
}
