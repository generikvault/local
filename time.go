package local

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Time struct {
	date time.Time
}

// NewTime returns a new date.Time.
func NewTime(year, month, day, hour, min, sec int) Time {
	return Time{
		date: time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC),
	}
}

// NowLocal returns a new local.Date with the current date.
func NowTime() Time {
	now := time.Now()
	return NewTime(now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second())
}

// TimeFromString returns a new local.Time from a string.
func TimeFromString(s string) (Time, error) {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return Time{}, err
	}
	return NewTime(t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second()), nil
}

// GormDataType returns gorm common data type. This type is used for the field's column type.
func (Time) GormDataType() string {
	return "time"
}

// GormDBDataType returns gorm DB data type based on the current using database.
func (Time) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "TIME"
	case "postgres":
		return "TIME"
	case "sqlserver":
		return "TIME"
	case "sqlite":
		return "TEXT"
	default:
		return ""
	}
}

// Scan implements sql.Scanner interface and scans value into Time,
func (t *Time) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return t.setFromString(string(v))
	case string:
		return t.setFromString(v)
	case time.Time:
		t.setFromTime(v)
	default:
		return fmt.Errorf("failed to scan value: %v", v)
	}

	return nil
}

func (t *Time) setFromString(str string) error {
	// hacked fix for https://github.com/go-gorm/datatypes/issues/279
	if len(str) > 19 {
		str = str[:19]
	}
	time, err := time.Parse(time.DateTime, str)
	if err != nil {
		return err
	}
	t.date = time
	return nil
}

func (t *Time) setFromTime(src time.Time) {
	*t = NewTime(src.Year(), int(src.Month()), src.Day(), src.Hour(), src.Minute(), src.Second())
}

// Value implements driver.Valuer interface and returns string format of Time.
func (t Time) Value() (driver.Value, error) {
	return t.date.Format(time.DateTime), nil
}

// Before returns true if l is before b.
func (l Time) Before(b Time) bool {
	return l.date.Before(b.date)
}

// After returns true if l is after b.
func (l Time) After(b Time) bool {
	return l.date.After(b.date)
}

// AddDays returns a new local.Date with the specified number of days added to l.
func (l Time) AddDays(days int) Time {
	return Time{
		date: l.date.AddDate(0, 0, days),
	}
}

// Add returns a new local.Date with the specified duration added to l.
func (l Time) Add(d time.Duration) Time {
	return Time{
		date: l.date.Add(d),
	}
}

// Sub returns the duration between l and b.
func (l Time) Sub(b Time) time.Duration {
	return l.date.Sub(b.date)
}

// MarshalJSON implements json.Marshaler to convert Time to json serialization.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements json.Unmarshaler to deserialize json data.
func (t *Time) UnmarshalJSON(data []byte) error {
	// ignore null
	if string(data) == "null" {
		return nil
	}
	return t.setFromString(strings.Trim(string(data), `"`))
}

// String implements fmt.Stringer interface.
func (t Time) String() string {
	return t.date.Format("02.01.2006 15:04:05")
}

// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (t Time) IsZero() bool {
	return t.date.IsZero()
}
