package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
	"github.com/go-chi/chi"
)

// update - Update webhook by id
// @Summary Update a webhook by id
// @Description Update webhook by ID
// @Tags Webhooks
// @ID update-webhook-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param webhook_id path string true "Webhook ID"
// @Param Webhook body webhook false "Webhook Object"
// @Success 200 {object} model.Webhook
// @Router /webhooks/{webhook_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "webhook_id")
	id, err := strconv.Atoi(webhookID)

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

	result := &model.Webhook{}
	result.ID = uint(id)

	// check record exists or not
	if err = config.DB.First(&result).Error; err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	tx := config.DB.WithContext(context.WithValue(r.Context(), userContext, uID)).Begin()

	newEvents := make([]model.Event, 0)
	if len(webhook.EventIDs) > 0 {
		config.DB.Model(&model.Event{}).Where(webhook.EventIDs).Find(&newEvents)
		if err = tx.Model(&result).Association("Events").Replace(&newEvents); err != nil {
			tx.Rollback()
			loggerx.Error(err)
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	} else {
		_ = config.DB.Model(&result).Association("Events").Clear()
	}

	tx.Model(&result).Select("Enabled").Omit("Events").Updates(model.Webhook{Enabled: webhook.Enabled})

	updatedWebhook := model.Webhook{
		Base: model.Base{UpdatedByID: uint(uID)},
		URL:  webhook.URL,
	}

	if err = tx.Model(&result).Updates(updatedWebhook).Preload("Events").First(&result).Error; err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusOK, result)
}
