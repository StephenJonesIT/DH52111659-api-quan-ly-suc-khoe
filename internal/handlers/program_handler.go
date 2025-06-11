package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
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
// @Param program body models.Program true "Program data"
// @Success 201 {object} common.ResponseNormal{data=object} "Program created successfully"
// @Failure 400 {object} common.ResponseError "Invalid request body"
// @Failure 401 {object} common.ResponseError "Unauthorized"
// @Failure 403 {object} common.ResponseError "Forbidden"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /expert/programs [post]
func (h *ProgramHandler) CreateProgramHandler(ctx *gin.Context) {
	var programRequest models.Program

	if err := ctx.ShouldBindJSON(&programRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
		return
	}

	if err := common.ValidateRequest(programRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
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
