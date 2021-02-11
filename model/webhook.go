package model

// Webhook model
type Webhook struct {
	Base
	URL    string  `gorm:"column:url" json:"url"`
	Events []Event `gorm:"many2many:webhook_events;" json:"events"`
}
