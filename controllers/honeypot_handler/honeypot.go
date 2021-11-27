package honeypot_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/internal/cluster"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
)

// CreateHoneypot 创建蜜罐
// @Summary 创建蜜罐
// @Description 创建蜜罐
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param HoneypotName body models.Honeypot true "HoneypotName"
// @Param ImageAddress body models.Honeypot true "ImageAddress"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 5001 {string} json "{"code":5001,"msg":"蜜罐服务器不存在、请检测蜜罐服务状态","data":{}}"
// @Failure 5002 {string} json "{"code":5002,"msg":"K8S创建镜像失败","data":{}}"
// @Failure 5006 {string} json "{"code":5006,"msg":"蜜罐名称存在","data":{}}"
// @Failure 7001 {string} json "{"code":7001,"msg":"镜像端口或是服务未配置、请先进行配置","data":{}}"
// @Router /api/v1/honeypot [post]
func CreateHoneypot(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypot models.Honeypot
	var server models.HoneypotServers
	var image models.Images
	err := c.ShouldBindJSON(&honeypot)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	honeypot.CreateTime = util.GetCurrentTime()
	honeypot.Creator = *(currentUser.(*string))
	r, err := server.GetFirstHoneypotServer()
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, nil)
		return
	}
	honeypot.ServersID = r.ID
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
		cluster.DeleteDeployment(honeypot.HoneypotName)
		appG.Response(http.StatusOK, app.ErrorHoneypotCreate, err)
		return
	}
	honeypot.PodName = pod.Name
	honeypot.ServerPort = s.ImagePort
	honeypot.ServerType = s.ImageType
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

// GetHoneypots 查找蜜罐接口
// @Summary 查找蜜罐接口
// @Description 查找蜜罐接口
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param Payload body comm.HoneypotSelectPayload false "Payload"
// @Param ProtocolType body comm.HoneypotSelectPayload false "ProtocolType"
// @Param PageNumber body comm.HoneypotSelectPayload true "PageNumber"
// @Param PageSize body comm.HoneypotSelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.HoneypotSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/honeypot/set [post]
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

// DeleteHoneypot 删除蜜罐接口
// @Summary 删除蜜罐接口
// @Description 删除蜜罐接口
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 5005 {string} json "{"code":5005,"msg":"K8S删除蜜罐错误","data":{}}"
// @Router /api/v1/honeypot/:id [delete]
func DeleteHoneypot(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypot models.Honeypot
	var protocolProxy models.ProtocolProxy
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
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
		data.CreateHoneypot()
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorHoneypotDelete, err.Error())
		return
	}
	if !flag {
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}
	err = cluster.DeleteDeployment(data.HoneypotName)
	if err != nil {
		data.CreateHoneypot()
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorHoneypotDelete, err.Error())
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetHoneypotDetail 获得蜜罐详细信息
// @Summary 获得蜜罐详细信息
// @Description 获得蜜罐详细信息
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 5005 {string} json "{"code":5005,"msg":"K8S删除蜜罐错误","data":{}}"
// @Router /api/v1/honeypot/:id [get]
func GetHoneypotDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypot models.Honeypot
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
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

// GetProtocolProxyHoneypots 协议代理蜜罐
// @Summary 查看已经进行协议转发的蜜罐
// @Description 查看已经进行协议转发的蜜罐
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"请求参数错误","data":["honeypot-10.42.0.41"]}"
// @Router /api/v1/honeypot/protocol [get]
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
	honeypots, _ := record.GetPodNameList()
	for _, honeypot := range honeypots {
		r := cluster.GetPodDetailInfo(honeypot)
		if r.Status == "Running" {
			record.RefreshServerStatusByPodName(honeypot, comm.SUCCESS, r.PodIP)
		} else if r.Status == "Failed" {
			record.RefreshServerStatusByPodName(honeypot, comm.FAILED, "")
		} else {
			continue
		}
	}
	return
}
