package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Program struct {
	ProgramID    uuid.UUID  `json:"program_id,omitempty" gorm:"primaryKey"`
	ProgramName  string   	`json:"name" gorm:"column:program_name;not null" validate:"required"`
	Description  string   	`json:"description" gorm:"column:description;not null" validate:"required"`
	TotalDays    int      	`json:"total_day" gorm:"column:total_days;not null"`
	DurationType string   	`json:"duration_type" gorm:"column:duration_type;not null" validate:"required"`
	ExpertID     int      	`json:"expert_id,omitempty" gorm:"column:expert_id;not null"`
	CreatedAt    *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
}

func(Program) TableName() string{
	return "programs"
}

func (p *Program) BeforeCreate(tx *gorm.DB) error {
	if p.ProgramID == uuid.Nil {
		p.ProgramID = uuid.New()
	}
	if p.CreatedAt == nil {
		now := time.Now()
		p.CreatedAt = &now
	}
	return nil
}
