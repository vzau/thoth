package events

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dhawton/log4g"
	"github.com/gin-gonic/gin"
	"github.com/vzau/thoth/internal/server/response"
	"github.com/vzau/thoth/pkg/database"
	dbTypes "github.com/vzau/types/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var log = log4g.Category("controllers/events")

func GetEvents(c *gin.Context) {
	events := []dbTypes.Event{}
	_, exists := c.GetQuery("past")

	if exists {
		if err := database.DB.Preload(clause.Associations).Find(&events).Order("end DESC").Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.RespondError(c, http.StatusNotFound, "Not Found")
				return
			}
			log.Error("Error querying for events: %s", err)
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	} else {
		if err := database.DB.Preload(clause.Associations).Where("end > ?", time.Now().Add(time.Hour*3)).Find(&events).Order("end DESC").Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				response.RespondError(c, http.StatusNotFound, "Not Found")
				return
			}
			log.Error("Error querying for events: %s", err)
			response.RespondError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	response.Respond(c, http.StatusOK, struct {
		Events []EventDTO `json:"events" xml:"events" yaml:"events"`
	}{changeBannerUrl(events)})
}

func changeBannerUrl(events []dbTypes.Event) []EventDTO {
	ret := []EventDTO{}

	for i := range events {
		ret = append(ret, EventDTO{
			ID:          events[i].ID,
			Title:       events[i].Title,
			Description: events[i].Description,
			Start:       events[i].Start,
			End:         events[i].End,
			Positions:   events[i].Positions,
			SignUps:     events[i].SignUps,
			Banner:      fmt.Sprintf("/v1/cdn/%d", events[i].Banner.ID),
			CreatedAt:   events[i].CreatedAt,
			UpdatedAt:   events[i].UpdatedAt,
		})
	}
	return ret
}
