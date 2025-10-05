package response

type CreateCreditTransactionResponse struct {
	TransactionID int32   `json:"transaction_id"`
	NewBalance    float64 `json:"new_balance"`
}

type CreateDebitTransactionResponse struct {
	TransactionID int32   `json:"transaction_id"`
	NewBalance    float64 `json:"new_balance"`
}

func NewCreateCreditTransactionResponse(transactionID int32, newBalance float64) CreateCreditTransactionResponse {
	return CreateCreditTransactionResponse{
		TransactionID: transactionID,
		NewBalance:    newBalance,
	}
}

func NewCreateDebitTransactionResponse(transactionID int32, newBalance float64) CreateDebitTransactionResponse {
	return CreateDebitTransactionResponse{
		TransactionID: transactionID,
		NewBalance:    newBalance,
	}
}
