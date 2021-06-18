package webhook

import (
	"database/sql/driver"
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

func TestWebhookCreate(t *testing.T) {

	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("Unprocessable webhook", func(t *testing.T) {
		e.POST(basePath).
			WithJSON(invalidData).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("Undecodable webhook", func(t *testing.T) {
		e.POST(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)
	})

	t.Run("Invalid tags", func(t *testing.T) {
		Data["tags"] = postgres.Jsonb{
			RawMessage: []byte(`"test"`),
		}

		e.POST(basePath).
			WithJSON(Data).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusUnprocessableEntity)

		Data["tags"] = postgres.Jsonb{
			RawMessage: []byte(`{"test":"tag"}`),
		}
	})

	t.Run("create webhook", func(t *testing.T) {

		event.SelectMock(mock)

		mock.ExpectBegin()

		mock.ExpectQuery(`INSERT INTO "webhooks"`).
			WithArgs(tests.AnyTime{}, tests.AnyTime{}, nil, 1, 1, Data["name"], Data["url"], Data["enabled"], Data["tags"]).
			WillReturnRows(sqlmock.
				NewRows([]string{"id"}).
				AddRow(1))

		mock.ExpectQuery(`INSERT INTO "events"`).
			WithArgs(tests.AnyTime{}, tests.AnyTime{}, nil, 1, 1, "event.done", postgres.Jsonb{RawMessage: []byte(`{"test":"tag"}`)}, 1).
			WillReturnRows(sqlmock.
				NewRows([]string{"id"}).
				AddRow(1))

		mock.ExpectExec(`INSERT INTO "webhook_events"`).
			WithArgs(1, 1).
			WillReturnResult(driver.ResultNoRows)

		mock.ExpectCommit()

		SelectMock(mock)
		EventAssociationMock(mock)

		e.POST(basePath).
			WithJSON(Data).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusCreated)
		tests.ExpectationsMet(t, mock)
	})
}
