package enum

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Status represents the status of an entity (e.g., activity completion)
type Status string

const (
	PendingStatus  Status = "pending"  // Waiting for action
	CompleteStatus Status = "complete" // Successfully completed
	SkipStatus     Status = "skip"     // Skipped by user
)

var allStatus = map[Status]string{
	PendingStatus:  "pending",
	CompleteStatus: "complete",
	SkipStatus:     "skip",
}

var stringToStatus = map[string]Status{
	"pending":  PendingStatus,
	"complete": CompleteStatus,
	"skip":     SkipStatus,
}

// String returns the string representation of Status
func (s Status) String() string {
	if name, ok := allStatus[s]; ok {
		return name
	}
	return "Unknown"
}

// ParseStr2Status converts a string to Status
func ParseStr2Status(str string) (Status, error) {
	str = strings.TrimSpace(strings.ToLower(str))
	if s, ok := stringToStatus[str]; ok {
		return s, nil
	}
	return "", fmt.Errorf("invalid status: %s", str)
}

// Scan implements sql.Scanner for Status
func (s *Status) Scan(value interface{}) error {
	if value == nil {
		*s = PendingStatus // Default to pending if NULL
		return nil
	}
	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return fmt.Errorf("invalid type for Status: %T", value)
	}
	v, err := ParseStr2Status(str)
	if err != nil {
		return err
	}
	*s = v
	return nil
}

// MarshalJSON implements json.Marshaler for Status
func (s Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler for Status
func (s *Status) UnmarshalJSON(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("invalid JSON data for Status")
	}
	str := string(data[1 : len(data)-1])
	v, err := ParseStr2Status(str)
	if err != nil {
		return err
	}
	*s = v
	return nil
}

// Value implements driver.Valuer for Status
func (s Status) Value() (driver.Value, error) {
	return s.String(), nil
}