package agent_handler

import (
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

type EngineUploadPayload struct {
	Version string `form:"version" json:"version" binding:"required"` //诱饵类型
}

// DownloadLinuxAgent 下载linux Agent
// @Summary 下载linux Agent
// @Description 下载linux Agent
// @Tags agent支持
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Router /api/v1/agent/linux [get]
func DownloadLinuxAgent(c *gin.Context) {
	appG := app.Gin{C: c}
	var URL = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + "agent/decept-agent.tar.gz"
	appG.Response(http.StatusOK, app.SUCCESS, URL)
}

func QueryLinuxEngineVersion(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.SUCCESS, configs.GetSetting().App.EngineVersion)
}

func UploadLinuxEngine(c *gin.Context) {
	appG := app.Gin{C: c}

	var payload EngineUploadPayload
	err := c.ShouldBind(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	if payload.Version == configs.GetSetting().App.EngineVersion {
		appG.Response(http.StatusOK, app.EngineVersionSame, nil)
		return
	}

	file, err := c.FormFile("file")

	if err != nil {
		appG.Response(http.StatusOK, app.ErrorFileUpload, nil)
		return
	}

	if !strings.HasSuffix(file.Filename, "tar.gz") {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	filePath := path.Join(util.WorkingPath(), "agent/engine/engine.tar.gz")

	_, err = os.Open(filePath)

	if err == nil {
		err = os.Remove(filePath)
		if err != nil {
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	zap.L().Info(fmt.Sprintf("%v", configs.GetSetting()))

	viper.Set("app.engineversion", payload.Version)
	err = viper.WriteConfig()
	zap.L().Info(fmt.Sprintf("%v", configs.GetSetting()))

	if err != nil {
		appG.Response(http.StatusOK, app.ErrorFileUpload, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, "")
}

func DownloadLinuxEngine(c *gin.Context) {
	appG := app.Gin{C: c}
	var URL = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + "agent/engine/engine.tar.gz"
	appG.Response(http.StatusOK, app.SUCCESS, URL)
}

// DownloadWindowsAgent 下载windows Agent
// @Summary 下载windows Agent
// @Description 下载windows Agent
// @Tags agent支持
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Router /api/v1/agent/windows [get]
func DownloadWindowsAgent(c *gin.Context) {
	appG := app.Gin{C: c}

	var URL string

	if configs.GetSetting().App.Extranet != "" {
		URL = "http:" + "//" + configs.GetSetting().App.Extranet + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + "agent/decept-agent-win.tar.gz"

	} else {
		URL = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + "agent/decept-agent-win.tar.gz"
	}

	appG.Response(http.StatusOK, app.SUCCESS, URL)
}
