package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser godoc
//	@Summary		Create a new user account
//	@Description	Create a new user account with email and password
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string										true	"Bearer Token"
//	@Param			account			body		models.AccountCreate						true	"Account information"
//	@Success		201				{object}	common.ResponseNormal{data=models.Account}	"Account created successfully"
//	@Failure		400				{object}	common.ResponseError						"Invalid request body"
//	@Failure		500				{object}	common.ResponseError						"Internal server error"
//	@Router			/admin/user [post]
func (h *UserHandler) CreateUserHandler(ctx *gin.Context) {
	var accountRequest models.AccountCreate

	if err := ctx.ShouldBindJSON(&accountRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := common.ValidateRequest(accountRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	account, err := h.userService.CreateAccount(ctx, &accountRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, common.NewResponseNormal("Create account successfully", account))
}

// ResetPasswordUser godoc
//	@Summary		Reset user password
//	@Description	Reset user password with email and new password
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization			header		string					true	"Bearer Token"
//	@Param			resetPasswordRequest	body		common.RequestAuth		true	"Reset password request"
//	@Success		200						{object}	common.ResponseNormal	"Password reset successfully"
//	@Failure		400						{object}	common.ResponseError	"Invalid request body"
//	@Failure		500						{object}	common.ResponseError	"Internal server error"
//	@Router			/admin/user/reset-password [post]
func (h *UserHandler) ResetPasswordUserHandler(ctx *gin.Context) {
	var resetPasswordRequest common.RequestAuth

	if err := ctx.ShouldBindJSON(&resetPasswordRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := common.ValidateRequest(resetPasswordRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	if err := h.userService.ResetPassword(ctx, &resetPasswordRequest); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseNormal("Reset password successfully", nil))
}

// GetListUser godoc
//	@Summary		Get list of users
//	@Description	Get a list of users with pagination
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string											true	"Bearer Token"
//	@Param			page			query		int												false	"Page number (default is 1)"
//	@Param			limit			query		int												false	"Number of users per page (default is 10)"
//	@Success		200				{object}	common.ResponseNormal{data=[]models.Account}	"List of users"
//	@Failure		400				{object}	common.ResponseError							"Invalid query parameters"
//	@Failure		500				{object}	common.ResponseError							"Internal server error"
//	@Router			/admin/users [get]
func (h *UserHandler) GetListUserHandler(ctx *gin.Context) {
	var paging common.Paging
	if err := ctx.ShouldBindQuery(&paging); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	accounts, err := h.userService.GetListAccounts(
		ctx,
		&paging,
		map[string]interface{}{
			"account_status": true,   // Only get active accounts
			"role":           "user", // Only get user accounts
			"is_verified":    true,   // Only get verified accounts
		},
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponsePaging("Get list accounts successfully", accounts, paging))
}

// GetUserById godoc
//	@Summary		Get user by ID
//	@Description	Get user details by user ID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string										true	"Bearer Token"
//	@Param			id				path		string										true	"User ID"
//	@Success		200				{object}	common.ResponseNormal{data=models.Account}	"User details"
//	@Failure		400				{object}	common.ResponseError						"Invalid user ID"
//	@Failure		404				{object}	common.ResponseError						"User not found"
//	@Failure		500				{object}	common.ResponseError						"Internal server error"
//	@Router			/admin/user/{id} [get]
func (h *UserHandler) GetUserByIdHandler(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("User ID is required"))
		return
	}

	account, err := h.userService.GetAccountById(ctx, userId)
	if account == nil {
		ctx.JSON(http.StatusNotFound, common.NewResponseError("User not found"))
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseNormal("Get user by ID successfully", account))
}

// LockUserAccount godoc
//	@Summary		Lock user account
//	@Description	Lock user account by user ID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Bearer Token"
//	@Param			id				path		string					true	"User ID"
//	@Success		200				{object}	common.ResponseNormal	"User account locked successfully"
//	@Failure		400				{object}	common.ResponseError	"Invalid user ID"
//	@Failure		500				{object}	common.ResponseError	"Internal server error"
//	@Router			/admin/user/{id}/lock [patch]
func (h *UserHandler) LockUserAccountHandler(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("User ID is required"))
		return
	}

	if err := h.userService.LockAccount(ctx, userId); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseNormal("User account locked successfully", nil))
}

// UnlockUserAccount godoc
//	@Summary		Unlock user account
//	@Description	Unlock user account by user ID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Bearer Token"
//	@Param			id				path		string					true	"User ID"
//	@Success		200				{object}	common.ResponseNormal	"User account unlocked successfully"
//	@Failure		400				{object}	common.ResponseError	"Invalid user ID"
//	@Failure		500				{object}	common.ResponseError	"Internal server error"
//	@Router			/admin/user/{id}/unlock [patch]
func (h *UserHandler) UnlockUserAccountHandler(ctx *gin.Context) {
	userId := ctx.Param("id")
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("User ID is required"))
		return
	}

	if err := h.userService.UnlockAccount(ctx, userId); err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseNormal("User account unlocked successfully", nil))
}
