package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ActivityHandler struct {
	service services.ActivityService
}

func NewActivityHandler(service services.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		service: service,
	}
}

// CreateActivityHandler godoc
// @Summary Create a new activity
// @Description Create a new activity
// @Tags Activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param activity body models.Activity true "Activity data"
// @Success 201 {object} common.ResponseNormal{data=models.Activity} "Activity created successfully"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 401 {object} common.ResponseError "Unauthorized"
// @Failure 403 {object} common.ResponseError "Forbidden"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /expert/activities [post]
func(h *ActivityHandler) CreateActivityHandler(ctx *gin.Context){
	var activityRequest models.Activity

	if err := ctx.ShouldBindJSON(&activityRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.NewResponseError("User ID not found in context"))
		return
	}

	// Kiểm tra kiểu dữ liệu trước khi ép kiểu
	uuidValue, ok := uuid.Parse(userID.(string))
	if ok != nil {
		// Nếu không thể ép kiểu, trả về lỗi
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError("Invalid user ID format"))
		return
	}
	
	actitity, err := h.service.CreateActivity(ctx, uuidValue,&activityRequest);
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, common.NewResponseNormal("Created a new activity successfully",actitity))
}

// GetAtivitiesExpertHandler godoc
// @Summary Get activities for an expert
// @Description Get activities for an expert
// @Tags Activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} common.ResponseNormal{data=[]models.Activity} "Activities retrieved successfully"
// @Failure 400 {object} common.ResponseError "Invalid request parameters"
// @Failure 401 {object} common.ResponseError "Unauthorized"
// @Failure 403 {object} common.ResponseError "Forbidden"
// @Failure 404 {object} common.ResponseError "No activities found"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /expert/activities [get]
func(h *ActivityHandler) GetAtivitiesExpertHandler(ctx *gin.Context){
	var paging common.Paging
	if err := ctx.ShouldBindQuery(&paging); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.NewResponseError("User ID not found in context"))
		return
	}
	// Kiểm tra kiểu dữ liệu trước khi ép kiểu
	uuidValue, ok := uuid.Parse(userID.(string))
	if ok != nil {
		// Nếu không thể ép kiểu, trả về lỗi
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError("Invalid user ID format"))
		return
	}

	activities, err := h.service.GetActivities(ctx, uuidValue, &paging)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	if len(activities) == 0 {
		ctx.JSON(http.StatusNotFound, common.NewResponseError("No activities found"))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponsePaging("Get activities successfully", activities, paging))
}