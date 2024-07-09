package webhook

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/paginationx"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
)

// list - Get all webhooks logs
// @Summary Show all webhooks logs
// @Description Get all webhooks logs
// @Tags Webhooks
// @ID get-all-webhooks-logs
// @Produce json
// @Param X-User header string true "User ID"
// @Param limit query string false "limit per page"
// @Param page query string false "page number"
// @Param tag query string false "tags"
// @Success 200 {object} paging
// @Router /webhooks/logs [get]
func webhooklogs(w http.ResponseWriter, r *http.Request) {
	spaceID := chi.URLParam(r, "space_id")
	sID, err := strconv.Atoi(spaceID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	webhookID := chi.URLParam(r, "webhook_id")
	wID, err := strconv.Atoi(webhookID)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	result := logPaging{}
	result.Nodes = make([]model.WebhookLog, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	queryMap := r.URL.Query()

	webhookLogsList := make([]model.WebhookLog, 0)
	err = config.DB.Model(&model.WebhookLog{}).
		Where(&model.WebhookLog{
			WebhookID: uint(wID),
			SpaceID:   uint(sID),
		}).Order("created_at DESC").Find(&webhookLogsList).Error
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}
	spew.Dump(webhookLogsList)
	tags := queryMap["tag"]
	if tags != nil {
		for _, webhook := range webhookLogsList {
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
		result.Nodes = webhookLogsList
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
