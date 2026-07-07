package persistence

import (
	"errors"
	"strings"

	"bluespot/internal/config"
	"bluespot/internal/model"
	"bluespot/internal/repository"
)

var ErrConfigNotInitialized = errors.New("config not initialized")

type UserRepositoryImpl struct{}

func NewUserRepository() repository.UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) GetCurrentUserInfo() (*model.UserInfoResponse, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, ErrConfigNotInitialized
	}

	user := cfg.Auth.DefaultUser
	return &model.UserInfoResponse{
		UserID:   user.UserID,
		Username: strings.TrimSpace(user.Username),
		Nickname: strings.TrimSpace(user.Nickname),
	}, nil
}
