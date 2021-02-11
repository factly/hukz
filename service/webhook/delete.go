package webhook

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

// delete - Delete webhook by id
// @Summary Delete webhook by id
// @Description Delete webhook by id
// @Tags Webhooks
// @ID delete-webhook-by-id
// @Param X-User header string true "User ID"
// @Param webhook_id path string true "Webhook ID"
// @Success 200
// @Failure 400 {array} string
// @Router  /webhooks/{webhook_id} [delete]
func delete(w http.ResponseWriter, r *http.Request) {
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
	if err = config.DB.First(&result).Error; err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.RecordNotFound()))
		return
	}

	config.DB.Delete(&result)

	renderx.JSON(w, http.StatusOK, nil)
}
