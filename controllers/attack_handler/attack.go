package attack_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/internal/cluster"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"math"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func GetFalcoAttackList(c *gin.Context) {
	appG := app.Gin{C: c}
	var queryMap map[string]interface{}
	var attack models.FalcoAttackEvent
	err := c.ShouldBind(&queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := attack.GetFalcoEvent(queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

func GetTokenTraceLog(c *gin.Context) {
	appG := app.Gin{C: c}
	var queryMap map[string]interface{}
	var tokenTraceLogs models.TokenTraceLog
	err := c.ShouldBind(&queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := tokenTraceLogs.GetTokenTraceLog(queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

func GetAttackSource(c *gin.Context) {
	appG := app.Gin{C: c}
	var queryMap map[string]interface{}
	var falcoEvents models.FalcoAttackEvent
	var attackEvents models.AttackEvent
	var traceSources []comm.TraceSourceVo
	err := c.ShouldBind(&queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	if queryMap["StartTime"] != nil && queryMap["StartTime"].(string) != "" {
		timeP, err := time.Parse("2006-01-02 15:04:05", queryMap["StartTime"].(string))
		if err == nil {
			queryMap["StartTime"] = strconv.FormatInt(timeP.Unix(), 10)
		}
	}
	if queryMap["EndTime"] != nil && queryMap["EndTime"].(string) != "" {
		timeP, err := time.Parse("2006-01-02 15:04:05", queryMap["EndTime"].(string))
		if err == nil {
			queryMap["EndTime"] = strconv.FormatInt(timeP.Unix(), 10)
		}
	}

	falcoList, falcoTotaol, err := falcoEvents.GetFalcoEvent(queryMap)

	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}

	attackList, attckTotaol, err := attackEvents.GetAttackEvent(queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}

	if queryMap["Type"].(string) == "Falco" {
		if falcoTotaol > 0 {
			for _, i := range *falcoList {
				falcoTraceSourceVo := buildTraceSourceVoByFalcoEvent(i)
				traceSources = append(traceSources, falcoTraceSourceVo)
			}
		}
	} else if queryMap["Type"].(string) == "Attack" {
		if attckTotaol > 0 {
			for _, i := range *attackList {
				attackTraceSourceVo := buildTraceSourceVoByAttackEvent(i)
				traceSources = append(traceSources, attackTraceSourceVo)
			}
		}
	} else {
		if falcoTotaol > 0 {
			for _, i := range *falcoList {
				falcoTraceSourceVo := buildTraceSourceVoByFalcoEvent(i)
				traceSources = append(traceSources, falcoTraceSourceVo)
			}
		}
		if attckTotaol > 0 {
			for _, i := range *attackList {
				attackTraceSourceVo := buildTraceSourceVoByAttackEvent(i)
				traceSources = append(traceSources, attackTraceSourceVo)
			}
		}
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: falcoTotaol + attckTotaol, List: traceSources})
}

func buildTraceSourceVoByFalcoEvent(falcEvent models.FalcoAttackEvent) comm.TraceSourceVo {
	bytes, _ := json.Marshal(falcEvent.OutputFields)

	traceSourceVo := comm.TraceSourceVo{
		Id:           falcEvent.FalcoAttackEventId,
		Type:         "Falco事件",
		AttackIp:     "Unknown",
		ProtocolType: falcEvent.OutputFields.Repository,
		HoneypotName: falcEvent.HoneypotName,
		Log:          falcEvent.Output,
		Time:         falcEvent.Time,
		EventTime:    falcEvent.CreateTime,
		Detail:       string(bytes),
	}
	return traceSourceVo
}

func buildTraceSourceVoByAttackEvent(attackEvent models.AttackEvent) comm.TraceSourceVo {
	traceSourceVo := comm.TraceSourceVo{
		Id:           attackEvent.ProtocolEventId,
		Type:         "攻击事件",
		AttackIp:     attackEvent.AttackIp,
		ProtocolType: attackEvent.ProtocolType,
		HoneypotName: attackEvent.HoneypotIp,
		Log:          attackEvent.AttackDetail,
		Time:         "",
		EventTime:    attackEvent.CreateTime,
		Detail:       attackEvent.AttackDetail,
	}
	return traceSourceVo
}

func GetFalcoAttackDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	var data models.FalcoAttackEvent
	d, err := data.GetFalcoEventByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}

	appG.Response(http.StatusOK, app.SUCCESS, d.OutputFields)
}

func GetAttackList(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload map[string]interface{}
	var attack models.AttackEvent
	err := c.ShouldBind(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	ret, count, err := attack.GetAttackEvent(payload)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: ret})
}

func CreateTransparentEventEvent(c *gin.Context) {
	appG := app.Gin{C: c}
	var attackEvent models.TransparentEvent

	err := c.ShouldBind(&attackEvent)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	if util.IsLocalIP(attackEvent.AttackIp) {
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}

	attackEvent.AttackLocation = util.FindLocationByIp(attackEvent.AttackIp)

	if err := attackEvent.CreateEvent(); err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	waning := `Agent:  ` + attackEvent.AgentToken + `\n\n > 攻击IP:  ` + attackEvent.AttackIp + `\n\n > 攻击端口:  ` + strconv.Itoa(int(attackEvent.AttackPort)) + `\n\n > 代理IP:  ` + attackEvent.ProxyIp + `\n\n > 代理端口:  ` + strconv.Itoa(int(attackEvent.ProxyPort)) + `\n\n > 蜜罐IP:  ` + attackEvent.DestIp + `\n\n > 蜜罐端口:  ` + strconv.Itoa(int(attackEvent.DestPort)) + `\n\n > 创建时间:  ` + util.Sec2TimeStr(attackEvent.CreateTime, "") + ``
	go util.SendDingMsg("欺骗防御告警", "透明代理告警", waning)
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func CreateProtocolAttackEvent(c *gin.Context) {
	appG := app.Gin{C: c}
	var attackEvent models.ProtocolEvent
	err := c.ShouldBind(&attackEvent)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	if attackEvent.AttackIp == "127.0.0.1" ||  attackEvent.AttackIp == configs.GetSetting().Server.AppHost {
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}

	if attackEvent.ProtocolProxyId != "" {
		attackEvent.ProtocolProxyId = strings.ReplaceAll(attackEvent.ProtocolProxyId, ":", "")
	}

	address := net.ParseIP(attackEvent.AttackIp)
	if address == nil {
		attackEvent.AttackIp = "Unknown"
	}

	if err := attackEvent.CreateEvent(); err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	waning := `攻击类型: ` + attackEvent.ProtocolType + `协议代理` + `\n\n > 攻击IP:  ` + attackEvent.AttackIp + `\n\n > 攻击端口:  ` + strconv.Itoa(int(attackEvent.AttackPort)) + `\n\n > 代理IP:  ` + attackEvent.ProxyIp + `\n\n > 代理端口:  ` + strconv.Itoa(int(attackEvent.ProxyPort)) + `\n\n > 蜜罐IP:  ` + attackEvent.DestIp + `\n\n > 蜜罐端口:  ` + strconv.Itoa(int(attackEvent.DestPort)) + `\n\n > 协议类型:  ` + attackEvent.ProtocolType + `\n\n > 创建时间:  ` + attackEvent.EventTime + ``
	util.SendDingMsg("欺骗防御告警", "协议代理告警", waning)
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func CreateFalcoAttackEvent(c *gin.Context) {
	appG := app.Gin{C: c}
	var event models.FalcoAttackEvent
	err := c.ShouldBind(&event)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	shouldAlarm := false
	var honeypot models.Honeypot
	data, _ := honeypot.GetHoneypotByPodName(event.OutputFields.PodName)
	if data.ImageAddress == "ehoney/smb:v1" && (event.OutputFields.Cmdline == "lpqd -FS --no-process-group" || event.OutputFields.Cmdline == "nmbd -D" || event.OutputFields.Cmdline == "smbd -FS --no-process-group") {
		zap.L().Info(fmt.Sprintf("ignore smb falco event: %v", event))
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}

	if data.ImageAddress == "ehoney/ftp:v1" && (event.OutputFields.Cmdline == "vsftpd /etc/vsftpd/vsftpd.conf" || event.OutputFields.Cmdline == "nmbd -D" || event.OutputFields.Cmdline == "smbd -FS --no-process-group") {
		zap.L().Info(fmt.Sprintf("ignore smb falco event: %v", event))
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}
	if event.Rule == "Create files below container any dir" {
		event.FileFlag = true
		code, err := util.GetUniqueID()
		if err != nil {
			zap.L().Error(err.Error())
			appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
			return
		}
		dest := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "falco", code)
		cluster.CopyFromPod(event.OutputFields.PodName, data.HoneypotName, event.OutputFields.FilePath, dest)
		event.DownloadPath = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/falco/" + code + "/" + path.Base(event.OutputFields.FilePath)
		shouldAlarm = true
	} else if event.Rule == "Any listen" {
		shouldAlarm = true
	} else if event.Rule == "Any accept" {
		shouldAlarm = true
	} else if event.Rule == "Any outbound connect" {
		shouldAlarm = true
	} else if event.Priority == "Error" || event.Priority == "ERROR" {
		shouldAlarm = true
	} else if event.Priority == "CRITICAL" || event.Priority == "Critical" {
		shouldAlarm = true
	}

	if shouldAlarm && data != nil {
		waning := `攻击类型: ` + event.Rule + "(falco)" + `\n\n > 蜜罐:  ` + data.HoneypotIp + `\n\n > 等级:  ` + event.Priority + "(" + data.HoneypotName + ")" + `\n\n > 攻击IP:  ` + event.OutputFields.Connection + `\n\n > 路径:  ` + event.OutputFields.FilePath + `\n\n > 用户:  ` + event.OutputFields.UserName + `\n\n > CMD:  ` + event.OutputFields.Cmdline + `\n\n > 进程/父进程:  ` + event.OutputFields.ProcessName + "/" + event.OutputFields.ProcessPName + `\n\n > 蜜罐名称:  ` + event.OutputFields.ContainerName + ``
		util.SendDingMsg("欺骗防御告警", "Falco日志告警", waning)
	}
	event.HoneypotName = event.OutputFields.PodName
	event.Time = time.Unix(event.OutputFields.EventTime/int64(math.Pow10(9)), 0).In(time.FixedZone("CST", 8*3600)).String()
	if err := event.CreateFalcoEvent(); err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetAttackIPStatistics(c *gin.Context) {
	appG := app.Gin{C: c}
	var attackEvent models.TransparentEvent

	data, err := attackEvent.GetAttackStatisticsByIP()
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, data)
}

func GetProbeIPStatistics(c *gin.Context) {
	appG := app.Gin{C: c}
	var attackEvent models.TransparentEvent

	data, err := attackEvent.GetAttackStatisticsByProbeIP()
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, data)
}

func GetAttackProtocolStatistics(c *gin.Context) {
	appG := app.Gin{C: c}
	var attackEvent models.ProtocolEvent

	data, err := attackEvent.GetAttackStatisticsByProtocol()
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, data)
}

func GetAttackLocationStatistics(c *gin.Context) {
	appG := app.Gin{C: c}
	var attackEvent models.TransparentEvent

	data, err := attackEvent.GetAttackStatisticsByIP()
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	for index, d := range *data {
		d, _ := util.GetLocationByIP(d.Data)
		if d.City == "-" || d.Country_long == "-" {
			(*data)[index].Data = "LAN"
		} else {
			(*data)[index].Data = d.City + "-" + d.Country_long
		}
	}
	appG.Response(http.StatusOK, app.SUCCESS, data)
}

func CreateCountEvent(c *gin.Context) {
	appG := app.Gin{C: c}

	var event models.CounterEvent

	counterInfo := c.Query("sid")
	ret := util.Base64Decode(counterInfo)

	ip := gjson.Get(ret, "ip").String()
	protocolType := gjson.Get(ret, "type").String()
	token := gjson.Get(ret, "token_builder").String()
	info := gjson.Get(ret, "info").String()

	data, _ := event.GetCounterEvent(info, protocolType, ip, token)
	if data == nil {
		event.IP = ip
		event.Type = protocolType
		event.Info = info
		event.Token = token
		event.CreateCountEvent()
	}

	appG.Response(http.StatusOK, app.SUCCESS, ret)
}
