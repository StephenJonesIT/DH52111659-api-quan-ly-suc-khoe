package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	dtos "DH52111659-api-quan-ly-suc-khoe/internal/data/DTOs"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LevelHandler struct {
	levelService services.LevelService
}

func NewLevelHandler(levelService services.LevelService) *LevelHandler{
	return &LevelHandler{levelService: levelService}
}

// CreateLevelHandler godoc
// @Summary Create a new level for a program
// @Description Creates a new level for a specified program owned by the expert. Requires valid JWT token.
// @Tags Level
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param level body dtos.CreateLevelRequest true "Level details to create"
// @Success 200 {object} common.ResponseNormal{} "Level created successfully"
// @Failure 400 {object} common.ResponseError "invalid level data or invalid program ID"
// @Failure 401 {object} common.ResponseError "token expired or invalid"
// @Failure 404 {object} common.ResponseError "program not found"
// @Failure 500 {object} common.ResponseError "failed to create level"
// @Router /expert/levels [post]
func (h *LevelHandler) CreateLevelHandler(ctx *gin.Context) {

	// Bind JSON body vào CreateLevelRequest
	var req dtos.CreateLevelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(fmt.Sprintf("invalid level data: %v", err)))
		return
	}

	// Gọi service
	createdLevel, err := h.levelService.CreateLevel(ctx, &req)
	if err != nil {
		switch err.Error() {
		case "program not found":
			ctx.JSON(http.StatusNotFound, common.NewResponseError("program not found"))
		case "failed to parse program id":
			ctx.JSON(http.StatusBadRequest, common.NewResponseError("invalid program ID"))
		default:
			ctx.JSON(http.StatusInternalServerError, common.NewResponseError(fmt.Sprintf("failed to create level: %v", err)))
		}
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseNormal("level created successfully", createdLevel))
}