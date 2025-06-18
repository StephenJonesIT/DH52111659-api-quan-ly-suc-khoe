package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GoalRepository interface {
	GetAllGoals(ctx context.Context, paging *common.Paging) ([]*models.Goal, error)
	GetGoalsByProgramID(ctx context.Context, cond uuid.UUID, paging *common.Paging) ([]*models.Goal, error)
}

type goalRepositoryImpl struct {
	DB *gorm.DB
}

func NewGoalRepository(db *gorm.DB) GoalRepository {
	return &goalRepositoryImpl{
		DB: db,
	}
}

func (r *goalRepositoryImpl) GetAllGoals(ctx context.Context, paging *common.Paging) ([]*models.Goal, error) {
	var goals []*models.Goal

	query := r.DB.WithContext(ctx).Table(models.Goal{}.TableName())

	if err := query.Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if err := query.Offset((paging.Page-1)*paging.Limit).Limit(paging.Limit).Find(&goals).Error; err != nil {
		return nil, err
	}

	return goals, nil
}

func(r *goalRepositoryImpl) GetGoalsByProgramID(ctx context.Context, cond uuid.UUID, paging *common.Paging) ([]*models.Goal, error){
	return nil, nil
}