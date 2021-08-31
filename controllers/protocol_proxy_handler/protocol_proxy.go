package protocol_proxy_handler

import (
    "decept-defense/controllers/comm"
    "decept-defense/internal/message_client"
    "decept-defense/models"
    "decept-defense/pkg/app"
    "decept-defense/pkg/configs"
    "decept-defense/pkg/util"
    "encoding/json"
    "github.com/astaxie/beego/validation"
    "github.com/gin-gonic/gin"
    "github.com/unknwon/com"
    "go.uber.org/zap"
    "net/http"
    "path"
    "strconv"
    "strings"
    "time"
)

type ProtocolProxyCreatePayload struct {
    ProxyPort           int32             `json:"ProxyPort" binding:"required"`
    HoneypotID          int64             `json:"HoneypotID" binding:"required"`
}

// CreateProtocolProxy 创建协议代理
// @Summary 创建协议代理
// @Description 创建协议代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param ProxyPort body ProtocolProxyCreatePayload true "ProxyPort"
// @Param HoneypotID body ProtocolProxyCreatePayload true "HoneypotID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 1007 {string} json "{"code":1007,"msg":"Redis异常","data":{}}"
// @Failure 6001 {string} json "{"code":6001,"msg":"协议代理创建失败","data":{}}"
// @Failure 6006 {string} json "{"code":6006,"msg":"代理端口重复","data":{}}"
// @Router /api/v1/proxy/protocol [post]
func CreateProtocolProxy(c *gin.Context) {
    appG := app.Gin{C: c}
    var protocolProxy models.ProtocolProxy
    var taskPayload comm.ProxyTaskPayload
    var protocol models.Protocols
    var honeypot models.Honeypot
    var payload ProtocolProxyCreatePayload
    err := c.ShouldBindJSON(&payload)
    if err != nil{
      appG.Response(http.StatusOK, app.InvalidParams, nil)
      return
    }
    currentUser, exist := c.Get("currentUser")
    if !exist{
      appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
      return
    }

    p, _ :=  protocolProxy.GetProtocolProxyByProxyPort(payload.ProxyPort)
    if p != nil{
        zap.L().Error("协议代理端口重复")
        appG.Response(http.StatusOK, app.ErrorProxyPortDup, nil)
        return
    }

    protocolProxy.CreateTime = util.GetCurrentTime()
    protocolProxy.Creator =  *(currentUser.(*string))
    protocolProxy.ProxyPort = payload.ProxyPort

    id, err := util.GetUniqueID()
    if err != nil{
      appG.Response(http.StatusOK, app.ErrorUUID, nil)
      return
    }
    protocolProxy.TaskID = id
    protocolProxy.HoneypotID = payload.HoneypotID
    r, err := honeypot.GetHoneypotByID(payload.HoneypotID)
    if err != nil{
      appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, nil)
      return
    }
    s, err := protocol.GetProtocolByType(r.ServerType)
    if err != nil{
      appG.Response(http.StatusOK, app.ErrorProtocolNotExist, nil)
      return
    }
    server, err :=  (&(models.HoneypotServers{})).GetServerStatusByID(r.ServersID)
    if err != nil{
        zap.L().Error(err.Error())
        appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, err.Error())
        return
    }
    protocolProxy.ProtocolID = s.ID
    {
        taskPayload.TaskType = comm.ProtocolProxy
        taskPayload.OperatorType = comm.DEPLOY
        taskPayload.HoneypotPort = r.ServerPort
        taskPayload.ProxyPort = payload.ProxyPort
        taskPayload.TaskID = id
        taskPayload.HoneypotIP = r.HoneypotIP
        taskPayload.DeployPath = path.Join(s.DeployPath, s.FileName)
        taskPayload.AgentID = server.AgentID
        jsonByte, _ :=  json.Marshal(taskPayload)
        err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
        if err != nil{
            appG.Response(http.StatusOK, app.ErrorRedis, nil)
            return
        }
    }
    protocolProxy.AgentID = server.AgentID
    protocolProxy.Status = comm.RUNNING
    if err := protocolProxy.CreateProtocolProxy(); err != nil{
      appG.Response(http.StatusOK, app.ErrorProtocolProxyCreate, nil)
      return
    }
    timeTicker := time.NewTicker(50 * time.Microsecond)
    count := 0
    for {
        <-timeTicker.C
        count++
        proxy, err :=  protocolProxy.GetProtocolProxyByTaskID(protocolProxy.TaskID)
        if err == nil{
            if proxy.Status == comm.FAILED{
                timeTicker.Stop()
                protocolProxy.DeleteProtocolProxyByID(proxy.ID)
                appG.Response(http.StatusOK, app.ErrorProtocolProxyFail, nil)
                return

            }else if  proxy.Status == comm.SUCCESS{
                timeTicker.Stop()
                appG.Response(http.StatusOK, app.SUCCESS, nil)
                return
            }
        }
        if count >= 1000 {
            timeTicker.Stop()
            protocolProxy.DeleteProtocolProxyByID(proxy.ID)
            appG.Response(http.StatusOK, app.ErrorProtocolProxyFail, nil)
            return
        }
    }
    appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetProtocolProxy 查找协议代理
// @Summary 查找协议代理
// @Description 查找协议代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param Payload body comm.SelectPayload false "Payload"
// @Param PageNumber body comm.SelectPayload true "PageNumber"
// @Param PageSize body comm.SelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.ProtocolProxySelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/protocol/set [post]
func GetProtocolProxy(c *gin.Context) {
    appG := app.Gin{C: c}
    var payload comm.SelectPayload
    var record models.ProtocolProxy
    err := c.ShouldBindJSON(&payload)
    if err != nil{
        appG.Response(http.StatusOK, app.InvalidParams, nil)
        return
    }
    data, count, err := record.GetProtocolProxy(&payload)
    if err != nil{
        appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
        return
    }
    appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// DeleteProtocolProxy 删除协议代理
// @Summary 删除协议代理
// @Description 删除协议代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/protocol/:id [delete]
func DeleteProtocolProxy(c *gin.Context) {
    appG := app.Gin{C: c}
    valid := validation.Validation{}
    id := com.StrTo(c.Param("id")).MustInt64()
    valid.Min(id, 1, "id").Message("ID必须大于0")
    if valid.HasErrors() {
        appG.Response(http.StatusOK, app.InvalidParams, nil)
        return
    }

    var protocolProxy models.ProtocolProxy
    var honeypot      models.Honeypot
    var honeypotServer models.HoneypotServers
    r, err := protocolProxy.GetProtocolProxyByID(id)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
        return
    }

    s, err := honeypot.GetHoneypotByID(r.HoneypotID)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, nil)
        return
    }
    p, err := honeypotServer.GetServerByID(s.ServersID)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, nil)
        return
    }
    var taskPayload comm.ProxyTaskPayload
    {
        taskPayload.TaskType = comm.ProtocolProxy
        taskPayload.OperatorType = comm.WITHDRAW
        taskPayload.HoneypotPort = s.ServerPort
        taskPayload.ProxyPort = r.ProxyPort
        taskPayload.TaskID = r.TaskID
        taskPayload.HoneypotIP = s.HoneypotIP
        taskPayload.AgentID = p.AgentID
        jsonByte, _ :=  json.Marshal(taskPayload)
        err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
        if err != nil{
            appG.Response(http.StatusOK, app.ErrorRedis, nil)
            return
        }
    }
    if err := protocolProxy.DeleteProtocolProxyByID(id); err != nil{
        appG.Response(http.StatusOK, app.ErrorDatabase, nil)
        return
    }
    appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// OfflineProtocolProxy 下线协议代理
// @Summary 下线协议代理
// @Description 下线协议代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/protocol/offline/:id [post]
func OfflineProtocolProxy(c *gin.Context) {
    appG := app.Gin{C: c}
    valid := validation.Validation{}
    id := com.StrTo(c.Param("id")).MustInt64()
    valid.Min(id, 1, "id").Message("ID必须大于0")
    if valid.HasErrors() {
        appG.Response(http.StatusOK, app.InvalidParams, nil)
        return
    }

    var protocolProxy models.ProtocolProxy
    var honeypot      models.Honeypot
    var honeypotServer models.HoneypotServers
    var protocol      models.Protocols
    r, err := protocolProxy.GetProtocolProxyByID(id)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
        return
    }

    if r.Status == comm.FAILED{
        appG.Response(http.StatusOK, app.ErrorProtocolOffline, nil)
        return
    }

    s, err := honeypot.GetHoneypotByID(r.HoneypotID)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, nil)
        return
    }
    p, err := honeypotServer.GetServerByID(s.ServersID)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, nil)
        return
    }
    t, err := protocol.GetProtocolByType(s.ServerType)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorProtocolNotExist, nil)
        return
    }
    var taskPayload comm.ProxyTaskPayload
    {
        taskPayload.TaskType = comm.ProtocolProxy
        taskPayload.OperatorType = comm.WITHDRAW
        taskPayload.HoneypotPort = s.ServerPort
        taskPayload.ProxyPort = r.ProxyPort
        taskPayload.TaskID = r.TaskID
        taskPayload.DeployPath = path.Join(t.DeployPath, t.FileName)
        taskPayload.HoneypotIP = s.HoneypotIP
        taskPayload.AgentID = p.AgentID
        jsonByte, _ :=  json.Marshal(taskPayload)
        err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
        if err != nil{
            appG.Response(http.StatusOK, app.ErrorRedis, nil)
            return
        }
    }
    timeTicker := time.NewTicker(50 * time.Microsecond)
    count := 0
    for {
        <-timeTicker.C
        count++
        proxy, err :=  protocolProxy.GetProtocolProxyByTaskID(taskPayload.TaskID)
        if err == nil{
            if proxy.Status == comm.FAILED{
                timeTicker.Stop()
                appG.Response(http.StatusOK, app.SUCCESS, nil)
                return

            }
        }
        if count >= 1000 {
            timeTicker.Stop()
            appG.Response(http.StatusOK, app.ErrorProtocolProxyOfflineFail, nil)
            return
        }
    }
    appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// OnlineProtocolProxy 上线协议代理
// @Summary 上线协议代理
// @Description 上线协议代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/protocol/online/:id [post]
func OnlineProtocolProxy(c *gin.Context) {
    appG := app.Gin{C: c}
    valid := validation.Validation{}
    id := com.StrTo(c.Param("id")).MustInt64()
    valid.Min(id, 1, "id").Message("ID必须大于0")
    if valid.HasErrors() {
        appG.Response(http.StatusOK, app.InvalidParams, nil)
        return
    }

    var protocolProxy models.ProtocolProxy
    var honeypot      models.Honeypot
    var honeypotServer models.HoneypotServers
    var protocol      models.Protocols
    r, err := protocolProxy.GetProtocolProxyByID(id)
    if err != nil{
        appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
        return
    }
    if r.Status == comm.SUCCESS{
        appG.Response(http.StatusOK, app.ErrorProtocolOnline, nil)
        return
    }

    s, err := honeypot.GetHoneypotByID(r.HoneypotID)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, nil)
        return
    }
    p, err := honeypotServer.GetServerByID(s.ServersID)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, nil)
        return
    }
    t, err := protocol.GetProtocolByType(s.ServerType)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorProtocolNotExist, nil)
        return
    }
    var taskPayload comm.ProxyTaskPayload
    {
        taskPayload.TaskType = comm.ProtocolProxy
        taskPayload.OperatorType = comm.DEPLOY
        taskPayload.HoneypotPort = s.ServerPort
        taskPayload.ProxyPort = r.ProxyPort
        taskPayload.TaskID = r.TaskID
        taskPayload.HoneypotIP = s.HoneypotIP
        taskPayload.DeployPath = path.Join(t.DeployPath, t.FileName)
        taskPayload.AgentID = p.AgentID
        jsonByte, _ :=  json.Marshal(taskPayload)
        err = message_client.PublishMessage(configs.GetSetting().App.TaskChannel, string(jsonByte))
        if err != nil{
            appG.Response(http.StatusOK, app.ErrorRedis, nil)
            return
        }
    }
    timeTicker := time.NewTicker(50 * time.Microsecond)
    count := 0
    for {
        <-timeTicker.C
        count++
        proxy, err :=  protocolProxy.GetProtocolProxyByTaskID(taskPayload.TaskID)
        if err == nil{
            if proxy.Status == comm.SUCCESS{
                timeTicker.Stop()
                appG.Response(http.StatusOK, app.SUCCESS, nil)
                return

            }
        }
        if count >= 1000 {
            timeTicker.Stop()
            appG.Response(http.StatusOK, app.ErrorProtocolProxyOnlineFail, nil)
            return
        }
    }
    appG.Response(http.StatusOK, app.SUCCESS, nil)
}


// TestProtocolProxy 测试协议代理
// @Summary 测试协议代理
// @Description 测试协议代理
// @Tags 影子代理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/proxy/protocol/test/:id [get]
func TestProtocolProxy(c *gin.Context) {
    appG := app.Gin{C: c}
    valid := validation.Validation{}
    id := com.StrTo(c.Param("id")).MustInt64()
    valid.Min(id, 1, "id").Message("ID必须大于0")
    if valid.HasErrors() {
        appG.Response(http.StatusOK, app.InvalidParams, nil)
        return
    }

    var protocolProxy models.ProtocolProxy
    var honeypot      models.Honeypot
    var honeypotServer models.HoneypotServers
    r, err := protocolProxy.GetProtocolProxyByID(id)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorProtocolProxyNotExist, nil)
        return
    }

    s, err := honeypot.GetHoneypotByID(r.HoneypotID)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorHoneypotNotExist, nil)
        return
    }
    p, err := honeypotServer.GetServerByID(s.ServersID)
    if err != nil{
        appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, nil)
        return
    }
    if util.TcpGather(strings.Split(p.ServerIP, ","),   strconv.Itoa(int(r.ProxyPort))){
        protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.SUCCESS, r.TaskID)
        appG.Response(http.StatusOK, app.SUCCESS, nil)
        return
    }else{
        appG.Response(http.StatusOK, app.ErrorConnectTest, nil)
        protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.FAILED, r.TaskID)
        return
    }
}

func UpdateProtocolProxyStatus() {
    var protocolProxy    models.ProtocolProxy
    var honeypot         models.Honeypot
    var honeypotServer   models.HoneypotServers
    r, err := protocolProxy.GetProtocolProxies()
    if err != nil{
        return
    }
    for _, d := range *r{
        if d.Status == comm.RUNNING{
            continue
        }
        s, err := honeypot.GetHoneypotByID(d.HoneypotID)
        if err != nil{
            continue
        }
        p, err := honeypotServer.GetServerByID(s.ServersID)
        if err != nil{
            continue
        }
        if util.TcpGather(strings.Split(p.ServerIP, ","),   strconv.Itoa(int(d.ProxyPort))){
            protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.SUCCESS, d.TaskID)
        }else{
            protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.FAILED, d.TaskID)
            return
        }
    }
}

