package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-astral/pkg/logger"
)

type Response struct {
	Error    *errorResponse `json:"error,omitempty"`
	Response interface{}    `json:"response,omitempty"`
	Data     interface{}    `json:"data,omitempty"`
}

type errorResponse struct {
	Code int    `json:"code,omitempty"`
	Text string `json:"text,omitempty"`
}

func errResponse(с *gin.Context, statusCode int, err, errResp string) {
	с.AbortWithStatusJSON(statusCode, Response{
		Error: &errorResponse{
			Code: statusCode,
			Text: errResp,
		},
	})

	logger.Warn(err)
}

func newResponse(c *gin.Context, statusCode int, data, response interface{}) {
	c.JSON(statusCode, Response{
		Response: response,
		Data:     data,
	})

	logger.Debugf("data: %v response: %v", data, response)
}
