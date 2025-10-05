package controller

import (
	"kc-ewallet/domains/usecase"
	requesthelper "kc-ewallet/internals/helpers/request"
	"kc-ewallet/protocols/http/request"
	"kc-ewallet/protocols/http/response"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	usecase usecase.ITransactionUsecase
}

func NewTransactionController(usecase usecase.ITransactionUsecase) *TransactionController {
	return &TransactionController{
		usecase: usecase,
	}
}

func (ctl *TransactionController) CreateCreditTransaction(ctx *gin.Context) {
	reqHelper := requesthelper.InitRequest(ctx)

	body := request.CreateCreditTransactionRequest{
		UserID: reqHelper.Auth.UserID,
	}
	if err := reqHelper.SetPostParams(&body); err != nil {
		return
	}

	transactionID, newBalance, err := ctl.usecase.CreateCreditTransaction(ctx.Request.Context(), body)
	if err != nil {
		response.RespondError(ctx, err)
		return
	}

	response.RespondSuccess(ctx, response.NewCreateCreditTransactionResponse(transactionID, newBalance), "success")
}

func (ctl *TransactionController) CreateDebitTransaction(ctx *gin.Context) {
	reqHelper := requesthelper.InitRequest(ctx)

	body := request.CreateDebitTransactionRequest{
		UserID: reqHelper.Auth.UserID,
	}
	if err := reqHelper.SetPostParams(&body); err != nil {
		return
	}

	transactionID, newBalance, err := ctl.usecase.CreateDebitTransaction(ctx.Request.Context(), body)
	if err != nil {
		response.RespondError(ctx, err)
		return
	}

	response.RespondSuccess(ctx, response.NewCreateDebitTransactionResponse(transactionID, newBalance), "success")
}
