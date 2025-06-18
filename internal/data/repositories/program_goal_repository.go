package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProgramGoalRepository interface {
	CreateProgramGoal(ctx context.Context, tx *gorm.DB, programGoal *[]models.ProgramGoal) error
	DeleteProgramGoal(ctx context.Context, tx *gorm.DB, programID uuid.UUID) error
	UpdateProgramGoal(ctx context.Context, tx *gorm.DB, programID uuid.UUID, programGoal *[]models.ProgramGoal) error
}

type programGoalRepository struct {
	db *gorm.DB
}

func NewProgramGoalRepository(db *gorm.DB) ProgramGoalRepository {
	return &programGoalRepository{db: db}
}

func (r *programGoalRepository) CreateProgramGoal(ctx context.Context, tx *gorm.DB, programGoal *[]models.ProgramGoal) error {
	db := tx
	if db == nil {
		db = r.db
	}

	return db.WithContext(ctx).Table(models.ProgramGoal{}.TableName()).Create(programGoal).Error
}

func (r *programGoalRepository) DeleteProgramGoal(ctx context.Context, tx *gorm.DB, programID uuid.UUID) error {
	DB := tx
	if DB == nil {
		DB = r.db
	}

	return DB.WithContext(ctx).
		Table(models.ProgramGoal{}.TableName()).
		Where("program_id = ?", programID).
		Delete(models.ProgramGoal{}).Error
}


func (r *programGoalRepository) UpdateProgramGoal(
	ctx context.Context,
	tx *gorm.DB,
	programID uuid.UUID,
	programGoal *[]models.ProgramGoal,
) error {
	DB := tx
	if DB == nil {
		DB = r.db
	}

	query := DB.WithContext(ctx).Table(models.ProgramGoal{}.TableName())

	if err := query.Where("program_id = ?", programID).Delete(models.ProgramGoal{}).Error; err != nil {
		return err
	}

	if err := query.Create(programGoal).Error; err != nil {
		return err
	}

	return nil
}
