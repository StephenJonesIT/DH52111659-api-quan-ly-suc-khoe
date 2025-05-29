package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/config"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"DH52111659-api-quan-ly-suc-khoe/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExpertHandler struct {
	expertService services.ExpertService
}

func NewExpertHandler(service services.ExpertService) *ExpertHandler{
	return &ExpertHandler{expertService: service}
}

// CreateExpert godoc
// @Summary      Create a new expert
// @Description  Create expert profile with file image and json expert data
// @Tags         Expert
// @Accept       multipart/form-data
// @Produce      json
// @Security     ApiKeyAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        image         formData  file    false  "Expert image file (max 10MB)"
// @Param        metadata      formData  string  true   "Expert data in JSON format" 
// @Success      201  {object}  common.ResponseNormal{data=object}
// @Failure      400  {object}  common.ResponseError
// @Failure      401  {object}  common.ResponseError
// @Failure      403  {object}  common.ResponseError
// @Failure      500  {object}  common.ResponseError
// @Router       /admin/expert [post]
func(h *ExpertHandler) CreateExpertHandler(ctx *gin.Context) {
	var createExpertRequest models.ExpertCreate

	avatarURL, err, isUploadFile := utils.HandleFileUpload(ctx, "image", config.AppConfig.UploadDir);
	if err != nil && isUploadFile{
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return 
	}

	if err := utils.UnmarshalFormValue(ctx, "metadata", &createExpertRequest); err != nil {
		utils.HandleFileDeleted(avatarURL, config.AppConfig.UploadDir)
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return 
	}

	if err := common.ValidateRequest(createExpertRequest); err != nil{
		utils.HandleFileDeleted(avatarURL, config.AppConfig.UploadDir)
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return 
	}

	createExpertRequest.AvatarURL = avatarURL

	expert, err := h.expertService.CreateExpert(ctx, &createExpertRequest)
	if err != nil {
		utils.HandleFileDeleted(avatarURL, config.AppConfig.UploadDir)
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseNormal("Expert created successfully", expert))
}

