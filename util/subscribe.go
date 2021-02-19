package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/x/requestx"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nats-io/nats.go"
)

// SubscribeEvents subscribe one or more events
func SubscribeEvents(events ...string) error {
	for _, event := range events {
		_, err := NC.Subscribe(event, FireWebhooks)
		if err != nil {
			return err
		}
	}
	return nil
}

// SubscribeExistingEvents subscribe existing events
func SubscribeExistingEvents() error {
	events := make([]model.Event, 0)

	config.DB.Model(&model.Event{}).Find(&events)

	eventNames := make([]string, 0)
	if len(events) > 0 {
		for _, event := range events {
			eventNames = append(eventNames, event.Name)
		}
	}

	return SubscribeEvents(eventNames...)
}

// FireWebhooks fires webhook for a event
func FireWebhooks(m *nats.Msg) {
	var payload map[string]interface{}
	_ = json.Unmarshal(m.Data, &payload)

	whData := model.WebhookData{}

	whData.Event = m.Subject
	whData.CreatedAt = time.Now()
	whData.Contains = []string{strings.Split(m.Subject, ".")[0]}
	whData.Payload = payload

	fmt.Printf("Received a [%v] event with data: %v\n", m.Subject, payload)

	// Fetch event id
	event := model.Event{}
	err := config.DB.Model(&model.Event{}).Where("name = ?", m.Subject).First(&event).Error
	if err != nil {
		return
	}

	// find all the registered webhooks for given event
	webhooks := make([]model.Webhook, 0)
	config.DB.Model(&model.Webhook{}).Joins("JOIN webhook_events ON webhooks.id = webhook_events.webhook_id AND event_id = ?", event.ID).Where("enabled = true").Find(&webhooks)

	for _, webhook := range webhooks {
		go PostWebhook(webhook, event.Name, whData)
	}
}

// PostWebhook does POST request to given URL
func PostWebhook(wh model.Webhook, event string, whData model.WebhookData) {
	bArr, _ := json.Marshal(whData)

	webHookLog := model.WebhookLog{
		CreatedAt:   time.Now(),
		Event:       event,
		URL:         wh.URL,
		Data:        postgres.Jsonb{RawMessage: bArr},
		CreatedByID: wh.CreatedByID,
		Tags:        wh.Tags,
	}

	resp, err := requestx.Request("POST", wh.URL, whData, nil)
	if err != nil {
		fmt.Println("webhook at ", wh.URL, "failed...")
		return
	}

	webHookLog.ResponseStatusCode = resp.StatusCode

	// Create a log entry for webhook
	config.DB.Create(&webHookLog)
}
