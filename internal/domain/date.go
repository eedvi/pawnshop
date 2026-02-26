package domain

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Date represents a date without time component
// Serializes as YYYY-MM-DD in JSON (no timezone conversion)
type Date struct {
	time.Time
}

// DateFormat is the format used for date serialization
const DateFormat = "2006-01-02"

// NewDate creates a new Date from year, month, day
func NewDate(year int, month time.Month, day int) Date {
	return Date{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

// DateFromTime creates a Date from a time.Time (strips time component)
func DateFromTime(t time.Time) Date {
	return Date{time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)}
}

// Today returns the current date
func Today() Date {
	return DateFromTime(time.Now())
}

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(DateFormat, s)
	if err != nil {
		return Date{}, err
	}
	return Date{t}, nil
}

// MarshalJSON implements json.Marshaler
// Returns date in YYYY-MM-DD format without timezone
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Format(DateFormat))), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (d *Date) UnmarshalJSON(data []byte) error {
	// Remove quotes
	str := string(data)
	if str == "null" || str == `""` {
		d.Time = time.Time{}
		return nil
	}

	// Remove quotes if present
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	t, err := time.Parse(DateFormat, str)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// Scan implements sql.Scanner
func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		t, err := time.Parse(DateFormat, string(v))
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	case string:
		t, err := time.Parse(DateFormat, v)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	}

	return fmt.Errorf("cannot scan %T into Date", value)
}

// Value implements driver.Valuer
func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

// String returns the date in YYYY-MM-DD format
func (d Date) String() string {
	if d.IsZero() {
		return ""
	}
	return d.Format(DateFormat)
}
