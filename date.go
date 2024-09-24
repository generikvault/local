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

type Date struct {
	date time.Time
}

// NewDate returns a new date.Date.
func NewDate(year, month, day int) Date {
	return Date{
		date: time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC),
	}
}

// NowDate returns a new date.Date with the current date.
func NowDate() Date {
	now := time.Now()
	return NewDate(now.Year(), int(now.Month()), now.Day())
}

// Equal returns true if l and b are the same date.
func (l Date) Equal(b Date) bool {
	return l.date.Equal(b.date)
}

// Year returns the year specified by l.
func (l Date) Year() int {
	return l.date.Year()
}

// Month returns the month specified by l.
func (l Date) Month() int {
	return int(l.date.Month())
}

// Day returns the day of the month specified by l.
func (l Date) Day() int {
	return l.date.Day()
}

// Quarter returns the first day of the quarter of the year specified by l.
func (l Date) Quarter() Date {
	month := l.Month()
	month -= (month - 1) % 3
	return NewDate(l.Year(), month, 1)
}

// EqualMonth returns true if l and b are the same month and year.
func (l Date) EqualMonth(b Date) bool {
	return l.Year() == b.Year() && l.Month() == b.Month()
}

// EqualQuarter returns true if l and b are the same quarter and year.
func (l Date) EqualQuarter(b Date) bool {
	return l.Year() == b.Year() && (l.Month()-1)/3 == (b.Month()-1)/3
}

// EqualYear returns true if l and b are the same year.
func (l Date) EqualYear(b Date) bool {
	return l.Year() == b.Year()
}

// Before returns true if l is before b.
func (l Date) Before(b Date) bool {
	return l.date.Before(b.date)
}

// After returns true if l is after b.
func (l Date) After(b Date) bool {
	return l.date.After(b.date)
}

// AddDays returns a new date.Date with the specified number of days added to l.
func (l Date) AddDays(days int) Date {
	return Date{
		date: l.date.AddDate(0, 0, days),
	}
}

// AddMonths returns a new date.Date with the specified number of months added to l.
func (l Date) AddMonths(months int) Date {
	return Date{
		date: l.date.AddDate(0, months, 0),
	}
}

// AddYears returns a new date.Date with the specified number of years added to l.
func (l Date) AddYears(years int) Date {
	return Date{
		date: l.date.AddDate(years, 0, 0),
	}
}

// MinusDays returns a new date.Date with the specified number of days subtracted from l.
func (l Date) MinusDays(days int) Date {
	return l.AddDays(-days)
}

// MinusMonths returns a new date.Date with the specified number of months subtracted from l.
func (l Date) MinusMonths(months int) Date {
	return l.AddMonths(-months)
}

// MinusYears returns a new date.Date with the specified number of years subtracted from l.
func (l Date) MinusYears(years int) Date {
	return l.AddYears(-years)
}

// GormDataType returns gorm common data type. This type is used for the field's column type.
func (Date) GormDataType() string {
	return "time"
}

// GormDBDataType returns gorm DB data type based on the current using database.
func (Date) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
func (t *Date) Scan(src interface{}) error {
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

func (t *Date) setFromString(str string) error {
	time, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return err
	}
	t.setFromTime(time)
	return nil
}

func FromString(str string) (Date, error) {
	var t Date
	err := t.setFromString(str)
	return t, err
}

func (t *Date) setFromTime(src time.Time) {
	*t = NewDate(src.Year(), int(src.Month()), src.Day())
}

// Value implements driver.Valuer interface and returns string format of Time.
func (t Date) Value() (driver.Value, error) {
	return t.date.Format(time.DateOnly), nil
}

// String implements fmt.Stringer interface.
func (t Date) String() string {
	return t.date.Format("2.1.2006")
}

// MarshalJSON implements json.Marshaler to convert Time to json serialization.
func (t Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.date.Format(time.DateOnly))
}

// UnmarshalJSON implements json.Unmarshaler to deserialize json data.
func (t *Date) UnmarshalJSON(data []byte) error {
	// ignore null
	if string(data) == "null" {
		return nil
	}
	return t.setFromString(strings.Trim(string(data), `"`))
}
