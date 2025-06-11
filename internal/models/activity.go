package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Activity struct {
	ID          		string `json:"id,omitempty" gorm:"column:activity_id;primaryKey;type:uuid"`
	Name        		string `json:"name" gorm:"column:activity_name;not null"`
	Description 		string `json:"description" gorm:"column:description;type:text"`
	Duration    		int    `json:"duration" gorm:"column:duration;not null"`
	Point       		int    `json:"point,omitempty" gorm:"column:point;default:10;not null"`
	ExpertID    		int    `json:"expert_id,omitempty" gorm:"column:expert_id;not null"`
}

func (Activity) TableName() string {
	return "activities"
}

func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}
