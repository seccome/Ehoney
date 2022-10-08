package protocol_proxy_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/proxy"
	"decept-defense/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
)

type ProtocolProxyCreatePayload struct {
	HoneypotId   string `json:"HoneypotId" binding:"required"`
	ProtocolType string `json:"ProtocolType"`
	ProxyPort    int32  `json:"ProxyPort" binding:"required"`
}

func CreateProtocolProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	var protocolProxy models.ProtocolProxy
	var Protocols models.Protocols
	var honeypots models.Honeypot
	var payload ProtocolProxyCreatePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	p, _ := protocolProxy.GetProtocolProxyByProxyPort(payload.ProxyPort)
	if p != nil {
		zap.L().Error("协议代理端口重复")
		appG.Response(http.StatusOK, app.ErrorProxyPortDup, nil)
		return
	}
	honeypot, err := honeypots.GetHoneypotByID(payload.HoneypotId)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, nil)
		return
	}
	protocol, err := Protocols.GetProtocolByType(honeypot.ProtocolType)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolNotExist, nil)
		return
	}
	protocolProxy.ProtocolProxyId = util.GenerateId()
	protocolProxy.CreateTime = util.GetCurrentIntTime()
	protocolProxy.ProtocolProxyName = fmt.Sprintf("%s-%d", honeypot.ProtocolType, payload.ProxyPort)
	protocolProxy.HoneypotId = honeypot.HoneypotId
	protocolProxy.HoneypotIp = honeypot.HoneypotIp
	protocolProxy.ProtocolId = protocol.ProtocolId
	protocolProxy.ProtocolType = protocol.ProtocolType
	protocolProxy.HoneypotPort = honeypot.ServerPort
	protocolProxy.ProtocolPath = generateDeployPath(protocol.LocalPath, protocol.FileName)
	protocolProxy.ProxyPort = payload.ProxyPort
	protocolProxy.Status = comm.RUNNING
	if processId, err := proxy.StartProxyProtocol(protocolProxy); err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyCreate, nil)
		return
	} else {
		protocolProxy.ProcessPid = processId
	}

	if err := protocolProxy.CreateProtocolProxy(); err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyCreate, nil)
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetProtocolProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.SelectPayload
	var record models.ProtocolProxy
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := record.GetProtocolProxy(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

func DeleteProtocolProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	var protocols models.ProtocolProxy
	protocol, err := protocols.GetProtocolProxyByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
		return
	}
	if protocol.ProcessPid > 100 {
		if err := util.KillProcess(protocol.ProcessPid); err != nil {
			appG.Response(http.StatusOK, app.ErrorProtocolProxyFail, nil)
			return
		}
	}

	if err := protocols.DeleteProtocolProxyByID(protocol.ProtocolProxyId); err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func OfflineProtocolProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	var protocols models.ProtocolProxy
	protocol, err := protocols.GetProtocolProxyByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
		return
	}
	if err := util.KillProcess(protocol.ProcessPid); err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyFail, nil)
		return
	}

	protocol.Status = comm.FAILED

	if err := protocol.UpdateStatus(); err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func OnlineProtocolProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	var protocols models.ProtocolProxy
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	protocolProxy, err := protocols.GetProtocolProxyByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	if processId, err := proxy.StartProxyProtocol(*protocolProxy); err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolOnline, nil)
		return
	} else {
		protocolProxy.ProcessPid = processId
	}
	protocolProxy.Status = comm.SUCCESS

	if err := protocolProxy.UpdateStatus(); err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func TestProtocolProxy(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	var protocols models.ProtocolProxy
	protocolProxy, err := protocols.GetProtocolProxyByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
		return
	}

	if proxy.TestLocalPortConnection(protocolProxy.ProxyPort) {
		protocolProxy.Status = comm.SUCCESS
		_ = protocols.UpdateStatus()
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}
	protocolProxy.Status = comm.FAILED
	_ = protocols.UpdateStatus()
	appG.Response(http.StatusOK, app.ErrorConnectTest, nil)
}

func generateDeployPath(localPath, fileName string) string {
	projectPath, _ := os.Getwd()
	return path.Join(projectPath, localPath, fileName)
}

func UpdateDeployedProtocolProxyStatus() {
	var protocols models.ProtocolProxy
	protocolProxies, err := protocols.GetAllProtocolProxies()
	if err != nil {
		return
	}
	for _, protocolProxy := range *protocolProxies {
		if proxy.TestLocalPortConnection(protocolProxy.ProxyPort) {
			protocolProxy.Status = comm.SUCCESS
		} else {
			protocolProxy.Status = comm.FAILED
		}
		_ = protocolProxy.UpdateStatus()
	}
}
