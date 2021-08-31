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

// UpdateHarborConfig 设置harbor镜像源
// @Summary 设置harbor镜像源
// @Description 设置harbor镜像源
// @Tags 系统设置
// @Produce application/json
// @Accept application/json
// @Param HarborURL body HarborConfigPayload true "HarborURL"
// @Param HarborProjectName body HarborConfigPayload true "HarborProjectName"
// @Param AuthenticateUser body HarborConfigPayload true "AuthenticateUser"
// @Param AuthenticatePass body HarborConfigPayload true "AuthenticatePass"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/harbor [put]
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

// GetHarborConfig 获取harbor镜像源
// @Summary 获取harbor镜像源
// @Description 获取harbor镜像源
// @Tags 系统设置
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/harbor [get]
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

// TestHarborConnection 测试harbor连接
// @Summary 测试harbor连接
// @Description 测试harbor连接
// @Tags 系统设置
// @Produce application/json
// @Accept application/json
// @Param HarborURL body HarborConfigPayload true "HarborURL"
// @Param HarborProjectName body HarborConfigPayload true "HarborProjectName"
// @Param AuthenticateUser body HarborConfigPayload true "AuthenticateUser"
// @Param AuthenticatePass body HarborConfigPayload true "AuthenticatePass"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.Images
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/harbor/health [get]
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
