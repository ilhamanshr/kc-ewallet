package user

import (
	"context"
	"errors"
	mock_configuration "kc-ewallet/configurations/mocks"
	mock_repository "kc-ewallet/domains/repository/mocks"
	"kc-ewallet/protocols/http/request"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRepo := mock_repository.NewMockIRepository(ctrl)
	mockJwtConfig := mock_configuration.NewMockIJWTConfiguration(ctrl)
	sqlDB, _, _ := sqlmock.New()

	usecase := NewUserUsecase(sqlDB, mockRepo, mockJwtConfig, nil)

	testCases := []struct {
		name          string
		request       request.RegisterUserRequest
		mock          func()
		expectedError error
	}{
		{
			name: "success",
			request: request.RegisterUserRequest{
				Username: "luffy",
			},
			mock: func() {
				mockRepo.EXPECT().CreateUser(gomock.Any(), "luffy").Return(int32(1), nil)
			},
		},
		{
			name: "should error when duplicate username",
			request: request.RegisterUserRequest{
				Username: "zore",
			},
			mock: func() {
				mockRepo.EXPECT().CreateUser(gomock.Any(), "zore").Return(int32(1), &pq.Error{Code: "23505"})
			},
			expectedError: errors.New("username already exists"),
		},
		{
			name: "should error when db fails",
			request: request.RegisterUserRequest{
				Username: "sanji",
			},
			mock: func() {
				mockRepo.EXPECT().CreateUser(gomock.Any(), "sanji").Return(int32(1), assert.AnError)
			},
			expectedError: errors.New("failed to create user"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock()

			if err := usecase.CreateUser(context.Background(), tc.request); err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			}
		})
	}
}
