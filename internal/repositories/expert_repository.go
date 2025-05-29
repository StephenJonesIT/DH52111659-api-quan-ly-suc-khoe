package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"context"

	"gorm.io/gorm"
)

type ExpertRepository interface {
	Create(ctx context.Context, expert *models.ExpertCreate) error
	GetExperts(ctx context.Context, cond map[string]interface{}, paging *common.Paging) ([]*models.Expert, error)
}

type ExpertRepositoryImpl struct {
	DB *gorm.DB
}

func NewExpertRepositoryImpl(db *gorm.DB) *ExpertRepositoryImpl {
	return &ExpertRepositoryImpl{DB: db}
}

func (repo *ExpertRepositoryImpl) Create(ctx context.Context, expert *models.ExpertCreate) error {
	if err := repo.DB.WithContext(ctx).
		Table(models.ExpertCreate{}.TableName()).
		Create(expert).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ExpertRepositoryImpl) GetExperts(ctx context.Context, cond map[string]interface{}, paging *common.Paging) (
	[]*models.Expert,
	error,
) {
	var experts []*models.Expert
	query := repo.DB.WithContext(ctx).Table(models.Expert{}.TableName()).Where(cond)
	if err := query.Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if err := query.Order("expert_id DESC").
		Offset((paging.Page - 1) * paging.Limit).
		Limit(paging.Limit).Find(&experts).Error; err != nil {
		return nil, err
	}

	return experts, nil
}
