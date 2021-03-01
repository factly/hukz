package event

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/hukz/service"
	"github.com/factly/hukz/tests"
	"github.com/gavv/httpexpect"
)

func TestEventDetails(t *testing.T) {
	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("invalid event id", func(t *testing.T) {
		e.GET(path).
			WithPath("event_id", "invalid_id").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("event record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(Columns))

		e.GET(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusNotFound)

		tests.ExpectationsMet(t, mock)
	})

	t.Run("get event by id", func(t *testing.T) {
		SelectMock(mock, 1)

		e.GET(path).
			WithPath("event_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(Data)

		tests.ExpectationsMet(t, mock)
	})
}
