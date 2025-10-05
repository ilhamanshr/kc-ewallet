package repository

import (
	"context"
	"database/sql"
	"kc-ewallet/domains/repository/postgres"
)

//go:generate mockgen -destination=mocks/mock_repository.go -source=repository.go IRepository,INats,IInternalService
type IRepository interface {
	// TX
	WithTx(tx *sql.Tx) *postgres.Queries

	// User
	CreateUser(ctx context.Context, arg postgres.CreateUserParams) (int32, error)
	GetUserByIDLock(ctx context.Context, id int32) (postgres.User, error)
	GetUserByUsername(ctx context.Context, username string) (postgres.User, error)
	UpdateUserBalanceByID(ctx context.Context, arg postgres.UpdateUserBalanceByIDParams) error

	// Transaction
	CreateTransaction(ctx context.Context, arg postgres.CreateTransactionParams) (int32, error)
}
