package event

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/factly/hukz/config"
	"github.com/factly/hukz/model"
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
// @Param tag query string false "tags"
// @Success 200 {object} paging
// @Router /events [get]
func list(w http.ResponseWriter, r *http.Request) {

	result := paging{}
	result.Nodes = make([]model.Event, 0)

	offset, limit := paginationx.Parse(r.URL.Query())

	queryMap := r.URL.Query()

	eventList := make([]model.Event, 0)
	config.DB.Model(&model.Event{}).Offset(offset).Limit(limit).Find(&eventList)

	tags := queryMap["tag"]
	if tags != nil {
		for _, event := range eventList {
			var tagMap map[string]string
			_ = json.Unmarshal(event.Tags.RawMessage, &tagMap)

			count := 0
			for _, t := range tags {
				toks := strings.Split(t, ":")
				if val, found := tagMap[toks[0]]; found && val == toks[1] {
					count++
				}
			}
			if count == len(tags) {
				result.Nodes = append(result.Nodes, event)
			}
		}
	} else {
		result.Nodes = eventList
	}

	result.Total = int64(len(result.Nodes))

	renderx.JSON(w, http.StatusOK, result)
}
