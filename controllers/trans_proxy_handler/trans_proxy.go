package trans_proxy_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/agent_client"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"errors"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type TransparentProxyCreatePayload struct {
	ProxyPort       int32  `json:"ProxyPort" binding:"required"`       //代理端口
	AgentId         string `json:"AgentId" binding:"required"`         //探针ID
	ProtocolProxyId string `json:"ProtocolProxyId" binding:"required"` //协议代理ID
}

func CreateTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	var transparentProxy models.TransparentProxy
	var protocolProxies models.ProtocolProxy
	var payload TransparentProxyCreatePayload
	var protocols models.Protocols
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	p, _ := transparentProxy.GetTransparentProxyByProxyPort(payload.ProxyPort, payload.AgentId)
	if p != nil {
		zap.L().Error("透明代理端口重复")
		appG.Response(http.StatusOK, app.ErrorProxyPortDup, nil)
		return
	}

	protocolProxy, err := protocolProxies.GetProtocolProxyByID(payload.ProtocolProxyId)

	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
		return
	}
	s, err := protocols.GetProtocolByID(protocolProxy.ProtocolId)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolNotExist, nil)
		return
	}
	if payload.ProxyPort < s.MinPort || payload.ProxyPort > s.MaxPort {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyPortRange, nil)
		return
	}
	agent := (&(models.Agent{})).GetAgentByAgentId(payload.AgentId)

	if agent == nil {
		appG.Response(http.StatusOK, app.ErrorProbeServerNotExist, nil)
		return
	}
	transparentProxy.ProxyPort = payload.ProxyPort
	transparentProxy.DestPort = protocolProxy.ProxyPort
	transparentProxy.DestIp = configs.GetSetting().Server.AppHost
	transparentProxy.AgentId = agent.AgentId
	transparentProxy.AgentToken = agent.AgentToken
	transparentProxy.AgentIp = agent.AgentIp
	transparentProxy.AgentHost = agent.HostName
	transparentProxy.ProtocolProxyId = protocolProxy.ProtocolProxyId
	transparentProxy.ProtocolType = protocolProxy.ProtocolType
	transparentProxy.Status = comm.RUNNING

	if err := transparentProxy.CreateTransparentProxy(); err != nil {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyCreate, err.Error())
		return
	}

	_ = agent_client.RegisterTransparentProxy(transparentProxy, comm.DEPLOY)

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

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

func GetTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload map[string]interface{}
	var record models.TransparentProxy
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := record.GetTransparentProxy(payload)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

func DeleteTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	var transparentProxies models.TransparentProxy
	transparentProxy, err := transparentProxies.GetTransparentProxyByID(id)
	if err != nil || transparentProxy == nil {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyNotExist, nil)
		return
	}
	if err := transparentProxies.DeleteTransparentProxyByID(transparentProxy.TransparentProxyId); err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func OfflineTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	//valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).String()

	var transparentProxies models.TransparentProxy
	transparentProxy, err := transparentProxies.GetTransparentProxyByID(id)
	if err != nil || transparentProxy == nil {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyNotExist, nil)
		return
	}
	err = offlineTransparent(transparentProxy.TransparentProxyId)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorTransparentOffline, nil)
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

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
		err = offlineTransparent(id)
		if err != nil {
			appG.Response(http.StatusOK, app.ErrorTransparentOffline, nil)
			return
		}
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func offlineTransparent(id string) error {
	var transparentProxies models.TransparentProxy
	transparentProxy, err := transparentProxies.GetTransparentProxyByID(id)
	if err != nil {
		return err
	}
	if transparentProxy.Status != comm.FAILED {
		err = transparentProxies.UpdateTransparentProxyStatusByID(comm.FAILED, id)
		if err != nil {
			return err
		}
	}
	_ = agent_client.RegisterTransparentProxy(*transparentProxy, comm.WITHDRAW)
	return nil
}

func onlineTransparentProxyById(transparentId string) error {
	var transparentProxies models.TransparentProxy
	transparentProxy, err := transparentProxies.GetTransparentProxyByID(transparentId)
	if err != nil {
		return err
	}
	if transparentProxy.Status != comm.FAILED {
		err = transparentProxies.UpdateTransparentProxyStatusByID(comm.SUCCESS, transparentId)
		if err != nil {
			return err
		}
	}

	transparentProxy.Status = comm.SUCCESS

	agent_client.RegisterTransparentProxy(*transparentProxy, comm.DEPLOY)
	return nil
}

func OnlineTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	err := onlineTransparentProxyById(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProbeServerNotExist, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

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
		err = onlineTransparentProxyById(id)
		if err != nil {
			appG.Response(http.StatusOK, app.ErrorTransparentOnline, nil)
			return
		}
	}

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func TestTransparentProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	var transparentProxy models.TransparentProxy
	var probe models.Agent
	r, err := transparentProxy.GetTransparentProxyByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorTransparentProxyNotExist, nil)
		return
	}
	s := probe.GetAgentByAgentId(r.AgentId)
	if s == nil {
		appG.Response(http.StatusOK, app.ErrorProbeServerNotExist, nil)
		return
	}
	if util.TcpGather(strings.Split(s.AgentIp, ","), strconv.Itoa(int(r.ProxyPort))) {
		_ = transparentProxy.UpdateTransparentProxyStatusByID(comm.SUCCESS, r.TransparentProxyId)
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	} else {
		appG.Response(http.StatusOK, app.ErrorConnectTest, nil)
		_ = transparentProxy.UpdateTransparentProxyStatusByID(comm.FAILED, r.TransparentProxyId)
		return
	}
}

func testTransparentProxyStatus(transparentProxy models.TransparentProxy) error {
	var probe models.Agent

	s := probe.GetAgentByAgentId(transparentProxy.AgentId)
	if s == nil {
		return errors.New("agent not exist")
	}
	if util.TcpGather(strings.Split(s.AgentIp, ","), strconv.Itoa(int(transparentProxy.ProxyPort))) {
		_ = transparentProxy.UpdateTransparentProxyStatusByID(comm.SUCCESS, transparentProxy.TransparentProxyId)
		return nil
	} else {
		_ = transparentProxy.UpdateTransparentProxyStatusByID(comm.FAILED, transparentProxy.TransparentProxyId)
		return errors.New("update transparent proxy status error ")
	}
}

func UpdateTransparentProxyStatus(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	var payload struct {
		Status comm.TaskStatus `json:"Status" binding:"required"` //代理状态
	}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	var transparentProxy models.TransparentProxy
	_ = transparentProxy.UpdateTransparentProxyStatusByID(payload.Status, id)
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func UpdateTransparentProxiesStatus() {
	var transparentProxies models.TransparentProxy
	transparentProxyArray, err := transparentProxies.GetTransparentProxies()
	if err != nil {
		return
	}
	for _, transparentProxy := range *transparentProxyArray {
		testTransparentProxyStatus(transparentProxy)
	}
}
