package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"context"
	"gorm.io/gorm"
)

type ExpertRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	Create(ctx context.Context, expert *models.ExpertRequest) error
	GetExperts(ctx context.Context,  cond map[string]interface{}, paging *common.Paging) ([]*models.Expert, error)
	Update(ctx context.Context, cond map[string]interface{}, updateValue models.Expert) error
	GetExpertByID(ctx context.Context, expertID int) (*models.Expert, error)
	UpdateIsDeleted(ctx context.Context, tx *gorm.DB, expertID int) error
}

type ExpertRepositoryImpl struct {
	DB *gorm.DB
}

func NewExpertRepositoryImpl(db *gorm.DB) *ExpertRepositoryImpl {
	return &ExpertRepositoryImpl{DB: db}
}

func (repo *ExpertRepositoryImpl) Create(ctx context.Context, expert *models.ExpertRequest) error {
	if err := repo.DB.WithContext(ctx).
		Table(models.ExpertRequest{}.TableName()).
		Create(expert).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ExpertRepositoryImpl) GetExperts(
	ctx context.Context,
	cond map[string]interface{},
	paging *common.Paging,
) (
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

func (repo *ExpertRepositoryImpl) Update(
	ctx context.Context,
	cond map[string]interface{},
	updateValue models.Expert,
) error {
	query := repo.DB.WithContext(ctx).Table(models.Expert{}.TableName()).Where(cond)
	if err := query.Updates(updateValue).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ExpertRepositoryImpl) GetExpertByID(ctx context.Context, expertID int) (*models.Expert, error) {
	var expert models.Expert
	if err := repo.DB.WithContext(ctx).
		Table(models.Expert{}.TableName()).
		Where("expert_id = ?", expertID).
		First(&expert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &expert, nil
}

func (repo *ExpertRepositoryImpl) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := repo.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tx, nil
}

func (repo *ExpertRepositoryImpl) UpdateIsDeleted(ctx context.Context, tx *gorm.DB, expertID int) error {
	if err := tx.Table(models.Expert{}.TableName()).
		Where("expert_id = ?", expertID).
		Update("is_deleted", true).Error; err != nil {
		return err
	}
	return nil
}
