package usecase

import (
	"context"
	"kc-ewallet/domains/repository/postgres"
	"kc-ewallet/protocols/http/request"
)

//go:generate mockgen -destination=mocks/mock_usecase.go -source=usecase.go IUserUsecase,ITransactionUsecase
type IUserUsecase interface {
	CreateUser(ctx context.Context, request request.RegisterUserRequest) error
	GetUserByID(ctx context.Context, userID int32) (*postgres.User, error)
	Login(ctx context.Context, request request.LoginRequest) (string, *postgres.User, error)
}

type ITransactionUsecase interface {
	CreateCreditTransaction(ctx context.Context, request request.CreateCreditTransactionRequest) (int32, float64, error)
	CreateDebitTransaction(ctx context.Context, request request.CreateDebitTransactionRequest) (int32, float64, error)
}

type GetUserByIDResponse struct {
	ID       int32   `json:"id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}
