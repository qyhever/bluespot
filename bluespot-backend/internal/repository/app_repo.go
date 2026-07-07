package repository

import (
	"bluespot/internal/model"
)

type AppRepository interface {
	GetHelloInfo(param *model.GetHelloInfoRequest) (*model.GetHelloInfoResponse, error)
}
