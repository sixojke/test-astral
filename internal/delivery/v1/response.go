package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sixojke/test-astral/pkg/logger"
)

type Response struct {
	Error    errorResponse `json:"error"`
	Response interface{}   `json:"response"`
	Data     interface{}   `json:"data"`
}

type errorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func errResponse(с *gin.Context, statusCode int, err, errResp string) {
	с.AbortWithStatusJSON(statusCode, Response{
		Error: errorResponse{
			Code: statusCode,
			Text: errResp,
		},
	})

	logger.Error(err)
}

func newResponse(c *gin.Context, statusCode int, data, response interface{}) {
	c.AbortWithStatusJSON(statusCode, Response{
		Response: response,
		Data:     data,
	})

	logger.Debugf("data: %v response: %v", data, response)
}
