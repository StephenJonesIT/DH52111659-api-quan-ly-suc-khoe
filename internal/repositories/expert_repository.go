package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"context"

	"gorm.io/gorm"
)

type ExpertRepository interface {
	Create(ctx context.Context, expert *models.ExpertCreate) (error)
}

type ExpertRepositoryImpl struct {
	DB *gorm.DB
}

func NewExpertRepositoryImpl(db *gorm.DB) *ExpertRepositoryImpl{
	return &ExpertRepositoryImpl{DB: db}
}

func(repo *ExpertRepositoryImpl) Create(ctx context.Context, expert *models.ExpertCreate) (error) {
	if err := repo.DB.WithContext(ctx).
		Table(models.ExpertCreate{}.TableName()).
		Create(expert).Error; err != nil {
			return err
	}
	return nil
}