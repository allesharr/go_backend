package main

import (
	"fmt"
	"go_backend/db"
	"go_backend/logger"
	"go_backend/prop_manager"
	"go_backend/routers"
	"go_backend/stats"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Logger{}.Log(logger.INFO, "Starting application...")

	// PROPERTIES
	props := prop_manager.AppProperties{}
	err := props.Read()
	if err != nil {
		logger.Logger{}.Log(logger.ERR, "Cannot init properties from properties.json. Exit...")
		os.Exit(1)
	}

	// DATABASE
	if !db.InitDB(props.Props) || !db.XormEnabled() {
		logger.Logger{}.Log(logger.ERR, "Cannot init database. Exit...")
		os.Exit(1)
	}

	// REST
	var router *gin.Engine
	if !props.Props.APP_DEBUG_MODE {
		gin.SetMode(gin.ReleaseMode)

		router = gin.New()
		router.Use(gin.Recovery())
	} else {
		router = gin.Default()
	}

	router.Use(CORSMiddleware())
	initRouters(router, &props)

	router.Run(fmt.Sprintf(":%d", props.Props.APP_LISTEN_PORT))
}

func initRouters(router *gin.Engine, props *prop_manager.AppProperties) {
	statManager := stats.StatManager{}
	statManager.Init()

	authRouter := routers.AuthRouter{Gin: router, Props: props, StatsManager: &statManager}
	router.Use(authRouter.AuthMiddleware)
	authRouter.Init()

	eventRouter := routers.EventRouter{Gin: router, StatsManager: &statManager}
	eventRouter.Init()

	userRouter := routers.UserRouter{Gin: router, StatsManager: &statManager}
	userRouter.Init()

	AuktRouter := routers.Aukt{Gin: router, StatsManager: &statManager}
	AuktRouter.Init()

	sysUserRouter := routers.SysUserRouter{Gin: router, StatsManager: &statManager}
	sysUserRouter.Init()

	reportsRouter := routers.ReportRouter{Gin: router, StatsManager: &statManager}
	reportsRouter.Init()

	statsRouter := routers.StatsRouter{Gin: router, StatsManager: &statManager}
	statsRouter.Init()

	Lateters_router := routers.LatersRouter{Gin: router, Props: props, StatsManager: &statManager}
	Lateters_router.Init()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
