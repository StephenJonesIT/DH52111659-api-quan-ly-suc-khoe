package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleService services.ScheduleService
}

func NewScheduleHandler(schedule services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: schedule,
	}
}

// CreateScheduleHandler godoc
// @Summary Create a new schedule
// @Description Create a new schedule for a program
// @Tags Schedule
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param schedule body models.ScheduleCreate true "Schedule data"
// @Success 201 {object} common.ResponseNormal{data=models.ScheduleCreate} "Schedule created successfully"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 401 {object} common.ResponseError "Unauthorized"
// @Failure 403 {object} common.ResponseError "Forbidden"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router  /expert/programs/{program_id}/schedules [post]
func(h *ScheduleHandler) CreateScheduleHandler(ctx *gin.Context){
	var scheduleRequest models.ScheduleCreate
	if err := ctx.ShouldBindJSON(&scheduleRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := common.ValidateRequest(scheduleRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	scheduleResponse, err := h.scheduleService.CreateSchedule(ctx, &scheduleRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, common.NewResponseNormal("Created a new schedule for program",scheduleResponse))
}