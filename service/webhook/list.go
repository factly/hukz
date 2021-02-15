package webhook

import (
	"net/http"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

type paging struct {
	Total int64           `json:"total"`
	Nodes []model.Webhook `json:"nodes"`
}

// list - Get all webhooks
// @Summary Show all webhooks
// @Description Get all webhooks
// @Tags Webhooks
// @ID get-all-webhooks
// @Produce json
// @Param X-User header string true "User ID"
// @Param limit query string false "limit per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /webhooks [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}
	result.Nodes = make([]model.Webhook, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	config.DB.Model(&model.Webhook{}).Count(&result.Total).Offset(offset).Limit(limit).Preload("Events").Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
