package routes

import (
	"kc-ewallet/constants"
	rate_limit "kc-ewallet/internals/helpers/rate_limiter"
	"kc-ewallet/protocols/http/controller"
	"kc-ewallet/protocols/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTransactionRoutes(router *gin.Engine, jwtSigningKey string, ctrl *controller.TransactionController) {
	v1RouterGroup := router.Group(constants.ApiV1BasePath)
	v1RouterGroup.Use(
		middleware.AuthorizeToken(
			jwtSigningKey,
			middleware.RegisterHandlers(
				map[string]bool{
					"CreateCreditTransaction": true,
					"CreateDebitTransaction":  true,
				},
			),
		),
		middleware.CheckRateLimit(
			middleware.NewRateLimiter(rate_limit.NewCacheService(), []string{}),
			middleware.RegisterHandlers(
				map[string]bool{
					"CreateCreditTransaction": true,
					"CreateDebitTransaction":  true,
				},
			),
		),
	)

	TransactionV1Routes(v1RouterGroup, ctrl)
}

func TransactionV1Routes(v1Router *gin.RouterGroup, ctrl *controller.TransactionController) {
	routes := v1Router.Group(constants.TransactionPath)

	routes.POST("/credit", ctrl.CreateCreditTransaction)
	routes.POST("/debit", ctrl.CreateDebitTransaction)
}
