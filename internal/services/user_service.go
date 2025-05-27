package services

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/repositories"
	"DH52111659-api-quan-ly-suc-khoe/utils"
	"context"
	"fmt"
)

type UserService interface {
	CreateAccount(ctx context.Context, account *models.AccountCreate) (*models.Account, error)
}

type UserServiceImpl struct {
	accountRepository repositories.AccountRepository
}

func NewUserServiceImpl(accountRepo repositories.AccountRepository) *UserServiceImpl {
	return &UserServiceImpl{accountRepository: accountRepo}
}

func(s *UserServiceImpl) CreateAccount(ctx context.Context, accountRequest *models.AccountCreate) (*models.Account, error){
	// Check if the account already exists
	existsAccount, err := s.accountRepository.GetByEmail(ctx, accountRequest.Email)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi kiểm tra tài khoản: %w", err)
	}

	if existsAccount != nil {
		return nil, fmt.Errorf("tài khoản đã tồn tại")
	}

	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(accountRequest.Password)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi mã hóa mật khẩu: %w", err)
	}

	
	account := &models.Account{
		Email:    accountRequest.Email,
		Password: hashedPassword,
		Role:     "user",
		IsVerified: true,
		AccountStatus: true,
	}

	if err := s.accountRepository.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("lỗi khi tạo tài khoản: %w", err)
	}

	return account, nil
}