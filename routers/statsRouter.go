package routers

import (
	"go_backend/stats"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatsRouter struct {
	Gin          *gin.Engine
	StatsManager *stats.StatManager
}

func (n *StatsRouter) Init() {
	n.Gin.GET("/stats", n.Get)
}

func (n *StatsRouter) Get(c *gin.Context) {
	c.JSON(http.StatusOK, n.StatsManager.StatList)
}
