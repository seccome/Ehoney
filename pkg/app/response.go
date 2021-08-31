package app

import (
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, Code int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: Code,
		Msg:  GetMsg(Code),
		Data: data,
	})
	return
}
