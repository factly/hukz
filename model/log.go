package model

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// WebhookLog model
type WebhookLog struct {
	ID                 uint           `gorm:"primary_key" json:"id"`
	CreatedAt          time.Time      `json:"created_at"`
	CreatedByID        uint           `gorm:"column:created_by_id" json:"created_by_id"`
	Event              string         `gorm:"column:event" json:"event"`
	URL                string         `gorm:"column:url" json:"url"`
	ResponseStatusCode int            `gorm:"column:response_status_code" json:"response_status_code"`
	Data               postgres.Jsonb `gorm:"column:data" json:"data" swaggertype:"primitive,string"`
}
