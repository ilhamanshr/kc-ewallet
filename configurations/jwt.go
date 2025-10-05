package configurations

import (
	"os"
	"strconv"
)

type jwtConfiguration struct {
	signingKey      string
	issuer          string
	expiresInMinute string
}

//go:generate mockgen -destination=mocks/mock_jwt.go -source=jwt.go IJWTConfiguration
type IJWTConfiguration interface {
	GetSigningKey() string
	GetIssuer() string
	GetExpireInMinute() int
}

func NewJWTConfiguration() *jwtConfiguration {
	return &jwtConfiguration{
		signingKey:      os.Getenv("JWT_SECRET"),
		issuer:          os.Getenv("JWT_ISSUER"),
		expiresInMinute: os.Getenv("JWT_EXPIRES_IN_MINUTE"),
	}
}

func (c *jwtConfiguration) GetSigningKey() string {
	return c.signingKey
}

func (c *jwtConfiguration) GetIssuer() string {
	if c.issuer == "" {
		return "kc-ewallet"
	}
	return c.issuer
}

func (c *jwtConfiguration) GetExpireInMinute() int {
	expiresInMinute, err := strconv.ParseInt(c.expiresInMinute, 10, 64)
	if err != nil {
		return 60 // default 60 minutes
	}

	return int(expiresInMinute)
}
