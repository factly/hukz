package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/factly/web-hooks-service/config"
	"github.com/factly/web-hooks-service/model"
	"github.com/factly/x/requestx"
	"github.com/nats-io/nats.go"
)

// SubscribeEvents subscribe one or more events
func SubscribeEvents(events ...string) {
	for _, event := range events {
		NC.Subscribe(event, FireWebhooks)
	}
}

// SubscribeExistingEvents subscribe existing events
func SubscribeExistingEvents() {
	events := make([]model.Event, 0)

	config.DB.Model(&model.Event{}).Find(&events)

	eventNames := make([]string, 0)
	if len(events) > 0 {
		for _, event := range events {
			eventNames = append(eventNames, event.Name)
		}
	}

	SubscribeEvents(eventNames...)
}

// FireWebhooks fires webhook for a event
func FireWebhooks(m *nats.Msg) {
	var payload map[string]interface{}
	json.Unmarshal(m.Data, &payload)

	whData := model.WebhookData{}

	whData.Event = m.Subject
	whData.CreatedAt = time.Now()
	whData.Contains = []string{strings.Split(m.Subject, ".")[0]}
	whData.Payload = payload

	fmt.Printf("Received a [%v] event with data: %v\n", m.Subject, payload)

	// TODO: find all the registered webhooks for given event and fire on each url through go routines

	url := "https://cffaf321acf0589eee315b01d3087c70.m.pipedream.net" // test url

	// Create a log entry in WebhookLog

	resp, err := requestx.Request("POST", url, whData, nil)
	if err != nil {
		fmt.Println("webhook at ", url, "failed...")
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("webhook at ", url, "failed...")
		return
	}
}
