package services

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/models"
	"DH52111659-api-quan-ly-suc-khoe/internal/repositories"
	"context"
	"fmt"
)

type ProfileService interface {
	GetProfileByID(ctx context.Context, id string) (*models.Profile, error)
	CreateProfile(ctx context.Context, profileRequest *models.Profile) (*models.Profile, error)
}

type ProfileServiceImpl struct {
	repo repositories.ProfileRepository
}

func NewProfileServiceImpl(repo repositories.ProfileRepository) *ProfileServiceImpl {
	return &ProfileServiceImpl{
		repo: repo,
	}
}

func(s *ProfileServiceImpl) GetProfileByID(ctx context.Context, id string) (*models.Profile, error){
	profile, err := s.repo.GetProfileByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, fmt.Errorf("profile not found")
	}
	return profile, nil
}


func(s *ProfileServiceImpl) CreateProfile(ctx context.Context, profileRequest *models.Profile) (*models.Profile, error){
	//check if the profile exists 
	if _, err := s.repo.GetProfileByID(ctx, profileRequest.UserID.String()); err != nil {
		return nil, err
	}

	profile, err := s.repo.Create(ctx, profileRequest)
	if err != nil{
		return nil, err
	}

	return profile, nil
}