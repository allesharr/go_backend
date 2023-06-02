package routers

import (
	"fmt"
	"go_backend/api"
	"go_backend/db"
	"go_backend/logger"
	"go_backend/prop_manager"
	"go_backend/stats"
	"go_backend/utils"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	Gin          *gin.Engine
	Props        *prop_manager.AppProperties
	StatsManager *stats.StatManager
}

func (n *AuthRouter) Init() {
	n.Gin.POST("/login", n.Login)
	n.Gin.POST("/check_auth", n.CheckAuth)
	n.Gin.POST("/register/:data", n.Register_system_user)
}

func (n *AuthRouter) AuthMiddleware(c *gin.Context) {
	urlPath := c.Request.URL.Path
	if urlPath == "/login" || urlPath == "/check_auth" {
		c.Next()
		return
	}

	if !n.Props.Props.APP_AUTH_ENABLED {
		c.Next()
		return
	}

	hasCorrectSession, _ := n.CheckSession(c)
	if !hasCorrectSession {
		c.AbortWithStatusJSON(401, gin.H{"message": "Invalid auth"})
		return
	}

	c.Next()
}

func (n *AuthRouter) CheckAuth(c *gin.Context) {
	var requestObject *api.SysUser

	if err := c.BindJSON(&requestObject); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseMessage{Message: fmt.Sprintf("Ошибка парсинга объекта: %s", err.Error())})
		return
	}

	var sysUserList []api.SysUser
	db.GetConection().Model(&api.SysUser{}).Where("id = ? AND session_key = ?", requestObject.ID, requestObject.SessionKey).Limit(1).Find(&sysUserList)

	if len(sysUserList) == 0 {
		c.AbortWithStatusJSON(401, gin.H{"message": "Invalid auth"})
		return
	}

	c.JSON(http.StatusOK, api.AuthResponse{SessionKey: requestObject.SessionKey.String, UserId: int(sysUserList[0].ID), Login: sysUserList[0].Login})
}

func (n *AuthRouter) CheckSession(c *gin.Context) (bool, *api.SysUser) {
	headerSessionKey := c.Request.Header["SessionKey"]
	login := c.Request.Header["Login"]

	if len(headerSessionKey) == 0 {
		return false, nil
	}

	var userList []api.SysUser
	db.GetConection().Model(&api.SysUser{}).Where("login = ?", login).Find(&userList)

	if len(userList) == 0 || userList[0].SessionKey.IsZero() || userList[0].SessionCreationDate == nil {
		return false, nil
	}

	user := userList[0]
	userSessionKey, err := user.SessionKey.Value()
	if err != nil {
		return false, &user
	}

	if userSessionKey != headerSessionKey[0] {
		return false, &user
	}

	timeDelta := time.Since(*user.SessionCreationDate).Hours()
	return timeDelta <= float64(n.Props.Props.APP_SESSION_TIMEOUT_HOURS), &user
}

func (n *AuthRouter) Login(c *gin.Context) {
	var requestObject *api.SysUser

	if err := c.BindJSON(&requestObject); err != nil {
		c.JSON(http.StatusBadRequest, api.ResponseMessage{Message: fmt.Sprintf("Ошибка парсинга объекта: %s", err.Error())})
		return
	}

	var sysUserList []api.SysUser
	db.GetConection().Model(&api.SysUser{}).Where("login = ? AND password_hash = ?", requestObject.Login, requestObject.PasswordHash).Limit(1).Find(&sysUserList)

	if len(sysUserList) == 0 {
		c.JSON(http.StatusUnauthorized, api.ResponseMessage{Message: "Неверное имя пользователя или пароль"})
		*n.StatsManager.GetStat("login").Errors += 1
		return
	}

	randomString := utils.GetRandomString(30)
	result := db.GetConection().Model(api.SysUser{}).Where("id = ?", sysUserList[0].ID).Updates(map[string]interface{}{"session_key": randomString, "session_creation_date": time.Now()})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, api.ResponseMessage{Message: "Ошибка авторизации"})
		logger.Logger{}.Log(logger.ERR, fmt.Sprintf("Ошибка авторизации: %s", result.Error))
		*n.StatsManager.GetStat("login").Errors += 1
	} else {
		c.JSON(http.StatusOK, api.AuthResponse{SessionKey: randomString, UserId: int(sysUserList[0].ID), Login: sysUserList[0].Login})
		*n.StatsManager.GetStat("login").SuccessAuthRequests += 1
	}

}

func (n *AuthRouter) Register_system_user(c *gin.Context) {

	data := c.Param("data")
	login, hash := strings.Split(data, "-")[0], strings.Split(data, "-")[1]

	fmt.Println(login, hash)
	isLoginCorrect, err := regexp.MatchString("[A-Za-z0-9]+", login)
	if err != nil {
		fmt.Println(err)
	}
	if !isLoginCorrect {
		c.JSON(http.StatusBadRequest, api.ResponseMessage{Message: "Ошибка авторизации, неверный логин и пароль"})
		logger.Logger{}.Log(logger.ERR, fmt.Sprintf("Ошибка авторизации попытка входа с логином : %s", login))

	} else {
		var user api.SysUser
		//____________________________________________________________
		db.GetConection().Where("login = ?", login).First(&user)
		fmt.Println(user)
		if user.Login == "" {
			db.GetConection().Create(&api.SysUser{
				Login:        login,
				PasswordHash: hash,
			})
			c.JSON(http.StatusOK, api.ResponseMessage{Message: "Успешная регистрация пользователя"})

		} else {
			c.JSON(http.StatusBadRequest, api.ResponseMessage{Message: "Пользователь уже существует"})
		}
		fmt.Println(user)

	}
}
