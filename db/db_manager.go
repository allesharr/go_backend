package db

import (
	"fmt"
	"go_backend/api"
	"go_backend/logger"
	"go_backend/prop_manager"

	"github.com/go-xorm/xorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB
var engine *xorm.Engine

func InitDB(props api.PropertiesStruct) bool {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Europe%%2FMoscow&charset=utf8", props.DB_USERNAME, props.DB_PASSWORD, props.DB_HOST, props.DB_PORT, props.DB_DATABASE_NAME)

	gormConfig := gorm.Config{}
	if props.APP_DEBUG_MODE {
		gormConfig.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}

	var err error
	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: connectionString,
	}), &gormConfig)

	if err != nil {
		logger.Logger{}.Log("ERR", err.Error())
		return false
	}

	return true
}

func GetConection() *gorm.DB {
	return db
}

func XormEnabled() bool {
	// PROPERTIES ARE READED HERE BECUSE U NEED MORE BASES IN FUTURE. THERE WILL BE MUCH MORE PROPS
	information := prop_manager.AppProperties{}
	err := information.Read()
	props := information.Props
	if err != nil {
		logger.Logger{}.Log(logger.ERR, "Cannot init properties from properties.json. Exit...")
		return false
	}
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Europe%%2FMoscow&charset=utf8", props.DB_USERNAME, props.DB_PASSWORD, props.DB_HOST, props.DB_PORT, props.DB_DATABASE_NAME)

	engine, err = xorm.NewEngine("mysql", connectionString)

	if err != nil {
		logger.Logger{}.Log(logger.ERR, "Cannot init database..")
		fmt.Println(err)
		return false
	}
	return true
}

func GetXORM() *xorm.Engine {
	return engine
}
