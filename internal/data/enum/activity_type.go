package enum

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type ActivityType string

const (
	TypeActivity  ActivityType = "Activity"
	TypeMiniGame  ActivityType = "MiniGame"
	TypeChallenge ActivityType = "Challenge"
)

var allTypeActivity = map[ActivityType]string{
	TypeActivity:  "Activity",
	TypeMiniGame:  "MiniGame",
	TypeChallenge: "Challenge",
}

var stringToActivityType = map[string]ActivityType{
	"Activity":  TypeActivity,
	"MiniGame":  TypeMiniGame,
	"Challenge": TypeChallenge,
}

func (t ActivityType) String() string {
	if name, ok := allTypeActivity[t]; ok {
		return name
	}
	return "Unknown"
}

func ParseStr2ActivityType(s string) (ActivityType, error) {
	s = strings.TrimSpace(s)
	if t, ok := stringToActivityType[s]; ok {
		return t, nil
	}
	return "", fmt.Errorf("invalid activity type: %s", s)
}

func (t *ActivityType) Scan(value interface{}) error {
	if value == nil {
		*t = TypeActivity
		return nil
	}
	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("invalid type for ActivityType: %T", value)
	}
	v, err := ParseStr2ActivityType(s)
	if err != nil {
		return err
	}
	*t = v
	return nil
}

func (t ActivityType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

func (t *ActivityType) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("invalid JSON data for ActivityType")
	}
	s := string(data[1 : len(data)-1])
	v, err := ParseStr2ActivityType(s)
	if err != nil {
		return err
	}
	*t = v
	return nil
}

func (t ActivityType) Value() (driver.Value, error) {
	return t.String(), nil
}