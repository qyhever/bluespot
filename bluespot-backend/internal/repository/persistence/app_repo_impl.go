package persistence

import (
	"bluespot/internal/model"
	"bluespot/internal/repository"
)

type AppRepositoryImpl struct{}

func NewAppRepository() repository.AppRepository {
	return &AppRepositoryImpl{}
}

func (r *AppRepositoryImpl) GetHelloInfo(req *model.GetHelloInfoRequest) (*model.GetHelloInfoResponse, error) {
	res := &model.GetHelloInfoResponse{
		Name: req.Name,
	}
	return res, nil
}
