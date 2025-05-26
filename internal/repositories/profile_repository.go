package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"context"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	GetProfileByID(ctx context.Context, profileID string) (*models.Profile, error)
	Create(ctx context.Context, profile *models.Profile) (*models.Profile, error)
}

type ProfileRepositoryImpl struct {
	DB *gorm.DB
}

func NewProfileRepoImpl(db *gorm.DB) *ProfileRepositoryImpl {
	return &ProfileRepositoryImpl{
		DB: db,
	}
}

func(r *ProfileRepositoryImpl) GetProfileByID(ctx context.Context, profileID string) (*models.Profile, error){
	var profile models.Profile

	if err := r.DB.WithContext(ctx).
		Table(models.Profile{}.TableName()).
		Where("user_id = ?", profileID).
		First(&profile).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, nil
			}
			return nil, err
	}

	return &profile, nil
}

func(r *ProfileRepositoryImpl) Create(ctx context.Context, profile *models.Profile) (*models.Profile, error){

	if err := r.DB.WithContext(ctx).
		Table(models.Profile{}.TableName()).
		Create(profile).Error; err != nil {
		return nil, err
	}

	return profile, nil
}

