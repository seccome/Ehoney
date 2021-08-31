package protocol_handler

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

type ProtocolUpdatePayload struct {
	MinPort int32 `json:"MinPort" form:"MinPort" binding:"required"`
	MaxPort int32 `json:"MaxPort" form:"MaxPort" binding:"required"`
}

// CreateProtocol 创建协议服务
// @Summary 创建协议服务
// @Description 创建协议接口
// @Tags 影子代理
// @Produce application/json
// @Accept multipart/form-data
// @Param ProtocolType body models.Protocols true "ProtocolType"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 5001 {string} json "{"code":5001,"msg":"蜜罐服务器不存在、请检测蜜罐服务状态","data":{}}"
// @Failure 6001 {string} json "{"code":6001,"msg":"协议创建失败","data":{}}"
// @Failure 6004 {string} json "{"code":6004,"msg":"协议名称重复","data":{}}"
// @Router /api/v1/protocol [post]
func CreateProtocol(c *gin.Context) {
	appG := app.Gin{C: c}
	var protocol models.Protocols
	var taskPayload comm.FileTaskPayload
	var server models.HoneypotServers
	file, err := c.FormFile("file")
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	err = c.ShouldBind(&protocol)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	if protocol.MinPort > protocol.MaxPort || protocol.MinPort <= 0 || protocol.MinPort > 65535 || protocol.MaxPort <= 0 || protocol.MaxPort > 65535{
		zap.L().Error("端口范围异常")
		appG.Response(http.StatusOK, app.ErrorProtocolPortRange, nil)
		return
	}
	_, err = protocol.GetProtocolByType(protocol.ProtocolType)
	if err == nil {
		appG.Response(http.StatusOK, app.ErrorProtocolDup, nil)
		return
	}
	protocolUploadFileBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "protocol")
	protocolScriptBasePath := path.Join(util.WorkingPath(), configs.GetSetting().App.ScriptPath, "protocol", "deploy.sh")

	currentUser, exist := c.Get("currentUser")
	if !exist {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	protocol.CreateTime = util.GetCurrentTime()
	protocol.Creator = *(currentUser.(*string))
	protocol.FileName = file.Filename
	savePath := path.Join(protocolUploadFileBasePath, protocol.ProtocolType)
	if err := util.CreateDir(savePath); err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	protocol.LocalPath = savePath
	if err := c.SaveUploadedFile(file, path.Join(savePath, protocol.FileName)); err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	tarPath := path.Join(protocolUploadFileBasePath, strings.Join([]string{protocol.ProtocolType, ".tar.gz"}, ""))
	err = util.CompressTarGz(tarPath, protocolScriptBasePath, path.Join(savePath, protocol.FileName))
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	protocol.DeployPath = configs.GetSetting().App.ProtocolDeployPath

	id, err := util.GetUniqueID()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	protocol.TaskID = id
	{
		taskPayload.TaskID = id
		taskPayload.URL = util.Base64Encode("http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/protocol/" + strings.Join([]string{protocol.ProtocolType, ".tar.gz"}, ""))
		md5, err := util.GetFileMD5(tarPath)
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
			return
		}
		taskPayload.FileMD5 = md5
		honeypot, err := server.GetFirstHoneypotServer()
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, err.Error())
			return
		}
		taskPayload.TaskType = comm.PROTOCOL
		taskPayload.OperatorType = comm.DEPLOY
		taskPayload.AgentID = honeypot.AgentID
		taskPayload.ScriptName = path.Base(protocolScriptBasePath)
		taskPayload.CommandParameters = map[string]string{"-d": protocol.DeployPath, "-s": protocol.FileName}
		jsonByte, _ := json.Marshal(taskPayload)
		err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.ErrorRedis, err.Error())
			return
		}
	}
	protocol.Status = comm.RUNNING
	err = protocol.CreateProtocol()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorProtocolCreate, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetProtocol 查找协议
// @Summary 查找协议
// @Description 查找协议
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param Payload body comm.SelectPayload false "Payload"
// @Param PageNumber body comm.SelectPayload true "PageNumber"
// @Param PageSize body comm.SelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.ProtocolSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 6002 {string} json "{"code":6002,"msg":"协议获取失败","data":{}}"
// @Router /api/v1/protocol/set [post]
func GetProtocol(c *gin.Context) {
	appG := app.Gin{C: c}
	var record models.Protocols
	var payload comm.SelectPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error("请求参数异常")
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := record.GetProtocol(&payload)
	if err != nil {
		zap.L().Error("协议服务列表请求异常")
		appG.Response(http.StatusOK, app.ErrorProtocolGet, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// GetProtocolType 查找协议类型接口
// @Summary 用户查看自定义协议类型接口
// @Description 用户查看自定义协议类型接口
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"OK","data":[""]}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/protocol/type [get]
func GetProtocolType(c *gin.Context) {
	appG := app.Gin{C: c}
	var protocol models.Protocols
	data, err := protocol.GetProtocolTypeList()
	if err != nil {
		zap.L().Error("获取协议服务类型异常")
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, data)
}

// DeleteProtocol 删除协议服务
// @Summary 删除协议服务
// @Description 删除协议服务
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 6003 {string} json "{"code":6003,"msg":"协议删除失败","data":{}}"
// @Router /api/v1/protocol/:id [delete]
func DeleteProtocol(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var protocol models.Protocols
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		zap.L().Error("请求参数异常")
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	if err := protocol.DeleteProtocolByID(id); err != nil {
		zap.L().Error("删除协议异常")
		appG.Response(http.StatusOK, app.ErrorProtocolDel, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}


func CreateSSHKey(c *gin.Context){
	appG := app.Gin{C: c}
	type SSHPayload struct {
		SSHKey  string  `json:"ssh_key"`
		AgentID string  `json:"agentid"`
	}
	var payload SSHPayload
	var setting models.Setting
	err := c.ShouldBind(payload)
	if err != nil{
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorProtocolDel, err.Error())
		return
	}
	setting.ConfigName = "SSHKey"
	setting.ConfigValue = payload.SSHKey
	setting.CreateSetting()
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// UpdateProtocolPortRange 更新协议端口范围
// @Summary 更新协议端口范围
// @Description 更新协议端口范围
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 6003 {string} json "{"code":6003,"msg":"协议删除失败","data":{}}"
// @Router /api/v1/protocol/port/:id [put]
func UpdateProtocolPortRange(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var protocol models.Protocols
	var payload ProtocolUpdatePayload
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		zap.L().Error("请求参数异常")
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	err := c.ShouldBind(&payload)
	if err != nil{
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	if payload.MinPort > payload.MaxPort || payload.MinPort <= 0 || payload.MinPort > 65535 || payload.MaxPort <= 0 || payload.MaxPort > 65535{
		zap.L().Error("端口范围异常")
		appG.Response(http.StatusOK, app.ErrorProtocolPortRange, nil)
		return
	}
	if err := protocol.UpdateProtocolPortRange(payload.MinPort, payload.MaxPort, id); err != nil {
		zap.L().Error("更新协议")
		appG.Response(http.StatusOK, app.ErrorProtocolUpdate, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}
