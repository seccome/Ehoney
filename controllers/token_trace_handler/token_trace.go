package token_trace_handler

import (
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)


type TraceHostConfigPayload struct{
	TraceHost  string   `gorm:"not null;size:256" form:"TraceHost" json:"TraceHost" binding:"required"` //密签跟踪设置
}

// UpdateTraceHostConfig 设置密签跟踪URL
// @Summary 设置密签跟踪URL
// @Description 设置密签跟踪URL
// @Tags 系统设置
// @Produce application/json
// @Accept application/json
// @Param TraceHost body TraceHostConfigPayload true "TraceHost"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/token/trace [put]
func UpdateTraceHostConfig(c *gin.Context) {
	appG := app.Gin{C: c}

	var payload TraceHostConfigPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil{
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	viper.Set("app.tokentraceaddress", payload.TraceHost)
	viper.WriteConfig()
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetTraceHostConfig 获得密签跟踪URL
// @Summary 获得密签跟踪URL
// @Description 获得密签跟踪URL
// @Tags 系统设置
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/token/trace [get]
func GetTraceHostConfig(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.SUCCESS, configs.GetSetting().App.TokenTraceAddress)
}
