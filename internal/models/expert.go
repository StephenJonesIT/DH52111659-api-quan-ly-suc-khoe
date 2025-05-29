package models

import (
	"time"

	"github.com/google/uuid"
)

type Expert struct {
	ExpertID    int    `json:"expert_id" gorm:"column:expert_id;primaryKey"`
	FullName    string `json:"full_name" gorm:"column:full_name;not null"`
	DateOfBirth *time.Time	`json:"date_of_birth" gorm:"column:date_of_birth;not null"`
	ExpertCreate
	AccountID	uuid.UUID	`json:"account_id" gorm:"column:account_id"`
}

type ExpertCreate struct {
	FullName    	string 		`json:"full_name" gorm:"column:full_name;not null" validate:"required"`
	DateOfBirth 	*time.Time	`json:"date_of_birth" gorm:"column:date_of_birth;not null" validate:"required"`
	Gender			bool		`json:"gender" gorm:"column:gender;default:true"`
	TelephoneNumber string		`json:"telephone_number" gorm:"column:telephone_number"`
	Email			string		`json:"email" gorm:"column:email;not null" validate:"required,email"`
	AvatarURL		string		`json:"avatar_url" gorm:"column:avatar_url"`
	Verified		bool		`json:"verified" gorm:"column:verified;default:true"`
	IsDeleted		bool		`json:"is_deleted" gorm:"column:is_deleted"`
}

func(Expert) TableName() string {
	return "experts"
}

func(ExpertCreate) TableName() string {
	return Expert{}.TableName()
}


