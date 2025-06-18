package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProgramRepository interface {
	BeginTx(ctx context.Context) (*gorm.DB, error)
	FindProgramByID(ctx context.Context, cond uuid.UUID)(*models.Program, error)
	FindProgramsByExpertID(ctx context.Context, cond map[string]interface{}, paging *common.Paging) ([]*models.Program, error)
	DeactiveProgram(ctx context.Context, programID uuid.UUID) error
	DeleteProgramByID(ctx context.Context, tx *gorm.DB,programID uuid.UUID) error 
	UpdateProgramByID(ctx context.Context, tx *gorm.DB, cond uuid.UUID, item *models.Program) error
	InsertProgram(ctx context.Context,tx *gorm.DB, item *models.Program) (error)
}

type programRepositoryImpl struct {
	DB *gorm.DB
}

func NewProgramRepository(db *gorm.DB) ProgramRepository{
	return &programRepositoryImpl{DB: db}
}

func (r *programRepositoryImpl) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tx, nil
}

func(r *programRepositoryImpl) InsertProgram(ctx context.Context,tx *gorm.DB, item *models.Program) (error){
	db := tx
	if db == nil {
		db = r.DB
	}

	if err := db.WithContext(ctx).
		Table(models.Program{}.TableName()).
		Create(item).Error; err != nil {
		return err
	}
	return nil
}

func(r *programRepositoryImpl) FindProgramByID(ctx context.Context, cond uuid.UUID)(*models.Program, error){
	var program models.Program
	if err := r.DB.WithContext(ctx).
		Table(models.Program{}.TableName()).
		Where("program_id = ? AND is_active = ?", cond, true).
		First(&program).Error; err != nil{
			return nil, err
	}
	return &program, nil
}

func(r *programRepositoryImpl) FindProgramsByExpertID(
	ctx context.Context,
	cond map[string]interface{},
	paging *common.Paging,
)([]*models.Program, error){
	var programs []*models.Program

	query := r.DB.WithContext(ctx).Table(models.Program{}.TableName()).Where(cond)

	if err := query.Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	if err := query.Order("created_at DESC").
		Offset((paging.Page - 1) *paging.Limit).
		Limit(paging.Limit).
		Find(&programs).Error; err != nil {
		return nil, err
	}

	return programs, nil
}

func(r *programRepositoryImpl)DeactiveProgram(ctx context.Context, programID uuid.UUID) error {
	return r.DB.WithContext(ctx).
		Table(models.Program{}.TableName()).
		Where("program_id = ?",programID).
		Update("is_active = ?", false).Error
}

func(r *programRepositoryImpl)DeleteProgramByID(ctx context.Context, tx *gorm.DB,programID uuid.UUID) error {
	DB := tx
	if DB == nil {
		DB = r.DB
	}
	return DB.WithContext(ctx).
		Table(models.Program{}.TableName()).
		Where("program_id = ?", programID).
		Delete(models.Program{}).Error
}

func(r *programRepositoryImpl)UpdateProgramByID(ctx context.Context, tx *gorm.DB, cond uuid.UUID, item *models.Program) error {
	DB := tx 

	if DB == nil {
		DB = r.DB
	}

	return DB.WithContext(ctx).Table(models.Program{}.TableName()).Where("program_id = ?", cond).Updates(item).Error
}