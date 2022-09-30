package agent_bait_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/agent_client"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type AgentBaitTaskCreatePayload struct {
	AgentId    string `form:"AgentId" binding:"required"`
	BaitId     string `form:"BaitId" binding:"required"`
	BaitType   string `form:"BaitType" json:"BaitType"`
	BaitData   string `form:"BaitData" json:"BaitData"`
	DeployPath string `form:"DeployPath" json:"DeployPath"`
}

func CreateProbeBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var baits models.Bait
	var agents models.Agent
	var payload AgentBaitTaskCreatePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error("请求参数异常")
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	bait, err := baits.GetBaitById(payload.BaitId)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorBaitNotExist, err.Error())
		return
	}
	agent := agents.GetAgentByAgentToken(payload.AgentId)
	if agent == nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorProbeNotExist, err.Error())
		return
	}

	var baitTask models.BaitTask
	baitTask.BaitTaskId = util.GenerateId()
	baitTask.CreateTime = util.GetCurrentIntTime()
	baitTask.DeployPath = payload.DeployPath
	baitTask.AgentToken = agent.AgentToken
	baitTask.BaitName = bait.BaitName
	baitTask.BaitId = bait.BaitId
	baitTask.BaitType = bait.BaitType
	baitTask.LocalPath = bait.UploadPath
	baitTask.OperatorType = string(comm.DEPLOY)
	if bait.BaitType == "FILE" {
		if payload.DeployPath == "" {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		baitUploadFileBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "bait")
		baitScriptPath := path.Join(util.WorkingPath(), configs.GetSetting().App.ScriptPath)
		if agent.SystemType == "Linux" {
			baitScriptPath = path.Join(baitScriptPath, "linux", "bait", "deploy.sh")
		} else {
			baitScriptPath = path.Join(baitScriptPath, "windows", "bait", "deploy.bat")
		}
		tarPath := path.Join(baitUploadFileBasePath, strings.Join([]string{bait.BaitName, ".tar.gz"}, ""))
		err = util.CompressTarGz(tarPath, baitScriptPath, bait.UploadPath)
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
			return
		}

		baitTask.URL = util.Base64Encode("http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/bait/" + strings.Join([]string{bait.BaitName, ".tar.gz"}, ""))
		md5, err := util.GetFileMD5(tarPath)
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
			return
		}
		baitTask.FileMD5 = md5
		baitTask.ScriptName = path.Base(baitScriptPath)
		baitTask.CommandParameters = fmt.Sprintf("-d %s -s %s", payload.DeployPath, bait.FileName)
	} else if bait.BaitType == "HISTORY" {
		baitTask.BaitData = util.Base64Encode(bait.BaitData)
	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	if err := baitTask.CreateBaitTask(); err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorProbeBaitCreate, err.Error())
		return
	}

	_ = agent_client.RegisterBaitTask(baitTask, comm.DEPLOY)

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetProbeBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var queryMap map[string]interface{}
	var baitTasks models.BaitTask
	err := c.ShouldBindJSON(&queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := baitTasks.QueryBaitTaskPage(queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

func DeleteProbeBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	var probeBaits models.BaitTask
	var agents models.Agent
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	baitTask, err := probeBaits.GetBaitTaskById(id)

	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	s := agents.GetAgentByAgentToken(baitTask.AgentToken)
	if s == nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	baitTask.BaitTaskId = util.GenerateId()
	baitTask.CreateTime = util.GetCurrentIntTime()
	baitTask.OperatorType = string(comm.WITHDRAW)

	if baitTask.BaitType == "FILE" {
		baitUploadFileBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "bait")
		tarPath := path.Join(baitUploadFileBasePath, strings.Join([]string{"un_" + baitTask.BaitName, ".tar.gz"}, ""))
		baitScriptPath := path.Join(util.WorkingPath(), configs.GetSetting().App.ScriptPath)
		if s.SystemType == "Linux" {
			baitScriptPath = path.Join(baitScriptPath, "linux", "bait", "withdraw.sh")
		} else {
			baitScriptPath = path.Join(baitScriptPath, "windows", "bait", "withdraw.bat")
		}
		err = util.CompressTarGz(tarPath, baitScriptPath)

		{
			baitTask.URL = util.Base64Encode("http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/bait/" + strings.Join([]string{"un_", baitTask.BaitName, ".tar.gz"}, ""))
			md5, err := util.GetFileMD5(tarPath)
			if err != nil {
				appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
				return
			}
			baitTask.FileMD5 = md5
			baitTask.AgentToken = s.AgentToken
			baitTask.ScriptName = path.Base(baitScriptPath)
			baitTask.CommandParameters = fmt.Sprintf("-d %s ", path.Join(baitTask.DeployPath, path.Base(baitTask.LocalPath)))
			baitTask.Status = comm.RUNNING
		}
	} else if baitTask.BaitType == "HISTORY" {

	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	if err := baitTask.CreateBaitTask(); err != nil {
		appG.Response(http.StatusOK, app.ErrorDeleteBait, nil)
		return
	}

	_ = agent_client.RegisterBaitTask(*baitTask, comm.WITHDRAW)

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func DownloadProbeBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	var probeBait models.BaitTask
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	p, err := probeBait.GetBaitTaskById(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	URL := "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/bait/" + p.BaitName + "/" + path.Base(p.LocalPath)

	appG.Response(http.StatusOK, app.SUCCESS, URL)
}
