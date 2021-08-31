package heartbeat_handler

import (
	"decept-defense/pkg/app"
	"github.com/gin-gonic/gin"
	"net/http"
)


// Heartbeat 心跳
// @Summary 心跳接口
// @Description 心跳接口
// @Tags 心跳
// @Produce application/json
// @Accept application/json
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Router /api/public/heartbeat [get]
func Heartbeat(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.OKHeartBeat, nil)
}