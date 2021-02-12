package util

import (
	"encoding/json"
	"fmt"

	"github.com/factly/web-hooks-service/config"
	"github.com/factly/web-hooks-service/model"
	"github.com/nats-io/nats.go"
)

// SubscribeEvents subscribe one or more events
func SubscribeEvents(events ...string) {
	for _, event := range events {
		NC.Subscribe(event, func(m *nats.Msg) {
			var payload map[string]interface{}
			json.Unmarshal(m.Data, &payload)

			fmt.Printf("Received a [%v] event with data: %v\n", m.Subject, payload)
		})
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
