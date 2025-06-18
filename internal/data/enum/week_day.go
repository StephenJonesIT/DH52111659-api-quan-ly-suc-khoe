package enum

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type WeekDay string

const (
	Monday    WeekDay = "monday"
	Tuesday   WeekDay = "tuesday"
	Wednesday WeekDay = "wednesday"
	Thursday  WeekDay = "thursday"
	Friday    WeekDay = "friday"
	Saturday  WeekDay = "saturday"
	Sunday    WeekDay = "sunday"
)

var weekDayToString = map[WeekDay]string{
	Monday:    "monday",
	Tuesday:   "tuesday",
	Wednesday: "wednesday",
	Thursday:  "thursday",
	Friday:    "friday",
	Saturday:  "saturday",
	Sunday:    "sunday",
}

var stringToWeekDay = map[string]WeekDay{
	"monday":    Monday,
	"tuesday":   Tuesday,
	"wednesday": Wednesday,
	"thursday":  Thursday,
	"friday":    Friday,
	"saturday":  Saturday,
	"sunday":    Sunday,
}

func (w WeekDay) String() string {
	if str, ok := weekDayToString[w]; ok {
		return str
	}
	return "unknown"
}

func ParseStr2WeekDay(s string) (WeekDay, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if w, ok := stringToWeekDay[s]; ok {
		return w, nil
	}
	return "", fmt.Errorf("invalid weekday: %s", s)
}

func (w *WeekDay) Scan(value interface{}) error {
	if value == nil {
		*w = Monday
		return nil
	}
	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("invalid type for WeekDay: %T", value)
	}
	v, err := ParseStr2WeekDay(s)
	if err != nil {
		return err
	}
	*w = v
	return nil
}

func (w WeekDay) Value() (driver.Value, error) {
	return w.String(), nil
}

func (w WeekDay) MarshalJSON() ([]byte, error) {
	return []byte(`"` + w.String() + `"`), nil
}

func (w *WeekDay) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("invalid JSON data for WeekDay")
	}
	s := string(data[1 : len(data)-1])
	v, err := ParseStr2WeekDay(s)
	if err != nil {
		return err
	}
	*w = v
	return nil
}