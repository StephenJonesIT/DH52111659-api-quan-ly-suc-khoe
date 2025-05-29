package services

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/repositories"
	"DH52111659-api-quan-ly-suc-khoe/utils"
	"context"
	"fmt"
)

type ExpertService interface {
	CreateExpert(ctx context.Context, createExpertRequest *models.ExpertCreate) (*models.ExpertCreate, error)
	GetAllExperts(ctx context.Context, paging *common.Paging) ([]*models.Expert, error)
}

type ExpertServiceImpl struct {
	expertRepo repositories.ExpertRepository
}

func NewExpertService(repo repositories.ExpertRepository) *ExpertServiceImpl {
	return &ExpertServiceImpl{expertRepo: repo}
}

func (s *ExpertServiceImpl) CreateExpert(ctx context.Context, createExpertRequest *models.ExpertCreate) (
	*models.ExpertCreate,
	error,
) {
	if createExpertRequest.TelephoneNumber != "" {
		if !utils.IsValidVietnamesePhoneNumber(createExpertRequest.TelephoneNumber) {
			return nil, fmt.Errorf("số điện thoại không hợp lệ")
		}
	}

	if err := s.expertRepo.Create(ctx, createExpertRequest); err != nil {
		return nil, err
	}

	return createExpertRequest, nil
}

func (s *ExpertServiceImpl) GetAllExperts(
	ctx context.Context,
	paging *common.Paging,
) ([]*models.Expert, error) {
	paging.ProcessPaging()

	experts, err := s.expertRepo.GetExperts(
		ctx,
		map[string]interface{}{"is_deleted":false},
		paging,
	)

	if err != nil {
		return nil, err
	}

	return experts, nil
}
