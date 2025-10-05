package controller

import (
	"kc-ewallet/domains/usecase"
	requesthelper "kc-ewallet/internals/helpers/request"
	"kc-ewallet/protocols/http/request"
	"kc-ewallet/protocols/http/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	usecase usecase.IUserUsecase
}

func NewUserController(usecase usecase.IUserUsecase) *UserController {
	return &UserController{
		usecase: usecase,
	}
}

func (ctl *UserController) RegisterUser(ctx *gin.Context) {
	reqHelper := requesthelper.InitRequest(ctx)

	var body request.RegisterUserRequest
	if err := reqHelper.SetPostParams(&body); err != nil {
		return
	}

	err := ctl.usecase.CreateUser(ctx.Request.Context(), body)
	if err != nil {
		response.RespondError(ctx, err)
		return
	}

	response.RespondSuccess(ctx, nil, "success")
}

func (ctl *UserController) Login(ctx *gin.Context) {
	reqHelper := requesthelper.InitRequest(ctx)

	var body request.LoginRequest
	if err := reqHelper.SetPostParams(&body); err != nil {
		return
	}

	accessToken, user, err := ctl.usecase.Login(ctx.Request.Context(), body)
	if err != nil {
		response.RespondError(ctx, err)
		return
	}

	res := response.NewLoginResponse(accessToken, *user)
	response.RespondSuccess(ctx, res, "success")
}

func (ctl *UserController) GetUserByID(ctx *gin.Context) {
	reqHelper := requesthelper.InitRequest(ctx)

	user, err := ctl.usecase.GetUserByID(ctx.Request.Context(), reqHelper.Auth.UserID)
	if err != nil {
		response.RespondError(ctx, err)
		return
	}

	res := response.NewGetUserByIDResponse(*user)
	response.RespondSuccess(ctx, res, "success")
}
