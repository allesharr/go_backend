package routers

import (
	"go_backend/db"
	"go_backend/prop_manager"
	"go_backend/stats"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Aukt struct {
	Gin          *gin.Engine
	Props        *prop_manager.AppProperties
	StatsManager *stats.StatManager
}

type Aukt_Table_Row struct {
	Number     string `gorm:"column:number" json:"number"`
	Seller     string `gorm:"column:seller" json:"seller"`
	Object     string `gorm:"column:object" json:"object"`
	WhoGaveMax string `gorm:"column:whogavemax" json:"WhoGaveMax"`
	Money      int    `gorm:"column:money" json:"money"`
	TimeToOut  int    `gorm:"column:timetoout" json:"timeToOut"`
	IsActive   bool   `gorm:"column:isactive" json:"isActive"`
}

func (Aukt_Table_Row) TableName() string {
	return "aukst"
}

// var rows []Aukt_Table_Row = make([]Aukt_Table_Row, 0)

func (n *Aukt) Init() {
	n.Gin.GET("/get_table_data", n.Aukt_All)
	n.Gin.GET("/aukt/:id", n.Aukt_By_ID)
}

func (n *Aukt) Aukt_All(c *gin.Context) {
	var current_rows []Aukt_Table_Row
	// db.GetConection().Model(&api.SysUser{}).Where("id = ? AND session_key = ?", requestObject.ID, requestObject.SessionKey).Limit(1).Find(&sysUserList)
	db.GetConection().Model(&Aukt_Table_Row{}).Where("is_active = ?", true).Find(&current_rows)
	c.JSON(http.StatusOK, current_rows)
}

func (n *Aukt) Aukt_By_ID(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, "")
	}
	var current_rows []Aukt_Table_Row
	db.GetConection().Model(&Aukt_Table_Row{}).Where("is_active = ? AND number = ?", true, id).Find(&current_rows)
}
