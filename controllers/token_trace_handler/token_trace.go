package token_trace_handler

import (
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"regexp"
)

type TraceHostConfigPayload struct {
	TraceHost string `gorm:"not null;size:256" form:"TraceHost" json:"TraceHost" binding:"required"` //密签跟踪设置
}

func UpdateTraceHostConfig(c *gin.Context) {
	appG := app.Gin{C: c}

	var payload TraceHostConfigPayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	viper.Set("app.tokentraceaddress", payload.TraceHost)
	viper.WriteConfig()
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// TraceMsgReceive 接收token 的上报信息
func TraceMsgReceive(c *gin.Context) {
	appG := app.Gin{C: c}
	var baits models.Bait
	traceCodeA := c.Query("tracecode")
	if len(traceCodeA) <= 0 {
		appG.Response(http.StatusNotAcceptable, app.InvalidParams, nil)
		return
	}
	bait, err := baits.GetBaitById(traceCodeA)
	if err != nil {
		zap.L().Error("can not find bait for trace code  : " + string(traceCodeA))
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}
	var tokenTraceLog models.TokenTraceLog
	tokenTraceLog.TokenTraceLogId = util.GenerateId()
	tokenTraceLog.TraceCode = traceCodeA
	tokenTraceLog.TokenType = bait.BaitType
	tokenTraceLog.BaitName = bait.BaitName
	tokenTraceLog.OpenIP = GetIP(c.Request.RemoteAddr)
	tokenTraceLog.OpenTime = util.GetCurrentIntTime()
	tokenTraceLog.UserAgent = c.Request.Header.Get("User-Agent")
	tokenTraceLog.Location = util.FindLocationByIp(tokenTraceLog.OpenIP)
	err = tokenTraceLog.CreateTokenTraceLog()
	if err != nil {
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, "success")
}

func GetIP(srcIP string) string {
	resultIP := ""
	ipRegexp := regexp.MustCompile(`^?([^:]*)`)
	params := ipRegexp.FindStringSubmatch(srcIP)
	if len(params) > 0 {
		resultIP = params[0]
	}
	return resultIP
}

func GetTraceHostConfig(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.SUCCESS, configs.GetSetting().Server.AppHost)
}
