package common

type RequestOTP struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type RequestAuth struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type RequestForgotPassword struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestChangePassword struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=100"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

type RequestRefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}