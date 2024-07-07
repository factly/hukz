package model

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

// Webhook model
type Webhook struct {
	Base
	Name    string         `gorm:"column:name" json:"name"`
	URL     string         `gorm:"column:url" json:"url"`
	Enabled bool           `gorm:"column:enabled" json:"enabled"`
	Events  []Event        `gorm:"many2many:webhook_events;" json:"events"`
	Tags    postgres.Jsonb `gorm:"column:tags" json:"tags" swaggertype:"primitive,string"`
	SpaceID uuid.UUID      `gorm:"column:space_id" json:"space_id"`
}

var webhookUser ContextKey = "webhook_user"

// BeforeCreate hook
func (wh *Webhook) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(webhookUser)

	if userID == nil {
		return nil
	}
	uID := userID.(string)

	wh.CreatedByID = uID
	wh.UpdatedByID = uID
	wh.ID = uuid.New()

	return nil
}
