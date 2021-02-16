package webhook

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
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
// @Param tag query string false "tags"
// @Success 200 {object} paging
// @Router /webhooks/logs [get]
func logs(w http.ResponseWriter, r *http.Request) {

	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	result := logPaging{}
	result.Nodes = make([]model.WebhookLog, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	queryMap := r.URL.Query()

	webhookLogsList := make([]model.WebhookLog, 0)
	config.DB.Model(&model.WebhookLog{}).Where(&model.WebhookLog{
		CreatedByID: uint(uID),
	}).Count(&result.Total).Offset(offset).Limit(limit).Order("created_at DESC").Find(&webhookLogsList)

	tags := queryMap["tag"]
	if tags != nil {
		for _, webhook := range webhookLogsList {
			var tagMap map[string]string
			_ = json.Unmarshal(webhook.Tags.RawMessage, &tagMap)

			for _, t := range tags {
				toks := strings.Split(t, ":")
				if val, found := tagMap[toks[0]]; found && val == toks[1] {
					result.Nodes = append(result.Nodes, webhook)
					break
				}
			}
		}
	} else {
		result.Nodes = webhookLogsList
	}

	result.Total = int64(len(result.Nodes))

	renderx.JSON(w, http.StatusOK, result)
}
