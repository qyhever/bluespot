package controller

import (
	"errors"

	"bluespot/internal/model"
	"bluespot/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(
	authService *service.AuthService,
) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login godoc
// @Summary 管理员登录
// @Description 使用管理员用户名和密码换取 accessToken 与 refreshToken。业务失败同样返回 HTTP 200，通过 code/message 区分。
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.AuthLoginRequest true "登录参数"
// @Success 200 {object} SwaggerAuthLoginResponse
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var param model.AuthLoginRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		ResponseFailed(c, CodeInvalidParam)
		return
	}

	result, err := ac.authService.AuthLogin(c.Request.Context(), param)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			ResponseFailed(c, CodeInvalidPassword)
			return
		}
		ResponseFailed(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, result)
}

// RefreshToken godoc
// @Summary 刷新管理员 Token
// @Description 使用 refreshToken 换取新的 accessToken 与 refreshToken。业务失败同样返回 HTTP 200，通过 code/message 区分。
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.AuthRefreshTokenRequest true "刷新参数"
// @Success 200 {object} SwaggerAuthLoginResponse
// @Router /auth/refresh [post]
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var param model.AuthRefreshTokenRequest
	if err := c.ShouldBindJSON(&param); err != nil {
		ResponseFailed(c, CodeInvalidParam)
		return
	}

	result, err := ac.authService.AuthRefreshToken(c.Request.Context(), param)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			ResponseFailed(c, CodeInvalidToken)
			return
		}
		ResponseFailed(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, result)
}
