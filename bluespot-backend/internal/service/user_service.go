package service

import (
	"bluespot/internal/model"
	"bluespot/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetCurrentUserInfo() (*model.UserInfoResponse, error) {
	return s.repo.GetCurrentUserInfo()
}
