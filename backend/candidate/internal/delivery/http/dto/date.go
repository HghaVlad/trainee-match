package dto

import (
	"encoding/json"
	"time"
)

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	t, err := time.Parse("02.01.2006", s)
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// MarshalJSON Should be marshaled by value
func (d Date) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return json.Marshal(t.Format("02.01.2006"))
}

func TimeToDate(t time.Time) Date {
	return Date(t)
}
func DateToTime(d Date) time.Time {
	return time.Time(d)
}
