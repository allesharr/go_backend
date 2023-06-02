package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type ShortTime time.Time

func (mt *ShortTime) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return err
	}
	*mt = ShortTime(t)
	return nil
}

func (mt *ShortTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*mt).Format("2006-01-02"))
}

func (date *ShortTime) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = ShortTime(nullTime.Time)
	return
}

func (date ShortTime) Value() (driver.Value, error) {
	y, m, d := time.Time(date).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Time(date).Location()), nil
}
