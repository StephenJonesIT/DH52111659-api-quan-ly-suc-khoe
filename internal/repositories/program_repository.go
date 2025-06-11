package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"context"

	"gorm.io/gorm"
)

type ProgramRepository interface {
	GetProgramByID(ctx context.Context, cond int)(*models.Program, error)
	Create(ctx context.Context, item *models.Program) (*models.Program, error)
}

type ProgramRepositoryImpl struct {
	DB *gorm.DB
}

func NewProgramRepository(db *gorm.DB) *ProgramRepositoryImpl{
	return &ProgramRepositoryImpl{DB: db}
}

func(r *ProgramRepositoryImpl) Create(ctx context.Context, item *models.Program) (*models.Program, error){
	if err := r.DB.WithContext(ctx).
		Table(models.Program{}.TableName()).
		Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func(r *ProgramRepositoryImpl) GetProgramByID(ctx context.Context, cond int)(*models.Program, error){
	var program models.Program
	if err := r.DB.WithContext(ctx).
		Table(models.Program{}.TableName()).
		Where("program_id = ?", cond).
		First(&program).Error; err != nil{
			return nil, err
	}
	return &program, nil
}