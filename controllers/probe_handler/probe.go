package probe_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// GetProbes 探针列表
// @Summary 探针列表
// @Description 探针列表
// @Tags 探针管理
// @Produce application/json
// @Accept application/json
// @Param Payload body comm.SelectPayload false "Payload"
// @Param PageNumber body comm.SelectPayload true "PageNumber"
// @Param PageSize body comm.SelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.ProbeSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/probe/set [post]
func GetProbes(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.SelectPayload
	var record models.Probes
	err := c.ShouldBindJSON(&payload)
	if err != nil{
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	data, count, err := record.GetProbe(&payload)
	if err != nil{
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

func RefreshProbeStatus(){
	var record models.Probes
	var honeypotServer models.HoneypotServers
	err := record.RefreshServerStatus()
	if err != nil{
		zap.L().Error(err.Error())
		return
	}
	err = honeypotServer.RefreshServerStatus()
	if err != nil{
		zap.L().Error(err.Error())
		return
	}
	return
}
