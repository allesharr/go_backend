package routers

import (
	"fmt"
	"go_backend/api"
	"go_backend/db"
	"go_backend/logger"
	"go_backend/stats"
	"go_backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

type SysUserRouter struct {
	Gin          *gin.Engine
	StatsManager *stats.StatManager
}

func (n *SysUserRouter) Init() {
	n.Gin.GET("/sys_user", n.Get)
	n.Gin.GET("/sys_user/:id", n.GetObject)
	n.Gin.POST("/sys_user", n.Create)
	n.Gin.PUT("/sys_user/:id", n.Update)
	n.Gin.DELETE("/sys_user/:id", n.Delete)
}

func (n *SysUserRouter) Get(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	countReq := db.GetConection().Model(&api.SysUser{})
	api.SysUser{}.SetFilterExpression(&queryParams, countReq)

	var totalRecords int64
	countReq.Count(&totalRecords)
	c.Writer.Header()["total_records"] = []string{strconv.FormatInt(totalRecords, 10)}

	resReq := db.GetConection().Model(&api.SysUser{})
	api.SysUser{}.SetFilterExpression(&queryParams, resReq)
	resReq = resReq.Order(api.SysUser{}.GetOrderExpression(&queryParams))
	utils.SetPaging(&queryParams, resReq)

	var sysUserList []api.SysUser = make([]api.SysUser, 0)
	resReq.Find(&sysUserList)
	c.JSON(http.StatusOK, sysUserList)

	*n.StatsManager.GetStat("sys_user").GetRequests += 1
}

func (n *SysUserRouter) GetObject(c *gin.Context) {
	var dbObject []api.SysUser
	db.GetConection().Model(&api.SysUser{}).Where("system_users.id = ?", c.Param("id")).Select("system_users.id", "login").Limit(1).Find(&dbObject)

	if len(dbObject) == 0 {
		c.JSON(http.StatusNotFound, api.ResponseMessage{Message: "Объект не найден"})
		*n.StatsManager.GetStat("sys_user").Errors += 1
	} else {
		c.JSON(http.StatusOK, dbObject[0])
		*n.StatsManager.GetStat("sys_user").GetByIdRequests += 1
	}
}

func (n *SysUserRouter) Create(c *gin.Context) {
	var requestObject *api.SysUser

	if err := c.BindJSON(&requestObject); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseMessage{Message: fmt.Sprintf("Ошибка парсинга объекта: %s", err.Error())})
		*n.StatsManager.GetStat("sys_user").Errors += 1
		return
	}

	var checkLogin []api.SysUser
	db.GetConection().Model(&api.SysUser{}).Where("login = ?", requestObject.Login).Limit(1).Find(&checkLogin)
	if len(checkLogin) > 0 {
		c.JSON(http.StatusConflict, api.ResponseMessage{Message: "Пользователь с таким логином уже существует"})
		*n.StatsManager.GetStat("sys_user").Errors += 1
		return
	}

	result := db.GetConection().Omit(clause.Associations).Create(requestObject)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseMessage{Message: "Ошибка создания"})
		logger.Logger{}.Log(logger.ERR, fmt.Sprintf("Ошибка создания объекта system user: %s", result.Error))
		*n.StatsManager.GetStat("sys_user").Errors += 1

		return
	}

	c.JSON(http.StatusOK, "")
	*n.StatsManager.GetStat("sys_user").CreateRequests += 1
}

func (n *SysUserRouter) Update(c *gin.Context) {
	var requestObject *api.SysUser

	if err := c.BindJSON(&requestObject); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseMessage{Message: fmt.Sprintf("Ошибка парсинга объекта: %s", err.Error())})
		*n.StatsManager.GetStat("sys_user").Errors += 1
		return
	}

	requestObjectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseMessage{Message: fmt.Sprintf("Ошибка получения идентификатора объекта: %s", err.Error())})
		*n.StatsManager.GetStat("sys_user").Errors += 1
		return
	}

	requestObject.ID = uint(requestObjectID)

	var previousUserList []api.SysUser
	db.GetConection().Model(&api.SysUser{}).Where("id = ?", requestObject.ID).Limit(1).Find(&previousUserList)
	if len(previousUserList) == 0 {
		c.JSON(http.StatusConflict, api.ResponseMessage{Message: fmt.Sprintf("Не найдено объекта с идентификатором %d", requestObject.ID)})
		*n.StatsManager.GetStat("sys_user").Errors += 1
		return
	}

	var checkLogin []api.SysUser
	db.GetConection().Model(&api.SysUser{}).Where("login = ? and id != ?", requestObject.Login, requestObject.ID).Limit(1).Find(&checkLogin)
	if len(checkLogin) > 0 {
		c.JSON(http.StatusConflict, api.ResponseMessage{Message: "Пользователь с таким логином уже существует"})
		*n.StatsManager.GetStat("sys_user").Errors += 1
		return
	}

	updateMap := map[string]interface{}{
		"login": requestObject.Login}

	if requestObject.PasswordHash != "" {
		updateMap["password_hash"] = requestObject.PasswordHash
	}

	result := db.GetConection().Model(api.SysUser{}).Where("id = ?", requestObject.ID).Updates(updateMap)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseMessage{Message: "Ошибка обновления"})
		logger.Logger{}.Log(logger.ERR, fmt.Sprintf("Ошибка обновления объекта system user: %s", result.Error))
		*n.StatsManager.GetStat("sys_user").Errors += 1

		return
	}

	c.JSON(http.StatusOK, "")
	*n.StatsManager.GetStat("sys_user").UpdateRequests += 1
}

func (n *SysUserRouter) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "")
}
