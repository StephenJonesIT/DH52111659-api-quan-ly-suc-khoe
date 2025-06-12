package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	CreateActivity(ctx context.Context, activity *models.Activity) (*models.Activity, error)
	GetActivities(ctx context.Context, cond map[string]interface{}, paging *common.Paging) ([]*models.Activity, error)
	GetActivityByID(ctx context.Context, cond uuid.UUID) (*models.Activity, error)
}

type ActivityRepositoryImpl struct {
	db *gorm.DB
}

func NewActivityRepositoryImpl(db *gorm.DB) *ActivityRepositoryImpl {
	return &ActivityRepositoryImpl{
		db: db,
	}
}

func (r *ActivityRepositoryImpl) CreateActivity(ctx context.Context, activity *models.Activity) (*models.Activity, error) {
	if err := r.db.WithContext(ctx).
		Table(models.Activity{}.TableName()).
		Create(activity).Error; err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *ActivityRepositoryImpl) GetActivities(
	ctx context.Context,
	cond map[string]interface{},
	paging *common.Paging,
) ([]*models.Activity, error) {
	var activities []*models.Activity

	query := r.db.WithContext(ctx).Table(models.Activity{}.TableName()).Where(cond)

	if err := query.Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if err := query.
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).
		Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (r *ActivityRepositoryImpl) GetActivityByID(
	ctx context.Context,
	cond uuid.UUID,
) (*models.Activity, error) {
	var activity models.Activity
	if err := r.db.WithContext(ctx).
		Table(models.Activity{}.TableName()).
		Where("activity_id = ?", cond).
		First(&activity).Error; err != nil{
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No record found
		}
		return nil, err // Other errors
	}
	return &activity, nil
}
