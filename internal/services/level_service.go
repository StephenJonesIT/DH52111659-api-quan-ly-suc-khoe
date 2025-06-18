package services

import (
	dtos "DH52111659-api-quan-ly-suc-khoe/internal/data/DTOs"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/repositories"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type LevelService interface {
	CreateLevel(ctx context.Context, req *dtos.CreateLevelRequest) (*models.Level, error)
}

type levelServiceImpl struct{
	programRepo repositories.ProgramRepository
	levelRepo repositories.LevelRepository
}

func NewLevelService(programRepo repositories.ProgramRepository, levelRepo repositories.LevelRepository) LevelService {
	return &levelServiceImpl{programRepo: programRepo, levelRepo: levelRepo}
}

func(s *levelServiceImpl) CreateLevel(ctx context.Context, req *dtos.CreateLevelRequest) (*models.Level, error){
	//Parse program id
	cond, err := uuid.Parse(req.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse program id")
	}

	//Check the exists program
	program, err := s.programRepo.FindProgramByID(ctx, cond)
	if err != nil{
		return nil, fmt.Errorf("program not found")
	}
	
	levelID := uuid.New()
	var levelItem = &models.Level{
		LevelID: levelID,
		ProgramID: program.ProgramID,
		Name: req.Name,
		Description: &req.Description,
		PointRequire: req.PointRequire,
		IsActive: true,
	}

	if err := s.levelRepo.InsertLevel(ctx, nil, levelItem); err != nil {
		return nil, err
	}

	return levelItem, nil
}