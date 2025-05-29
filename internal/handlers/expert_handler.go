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
// @Success      201  {object}  common.ResponseNormal{data=object} "Expert created successfully"
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

// GetExperts godoc
// @Summary Retrieve a list of experts
// @Description Returns a list of experts available in the system
// @Tags Expert
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Page number (default is 1)"
// @Param limit query int false "Number of experts per page (default is 10)"
// @Success 200 {object} common.ResponseNormal{data=[]models.Expert} "List of expert"
// @Failure 400 {object} common.ResponseError "Invalid query parameters"
// @Failure 401 {object} common.ResponseError "Invalid token"
// @Failure 403 {object} common.ResponseError "You do not have permission to access this resource"
// @Failure 500 {object} common.ResponseError "Internal server error"
// @Router /admin/experts [get]
func(h *ExpertHandler) GetExpertsHandler(ctx *gin.Context){ 
    var Paging common.Paging
    if err := ctx.ShouldBindQuery(&Paging); err != nil {
        ctx.JSON(http.StatusBadRequest, common.NewResponseError(common.ErrBadRequestShouldBind))
        return
    }

    expertsResponse, err := h.expertService.GetAllExperts(ctx, &Paging)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError,common.NewResponseError(err.Error()))
        return 
    }

    ctx.JSON(http.StatusOK, common.NewResponseNormal("Retrieve a list of experts successfully", expertsResponse))
}