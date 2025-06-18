package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProgramDiseaseRepository interface {
	CreateProgramDisease(ctx context.Context, tx *gorm.DB, programDisease *[]models.ProgramDiseases) error
	DeleteProgramDisease(ctx context.Context, tx *gorm.DB, programID uuid.UUID) error
	UpdateProgramDisease(ctx context.Context, tx *gorm.DB, programID uuid.UUID, programDisease *[]models.ProgramDiseases) error
}

type programDiseaseRepository struct {
	db *gorm.DB
}

func NewProgramDiseaseRepository(db *gorm.DB) ProgramDiseaseRepository {
	return &programDiseaseRepository{db: db}
}

func (r *programDiseaseRepository) CreateProgramDisease(ctx context.Context, tx *gorm.DB, programDisease *[]models.ProgramDiseases) error {
	db := tx
	if db == nil {
		db = r.db
	}
	return db.WithContext(ctx).Table(models.ProgramDiseases{}.TableName()).Create(programDisease).Error
}

func (r *programDiseaseRepository) DeleteProgramDisease(ctx context.Context, tx *gorm.DB, programID uuid.UUID) error {
	db := tx
	if db == nil {
		db = r.db
	}

	return db.WithContext(ctx).
		Table(models.ProgramDiseases{}.TableName()).
		Where("program_id = ?", programID).
		Delete(models.ProgramDiseases{}).Error
}

func (r *programDiseaseRepository) UpdateProgramDisease(
	ctx context.Context,
	tx *gorm.DB,
	programID uuid.UUID,
	programDisease *[]models.ProgramDiseases,
) error {
	DB := tx
	if DB == nil {
		DB = r.db
	}

	query := DB.WithContext(ctx).Table(models.ProgramDiseases{}.TableName())

	if err := query.Where("program_id = ?", programID).Delete(models.ProgramDiseases{}).Error; err != nil {
		return err
	}

	if err := query.Create(programDisease).Error; err != nil {
		return err
	}

	return nil
}
