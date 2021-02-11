package model

import "gorm.io/gorm"

// Event model
type Event struct {
	Base
	Name     string    `gorm:"column:name" json:"name"`
	Webhooks []Webhook `gorm:"many2many:webhook_events" json:"events,omitempty"`
}

var eventUser ContextKey = "event_user"

// BeforeCreate hook
func (event *Event) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	userID := ctx.Value(eventUser)

	if userID == nil {
		return nil
	}
	uID := userID.(int)

	event.CreatedByID = uint(uID)
	event.UpdatedByID = uint(uID)
	return nil
}
