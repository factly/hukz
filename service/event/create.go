package event

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/hukz/util"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create Event
// @Summary Create Event
// @Description Create Event
// @Tags Events
// @ID add-event
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param Event body event true "Event Object"
// @Success 201 {object} model.Event
// @Failure 400 {array} string
// @Router /events [post]
func create(w http.ResponseWriter, r *http.Request) {
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

	// check if new event name exist in db
	newName := strings.ToLower(strings.TrimSpace(event.Event))
	var sameNameCount int64
	config.DB.Model(&model.Event{}).Where("event ILIKE (?)", newName).Count(&sameNameCount)

	if sameNameCount > 0 {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.GetMessage("event with same event string exist", http.StatusUnprocessableEntity)))
		return
	}

	// validate tags
	var tags map[string]string
	err = json.Unmarshal(event.Tags.RawMessage, &tags)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	result := &model.Event{
		Name:  event.Name,
		Event: event.Event,
		Tags:  event.Tags,
	}

	if err = config.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Create(&result).Error; err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	if err = util.SubscribeEvents(result.Event); err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, result)
}
