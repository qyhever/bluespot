package controller

import "bluespot/internal/model"

type SwaggerErrorResponse struct {
	Code    int         `json:"code" example:"1001"`
	Message string      `json:"message" example:"请求参数错误"`
	Data    interface{} `json:"data"`
}

type SwaggerHelloInfoResponse struct {
	Code    int                        `json:"code" example:"1000"`
	Message string                     `json:"message" example:"success"`
	Data    model.GetHelloInfoResponse `json:"data"`
}

type SwaggerAuthLoginResponse struct {
	Code    int                     `json:"code" example:"1000"`
	Message string                  `json:"message" example:"success"`
	Data    model.AuthLoginResponse `json:"data"`
}

type SwaggerUserInfoResponse struct {
	Code    int                    `json:"code" example:"1000"`
	Message string                 `json:"message" example:"success"`
	Data    model.UserInfoResponse `json:"data"`
}
