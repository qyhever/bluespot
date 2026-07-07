package repository

import "bluespot/internal/model"

type UserRepository interface {
	GetCurrentUserInfo() (*model.UserInfoResponse, error)
}
