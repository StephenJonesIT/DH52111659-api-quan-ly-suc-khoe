package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	UserID 		uuid.UUID 	`json:"account_id" gorm:"column:account_id;not null" validate:"required"`
	FullName 	string 		`json:"full_name" gorm:"column:full_name;not null" validate:"required"`
	DayOfBirth	int			`json:"year_old" gorm:"column:year_old;not null" validate:"required"`
	Gender 		bool 		`json:"gender" gorm:"column:gender;not null;default:true"`
	Weight 		int		 	`json:"weight" gorm:"column:weight;not null" validate:"required"`
	Height 		int 		`json:"height" gorm:"column:height;not null" validate:"required"`
	AvatarURL	string		`json:"avatar_url,omitempty" gorm:"column:avatar_url"`
	HealthGoal  *string		`json:"health_goal,omitempty" gorm:"column:health_goal"`
}

func(Profile) TableName() string {
	return "users"
}

type ProfileResponse struct {
	ProfileID 	uuid.UUID	`json:"profile_id,omitempty"`
	UserID 		uuid.UUID 	`json:"account_id"`
	FullName 	string 		`json:"full_name"`
	DayOfBirth	int			`json:"year_old"`
	Gender 		bool 		`json:"gender"`
	Weight 		int		 	`json:"weight"`
	Height 		int 		`json:"height"`
	BMI 		float64		`json:"bmi,omitempty"`
	AvatarURL	string		`json:"avatar_url,omitempty"`
	CreatedAt	*time.Time	`json:"created_at,omitempty"`
} 