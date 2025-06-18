package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"gorm.io/gorm"
)

type ActivityRepeatDayRepository interface {
	InsertRepeatDay(ctx context.Context, tx *gorm.DB, repeatItem *[]models.ActivityRepeatDay) error
}

type activityRepeatDayRepository struct {
	DB *gorm.DB
}

func NewActivityRepeatDayRepositor(db *gorm.DB) ActivityRepeatDayRepository{
	return &activityRepeatDayRepository{DB: db}
}

func(repo activityRepeatDayRepository) InsertRepeatDay(ctx context.Context, tx *gorm.DB, repeatItem *[]models.ActivityRepeatDay) error{
	db := tx
	
	if db == nil {
		db = repo.DB
	}

	return db.WithContext(ctx).
		Table(models.ActivityRepeatDay{}.TableName()).
		Create(&repeatItem).Error
}