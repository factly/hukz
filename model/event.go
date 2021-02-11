package model

// Event model
type Event struct {
	Base
	Name     string    `gorm:"column:name" json:"name"`
	Webhooks []Webhook `gorm:"many2many:webhook_events" json:"events,omitempty"`
}
