package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Schedule struct {
}

type ScheduleCreate struct {
	ScheduleID 		uuid.UUID `json:"schedule_id,omitempty" gorm:"primaryKey;type:uuid"`
	ProgramID  		uuid.UUID `json:"program_id,omitempty" gorm:"column:program_id;not null;type:uuid"`
	ActivityID 		uuid.UUID `json:"activity_id,omitempty" gorm:"column:activity_id;not null;type:uuid"`
	WeekNumber  	int       `json:"week_number" gorm:"column:week_number;not null" validate:"required"`
	DayNumber   	int       `json:"day_number" gorm:"column:day_number;not null" validate:"required"`
	RepeatInterval  int       `json:"repeat_interval,omitempty" gorm:"column:repeat_interval;not null;default:0"`
}

func (s Schedule) TableName() string {
	return "schedules"
}

func (s ScheduleCreate) TableName() string {
	return Schedule{}.TableName()
}

func (s *ScheduleCreate) BeforeCreate(tx *gorm.DB) error {
	if s.ScheduleID == uuid.Nil {
		s.ScheduleID = uuid.New()
	}
	if s.RepeatInterval < 0 {
		s.RepeatInterval = 0 // Ensure repeat interval is non-negative
	}
	return nil
}