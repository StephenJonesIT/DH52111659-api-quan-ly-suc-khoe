package services

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	dtos "DH52111659-api-quan-ly-suc-khoe/internal/data/DTOs"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/enum"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/repositories"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ProgramService interface {
	GetProgramByID(ctx context.Context, cond uuid.UUID) (*models.Program, error)
	CreateProgram(ctx context.Context, cond uuid.UUID, item *dtos.CreateProgramRequest) (*models.Program, error)
	RetrieveProgramsByExpertID(ctx context.Context, expertID uuid.UUID, paging *common.Paging) ([]*models.Program, error)
	DeleteProgram(ctx context.Context, expertID, programID uuid.UUID) error
	UpdateProgram(ctx context.Context, programID, expertID uuid.UUID, req *dtos.UpdateProgramRequest) (*models.Program, error)
}

type programServiceImpl struct {
	programRepo        repositories.ProgramRepository
	expertRepo         repositories.ExpertRepository
	levelRepo          repositories.LevelRepository
	activityRepo       repositories.ActivityRepository
	repeatDayRepo      repositories.ActivityRepeatDayRepository
	programDiseaseRepo repositories.ProgramDiseaseRepository
	programGoalRepo    repositories.ProgramGoalRepository
	userProgramRepo    repositories.UserProgramRepository
}

func NewProgramService(
	programRepo repositories.ProgramRepository,
	experRepo repositories.ExpertRepository,
	levelRepo repositories.LevelRepository,
	activityRepo repositories.ActivityRepository,
	repeatDayRepo repositories.ActivityRepeatDayRepository,
	programDiseaseRepo repositories.ProgramDiseaseRepository,
	programGoalsRepo repositories.ProgramGoalRepository,
	userProgramRepo repositories.UserProgramRepository,
) ProgramService {
	return &programServiceImpl{
		programRepo:        programRepo,
		expertRepo:         experRepo,
		levelRepo:          levelRepo,
		activityRepo:       activityRepo,
		repeatDayRepo:      repeatDayRepo,
		programDiseaseRepo: programDiseaseRepo,
		programGoalRepo:    programGoalsRepo,
		userProgramRepo:    userProgramRepo,
	}
}

func (s *programServiceImpl) GetProgramByID(ctx context.Context, cond uuid.UUID) (*models.Program, error) {
	program, err := s.programRepo.FindProgramByID(ctx, cond)
	if err != nil {
		return nil, err
	}
	if program == nil {
		return nil, nil // or return an error if preferred
	}
	return program, nil
}

func (s *programServiceImpl) CreateProgram(ctx context.Context, cond uuid.UUID, req *dtos.CreateProgramRequest) (*models.Program, error) {
	// Check if the expert exists
	expert, err := s.expertRepo.GetExpertByUserID(ctx, cond)
	if err != nil {
		return nil, err
	}

	if expert == nil {
		return nil, fmt.Errorf(common.ErrExpertNotFound) // or a custom error indicating the expert does not exist
	}

	// Init value program
	programID := uuid.New()
	programItem := models.Program{
		ProgramID:   programID,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		IsActive:    true,
		CreatedBy:   expert.ExpertID,
	}

	// Begin transaction
	tx, err := s.programRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", tx.Error)
	}

	if err := s.programRepo.InsertProgram(ctx, tx, &programItem); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create program: %v", err)
	}

	//Insert item on table program disease
	var programDiseases []models.ProgramDiseases
	for _, diseaseID := range req.DiseaseIDs {
		programDiseases = append(programDiseases, models.ProgramDiseases{
			ProgramID: programItem.ProgramID,
			DiseaseID: diseaseID,
		})
	}

	if err := s.programDiseaseRepo.CreateProgramDisease(ctx, tx, &programDiseases); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create program disease: %v", err)
	}

	//Insert goal item on table program goals
	var programGoals []models.ProgramGoal
	for _, goalID := range req.GoalIDs {
		programGoals = append(programGoals, models.ProgramGoal{
			GoalID:    goalID,
			ProgramID: programItem.ProgramID,
		})
	}

	if err := s.programGoalRepo.CreateProgramGoal(ctx, tx, &programGoals); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create program goal: %v", err)
	}

	//Insert level and activity for program
	for _, levelItem := range req.Levels {
		levelID := uuid.New()

		if err := s.levelRepo.InsertLevel(ctx, tx, &models.Level{
			LevelID:      levelID,
			ProgramID:    programItem.ProgramID,
			Name:         levelItem.Name,
			Description:  &levelItem.Description,
			PointRequire: levelItem.PointRequire,
			IsActive:     true,
		}); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create level for program: %v", err)
		}

		//Insert activity for program
		for _, activityItem := range levelItem.Activities {
			activityID := uuid.New()
			// Convert string to *enum.ActivityType
			activityType, err := enum.ParseStr2ActivityType(activityItem.Type)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("invalid activity type: %v", err)
			}

			if err := s.activityRepo.InsertActivity(ctx, tx, &models.Activity{
				ActivityID:  activityID,
				LevelID:     levelID,
				Title:       activityItem.Title,
				Description: activityItem.Description,
				Duration:    activityItem.Duration,
				PointReward: activityItem.PointReward,
				Type:        activityType,
			}); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create activity for level in program: %v", err)
			}

			//Insert repeat day for activity
			var repeatDays []models.ActivityRepeatDay
			for _, repeatDayItem := range activityItem.RepeatDay {
				day, err := enum.ParseStr2WeekDay(repeatDayItem)
				if err != nil {
					tx.Rollback()
					return nil, fmt.Errorf("invalid week day: %v", err)
				}
				repeatDays = append(repeatDays, models.ActivityRepeatDay{
					ActivityID: activityID,
					Repeat:     day,
				})
			}
			if err := s.repeatDayRepo.InsertRepeatDay(ctx, tx, &repeatDays); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create repeat day for activity: %v", err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &programItem, nil
}

func (s *programServiceImpl) RetrieveProgramsByExpertID(
	ctx context.Context,
	expertID uuid.UUID,
	paging *common.Paging,
) ([]*models.Program, error) {
	if expertID == uuid.Nil {
		return nil, fmt.Errorf("expert ID cannot be nil")
	}

	paging.ProcessPaging()

	expertExists, err := s.expertRepo.GetExpertByUserID(ctx, expertID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving expert: %w", err)
	}

	if expertExists == nil {
		return nil, fmt.Errorf("expert not found for user ID: %s", expertID)
	}

	cond := map[string]interface{}{
		"expert_id": expertExists.ExpertID,
	}

	programs, err := s.programRepo.FindProgramsByExpertID(ctx, cond, paging)
	if err != nil {
		return nil, err
	}

	if len(programs) == 0 {
		return nil, fmt.Errorf("no programs found for expert ID: %s", expertID)
	}

	return programs, nil
}

func (s *programServiceImpl) DeleteProgram(ctx context.Context, expertID, programID uuid.UUID) error {
	expertExists, err := s.expertRepo.GetExpertByUserID(ctx, expertID)
	if err != nil {
		return fmt.Errorf("failed to get expert: %v", err)
	}

	if expertExists == nil {
		return fmt.Errorf("expert not found")
	}

	program, err := s.programRepo.FindProgramByID(ctx, programID)
	if err != nil {
		return fmt.Errorf("program not found")
	}

	if program.CreatedBy != expertExists.ExpertID {
		return fmt.Errorf("unauthorized: not program owner")
	}

	//Check the number of participants
	count, err := s.userProgramRepo.CountParticipants(ctx, programID)
	if err != nil {
		return fmt.Errorf("failed to check participants: %v", err)
	}

	if count > 0 {
		if err := s.programRepo.DeactiveProgram(ctx, programID); err != nil {
			return err
		}
	} else {
		tx, err := s.programRepo.BeginTx(ctx)
		if err != nil {
			return fmt.Errorf("failed to start transaction: %v", tx.Error)
		}

		if err := s.programDiseaseRepo.DeleteProgramDisease(ctx, tx, programID); err != nil {
			tx.Rollback()
			return err
		}

		if err := s.programGoalRepo.DeleteProgramGoal(ctx, tx, programID); err != nil {
			tx.Rollback()
			return err
		}
		if err := s.programRepo.DeleteProgramByID(ctx, tx, programID); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit transaction: %v", err)
		}
	}

	return nil
}

func (s *programServiceImpl) UpdateProgram(
	ctx context.Context,
	programID uuid.UUID,
	expertID uuid.UUID,
	req *dtos.UpdateProgramRequest,
) (*models.Program, error) {
	expertExists, err := s.expertRepo.GetExpertByUserID(ctx, expertID)
	if err != nil {
		return nil, fmt.Errorf("failed to get expert: %v", err)
	}

	if expertExists == nil {
		return nil, fmt.Errorf("expert not found")
	}

	programExists, err := s.programRepo.FindProgramByID(ctx, programID)
	if err != nil {
		return nil, fmt.Errorf("program not found")
	}

	if expertExists.ExpertID != programExists.CreatedBy {
		return nil, fmt.Errorf("unauthorized: not program owner")
	}

	programExists.Title = req.Title
	programExists.Description = req.Description
	programExists.Duration = req.Duration

	tx, err := s.programRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction")
	}

	//Update program disease
	var programDiseases []models.ProgramDiseases
	for _, value := range req.DiseaseIDs {
		programDiseases = append(programDiseases, models.ProgramDiseases{
			ProgramID: programExists.ProgramID,
			DiseaseID: value,
		})
	}

	if err := s.programDiseaseRepo.UpdateProgramDisease(ctx, tx, programExists.ProgramID, &programDiseases); err != nil {
		tx.Rollback()
		return nil, err
	}

	//Update program goal
	var programGoals []models.ProgramGoal
	for _, value := range req.GoalIDs {
		programGoals = append(programGoals, models.ProgramGoal{
			ProgramID: programExists.ProgramID,
			GoalID:    value,
		})
	}
	if err := s.programGoalRepo.UpdateProgramGoal(ctx, tx, programExists.ProgramID, &programGoals); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.programRepo.UpdateProgramByID(ctx, tx, programExists.ProgramID, programExists); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return programExists, nil
}
