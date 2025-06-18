package models

import "github.com/google/uuid"

type Goal struct {
	GoalID int    `json:"goal_id" gorm:"column:goal_id;primaryKey"`
	Name   string `json:"name" gorm:"column:name;not null" validate:"required"`
}

func (Goal) TableName() string {
	return "goals"
}

type UserGoal struct {
	ID     int       `json:"id" gorm:"column:id;primaryKey"`
	GoalID int       `json:"goal_id" gorm:"column:goal_id;not null"`
	UserID uuid.UUID `json:"user_id" gorm:"column:user_id;not null"`
}

func (UserGoal) TableName() string {
	return "user_goals"
}

type ProgramGoal struct {
	ID     int       	`json:"id" gorm:"column:id;primaryKey"`
	GoalID int       	`json:"goal_id" gorm:"column:goal_id;not null"`
	ProgramID uuid.UUID `json:"program_id" gorm:"column:program_id;not null"`
}

func (ProgramGoal) TableName() string {
	return "program_goals"
}