package agent_handler

import (
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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
