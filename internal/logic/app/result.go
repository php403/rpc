package app

import (
	"github.com/gin-gonic/gin"
	"github.com/php403/im/pkg/log"
)

const (
	// OK ok
	OK = 0
	// RequestErr request error
	RequestErr = -400
	// ServerErr server error
	ServerErr = -500

	contextErrCode = "context/err/code"
)

type resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Errors(c *gin.Context, err *log.Error) {
	c.Set(contextErrCode, err.Code())
	c.JSON(200, resp{
		Code:    err.Code(),
		Message: err.Msg(),
	})
}

func Result(c *gin.Context, data interface{}, code int) {
	c.Set(contextErrCode, code)
	c.JSON(200, resp{
		Code: code,
		Data: data,
	})
}
