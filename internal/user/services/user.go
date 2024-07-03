package services

import (
	"cleanarch/boiler/internal/user/domain"
	"cleanarch/boiler/internal/utils/logger"
	"context"
)

type UserService struct {
	l              logger.Interface
	userRepository UserRepository
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error)
}

func NewUserService(l logger.Interface, userRepository UserRepository) *UserService {
	return &UserService{
		l:              l,
		userRepository: userRepository,
	}
}
func (s *UserService) GetUserByID(ctx context.Context, id string) (*domain.UserResponse, error) {
	return s.userRepository.GetUserByID(ctx, id)
}
