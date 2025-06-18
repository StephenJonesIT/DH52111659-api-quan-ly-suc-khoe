package services

import (
	"DH52111659-api-quan-ly-suc-khoe/config"
	"DH52111659-api-quan-ly-suc-khoe/internal/data/repositories"
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
)


// SendOTPService interface defines methods for sending and verifying OTPs.
type SendOTPService interface {
	SendOTPAndStore(ctx context.Context, email string) (error)
	VerifyOTPInRedis(ctx context.Context, email, otp string) (bool, error)
}

// EmailConfig holds the configuration for sending emails.
// It includes SMTP host, port, sender email, and password.
// This struct is used to configure the email sending service.
// It is initialized with values from the application configuration.
// The SMTP host and port are used to connect to the email server.
type EmailConfig struct {
	SMTPHost    string
	SMTPPort    string
	SenderEmail string
	SenderPass  string
}


type SendOTPServiceImpl struct {
	emailConfig   EmailConfig
	redisStore    repositories.RedisStore
}

func NewSendOTPServiceImpl(emailConfig EmailConfig, redisStore repositories.RedisStore) *SendOTPServiceImpl {
	return &SendOTPServiceImpl{
		emailConfig: emailConfig,
		redisStore:  redisStore,
	}
}

func(s *SendOTPServiceImpl) generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func(s *SendOTPServiceImpl) sendOTP(toEmail, otp string) error {
	subject := "Verifi email address with OTP"
	body := fmt.Sprintf("Your authentication code is: %s\nPlease use this code within 10 minutes before it expires.", otp)
	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body,
	)

	auth := smtp.PlainAuth("", config.AppConfig.SenderEmail, config.AppConfig.SenderPass, s.emailConfig.SMTPHost)
	addr := config.AppConfig.SMTPHost + ":" + config.AppConfig.SMTPPort
	err := smtp.SendMail(addr, auth, s.emailConfig.SenderEmail, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("gửi email thất bại: %w", err)
	}
	return nil
}

func(s *SendOTPServiceImpl) SendOTPAndStore(ctx context.Context, email string) (error) {
	otp := s.generateOTP()
	if err := s.redisStore.StoreOTP(ctx, email, otp); err != nil {
		return fmt.Errorf("lưu OTP vào Redis thất bại: %w", err)
	}

	if err := s.sendOTP(email, otp); err != nil {
		return fmt.Errorf("gửi OTP qua email thất bại: %w", err)
	}

	return nil
}

func(s *SendOTPServiceImpl) VerifyOTPInRedis(ctx context.Context, email, otp string) (bool, error) {
	isValid, err := s.redisStore.VerifyOTP(ctx, email, otp)
	if err != nil {
		return false, fmt.Errorf("xác thực OTP thất bại: %w", err)
	}
	if !isValid {
		return false, fmt.Errorf("OTP không hợp lệ hoặc đã hết hạn")
	}
	return true, nil
}

