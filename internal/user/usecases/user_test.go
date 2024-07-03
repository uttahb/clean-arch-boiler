package usecases

// import (
// 	"context"
// 	"errors"
// 	"cleanarch/boiler/internal/user/domain"
// 	"cleanarch/boiler/internal/utils/logger"
// 	"testing"
// )

// func TestUserUsecases_GetUserByID(t *testing.T) {
// 	ctx := context.Background()
// 	userID := "user123"

// 	tests := []struct {
// 		name         string
// 		mockResponse *domain.UserResponse
// 		mockErr      error
// 		wantErr      bool
// 	}{
// 		{
// 			name:         "Success",
// 			mockResponse: &domain.UserResponse{ID: userID, Email: "johndoe@gmail.com"},
// 			mockErr:      nil,
// 			wantErr:      false,
// 		},
// 		{
// 			name:         "Error",
// 			mockResponse: nil,
// 			mockErr:      errors.New("failed to get user"),
// 			wantErr:      true,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			mockUserService := &mocks.UserService{}
// 			mockUserService.On("GetUserByID", ctx, userID).Return(tc.mockResponse, tc.mockErr)

// 			l := logger.NewLogger("")
// 			u := NewUserUsecases(l, mockUserService)

// 			resp, err := u.GetUserByID(ctx, userID)
// 			if tc.wantErr {
// 				assert.Error(t, err)
// 				return
// 			}

// 			assert.NoError(t, err)
// 			assert.Equal(t, tc.mockResponse, resp)
// 			mockUserService.AssertExpectations(t)
// 		})
// 	}
// }

// func TestNewUserUsecases(t *testing.T) {
// 	l := logger.NewMockLogger()
// 	mockUserService := &mocks.UserService{}

// 	u := NewUserUsecases(l, mockUserService)

// 	assert.NotNil(t, u)
// 	assert.Equal(t, l, u.l)
// 	assert.Equal(t, mockUserService, u.userService)
// }
