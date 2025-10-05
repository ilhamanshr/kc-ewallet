package requesthelper

import (
	"kc-ewallet/internals/helpers/container"
	"kc-ewallet/protocols/http/middleware"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	UserID               int32
	AccessToken          string
	IsUsingInternalToken bool
	IsLoggedIn           bool
}

type ExtraData struct {
	// Platform string
}

type Request struct {
	*container.Request
	Auth
	ExtraData
}

func InitRequest(c *gin.Context) *Request {
	coreRequest := container.InitCoreRequest(c)
	r := Request{
		Request: coreRequest,
	}

	initAuth(c, &r)
	initExtraData(&r)
	return &r
}

func initAuth(c *gin.Context, r *Request) {
	if c == nil {
		return
	}

	actor, err := middleware.NewActorFromContext(c)
	if err == nil {
		r.Auth.UserID = actor.UserID
	}

	r.Auth.IsLoggedIn = false
	r.Auth.IsUsingInternalToken = false
	userID := c.GetString("user_id")
	r.Auth.AccessToken = c.GetString("AccessToken")

	if userID == "" {
		return
	}

	r.Auth.IsLoggedIn = true
	if userID == "-1" {
		r.Auth.IsUsingInternalToken = true
	}
}

func initExtraData(r *Request) {
	// r.ExtraData.Platform = platform
}
