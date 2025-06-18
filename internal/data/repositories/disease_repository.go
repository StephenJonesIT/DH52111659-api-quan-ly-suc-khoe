package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"gorm.io/gorm"
)

type DiseaseRepository interface {
	CreateDisease(ctx context.Context, tx *gorm.DB, itemDiseaseProgram *models.ProgramDiseases)(error)
}

type diseaseRepositoryImpl struct {
	DB *gorm.DB
}

func NewDiseaseRepository(db *gorm.DB) DiseaseRepository {
	return &diseaseRepositoryImpl{
		DB: db,
	}
}

func(repo *diseaseRepositoryImpl) CreateDisease(ctx context.Context, tx *gorm.DB, itemDiseaseProgram *models.ProgramDiseases)(error){
	return nil
}