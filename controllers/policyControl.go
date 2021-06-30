package controllers

import (
	_ "decept-defense/models/util"
)

//type BaitPolicy struct {
//	beego.Controller
//}
//type TransparentTransponderPolicy struct {
//	beego.Controller
//}
//type AttackLog struct {
//	beego.Controller
//}
//type Config struct {
//	beego.Controller
//}

// 查询诱饵策略传入的json
type SelectBaitPolicyJson struct {
	AgentId          string `json:"agentId"`
	PageSize         int    `json:"pageSize"`
	PageNum          int    `json:"pageNum"`
	BaitType         string `json:"baitType"`
	Creator          string `json:"creator"`
	CreateStartTime  string `json:"createStartTime"`
	CreateEndTime    string `json:"createEndTime"`
	OfflineStartTime string `json:"offlineStartTime"`
	OfflineEndTime   string `json:"offlineEndTime"`
	Status           int    `json:"status"`
}

// 查询透明转发策略传入的json
type SelectTransPolicyJson struct {
	AgentId          string `json:"agentId"`
	PageSize         int    `json:"pageSize"`
	PageNum          int    `json:"pageNum"`
	Creator          string `json:"creator"`
	CreateStartTime  string `json:"createStartTime"`
	CreateEndTime    string `json:"createEndTime"`
	OfflineStartTime string `json:"offlineStartTime"`
	OfflineEndTime   string `json:"offlineEndTime"`
	ForwardPort      int    `json:"forwardPort"`
	HoneyPotTypeId   string `json:"honeyPotTypeId"`
	HoneyPotPort     int    `json:"honeyPotPort"`
	Status           int    `json:"status"`
}

// 下发诱饵策略传入的json
type BaitJson struct {
	AgentId  string
	BaitId   int
	Creator  string
	Status   int
	Address  string
	Md5      string
	Baitinfo *BaitInfo
}

type BaitInfo struct {
}

// 诱饵策略（下发到Redis）
type BaitPolicyJson struct {
	TaskId  string
	AgentId string
	Address string
	Md5     string
	Status  int
}

// 上线透明转发策略传入的json
type TransparentTransponderJson struct {
	AgentId      string
	HoneyTypeId  string
	ForwardPort  int
	HoneyPotId   string
	HoneyPotPort int
	Creator      string
	Status       int
}

// 下线透明转发策略的json
type TransOfflineJson struct {
	AgentId string `json:"agentId"`
	TaskId  string `json:"taskId"`
	Status  int    `json:"status"`
}
type ConfJson struct {
	Id        int
	ConfName  string
	ConfValue string
	PageSize  int
	PageNum   int
}

// 透明转发策略（下发到Redis）
type TransparentTransponderPolicyJson struct {
	TaskId     string
	AgentId    string
	ListenPort int
	ServerType string
	HoneyIP    string
	HoneyPort  int
	Status     int
}

// 攻击日志列表-前端json
type AttackLogListJson struct {
	HoneyPotTypeId string
	SrcHost        string
	AttackIP       string
	HoneyIP        string
	StartTime      string
	EndTime        string
	PageSize       int
	PageNum        int
}

// 攻击日志详情 - 前端json
type AttackLogDetailJson struct {
	Id             int
	SrcHost        string
	HoneyPotPort   int
	HoneyPotTypeId string
	AttackIP       string
	HoneyIP        string
	StartTime      string
	EndTime        string
	EventDetail    string
	PageSize       int
	PageNum        int
}

// 新增诱饵策略，insert数据库，下发策略到Redis,返回taskid
//func (this *BaitPolicy) Add(){
//	// 前端json转换为结构体
//	var bait BaitJson
//	data := this.Ctx.Input.RequestBody
//	log.Println(string(data))
//	err := json.Unmarshal(data,&bait)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	if bait.Status == 0{
//		datas := util.ReturnJsonsError("请确认诱饵策略的状态")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	// 前端json转换成策略结构体
//	var baitPolicyJson BaitPolicyJson
//	err1 := json.Unmarshal(data,&baitPolicyJson)
//	if err1 != nil{
//		return
//	}
//	baitInfo, err2 := json.Marshal(bait.Baitinfo)
//	if err2 != nil {
//		fmt.Println("json.Marshal failed:", err)
//		datas := util.ReturnJsonsError("BaitInfo 解析失败")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	taskid :=util.GetUUID()
//	currentTime:=time.Now()
//	createTime := currentTime.Format("2006-01-02 15:04:05")
//	policyCenter.InsertBaitPolicy(taskid,bait.AgentId,bait.BaitId,string(baitInfo),createTime,bait.Creator,bait.Status)
//	// 策略结构体转换成策略json
//	baitPolicyJson.TaskId = taskid
//	baitPolicy, err3 := json.Marshal(baitPolicyJson)
//	if err3 != nil {
//		fmt.Println("json.Marshal failed:", err)
//		datas := util.ReturnJsonsError("策略结构体转换失败")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
//	result := map[string]interface{}{}
//	results := util.ReturnJsonsSuccess(result)
//	this.Data["json"] = results
//	this.ServeJSON()
//}
//
//// 删除诱饵策略，update数据库，下发策略到Redis
//func (this *BaitPolicy) Delete(){
//	// 前端json转换为结构体
//	var bait BaitPolicyJson
//	data := this.Ctx.Input.RequestBody
//	log.Println(string(data))
//	err := json.Unmarshal(data,&bait)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	if bait.Status == 1{
//		datas := util.ReturnJsonsError("请确认诱饵策略的状态")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	// 前端json转换成策略结构体
//	var baitPolicyJson BaitPolicyJson
//	err1 := json.Unmarshal(data,&baitPolicyJson)
//	if err1 != nil{
//		fmt.Println("json.Marshal failed:", err)
//		datas := util.ReturnJsonsError("json解析失败")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	currentTime:=time.Now()
//	offlineTime := currentTime.Format("2006-01-02 15:04:05")
//	policyCenter.UpdateBaitPolicy(bait.TaskId,offlineTime,bait.Status)
//	//策略结构体转换成策略json
//	baitPolicy, err2 := json.Marshal(baitPolicyJson)
//	if err2 != nil {
//		fmt.Println("json.Marshal failed:", err)
//		datas := util.ReturnJsonsError("策略结构体转换失败")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
//	result := map[string]interface{}{}
//	results := util.ReturnJsonsSuccess(result)
//	this.Data["json"] = results
//	this.ServeJSON()
//}
//
//// 新增转发策略，insert数据库，下发策略到Redis
//func (this *TransparentTransponderPolicy) Add(){
//	// 前端json转换为结构体
//	var transJson TransparentTransponderJson
//	data := this.Ctx.Input.RequestBody
//	log.Println(string(data))
//	err := json.Unmarshal(data,&transJson)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		log.Println(transJson)
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	if transJson.Status == 0{
//		datas := util.ReturnJsonsError("请确认诱饵策略的状态")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	// 前端json转换成策略结构体
//	var transPolicyJson TransparentTransponderPolicyJson
//	taskid :=util.GetUUID()
//	currentTime:=time.Now()
//	createTime := currentTime.Format("2006-01-02 15:04:05")
//	policyCenter.InsertTransPolicy(taskid,transJson.AgentId,transJson.ForwardPort,transJson.HoneyPotPort,transJson.HoneyPotId,createTime,transJson.Creator,transJson.Status)
//	transPolicyJson.ServerType = policyCenter.SelectHoneyPotType(transJson.HoneyTypeId)
//	transPolicyJson.HoneyIP = policyCenter.SelectHoneyPotIP(transJson.HoneyPotId)
//	transPolicyJson.TaskId = taskid
//	transPolicyJson.HoneyPort = transJson.HoneyPotPort
//	transPolicyJson.ListenPort = transJson.ForwardPort
//	transPolicyJson.Status = transJson.Status
//	transPolicyJson.AgentId = transJson.AgentId
//	// 策略结构体转换成策略json
//	//transPolicy, err2 := json.Marshal(transPolicyJson)
//	//if err2 != nil {
//	//	fmt.Println("json.Marshal failed:", err)
//	//}
//	//redisCenter.RedisSubProducerTransPolicy(string(transPolicy))
//	result := map[string]interface{}{}
//	results := util.ReturnJsonsSuccess(result)
//	this.Data["json"] = results
//	this.ServeJSON()
//}
//
//// 删除转发策略，insert数据库，下发策略到Redis
//func (this *TransparentTransponderPolicy) Delete(){
//	// 前端json转换为结构体
//	var transJson TransOfflineJson
//	data := this.Ctx.Input.RequestBody
//	log.Println(string(data))
//	err := json.Unmarshal(data,&transJson)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		log.Println(transJson)
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	if transJson.Status == 1{
//		datas := util.ReturnJsonsError("请确认诱饵策略的状态")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	currentTime:=time.Now()
//	offlineTime := currentTime.Format("2006-01-02 15:04:05")
//	policyCenter.UpdateTransPolicy(transJson.TaskId,offlineTime,transJson.Status)
//	policy := policyCenter.SelectHoneyPotInfoByTaskId(transJson.TaskId)
//	// 策略结构体转换成策略json
//	//transPolicy, err2 := json.Marshal(transPolicyJson)
//	//if err2 != nil {
//	//	fmt.Println("json.Marshal failed:", err)
//	//}
//	redisCenter.RedisSubProducerTransPolicy(policy)
//	result := map[string]interface{}{}
//	results := util.ReturnJsonsSuccess(result)
//	this.Data["json"] = results
//	this.ServeJSON()
//}
//
//// 查看指定agent的所有诱饵策略
//func (this *BaitPolicy) Select(){
//	var selectBait SelectBaitPolicyJson
//	data := this.Ctx.Input.RequestBody
//	err := json.Unmarshal(data,&selectBait)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	datas := policyCenter.SelectBaitPolicy(selectBait.AgentId,selectBait.BaitType,selectBait.Creator,selectBait.CreateStartTime,selectBait.CreateEndTime,selectBait.OfflineStartTime,selectBait.OfflineEndTime,selectBait.Status,selectBait.PageSize,selectBait.PageNum)
//	this.Data["json"] = datas
//	this.ServeJSON()
//	return
//}
//
//// 查看指定agent的所有转发策略
//func (this *TransparentTransponderPolicy) SelectAgentId() {
//	var agentId SelectTransPolicyJson
//	data := this.Ctx.Input.RequestBody
//	err := json.Unmarshal(data,&agentId)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	datas := policyCenter.SelectTransPolicy(agentId.AgentId,agentId.ForwardPort,agentId.HoneyPotPort,agentId.CreateStartTime,agentId.CreateEndTime,agentId.OfflineStartTime,agentId.OfflineEndTime,agentId.Creator,agentId.Status,agentId.HoneyPotTypeId,agentId.PageSize,agentId.PageNum)
//	this.Data["json"] = datas
//	this.ServeJSON()
//	return
//}

// 查询攻击事件列表
//func (this *AttackLog) SelectList() {
//	var attackLogList AttackLogListJson
//	data := this.Ctx.Input.RequestBody
//	err := json.Unmarshal(data,&attackLogList)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	datas := policyCenter.SelectAttackLogList(attackLogList.HoneyPotTypeId,attackLogList.SrcHost,attackLogList.AttackIP,attackLogList.StartTime,attackLogList.EndTime,attackLogList.PageSize,attackLogList.PageNum,attackLogList.HoneyIP)
//	this.Data["json"] = datas
//	this.ServeJSON()
//	return
//}
//
//// 查询攻击事件详情
//func (this *AttackLog) SelectDetail() {
//	var attackLogDetail AttackLogDetailJson
//	data := this.Ctx.Input.RequestBody
//	err := json.Unmarshal(data,&attackLogDetail)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	datas := policyCenter.SelectAttackLogDetail(attackLogDetail.Id,attackLogDetail.HoneyPotTypeId,attackLogDetail.SrcHost,attackLogDetail.AttackIP,attackLogDetail.StartTime,attackLogDetail.EndTime,attackLogDetail.HoneyPotPort,attackLogDetail.HoneyIP,attackLogDetail.PageSize,attackLogDetail.PageNum)
//	this.Data["json"] = datas
//	this.ServeJSON()
//	return
//}
//
//// 添加钉钉告警
//func (this *Config) AddConf() {
//	var config ConfJson
//	data := this.Ctx.Input.RequestBody
//	err := json.Unmarshal(data,&config)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	if config.ConfName != "DingDing" && config.ConfName != "Message"{
//		datas := util.ReturnJsonsError("类型错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}else {
//		result := policyCenter.InsertConfig(config.ConfName,config.ConfValue)
//		this.Data["json"] = result
//		this.ServeJSON()
//		return
//	}
//
//}
//
//// 查询钉钉告警
//func (this *Config) SelectConf() {
//	var config ConfJson
//	data := this.Ctx.Input.RequestBody
//	err := json.Unmarshal(data,&config)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	if config.ConfName != "DingDing" && config.ConfName != "Message" {
//		datas := util.ReturnJsonsError("类型错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}else {
//		result := policyCenter.SelectConfig(config.Id,config.ConfName,config.PageSize,config.PageNum)
//		this.Data["json"] = result
//		this.ServeJSON()
//		return
//	}
//}
//
//func (this *Config) DeleteConf(){
//	var config ConfJson
//	data := this.Ctx.Input.RequestBody
//	err := json.Unmarshal(data,&config)
//	if err != nil{
//		datas := util.ReturnJsonsError("数据输入错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}
//	if config.ConfName != "DingDing" && config.ConfName != "Message" {
//		datas := util.ReturnJsonsError("类型错误")
//		this.Data["json"] = datas
//		this.ServeJSON()
//		return
//	}else {
//		result,msg := policyCenter.DeleteConfig(config.Id,config.ConfName)
//		this.Data["json"] = result
//		this.ServeJSON()
//		return
//	}
//}
