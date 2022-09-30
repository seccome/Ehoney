package harbor_handler

import (
	"decept-defense/internal/harbor"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type HarborConfigPayload struct {
	HarborURL         string `gorm:"not null;size:256" form:"HarborURL" json:"HarborURL" binding:"required"`                 //URL
	HarborProjectName string `gorm:"not null;size:256" form:"HarborProjectName" json:"HarborProjectName" binding:"required"` //项目名称
	HarborAPIVersion  string `gorm:"not null;size:256" form:"HarborAPIVersion" json:"HarborAPIVersion"`
	AuthenticateUser  string `gorm:"not null;size:256" form:"AuthenticateUser"  json:"AuthenticateUser" binding:"required"` //用户名
	AuthenticatePass  string `gorm:"not null;size:256" form:"AuthenticatePass"  json:"AuthenticatePass" binding:"required"` //密码
}

func UpdateHarborConfig(c *gin.Context) {
	appG := app.Gin{C: c}

	var payload HarborConfigPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	viper.Set("harbor.HarborURL", payload.HarborURL)
	viper.Set("harbor.HarborProject", payload.HarborProjectName)
	viper.Set("harbor.user", payload.AuthenticateUser)
	viper.Set("harbor.password", payload.AuthenticatePass)
	err = viper.WriteConfig()
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	harbor.RefreshImages()
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetHarborConfig(c *gin.Context) {
	appG := app.Gin{C: c}
	var harborConfigPayload HarborConfigPayload
	harborConfigPayload.HarborAPIVersion = configs.GetSetting().Harbor.APIVersion
	harborConfigPayload.HarborURL = configs.GetSetting().Harbor.HarborURL
	harborConfigPayload.AuthenticatePass = configs.GetSetting().Harbor.Password
	harborConfigPayload.AuthenticateUser = configs.GetSetting().Harbor.User
	harborConfigPayload.HarborProjectName = configs.GetSetting().Harbor.HarborProject
	appG.Response(http.StatusOK, app.SUCCESS, harborConfigPayload)
}

func TestHarborConnection(c *gin.Context) {
	appG := app.Gin{C: c}

	var payload HarborConfigPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	viper.Set("harbor.HarborURL", payload.HarborURL)
	viper.Set("harbor.HarborProject", payload.HarborProjectName)
	viper.Set("harbor.user", payload.AuthenticateUser)
	viper.Set("harbor.password", payload.AuthenticatePass)

	header := map[string]string{
		"authorization": "Basic " + payload.AuthenticateUser + ":" + payload.AuthenticatePass,
	}
	requestURI := strings.Join([]string{payload.HarborURL, "api", "v2", "projects", payload.HarborProjectName, "repositories"}, "/")
	_, err = util.SendGETRequest(header, requestURI)

	if err != nil {
		zap.L().Info("harbor连接测试失败")
		appG.Response(http.StatusOK, app.ErrorConnectTest, nil)
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}
