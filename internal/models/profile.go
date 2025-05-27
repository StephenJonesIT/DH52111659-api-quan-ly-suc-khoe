package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserID 		uuid.UUID 	`json:"user_id" gorm:"column:user_id;primaryKey" validate:"required"`
	FullName 	string 		`json:"full_name" gorm:"column:full_name;not null" validate:"required"`
	DayOfBirth	*time.Time	`json:"day_of_birth" gorm:"column:day_of_birth;not null" validate:"required"`
	Gender 		bool 		`json:"gender" gorm:"column:gender;not null;default:true"`
	AvatarURL	string		`json:"avatar_url,omitempty" gorm:"column:avatar_url"`
	CreatedAt	*time.Time	`json:"created_at,omitempty" gorm:"column:created_at"`
}

func(Profile) TableName() string {
	return "profiles"
}
