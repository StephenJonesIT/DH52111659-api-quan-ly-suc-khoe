package services

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/repositories"
	"context"
	"fmt"
)

type ScheduleService interface {
	CreateSchedule(ctx context.Context, schedule *models.ScheduleCreate) (*models.ScheduleCreate, error)
}

type scheduleServiceImpl struct {
	scheduleRepo repositories.ScheduleRepository
	programRepo  repositories.ProgramRepository
	activityRepo repositories.ActivityRepository
}

func NewScheduleServiceImpl(
	scheduleRepo repositories.ScheduleRepository,
	programRepo repositories.ProgramRepository,
	activityRepo repositories.ActivityRepository,
) *scheduleServiceImpl {
	return &scheduleServiceImpl{
		scheduleRepo: scheduleRepo,
		programRepo:  programRepo,
		activityRepo: activityRepo,
	}
}

func (s *scheduleServiceImpl) CreateSchedule(
	ctx context.Context,
	schedule *models.ScheduleCreate,
) (*models.ScheduleCreate, error) {
	if schedule == nil {
		return nil, fmt.Errorf("schedule cannot be nil")
	}

	// Check if the program exists
	program, err := s.programRepo.GetProgramByID(ctx, schedule.ProgramID);
	if err != nil {
		return nil, err 
	}

	if program == nil {
		return nil, fmt.Errorf("program with ID %s not found", schedule.ProgramID)
	}

	// Check if the activity exists
	activity, err := s.activityRepo.GetActivityByID(ctx, schedule.ActivityID)
	if err != nil {
		return nil, err
	}

	if activity == nil {
		return nil, fmt.Errorf("activity with ID %s not found", schedule.ActivityID)
	}

	// Create the schedule
	savedSchedule, err := s.scheduleRepo.Create(ctx, schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}
	return savedSchedule, nil
}
