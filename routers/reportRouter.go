package routers

import (
	"go_backend/api"
	"go_backend/db"
	"go_backend/stats"
	"go_backend/utils"
	"net/http"
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type ReportRouter struct {
	Gin          *gin.Engine
	StatsManager *stats.StatManager
}

func (n *ReportRouter) Init() {
	n.Gin.GET("/reports/in_out_for_person", n.GetInOutForPerson)
}

func (n *ReportRouter) GetInOutForPerson(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	resReq := db.GetConection().Model(&api.Event{}).Joins("User")
	api.Event{}.SetFilterExpression(&queryParams, resReq)
	resReq = resReq.Where(`door_index = '1'`)
	resReq = resReq.Where(`event_type = '32'`)

	order, field := api.Event{}.GetOrderExpression(&queryParams)
	if (!utils.SliceContains([]string{"fullname", "event_1_date"}, field)) {
		order = "event_time asc"
		field = "event_1_date"
	}
	resReq = resReq.Order(order)

	var eventList []api.Event = make([]api.Event, 0)
	resReq.Find(&eventList)


	// // map user_id -> event
	// userMap := make(map[int][]api.Event)
	// for _, event := range eventList {
	// 	_, ok := userMap[event.UserId]
	// 	if !ok {
	// 		userMap[event.UserId] = []api.Event{event}
	// 	} else {
	// 		userMap[event.UserId] = append(userMap[event.UserId], event)
	// 	}
	// }

	// outArr := make([][]api.Event, 0)
	// for userId, eventList := range userMap {
	// 	if len(eventList) > 2 {
	// 		firstEvent := userMap[userId][0]
	// 		lastEvent := userMap[userId][len(userMap[userId])-1]

	// 		outArr = append(outArr, []api.Event{firstEvent, lastEvent})
	// 	} else {
	// 		outArr = append(outArr, eventList)
	// 	}
	// }

	keyMap := make(map[string]string)

	// map user_id -> event
	userMap := make(map[string][]api.Event)
	for _, event := range eventList {
		var key string = fmt.Sprintf("%d", event.User.ID)
		
		_, ok := userMap[key]
		if !ok {
			userMap[key] = []api.Event{event}
		} else {
			userMap[key] = append(userMap[key], event)
		}
		
		var sortKey string = ""
		if (field == "fullname") {
			sortKey = fmt.Sprintf("%s %s %s", event.User.Surname, event.User.FirstName, event.User.MiddleName)
			sortKey = strings.TrimSpace(sortKey)
		} else {
			sortKey = fmt.Sprintf("%d", userMap[key][0].EventTime.Unix())
		}
		keyMap[sortKey] = key
	}



	keys := make([]string, 0, len(keyMap))
	for k := range keyMap{
        keys = append(keys, k)
    }
    sort.Strings(keys)



	outArr := make([][]api.Event, 0)
	for _, key := range keys {
		var userId string = keyMap[key]
		var eventList []api.Event = userMap[userId]

		if len(eventList) > 2 {
			firstEvent := userMap[userId][0]
			lastEvent := userMap[userId][len(userMap[userId])-1]

			outArr = append(outArr, []api.Event{firstEvent, lastEvent})
		} else {
			outArr = append(outArr, eventList)
		}
	}

	//fmt.Printf("\n")
	//fmt.Printf("%+v\n", outArr)

	c.JSON(http.StatusOK, outArr)

	*n.StatsManager.GetStat("reports").GetRequests += 1
}
