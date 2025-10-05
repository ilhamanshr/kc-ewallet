package middleware

import (
	"context"
	"kc-ewallet/internals/errors"
	"kc-ewallet/protocols/http/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	jwtHelper "kc-ewallet/internals/helpers/jwt"
)

type ActorKey struct{}

const (
	Customer               string = "customer"
	Admin                  string = "admin"
	Lender                 string = "lender"
	RelationshipManager    string = "relationship_manager"
	CreditAnalyst          string = "credit_analyst"
	FieldOfficer           string = "field_officer"
	BearerScheme           string = "Bearer"
	AccessTokenCookieName  string = "access_token"
	RefreshTokenCookieName string = "refresh_token"
)

var (
	ErrUnauthorized = errors.NewExtError(errors.ExtErrorArg{
		Code:      http.StatusUnauthorized,
		IdMessage: "Harap login kembali",
		ErrCode:   "ER006",
		Err:       "invalid token",
	})
)

func AuthorizeToken(secret string, opts ...middlewareOptionFn) gin.HandlerFunc {
	return func(c *gin.Context) {
		opt := defaultMiddlewareOption()
		for _, o := range opts {
			o(opt)
		}

		if strings.Contains(c.FullPath(), "private") {
			return
		}

		handlerName := getHandlerNameFromGinContext(c.HandlerNames())

		registered := true

		if opt.registerHandlers != nil {
			registered = opt.registerHandlers[handlerName]
		}

		if opt.excludedHandlers != nil {
			if opt.excludedHandlers[handlerName] {
				registered = false
			}
		}

		if !registered {
			return
		}

		authCookie, err := c.Cookie(AccessTokenCookieName)
		if err == nil && authCookie != "" {
			c.Request.Header["Authorization"] = []string{"Bearer " + authCookie}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.RespondError(c, ErrUnauthorized)
			c.Abort()
			return
		}

		tokenString := authHeader[len(BearerScheme)+1:]

		if claims, err := jwtHelper.VerifyHMACToken(tokenString, secret, jwt.SigningMethodHS256); err != nil {
			response.RespondError(c, err)
			c.Abort()
			return
		} else {
			var actor Actor
			if errSetActor := actor.SetActorFromClaims(claims); errSetActor != nil {
				response.RespondError(c, errSetActor)
				c.Abort()
				return
			}
			actor.OriginToken = tokenString
			actor.SetToContext(c)

			ctxRequestWithActor := context.WithValue(c.Request.Context(), ActorKey{}, actor)
			c.Request = c.Request.WithContext(ctxRequestWithActor)
		}
	}
}
