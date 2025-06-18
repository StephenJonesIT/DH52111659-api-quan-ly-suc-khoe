package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LevelRepository interface {
	InsertLevel(ctx context.Context, tx *gorm.DB, level *models.Level) error
	FindLevelByID(ctx context.Context, cond uuid.UUID) (*models.Level, error)
}

type levelRepository struct {
	db *gorm.DB
}

func NewLevelRepository(db *gorm.DB) LevelRepository {
	return &levelRepository{db: db}
}

func (r *levelRepository) InsertLevel(ctx context.Context, tx *gorm.DB, level *models.Level) error {
	db := tx

	if db == nil {
		db = r.db
	}

	return db.WithContext(ctx).Table(models.Level{}.TableName()).Create(level).Error
}

func (r *levelRepository) FindLevelByID(ctx context.Context, cond uuid.UUID) (*models.Level, error) {
	var levelItem models.Level

	if err := r.db.WithContext(ctx).
		Table(models.Level{}.TableName()).
		Where("level_id = ?", cond).
		First(&levelItem).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return nil, nil
		}
		return nil, err
	}

	return &levelItem, nil
}
