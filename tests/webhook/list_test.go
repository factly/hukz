package webhook

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/hukz/service"
	"github.com/factly/hukz/tests"
	"github.com/factly/hukz/tests/event"
	"github.com/gavv/httpexpect"
)

func TestWebhookList(t *testing.T) {
	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("get empty list of webhooks", func(t *testing.T) {

		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(Columns))

		e.GET(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})

		tests.ExpectationsMet(t, mock)
	})

	t.Run("get list of webhooks", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(Columns).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, webhookList[0]["name"], webhookList[0]["url"], webhookList[0]["enabled"], webhookList[0]["tags"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, webhookList[1]["name"], webhookList[1]["url"], webhookList[1]["enabled"], webhookList[1]["tags"]))

		EventAssociationMock(mock)

		e.GET(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 2})
		tests.ExpectationsMet(t, mock)
	})

	t.Run("get filtered list of webhooks for tags", func(t *testing.T) {
		mock.ExpectQuery(countQuery).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(Columns).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, webhookList[0]["name"], webhookList[0]["url"], webhookList[0]["enabled"], webhookList[0]["tags"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, webhookList[1]["name"], webhookList[1]["url"], webhookList[1]["enabled"], webhookList[1]["tags"]))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "webhook_events"`)).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"webhook_id", "event_id"}).
				AddRow(1, 1))

		event.SelectMock(mock, 1)

		e.GET(basePath).
			WithHeader("X-User", "1").
			WithQuery("tag", "test:tag2").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 1})
		tests.ExpectationsMet(t, mock)
	})
}
