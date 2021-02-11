package event

import (
	"net/http"
	"strconv"

	"github.com/factly/web-hooks-service/config"
	"github.com/factly/web-hooks-service/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get event by id
// @Summary Show a event by id
// @Description Get event by ID
// @Tags Events
// @ID get-event-by-id
// @Produce json
// @Param X-User header string true "User ID"
// @Param event_id path string true "Event ID"
// @Success 200 {object} model.Event
// @Router /events/{event_id} [get]
func details(w http.ResponseWriter, r *http.Request) {
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

	renderx.JSON(w, http.StatusOK, result)
}
