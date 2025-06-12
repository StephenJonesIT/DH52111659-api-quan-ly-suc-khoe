package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"context"

	"gorm.io/gorm"
)

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *models.ScheduleCreate) (*models.ScheduleCreate, error)
}

type scheduleRepositoryImpl struct {
	DB *gorm.DB
}

func NewScheduleRepositoryImpl(db *gorm.DB) *scheduleRepositoryImpl {
	return &scheduleRepositoryImpl{
		DB: db,
	}
}

func (r *scheduleRepositoryImpl) Create(
	ctx context.Context,
	schedule *models.ScheduleCreate,
) (*models.ScheduleCreate, error){
	if err := r.DB.WithContext(ctx).Table(models.ScheduleCreate{}.TableName()).Create(schedule).Error; err != nil {
		return nil, err 
	}
	return schedule, nil
}
