package services

import (
	"cleanarch/boiler/internal/user/domain"
	"context"
)

type AuthSerivce struct {
	authRepository AuthRepository
}
type AuthRepository interface {
	SignUp(ctx context.Context, user *domain.AddUserRequest, tenantId string) error
	GetAuthenticatedUser(ctx context.Context, user *domain.AddUserRequest) (*domain.UserResponse, error)
}

func NewAuthService(repo AuthRepository) *AuthSerivce {
	return &AuthSerivce{
		authRepository: repo,
	}
}
func (a *AuthSerivce) SignUp(ctx context.Context, user *domain.AddUserRequest, tenantId string) error {
	return a.authRepository.SignUp(ctx, user, tenantId)
}
func (a *AuthSerivce) GetAuthenticatedUser(ctx context.Context, user *domain.AddUserRequest) (*domain.UserResponse, error) {
	return a.authRepository.GetAuthenticatedUser(ctx, user)
}
