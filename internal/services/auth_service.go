package services

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/config"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/repositories"
	"DH52111659-api-quan-ly-suc-khoe/utils"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type AuthService interface {
	RegisterAccount(ctx context.Context, account *models.Account) error
	VerifyOTP(ctx context.Context, toEmail, otp string, isVerified bool) (bool, error)
	Login(ctx context.Context, loginRequest *common.RequestAuth) (*models.Account, string, string, error)
	ForgotPassowrd(ctx context.Context, email string) (bool, error)
	ResetPassword(ctx context.Context, resetPasswordRequest *common.RequestAuth) error
	ChangePassword(ctx context.Context, id string, changePasswordRequest *common.RequestChangePassword) error
	RefreshToken(ctx context.Context, requestRefreshToken *common.RequestRefreshToken) (string, error)
}

type AuthServiceImpl struct {
	accountRepository repositories.AccountRepository
	redisStore        repositories.RedisStore
	emailConfig       EmailConfig
	tokenService      utils.TokenService
}

func NewAuthServiceImpl(accountRepo repositories.AccountRepository, redis repositories.RedisStore) *AuthServiceImpl {
	return &AuthServiceImpl{
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

func (s *AuthServiceImpl) RegisterAccount(ctx context.Context, account *models.Account) error {
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

func (s *AuthServiceImpl) VerifyOTP(ctx context.Context, toEmail, otp string, isVerified bool) (bool, error) {
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
			map[string]interface{}{"email": toEmail},
			map[string]interface{}{"is_verified": isVerified}); err != nil {
			return false, fmt.Errorf("lỗi khi cập nhật trạng thái tài khoản: %w", err)
		}
	}

	return result, nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, loginRequest *common.RequestAuth) (*models.Account, string, string, error) {
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

	if !account.AccountStatus {
		return nil, "", "", fmt.Errorf("tài khoản đã bị khóa")
	}

	if !utils.ComparePasswordHash(account.Password, loginRequest.Password) {
		return nil, "", "", fmt.Errorf("mật khẩu không chính xác")
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokens(*account)
	if err != nil {
		return nil, "", "", fmt.Errorf("lỗi khi tạo token: %w", err)
	}

	return account, accessToken, refreshToken, nil
}

func (s *AuthServiceImpl) ForgotPassowrd(ctx context.Context, email string) (bool, error) {
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

func (s *AuthServiceImpl) ResetPassword(ctx context.Context, resetPasswordRequest *common.RequestAuth) error {
	// Check if the account exists
	account, err := s.accountRepository.GetByEmail(ctx, resetPasswordRequest.Email)
	if err != nil {
		return fmt.Errorf("lỗi khi lấy tài khoản: %w", err)
	}

	if account == nil {
		return fmt.Errorf("tài khoản không tồn tại")
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(resetPasswordRequest.Password)
	if err != nil {
		return fmt.Errorf("lỗi khi mã hóa mật khẩu: %w", err)
	}

	// Update the account's password
	if err := s.accountRepository.Update(
		ctx,
		map[string]interface{}{"email": resetPasswordRequest.Email},
		map[string]interface{}{"password_hash": hashedPassword},
	); err != nil {
		return fmt.Errorf("lỗi khi cập nhật mật khẩu: %w", err)
	}

	return nil
}

func (s *AuthServiceImpl) ChangePassword(ctx context.Context, id string, changePasswordRequest *common.RequestChangePassword) error {
	// Get the account by ID
	account, err := s.accountRepository.GetAccountById(ctx, id)
	if err != nil {
		return fmt.Errorf("lỗi khi lấy tài khoản: %w", err)
	}

	if account == nil {
		return fmt.Errorf("tài khoản không tồn tại")
	}

	// Check if the old password matches
	if !utils.ComparePasswordHash(account.Password, changePasswordRequest.OldPassword) {
		return fmt.Errorf("mật khẩu cũ không chính xác")
	}
	
	// Hash the new password
	hashedPassword, err := utils.HashPassword(changePasswordRequest.NewPassword)
	if err != nil {
		return fmt.Errorf("lỗi khi mã hóa mật khẩu: %w", err)
	}

	// Update the account's password
	if err := s.accountRepository.Update(
		ctx,
		map[string]interface{}{"id": id},
		map[string]interface{}{"password_hash": hashedPassword},
	); err != nil {
		return fmt.Errorf("lỗi khi cập nhật mật khẩu: %w", err)
	}

	return nil
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, requestRefreshToken *common.RequestRefreshToken) (string, error) {
	// Verify the refresh token
	account, err := s.tokenService.VerifyToken(requestRefreshToken.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("lỗi khi xác thực refresh token: %w", err)
	}

	// If the account is nil, it means the token is invalid or expired
	if account == nil {
		return "", fmt.Errorf("refresh token không hợp lệ hoặc đã hết hạn")
	}

	var accountTemp models.Account
	// If the account is nil, try to get it by ID
	parsedUUID, err := uuid.Parse(account.ID)
	if err != nil {
		return "", fmt.Errorf("lỗi khi chuyển đổi ID: %w", err)
	}
	accountTemp.ID = parsedUUID
	accountTemp.Role = account.Role

	// Generate new access and refresh tokens
	newAccessToken, _, err := s.tokenService.GenerateTokens(accountTemp)
	if err != nil {
		return "", fmt.Errorf("lỗi khi tạo token mới: %w", err)
	}

	return newAccessToken, nil
}
