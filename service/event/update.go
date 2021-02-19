package event

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/hukz/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update event by id
// @Summary Update a event by id
// @Description Update event by ID
// @Tags Events
// @ID update-event-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param event_id path string true "Event ID"
// @Param Event body event false "Event Object"
// @Success 200 {object} model.Event
// @Router /events/{event_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "event_id")
	id, err := strconv.Atoi(eventID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	event := &event{}

	if err = json.NewDecoder(r.Body).Decode(&event); err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	if validationError := validationx.Check(event); validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
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

	// check if new event name exist in db
	if event.Name != result.Name {
		newName := strings.ToLower(strings.TrimSpace(event.Name))
		var sameNameCount int64
		config.DB.Model(&model.Event{}).Where("name ILIKE (?)", newName).Count(&sameNameCount)

		if sameNameCount > 0 {
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.GetMessage("event with same name exist", http.StatusUnprocessableEntity)))
			return
		}
	}

	// validate tags
	var tags map[string]string
	err = json.Unmarshal(event.Tags.RawMessage, &tags)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	updatedEvent := model.Event{
		Base: model.Base{UpdatedByID: uint(uID)},
		Name: event.Name,
		Tags: event.Tags,
	}

	if err = config.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Model(&result).Updates(updatedEvent).Error; err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	if err = util.SubscribeEvents(result.Name); err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
