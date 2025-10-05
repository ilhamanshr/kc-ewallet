package transaction

import (
	"context"
	"database/sql"
	"kc-ewallet/constants"
	"kc-ewallet/domains/repository"
	"kc-ewallet/domains/repository/postgres"
	"kc-ewallet/internals/errors"
	log_color "kc-ewallet/internals/helpers/color"
	"kc-ewallet/protocols/http/request"

	goerrors "errors"

	"go.opentelemetry.io/otel/trace"
)

type transactionUscase struct {
	db         *sql.DB
	repository repository.IRepository
	trace      trace.Tracer
}

func NewTransactionUsecase(
	db *sql.DB,
	repository repository.IRepository,
	trace trace.Tracer,
) *transactionUscase {
	return &transactionUscase{
		db:         db,
		repository: repository,
		trace:      trace,
	}
}

func (t *transactionUscase) CreateCreditTransaction(ctx context.Context, request request.CreateCreditTransactionRequest) (int32, float64, error) {
	var (
		tx  *sql.Tx
		err error
	)

	// Begin transaction
	if t.db != nil {
		tx, err = t.db.Begin()
		if err != nil {
			log_color.PrintRedf("Failed to begin transaction: %v\n", err)
			return 0, 0, errors.InternalServer.NewWithUserMsg(err, "failed to begin transaction")
		}
	}

	// Ensure to commit or rollback transaction at the end
	defer func() {
		if tx == nil {
			log_color.PrintRedf("Transaction is nil, cannot rollback\n")
			return
		}
		if err != nil {
			errRollback := tx.Rollback()
			log_color.PrintRedf("Transaction rollback due to error: %v, rollback error: %v\n", err, errRollback)
			return
		}
		tx.Commit()
	}()

	// Use transaction if available
	query := t.repository
	if tx != nil {
		query = t.repository.WithTx(tx)
	}

	// Lock the row for update
	user, err := query.GetUserByIDLock(ctx, request.UserID)
	if err != nil {
		log_color.PrintRedf("error get user by id: %v", err)
		if goerrors.Is(err, sql.ErrNoRows) {
			return 0, 0, errors.NotFound.NewWithUserMsg(err, "user not found")
		}
		return 0, 0, errors.InternalServer.NewWithUserMsg(err, "failed to get user by id")
	}

	// Update balance
	newBalance := user.Balance + request.Amount
	if err := query.UpdateUserBalanceByID(ctx, postgres.UpdateUserBalanceByIDParams{
		ID:      user.ID,
		Balance: newBalance,
	}); err != nil {
		return 0, 0, errors.InternalServer.NewWithUserMsg(err, "failed to update balance")
	}

	// Create transaction record
	transactionID, err := query.CreateTransaction(ctx, postgres.CreateTransactionParams{
		UserID: sql.NullInt32{Int32: request.UserID, Valid: true},
		Amount: request.Amount,
		Type:   constants.TransactionTypeCredit,
	})

	return transactionID, newBalance, nil
}

func (t *transactionUscase) CreateDebitTransaction(ctx context.Context, request request.CreateDebitTransactionRequest) (int32, float64, error) {
	var (
		tx  *sql.Tx
		err error
	)

	// Begin transaction
	if t.db != nil {
		tx, err = t.db.Begin()
		if err != nil {
			log_color.PrintRedf("Failed to begin transaction: %v\n", err)
			return 0, 0, errors.InternalServer.NewWithUserMsg(err, "failed to begin transaction")
		}
	}

	// Ensure to commit or rollback transaction at the end
	defer func() {
		if tx == nil {
			log_color.PrintRedf("Transaction is nil, cannot rollback\n")
			return
		}
		if err != nil {
			errRollback := tx.Rollback()
			log_color.PrintRedf("Transaction rollback due to error: %v, rollback error: %v\n", err, errRollback)
			return
		}
		tx.Commit()
	}()

	// Use transaction if available
	query := t.repository
	if tx != nil {
		query = t.repository.WithTx(tx)
	}

	// Lock the row for update
	user, err := query.GetUserByIDLock(ctx, request.UserID)
	if err != nil {
		log_color.PrintRedf("error get user by id: %v", err)
		if goerrors.Is(err, sql.ErrNoRows) {
			return 0, 0, errors.NotFound.NewWithUserMsg(err, "user not found")
		}
		return 0, 0, errors.InternalServer.NewWithUserMsg(err, "failed to get user by id")
	}

	// Check if balance is sufficient
	if user.Balance < request.Amount {
		return 0, 0, errors.BadRequest.NewWithUserMsg(nil, "Insufficient funds")
	}

	// Update balance
	newBalance := user.Balance - request.Amount
	if err := query.UpdateUserBalanceByID(ctx, postgres.UpdateUserBalanceByIDParams{
		ID:      user.ID,
		Balance: newBalance,
	}); err != nil {
		return 0, 0, errors.InternalServer.NewWithUserMsg(err, "failed to update balance")
	}

	// Create transaction record
	transactionID, err := query.CreateTransaction(ctx, postgres.CreateTransactionParams{
		UserID: sql.NullInt32{Int32: request.UserID, Valid: true},
		Amount: request.Amount,
		Type:   constants.TransactionTypeDebit,
	})

	return transactionID, newBalance, nil
}
