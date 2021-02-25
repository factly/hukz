package event

import (
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm/dialects/postgres"
)

var Data = map[string]interface{}{
	"name": "event.done",
	"tags": postgres.Jsonb{
		RawMessage: []byte(`{"test":"tag"}`),
	},
}

var eventList = []map[string]interface{}{
	{
		"name": "event.done1",
		"tags": postgres.Jsonb{
			RawMessage: []byte(`{"test":"tag1"}`),
		},
	},
	{
		"name": "event.done2",
		"tags": postgres.Jsonb{
			RawMessage: []byte(`{"test":"tag2"}`),
		},
	},
}

var invalidData = map[string]interface{}{
	"nae": "a",
}

var Columns = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "name", "tags"}

var selectQuery = regexp.QuoteMeta(`SELECT * FROM "events"`)
var deleteQuery = regexp.QuoteMeta(`UPDATE "events" SET "deleted_at"=`)

var basePath = "/events"
var path = "/events/{event_id}"

func SelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(Columns).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, Data["name"], Data["tags"]))

}
