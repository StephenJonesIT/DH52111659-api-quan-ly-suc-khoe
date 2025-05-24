package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(accountService services.AccountService) *AccountHandler{
	return &AccountHandler{
		accountService: accountService,
	}
}

// RegisterAccount godoc
// @Summary Register a new account
// @Description Register a new account with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param account body models.Account true "Account information"
// @Success 201 {object} common.ResponseNormal{email=string} "Account created successfully"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /auth/register [post]
func(h *AccountHandler) RegisterAccountHandler(ctx *gin.Context) {
	var request models.Account

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := h.accountService.RegisterAccount(ctx, &request); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, common.NewResponseRegister("Account created successfully", request.Email))
}

// VerifyOTP godoc
// @Summary Verify OTP for account
// @Description Verify OTP for account registration or password reset
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body common.RequestOTP true "Request OTP information"
// @Success 200 {object} common.ResponseNormal{result=bool} "OTP verified successfully"
// @Failure 400 {object} common.ResponseError "Invalid request parameters"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /auth/verify-email [post]
func(h *AccountHandler) RegisterVerifyOTPHandler(ctx *gin.Context) {
	var request common.RequestOTP

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("Invalid request parameters"))
		return
	}

	isVerified, err := h.accountService.VerifyOTP(ctx, request.Email, request.OTP, false)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseResult("OTP verified successfully", isVerified))
}


// Login godoc
// @Summary Login to the system
// @Description Login to the system with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body common.RequestLogin true "Login request information"
// @Success 200 {object} common.ResponseLogin "Login successful"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /auth/login [post]
func(h *AccountHandler) LoginHandler(ctx *gin.Context) {
	var loginRequest common.RequestLogin

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := common.ValidateRequest(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	account, accessToken, refreshToken, err := h.accountService.Login(ctx, &loginRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseLogin(
		account.ID.String(),
		account.Role,
		accessToken,
		refreshToken,
	))
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Handle forgot password request
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body common.RequestForgotPassword true "Forgot password request information"
// @Success 200 {object} common.ResponseNormal "OTP sent successfully"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /auth/forgot-password [post]
func(h *AccountHandler) ForgotPasswordHandler(ctx *gin.Context){
	var request common.RequestForgotPassword
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := common.ValidateRequest(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	result, err := h.accountService.ForgotPassowrd(ctx, request.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseForgotPassword("OTP sent successfully", request.Email ,result))
}

// VerifyOTPHandler godoc
// @Summary Verify OTP for forgot password
// @Description Verify OTP for forgot password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body common.RequestOTP true "Request OTP information"
// @Success 200 {object} common.ResponseNormal{result=bool} "OTP verified successfully"
// @Failure 400 {object} common.ResponseError "Invalid request parameters"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /auth/verify-otp [post]
func(h *AccountHandler) VerifyOTPHandler(ctx *gin.Context){
	var request common.RequestOTP

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("Invalid request parameters"))
		return
	}

	isVerified, err := h.accountService.VerifyOTP(ctx, request.Email, request.OTP, true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseResult("OTP verified successfully", isVerified))
}

// ResetPasswordHandler godoc
// @Summary Reset password
// @Description Reset password after verifying OTP
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body common.RequestLogin true "Reset password request information"
// @Success 200 {object} common.ResponseNormal "Password reset successfully"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /auth/reset-password [post]
func(h *AccountHandler) ResetPasswordHandler(ctx *gin.Context){
	var request common.RequestLogin
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := common.ValidateRequest(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	if err := h.accountService.ResetPassword(ctx, &request); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseResult("Password reset successfully", true))
}

// ChangPassword godoc
// @Summary Change password
// @Description Change password for the logged-in user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token for authentication"
// @Param request body common.RequestChangePassword true "Change password request information"
// @Success 200 {object} common.ResponseNormal{result=bool} "Password changed successfully"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 401 {object} common.ResponseError "Token must be in Bearer format"
// @Failure 403 {object} common.ResponseError "You do not have permission to access this resource""
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /auth/change-password [post]
func(h *AccountHandler) ChangePasswordHandler(ctx *gin.Context){
	var request common.RequestChangePassword
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := common.ValidateRequest(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.NewResponseError("User ID not found in context"))
		return
	}

	if err := h.accountService.ChangePassword(ctx, userID.(string), &request); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseResult("Password changed successfully", true))
}


// RefreshTokenHandler godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body common.RequestRefreshToken true "Refresh token request information"
// @Success 200 {object} common.ResponseAccessToken
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /auth/refresh-token [post]
func(h *AccountHandler) RefreshTokenHandler(ctx *gin.Context) {
	var request common.RequestRefreshToken
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := common.ValidateRequest(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	accessToken,  err := h.accountService.RefreshToken(ctx, &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseAccessToken(accessToken))
}