package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	InsertActivity(ctx context.Context, tx *gorm.DB, activity *models.Activity) error
	FindActivities(ctx context.Context, cond map[string]interface{}, paging *common.Paging) ([]*models.Activity, error)
	FindActivityByID(ctx context.Context, cond uuid.UUID) (*models.Activity, error)
	UpdateActivityByID(ctx context.Context, tx *gorm.DB, req *models.Activity) error
	DeleteActivityByID(ctx context.Context, tx *gorm.DB, cond uuid.UUID) error
	DeactiveActivity(ctx context.Context, cond uuid.UUID) error
}

type activityRepositoryImpl struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepositoryImpl{
		db: db,
	}
}

func (r *activityRepositoryImpl) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tx, nil
}

func (r *activityRepositoryImpl) InsertActivity(ctx context.Context, tx *gorm.DB, activity *models.Activity) error {
	db := tx
	if db == nil {
		db = r.db
	}

	if err := db.WithContext(ctx).
		Table(models.Activity{}.TableName()).
		Create(activity).Error; err != nil {
		return err
	}
	return nil
}

func (r *activityRepositoryImpl) FindActivities(
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

func (r *activityRepositoryImpl) FindActivityByID(
	ctx context.Context,
	cond uuid.UUID,
) (*models.Activity, error) {
	var activity models.Activity
	if err := r.db.WithContext(ctx).
		Table(models.Activity{}.TableName()).
		Where("activity_id = ?", cond).
		First(&activity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No record found
		}
		return nil, err // Other errors
	}
	return &activity, nil
}

func (r *activityRepositoryImpl) UpdateActivityByID(ctx context.Context, tx *gorm.DB, req *models.Activity) error {
	DB := tx

	if DB == nil {
		DB = r.db
	}

	return DB.WithContext(ctx).
		Table(models.Activity{}.TableName()).
		Where("activity_id = ?", req.ActivityID).
		Updates(req).Error
}

func (r *activityRepositoryImpl) DeleteActivityByID(ctx context.Context, tx *gorm.DB, cond uuid.UUID) error {
	DB := tx

	if DB == nil {
		DB = r.db
	}

	return DB.WithContext(ctx).
		Table(models.Activity{}.TableName()).
		Where("activity_id = ?", cond).
		Delete(models.Activity{}).Error
}

func (r *activityRepositoryImpl) DeactiveActivity(ctx context.Context, cond uuid.UUID) error{
	return r.db.WithContext(ctx).
		Table(models.Activity{}.TableName()).
		Where("activity_id = ?", cond).
		Update("is_active = ?", false).Error
}
