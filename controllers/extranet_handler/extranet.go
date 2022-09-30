package extranet_handler

import (
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

type ExtranetConfigPayload struct {
	Extranet string `gorm:"null;size:256" form:"Extranet" json:"Extranet"` //外网IP
}

func UpdateExtranetConfig(c *gin.Context) {
	appG := app.Gin{C: c}

	var payload ExtranetConfigPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	viper.Set("app.extranet", payload.Extranet)
	viper.WriteConfig()
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetExtranetConfig(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.SUCCESS, configs.GetSetting().App.Extranet)
}
