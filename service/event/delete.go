package event

import (
	"net/http"
	"strconv"

	"github.com/factly/web-hooks-service/config"
	"github.com/factly/web-hooks-service/model"
	"github.com/factly/web-hooks-service/util"
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

	config.DB.Delete(&result)

	// reconnect nats server
	util.NC.Close()
	util.ConnectNats()
	util.SubscribeExistingEvents()

	renderx.JSON(w, http.StatusOK, nil)
}
