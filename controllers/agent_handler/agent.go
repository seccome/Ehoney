package agent_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/agent_client"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/unknwon/com"
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

type HeartBeat struct {
	AgentToken string
	Status     int    // 表示当前agent 状态
	Version    string // 为之后升级
	AgentIp    string
	HostName   string
	Type       string
	System     string
}

func AgentHeartBeat(c *gin.Context) {
	appG := app.Gin{C: c}
	var heartBeat HeartBeat
	var agents models.Agent

	err := c.ShouldBind(&heartBeat)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	agent := agents.GetAgentByAgentToken(heartBeat.AgentToken)

	if agent == nil {
		agent = &models.Agent{
			AgentId:       util.GenerateId(),
			AgentToken:    heartBeat.AgentToken,
			AgentIp:       heartBeat.AgentIp,
			SystemType:    heartBeat.System,
			SubnetMask:    "",
			HostName:      heartBeat.HostName,
			Status:        heartBeat.Status,
			CreateTime:    util.GetCurrentIntTime(),
			HeartbeatTime: util.GetCurrentIntTime(),
		}
		err = agent.CreateAgent()
		if err != nil {
			appG.Response(http.StatusOK, app.ErrorDatabase, nil)
			return
		}
	} else {
		agent.HeartbeatTime = util.GetCurrentIntTime()
		agent.Status = heartBeat.Status
		err = agent.UpdateAgentHeartBeat()
		if err != nil {
			appG.Response(http.StatusOK, app.ErrorDatabase, nil)
			return
		}
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func LoadAgentTask(c *gin.Context) {
	appG := app.Gin{C: c}
	agentToken := com.StrTo(c.Param("id")).String()

	if !agent_client.HasTask(agentToken) {
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}
	agentTask := agent_client.LoadAgentTask(agentToken)
	appG.Response(http.StatusOK, app.SUCCESS, agentTask)
}

func AgentTaskCallBack(c *gin.Context) {

	appG := app.Gin{C: c}

	var agentTask comm.AgentTaskBase

	err := c.ShouldBind(&agentTask)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	if agentTask.TaskType == comm.TRANSPARENT_PROXY_TASK {
		proxies := models.TransparentProxy{}
		err = proxies.UpdateTransparentProxyStatusByID(agentTask.Status, agentTask.TaskId)
		if err != nil {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
	} else if agentTask.TaskType == comm.BAIT_TASK {
		baitTasks := models.BaitTask{}
		err = baitTasks.UpdateBaitTaskStatusById(agentTask.Status, agentTask.TaskId)
		if err != nil {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
	}
	agent_client.FinishTask(agentTask.TaskId)
	appG.Response(http.StatusOK, app.SUCCESS, agentTask)
}

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

func DownloadWindowsAgent(c *gin.Context) {
	appG := app.Gin{C: c}

	var URL string

	if configs.GetSetting().App.Extranet != "http://localhost:8082" {
		URL = configs.GetSetting().App.Extranet + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + "agent/decept-agent-win.tar.gz"
	} else {
		URL = configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + "agent/decept-agent-win.tar.gz"
	}
	appG.Response(http.StatusOK, app.SUCCESS, URL)
}

func AgentPage(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.SelectPayload
	var record models.Agent
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}

	agents, count, err := record.QueryAgentPage(&payload)

	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: agents})
}

func RefreshAgentStatus() {
	var record models.Agent
	err := record.RefreshAgentStatus()
	if err != nil {
		zap.L().Error(err.Error())
		return
	}
	return
}
