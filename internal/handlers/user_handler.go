package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"net/http"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler{
	return &UserHandler{userService: userService}
}

// CreateUser godoc
// @Summary Create a new user account
// @Description Create a new user account with email and password
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param account body models.AccountCreate true "Account information"
// @Success 201 {object} common.ResponseNormal{data=models.Account} "Account created successfully"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /admin/user [post]
func(h *UserHandler) CreateUserHandler(ctx *gin.Context) {
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

	ctx.JSON(http.StatusCreated, common.NewResponseNormal("Create account successfully",account))
}