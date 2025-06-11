package services

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/repositories"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ProgramService interface {
	GetProgramByID(ctx context.Context, cond int) (*models.Program, error)
	CreateProgram(ctx context.Context, cond uuid.UUID,item *models.Program) (*models.Program, error)
}

type ProgramServiceImpl struct {
	programRepo repositories.ProgramRepository
	expertRepo repositories.ExpertRepository
}

func NewProgramService(programRepo repositories.ProgramRepository, experRepo repositories.ExpertRepository) *ProgramServiceImpl {
	return &ProgramServiceImpl{
		programRepo: programRepo,
		expertRepo: experRepo,
	}
}

func (s *ProgramServiceImpl) GetProgramByID(ctx context.Context, cond int) (*models.Program, error) {
	program, err := s.programRepo.GetProgramByID(ctx, cond)
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, nil // or return an error if preferred
	}
	return program, nil
}

func (s *ProgramServiceImpl) CreateProgram(ctx context.Context, cond uuid.UUID,item *models.Program) (*models.Program, error) {
	// Check if the expert exists
	expert, err := s.expertRepo.GetExpertByUserID(ctx, cond)
	if err != nil {
		return nil, err
	}

	if expert == nil {
		return nil, fmt.Errorf(common.ErrExpertNotFound) // or a custom error indicating the expert does not exist
	}

	// Set the ExpertID in the Program model
	item.ExpertID = expert.ExpertID
	
	program, err := s.programRepo.Create(ctx, item)
	if err != nil {
		return nil, err
	}
	return program, nil
}



