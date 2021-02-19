package model

import "github.com/factly/hukz/config"

func Migration() {
	_ = config.DB.AutoMigrate(
		&Event{},
		&WebhookLog{},
		&Webhook{},
	)
}
