package webhook

import (
	"database/sql/driver"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/hukz/service"
	"github.com/factly/hukz/tests"
	"github.com/factly/hukz/tests/event"
	"github.com/gavv/httpexpect"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func TestWebhookUpdate(t *testing.T) {
	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("invalid webhook id", func(t *testing.T) {
		e.PUT(path).
			WithPath("webhook_id", "invalid_id").
			WithJSON(Data).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusBadRequest)
	})

	t.Run("webhook record not found", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows(Columns))

		e.PUT(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			WithJSON(Data).
			Expect().
			Status(http.StatusNotFound)
		tests.ExpectationsMet(t, mock)
	})

	t.Run("unprocessable webhook body", func(t *testing.T) {
		e.PUT(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			WithJSON(invalidData).
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("undecodable webhook body", func(t *testing.T) {
		e.PUT(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("replacing events fails", func(t *testing.T) {
		SelectMock(mock, 1)

		mock.ExpectBegin()
		event.SelectMock(mock)

		mock.ExpectExec(`UPDATE \"webhooks\"`).
			WithArgs(tests.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "events"`).
			WithArgs(tests.AnyTime{}, tests.AnyTime{}, nil, 1, 1, "event.done", postgres.Jsonb{RawMessage: []byte(`{"test":"tag"}`)}, 1).
			WillReturnRows(sqlmock.
				NewRows([]string{"id"}).
				AddRow(1))

		mock.ExpectExec(`INSERT INTO "webhook_events"`).
			WithArgs(1, 1).
			WillReturnError(errors.New(`cannot replace association`))

		mock.ExpectRollback()

		e.PUT(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			WithJSON(Data).
			Expect().
			Status(http.StatusInternalServerError)
		tests.ExpectationsMet(t, mock)
	})

	t.Run("invalid tags in request body", func(t *testing.T) {
		SelectMock(mock, 1)

		mock.ExpectBegin()

		mock.ExpectRollback()

		Data["tags"] = postgres.Jsonb{
			RawMessage: []byte(`"test"`),
		}
		e.PUT(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			WithJSON(Data).
			Expect().
			Status(http.StatusUnprocessableEntity)
		tests.ExpectationsMet(t, mock)
		Data["tags"] = postgres.Jsonb{
			RawMessage: []byte(`{"test":"tag"}`),
		}
	})

	t.Run("update webhook", func(t *testing.T) {
		SelectMock(mock, 1)

		mock.ExpectBegin()
		event.SelectMock(mock)

		mock.ExpectExec(`UPDATE \"webhooks\"`).
			WithArgs(tests.AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(`INSERT INTO "events"`).
			WithArgs(tests.AnyTime{}, tests.AnyTime{}, nil, 1, 1, "event.done", postgres.Jsonb{RawMessage: []byte(`{"test":"tag"}`)}, 1).
			WillReturnRows(sqlmock.
				NewRows([]string{"id"}).
				AddRow(1))

		mock.ExpectExec(`INSERT INTO "webhook_events"`).
			WithArgs(1, 1).
			WillReturnResult(driver.ResultNoRows)

		mock.ExpectExec(`DELETE FROM "webhook_events"`).
			WithArgs(1, 1).
			WillReturnResult(driver.ResultNoRows)

		mock.ExpectExec(`UPDATE \"webhooks\"`).
			WithArgs(tests.AnyTime{}, true, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(`UPDATE \"webhooks\"`).
			WithArgs(tests.AnyTime{}, 1, Data["name"], Data["url"], Data["tags"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		SelectMock(mock)
		EventAssociationMock(mock)

		mock.ExpectCommit()

		e.PUT(path).
			WithPath("webhook_id", "1").
			WithHeader("X-User", "1").
			WithJSON(Data).
			Expect().
			Status(http.StatusOK)
		tests.ExpectationsMet(t, mock)
	})
}
