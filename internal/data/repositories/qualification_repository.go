package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"gorm.io/gorm"
)

type QualificationRepository interface {
	Create(ctx context.Context, qualificationRequest *models.QualificationRequest) (*models.Qualification, error)
}

type QualificationRepoImpl struct {
	DB *gorm.DB
}

func NewQualificationRepository(db *gorm.DB) *QualificationRepoImpl {
	return &QualificationRepoImpl{
		DB: db,
	}
}

func (r *QualificationRepoImpl) Create(ctx context.Context, qualificationRequest *models.QualificationRequest) (error){
	return nil
}