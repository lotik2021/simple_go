package models

import (
	"database/sql/driver"
	"fmt"
	"github.com/go-pg/pg/v9/types"
	"strings"
	"time"
)

type Time struct {
	Time time.Time
}

func NewTime(t time.Time) *Time {
	return &Time{t}
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(time.RFC3339))), nil
}

func (t *Time) UnmarshalJSON(data []byte) error {
	temp, err := time.Parse(time.RFC3339, strings.ReplaceAll(string(data), "\"", ""))
	if err != nil {
		return err
	}

	t.Time = temp

	return nil
}

func (t Time) IsZero() bool {
	return t.Time.IsZero()
}

func (t Time) Add(d time.Duration) Time {
	return Time{t.Time.Add(d)}
}

func (t Time) Sub(u time.Time) time.Duration {
	return t.Time.Sub(u)
}

func (t Time) Before(u time.Time) bool {
	return t.Time.Before(u)
}

func (t Time) After(u time.Time) bool {
	return t.Time.After(u)
}

func (t Time) Format(layout string) string {
	return t.Time.Format(layout)
}

func (t Time) Unix() int64 {
	return t.Time.Unix()
}

func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t *Time) Scan(src interface{}) error {
	if src == nil {
		*t = Time{}
		return nil
	}

	tmp, err := types.ParseTime(src.([]byte))
	if err != nil {
		return err
	}

	*t = Time{tmp}

	return nil
}
