package trans_proxy_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/internal/message_client"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TransparentProxyCreatePayload struct {
	ProxyPort       int32 `json:"ProxyPort" binding:"required"`       //代理端口
	ProbeID         int64 `json:"ProbeID" binding:"required"`         //探针ID
	ProtocolProxyID int64 `json:"ProtocolProxyID" binding:"required"` //协议代理ID
}

// CreateTransparentProxy 创建透明代理
// @Summary 创建透明代理
// @Description 创建透明代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param ProxyPort body TransparentProxyCreatePayload true "ProxyPort"
// @Param ProtocolProxyID body TransparentProxyCreatePayload true "ProtocolProxyID"
// @Param ProbeID body TransparentProxyCreatePayload true "ProbeID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/transparent [post]
func CreateTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	var transparentProxy models.TransparentProxy
	var taskPayload comm.ProxyTaskPayload
	var protocolProxy models.ProtocolProxy
	var payload TransparentProxyCreatePayload
	var protocol models.Protocols
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	p, _ := transparentProxy.GetTransparentProxyByProxyPort(payload.ProxyPort, payload.ProbeID)
	if p != nil {
		zap.L().Error("透明代理端口重复")
		appG.Response(http.StatusOK, app.ErrorProxyPortDup, nil)
		return
	}

	currentUser, exist := c.Get("currentUser")
	if !exist {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	transparentProxy.CreateTime = util.GetCurrentTime()
	transparentProxy.Creator = *(currentUser.(*string))
	transparentProxy.ProxyPort = payload.ProxyPort

	id, err := util.GetUniqueID()
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorUUID, nil)
		return
	}
	transparentProxy.TaskID = id
	r, err := protocolProxy.GetProtocolProxyByID(payload.ProtocolProxyID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
		return
	}
	s, err := protocol.GetProtocolByID(r.ProtocolID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
		return
	}
	if payload.ProxyPort < s.MinPort || payload.ProxyPort > s.MaxPort {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyPortRange, nil)
		return
	}
	transparentProxy.ProtocolProxyID = r.ID
	probe, err := (&(models.Probes{})).GetServerStatusByID(payload.ProbeID)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorProbeServerNotExist, nil)
		return
	}

	honeypot, err := (&(models.Honeypot{})).GetHoneypotByID(r.HoneypotID)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, err.Error())
		return
	}
	honeypotServer, err := (&(models.HoneypotServers{})).GetServerByID(honeypot.ServersID)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, err.Error())
		return
	}
	transparentProxy.ServerID = probe.ID
	{
		taskPayload.TaskType = comm.TransparentProxy
		taskPayload.OperatorType = comm.DEPLOY
		taskPayload.HoneypotServerPort = r.ProxyPort
		taskPayload.ProxyPort = payload.ProxyPort
		taskPayload.TaskID = id
		taskPayload.ProbeIP = probe.ServerIP
		taskPayload.HoneypotServerIP = honeypotServer.ServerIP
		taskPayload.AgentID = probe.AgentID

	}
	transparentProxy.AgentID = probe.AgentID
	transparentProxy.Status = comm.RUNNING
	if err := transparentProxy.CreateTransparentProxy(); err != nil {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyCreate, err.Error())
		return
	}

	jsonByte, _ := json.Marshal(taskPayload)
	err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorRedis, err.Error())
		return
	}

	timeTicker := time.NewTicker(10 * time.Microsecond)
	count := 0
	for {
		<-timeTicker.C
		count++
		proxy, err := transparentProxy.GetTransparentProxyByTaskID(transparentProxy.TaskID)
		if err == nil {
			if proxy.Status == comm.SUCCESS {
				timeTicker.Stop()
				appG.Response(http.StatusOK, app.SUCCESS, nil)
				return
			} else if proxy.Status == comm.FAILED {
				timeTicker.Stop()
				transparentProxy.DeleteTransparentProxyByID(proxy.ID)
				appG.Response(http.StatusOK, app.ErrorTransparentProxyFail, nil)
				return
			}
		}
		if count >= 1000 {
			timeTicker.Stop()
			transparentProxy.DeleteTransparentProxyByID(proxy.ID)
			appG.Response(http.StatusOK, app.ErrorTransparentProxyFail, nil)
			return
		}
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetTransparentByAgent 查找Agent激活状态的透明代理
// @Summary 查找Agent激活状态的透明代理
// @Description 查找Agent激活状态的透明代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param Payload body comm.SelectTransparentProxyPayload false "Payload"
// @Param ProtocolProxyID body comm.SelectTransparentProxyPayload true "ProtocolProxyID"
// @Param PageNumber body comm.SelectTransparentProxyPayload true "PageNumber"
// @Param PageSize body comm.SelectTransparentProxyPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.TransparentProxySelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/transparent/set [post]
func GetTransparentByAgent(c *gin.Context) {
	appG := app.Gin{C: c}
	agentId := c.Query("agentId")
	var record models.TransparentProxy
	data, err := record.GetTransparentByAgent(agentId)
	if err != nil {
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, data)
}

// GetTransparentProxy 查找透明代理
// @Summary 查找透明代理
// @Description 查找透明代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param Payload body comm.SelectTransparentProxyPayload false "Payload"
// @Param ProtocolProxyID body comm.SelectTransparentProxyPayload true "ProtocolProxyID"
// @Param PageNumber body comm.SelectTransparentProxyPayload true "PageNumber"
// @Param PageSize body comm.SelectTransparentProxyPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.TransparentProxySelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/transparent/set [post]
func GetTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.SelectTransparentProxyPayload
	var record models.TransparentProxy
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := record.GetTransparentProxy(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// DeleteTransparentProxy 删除透明代理
// @Summary 删除透明代理
// @Description 删除透明代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/transparent/:id [delete]
func DeleteTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	var transparentProxy models.TransparentProxy
	var protocolProxy models.ProtocolProxy
	var probe models.Probes
	var honeypotServer models.HoneypotServers
	var honeypot models.Honeypot
	r, err := transparentProxy.GetTransparentProxyByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyNotExist, nil)
		return
	}
	r1, err := protocolProxy.GetProtocolProxyByID(r.ProtocolProxyID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
		return
	}
	r2, err := honeypot.GetHoneypotByID(r1.HoneypotID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, nil)
		return
	}
	s, err := probe.GetServerStatusByID(r.ServerID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProbeNotExist, nil)
		return
	}
	p, err := honeypotServer.GetServerByID(r2.ServersID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, nil)
		return
	}

	var taskPayload comm.ProxyTaskPayload
	{
		taskPayload.TaskType = comm.TransparentProxy
		taskPayload.OperatorType = comm.WITHDRAW
		taskPayload.HoneypotServerPort = r1.ProxyPort
		taskPayload.ProxyPort = r.ProxyPort
		taskPayload.HoneypotServerIP = p.ServerIP
		taskPayload.TaskID = r.TaskID
		taskPayload.ProbeIP = s.ServerIP
		taskPayload.AgentID = s.AgentID
		jsonByte, _ := json.Marshal(taskPayload)
		err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
		if err != nil {
			appG.Response(http.StatusOK, app.ErrorRedis, nil)
			return
		}
	}
	if err := transparentProxy.DeleteTransparentProxyByID(id); err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// OfflineTransparentProxy 下线透明代理
// @Summary 下线透明代理
// @Description 下线透明代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/transparent/offline/:id [post]
func OfflineTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	var transparentProxy models.TransparentProxy
	var protocolProxy models.ProtocolProxy
	var probe models.Probes
	var honeypotServer models.HoneypotServers
	var honeypot models.Honeypot
	r, err := transparentProxy.GetTransparentProxyByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyNotExist, nil)
		return
	}
	if r.Status == comm.FAILED {
		appG.Response(http.StatusOK, app.ErrorTransparentOffline, nil)
		return
	}
	r1, err := protocolProxy.GetProtocolProxyByID(r.ProtocolProxyID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
		return
	}
	r2, err := honeypot.GetHoneypotByID(r1.HoneypotID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, nil)
		return
	}
	s, err := probe.GetServerStatusByID2(r.ServerID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProbeServerNotExist, nil)
		return
	}
	p, err := honeypotServer.GetServerByID(r2.ServersID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProbeServerNotExist, nil)
		return
	}

	var taskPayload comm.ProxyTaskPayload
	{
		taskPayload.TaskType = comm.TransparentProxy
		taskPayload.OperatorType = comm.WITHDRAW
		taskPayload.HoneypotServerPort = r1.ProxyPort
		taskPayload.ProxyPort = r.ProxyPort
		taskPayload.HoneypotServerIP = p.ServerIP
		taskPayload.TaskID = r.TaskID
		taskPayload.ProbeIP = s.ServerIP
		taskPayload.AgentID = s.AgentID
		jsonByte, _ := json.Marshal(taskPayload)
		err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
		if err != nil {
			appG.Response(http.StatusOK, app.ErrorRedis, nil)
			return
		}
	}
	timeTicker := time.NewTicker(10 * time.Microsecond)
	count := 0
	for {
		<-timeTicker.C
		count++
		proxy, err := transparentProxy.GetTransparentProxyByTaskID(taskPayload.TaskID)
		if err == nil {
			if proxy.Status == comm.FAILED {
				timeTicker.Stop()
				appG.Response(http.StatusOK, app.SUCCESS, nil)
				return
			}
		}
		if count >= 1000 {
			timeTicker.Stop()
			appG.Response(http.StatusOK, app.ErrorTransparentOffline, nil)
			return
		}
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// BatchOfflineTransparentProxy 批量下线透明代理
// @Summary 批量下线透明代理
// @Description 批量下线透明代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/transparent/offline/:id [post]
func BatchOfflineTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	var payload comm.BatchSelectPayload

	err := c.ShouldBind(&payload)

	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	for _, id := range payload.Ids {
		valid.Min(id, 1, "id").Message("ID必须大于0")
		if valid.HasErrors() {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
	}
	for _, id := range payload.Ids {
		code := offlineTransparent(int64(id))
		if code != app.SUCCESS {
			appG.Response(http.StatusOK, code, nil)
			return
		}
	}

	appG.Response(http.StatusOK, app.SUCCESS, nil)

}

func offlineTransparent(id int64) int {
	var transparentProxy models.TransparentProxy
	var protocolProxy models.ProtocolProxy
	var probe models.Probes
	var honeypotServer models.HoneypotServers
	var honeypot models.Honeypot
	r, err := transparentProxy.GetTransparentProxyByID(id)
	if err != nil {

		return app.ErrorTransparentProxyNotExist
	}
	if r.Status == comm.FAILED {
		return app.ErrorTransparentOffline
	}
	r1, err := protocolProxy.GetProtocolProxyByID(r.ProtocolProxyID)
	if err != nil {
		return app.ErrorProtocolProxyNotExist
	}
	r2, err := honeypot.GetHoneypotByID(r1.HoneypotID)
	if err != nil {
		return app.ErrorHoneypotNotExist
	}
	s, err := probe.GetServerStatusByID2(r.ServerID)
	if err != nil {
		return app.ErrorProbeServerNotExist
	}
	p, err := honeypotServer.GetServerByID(r2.ServersID)
	if err != nil {
		return app.ErrorProbeServerNotExist
	}

	var taskPayload comm.ProxyTaskPayload
	{
		taskPayload.TaskType = comm.TransparentProxy
		taskPayload.OperatorType = comm.WITHDRAW
		taskPayload.HoneypotServerPort = r1.ProxyPort
		taskPayload.ProxyPort = r.ProxyPort
		taskPayload.HoneypotServerIP = p.ServerIP
		taskPayload.TaskID = r.TaskID
		taskPayload.ProbeIP = s.ServerIP
		taskPayload.AgentID = s.AgentID
		jsonByte, _ := json.Marshal(taskPayload)
		err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
		if err != nil {
			return app.ErrorRedis
		}
	}
	timeTicker := time.NewTicker(10 * time.Microsecond)
	count := 0
	for {
		<-timeTicker.C
		count++
		proxy, err := transparentProxy.GetTransparentProxyByTaskID(taskPayload.TaskID)
		if err == nil {
			if proxy.Status == comm.FAILED {
				timeTicker.Stop()
				return app.SUCCESS
			}
		}
		if count >= 100 {
			timeTicker.Stop()
			return app.ErrorTransparentOffline
		}
		time.Sleep(500)
	}
}

// OnlineTransparentProxy 上线透明代理
// @Summary 上线透明代理
// @Description 上线透明代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/transparent/online/:id [post]
func OnlineTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	var transparentProxy models.TransparentProxy
	var protocolProxy models.ProtocolProxy
	var probe models.Probes
	var honeypotServer models.HoneypotServers
	var honeypot models.Honeypot
	r, err := transparentProxy.GetTransparentProxyByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyNotExist, nil)
		return
	}
	if r.Status == comm.SUCCESS {
		appG.Response(http.StatusOK, app.ErrorTransparentOnline, nil)
		return
	}
	r1, err := protocolProxy.GetProtocolProxyByID(r.ProtocolProxyID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
		return
	}
	r2, err := honeypot.GetHoneypotByID(r1.HoneypotID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, nil)
		return
	}
	s, err := probe.GetServerStatusByID2(r.ServerID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProbeServerNotExist, nil)
		return
	}
	p, err := honeypotServer.GetServerByID(r2.ServersID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProbeServerNotExist, nil)
		return
	}

	var taskPayload comm.ProxyTaskPayload
	{
		taskPayload.TaskType = comm.TransparentProxy
		taskPayload.OperatorType = comm.DEPLOY
		taskPayload.HoneypotServerPort = r1.ProxyPort
		taskPayload.ProxyPort = r.ProxyPort
		taskPayload.HoneypotServerIP = p.ServerIP
		taskPayload.TaskID = r.TaskID
		taskPayload.ProbeIP = s.ServerIP
		taskPayload.AgentID = s.AgentID
		jsonByte, _ := json.Marshal(taskPayload)
		err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
		if err != nil {
			appG.Response(http.StatusOK, app.ErrorRedis, nil)
			return
		}
	}
	timeTicker := time.NewTicker(10 * time.Microsecond)
	count := 0
	for {
		<-timeTicker.C
		count++
		proxy, err := transparentProxy.GetTransparentProxyByTaskID(taskPayload.TaskID)
		if err == nil {
			if proxy.Status == comm.SUCCESS {
				timeTicker.Stop()
				appG.Response(http.StatusOK, app.SUCCESS, nil)
				return
			}
		}
		if count >= 1000 {
			timeTicker.Stop()
			appG.Response(http.StatusOK, app.ErrorTransparentProxyOnlineFail, nil)
			return
		}
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// BatchOnlineTransparentProxy 批量上线透明代理
// @Summary 批量上线透明代理
// @Description 批量上线透明代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/transparent/online/batch [post]
func BatchOnlineTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	var payload comm.BatchSelectPayload

	err := c.ShouldBind(&payload)

	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	for _, id := range payload.Ids {
		valid.Min(id, 1, "id").Message("ID必须大于0")
		if valid.HasErrors() {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
	}
	for _, id := range payload.Ids {
		code := onlineTransparentProxyById(int64(id))
		if code != app.SUCCESS {
			appG.Response(http.StatusOK, code, nil)
			return
		}
	}

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func onlineTransparentProxyById(transparentId int64) int {
	var transparentProxy models.TransparentProxy
	var protocolProxy models.ProtocolProxy
	var probe models.Probes
	var honeypotServer models.HoneypotServers
	var honeypot models.Honeypot
	r, err := transparentProxy.GetTransparentProxyByID(transparentId)
	if err != nil {
		return app.ErrorTransparentProxyNotExist
	}
	r1, err := protocolProxy.GetProtocolProxyByID(r.ProtocolProxyID)
	if err != nil {
		return app.ErrorProtocolProxyNotExist
	}
	r2, err := honeypot.GetHoneypotByID(r1.HoneypotID)
	if err != nil {
		return app.ErrorHoneypotNotExist
	}
	s, err := probe.GetServerStatusByID2(r.ServerID)
	if err != nil {
		return app.ErrorProbeServerNotExist
	}
	p, err := honeypotServer.GetServerByID(r2.ServersID)
	if err != nil {
		return app.ErrorProbeServerNotExist
	}

	var taskPayload comm.ProxyTaskPayload
	{
		taskPayload.TaskType = comm.TransparentProxy
		taskPayload.OperatorType = comm.DEPLOY
		taskPayload.HoneypotServerPort = r1.ProxyPort
		taskPayload.ProxyPort = r.ProxyPort
		taskPayload.HoneypotServerIP = p.ServerIP
		taskPayload.TaskID = r.TaskID
		taskPayload.ProbeIP = s.ServerIP
		taskPayload.AgentID = s.AgentID
		jsonByte, _ := json.Marshal(taskPayload)
		err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
		if err != nil {
			return app.ErrorRedis
		}
	}
	timeTicker := time.NewTicker(10 * time.Microsecond)
	count := 0
	for {
		<-timeTicker.C
		count++
		proxy, err := transparentProxy.GetTransparentProxyByTaskID(taskPayload.TaskID)
		if err == nil {
			if proxy.Status == comm.SUCCESS {
				timeTicker.Stop()
				return app.SUCCESS
			}
		}
		time.Sleep(500)
		if count >= 10 {
			timeTicker.Stop()
			return app.ErrorTransparentProxyOnlineFail
		}
	}
}

// TestTransparentProxy 测试透明代理
// @Summary 测试透明代理
// @Description 测试透明代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/transparent/test/:id [get]
func TestTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	var transparentProxy models.TransparentProxy
	var probe models.Probes
	r, err := transparentProxy.GetTransparentProxyByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyNotExist, nil)
		return
	}
	s, err := probe.GetServerStatusByID2(r.ServerID)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProbeServerNotExist, nil)
		return
	}
	fmt.Println("test probe ips:" + s.ServerIP)
	if util.TcpGather(strings.Split(s.ServerIP, ","), strconv.Itoa(int(r.ProxyPort))) {
		transparentProxy.UpdateTransparentProxyStatusByTaskID(comm.SUCCESS, r.TaskID)
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	} else {
		appG.Response(http.StatusOK, app.ErrorConnectTest, nil)
		transparentProxy.UpdateTransparentProxyStatusByTaskID(comm.FAILED, r.TaskID)
		return
	}
}

func UpdateTransparentProxyStatus(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	var payload struct {
		Status int64 `json:"Status" binding:"required"` //代理状态
	}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	var transparentProxy models.TransparentProxy
	transparentProxy.UpdateTransparentProxyStatusByID(payload.Status, id)
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

//func UpdateTransparentProxyStatus() {
//    var transparentProxy models.TransparentProxy
//    var protocolProxy    models.ProtocolProxy
//    var probe            models.Probes
//    r, err := transparentProxy.GetTransparentProxies()
//    if err != nil{
//        return
//    }
//    for _, d := range *r{
//        if d.Status == comm.RUNNING{
//            continue
//        }
//        s, err := probe.GetServerStatusByID(d.ServerID)
//        if err != nil{
//            continue
//        }
//        if util.TcpGather(strings.Split(s.ServerIP, ","),   strconv.Itoa(int(d.ProxyPort))){
//            protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.SUCCESS, d.TaskID)
//            return
//        }else{
//            protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.FAILED, d.TaskID)
//            return
//        }
//    }
//}
