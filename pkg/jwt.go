package pkg

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingKey     = errors.New("jwt key not found")
	ErrUnexpectedSign = errors.New("unexpected signing method")
)

type JWTClaims struct {
	ID   int32  `json:"user_id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTClaims(user_id int32, role string, expired int) *JWTClaims {
	return &JWTClaims{
		ID:        user_id,
		Role:      role,
		Issuer:    os.Getenv("JWT_ISSUER"),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expired) * time.Minute)),
	}
}

func (j *JWTClaims) GeneretToken() (string, error) {
	jwtKey := os.Getenv("JWT_SECREET_KEY")
	if jwtKey == "" {
		return "", ErrMissingKey
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, j)
	return token.SignedString([]byte(jwtKey))
}

func (j *JWTClaims) DecodeToken(tokenString string) error {
	jwtKey := os.Getenv("JWT_SECREET_KEY")
	if jwtKey == "" {
		return ErrMissingKey
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, j, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpectedSign, t.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return jwt.ErrTokenInvalidClaims
	}

	iss, err := jwtToken.Claims.GetIssuer()
	if err != nil {
		return err
	}
	if iss != os.Getenv("JWT_ISSUER") {
		return jwt.ErrTokenInvalidIssuer
	}

	return nil
}
