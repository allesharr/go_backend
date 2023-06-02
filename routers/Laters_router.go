package routers

import (
	"go_backend/api"
	"go_backend/db"
	"go_backend/prop_manager"
	"go_backend/stats"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type LatersRouter struct {
	Gin          *gin.Engine
	Props        *prop_manager.AppProperties
	StatsManager *stats.StatManager
}

func (n *LatersRouter) Init() {
	n.Gin.GET("/lates", n.Lates)
}

//Muy IMPORTANTE USAR COMO ESO
// current_date string = "2023-02-02" only this format can be there

func (n *LatersRouter) Lates(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	engine := db.GetXORM()
	WasFlag := false
	current_data := ""
	for key, value := range queryParams {
		if key == "data" {
			current_data = value[0]
			WasFlag = true
		}
	}

	if WasFlag {
		if isDataCorrect(current_data) {
			data := api.CalcLaters(engine, current_data)
			c.JSON(http.StatusOK, data)
		} else {
			c.JSON(http.StatusNotFound, api.ResponseMessage{Message: "Данные введены некорректно используйте формат yyyy-mm-dd"})
		}
	} else {
		c.JSON(http.StatusNotFound, api.ResponseMessage{Message: "Ключ data не найден, не могу искать все среди всех событий в мире"})
	}

}
func isDataCorrect(data string) bool {
	Correct_regexp := regexp.MustCompile(`20\d{2}-(?:(0[0-9])|(1[0-2]))-(?:(0[1-9])|(1[0-9])|(2[0-9])|(3[0-1]))`)
	isCorrent := Correct_regexp.MatchString(data)
	return isCorrent
}
