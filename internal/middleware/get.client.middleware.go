package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

type ClientInfo struct {
	IP        string
	UserAgent string
}

func (m *Middleware) GetClientInfo(c *gin.Context) (*ClientInfo, error) {
	ip := c.ClientIP()
	if ip == "" {
		return nil, errors.New("unable to determine client IP")
	}

	userAgent := strings.TrimSpace(c.GetHeader("User-Agent"))
	if userAgent == "" {
		userAgent = "Unknown Device"
	}

	return &ClientInfo{
		IP:        ip,
		UserAgent: userAgent,
	}, nil
}
