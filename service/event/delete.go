package event

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/hukz/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// delete - Delete Event by id
// @Summary Delete Event by id
// @Description Delete Event by id
// @Tags Events
// @ID delete-event-by-id
// @Param X-User header string true "User ID"
// @Param event_id path string true "Event ID"
// @Success 200
// @Failure 400 {array} string
// @Router  /events/{event_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")
	id, err := strconv.Atoi(eventID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Event{}
	result.ID = uint(id)

	// check record exists or not
	if err = config.DB.First(&result).Error; err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	// check if event is associated with webhook
	event := new(model.Event)
	event.ID = uint(id)
	totAssociated := config.DB.Model(event).Association("Webhooks").Count()

	if totAssociated != 0 {
		loggerx.Error(errors.New("event is associated with webhook"))
		errorx.Render(w, errorx.Parser(errorx.CannotDelete("event", "webhook")))
		return
	}

	config.DB.Delete(&result)

	// reconnect nats server
	if util.NC != nil {
		util.NC.Close()
		util.ConnectNats()
	}

	if err = util.SubscribeExistingEvents(); err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	renderx.JSON(w, http.StatusOK, nil)
}
