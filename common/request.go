package common

type RequestOTP struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type RequestForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestChangePassword struct {
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

type RequestRefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}