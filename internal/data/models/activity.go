package models

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/enum"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Activity struct {
	ActivityID  uuid.UUID  			`gorm:"column:activity_id;primaryKey" json:"activity_id"`
	LevelID     uuid.UUID  			`gorm:"column:level_id;not null" json:"level_id"`
	Title       string     			`gorm:"column:name;not null" json:"title" validate:"required"`
	Description string     			`gorm:"column:description;type:text" json:"description,omitempty"`
	Duration    int       			`gorm:"column:duration;not null" json:"duration" validate:"required"`
	PointReward int        			`gorm:"column:point_reward;not null" json:"point_reward" validate:"required"`
	Type	    enum.ActivityType 	`gorm:"column:type;default:Activity" json:"type"`	
	CreatedAt   *time.Time 			`gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt   *time.Time 			`gorm:"column:updated_at" json:"updated_at,omitempty"`	
}

func (Activity) TableName() string {
	return "activities"
}

func (a *Activity) BeforeCreate(tx *gorm.DB) error {
	if a.ActivityID == uuid.Nil {
		a.ActivityID = uuid.New()
	}

	if a.CreatedAt == nil {
		now := time.Now()
		a.CreatedAt = &now
	}

	if a.UpdatedAt == nil {
		now := time.Now()
		a.UpdatedAt = &now
	}

	return nil
}
