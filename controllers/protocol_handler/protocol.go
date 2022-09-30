package protocol_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
)

type ProtocolUpdatePayload struct {
	MinPort int32 `json:"MinPort" form:"MinPort" binding:"required"`
	MaxPort int32 `json:"MaxPort" form:"MaxPort" binding:"required"`
}

func CreateProtocol(c *gin.Context) {

}

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

func DeleteProtocol(c *gin.Context) {
	appG := app.Gin{C: c}
	var protocol models.Protocols
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
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

func UpdateProtocolPortRange(c *gin.Context) {
	appG := app.Gin{C: c}
	var protocol models.Protocols
	var payload ProtocolUpdatePayload
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	err := c.ShouldBind(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	if payload.MinPort > payload.MaxPort || payload.MinPort <= 0 || payload.MinPort > 65535 || payload.MaxPort <= 0 || payload.MaxPort > 65535 {
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
