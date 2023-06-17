package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type SysUser struct {
	ID                  uint        `gorm:"primaryKey" json:"id"`
	Login               string      `gorm:"column:login" json:"login"`
	PasswordHash        string      `gorm:"column:password_hash" json:"password_hash"`
	SessionKey          null.String `gorm:"column:session_key" json:"session_key"`
	SessionCreationDate *time.Time  `gorm:"column:session_creation_date" json:"session_creation_date"`
	IsAdmin             bool        `gorm:"column:is_admin" json:"is_admin"`
}

func (SysUser) TableName() string {
	return "system_users"
}

func (SysUser) GetOrderExpression(queryMap *url.Values) string {
	orderField := queryMap.Get("order_field")
	if orderField == "" {
		return ""
	}

	filterString := ""
	switch orderField {

	default:
		filterString = fmt.Sprintf("%s %s", orderField, queryMap.Get("order_op"))

	}

	return filterString
}

func (SysUser) SetFilterExpression(queryParams *url.Values, query *gorm.DB) {
	filterJson := queryParams.Get("filter")

	filter := QueryFilter{}
	filter.X = make(map[string]interface{})

	if filterJson != "" {
		json.Unmarshal([]byte(filterJson), &filter.X)
	}

	if len(filter.X) == 0 {
		return
	}

	for filterKey, filterValue := range filter.X {
		switch filterKey {
		case "id":
			query.Where(fmt.Sprintf(`%s = ?`, filterKey), filterValue)
			continue

		case "login":
			query.Where(fmt.Sprintf(`%s LIKE ?`, filterKey), fmt.Sprintf(`%%%s%%`, filterValue))
			continue

		default:
			return
		}

	}
}
