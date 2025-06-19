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

type ActivityService interface {
	CreateActivity(ctx context.Context, req *dtos.CreateActivityRequest) (*models.Activity, error)
	GetActivities(ctx context.Context, expertID uuid.UUID, paging *common.Paging) ([]*models.Activity, error)
	UpdateActivity(ctx context.Context, cond string, req *dtos.UpdateActivityRequest) (*models.Activity, error)
	DeleteActivity(ctx context.Context, cond string) error
}

type ActivityServiceImpl struct {
	activityRepo     repositories.ActivityRepository
	levelRepo        repositories.LevelRepository
	repeatDayRepo    repositories.ActivityRepeatDayRepository
	userActivityRepo repositories.UserActivityRepository
}

func NewActivityService(
	activityRepository repositories.ActivityRepository,
	levelRepo repositories.LevelRepository,
	repeatDayRepo repositories.ActivityRepeatDayRepository,
	userActivityRepo repositories.UserActivityRepository) ActivityService {
	return &ActivityServiceImpl{
		activityRepo:     activityRepository,
		levelRepo:        levelRepo,
		repeatDayRepo:    repeatDayRepo,
		userActivityRepo: userActivityRepo,
	}
}

func (s *ActivityServiceImpl) CreateActivity(ctx context.Context, req *dtos.CreateActivityRequest) (*models.Activity, error) {
	levelID, err := uuid.Parse(req.LevelID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse level id")
	}

	//Check level exists
	levelExists, err := s.levelRepo.FindLevelByID(ctx, levelID)
	if err != nil {
		return nil, fmt.Errorf("failed to find level")
	}

	if levelExists == nil {
		return nil, fmt.Errorf("level not found")
	}

	tx, err := s.activityRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction")
	}

	//Insert activity
	activityID := uuid.New()

	activityType, err := enum.ParseStr2ActivityType(req.Type)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("invalid activity type: %v", err)
	}

	activity := &models.Activity{
		LevelID:     levelExists.LevelID,
		ActivityID:  activityID,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		PointReward: req.PointReward,
		Type:        activityType,
	}

	if err := s.activityRepo.InsertActivity(ctx, tx, activity); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create activity")
	}

	//Insert repeat day for activity
	var repeatDays []models.ActivityRepeatDay
	for _, value := range req.RepeatDays {
		day, err := enum.ParseStr2WeekDay(value)
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
		return nil, fmt.Errorf("failed to create repeat day")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction")
	}

	return activity, nil
}

func (s *ActivityServiceImpl) GetActivities(
	ctx context.Context,
	expertID uuid.UUID,
	paging *common.Paging,
) ([]*models.Activity, error) {
	return nil, nil
}

func (s *ActivityServiceImpl) UpdateActivity(
	ctx context.Context,
	cond string,
	req *dtos.UpdateActivityRequest,
) (*models.Activity, error) {
	activityID, err := uuid.Parse(cond)
	if err != nil {
		return nil, fmt.Errorf("failed to parse uuid activity id")
	}

	// check activity exists
	activityExists, err := s.activityRepo.FindActivityByID(ctx, activityID)
	if err != nil {
		return nil, fmt.Errorf("failed to find activity with id")
	}

	if activityExists == nil {
		return nil, fmt.Errorf("activity not found")
	}

	tx, err := s.activityRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction")
	}

	activityType, err := enum.ParseStr2ActivityType(req.Type)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("invalid activity type: %v", err)
	}
	activityExists.Title = req.Title
	activityExists.Description = req.Description
	activityExists.Duration = req.Duration
	activityExists.PointReward = req.PointReward
	activityExists.Type = activityType

	if err := s.activityRepo.UpdateActivityByID(ctx, tx, activityExists); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update activity")
	}

	// Update repeat day for activity
	var repeatDays []models.ActivityRepeatDay
	for _, value := range req.RepeatDays {
		day, err := enum.ParseStr2WeekDay(value)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("faild to parse repeat day")
		}
		repeatDays = append(repeatDays, models.ActivityRepeatDay{
			ActivityID: activityExists.ActivityID,
			Repeat:     day,
		})
	}

	if err := s.repeatDayRepo.UpdateRepeatDayByID(ctx, tx, activityExists.ActivityID, &repeatDays); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update repeat day")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction")
	}

	return activityExists, nil
}

func (s *ActivityServiceImpl) DeleteActivity(ctx context.Context, cond string) error {
	activityID, err := uuid.Parse(cond)
	if err != nil {
		return fmt.Errorf("failed to parse uuid activity id")
	}

	// check activity exists
	activityExists, err := s.activityRepo.FindActivityByID(ctx, activityID)
	if err != nil {
		return fmt.Errorf("failed to find activity with id")
	}

	if activityExists == nil {
		return fmt.Errorf("activity not found")
	}

	count, err := s.userActivityRepo.CountParticipants(ctx, activityExists.ActivityID)
	if err != nil {
		return fmt.Errorf("failed to check participants: %v", err)
	}

	if count > 0 {
		if err := s.activityRepo.DeactiveActivity(ctx, activityExists.ActivityID); err != nil {
			return err
		}
	} else {
		tx, err := s.activityRepo.BeginTx(ctx)
		if err != nil {
			return fmt.Errorf("failed to start transaction: %v", tx.Error)
		}

		if err := s.repeatDayRepo.DeleteRepeatDayByID(ctx, tx, activityExists.ActivityID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete repeat day")
		}

		if err := s.activityRepo.DeleteActivityByID(ctx, tx, activityExists.ActivityID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete activity")
		}

		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit transaction")
		}
	}
	return nil
}
