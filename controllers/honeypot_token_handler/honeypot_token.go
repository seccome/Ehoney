package honeypot_token_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/internal/cluster"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type HoneypotTokenCreatePayload struct {
	TokenName  string `form:"TokenName"  json:"TokenName" binding:"required"`  //蜜罐密签名称
	DeployPath string `form:"DeployPath" json:"DeployPath"`                    //部署路径
	HoneypotID int64  `form:"HoneypotID" json:"HoneypotID" binding:"required"` //蜜罐ID
}

/*var TokenInjectCmdPattenMap = map[string][]string{
"exe": {"-t", "-b"},
"file": {"-t", "-b"}}*/

// CreateHoneypotToken 创建蜜罐密签
// @Summary 创建蜜罐密签
// @Description 创建蜜罐密签
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param TokenName body HoneypotTokenCreatePayload true "TokenName"
// @Param DeployPath body HoneypotTokenCreatePayload false "DeployPath"
// @Param HoneypotID body HoneypotTokenCreatePayload true "HoneypotID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 5001 {string} json "{"code":5001,"msg":"蜜罐服务器不存在、请检测蜜罐服务状态","data":{}}"
// @Failure 3005 {string} json "{"code":3005,"msg":"蜜罐诱饵创建异常","data":{}}"
// @Failure 3015 {string} json "{"code":3015,"msg":"密签不存在","data":{}}"
// @Failure 3012 {string} json "{"code":3012,"msg":"K8S拷贝异常","data":{}}"
// @Router /api/v1/token/honeypot [post]
func CreateHoneypotTokenNew(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypotToken models.HoneypotToken
	var token models.Token
	var honeypot models.Honeypot
	var payload HoneypotTokenCreatePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist {
		zap.L().Error("当前用户获取错误")
		appG.Response(http.StatusOK, app.INTERNAlERROR, "当前用户获取错误")
		return
	}
	value, ok := currentUser.(*string)
	if !ok {
		zap.L().Error("当前用户解析错误")
		appG.Response(http.StatusOK, app.INTERNAlERROR, "当前用户解析错误")
		return
	}
	if payload.DeployPath == "" {
		zap.L().Error("请求参数错误")
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	honeypotToken.Creator = *(value)
	honeypotToken.CreateTime = util.GetCurrentTime()

	code, err := util.GetUniqueID()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	traceCode := code
	honeypotToken.TraceCode = traceCode

	r, err := token.GetTokenByName(payload.TokenName)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorTokenNotExist, err.Error())
		return
	}
	s, err := honeypot.GetHoneypotByID(payload.HoneypotID)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, err.Error())
		return
	}

	tokenFileCreateBody := util.TokenFileCreateBody{
		TokenType:  r.TokenType,
		SourceFile: r.UploadPath,
		DestFile:   path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "honeypot_token", traceCode, path.Base(r.UploadPath)),
		TraceCode:  traceCode,
		TraceUrl:   strings.Join([]string{configs.GetSetting().App.TokenTraceAddress, configs.GetSetting().App.TokenTraceApiPath}, "/") + "?tracecode=" + traceCode,
	}

	if r.TokenType == "BrowserPDF" {
		tokenFileCreateBody.Content = r.TokenData
		fileName := r.TokenName + ".pdf"
		tokenFileCreateBody.DestFile = path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "honeypot_token", traceCode, fileName)
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

	honeypotToken.TokenName = r.TokenName
	honeypotToken.TokenType = r.TokenType
	honeypotToken.HoneypotID = payload.HoneypotID
	honeypotToken.TraceCode = traceCode
	honeypotToken.LocalPath = tokenFileCreateBody.DestFile
	honeypotToken.DeployPath = payload.DeployPath

	if err := util.CreateTokenFile(tokenFileCreateBody); err != nil {
		zap.L().Error("文件加签异常: " + err.Error())
		appG.Response(http.StatusOK, app.ErrorDoFileTokenTrace, err.Error())
		return
	}

	if err = cluster.CopyToPod(s.PodName, s.HoneypotName, tokenFileCreateBody.DestFile, payload.DeployPath); err != nil {
		zap.L().Error("k3s拷贝文件异常: " + err.Error())
		appG.Response(http.StatusOK, app.ErrorHoneypotK8SCP, err)
		return
	}

	honeypotToken.Status = comm.SUCCESS
	// 数据处理 save honey pod token record

	if err := honeypotToken.CreateHoneypotToken(); err != nil {
		cluster.RemoveFromPod(s.PodName, s.HoneypotName, path.Join(payload.DeployPath, path.Base(honeypotToken.LocalPath)))
		zap.L().Error("创建密签记录、数据库异常")
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// CreateHoneypotToken 创建蜜罐密签
// @Summary 创建蜜罐密签
// @Description 创建蜜罐密签
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param TokenName body HoneypotTokenCreatePayload true "TokenName"
// @Param DeployPath body HoneypotTokenCreatePayload false "DeployPath"
// @Param HoneypotID body HoneypotTokenCreatePayload true "HoneypotID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 5001 {string} json "{"code":5001,"msg":"蜜罐服务器不存在、请检测蜜罐服务状态","data":{}}"
// @Failure 3005 {string} json "{"code":3005,"msg":"蜜罐诱饵创建异常","data":{}}"
// @Failure 3015 {string} json "{"code":3015,"msg":"密签不存在","data":{}}"
// @Failure 3012 {string} json "{"code":3012,"msg":"K8S拷贝异常","data":{}}"
// @Router /api/v1/token/honeypot [post]
func CreateHoneypotToken(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypotToken models.HoneypotToken
	var token models.Token
	var honeypot models.Honeypot
	var payload HoneypotTokenCreatePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist {
		zap.L().Error("当前用户获取错误")
		appG.Response(http.StatusOK, app.INTERNAlERROR, "当前用户获取错误")
		return
	}
	value, ok := currentUser.(*string)
	if !ok {
		zap.L().Error("当前用户解析错误")
		appG.Response(http.StatusOK, app.INTERNAlERROR, "当前用户解析错误")
		return
	}
	honeypotToken.Creator = *(value)
	honeypotToken.CreateTime = util.GetCurrentTime()

	code, err := util.GetUniqueID()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	traceCode := code
	honeypotToken.TraceCode = traceCode

	r, err := token.GetTokenByName(payload.TokenName)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorTokenNotExist, err.Error())
		return
	}
	s, err := honeypot.GetHoneypotByID(payload.HoneypotID)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, err.Error())
		return
	}
	if r.TokenType == "FILE" || r.TokenType == "EXE" {
		if payload.DeployPath == "" {
			zap.L().Error("请求参数错误")
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		sourceDir := r.UploadPath
		destDir := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "honeypot_token", traceCode, path.Base(r.UploadPath))
		traceUrl := strings.Join([]string{configs.GetSetting().App.TokenTraceAddress, configs.GetSetting().App.TokenTraceApiPath}, "/") + "?tracecode=" + traceCode
		if r.TokenType == "FILE" {
			if err := util.DoFileTokenTrace(sourceDir, destDir, traceUrl); err != nil {
				zap.L().Error("文件加签异常: " + err.Error())
				appG.Response(http.StatusOK, app.ErrorDoFileTokenTrace, err.Error())
				return
			}
		} else if r.TokenType == "EXE" {
			if err := util.DoEXEToken(sourceDir, destDir, traceUrl); err != nil {
				zap.L().Error("文件加签异常: " + err.Error())
				appG.Response(http.StatusOK, app.ErrorDoFileTokenTrace, err.Error())
				return
			}
		}

		honeypotToken.HoneypotID = payload.HoneypotID
		honeypotToken.TokenName = r.TokenName
		honeypotToken.TokenType = r.TokenType
		honeypotToken.TraceCode = traceCode
		tokenPath := destDir
		honeypotToken.LocalPath = destDir
		honeypotToken.DeployPath = payload.DeployPath

		if err = cluster.CopyToPod(s.PodName, s.HoneypotName, tokenPath, payload.DeployPath); err != nil {
			zap.L().Error("k3s拷贝文件异常: " + err.Error())
			appG.Response(http.StatusOK, app.ErrorHoneypotK8SCP, err)
			return
		}
		honeypotToken.Status = comm.SUCCESS
	} else if r.TokenType == "BrowserPDF" {
		honeypotToken.TokenName = r.TokenName
		honeypotToken.TokenType = r.TokenType
		honeypotToken.HoneypotID = payload.HoneypotID
		honeypotToken.TraceCode = traceCode
		fileName := honeypotToken.TokenName + ".pdf"
		destDir := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "honeypot_token", traceCode, fileName)
		traceUrl := strings.Join([]string{configs.GetSetting().App.TokenTraceAddress, configs.GetSetting().App.TokenTraceApiPath}, "/") + "?tracecode=" + traceCode
		if err := util.DoBrowserPDFToken(r.TokenData, destDir, traceUrl); err != nil {
			zap.L().Error("浏览器PDF密签创建异常: " + err.Error())
			appG.Response(http.StatusOK, app.ErrorDoBrowserPDFTokenTrace, err.Error())
			return
		}
		honeypotToken.LocalPath = destDir
		honeypotToken.DeployPath = payload.DeployPath
		tokenPath := destDir
		if err = cluster.CopyToPod(s.PodName, s.HoneypotName, tokenPath, payload.DeployPath); err != nil {
			zap.L().Error("k3s拷贝文件异常: " + err.Error())
			appG.Response(http.StatusOK, app.ErrorHoneypotK8SCP, err)
			return
		}
		honeypotToken.Status = comm.SUCCESS
	} else {
		zap.L().Error("不支持的密签类型")
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if err := honeypotToken.CreateHoneypotToken(); err != nil {
		if r.TokenType == "FILE" || r.TokenType == "EXE" {
			cluster.RemoveFromPod(s.PodName, s.HoneypotName, path.Join(payload.DeployPath, path.Base(honeypotToken.LocalPath)))
		}
		zap.L().Error("创建密签记录、数据库异常")
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetHoneypotToken 查找蜜罐密签接口
// @Summary 用户查看自定义蜜罐密签接口
// @Description 用户查看自定义蜜罐密签接口
// @Tags 蜜罐管理
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
// @Router /api/v1/token/honeypot/set [post]
func GetHoneypotToken(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.ServerTokenSelectPayload
	var honeypotToken models.HoneypotToken
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := honeypotToken.GetHoneypotToken(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// DeleteHoneypotTokenByID 删除蜜罐密签
// @Summary 删除蜜罐密签
// @Description 删除蜜罐密签
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token/honeypot/:id [delete]
func DeleteHoneypotTokenByID(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypotToken models.HoneypotToken
	var honeypot models.Honeypot
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := honeypotToken.GetHoneypotTokenByID(id)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	s, err := honeypot.GetHoneypotByID(r.HoneypotID)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	if r.TokenType == "FILE" || r.TokenType == "EXE" {
		err = cluster.RemoveFromPod(s.PodName, s.HoneypotName, path.Join(r.DeployPath, path.Base(r.LocalPath)))
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.ErrorHoneypotBaitWithdraw, err.Error())
			return
		}
		os.RemoveAll(filepath.Dir(filepath.Dir(r.LocalPath)))
	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if err := honeypotToken.DeleteHoneypotTokenByID(id); err != nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotBaitDelete, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// DownloadHoneypotTokenByID 下载蜜罐密签接口
// @Summary 下载蜜罐密签接口
// @Description 下载蜜罐密签接口
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token/honeypot/:id [get]
func DownloadHoneypotTokenByID(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var token models.HoneypotToken
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := token.GetHoneypotTokenByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	var URL string
	if configs.GetSetting().App.Extranet != "" {
		URL = "http:" + "//" + configs.GetSetting().App.Extranet + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/honeypot_token/" + r.TraceCode + "/" + path.Base(r.LocalPath)

	} else {
		URL = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/honeypot_token/" + r.TraceCode + "/" + path.Base(r.LocalPath)
	}

	appG.Response(http.StatusOK, app.SUCCESS, URL)
}
