package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-astral/domain"
	"github.com/sixojke/test-astral/pkg/logger"
)

const (
	authHeader = "Authorization"
)

func (h *Handler) middlewareAuth(c *gin.Context) {
	token, err := h.parseAuthHeader(c)
	if err != nil {
		errResponse(c, http.StatusUnauthorized, err.Error(), err.Error())

		return
	}
	logger.Debugf("token=%v", token)

	userId, err := h.service.User.GetUserIdByToken(token)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			errResponse(c, http.StatusUnauthorized, err.Error(), domain.ErrUserUnauthorized.Error())

			return
		}
	}

	c.Set("userId", userId)

	c.Next()
}

func (h *Handler) parseAuthHeader(c *gin.Context) (token string, err error) {
	header := c.GetHeader(authHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts) == 0 {
		return "", domain.ErrInvalidToken
	}

	return headerParts[1], nil
}
