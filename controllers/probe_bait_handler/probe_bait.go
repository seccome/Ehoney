package probe_bait_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/internal/message_client"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type ProbeBaitCreatePayload struct {
	BaitID     int64  `form:"BaitID" binding:"required"`
	DeployPath string `form:"DeployPath" json:"DeployPath"`
	BaitData   string `form:"BaitData" json:"BaitData"`
	ProbeID    int64  `form:"ProbeID" binding:"required"`
}

// CreateProbeBait 创建探针诱饵
// @Summary 创建探针诱饵
// @Description 创建探针诱饵
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param BaitID body ProbeBaitCreatePayload true "BaitID"
// @Param DeployPath body ProbeBaitCreatePayload false "DeployPath"
// @Param BaitData body ProbeBaitCreatePayload false "BaitData"
// @Param ProbeID body ProbeBaitCreatePayload true "ProbeID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 1007 {string} json "{"code":1007,"msg":"Redis异常","data":{}}"
// @Failure 5007 {string} json "{"code":5007,"msg":"探针服务器不存在","data":{}}"
// @Failure 3005 {string} json "{"code":3005,"msg":"探针诱饵创建异常","data":{}}"
// @Failure 3011 {string} json "{"code":3011,"msg":"诱饵不存在","data":{}}"
// @Router /api/v1/bait/probe [post]
func CreateProbeBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var probeBait models.ProbeBaits
	var bait models.Baits
	var server models.Probes
	var payload ProbeBaitCreatePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error("请求参数异常")
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	//TODO .() will panic, best to use value, ok to accept
	probeBait.Creator = *(currentUser.(*string))
	probeBait.CreateTime = util.GetCurrentTime()

	r, err := bait.GetBaitByID(payload.BaitID)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorBaitNotExist, err.Error())
		return
	}
	s, err := server.GetServerStatusByID(payload.ProbeID)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorProbeNotExist, err.Error())
		return
	}
	id, err := util.GetUniqueID()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorUUID, nil)
		return
	}
	if r.BaitType == "FILE" {
		if payload.DeployPath == "" {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		var taskPayload comm.BaitFileTaskPayload
		baitUploadFileBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "bait")
		baitScriptPath := path.Join(util.WorkingPath(), configs.GetSetting().App.ScriptPath)
		if s.SystemType == "Linux" {
			baitScriptPath = path.Join(baitScriptPath, "linux", "bait", "deploy.sh")
		} else {
			baitScriptPath = path.Join(baitScriptPath, "windows", "bait", "deploy.bat")
		}
		tarPath := path.Join(baitUploadFileBasePath, strings.Join([]string{r.BaitName, ".tar.gz"}, ""))
		err = util.CompressTarGz(tarPath, baitScriptPath, r.UploadPath)
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
			return
		}
		probeBait.DeployPath = payload.DeployPath
		{
			taskPayload.TaskID = id
			taskPayload.URL = util.Base64Encode("http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/bait/" + strings.Join([]string{r.BaitName, ".tar.gz"}, ""))
			md5, err := util.GetFileMD5(tarPath)
			if err != nil {
				zap.L().Error(err.Error())
				appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
				return
			}
			taskPayload.FileMD5 = md5
			taskPayload.AgentID = s.AgentID
			taskPayload.TaskType = comm.BAIT
			taskPayload.OperatorType = comm.DEPLOY
			taskPayload.ScriptName = path.Base(baitScriptPath)
			taskPayload.BaitType = r.BaitType
			taskPayload.CommandParameters = map[string]string{"-d": payload.DeployPath, "-s": r.FileName}
			jsonByte, _ := json.Marshal(taskPayload)
			taskPayload.Status = comm.RUNNING
			err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
			if err != nil {
				appG.Response(http.StatusOK, app.ErrorRedis, nil)
				return
			}
		}
	} else if r.BaitType == "HISTORY" {
		var taskPayload comm.HistoryBaitDeployTaskPayload
		taskPayload.TaskID = id
		taskPayload.TaskType = comm.BAIT
		taskPayload.OperatorType = comm.DEPLOY
		taskPayload.BaitType = r.BaitType
		taskPayload.AgentID = s.AgentID
		taskPayload.BaitData = util.Base64Encode(r.BaitData)
		jsonByte, _ := json.Marshal(taskPayload)
		taskPayload.Status = comm.RUNNING
		err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
		if err != nil {
			appG.Response(http.StatusOK, app.ErrorRedis, nil)
			return
		}
	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	probeBait.ServerID = payload.ProbeID
	probeBait.BaitName = r.BaitName
	probeBait.BaitType = r.BaitType
	probeBait.TaskID = id
	probeBait.LocalPath = r.UploadPath

	if err := probeBait.CreateProbeBait(); err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorProbeBaitCreate, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetProbeBait 查找探针诱饵
// @Summary 查找探针诱饵
// @Description 查找探针诱饵
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param BaitType body comm.ServerBaitSelectPayload false "BaitType"
// @Param BaitName body comm.ServerBaitSelectPayload false "BaitName"
// @Param Status body comm.ServerBaitSelectPayload false "Status"
// @Param StartTimestamp body comm.ServerBaitSelectPayload false "StartTimestamp"
// @Param EndTimestamp body comm.ServerBaitSelectPayload false "EndTimestamp"
// @Param PageNumber body comm.ServerBaitSelectPayload true "PageNumber"
// @Param PageSize body comm.ServerBaitSelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.ServerBaitSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/bait/probe/set [post]
func GetProbeBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.ServerBaitSelectPayload
	var probeBaits models.ProbeBaits
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := probeBaits.GetProbeBait(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// DeleteProbeBaitByID 删除探针诱饵
// @Summary 删除探针诱饵
// @Description 删除探针诱饵
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 1007 {string} json "{"code":1007,"msg":"Redis异常","data":{}}"
// @Router /api/v1/bait/probe/:id [delete]
func DeleteProbeBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	var probeBait models.ProbeBaits
	var server models.Probes
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	p, err := probeBait.GetProbeBaitByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	s, err := server.GetServerStatusByID(p.ServerID)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if p.BaitType == "FILE" {
		var taskPayload comm.BaitFileTaskPayload
		baitUploadFileBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "bait")
		tarPath := path.Join(baitUploadFileBasePath, strings.Join([]string{"un_" + p.BaitName, ".tar.gz"}, ""))
		baitScriptPath := path.Join(util.WorkingPath(), configs.GetSetting().App.ScriptPath)
		if s.SystemType == "Linux" {
			baitScriptPath = path.Join(baitScriptPath, "linux", "bait", "withdraw.sh")
		} else {
			baitScriptPath = path.Join(baitScriptPath, "windows", "bait", "withdraw.bat")
		}
		err = util.CompressTarGz(tarPath, baitScriptPath)
		if err != nil {
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
		{
			taskPayload.TaskID = p.TaskID
			taskPayload.URL = util.Base64Encode("http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/bait/" + strings.Join([]string{"un_", p.BaitName, ".tar.gz"}, ""))
			md5, err := util.GetFileMD5(tarPath)
			if err != nil {
				appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
				return
			}
			taskPayload.FileMD5 = md5
			taskPayload.AgentID = s.AgentID
			taskPayload.TaskType = comm.BAIT
			taskPayload.OperatorType = comm.WITHDRAW
			taskPayload.BaitType = p.BaitType
			taskPayload.ScriptName = path.Base(baitScriptPath)
			taskPayload.CommandParameters = map[string]string{"-d": path.Join(p.DeployPath, path.Base(p.LocalPath))}
			jsonByte, _ := json.Marshal(taskPayload)
			taskPayload.Status = comm.RUNNING
			err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
			if err != nil {
				appG.Response(http.StatusOK, app.ErrorRedis, nil)
				return
			}
		}
		probeBait.Status = 1
	} else if p.BaitType == "HISTORY" {

	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	if err := probeBait.DeleteProbeBaitByID(id); err != nil {
		appG.Response(http.StatusOK, app.ErrorDeleteBait, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// DownloadProbeBaitByID 下载探针诱饵
// @Summary 下载探针诱饵
// @Description 下载探针诱饵
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/bait/probe/:id [get]
func DownloadProbeBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var probeBait models.ProbeBaits
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	p, err := probeBait.GetProbeBaitByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	var URL string = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/bait/" + p.BaitName + "/" + path.Base(p.LocalPath)

	appG.Response(http.StatusOK, app.SUCCESS, URL)
}
