package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID    int32  `json:"user_id"`    // user id on auth service
	RoleGroup string `json:"role_group"` // this is the actual role
	RoleID    string `json:"role_id"`
	Role      string `json:"role"` // this actually means platorm, to be deprecated
	Platform  string `json:"platform"`
	FullName  string `json:"full_name"`
}
