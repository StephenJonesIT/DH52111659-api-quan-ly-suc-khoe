package models

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/enum"

	"github.com/google/uuid"
)

type ActivityRepeatDay struct {
	ID         int  			`gorm:"column:id;primaryKey" json:"id"`
	ActivityID uuid.UUID  		`gorm:"column:activity_id;not null" json:"activity_id" validate:"required"`
	Repeat     enum.WeekDay 	`gorm:"column:repeat_day;not null;default:monday" json:"repeat" validate:"required"`
}

func (ActivityRepeatDay) TableName() string {
	return "activity_repeat_days"
}