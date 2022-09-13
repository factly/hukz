package webhook

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"

	//"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

type CheckWebhookEvent struct {
	Enabled bool `json:"event_enabled"`
}

// check - check whether a webhook is enabled for given event
// @Summary check whether a webhook is enabled for given event
// @Description check whether a webhook is enabled for given event
// @Produce json
// @Param X-User header string true "User ID"
// @Param event query string false "event"
// @Success 200 {object}
// @Router /webhooks/space/{space_id}/check [get]
func check(w http.ResponseWriter, r *http.Request) {
	result := CheckWebhookEvent{}
	result.Enabled = false

	event := r.URL.Query().Get("event")
	spaceID := chi.URLParam(r, "space_id")

	id, err := strconv.Atoi(spaceID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	webhookList := make([]model.Webhook, 0)
	config.DB.Model(&model.Webhook{}).Where(&model.Webhook{
		SpaceID: uint(id),
	}).Preload("Events").Find(&webhookList)
	for _, webhook := range webhookList {
		if result.Enabled {
			break
		}
		if webhook.Enabled {
			for _, webhookevent := range webhook.Events {
				if webhookevent.Event == event {
					fmt.Println("event is enabled")
					result.Enabled = true
					break
				}

			}
		}
	}

	renderx.JSON(w, http.StatusOK, result)
}
