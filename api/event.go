package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	//"go_skud_backend/logger"
	"gorm.io/gorm"
)

type Event struct {
	ID        uint   `json:"id"`
	SkudBid   int64  `gorm:"column:skud_bid" json:"skud_bid"`
	DoorIndex int64  `gorm:"column:door_index" json:"door_index"`
	EventName string `gorm:"column:event_name" json:"event_name"`
	EventType int64  `gorm:"column:event_type" json:"event_type"`
	// EventTime *ShortTime `gorm:"type:time;column:event_time" json:"event_time"`
	EventTime *time.Time `gorm:"type:time;column:event_time" json:"event_time"`

	UserId  int    `gorm:"column:user_id" json:"user_id"`
	User    User   `gorm:"references:SkudId" json:"user"`
	OrgName string `gorm:"-" json:"org_name"`
}

func (Event) TableName() string {
	return "skud_events"
}

func (Event) GetOrderExpression(queryMap *url.Values) (string, string) {
	orderField := queryMap.Get("order_field")
	if orderField == "" {
		return "", ""
	}

	filterString := ""
	switch orderField {
	case "fullname":
		filterString = fmt.Sprintf("User.surname, User.first_name, User.middle_name, event_time %s", queryMap.Get("order_op"))

	case "event_1_date":
		filterString = fmt.Sprintf("event_time %s", queryMap.Get("order_op"))

	case "user":
		filterString = fmt.Sprintf("User.surname %s", queryMap.Get("order_op"))

	default:
		filterString = fmt.Sprintf("%s %s", orderField, queryMap.Get("order_op"))

	}

	return filterString, orderField
}

func (Event) SetFilterExpression(queryParams *url.Values, query *gorm.DB) {
	filterJson := queryParams.Get("filter")

	filter := QueryFilter{}
	filter.X = make(map[string]interface{})

	if filterJson != "" {
		json.Unmarshal([]byte(filterJson), &filter.X)
	}

	networkTypeParam := queryParams.Get("door_index")
	if networkTypeParam != "" {
		filter.X["door_index"] = networkTypeParam
	}

	//logger.Logger{}.Log(logger.INFO, networkTypeParam)

	if len(filter.X) == 0 {
		return
	}

	for filterKey, filterValue := range filter.X {
		switch filterKey {
		case "report_type":
			if filterValue == "employees" {
				query.Where(`User.org_name IN ('ГУП "Электронный регион"', 'ГБУ "Электронный регион"')`)

			} else if filterValue == "third_party_employees" {
				query.Where(`User.org_name NOT IN ('ГУП "Электронный регион"', 'ГБУ "Электронный регион"')`)

			} else if filterValue == "all" {
				//

			} else {
				//

			}
			continue

		case "id":
			query.Where(`skud_events.id = ?`, filterValue)
			continue

		case "event_type", "door_index", "skud_bid":
			query.Where(fmt.Sprintf(`%s = ?`, filterKey), filterValue)
			continue

		case "event_name":
			query.Where(`event_name LIKE ?`, fmt.Sprintf(`%%%s%%`, filterValue))
			continue

		case "user":
			query.Where(`CONCAT(User.surname, ' ', User.first_name, ' ', User.middle_name) LIKE ?`, fmt.Sprintf(`%%%s%%`, filterValue))
			continue

		case "org_name":
			query.Where(`User.org_name LIKE ?`, fmt.Sprintf(`%%%s%%`, filterValue))
			continue

		case "event_time":
			v, ok := filterValue.(map[string]interface{})
			if !ok {
				continue
			}

			query.Where(fmt.Sprintf(`%s BETWEEN ? AND ?`, filterKey), v["date_from"], v["date_to"])
			continue

		default:
			return
		}

	}
}
