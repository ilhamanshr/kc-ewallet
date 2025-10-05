package jwt

import (
	"context"
	"errors"
	"fmt"
	"kc-ewallet/configurations"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type IJWTHelper interface {
	GetClaim(tokenString string) (*JsonWebTokenClaims, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetToken(ctx context.Context) (string, error)
}

type jwtHelper struct {
	expirationTimeInMinute int
	signingKey             string
	issuer                 string
}

type JsonWebTokenClaims struct {
	UserID string `json:"userId"`
	Role   string `json:"roleId"`
	NIP    int    `json:"nip"`
	jwt.RegisteredClaims
}

func NewJWTHelper(jwtConfiguration configurations.IJWTConfiguration) *jwtHelper {
	return &jwtHelper{
		expirationTimeInMinute: jwtConfiguration.GetExpireInMinute(),
		signingKey:             jwtConfiguration.GetSigningKey(),
		issuer:                 jwtConfiguration.GetIssuer(),
	}
}

func (j *jwtHelper) GetClaim(tokenString string) (*JsonWebTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JsonWebTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		} else if method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}

		return []byte(j.signingKey), nil

	}, jwt.WithExpirationRequired(), jwt.WithIssuer(j.issuer))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claim, ok := token.Claims.(*JsonWebTokenClaims)

	if !ok {
		return nil, ErrInvalidToken
	}

	return claim, nil
}

func (j *jwtHelper) GetToken(ctx context.Context) (string, error) {
	// Extract metadata from context (which includes cookies in gRPC calls)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata found")
	}

	// Look for the 'cookie' field in metadata (which contains cookies)

	cookieHeaders := md["cookie"]
	if len(cookieHeaders) == 0 {
		cookieHeaders = md["grpcgateway-cookie"]
		// return "ctx", fmt.Errorf("no cookie found in metadata")
	}

	// Iterate through cookies to find the token
	for _, cookieHeader := range cookieHeaders {
		// Parse cookies from the header
		cookies := strings.Split(cookieHeader, ";")
		for _, cookie := range cookies {
			cookie = strings.TrimSpace(cookie)
			parts := strings.SplitN(cookie, "=", 2)
			if len(parts) != 2 {
				continue
			}
			cookieName := parts[0]
			cookieValue := parts[1]

			// Check if the cookie contains the access token (e.g., cookie name "access_token")
			if cookieName == "access_token" {
				// You can validate the token or attach it to the context for further use
				return cookieValue, nil
			}
		}
	}

	// Look for the 'authorization' field in metadata
	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		authHeaders = md["grpcgateway-authorization"]
		if len(authHeaders) == 0 {
			return "", fmt.Errorf("no authorization header found")
		}
	}

	// Iterate through the authorization headers
	for _, authHeader := range authHeaders {
		// Ensure the header follows the format "Bearer <token>"
		if strings.HasPrefix(authHeader, "Bearer ") {
			// Extract the token from the header
			token := strings.TrimPrefix(authHeader, "Bearer ")
			return strings.TrimSpace(token), nil
		}
	}

	return "", fmt.Errorf("no access token found in cookies")
}

func (j *jwtHelper) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.signingKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
