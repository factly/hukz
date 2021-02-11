package model

import "github.com/factly/web-hooks-service/config"

func Migration() {
	_ = config.DB.AutoMigrate(
		&Event{},
		&WebhookLog{},
		&Webhook{},
	)
}
