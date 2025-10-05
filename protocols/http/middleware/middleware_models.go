package middleware

import (
	"context"
	"kc-ewallet/internals/errors"
	"kc-ewallet/internals/helpers/jwt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Actor struct {
	RoleID         uuid.UUID
	Role           string
	Platform       string
	RoleGroup      string
	FullName       string
	UserID         int32
	ReferenceID    uuid.UUID
	PermissionPage []string
	OriginToken    string
	CompanyID      uuid.UUID
}

func (a *Actor) IsPermit(page PagePermission) bool {
	for _, permit := range a.PermissionPage {
		if permit == string(page) {
			return true
		}
	}
	return false
}

func NewActorFromContext(ctx context.Context) (*Actor, error) {
	// Try to get actor using type as key first
	actor, ok := ctx.Value(ActorKey{}).(Actor)
	if ok {
		return &actor, nil
	}

	// If not found, try to get actor using string as key
	actor, ok = ctx.Value("actor").(Actor)
	if !ok {
		return nil, errors.BadRequest.New("actor not found")
	}

	return &actor, nil
}

func (a *Actor) SetToContext(c *gin.Context) {
	c.Set("actor", *a)
}

func (a *Actor) SetActorFromClaims(claims map[string]any) error {
	if userID, ok := claims["user_id"].(float64); ok {
		a.UserID = int32(userID)
	}

	if fullName, ok := claims["full_name"].(string); ok {
		a.FullName = fullName
	}

	if roleIdClaims, ok := claims["role_id"]; ok {
		if roleId, err := uuid.Parse(roleIdClaims.(string)); err == nil {
			a.RoleID = roleId
		}
	}

	platform := claims["role"].(string)
	a.Role = platform
	a.Platform = platform

	roleGroup := claims["role_group"].(string)
	a.RoleGroup = roleGroup

	jwt.ParsePermission(claims["permission_page"], &a.PermissionPage)

	return nil
}

func getHandlerNameFromGinContext(handlerList []string) string {
	lastHandler := handlerList[len(handlerList)-1]
	handlerNameSplit := strings.Split(lastHandler, ".")
	handlerName := handlerNameSplit[len(handlerNameSplit)-1]
	return strings.ReplaceAll(handlerName, "-fm", "")
}
