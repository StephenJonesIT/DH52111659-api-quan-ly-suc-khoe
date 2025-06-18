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

type ProgramHandler struct {
	programService services.ProgramService
}

func NewProgramHandler(programService services.ProgramService) *ProgramHandler {
	return &ProgramHandler{
		programService: programService,
	}
}

// CreateProgram godoc
// @Summary Create a new program
// @Description Create a new program
// @Tags Program
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param program body dtos.CreateProgramRequest true "Program request data"
// @Success 201 {object} common.ResponseNormal{data=object} "Program created successfully"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 401 {object} common.ResponseError "Unauthorized"
// @Failure 403 {object} common.ResponseError "Forbidden"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /expert/programs [post]
func (h *ProgramHandler) CreateProgramHandler(ctx *gin.Context) {
	var programRequest dtos.CreateProgramRequest

	if err := ctx.ShouldBindJSON(&programRequest); err != nil {
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
	
	// Assuming userID is of type int, adjust as necessary
	program, err := h.programService.CreateProgram(ctx, uuidValue, &programRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, common.NewResponseNormal("Program created successfully", program))
}

// GetProgramsByExpertIDHandler godoc
// @Summary Get programs by expert ID
// @Description Get programs by expert ID
// @Tags Program
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success 200 {object} common.ResponseNormal{data=[]models.Program} "Programs retrieved successfully"
// @Failure 400 {object} common.ResponseError "Invalid expert ID format"
// @Failure 401 {object} common.ResponseError "Unauthorized"
// @Failure 403 {object} common.ResponseError "Forbidden"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /expert/programs [get]
func(h *ProgramHandler) GetProgramsByExpertIDHandler(ctx *gin.Context){
	expertContext, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.NewResponseError("User ID not found in context"))
		return
	}

	expertID, err := uuid.Parse(expertContext.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("Invalid expert ID format"))
		return
	}

	var paging common.Paging
	if err := ctx.ShouldBindQuery(&paging); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
	}

	if expertID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("Expert ID cannot be nil"))
		return
	}

	programs, err := h.programService.RetrieveProgramsByExpertID(ctx, expertID, &paging)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, common.NewResponsePaging("Retrieved programs successfully", programs, paging))
}

// DeleteProgramHandler godoc
// @Summary Delete or deactivate a program by ID
// @Description  Deletes a program if it has no participants, or deactivates it (sets is_active = false) 
// if it has participants. Requires confirmation via query parameter.
// @Tags Program
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param program_id path string true "Program ID (UUID)" Format(uuid)
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} common.ResponseNormal{} "Programs retrieved successfully"
// @Failure 400 {object} common.ResponseError "invalid program ID or confirmation required"
// @Failure 401 {object} common.ResponseError "token expired or invalid"
// @Failure 403 {object} common.ResponseError "unauthorized: not program owner"
// @Failure 500 {object} common.ResponseError "failed to delete program"
// @Router /expert/programs/{program_id} [delete]
func(h *ProgramHandler) DeleteProgramHandler(ctx *gin.Context) {
	expertIDStr, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.NewResponseError("token expired or invalid"))
		return
	}

	expertID, err := uuid.Parse(expertIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("invalid user ID"))
		return
	}

	programIdStr := ctx.Param("program_id")
	programID, err := uuid.Parse(programIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("invalid program ID"))
		return
	}

	if err := h.programService.DeleteProgram(ctx, expertID, programID); err != nil {
		switch err.Error() {
		case "expert not found":
			ctx.JSON(http.StatusNotFound, common.NewResponseError("expert not found"))
		case "program not found":
			ctx.JSON(http.StatusNotFound, common.NewResponseError("program not found"))
		case "unauthorized: not program owner":
			ctx.JSON(http.StatusForbidden, common.NewResponseError("unauthorized: not program owner"))
		default:
			ctx.JSON(http.StatusInternalServerError, common.NewResponseError(fmt.Sprintf("failed to delete program: %v", err)))
		}
		return 
	}

	ctx.JSON(http.StatusOK, common.NewResponseMessage("program deleted successfully"))
}


// UpdateProgramHandler godoc
// @Summary Update a program by ID
// @Description Updates program details (title, description, duration, disease IDs, goal IDs) for a program owned by the expert. Requires valid JWT token.
// @Tags Program
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param program_id path string true "Program ID (UUID)" Format(uuid)
// @Param Authorization header string true "Bearer token"
// @Param program body dtos.UpdateProgramRequest true "Program details to update"
// @Success 200 {object} common.ResponseNormal{data=object} "Program updated successfully"
// @Failure 400 {object} common.ResponseError "invalid program ID or invalid program data"
// @Failure 401 {object} common.ResponseError "token expired or invalid"
// @Failure 403 {object} common.ResponseError "unauthorized: not program owner"
// @Failure 404 {object} common.ResponseError "expert not found or program not found"
// @Failure 500 {object} common.ResponseError "failed to update program"
// @Router /expert/programs/{program_id} [put]
func (h *ProgramHandler) UpdateProgramHandler(ctx *gin.Context) {
	// Lấy expertID từ token
	expertIDStr, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, common.NewResponseError("token expired or invalid"))
		return
	}
	expertID, err := uuid.Parse(expertIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("invalid user ID"))
		return
	}

	// Lấy programID từ URL
	programIDStr := ctx.Param("program_id")
	programID, err := uuid.Parse(programIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError("invalid program ID"))
		return
	}

	// Bind JSON body vào UpdateProgramRequest
	var req dtos.UpdateProgramRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(fmt.Sprintf("invalid program data: %v", err)))
		return
	}

	// Gọi service
	updatedProgram, err := h.programService.UpdateProgram(ctx, programID, expertID, &req)
	if err != nil {
		switch err.Error() {
		case "expert not found", "program not found":
			ctx.JSON(http.StatusNotFound, common.NewResponseError(err.Error()))
		case "unauthorized: not program owner":
			ctx.JSON(http.StatusForbidden, common.NewResponseError("unauthorized: not program owner"))
		case "failed to start transaction":
			ctx.JSON(http.StatusInternalServerError, common.NewResponseError("failed to start transaction"))
		default:
			if err.Error()[0:19] == "invalid program data" {
				ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
			} else {
				ctx.JSON(http.StatusInternalServerError, common.NewResponseError(fmt.Sprintf("failed to update program: %v", err)))
			}
		}
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseNormal("program updated successfully", updatedProgram))
}