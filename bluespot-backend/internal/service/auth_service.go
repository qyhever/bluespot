package service

import (
	"context"
	"strings"

	"bluespot/internal/config"
	"bluespot/internal/domain"
	"bluespot/internal/model"
	jwtpkg "bluespot/internal/pkg/jwt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

var (
	ErrInvalidCredentials  = domain.ErrInvalidAdminCredentials
	ErrInvalidRefreshToken = domain.ErrInvalidRefreshToken
)

func (s *AuthService) AuthLogin(ctx context.Context, req model.AuthLoginRequest) (*model.AuthLoginResponse, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	cfg := config.GetConfig()
	if cfg == nil {
		return nil, ErrInvalidCredentials
	}

	user := cfg.Auth.DefaultUser
	if user.UserID <= 0 ||
		strings.TrimSpace(user.Username) == "" ||
		user.Password == "" ||
		username != strings.TrimSpace(user.Username) ||
		password != user.Password {
		return nil, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := jwtpkg.GenToken(uint64(user.UserID))
	if err != nil {
		return nil, err
	}
	return &model.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) AuthRefreshToken(ctx context.Context, req model.AuthRefreshTokenRequest) (*model.AuthLoginResponse, error) {
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		return nil, ErrInvalidRefreshToken
	}

	claims, err := jwtpkg.ParseToken(refreshToken)
	if err != nil || !claims.IsRefreshToken() {
		return nil, ErrInvalidRefreshToken
	}

	accessToken, newRefreshToken, err := jwtpkg.GenToken(claims.UserID)
	if err != nil {
		return nil, err
	}
	return &model.AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
