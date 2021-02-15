package webhook

import (
	"net/http"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

type logPaging struct {
	Total int64              `json:"total"`
	Nodes []model.WebhookLog `json:"nodes"`
}

// list - Get all webhooks logs
// @Summary Show all webhooks logs
// @Description Get all webhooks logs
// @Tags Webhooks
// @ID get-all-webhooks-logs
// @Produce json
// @Param X-User header string true "User ID"
// @Param limit query string false "limit per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /webhooks/logs [get]
func logs(w http.ResponseWriter, r *http.Request) {

	result := logPaging{}
	result.Nodes = make([]model.WebhookLog, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	config.DB.Model(&model.WebhookLog{}).Count(&result.Total).Offset(offset).Limit(limit).Order("created_at DESC").Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
