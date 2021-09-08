package extranet_handler

import (
    "decept-defense/pkg/app"
    "decept-defense/pkg/configs"
    "github.com/gin-gonic/gin"
    "github.com/spf13/viper"
    "net/http"
)

type ExtranetConfigPayload struct{
    Extranet  string   `gorm:"null;size:256" form:"Extranet" json:"Extranet"` //外网IP
}

// UpdateExtranetConfig 设置外网ip
// @Summary 设置外网ip
// @Description 设置外网ip
// @Tags 系统设置
// @Produce application/json
// @Accept application/json
// @Param Extranet body ExtranetConfigPayload true "Extranet"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/extranet [put]
func UpdateExtranetConfig(c *gin.Context) {
    appG := app.Gin{C: c}

    var payload ExtranetConfigPayload
    err := c.ShouldBindJSON(&payload)
    if err != nil{
        appG.Response(http.StatusOK, app.InvalidParams, nil)
        return
    }
    viper.Set("app.extranet", payload.Extranet)
    viper.WriteConfig()
    appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetExtranetConfig 获取外网ip
// @Summary 获取外网ip
// @Description 获取外网ip
// @Tags 系统设置
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/extranet [get]
func GetExtranetConfig(c *gin.Context) {
    appG := app.Gin{C: c}
    appG.Response(http.StatusOK, app.SUCCESS, configs.GetSetting().App.Extranet)
}

