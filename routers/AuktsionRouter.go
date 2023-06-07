package routers

import (
	"fmt"
	"go_backend/db"
	"go_backend/logger"
	"go_backend/prop_manager"
	"go_backend/stats"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Aukt struct {
	Gin          *gin.Engine
	Props        *prop_manager.AppProperties
	StatsManager *stats.StatManager
}

type Aukt_Table_Row struct {
	Number     string    `gorm:"column:number primaryKey" json:"number"`
	DataOfSet  time.Time `gorm:"column:dataofset json:"dataofset"`
	Seller     string    `gorm:"column:seller" json:"seller"`
	Object     string    `gorm:"column:object" json:"object"`
	WhoGaveMax string    `gorm:"column:whogavemax" json:"WhoGaveMax"`
	Money      int       `gorm:"column:money" json:"money"`
	TimeToOut  int       `gorm:"column:timetoout" json:"timeToOut"`
	IsActive   bool      `gorm:"column:isactive" json:"isActive"`
}

var timer Time

func (Aukt_Table_Row) TableName() string {
	return "aukst"
}

// var rows []Aukt_Table_Row = make([]Aukt_Table_Row, 0)

func (n *Aukt) Init() {
	db.GetConection().AutoMigrate(&Aukt_Table_Row{})
	n.Gin.GET("/get_table_data", n.Aukt_All)
	n.Gin.GET("/aukt/:id", n.Aukt_By_ID)

	timer := time.NewTimer(time.Minute)
	select {
	case <-timer.C:
		n.TimeOut()
		timer = time.NewTimer(time.Minute)
	}

}

func (n *Aukt) Aukt_All(c *gin.Context) {
	var current_rows []Aukt_Table_Row
	// db.GetConection().Model(&api.SysUser{}).Where("id = ? AND session_key = ?", requestObject.ID, requestObject.SessionKey).Limit(1).Find(&sysUserList)
	db.GetConection().Model(&Aukt_Table_Row{}).Select("number, seller, object, money, whogavemax").Where("isactive = ?", 1).Find(&current_rows)
	fmt.Println(current_rows)
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

// needs localstorage.username and new Coast
// seller and object to find a nessesary object in the table
func (n *Aukt) Update_Coast(c *gin.Context) {
	newData, ok := c.Params.Get("data")
	mass := strings.Split(newData, ";")
	newMoney, err := strconv.Atoi(mass[0])
	if err != nil {
		logger.Logger.Log(logger.Logger{}, logger.INFO, "Cannot execute new money from request to update cost")
	}
	seller := mass[1]
	object := mass[2]
	whogavemax := mass[3]
	localname := mass[4]

	//Seller can't gave new Coast
	if localname == seller {
		c.JSON(http.StatusForbidden, "")
	}
	if !ok {
		c.JSON(http.StatusNotFound, "")
	}
	toFind := Aukt_Table_Row{
		Seller:   seller,
		Object:   object,
		IsActive: true,
	}
	db.GetConection().First(&toFind)
	if toFind.Money != 0 {
		toFind.Money = newMoney
		toFind.WhoGaveMax = whogavemax
	}
	db.GetConection().Save(&toFind)
	c.JSON(http.StatusAccepted, "")

}

// reading object and coast
// add the lolcalstorage.username as seller
// whogavemax is nil, TimeToOut = 3600, IsActive True, dataofset = time.Now()
func (n *Aukt) AddDataToTable() {

}

// cals every minite if timeofset + timeToOut > Max -> set isActive to false
// good for small database
func (n *Aukt) TimeOut() {
	var current_rows []Aukt_Table_Row
	db.GetConection().Model(&Aukt_Table_Row{}).Select("dataofset,timetoout,isactive").Where("isactive = ?", 1).Find(&current_rows)
	fmt.Println(current_rows)

	for _, elem := range current_rows {
		dateNow := time.Now()
		//calculate sub
		if dateNow.Sub(elem.DataOfSet) > time.Duration(elem.TimeToOut) {
			elem.IsActive = false
		}
	}
	db.GetConection().Save(&current_rows)

}
