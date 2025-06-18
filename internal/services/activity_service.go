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
}

type ActivityServiceImpl struct {
	activityRepo repositories.ActivityRepository
	levelRepo repositories.LevelRepository
	repeatDayRepo repositories.ActivityRepeatDayRepository
}

func NewActivityService(
	activityRepository repositories.ActivityRepository, 
	levelRepo repositories.LevelRepository,
	repeatDayRepo repositories.ActivityRepeatDayRepository) ActivityService {
	return &ActivityServiceImpl{
		activityRepo: activityRepository,
		levelRepo: levelRepo,
		repeatDayRepo: repeatDayRepo,
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
		LevelID: levelExists.LevelID,
		ActivityID: activityID,
		Title: req.Title,
		Description: req.Description,
		Duration: req.Duration,
		PointReward: req.PointReward,
		Type: activityType,
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
			Repeat: day,
		})
	}

	if err := s.repeatDayRepo.InsertRepeatDay(ctx, tx, &repeatDays); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create repeat day")
	}

	if err := tx.Commit().Error; err != nil{
		return nil, fmt.Errorf("failed to commit transaction")
	}

	return activity, nil
}

func (s *ActivityServiceImpl) GetActivities(
	ctx context.Context,
	expertID uuid.UUID,
	paging *common.Paging,
) ([]*models.Activity, error){
	return nil, nil
}
