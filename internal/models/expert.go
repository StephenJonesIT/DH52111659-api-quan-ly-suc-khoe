package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expert struct {
	ExpertID    int    `json:"expert_id" gorm:"column:expert_id;primaryKey"`
	FullName    string `json:"full_name" gorm:"column:full_name;not null"`
	DateOfBirth *time.Time	`json:"date_of_birth" gorm:"column:date_of_birth;not null"`
	ExpertRequest
	
}

type ExpertRequest struct {
	FullName    	string 		`json:"full_name" gorm:"column:full_name;not null" validate:"required"`
	DateOfBirth 	*time.Time	`json:"date_of_birth" gorm:"column:date_of_birth;not null" validate:"required"`
	Gender			bool		`json:"gender" gorm:"column:gender;default:true"`
	TelephoneNumber string		`json:"telephone_number" gorm:"column:telephone_number"`
	Email			string		`json:"email" gorm:"column:email;not null" validate:"required,email"`
	AvatarURL		string		`json:"avatar_url" gorm:"column:avatar_url"`
	Verified		bool		`json:"verified" gorm:"column:verified;default:true"`
	IsDeleted		bool		`json:"is_deleted" gorm:"column:is_deleted"`
	AccountID		uuid.UUID	`json:"account_id,omitempty" gorm:"column:account_id"`
}

func(Expert) TableName() string {
	return "experts"
}

func(ExpertRequest) TableName() string {
	return Expert{}.TableName()
}

func(e *ExpertRequest) BeforeCreate(tx *gorm.DB) error {
	if e.AccountID == uuid.Nil {
		e.AccountID = uuid.New()
	}
	return nil
}
