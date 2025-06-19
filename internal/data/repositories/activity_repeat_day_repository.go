package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityRepeatDayRepository interface {
	InsertRepeatDay(ctx context.Context, tx *gorm.DB, repeatItem *[]models.ActivityRepeatDay) error
	UpdateRepeatDayByID(ctx context.Context, tx *gorm.DB, cond uuid.UUID, repeatItems *[]models.ActivityRepeatDay) error
	DeleteRepeatDayByID(ctx context.Context, tx *gorm.DB, cond uuid.UUID) error
}

type activityRepeatDayRepository struct {
	DB *gorm.DB
}

func NewActivityRepeatDayRepositor(db *gorm.DB) ActivityRepeatDayRepository {
	return &activityRepeatDayRepository{DB: db}
}

func (repo *activityRepeatDayRepository) InsertRepeatDay(ctx context.Context, tx *gorm.DB, repeatItem *[]models.ActivityRepeatDay) error {
	db := tx

	if db == nil {
		db = repo.DB
	}

	return db.WithContext(ctx).
		Table(models.ActivityRepeatDay{}.TableName()).
		Create(&repeatItem).Error
}

func (r *activityRepeatDayRepository) UpdateRepeatDayByID(
	ctx context.Context,
	tx *gorm.DB,
	cond uuid.UUID,
	repeatItems *[]models.ActivityRepeatDay,
) error {
	DB := tx

	if DB == nil {
		DB = r.DB
	}

	query := DB.WithContext(ctx).Table(models.ActivityRepeatDay{}.TableName())

	if err := query.Where("activity_id = ?", cond).Delete(models.ActivityRepeatDay{}).Error; err != nil {
		return err
	}

	return query.Create(repeatItems).Error
}

func(r *activityRepeatDayRepository) DeleteRepeatDayByID(ctx context.Context, tx *gorm.DB, cond uuid.UUID) error{
	DB := tx

	if DB == nil {
		DB = r.DB
	}

	query := DB.WithContext(ctx).Table(models.ActivityRepeatDay{}.TableName())

	return query.Where("activity_id = ?", cond).Delete(models.ActivityRepeatDay{}).Error
}