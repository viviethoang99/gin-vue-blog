package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenExpired     = errors.New("Token has expired, please log in again")
	ErrTokenNotValidYet = errors.New("Token is not valid yet, please log in again")
	ErrTokenMalformed   = errors.New("Malformed token, please log in again")
	ErrTokenInvalid     = errors.New("Invalid token, please log in again")
)

type MyClaims struct {
	UserId  int   `json:"user_id"`
	RoleIds []int `json:"role_ids"`
	// UUID    string `json:"uuid"`
	jwt.RegisteredClaims
}

func GenToken(secret, issuer string, expireHour, userId int, roleIds []int) (string, error) {
	claims := MyClaims{
		UserId:  userId,
		RoleIds: roleIds,
		// UUID:    uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(secret, token string) (*MyClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

	if err != nil {
		switch vError, ok := err.(*jwt.ValidationError); ok {
		case vError.Errors&jwt.ValidationErrorMalformed != 0:
			return nil, ErrTokenMalformed
		case vError.Errors&jwt.ValidationErrorExpired != 0:
			return nil, ErrTokenExpired
		case vError.Errors&jwt.ValidationErrorNotValidYet != 0:
			return nil, ErrTokenNotValidYet
		default:
			return nil, ErrTokenInvalid
		}
	}

	if claims, ok := jwtToken.Claims.(*MyClaims); ok && jwtToken.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
