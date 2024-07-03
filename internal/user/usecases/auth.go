package usecases

import (
	"cleanarch/boiler/internal/user/domain"
	"cleanarch/boiler/internal/utils/logger"
	"context"
)

type AuthUseCases struct {
	l           logger.Interface
	authService AuthService
	jwtService  JwtService
	userService UserService
}

type AuthService interface {
	SignUp(ctx context.Context, user *domain.AddUserRequest, tenantId string) error
	GetAuthenticatedUser(ctx context.Context, user *domain.AddUserRequest) (*domain.UserResponse, error)
}
type JwtService interface {
	GenerateAccessToken(ctx context.Context, user *domain.UserResponse) (string, error)
	GenerateRefreshToken(ctx context.Context, user *domain.UserResponse, currentRefreshTokenID string) (string, error)
	ValidateAccessToken(ctx context.Context, accessToken string) (string, error)
	RefreshTokenAccess(ctx context.Context, refreshToken string) (string, string, error)
}

func NewAuthUseCases(l logger.Interface, authService AuthService, jwtService JwtService, userService UserService) *AuthUseCases {
	return &AuthUseCases{
		l:           l,
		authService: authService,
		jwtService:  jwtService,
		userService: userService,
	}
}

func (a *AuthUseCases) Login(ctx context.Context, user *domain.AddUserRequest) (*domain.UserTokens, error) {
	dbUser, error := a.authService.GetAuthenticatedUser(ctx, user)
	if error != nil {
		return nil, error
	}
	accessToken, err := a.jwtService.GenerateAccessToken(ctx, dbUser)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.jwtService.GenerateRefreshToken(ctx, dbUser, "")
	if err != nil {
		return nil, err
	}

	return &domain.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
func (a *AuthUseCases) ValidateAccessToken(ctx context.Context, token string) (string, error) {
	return a.jwtService.ValidateAccessToken(ctx, token)
}

func (a *AuthUseCases) RefreshTokenAccess(ctx context.Context, token string) (*domain.UserTokens, error) {
	userId, tokenId, err := a.jwtService.RefreshTokenAccess(ctx, token)
	if err != nil {
		return nil, err
	}
	dbUser, error := a.userService.GetUserByID(ctx, userId)
	if error != nil {
		return nil, error
	}
	accessToken, err := a.jwtService.GenerateAccessToken(ctx, dbUser)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.jwtService.GenerateRefreshToken(ctx, dbUser, tokenId)
	if err != nil {
		return nil, err
	}

	return &domain.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (a *AuthUseCases) SignUp(ctx context.Context, user *domain.AddUserRequest, tenantId string) error {
	return a.authService.SignUp(ctx, user, tenantId)
}
