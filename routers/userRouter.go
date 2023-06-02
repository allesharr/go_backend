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

type UserRouter struct {
	Gin          *gin.Engine
	StatsManager *stats.StatManager
}

func (n *UserRouter) Init() {
	n.Gin.GET("/user", n.Get)
	n.Gin.GET("/user/:id", n.GetObject)
}

func (n *UserRouter) Get(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	countReq := db.GetConection().Model(&api.User{})
	api.User{}.SetFilterExpression(&queryParams, countReq)

	var totalRecords int64
	countReq.Count(&totalRecords)
	c.Writer.Header()["total_records"] = []string{strconv.FormatInt(totalRecords, 10)}

	resReq := db.GetConection().Model(&api.User{})
	api.User{}.SetFilterExpression(&queryParams, resReq)
	resReq = resReq.Order(api.User{}.GetOrderExpression(&queryParams))
	utils.SetPaging(&queryParams, resReq)

	var userList []api.User = make([]api.User, 0)
	resReq.Find(&userList)
	c.JSON(http.StatusOK, userList)

	*n.StatsManager.GetStat("user").GetRequests += 1
}

func (n *UserRouter) GetObject(c *gin.Context) {
	var dbObject []api.User
	db.GetConection().Model(&api.User{}).Where("id = ?", c.Param("id")).Limit(1).Find(&dbObject)

	if len(dbObject) == 0 {
		c.JSON(http.StatusNotFound, api.ResponseMessage{Message: "Объект не найден"})
		*n.StatsManager.GetStat("user").Errors += 1
	} else {
		c.JSON(http.StatusOK, dbObject[0])
		*n.StatsManager.GetStat("user").GetByIdRequests += 1
	}
}
