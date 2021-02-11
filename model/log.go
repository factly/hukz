package model

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// WebhookLog model
type WebhookLog struct {
	ID                 uint           `gorm:"primary_key" json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	Event              string         `gorm:"column:event" json:"event"`
	URL                string         `gorm:"column:url" json:"url"`
	ResponseStatusCode int            `gorm:"column:response_status_code" json:"response_status_code"`
	Payload            postgres.Jsonb `gorm:"column:payload" json:"payload" swaggertype:"primitive,string"`
}
