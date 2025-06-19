package repositories

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserProgramRepository interface {
	CountParticipants(ctx context.Context, programID uuid.UUID) (int64, error)
}

type userProgramRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserProgramRepository(db *gorm.DB) UserProgramRepository{
	return &userProgramRepositoryImpl{
		DB: db,
	}
}

func(r *userProgramRepositoryImpl) CountParticipants(ctx context.Context, programID uuid.UUID) (int64, error){
	var count int64
	if err := r.DB.WithContext(ctx).
		Table(models.UserProgram{}.TableName()).
		Where("program_id = ?",programID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}