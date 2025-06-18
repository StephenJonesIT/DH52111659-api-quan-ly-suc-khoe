package models

import "github.com/google/uuid"

type Disease struct {
	DiseaseID int    `json:"disease_id" gorm:"column:disease_id;primaryKey"`
	Name      string `json:"name" gorm:"column:name;not null" validate:"required"`
}

func (Disease) TableName() string {
	return "diseases"
}

type ProgramDiseases struct {
	ID        int       `json:"id" gorm:"column:id;primaryKey"`
	ProgramID uuid.UUID `json:"program_id" gorm:"column:program_id;not null"`
	DiseaseID int 		`json:"disease_id" gorm:"column:disease_id;not null"`
}

func (ProgramDiseases) TableName() string {
	return "program_diseases"
}

type UserDiseases struct{
	ID        int       `json:"id" gorm:"column:id;primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"column:user_id;not null"`
	DiseaseID int       `json:"disease_id" gorm:"column:disease_id;not null"`
}

func (UserDiseases) TableName() string {
	return "user_diseases"
}