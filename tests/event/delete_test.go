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
)

func TestEventDelete(t *testing.T) {
	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("invalid event id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("event_id", "invalid_id").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("event record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(Columns))

		e.DELETE(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusNotFound)

		tests.ExpectationsMet(t, mock)
	})

	t.Run("event is associated with some webhooks", func(t *testing.T) {
		SelectMock(mock, 1)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "webhooks"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		util.SubscribeExistingEvents = func() error { return nil }

		e.DELETE(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		tests.ExpectationsMet(t, mock)
	})

	t.Run("delete event", func(t *testing.T) {
		SelectMock(mock, 1)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "webhooks"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectExec(deleteQuery).
			WithArgs(tests.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		util.SubscribeExistingEvents = func() error { return nil }

		e.DELETE(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK)

		tests.ExpectationsMet(t, mock)
	})
}
