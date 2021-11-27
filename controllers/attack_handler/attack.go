package attack_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/internal/cluster"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"math"
	"net/http"
	"path"
	"strconv"
	"time"
)

//func CreateAttackEvent(c *gin.Context) {
//	appG := app.Gin{C: c}
//
//	var attackEvent models.AttackEvent
//	var honeypot    models.Honeypot
//	err :=  c.ShouldBind(&attackEvent)
//	if err != nil {
//		appG.Response(http.StatusOK, app.InvalidParams, nil)
//		return
//	}
//	attackEvent.CreateTime = util.GetCurrentTime()
//	attackEvent.AttackAddress = "未知"
//	r, err :=  honeypot.GetHoneypotByAddress(attackEvent.AttackAddress)
//	if err != nil{
//		appG.Response(http.StatusOK, app.InvalidParams, nil)
//		return
//	}
//
//	attackEvent.ProtocolType = r.ServerType
//	if err := attackEvent.CreateAttackEvent(); err != nil{
//		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
//		return
//	}
//	appG.Response(http.StatusOK, app.SUCCESS, nil)
//}

// GetFalcoAttackList falco攻击
// @Summary falco攻击
// @Description falco攻击
// @Tags 威胁感知
// @Produce application/json
// @Param StartTime body comm.FalcoEventSelectPayload false "StartTime"
// @Param EndTime body comm.FalcoEventSelectPayload false "EndTime"
// @Param Payload body comm.FalcoEventSelectPayload false "Payload"
// @Param PageNumber body comm.FalcoEventSelectPayload true "PageNumber"
// @Param PageSize body comm.FalcoEventSelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.FalcoSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/attack/falco [post]
func GetFalcoAttackList(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.FalcoEventSelectPayload
	var attack models.FalcoAttackEvent
	err := c.ShouldBind(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := attack.GetFalcoEvent(payload)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// GetTokenTraceLog 密签跟踪日志
// @Summary 密签跟踪日志
// @Description 密签跟踪日志
// @Tags 威胁感知
// @Produce application/json
// @Param Payload body comm.TokenTraceSelectPayload false "Payload"
// @Param StartTime body comm.TokenTraceSelectPayload false "StartTime"
// @Param EndTime body comm.TokenTraceSelectPayload false "EndTime"
// @Param AttackIP body comm.TokenTraceSelectPayload false "AttackIP"
// @Param ServerType body comm.TokenTraceSelectPayload false "ServerType"
// @Param PageNumber body comm.TokenTraceSelectPayload true "PageNumber"
// @Param PageSize body comm.TokenTraceSelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.TokenTraceSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/attack/token [post]
func GetTokenTraceLog(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.TokenTraceSelectPayload
	var attack models.TokenTraceLog
	err := c.ShouldBind(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := attack.GetTokenTraceLog(payload)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// GetAttackSource 溯源
// @Summary 溯源
// @Description 溯源
// @Tags 威胁感知
// @Produce application/json
// @Param Payload body comm.AttackTraceSelectPayload false "Payload"
// @Param StartTime body comm.AttackTraceSelectPayload false "StartTime"
// @Param EndTime body comm.AttackTraceSelectPayload false "EndTime"
// @Param AttackIP body comm.AttackTraceSelectPayload false "AttackIP"
// @Param HoneypotIP body comm.AttackTraceSelectPayload false "HoneypotIP"
// @Param Type body comm.AttackTraceSelectPayload false "Type"
// @Param ProtocolType body comm.AttackTraceSelectPayload false "ProtocolType"
// @Param PageNumber body comm.AttackTraceSelectPayload true "PageNumber"
// @Param PageSize body comm.AttackTraceSelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.TraceSourceResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/attack/trace [post]
func GetAttackSource(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.AttackTraceSelectPayload
	var falcoEvent models.FalcoAttackEvent
	var attackEvent models.AttackEvent
	var ret []comm.TraceSourceResultPayload
	err := c.ShouldBind(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	d1, err := falcoEvent.GetFalcoEventForTraceSource(payload)
	d2, err := attackEvent.GetAttackEventForSource(payload)

	if payload.Type == "Falco事件" {
		if d1 != nil {
			for _, i := range *d1 {
				i.Type = "Falco事件"
				ret = append(ret, i)
			}
		}
	} else if payload.Type == "攻击事件" {
		if d2 != nil {
			for _, i := range *d2 {
				i.Type = "攻击事件"
				ret = append(ret, i)
			}
		}
	} else {
		if d1 != nil {
			for _, i := range *d1 {
				i.Type = "Falco事件"
				ret = append(ret, i)
			}
		}
		if d2 != nil {
			for _, i := range *d2 {
				i.Type = "攻击事件"
				ret = append(ret, i)
			}
		}
	}

	var count int64 = int64(len(ret))
	var start int = payload.PageSize * (payload.PageNumber - 1)
	var end int = payload.PageSize*(payload.PageNumber-1) + payload.PageSize
	if payload.PageSize*(payload.PageNumber-1) > len(ret) {
		start = len(ret)
	}
	if payload.PageSize*(payload.PageNumber-1)+payload.PageSize > len(ret) {
		end = len(ret)
	}
	ret = ret[start:end]

	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: ret})
}

// GetFalcoAttackDetail falco攻击详情
// @Summary falco攻击详情
// @Description falco攻击详情
// @Tags 威胁感知
// @Produce application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.OutputFields
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/attack/falco/:id [get]
func GetFalcoAttackDetail(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
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

// GetAttackList 流量攻击
// @Summary 流量攻击
// @Description 流量攻击
// @Tags 威胁感知
// @Produce application/json
// @Param AttackIP body comm.AttackEventSelectPayload false "AttackIP"
// @Param JumpIP body comm.AttackEventSelectPayload false "JumpIP"
// @Param ProbeIP body comm.AttackEventSelectPayload false "ProbeIP"
// @Param HoneypotIP body comm.AttackEventSelectPayload false "HoneypotIP"
// @Param ProtocolType body comm.AttackEventSelectPayload false "ProtocolType"
// @Param Payload body comm.AttackEventSelectPayload false "Payload"
// @Param PageNumber body comm.AttackEventSelectPayload true "PageNumber"
// @Param PageSize body comm.AttackEventSelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.AttackSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/v1/attack [post]

func GetAttackList(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.AttackEventSelectPayload
	var attack models.AttackEvent
	err := c.ShouldBind(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	ret, err := attack.GetAttackEvent(payload)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}

	var count = int64(len(*ret))

	var start int = payload.PageSize * (payload.PageNumber - 1)
	var end int = payload.PageSize*(payload.PageNumber-1) + payload.PageSize
	if payload.PageSize*(payload.PageNumber-1) > len(*ret) {
		start = len(*ret)
	}
	if payload.PageSize*(payload.PageNumber-1)+payload.PageSize > len(*ret) {
		end = len(*ret)
	}
	*ret = (*ret)[start:end]

	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: ret})
}

// CreateProtocolAttackEvent 创建协议代理攻击事件
// @Summary 创建协议代理攻击事件
// @Description 创建协议代理攻击事件
// @Tags 非前端接口
// @Produce application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/public/attack/protocol [post]
func CreateProtocolAttackEvent(c *gin.Context) {
	appG := app.Gin{C: c}
	var attackEvent models.ProtocolEvent
	err := c.ShouldBind(&attackEvent)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	if util.IsLocalIP(attackEvent.AttackIP) {
		appG.Response(http.StatusOK, app.SUCCESS, nil)
		return
	}
	if err := attackEvent.CreateEvent(); err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	waning := `攻击类型: ` + attackEvent.AttackType + `\n\n > AgentID:  ` + attackEvent.AgentID + `\n\n > 攻击IP:  ` + attackEvent.AttackIP + `\n\n > 攻击端口:  ` + strconv.Itoa(int(attackEvent.AttackPort)) + `\n\n > 代理IP:  ` + attackEvent.ProxyIP + `\n\n > 代理端口:  ` + strconv.Itoa(int(attackEvent.ProxyPort)) + `\n\n > 蜜罐IP:  ` + attackEvent.DestIP + `\n\n > 蜜罐端口:  ` + strconv.Itoa(int(attackEvent.DestPort)) + `\n\n > 协议类型:  ` + attackEvent.ProtocolType + `\n\n > 创建时间:  ` + attackEvent.EventTime + ``
	util.SendDingMsg("欺骗防御告警", "协议代理告警", waning)
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// CreateFalcoAttackEvent 创建falco异常事件
// @Summary 创建falco异常事件
// @Description 创建falco异常事件
// @Tags 非前端接口
// @Produce application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Router /api/public/attack/falco [post]
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
		waning := `攻击类型: ` + event.Rule + "(falco)" + `\n\n > 蜜罐:  ` + data.HoneypotIP + `\n\n > 等级:  ` + event.Priority + "(" + data.HoneypotName + ")" + `\n\n > 攻击IP:  ` + event.OutputFields.Connection + `\n\n > 路径:  ` + event.OutputFields.FilePath + `\n\n > 用户:  ` + event.OutputFields.UserName + `\n\n > CMD:  ` + event.OutputFields.Cmdline + `\n\n > 进程/父进程:  ` + event.OutputFields.ProcessName + "/" + event.OutputFields.ProcessPName + `\n\n > 蜜罐名称:  ` + event.OutputFields.ContainerName + ``
		util.SendDingMsg("欺骗防御告警", "Falco日志告警", waning)
	}

	event.Time = time.Unix(event.OutputFields.EventTime/int64(math.Pow10(9)), 0).In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")
	if err := event.CreateFalcoEvent(); err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorDatabase, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetAttackIPStatistics 攻击IP统计
// @Summary 攻击IP统计
// @Description 攻击IP统计
// @Tags 大屏展示
// @Produce application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.AttackStatistics
// @Router /api/v1/attack/display/attackIP [get]
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

// GetProbeIPStatistics 攻击业务IP统计
// @Summary 攻击业务IP统计
// @Description 攻击业务IP统计
// @Tags 大屏展示
// @Produce application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.AttackStatistics
// @Router /api/v1/attack/display/probeIP [get]
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

// GetAttackProtocolStatistics 攻击服务类型统计
// @Summary 攻击服务类型统计
// @Description 攻击服务类型统计
// @Tags 大屏展示
// @Produce application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.AttackStatistics
// @Router /api/v1/attack/display/protocol [get]
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

// GetAttackLocationStatistics 攻击位置统计
// @Summary 攻击位置统计
// @Description 攻击位置统计
// @Tags 大屏展示
// @Produce application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.AttackStatistics
// @Router /api/v1/attack/display/location [get]
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
	token := gjson.Get(ret, "token").String()
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
