package routers

import (
	"encoding/json"
	"fmt"
	"go_backend/logger"
	"go_backend/prop_manager"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type User_Event_Per_Day struct {
	Event_Time  time.Time `gorm:"column:event_time"`
	Event_Name  string    `gorm:"column:event_name"`
	Person_Name string    `gorm:"column:fullname"`
}
type LaterPerson struct {
	Name             string
	IsMorningLate    bool
	IsEveningEarly   bool
	IsLateFromDinner bool
	IsSmokingTooMuch bool
}

// var employers = make(map[string][]User_Event_Per_Day)
var names_of_employers []string
var current_date string = "2023-02-02"

// var Start_Day_Time time.Time
// var End_Day_Time time.Time
var err error

// Using xorm agains gorm
func TestSQL(t *testing.T) {

	Start_Day_Time, err := time.Parse(time.DateTime, current_date+" 8:30:05")
	End_Day_Time, er := time.Parse(time.DateTime, current_date+" 17:30:05")
	if err != nil {
		fmt.Println(err)
	}
	if er != nil {
		fmt.Println(er)
	}
	if err != nil {
		fmt.Println("Time of start and end had not parsed")
	}

	logger.Logger{}.Log(logger.INFO, "Starting application...")

	// PROPERTIES
	information := prop_manager.AppProperties{}
	err = information.Read()
	props := information.Props
	if err != nil {
		logger.Logger{}.Log(logger.ERR, "Cannot init properties from properties.json. Exit...")
		os.Exit(1)
	}
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Europe%%2FMoscow&charset=utf8", props.DB_USERNAME, props.DB_PASSWORD, props.DB_HOST, props.DB_PORT, props.DB_DATABASE_NAME)

	engine, err := xorm.NewEngine("mysql", connectionString)

	if err != nil {
		logger.Logger{}.Log(logger.ERR, "Cannot init database..")
		fmt.Println(err)
	}

	data, err := engine.Query("select concat(first_name,' ', middle_name, ' ', surname) as fullname  from users	where org_name like '%ГУП%'")
	if err != nil {
		fmt.Println("Error in request to database: ", err)
	}

	for _, elem := range data {
		names_of_employers = append(names_of_employers, string(elem["fullname"]))
	}
	defer engine.Close()
	//_______________________________________________________________________
	var laters []LaterPerson
	for _, name := range names_of_employers {

		data_events := Find_All_Events_In_Day_Por_Name(engine, current_date, name)
		later := Morning_OR_Evining_Late(data_events, name, Start_Day_Time, End_Day_Time)

		if later.IsEveningEarly || later.IsMorningLate || later.IsLateFromDinner {
			laters = append(laters, later)
		}
	}

	laters_json, err := json.Marshal(laters)

	fmt.Println(laters_json)

}

func Find_All_Events_In_Day_Por_Name(x *xorm.Engine, day string, name string) []User_Event_Per_Day {
	queryStr := fmt.Sprintf("select event_time,	event_name , concat(users.first_name, ' ', users.middle_name, ' ', users.surname) fullname from	skud_events se left join users on user_id = users.skud_id where 	users.org_name like '%%Электронный регион%%' and event_type = 32 and event_time like '%%%s%%'	and concat(users.first_name, ' ', users.middle_name, ' ', users.surname) = '%s'", day, name)
	data, err := x.Query(queryStr)
	if err != nil {
		fmt.Println("Error while fetching data", err)
	}
	eves := []User_Event_Per_Day{}
	for _, elem := range data {
		if string(elem["fullname"]) != "" {

			this_time, err := time.Parse(time.RFC3339, string(elem["event_time"]))
			if err == nil {
				res := User_Event_Per_Day{
					Person_Name: string(elem["fullname"]),
					Event_Name:  string(elem["event_name"]),
					Event_Time:  this_time,
				}
				eves = append(eves, res)
			}
		}
	}
	return eves

}

func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}

func Morning_OR_Evining_Late(data_events []User_Event_Per_Day, name string, Start_Day_Time, End_Day time.Time) LaterPerson {
	lp := LaterPerson{
		Name:             name,
		IsMorningLate:    false,
		IsEveningEarly:   false,
		IsLateFromDinner: false,
	}

	//work without dinner
	for index, elem := range data_events {
		//первый вход
		if strings.Contains(elem.Event_Name, "Вход") && index == 0 && len(data_events) == 1 {
			if elem.Event_Time.Sub(Start_Day_Time).Minutes() > 0 {
				lp.IsMorningLate = true
			}
		}
		//dinner or smoke
		if strings.Contains(elem.Event_Name, "Вход") && index > 0 {
			if strings.Contains(data_events[index-1].Event_Name, "Выход") {
				timeDinnerStart, err := time.Parse(time.DateTime, current_date+" 11:45:05")
				if err != nil {
					fmt.Println(err)
				}
				timeDinnerFinish, err := time.Parse(time.DateTime, current_date+" 14:00:05")
				if err != nil {
					fmt.Println(err)
				}

				if !inTimeSpan(timeDinnerStart, timeDinnerFinish, data_events[index-1].Event_Time) {
					timeOfLate := data_events[index].Event_Time.Sub(data_events[index-1].Event_Time)
					if timeOfLate.Minutes() > 48 {
						lp.IsLateFromDinner = true
					}
				} else {
					timeOfLate := data_events[index].Event_Time.Sub(data_events[index-1].Event_Time)
					if timeOfLate.Minutes() > 15 {
						lp.IsSmokingTooMuch = true
					}
				}
			}

		}

		//out
		if index == len(data_events)-1 && len(data_events) > 1 {
			if End_Day.Sub(elem.Event_Time).Minutes() > 0 {
				lp.IsEveningEarly = false
			} else {
				lp.IsEveningEarly = true
			}

		}

	}

	return lp
}
