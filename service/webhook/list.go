package webhook

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
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
// @Param tag query string false "tags"
// @Success 200 {object} paging
// @Router /webhooks [get]
func list(w http.ResponseWriter, r *http.Request) {
	spaceID := chi.URLParam(r, "space_id")
	id, err := strconv.Atoi(spaceID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := paging{}
	result.Nodes = make([]model.Webhook, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	queryMap := r.URL.Query()

	webhookList := make([]model.Webhook, 0)
	config.DB.Model(&model.Webhook{}).Where(&model.Webhook{
		SpaceID: uint(id),
	}).Count(&result.Total).Preload("Events").Find(&webhookList)

	tags := queryMap["tag"]
	if tags != nil {
		for _, webhook := range webhookList {
			var tagMap map[string]string
			_ = json.Unmarshal(webhook.Tags.RawMessage, &tagMap)

			count := 0
			for _, t := range tags {
				toks := strings.Split(t, ":")
				if val, found := tagMap[toks[0]]; found && val == toks[1] {
					count++
				}
			}
			if count == len(tags) {
				result.Nodes = append(result.Nodes, webhook)
			}
		}
	} else {
		result.Nodes = webhookList
	}

	var end int
	if limit+offset > len(result.Nodes) {
		end = len(result.Nodes)
	} else {
		end = limit + offset
	}
	if offset > len(result.Nodes) {
		offset = len(result.Nodes)
	} else if offset < 0 {
		offset = 0
	}

	result.Nodes = result.Nodes[offset:end]
	result.Total = int64(len(result.Nodes))

	renderx.JSON(w, http.StatusOK, result)
}
