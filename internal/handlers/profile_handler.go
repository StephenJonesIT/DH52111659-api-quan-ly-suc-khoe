package handlers

import (
	"DH52111659-api-quan-ly-suc-khoe/common"
	"DH52111659-api-quan-ly-suc-khoe/config"
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/services"
	"DH52111659-api-quan-ly-suc-khoe/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService services.ProfileService
}

func NewProfileHandler(profileService services.ProfileService) *ProfileHandler{
	return &ProfileHandler{profileService: profileService}
}

//CreateProfile godoc
//@Summary Create a new profile
//@Description Create user profile with file image and json profile
//@Tags Profile
//@Accept multipart/form-data
//@Produce json
//@Param Authorization header string true "Bearer Token"
//@Param image formData file false "Profile image file (max 10MB)"
//@Param metadata formData string true "Json body for profile"
//@Success 201 {object} common.ResponseNormal{data=object} "Profile created successfully"
//@Failure 400 {object} common.ResponseError "invalid request form-data"
//@Failure 401 {object} common.ResponseError "invalid token"
//@Failure 500 {object} common.ResponseError "Internal server error"
//@Router /profile [post] 
func(h *ProfileHandler) CreateProfileHandler(ctx *gin.Context) { 
	var profileRequest models.Profile

	avatarURL, err, isUploadFile := utils.HandleFileUpload(ctx, "image", config.AppConfig.UploadDir)
	if err != nil && isUploadFile{
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	if err := utils.UnmarshalFormValue(ctx, "metadata", &profileRequest); err != nil {
		utils.HandleFileDeleted(avatarURL, config.AppConfig.UploadDir)
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	if err := common.ValidateRequest(profileRequest); err != nil {
		utils.HandleFileDeleted(avatarURL, config.AppConfig.UploadDir)
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	profileRequest.AvatarURL = avatarURL

	profile, err := h.profileService.CreateProfile(ctx, &profileRequest)
	if err != nil {
		utils.HandleFileDeleted(avatarURL, config.AppConfig.UploadDir)
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}
	
	ctx.JSON(http.StatusCreated, common.NewResponseNormal("Profile created successfully", profile))
} 

// UpdateProfile godoc
//@Summary Update user profile
//@Description Update user profile with file image and json profile
//@Tags Profile
//@Accept multipart/form-data
//@Produce json
//@Param Authorization header string true "Bearer Token"
//@Param id path string true "User ID"
//@Param image formData file false "Profile image file (max 10MB)"
//@Param metadata formData string true "Json body for profile"
//@Success 200 {object} common.ResponseNormal{data=object} "Profile updated successfully"
//@Failure 400 {object} common.ResponseError "invalid request form-data"
//@Failure 401 {object} common.ResponseError "invalid token"
//@Failure 404 {object} common.ResponseError "Profile not found"
//@Failure 500 {object} common.ResponseError "Internal server error"
// @Router /profile/{id} [put]
func(h *ProfileHandler) UpdateProfileHandler(ctx *gin.Context) {
	var updateProfileRequest models.Profile

	avatarURL, err, isUploadFile := utils.HandleFileUpload(ctx, "image", config.AppConfig.UploadDir)
	if err != nil && isUploadFile{
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return 
	}

	if err := utils.UnmarshalFormValue(ctx, "metadata", &updateProfileRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	if err := common.ValidateRequest(updateProfileRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewResponseError(err.Error()))
		return
	}

	cond := ctx.Param("id")
	deleteURL := updateProfileRequest.AvatarURL
	
	log.Print(deleteURL)
	updateProfileRequest.AvatarURL = avatarURL
	 // Store the URL to delete later if needed

	profile, err := h.profileService.UpdateProfile(ctx, cond, &updateProfileRequest); 
	if err != nil {
		utils.HandleFileDeleted(avatarURL, config.AppConfig.UploadDir)
		ctx.JSON(http.StatusInternalServerError, common.NewResponseError(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, common.NewResponseNormal("Profile updated successfully", profile))

	// Clean up the uploaded file after use
	utils.HandleFileDeleted(deleteURL, config.AppConfig.UploadDir)
}
