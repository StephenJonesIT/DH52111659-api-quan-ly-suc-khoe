package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expert struct {
	ExpertRequest
}

type ExpertRequest struct {
	ExpertID    uuid.UUID    `json:"expert_id" gorm:"column:expert_id;primaryKey"`
	FullName    	string 		`json:"full_name" gorm:"column:full_name;not null" validate:"required"`
	DateOfBirth 	*time.Time	`json:"date_of_birth" gorm:"column:date_of_birth;not null" validate:"required"`
	Gender			bool		`json:"gender" gorm:"column:gender;default:true"`
	TelephoneNumber string		`json:"telephone_number" gorm:"column:telephone_number"`
	AvatarURL		string		`json:"avatar_url,omitempty" gorm:"column:avatar_url"`
	ExpertType		string		`json:"expert_type" gorm:"column:expert_type;not null" validate:"required"`
	AccountID		uuid.UUID	`json:"account_id,omitempty" gorm:"column:account_id"`
}

func(Expert) TableName() string {
	return "experts"
}

func(ExpertRequest) TableName() string {
	return Expert{}.TableName()
}

func(e *ExpertRequest) BeforeCreate(tx *gorm.DB) error {
	if e.ExpertID == uuid.Nil{
		e.ExpertID = uuid.New()
	}
	
	if e.AccountID == uuid.Nil {
		e.AccountID = uuid.New()
	}
	return nil
}
