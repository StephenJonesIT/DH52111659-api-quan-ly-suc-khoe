package services

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/repositories"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ActivityService interface {
	CreateActivity(ctx context.Context, cond uuid.UUID, activity *models.Activity) (*models.Activity, error)
	GetActivities(ctx context.Context, cond map[string]interface{}, paging *common.Paging) ([]*models.Activity, error)
}

type ActivityServiceImpl struct {
	activityRepository repositories.ActivityRepository
	expertRepository   repositories.ExpertRepository
}

func NewActivityService(activityRepository repositories.ActivityRepository, expertRepository repositories.ExpertRepository) *ActivityServiceImpl {
	return &ActivityServiceImpl{
		activityRepository: activityRepository,
		expertRepository:   expertRepository,
	}
}

func (s *ActivityServiceImpl) CreateActivity(ctx context.Context, cond uuid.UUID, activity *models.Activity) (*models.Activity, error) {
	if activity == nil {
		return nil, fmt.Errorf("activity cannot be nil")
	}

	// Check if the expert exists
	expert, err := s.expertRepository.GetExpertByUserID(ctx, cond)

	if err != nil {
		return nil, err
	}

	if expert == nil {
		return nil, fmt.Errorf("expert not found for user ID: %s", cond)
	}

	// Set the ExpertID in the Activity model
	activity.ExpertID = expert.ExpertID

	// Call the repository to create the activity
	createdActivity, err := s.activityRepository.CreateActivity(ctx, activity)
	if err != nil {
		return nil, err
	}

	return createdActivity, nil
}

func (s *ActivityServiceImpl) GetActivities(
	ctx context.Context,
	cond map[string]interface{},
	paging *common.Paging,
) ([]*models.Activity, error){
	paging.ProcessPaging()

	activities, err := s.activityRepository.GetActivities(ctx, cond, paging)
	if err != nil {
		return nil, err
	}
	if len(activities) == 0 {
		return nil, fmt.Errorf("no activities found")
	}
	return activities, nil
}
