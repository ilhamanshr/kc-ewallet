package jwt

import (
	"kc-ewallet/internals/errors"

	goerrors "errors"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyHMACToken(tokenStr string, secret string, signMethod *jwt.SigningMethodHMAC) (claims map[string]any, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = errors.BadRequest.New("Signing method invalid")
			return nil, err
		} else if method != signMethod {
			err = errors.BadRequest.New("Signing method invalid")
			return nil, err
		}

		return []byte(secret), nil
	})

	if goerrors.Is(err, jwt.ErrTokenExpired) {
		err = errors.Unauthorized.New("expired token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		err = errors.Unauthorized.New("invalid token")
		return
	}

	return
}
