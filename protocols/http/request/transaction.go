package request

type CreateCreditTransactionRequest struct {
	UserID int32   `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type CreateDebitTransactionRequest struct {
	UserID int32   `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
