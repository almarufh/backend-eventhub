package middleware

import (
	"backend/EventHub/pkg"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) GetPayload(c *gin.Context) (*Payload, error) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		log.Printf("[GET PAYLOAD] Error: missing Authorization header from IP: %s\n", c.ClientIP())
		return nil, errors.New("Authorization header is required")
	}

	result := strings.Split(bearer, " ")
	if len(result) != 2 {
		log.Printf("[GET PAYLOAD] Error: malformed Authorization header length (%d parts)\n", len(result))
		return nil, errors.New("Invalid authorization header format")
	}

	if result[0] != "Bearer" {
		log.Printf("[GET PAYLOAD] Error: invalid scheme '%s', expected 'Bearer'\n", result[0])
		return nil, errors.New("Authorization scheme must be Bearer")
	}

	var token pkg.JWTClaims
	err := token.DecodeToken(result[1])
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			log.Printf("[GET PAYLOAD] Warning: token validation failed: %v\n", err)
			return nil, errors.New("Token is invalid or expired")
		}

		log.Printf("[GET PAYLOAD] Critical Error: unexpected error while decoding token: %v\n", err)
		return nil, errors.New("Internal server error")
	}

	var expiresInSec int64 = 0
	if token.ExpiresAt != nil {
		diff := time.Until(token.ExpiresAt.Time)
		if diff > 0 {
			expiresInSec = int64(diff.Seconds())
		}
	}

	return &Payload{
		ID:        token.ID,
		Role:      token.Role,
		ExpiresIn: expiresInSec,
	}, nil
}
