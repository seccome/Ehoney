package webhook_handler

import (
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type WebHookConfigPayload struct {
	WebHook string `gorm:"not null;size:256" form:"WebHook" json:"WebHook" binding:"required"` //钉钉接口
}

func UpdateWebHookConfig(c *gin.Context) {
	appG := app.Gin{C: c}

	var payload WebHookConfigPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	viper.Set("app.webhook", payload.WebHook)
	viper.WriteConfig()
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetWebHookConfig(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.SUCCESS, configs.GetSetting().App.WebHook)
}
