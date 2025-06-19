package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserActivityRepository interface {
	CountParticipants(ctx context.Context, cond uuid.UUID)(int64, error)
}

type userActivityRepository struct{
	DB *gorm.DB
}

func NewUserActivityRepository(db *gorm.DB) UserActivityRepository{
	return &userActivityRepository{DB: db}
}

func(r *userActivityRepository) CountParticipants(ctx context.Context, cond uuid.UUID)(int64, error){
	var count int64

	if err := r.DB.WithContext(ctx).
		Table(models.UserActivity{}.TableName()).
		Where("activity_id = ?", cond).
		Count(&count).Error; err != nil{
			return 0, err
	}

	return count, nil
}