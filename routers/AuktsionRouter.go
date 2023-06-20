package routers

import (
	"encoding/json"
	"fmt"
	"go_backend/db"
	"go_backend/prop_manager"
	"go_backend/stats"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Aukt struct {
	Gin          *gin.Engine
	Props        *prop_manager.AppProperties
	StatsManager *stats.StatManager
}

type Aukt_Table_Row struct {
	Number     string    `gorm:"primaryKey;column:number" json:"number"`
	DataOfSet  time.Time `gorm:"column:dataofset json:"dataofset"`
	Seller     string    `gorm:"column:seller" json:"seller"`
	Object     string    `gorm:"column:object" json:"object"`
	WhoGaveMax string    `gorm:"column:whogavemax" json:"whogavemax"`
	Money      int       `gorm:"column:money" json:"money"`
	TimeToOut  int       `gorm:"column:timetoout" json:"tto"`
	IsActive   bool      `gorm:"column:isactive" json:"isActive"`
}

var timer *time.Timer

func (Aukt_Table_Row) TableName() string {
	return "aukst"
}


// there is no goroutine exit
func (n *Aukt) Init() {
	db.GetConection().AutoMigrate(&Aukt_Table_Row{})
	n.Gin.GET("/get_table_data", n.Aukt_All)
	n.Gin.GET("/aukt/:id", n.Aukt_By_ID)
	n.Gin.POST("/update", n.Update_Coast)
	n.Gin.POST("/set_lot", n.Set_Lot)
	ticker := time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				n.TimeOut()
			}
		}
	}()

}

func (n *Aukt) Set_Lot(c *gin.Context) {
	data := Aukt_Table_Row{}
	b, err := c.GetRawData()
	if err != nil {
		fmt.Println("Post data is not correct")
	}
	json.Unmarshal(b, &data)
	number := strconv.Itoa(int(time.Now().Unix()))
	date := time.Now()
	data.Number = number
	data.DataOfSet = date
	data.IsActive = true

	db.GetConection().Create(&data)

}

func (n *Aukt) Aukt_All(c *gin.Context) {
	var current_rows []Aukt_Table_Row
	
	db.GetConection().Model(&Aukt_Table_Row{}).Select("*").Where("isactive = ?", 1).Find(&current_rows)
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
	toUpdate := Aukt_Table_Row{}
	b, err := c.GetRawData()
	if err != nil {
		fmt.Println("Post data is not correct")
	}
	json.Unmarshal(b, &toUpdate)
	fmt.Println("My collected data is ", toUpdate)
	fmt.Println("number", toUpdate.Number, "whoGaveMax", toUpdate.WhoGaveMax, "money", toUpdate.Money)

	whogavemax := toUpdate.WhoGaveMax
	money := toUpdate.Money

	db.GetConection().Model(&toUpdate).Update("whogavemax", whogavemax)
	db.GetConection().Model(&toUpdate).Update("money", money)

}

// reading object and coast
// add the lolcalstorage.username as seller
// whogavemax is nil, TimeToOut = 3600, IsActive True, dataofset = time.Now()
func (n *Aukt) AddDataToTable() {

}

// cals every minite if timeofset + timeToOut > Max -> set isActive to false
// good for small database
func (n *Aukt) TimeOut() {
	fmt.Println("Tick")
	var current_rows []Aukt_Table_Row
	db.GetConection().Model(&Aukt_Table_Row{}).Select("*").Where("isactive = ?", 1).Find(&current_rows)
	fmt.Println(current_rows)
	param := 1000000000 //get seconds from micro
	for _, elem := range current_rows {

		dateNow := time.Now()

		//calculate sub
		if dateNow.Sub(elem.DataOfSet) > time.Duration(elem.TimeToOut*param) {
			elem.IsActive = false
			db.GetConection().Model(Aukt_Table_Row{Number: elem.Number}).Update("isactive", elem.IsActive)
		}
	}

}
