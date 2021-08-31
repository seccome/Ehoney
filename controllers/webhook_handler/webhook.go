package webhook_handler

import (
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type WebHookConfigPayload struct{
	WebHook  string   `gorm:"not null;size:256" form:"WebHook" json:"WebHook" binding:"required"` //钉钉接口
}

// UpdateWebHookConfig 设置钉钉webhook
// @Summary 设置钉钉webhook
// @Description 设置钉钉webhook
// @Tags 系统设置
// @Produce application/json
// @Accept application/json
// @Param WebHook body WebHookConfigPayload true "WebHook"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/webhook [put]
func UpdateWebHookConfig(c *gin.Context) {
	appG := app.Gin{C: c}

	var payload WebHookConfigPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil{
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	viper.Set("app.webhook", payload.WebHook)
	viper.WriteConfig()
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetWebHookConfig 获取钉钉webhook
// @Summary 获取钉钉webhook
// @Description 获取钉钉webhook
// @Tags 系统设置
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/webhook [get]
func GetWebHookConfig(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.SUCCESS, configs.GetSetting().App.WebHook)
}
