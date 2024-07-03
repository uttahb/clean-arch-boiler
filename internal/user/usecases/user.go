package usecases

import (
	"cleanarch/boiler/internal/user/domain"
	"cleanarch/boiler/internal/utils/logger"
	"context"
)

type UserUsecases struct {
	l           logger.Interface
	userService UserService
}

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error)
}

func NewUserUsecases(l logger.Interface, userService UserService) *UserUsecases {
	return &UserUsecases{
		l:           l,
		userService: userService,
	}
}

func (u *UserUsecases) GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error) {
	return u.userService.GetUserByID(ctx, id)
}
