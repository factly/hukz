package event

import (
	"net/http"

	"github.com/factly/web-hooks-service/config"
	"github.com/factly/web-hooks-service/model"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
)

type paging struct {
	Total int64         `json:"total"`
	Nodes []model.Event `json:"nodes"`
}

// list - Get all events
// @Summary Show all events
// @Description Get all events
// @Tags Events
// @ID get-all-events
// @Produce json
// @Param X-User header string true "User ID"
// @Param limit query string false "limit per page"
// @Param page query string false "page number"
// @Success 200 {object} paging
// @Router /events [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}
	result.Nodes = make([]model.Event, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	config.DB.Model(&model.Event{}).Count(&result.Total).Offset(offset).Limit(limit).Find(&result.Nodes)

	renderx.JSON(w, http.StatusOK, result)
}
