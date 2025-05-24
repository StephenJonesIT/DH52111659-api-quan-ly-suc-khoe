package services

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
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
	Login(ctx context.Context, loginRequest *common.RequestLogin) (*models.Account, string, string, error)
	ForgotPassowrd(ctx context.Context, email string) (bool, error)
}

type AccountServiceImpl struct {
	accountRepository repositories.AccountRepository
	redisStore        repositories.RedisStore
	emailConfig       EmailConfig
	tokenService	  utils.TokenService
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
		tokenService: *utils.NewTokenService(config.AppConfig.SECRET_KEY),
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

func(s *AccountServiceImpl) Login(ctx context.Context, loginRequest *common.RequestLogin) (*models.Account, string, string, error){
	account, err := s.accountRepository.GetByEmail(ctx, loginRequest.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("lỗi khi lấy tài khoản: %w", err)
	}

	if account == nil {
		return nil, "", "", fmt.Errorf("tài khoản không tồn tại")
	}

	if !account.IsVerified {
		return nil, "", "", fmt.Errorf("tài khoản chưa được xác thực")
	}

	if !utils.ComparePasswordHash(account.Password, loginRequest.Password) {
		return nil, "", "", fmt.Errorf("mật khẩu không đúng")
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokens(*account)
	if err != nil {
		return nil, "", "", fmt.Errorf("lỗi khi tạo token: %w", err)
	}

	return account, accessToken, refreshToken, nil
}

func(s *AccountServiceImpl) ForgotPassowrd(ctx context.Context, email string) (bool, error){
	// Check if the account exists
	account, err := s.accountRepository.GetByEmail(ctx, email)
	
	if err != nil {
		return false, fmt.Errorf("lỗi khi lấy tài khoản: %w", err)
	}

	if account == nil {
		return false, fmt.Errorf("tài khoản không tồn tại")
	}

	// Generate OTP and send it to the email
	sendOTPService := NewSendOTPServiceImpl(s.emailConfig, s.redisStore)
	if err := sendOTPService.SendOTPAndStore(ctx, email); err != nil {
		return false, fmt.Errorf("lỗi khi gửi OTP: %w", err)
	}

	return true, nil
}