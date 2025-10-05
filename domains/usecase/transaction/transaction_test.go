package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"kc-ewallet/configurations"
	"kc-ewallet/domains/repository/postgres"
	log_color "kc-ewallet/internals/helpers/color"
	"kc-ewallet/protocols/http/request"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func setupTestDB(t *testing.T) *sql.DB {
	godotenv.Load(filepath.Join("..", "..", "..", ".env"))

	dbConfig := configurations.NewDatabaseWriter()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		dbConfig.GetUser(),
		dbConfig.GetPassword(),
		dbConfig.GetHost(),
		dbConfig.GetPort(),
		dbConfig.GetDBName(),
		dbConfig.GetDBParam(),
	)

	db, err := sql.Open("postgres", connStr)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)
	if err != nil {
		t.Fatal("failed to connect db:", err)
	}
	return db
}

func TestTransactionUsecase_ConcurrentTransactions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := postgres.New(db)
	usecase := NewTransactionUsecase(db, repo, nil)

	// create user
	userID, err := repo.CreateUser(ctx, postgres.CreateUserParams{
		Username: fmt.Sprintf("concurrency_test_%d", time.Now().UnixNano()),
		Password: "password",
	})
	assert.NoError(t, err)

	// concurrency params
	numThreads := 100
	numTransactions := 100
	amount := float64(1000)

	wg := sync.WaitGroup{}
	wg.Add(numThreads)

	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numTransactions; j++ {
				_, _, err := usecase.CreateCreditTransaction(ctx, request.CreateCreditTransactionRequest{
					UserID: userID,
					Amount: amount,
				})
				if err != nil {
					// fail test if any transaction error
					t.Errorf("credit failed: %v", err)
				}
			}
		}()
	}

	wg.Wait()

	// check final balance
	user, err := repo.GetUserByIDLock(ctx, userID)
	assert.NoError(t, err)

	expected := float64(numThreads * numTransactions * 1000)
	assert.Equal(
		t, expected,
		user.Balance,
		fmt.Sprintf("expected final balance %f, got %f", expected, user.Balance),
	)
	log_color.PrintGreenf("expected final balance %f, got %f", expected, user.Balance)

	// clean up
	defer func() {
		db.Exec("DELETE FROM transactions WHERE user_id = $1", userID)
		db.Exec("DELETE FROM users WHERE id = $1", userID)
	}()
}
