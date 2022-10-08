package honeypot_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/internal/cluster"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/proxy"
	"decept-defense/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
)

func CreateHoneypot(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypot models.Honeypot
	var image models.Images
	err := c.ShouldBindJSON(&honeypot)

	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	s, err := image.GetImageByAddress(honeypot.ImageAddress)

	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if err = honeypot.GetHoneypotByName(honeypot.HoneypotName); err == nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotNameExist, nil)
		return
	}
	flag, err := cluster.DeploymentIsExist(honeypot.HoneypotName)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, err)
		return
	}
	if flag {
		appG.Response(http.StatusOK, app.ErrorHoneypotPodExist, nil)
		return
	}
	pod, err := cluster.CreateDeployment(honeypot.HoneypotName, s.ImageAddress, s.ImagePort)
	if err != nil {
		_ = cluster.DeleteDeployment(honeypot.HoneypotName)
		appG.Response(http.StatusOK, app.ErrorHoneypotCreate, err)
		return
	}

	honeypot.PodName = pod.Name
	honeypot.ServerPort = s.ImagePort
	honeypot.ImageId = s.ImageId
	honeypot.ProtocolType = s.ProtocolType
	honeypot.Status = comm.RUNNING

	err = honeypot.CreateHoneypot()
	if err != nil {
		//exception remove created honeypot
		err = cluster.DeleteDeployment(honeypot.HoneypotName)
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetHoneypots(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.HoneypotSelectPayload
	var record models.Honeypot
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := record.GetHoneypot(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

func DeleteHoneypot(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypot models.Honeypot
	var protocolProxy models.ProtocolProxy
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	_, err := protocolProxy.GetProtocolProxyByHoneypotID(id)
	if err == nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotProtocolProxyExist, nil)
		return
	}

	data, err := honeypot.GetHoneypotByID(id)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	if err = honeypot.DeleteHoneypotByID(id); err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	flag, err := cluster.DeploymentIsExist(data.HoneypotName)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorDeploymentExist, err.Error())
		return
	}
	if !flag {
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}
	err = cluster.DeleteDeployment(data.HoneypotName)
	if err != nil {
		_ = data.CreateHoneypot()
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorHoneypotDelete, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetHoneypotDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypot models.Honeypot
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, err := honeypot.GetHoneypotByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	result := cluster.GetPodDetailInfo(data.PodName)
	appG.Response(http.StatusOK, app.SUCCESS, result)
}

func GetProtocolProxyHoneypots(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypot models.Honeypot
	data, err := honeypot.GetProtocolProxyHoneypot()
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, data)
}

func RefreshHoneypotStatus() {
	var record models.Honeypot
	honeypots, _ := record.GetHoneypots()

	zap.L().Info(fmt.Sprintf("load honeypots size [%d]", len(*honeypots)))

	for _, honeypot := range *honeypots {

		r := cluster.GetPodDetailInfo(honeypot.PodName)

		zap.L().Info(fmt.Sprintf("%v", r))

		zap.L().Info(fmt.Sprintf("load honeypot {%s} ip [%s] -> [%s]", honeypot.PodName, honeypot.HoneypotIp, r.PodIP))

		if r.PodIP == honeypot.HoneypotIp && (r.Status == "Running" && honeypot.Status == comm.SUCCESS) {
			continue
		}

		if r.Status == "Running" {
			_ = record.UpdatePodInfoByPodName(honeypot.PodName, comm.SUCCESS, r.PodIP)
		} else if r.Status == "Creating" {
			_ = record.UpdatePodInfoByPodName(honeypot.PodName, comm.RUNNING, "")
			continue
		} else {
			_ = record.UpdatePodInfoByPodName(honeypot.PodName, comm.FAILED, "")
			continue
		}
		refreshedHoneypot, err := record.GetHoneypotByID(honeypot.HoneypotId)
		if err != nil {
			return
		}
		if honeypot.HoneypotIp != refreshedHoneypot.HoneypotIp {
			zap.L().Info(fmt.Sprintf("honeypot ip changed [%s]->[%s] start restart protocol proxies", honeypot.HoneypotIp, refreshedHoneypot.HoneypotIp))
			updateProtocolProxiesByHoneypotUpdate(*refreshedHoneypot)
		}
	}
	return
}

func updateProtocolProxiesByHoneypotUpdate(honeypot models.Honeypot) {
	var protocolProxies models.ProtocolProxy
	honeypotProtocolProxies, err := protocolProxies.QueryProtocolProxyByHoneypot(honeypot)
	if err != nil {
		zap.L().Error(fmt.Sprintf("query protocol proxies on honeypot [%s] err", honeypot.HoneypotName))
		return
	}

	for _, honeypotProtocolProxy := range *honeypotProtocolProxies {

		if honeypotProtocolProxy.HoneypotIp == honeypot.HoneypotIp {
			continue
		}

		honeypotProtocolProxy.HoneypotIp = honeypot.HoneypotIp
		err = honeypotProtocolProxy.UpdateHoneypot()
		if err != nil {
			zap.L().Error(fmt.Sprintf("update protocol proxy honeypot [%s] ip err", honeypot.HoneypotName))
			continue
		}
		if err = util.KillProcess(honeypotProtocolProxy.ProcessPid); err != nil {
			continue
		}

		if processId, err := proxy.StartProxyProtocol(honeypotProtocolProxy); err != nil {
			zap.L().Error(fmt.Sprintf("start protocol proxy [%s] err", honeypotProtocolProxy.ProtocolProxyName))
			continue
		} else {
			honeypotProtocolProxy.ProcessPid = processId
		}
		err = honeypotProtocolProxy.UpdateStatus()
		if err != nil {
			zap.L().Error(fmt.Sprintf("update protocol proxy [%s] status err", honeypotProtocolProxy.ProtocolProxyName))
			continue
		}
	}

}
