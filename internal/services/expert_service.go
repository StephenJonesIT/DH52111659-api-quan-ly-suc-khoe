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
	CreateExpert(ctx context.Context, createExpertRequest *models.ExpertRequest) (*models.ExpertRequest, error)
	GetAllExperts(ctx context.Context, paging *common.Paging) ([]*models.Expert, error)
	UpdateExpert(ctx context.Context, expertID int, expert *models.ExpertRequest) (*models.Expert, error)
	DeleteExpert(ctx context.Context, expertID int) (error)
}

type ExpertServiceImpl struct {
	expertRepo repositories.ExpertRepository
	accountRepo repositories.AccountRepository
}

func NewExpertService(expertRepo repositories.ExpertRepository, accountRepo repositories.AccountRepository) *ExpertServiceImpl {
	return &ExpertServiceImpl{expertRepo: expertRepo, accountRepo: accountRepo}
}

func (s *ExpertServiceImpl) CreateExpert(ctx context.Context, createExpertRequest *models.ExpertRequest) (
	*models.ExpertRequest,
	error,
) {
	if createExpertRequest.TelephoneNumber != "" {
		if !utils.IsValidVietnamesePhoneNumber(createExpertRequest.TelephoneNumber) {
			return nil, fmt.Errorf(common.ErrInvalidPhoneNumber)
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
		map[string]interface{}{"is_deleted": false},
		paging,
	)

	if err != nil {
		return nil, err
	}

	return experts, nil
}

func (s *ExpertServiceImpl) UpdateExpert(
	ctx context.Context,
	expertID int,
	expertRequest *models.ExpertRequest,
) (*models.Expert, error) {
	expertExists, err := s.expertRepo.GetExpertByID(ctx, expertID)
	if err != nil {
		return nil, err
	}

	if expertExists == nil {
		return nil, fmt.Errorf(common.ErrRecordNotFound)
	}

	if expertRequest.TelephoneNumber != "" {
		if !utils.IsValidVietnamesePhoneNumber(expertRequest.TelephoneNumber) {
			return nil, fmt.Errorf(common.ErrInvalidPhoneNumber)
		}
	}

	//Assign value to expert
	expertExists.FullName 			= expertRequest.FullName
	expertExists.DateOfBirth 		= expertRequest.DateOfBirth
	expertExists.Email 				= expertRequest.Email
	expertExists.Gender 			= expertRequest.Gender
	expertExists.TelephoneNumber 	= expertRequest.TelephoneNumber
	expertExists.AvatarURL			= expertRequest.AvatarURL

	if err := s.expertRepo.Update(
		ctx, 
		map[string]interface{}{"expert_id":expertID},
		*expertExists,
	); err != nil {
		return nil, err
	}

	return expertExists, nil
}

func(s *ExpertServiceImpl) DeleteExpert(ctx context.Context, expertID int) (error){
	//Check if the profile expert exists
	expertExists, err := s.expertRepo.GetExpertByID(ctx, expertID)
	if err != nil {
		return err
	}

	if expertExists == nil {
		return fmt.Errorf(common.ErrRecordNotFound)
	}

	// Check if the account expert exists
	accountExists, err := s.accountRepo.GetAccountById(ctx, expertExists.AccountID.String())
	if err != nil {
		return err
	}

	tx, err := s.expertRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {if err != nil {tx.Rollback()}}()

	//Handle set deleted status
	if err := s.expertRepo.UpdateIsDeleted(ctx, tx, expertID); err != nil {
		return err
	}

	if accountExists != nil {
		if err := s.accountRepo.DeactivateAccount(ctx, tx, accountExists.ID.String()); err != nil {
			return err
		}
	}

	return tx.Commit().Error
}