package routers

import (
	"go_backend/api"
	"go_backend/db"
	"go_backend/stats"
	"go_backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventRouter struct {
	Gin          *gin.Engine
	StatsManager *stats.StatManager
}

func (n *EventRouter) Init() {
	n.Gin.GET("/event", n.Get)
	n.Gin.GET("/event/:id", n.GetObject)
}

func (n *EventRouter) Get(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	countReq := db.GetConection().Model(&api.Event{}).Select("skud_events.id").Joins("User")
	api.Event{}.SetFilterExpression(&queryParams, countReq)

	var totalRecords int64
	countReq.Count(&totalRecords)
	c.Writer.Header()["total_records"] = []string{strconv.FormatInt(totalRecords, 10)}

	resReq := db.GetConection().Model(&api.Event{}).Joins("User")
	api.Event{}.SetFilterExpression(&queryParams, resReq)

	order, _ := api.Event{}.GetOrderExpression(&queryParams)
	resReq = resReq.Order(order)
	utils.SetPaging(&queryParams, resReq)

	var eventList []api.Event = make([]api.Event, 0)
	resReq.Find(&eventList)
	for k, v := range eventList {
		eventList[k].OrgName = v.User.OrgName
	}

	c.JSON(http.StatusOK, eventList)

	*n.StatsManager.GetStat("event").GetRequests += 1
}

func (n *EventRouter) GetObject(c *gin.Context) {
	var dbObject []api.Event
	db.GetConection().Model(&api.Event{}).Joins("User").Where("skud_events.id = ?", c.Param("id")).Limit(1).Find(&dbObject)

	if len(dbObject) == 0 {
		c.JSON(http.StatusNotFound, api.ResponseMessage{Message: "Объект не найден"})
		*n.StatsManager.GetStat("event").Errors += 1
	} else {
		c.JSON(http.StatusOK, dbObject[0])
		*n.StatsManager.GetStat("event").GetByIdRequests += 1
	}
}
