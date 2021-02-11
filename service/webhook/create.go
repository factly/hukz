package webhook

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

// create - Create Webhook
// @Summary Create Webhook
// @Description Create Webhook
// @Tags Webhooks
// @ID add-webhook
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param Webhook body webhook true "Webhook Object"
// @Success 201 {object} model.Webhook
// @Failure 400 {array} string
// @Router /webhooks [post]
func create(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	webhook := &webhook{}

	if err = json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	if validationError := validationx.Check(webhook); validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	result := &model.Webhook{
		URL:     webhook.URL,
		Enabled: webhook.Enabled,
	}

	if len(webhook.EventIDs) > 0 {
		config.DB.Model(&model.Event{}).Where(webhook.EventIDs).Find(&result.Events)
	}

	if err = config.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Create(&result).Error; err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	config.DB.Model(&model.Webhook{}).Preload("Events").First(&result)

	renderx.JSON(w, http.StatusCreated, result)
}
