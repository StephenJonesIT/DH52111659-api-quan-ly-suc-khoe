package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID 				uuid.UUID 	`json:"id,omitempty" gorm:"column:id;primaryKey"`
	Email 			string		`json:"email" gorm:"column:email;unique;not null" validate:"required,email"`
	Password 		string 		`json:"password" gorm:"column:password_hash;not null" validate:"required,min=8,max=100"`
	Role 			string 		`json:"role,omitempty" gorm:"column:role;default:'user'"`
	CreatedAt 		*time.Time 	`json:"created_at,omitempty" gorm:"column:created_at"`
	IsVerified 		bool 		`json:"is_verified,omitempty" gorm:"column:is_verified;default:false"`
	AccountStatus 	bool 		`json:"account_status,omitempty" gorm:"column:account_status;default:true"`
}

func (Account) TableName() string {
	return "accounts"
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
    if a.ID == uuid.Nil {
        a.ID = uuid.New()
    }
    return nil
}

type AccountCreate struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=admin user"`
	IsVerified bool   `json:"is_verified,omitempty" validate:"omitempty"`
}

