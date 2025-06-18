package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	dtos "DH52111659-api-quan-ly-suc-khoe/internal/data/DTOs"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"fmt"
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
// @Summary Create a new activity for a level
// @Description Creates a new activity for a specified level, including repeat days. Requires valid JWT token.
// @Tags Activity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param activity body dtos.CreateActivityRequest true "Activity details to create"
// @Success 200 {object} common.ResponseNormal{} "Activity created successfully"
// @Failure 400 {object} common.ResponseError "invalid activity data, invalid level ID, invalid activity type, or invalid week day"
// @Failure 401 {object} common.ResponseError "token expired or invalid"
// @Failure 404 {object} common.ResponseError "level not found"
// @Failure 500 {object} common.ResponseError "failed to create activity"
// @Router /expert/activities [post]
func (h *ActivityHandler) CreateActivityHandler(ctx *gin.Context) {
	// Bind JSON body vào CreateActivityRequest
	var req *dtos.CreateActivityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(fmt.Sprintf("invalid activity data: %v", err)))
		return
	}

	// Gọi service
	createdActivity, err := h.service.CreateActivity(ctx, req)
	if err != nil {
		switch err.Error() {
		case "level not found":
			ctx.JSON(http.StatusNotFound, common.NewResponseError("level not found"))
		case "failed to parse level id":
			ctx.JSON(http.StatusBadRequest, common.NewResponseError("invalid level ID"))
		case "failed to start transaction":
			ctx.JSON(http.StatusInternalServerError, common.NewResponseError("failed to start transaction"))
		default:
			if err.Error()[0:19] == "invalid activity type" || err.Error()[0:15] == "invalid week day" {
				ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
			} else {
				ctx.JSON(http.StatusInternalServerError, common.NewResponseError(fmt.Sprintf("failed to create activity: %v", err)))
			}
		}
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseNormal("activity created successfully", createdActivity))
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