package model

import "time"

// WebhookData general struct to fire webhook
type WebhookData struct {
	Event     string      `json:"event,omitempty"`
	Contains  []string    `json:"contains,omitempty"`
	CreatedAt time.Time   `json:"created_at,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
}
