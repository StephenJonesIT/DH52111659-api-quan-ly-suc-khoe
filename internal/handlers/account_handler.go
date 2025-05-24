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
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("Invalid request body"))
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
