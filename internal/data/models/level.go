package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Level struct {
	LevelID 		uuid.UUID   `gorm:"column:level_id;primaryKey" json:"level_id"`
	ProgramID		uuid.UUID   `gorm:"column:program_id;not null" json:"program_id"`
	Name 			string      `gorm:"column:name;not null" json:"name" validate:"required"`
	Description 	*string     `gorm:"column:description;type:text" json:"description,omitempty"`
	PointRequire 	int  	 	`gorm:"column:point_require;not null" json:"point_require" validate:"required"`
	IsActive		bool        `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt       *time.Time  `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt       *time.Time  `gorm:"column:update_at" json:"updated_at,omitempty"`
}

func (Level) TableName() string {
	return "levels"
}


func (l *Level) BeforeCreate(tx *gorm.DB) error {
	if l.LevelID == uuid.Nil {
		l.LevelID = uuid.New()
	}

	if l.CreatedAt == nil {
		now := time.Now()
		l.CreatedAt = &now
	}

	if l.UpdatedAt == nil {
		now := time.Now()
		l.UpdatedAt = &now
	}

	return nil
}