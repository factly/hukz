package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	googlechat "github.com/factly/x/hukzx/google_chat"
	"github.com/factly/x/hukzx/slack"
	"github.com/factly/x/requestx"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

// SubscribeEvents subscribe one or more events
var SubscribeEvents = func(events ...string) error {
	for _, event := range events {
		_, err := NC.QueueSubscribe(event, viper.GetString("queue_group"), FireWebhooks)
		if err != nil {
			return err
		}
	}
	return nil
}

// SubscribeExistingEvents subscribe existing events
var SubscribeExistingEvents = func() error {
	events := make([]model.Event, 0)

	config.DB.Model(&model.Event{}).Find(&events)

	eventNames := make([]string, 0)
	if len(events) > 0 {
		for _, event := range events {
			eventNames = append(eventNames, event.Event)
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
	err := config.DB.Model(&model.Event{}).Where("event = ?", m.Subject).First(&event).Error
	if err != nil {
		return
	}

	// find all the registered webhooks for given event
	webhooks := make([]model.Webhook, 0)
	config.DB.Model(&model.Webhook{}).Joins("JOIN webhook_events ON webhooks.id = webhook_events.webhook_id AND event_id = ?", event.ID).Where("enabled = true").Find(&webhooks)

	for _, webhook := range webhooks {
		go PostWebhook(webhook, event.Event, whData)
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

	var resp *http.Response
	var err error
	if strings.Contains(wh.URL, "chat.googleapis.com") && config.DegaToGoogleChat() {
		data, err := googlechat.ToMessage(whData)
		if err != nil {
			fmt.Println("webhook at ", wh.URL, "failed...")
			return
		}

		bArr, _ = json.Marshal(data)
		webHookLog.Data = postgres.Jsonb{RawMessage: bArr}

		if resp, err = requestx.Request("POST", wh.URL, data, nil); err != nil {
			fmt.Println("webhook at ", wh.URL, "failed...")
			return
		}
	} else if strings.Contains(wh.URL, "hooks.slack.com") && config.DegaToSlack() {
		data, err := slack.ToMessage(whData)
		if err != nil {
			fmt.Println("webhook at ", wh.URL, "failed...")
			return
		}

		bArr, _ = json.Marshal(data)
		webHookLog.Data = postgres.Jsonb{RawMessage: bArr}

		if resp, err = requestx.Request("POST", wh.URL, data, nil); err != nil {
			fmt.Println("webhook at ", wh.URL, "failed...")
			return
		}

	} else if strings.Contains(wh.URL, "chat.googleapis.com") {
		message := googlechat.DefaultMessage(whData)
		if resp, err = requestx.Request("POST", wh.URL, message, nil); err != nil {
			fmt.Println("webhook at ", wh.URL, "failed...")
			return
		}
	} else if strings.Contains(wh.URL, "hooks.slack.com") {
		message, _ := slack.DefaultMessage(whData)
		if resp, err = requestx.Request("POST", wh.URL, message, nil); err != nil {
			fmt.Println("webhook at ", wh.URL, "failed...")
			return
		}

	} else {
		if resp, err = requestx.Request("POST", wh.URL, whData, nil); err != nil {
			fmt.Println("webhook at ", wh.URL, "failed...")
			return
		}
	}

	defer resp.Body.Close()

	webHookLog.ResponseStatusCode = resp.StatusCode

	body_bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if strings.Contains(wh.URL, "hooks.slack.com") {
		body_bytes = []byte(fmt.Sprint(`{"data":"`, string(body_bytes), `"}`))
	}

	webHookLog.ResponseBody = postgres.Jsonb{
		RawMessage: body_bytes,
	}

	// Create a log entry for webhook
	config.DB.Create(&webHookLog)
}
