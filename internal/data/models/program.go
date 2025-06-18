package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Program struct {
	ProgramID   uuid.UUID    `gorm:"column:program_id;primaryKey" json:"program_id" `
	Title	   string        `gorm:"column:title;not null" json:"title" validate:"required"`
	Description string       `gorm:"column:description;type:text" json:"description,omitempty"`
	Duration    int          `gorm:"column:duration;not null" json:"duration" validate:"required"`
	IsActive	bool         `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt   *time.Time   `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt   *time.Time   `gorm:"column:updated_at" json:"updated_at,omitempty"`
	CreatedBy   uuid.UUID    `gorm:"column:created_by;not null" json:"created_by"`
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
	
	if p.UpdatedAt == nil {
		now := time.Now()
		p.UpdatedAt = &now
	}


	return nil
}

func(p *Program) BeforeUpdate(tx *gorm.DB) error{
	now := time.Now()
	p.UpdatedAt = &now

	return nil
}

type UserProgram struct {
	ID 			int 		`column:"id;primaryKey" json:"id"`
	UserID  	uuid.UUID 	`column:"user_id;not null" json:"user_id"`
	ProgramID 	uuid.UUID 	`column:"program_id;not null" json:"program_id"`
	LevelID	 	uuid.UUID 	`column:"level_id;not null" json:"level_id"`
	StartDate 	*time.Time 	`column:"start_date" json:"start_date"`
	TotalPoints int			`column:"total_points" json:"total_points"`
	CreatedAt 	*time.Time	`column:"created_at" json:"created_at"`
	UpdatedAt 	*time.Time	`column:"updated_at;not null" json:"updated_at"`
}

func(UserProgram) TableName() string{
	return "user_programs"
}