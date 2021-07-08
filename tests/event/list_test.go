package event

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/hukz/service"
	"github.com/factly/hukz/tests"
	"github.com/gavv/httpexpect"
)

func TestEventList(t *testing.T) {
	mock := tests.SetupMockDB()

	testServer := httptest.NewServer(service.RegisterRoutes())
	defer testServer.Close()

	// create httpexpect instance
	e := httpexpect.New(t, testServer.URL)

	t.Run("get empty list of events", func(t *testing.T) {
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

	t.Run("get list of events", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(Columns).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, eventList[0]["name"], eventList[0]["event"], eventList[0]["tags"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, eventList[1]["name"], eventList[1]["name"], eventList[1]["tags"]))

		e.GET(basePath).
			WithHeader("X-User", "1").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 2}).
			Value("nodes").
			Array().Element(0).Object().ContainsMap(eventList[0])

		tests.ExpectationsMet(t, mock)
	})

	t.Run("get list of events for tags", func(t *testing.T) {
		mock.ExpectQuery(selectQuery).
			WillReturnRows(sqlmock.NewRows(Columns).
				AddRow(1, time.Now(), time.Now(), nil, 1, 1, eventList[0]["name"], eventList[0]["event"], eventList[0]["tags"]).
				AddRow(2, time.Now(), time.Now(), nil, 1, 1, eventList[1]["name"], eventList[1]["event"], eventList[1]["tags"]))

		e.GET(basePath).
			WithHeader("X-User", "1").
			WithQuery("tag", "test:tag2").
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 1}).
			Value("nodes").
			Array().Element(0).Object().ContainsMap(eventList[1])

		tests.ExpectationsMet(t, mock)
	})
}
