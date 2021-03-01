package webhook

import (
	"database/sql/driver"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/hukz/tests/event"
	"github.com/jinzhu/gorm/dialects/postgres"
)

var Data = map[string]interface{}{
	"url":       "testurl.com",
	"enabled":   true,
	"event_ids": []uint{1},
	"tags": postgres.Jsonb{
		RawMessage: []byte(`{"test":"tag"}`),
	},
}

var respData = map[string]interface{}{
	"url":     "testurl.com",
	"enabled": true,
	"tags": postgres.Jsonb{
		RawMessage: []byte(`{"test":"tag"}`),
	},
}

var webhookList = []map[string]interface{}{
	{
		"url":       "testurl1.com",
		"enabled":   true,
		"event_ids": []uint{1},
		"tags": postgres.Jsonb{
			RawMessage: []byte(`{"test":"tag1"}`),
		},
	},
	{
		"url":       "testurl2.com",
		"enabled":   true,
		"event_ids": []uint{1},
		"tags": postgres.Jsonb{
			RawMessage: []byte(`{"test":"tag2"}`),
		},
	},
}

var invalidData = map[string]interface{}{
	"nae": "a",
}

var Columns = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "url", "enabled", "tags"}

var selectQuery = `SELECT (.+) FROM "webhooks"`
var deleteQuery = regexp.QuoteMeta(`UPDATE "webhooks" SET "deleted_at"=`)
var countQuery = regexp.QuoteMeta(`SELECT count(1) FROM "webhooks"`)

var basePath = "/webhooks"
var path = "/webhooks/{webhook_id}"

func SelectMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(selectQuery).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows(Columns).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, Data["url"], Data["enabled"], Data["tags"]))
}

func EventAssociationMock(mock sqlmock.Sqlmock, args ...driver.Value) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "webhook_events"`)).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"webhook_id", "event_id"}).
			AddRow(1, 1))

	event.SelectMock(mock, 1)
}
