package controller

import (
	"bluespot/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetCurrentUserInfo godoc
// @Summary 获取当前管理员信息
// @Description 获取当前登录管理员的基础信息，不返回密码。
// @Tags user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SwaggerUserInfoResponse
// @Router /user/info [get]
func (uc *UserController) GetCurrentUserInfo(c *gin.Context) {
	result, err := uc.userService.GetCurrentUserInfo()
	if err != nil {
		zap.L().Error("get current user info failed", zap.Error(err))
		ResponseFailed(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, result)
}
