package event

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/factly/web-hooks-service/config"
	"github.com/factly/web-hooks-service/model"
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

	result := &model.Event{
		Name: event.Name,
	}

	if err = config.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Create(&result).Error; err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	renderx.JSON(w, http.StatusCreated, result)
}
