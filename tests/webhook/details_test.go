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

func TestWebhookDetails(t *testing.T) {
	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("invalid webhook id", func(t *testing.T) {
		e.GET(path).
			WithPath("webhook_id", "invalid_id").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("webhook record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(Columns))

		e.GET(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusNotFound)

		tests.ExpectationsMet(t, mock)
	})

	t.Run("return webhook with id", func(t *testing.T) {
		SelectMock(mock, 1)

		EventAssociationMock(mock)

		e.GET(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(respData)

		tests.ExpectationsMet(t, mock)
	})
}
