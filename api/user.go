package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"gorm.io/gorm"
)

type User struct {
	ID          uint   `json:"id"`
	SkudId      int64  `gorm:"column:skud_id" json:"skud_id"`
	TableNumber int64  `gorm:"column:table_number" json:"table_number"`
	FirstName   string `gorm:"column:first_name" json:"first_name"`
	MiddleName  string `gorm:"column:middle_name" json:"middle_name"`
	Surname     string `gorm:"column:surname" json:"surname"`
	OrgName     string `gorm:"column:org_name" json:"org_name"`
}

func (User) TableName() string {
	return "users"
}

func (User) GetOrderExpression(queryMap *url.Values) string {
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

func (User) SetFilterExpression(queryParams *url.Values, query *gorm.DB) {
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
		case "id", "skud_id", "table_number":
			query.Where(fmt.Sprintf(`%s = ?`, filterKey), filterValue)
			continue

		case "first_name", "middle_name", "surname":
			query.Where(fmt.Sprintf(`%s LIKE ?`, filterKey), fmt.Sprintf(`%%%s%%`, filterValue))
			continue

		case "org_name":
			query.Where(`org_name = ?`, filterValue)
			continue

		default:
			return
		}

	}
}
