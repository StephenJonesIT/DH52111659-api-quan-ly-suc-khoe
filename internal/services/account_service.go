package services

import (
	"DH52111659-api-quan-ly-suc-khoe/config"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/repositories"
	"DH52111659-api-quan-ly-suc-khoe/utils"
	"context"
	"fmt"
)

type AccountService interface {
	RegisterAccount(ctx context.Context, account *models.Account) error
	VerifyOTP(ctx context.Context, toEmail, otp string, isVerified bool) (bool, error)
}

type AccountServiceImpl struct {
	accountRepository repositories.AccountRepository
	redisStore        repositories.RedisStore
	emailConfig       EmailConfig
}

func NewAccountServiceImpl(accountRepo repositories.AccountRepository, redis repositories.RedisStore) *AccountServiceImpl {
	return &AccountServiceImpl{
		accountRepository: accountRepo,
		redisStore:        redis,
		emailConfig: EmailConfig{
			SMTPHost:    config.AppConfig.SMTPHost,
			SMTPPort:    config.AppConfig.SMTPPort,
			SenderEmail: config.AppConfig.SenderEmail,
			SenderPass:  config.AppConfig.SenderPass,
		},
	}
}

func (s *AccountServiceImpl) RegisterAccount(ctx context.Context, account *models.Account) error {
	// Check if the account already exists
	existsAccount, err := s.accountRepository.GetByEmail(ctx, account.Email)
	if err != nil {
		return err
	}

	if existsAccount != nil {
		return fmt.Errorf("tài khoản đã tồn tại")
	}

	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(account.Password)
	if err != nil {
		return fmt.Errorf("lỗi khi mã hóa mật khẩu: %w", err)
	}

	// Set the hashed password back to the account
	account.Password = hashedPassword

	// Create the account in the database
	if err := s.accountRepository.Create(ctx, account); err != nil {
		return fmt.Errorf("lỗi khi tạo tài khoản: %w", err)
	}

	// Generate OTP
	sendOTPService := NewSendOTPServiceImpl(s.emailConfig, s.redisStore)
	if err := sendOTPService.SendOTPAndStore(ctx, account.Email); err != nil {
		return err
	}

	return nil
}

func (s *AccountServiceImpl) VerifyOTP(ctx context.Context, toEmail, otp string, isVerified bool) (bool, error) {
	// Verify OTP
	sendOTPService := NewSendOTPServiceImpl(s.emailConfig, s.redisStore)
	result, err := sendOTPService.VerifyOTPInRedis(ctx, toEmail, otp)
	if err != nil {
		return result, err
	}

	// Update account status if OTP is valid and not already verified
	if !isVerified {
		isVerified = true
		if err := s.accountRepository.Update(
			ctx, 
			map[string]interface{}{"email":toEmail}, 
			map[string]interface{}{"is_verified": isVerified,}); err != nil {
			return false, fmt.Errorf("lỗi khi cập nhật trạng thái tài khoản: %w", err)
		}
	}

	return result, nil
}
