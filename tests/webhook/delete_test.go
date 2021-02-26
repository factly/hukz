package webhook

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/hukz/service"
	"github.com/factly/hukz/tests"
	"github.com/gavv/httpexpect"
)

func TestWebhookDelete(t *testing.T) {
	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("invalid webhook id", func(t *testing.T) {
		e.DELETE(path).
			WithPath("webhook_id", "invalid_id").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("webhook record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(Columns))

		e.DELETE(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusNotFound)

		tests.ExpectationsMet(t, mock)
	})

	t.Run("delete webhook by id", func(t *testing.T) {
		SelectMock(mock, 1)

		mock.ExpectExec(deleteQuery).
			WithArgs(tests.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		e.DELETE(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK)

		tests.ExpectationsMet(t, mock)
	})
}
