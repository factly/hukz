package webhook

import (
	"net/http"
	"strconv"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// details - Get webhook by id
// @Summary Show a webhook by id
// @Description Get webhook by ID
// @Tags Webhooks
// @ID get-webhook-by-id
// @Produce json
// @Param X-User header string true "User ID"
// @Param webhook_id path string true "Webhook ID"
// @Success 200 {object} model.Webhook
// @Router /webhooks/{webhook_id} [get]
func details(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "webhook_id")
	id, err := strconv.Atoi(webhookID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := &model.Webhook{}
	result.ID = uint(id)

	// check record exists or not
	if err = config.DB.Preload("Events").First(&result).Error; err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	renderx.JSON(w, http.StatusOK, result)
}
