package middleware

import (
	"kc-ewallet/internals/errors"
	"kc-ewallet/protocols/http/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type PagePermission string

const (
	UserPage        PagePermission = "user"
	TransactionPage PagePermission = "transaction"
)

var (
	ErrInvalidRolePermissions error = errors.Unauthorized.New("invalid role permission")
)

func CheckPermission(pages []PagePermission, opts ...middlewareOptionFn) gin.HandlerFunc {
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

		actor, err := NewActorFromContext(c)
		if err != nil {
			response.RespondError(c, ErrInvalidRolePermissions)
			c.Abort()
			return
		}

		var permitted bool

		for _, page := range pages {
			permitted = actor.IsPermit(page)
			if permitted {
				break
			}
		}

		if !permitted {
			response.RespondError(c, ErrInvalidRolePermissions)
			c.Abort()
		}
	}
}
