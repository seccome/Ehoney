package probe_token_handler

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

type ProbeTokenCreatePayload struct {
	TokenName  string `form:"TokenName" json:"TokenName" binding:"required"` //探针密签名称
	DeployPath string `form:"DeployPath" json:"DeployPath"`                  //部署路径
	ProbeID    int64  `form:"ServerID"  json:"ProbeID" binding:"required"`   //探针ID
}

// CreateProbeToken 创建探针密签
// @Summary 创建探针密签
// @Description 创建探针密签
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param TokenName body ProbeTokenCreatePayload true "TokenName"
// @Param DeployPath body ProbeTokenCreatePayload true "DeployPath"
// @Param ProbeID body ProbeTokenCreatePayload true "ProbeID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 1007 {string} json "{"code":1007,"msg":"Redis异常","data":{}}"
// @Failure 5007 {string} json "{"code":5007,"msg":"探针服务器不存在","data":{}}"
// @Failure 3005 {string} json "{"code":3005,"msg":"探针诱饵创建异常","data":{}}"
// @Failure 3015 {string} json "{"code":3015,"msg":"密签不存在","data":{}}"
// @Router /api/v1/token/probe [post]
func CreateProbeTokenNew(c *gin.Context) {
	appG := app.Gin{C: c}
	var probeToken models.ProbeToken
	var token models.Token
	var server models.Probes
	var payload ProbeTokenCreatePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	//TODO .() will panic, best to use value, ok to accept
	probeToken.Creator = *(currentUser.(*string))
	probeToken.CreateTime = util.GetCurrentTime()

	r, err := token.GetTokenByName(payload.TokenName)
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
	code, err := util.GetUniqueID()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	traceCode := code
	probeToken.TraceCode = traceCode
	id, err := util.GetUniqueID()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorUUID, err.Error())
		return
	}

	tokenFileCreateBody := util.TokenFileCreateBody{
		SourceFile: r.UploadPath,
		DestFile:   path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "probe_token", traceCode, path.Base(r.UploadPath)),
		TraceCode:  traceCode,
		TraceUrl:   strings.Join([]string{configs.GetSetting().App.TokenTraceAddress, configs.GetSetting().App.TokenTraceApiPath}, "/") + "?tracecode=" + traceCode,
	}

	if r.TokenType == "BrowserPDF" {
		tokenFileCreateBody.Content = r.TokenData
		fileName := probeToken.TokenName + ".pdf"
		tokenFileCreateBody.DestFile = path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "probe_token", traceCode, fileName)
	}

	if r.TokenType == "WPS" {
		tokenFileCreateBody.Content = r.TokenData
		fileName := r.TokenName

		if !strings.HasSuffix(fileName, ".doc") && !strings.HasSuffix(fileName, ".docs") && !strings.HasSuffix(fileName, ".xls") && !strings.HasSuffix(fileName, ".xlsx") && !strings.HasSuffix(fileName, ".ppt") && !strings.HasSuffix(fileName, ".pptx") {
			errMsg := "WPS蜜签文件名称不符合规范: " + "[" + fileName + "]"
			zap.L().Error(errMsg)
			appG.Response(http.StatusOK, app.ErrorDoFileTokenTrace, errMsg)
		}

		tokenFileCreateBody.DestFile = path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "honeypot_token", traceCode, fileName)
	}

	if err := util.CreateTokenFile(tokenFileCreateBody); err != nil {
		zap.L().Error("文件加签异常: " + err.Error())
		appG.Response(http.StatusOK, app.ErrorDoFileTokenTrace, err.Error())
		return
	}

	var taskPayload comm.TokenFileTaskPayload

	tokenUploadFileBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "probe_token", traceCode)
	tarPath := path.Join(tokenUploadFileBasePath, strings.Join([]string{r.TokenName, ".tar.gz"}, ""))
	tokenScriptPath := path.Join(util.WorkingPath(), configs.GetSetting().App.ScriptPath)
	if s.SystemType == "Linux" {
		tokenScriptPath = path.Join(tokenScriptPath, "linux", "token", "deploy.sh")
	} else {
		tokenScriptPath = path.Join(tokenScriptPath, "windows", "token", "deploy.bat")
	}
	err = util.CompressTarGz(tarPath, tokenScriptPath, r.UploadPath)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	probeToken.DeployPath = payload.DeployPath
	probeToken.LocalPath = tokenFileCreateBody.DestFile
	probeToken.TraceCode = traceCode
	{
		taskPayload.TaskID = id
		taskPayload.URL = util.Base64Encode("http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/probe_token/" + traceCode + "/" + strings.Join([]string{r.TokenName, ".tar.gz"}, ""))
		md5, err := util.GetFileMD5(tarPath)
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
			return
		}
		taskPayload.FileMD5 = md5
		taskPayload.AgentID = s.AgentID
		taskPayload.TaskType = comm.TOKEN
		taskPayload.OperatorType = comm.DEPLOY
		taskPayload.TokenType = r.TokenType
		taskPayload.ScriptName = path.Base(tokenScriptPath)
		taskPayload.CommandParameters = map[string]string{"-d": payload.DeployPath, "-s": r.FileName}
		jsonByte, _ := json.Marshal(taskPayload)
		taskPayload.Status = comm.RUNNING
		err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.ErrorRedis, err.Error())
			return
		}
	}
	probeToken.ServerID = payload.ProbeID
	probeToken.TokenName = r.TokenName
	probeToken.TokenType = r.TokenType
	probeToken.TaskID = id
	probeToken.Status = comm.RUNNING

	if err := probeToken.CreateProbeToken(); err != nil {
		//TODO if exception you should rollback
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorProbeBaitCreate, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// CreateProbeToken 创建探针密签
// @Summary 创建探针密签
// @Description 创建探针密签
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param TokenName body ProbeTokenCreatePayload true "TokenName"
// @Param DeployPath body ProbeTokenCreatePayload true "DeployPath"
// @Param ProbeID body ProbeTokenCreatePayload true "ProbeID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 1007 {string} json "{"code":1007,"msg":"Redis异常","data":{}}"
// @Failure 5007 {string} json "{"code":5007,"msg":"探针服务器不存在","data":{}}"
// @Failure 3005 {string} json "{"code":3005,"msg":"探针诱饵创建异常","data":{}}"
// @Failure 3015 {string} json "{"code":3015,"msg":"密签不存在","data":{}}"
// @Router /api/v1/token/probe [post]
func CreateProbeToken(c *gin.Context) {
	appG := app.Gin{C: c}
	var probeToken models.ProbeToken
	var token models.Token
	var server models.Probes
	var payload ProbeTokenCreatePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	//TODO .() will panic, best to use value, ok to accept
	probeToken.Creator = *(currentUser.(*string))
	probeToken.CreateTime = util.GetCurrentTime()

	r, err := token.GetTokenByName(payload.TokenName)
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
	code, err := util.GetUniqueID()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	traceCode := code
	probeToken.TraceCode = traceCode
	id, err := util.GetUniqueID()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorUUID, err.Error())
		return
	}
	if r.TokenType == "FILE" || r.TokenType == "EXE" {
		if payload.DeployPath == "" {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		var taskPayload comm.TokenFileTaskPayload
		sourceDir := r.UploadPath
		destDir := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "probe_token", traceCode, path.Base(r.UploadPath))
		traceUrl := strings.Join([]string{configs.GetSetting().App.TokenTraceAddress, configs.GetSetting().App.TokenTraceApiPath}, "/") + "?tracecode=" + traceCode
		if r.TokenType == "FILE" {
			if err := util.DoFileTokenTrace(sourceDir, destDir, traceUrl); err != nil {
				zap.L().Error("文件加签异常: " + err.Error())
				appG.Response(http.StatusOK, app.ErrorDoFileTokenTrace, err.Error())
				return
			}
		} else if r.TokenType == "EXE" {
			if err := util.DoEXEToken(sourceDir, destDir, traceUrl); err != nil {
				zap.L().Error("EXE加签异常: " + err.Error())
				appG.Response(http.StatusOK, app.ErrorDoFileTokenTrace, err.Error())
				return
			}
		}
		tokenUploadFileBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "probe_token", traceCode)
		tarPath := path.Join(tokenUploadFileBasePath, strings.Join([]string{r.TokenName, ".tar.gz"}, ""))
		tokenScriptPath := path.Join(util.WorkingPath(), configs.GetSetting().App.ScriptPath)
		if s.SystemType == "Linux" {
			tokenScriptPath = path.Join(tokenScriptPath, "linux", "token", "deploy.sh")
		} else {
			tokenScriptPath = path.Join(tokenScriptPath, "windows", "token", "deploy.bat")
		}
		err = util.CompressTarGz(tarPath, tokenScriptPath, r.UploadPath)
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
			return
		}
		probeToken.DeployPath = payload.DeployPath
		probeToken.LocalPath = destDir
		probeToken.TraceCode = traceCode
		{
			taskPayload.TaskID = id
			taskPayload.URL = util.Base64Encode("http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/probe_token/" + traceCode + "/" + strings.Join([]string{r.TokenName, ".tar.gz"}, ""))
			md5, err := util.GetFileMD5(tarPath)
			if err != nil {
				zap.L().Error(err.Error())
				appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
				return
			}
			taskPayload.FileMD5 = md5
			taskPayload.AgentID = s.AgentID
			taskPayload.TaskType = comm.TOKEN
			taskPayload.OperatorType = comm.DEPLOY
			taskPayload.TokenType = r.TokenType
			taskPayload.ScriptName = path.Base(tokenScriptPath)
			taskPayload.CommandParameters = map[string]string{"-d": payload.DeployPath, "-s": r.FileName}
			jsonByte, _ := json.Marshal(taskPayload)
			taskPayload.Status = comm.RUNNING
			err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
			if err != nil {
				zap.L().Error(err.Error())
				appG.Response(http.StatusOK, app.ErrorRedis, err.Error())
				return
			}
		}
	} else if r.TokenType == "BrowserPDF" {
		if payload.DeployPath == "" {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		var taskPayload comm.TokenFileTaskPayload
		fileName := probeToken.TokenName + ".pdf"
		destDir := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "honeypot_token", traceCode, fileName)
		traceUrl := strings.Join([]string{configs.GetSetting().App.TokenTraceAddress, configs.GetSetting().App.TokenTraceApiPath}, "/") + "?tracecode=" + traceCode
		if err := util.DoBrowserPDFToken(r.TokenData, destDir, traceUrl); err != nil {
			zap.L().Error("浏览器PDF密签创建异常: " + err.Error())
			appG.Response(http.StatusOK, app.ErrorDoFileTokenTrace, err.Error())
			return
		}
		tokenUploadFileBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "probe_token", traceCode)
		tarPath := path.Join(tokenUploadFileBasePath, strings.Join([]string{r.TokenName, ".tar.gz"}, ""))
		tokenScriptPath := path.Join(util.WorkingPath(), configs.GetSetting().App.ScriptPath)
		if s.SystemType == "Linux" {
			tokenScriptPath = path.Join(tokenScriptPath, "linux", "token", "deploy.sh")
		} else {
			tokenScriptPath = path.Join(tokenScriptPath, "windows", "token", "deploy.bat")
		}
		err = util.CompressTarGz(tarPath, tokenScriptPath, r.UploadPath)
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
			return
		}
		probeToken.DeployPath = payload.DeployPath
		probeToken.LocalPath = destDir
		probeToken.TraceCode = traceCode
		{
			taskPayload.TaskID = id
			taskPayload.URL = util.Base64Encode("http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/probe_token/" + traceCode + "/" + strings.Join([]string{r.TokenName, ".tar.gz"}, ""))
			md5, err := util.GetFileMD5(tarPath)
			if err != nil {
				zap.L().Error(err.Error())
				appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
				return
			}
			taskPayload.FileMD5 = md5
			taskPayload.AgentID = s.AgentID
			taskPayload.TaskType = comm.TOKEN
			taskPayload.OperatorType = comm.DEPLOY
			taskPayload.TokenType = r.TokenType
			taskPayload.ScriptName = path.Base(tokenScriptPath)
			taskPayload.CommandParameters = map[string]string{"-d": payload.DeployPath, "-s": r.FileName}
			jsonByte, _ := json.Marshal(taskPayload)
			taskPayload.Status = comm.RUNNING
			err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
			if err != nil {
				zap.L().Error(err.Error())
				appG.Response(http.StatusOK, app.ErrorRedis, err.Error())
				return
			}
		}
	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	probeToken.ServerID = payload.ProbeID
	probeToken.TokenName = r.TokenName
	probeToken.TokenType = r.TokenType
	probeToken.TaskID = id
	probeToken.Status = comm.RUNNING

	if err := probeToken.CreateProbeToken(); err != nil {
		//TODO if exception you should rollback
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorProbeBaitCreate, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetProbeToken 查看探针密签
// @Summary 查看探针密签
// @Description 查看探针密签
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param ServerID body comm.ServerTokenSelectPayload true "ServerID"
// @Param Payload body comm.ServerTokenSelectPayload false "Payload"
// @Param PageNumber body comm.ServerTokenSelectPayload true "PageNumber"
// @Param PageSize body comm.ServerTokenSelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.ServerTokenSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token/probe/set [post]
func GetProbeToken(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.ServerTokenSelectPayload
	var probeSigns models.ProbeToken
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	data, count, err := probeSigns.GetProbeToken(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// DeleteProbeTokenByID 删除探针密签
// @Summary 删除探针密签
// @Description 删除探针密签
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token/probe/:id [delete]
func DeleteProbeTokenByID(c *gin.Context) {
	appG := app.Gin{C: c}
	var probeSign models.ProbeToken
	var server models.Probes
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	p, err := probeSign.GetProbeTokenByID(id)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	s, err := server.GetServerStatusByID(p.ServerID)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	if p.TokenType == "FILE" || p.TokenType == "EXE" {
		var taskPayload comm.TokenFileTaskPayload
		tokenUploadFileBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "probe_token", p.TraceCode)
		tarPath := path.Join(tokenUploadFileBasePath, strings.Join([]string{"un_" + p.TokenName, ".tar.gz"}, ""))
		tokenScriptPath := path.Join(util.WorkingPath(), configs.GetSetting().App.ScriptPath)
		if s.SystemType == "Linux" {
			tokenScriptPath = path.Join(tokenScriptPath, "linux", "token", "withdraw.sh")
		} else {
			tokenScriptPath = path.Join(tokenScriptPath, "windows", "token", "withdraw.bat")
		}
		err = util.CompressTarGz(tarPath, tokenScriptPath)
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
			return
		}
		{
			taskPayload.TaskID = p.TaskID
			taskPayload.URL = util.Base64Encode("http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/probe_token/" + p.TraceCode + "/" + strings.Join([]string{"un_", p.TokenName, ".tar.gz"}, ""))
			md5, err := util.GetFileMD5(tarPath)
			if err != nil {
				zap.L().Error(err.Error())
				appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
				return
			}
			taskPayload.FileMD5 = md5
			taskPayload.AgentID = s.AgentID
			taskPayload.TaskType = comm.TOKEN
			taskPayload.OperatorType = comm.WITHDRAW
			taskPayload.TokenType = p.TokenType
			taskPayload.ScriptName = path.Base(tokenScriptPath)
			taskPayload.CommandParameters = map[string]string{"-d": path.Join(p.DeployPath, path.Base(p.LocalPath))}
			jsonByte, _ := json.Marshal(taskPayload)
			err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
			if err != nil {
				appG.Response(http.StatusOK, app.ErrorRedis, err.Error())
				return
			}
		}
	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if err := probeSign.DeleteProbeTokenByID(id); err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorDeleteBait, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// DownloadProbeTokenByID 下载探针密签
// @Summary 下载探针密签
// @Description 下载探针密签
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token/honeypot/:id [get]
func DownloadProbeTokenByID(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var token models.ProbeToken
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := token.GetProbeTokenByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	var URL = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/probe_token/" + r.TraceCode + "/" + path.Base(r.LocalPath)
	appG.Response(http.StatusOK, app.SUCCESS, URL)
}
