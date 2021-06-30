package policyCenter

import (
	"database/sql"
	"decept-defense/models/util"
	"decept-defense/models/util/comm"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	"log"
	"strconv"
	"time"
)

type TblPolicyBait struct {
	id          int       `orm:"column(id);auto" description:"诱饵策略id"`
	AgentId     string    `orm:"column(agentid)" description:"agentId"`
	BaitId      string    `orm:"column(baitid);rel(fk)" description:"诱饵ID"`
	TaskId      string    `orm:"column(taskid)" description:"任务ID"`
	CreateTime  time.Time `orm:"column(createtime)" description:"操作时间"`
	Creator     string    `orm:"column(creator)" description:"操作人员"`
	Status      int       `orm:"column(status)" description:"是否启用"`
	OfflineTime time.Time `orm:"column(offlinetime)" description:"下线时间"`
	BaitInfo    string    `orm:"column(baitinfo)"`
	Address     string    `orm:"column(address)"`
	Md5         string    `orm:"column(md5)"`
}
type TblPolicyTrans struct {
	TaskId       string    `orm:"column(taskid)"`
	AgentId      string    `orm:"column(agentid)" description:"agentId"`
	ForwardPort  string    `orm:"column(forwardport)"`
	HoneyPotPort string    `orm:"column(honeypotport)"`
	HoneyPotId   string    `orm:"column(honeypotid);rel(fk)"`
	Status       int       `orm:"column(status)" description:"是否启用"`
	Creator      string    `orm:"column(creator)" description:"创建人员"`
	CreateTime   time.Time `orm:"column(createtime)" description:"创建时间"`
	OfflineTime  time.Time `orm:"column(offlinetime)"`
	Type         int       `orm:"column(type)"`
	Path         string    `orm:"column(path)"`
}
type TblPolicyHoneyTrans struct {
	TaskId       string    `orm:"column(taskid)"`
	AgentId      string    `orm:"column(agentid)" description:"agentId"`
	ForwardPort  string    `orm:"column(forwardport)"`
	HoneyPotPort string    `orm:"column(honeypotport)"`
	HoneyPotId   string    `orm:"column(honeypotid);rel(fk)"`
	Status       int       `orm:"column(status)" description:"是否启用"`
	Creator      string    `orm:"column(creator)" description:"创建人员"`
	CreateTime   time.Time `orm:"column(createtime)" description:"创建时间"`
	OfflineTime  time.Time `orm:"column(offlinetime)"`
	Type         int       `orm:"column(type)"`
	Path         string    `orm:"column(path)"`
}
type TblBaits struct {
	BaitType   string `orm:"column(baittype)"`
	Filename   string `orm:"column(filename)"`
	CreateTime string `orm:"column(createtime)"`
	Creator    string `orm:"column(creator)"`
	Baitid     string `orm:"column(baitid);pk"`
}

type TblHoneyPotsType struct {
	HoneyPotType string `orm:"column(honeypottype)"`
	TypeId       string `orm:"column(typeid);pk"`
}
type TblSysType struct {
	SysId   string `orm:"column(sysid);pk"`
	SysType string `orm:"column(systype)"`
}
type TblHoneyPots struct {
	HoneyName   string    `orm:"column(honeyname)"`
	HoneyPotId  string    `orm:"column(honeypotid);pk"`
	HoneyIP     string    `orm:"column(honeyip)"`
	HoneyTypeId string    `orm:"column(honeytypeid);rel(fk)"`
	ServerId    string    `orm:"column(serverid);rel(fk)"`
	SysId       string    `orm:"column(sysid);rel(fk)"`
	Status      string    `orm:"column(status)"`
	Creator     string    `orm:"column(creator)"`
	CreateTime  time.Time `orm:"column(createtime)" description:"创建时间"`
}
type TblHoneyPotServers struct {
	HoneyName string `orm:"column(servername)"`
	ServerId  string `orm:"column(serverid);pk"`
	ServerIP  string `orm:"column(serverip)"`
	Status    string `orm:"column(status)"`
}

type TransPolicy struct {
	TaskId     string `json:"taskid"`
	AgentId    string `json:"agentid"`
	ListenPort string `json:"forwardport"`
	ServerType string `json:"honeypottype"`
	Serverip   string `json:"serverip"`
	HoneyPort  string `json:"honeypotport"`
	Status     string `json:"status"`
}
type TransPolicyJson struct {
	TaskId     string
	AgentId    string
	ListenPort int
	ServerType string
	HoneyIP    string
	HoneyPort  int
	Status     int
	Type       string
	Path       string
}
type TblAttackLog struct {
	Id           int        `orm:"column(id)"`
	SrcHost      string     `orm:"column(srchost)"`
	SrcPost      int        `orm:"column(srcport)"`
	HoneyPotId   string     `orm:"column(honeypotid)"`
	HoneyPotPort int        `orm:"column(honeypotport)"`
	AttackIP     string     `orm:"column(attackip)"`
	AttackTime   *time.Time `orm:"column(attacktime)"`
	EventDetail  string     `orm:"column(eventdetail)"`
	ProxyType    string     `orm:"column(proxytype)"`
	SourceType   int        `orm:"column(sourcetype)"`
}

func (t *TblHoneyPots) TableName() string {
	return "honeypots"
}
func (t *TblPolicyBait) TableName() string {
	return "server_bait"
}
func (t *TblPolicyTrans) TableName() string {
	return "fowards"
}
func (t *TblBaits) TableName() string {
	return "baits"
}
func (t *TblPolicyHoneyTrans) TableName() string {
	return "honeyfowards"
}
func (t *TblHoneyPotsType) TableName() string {
	return "honeypotstype"
}
func (t *TblHoneyPotServers) TableName() string {
	return "honeypotservers"
}
func (t *TblSysType) TableName() string {
	return "systemtype"
}

var (
	dbhost     = beego.AppConfig.String("dbhost")
	dbport     = beego.AppConfig.String("dbport")
	dbuser     = beego.AppConfig.String("dbuser")
	dbpassword = beego.AppConfig.String("dbpassword")
	dbname     = beego.AppConfig.String("dbname")
)

func InsertBaitPolicy(taskid string, agentid string, baitid string, baitinfo string, createtime int64, creator string, status int, data string, md5 string, types string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into server_bait (taskid,agentid,baitid,baitinfo,createtime,creator,data,md5,type) VALUES (?,?,?,?,?,?,?,?,?)", taskid, agentid, baitid, baitinfo, createtime, creator, data, md5, types).Values(&maps)
	if err != nil {
		logs.Error("[InsertBaitPolicy] insert bait policy error,%s", err)
	}
}

func InsertSignPolicy(taskid string, agentid string, baitid string, baitinfo string, createtime int64, creator string, status int, data string, md5 string, types string, tracecode string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into server_sign (taskid,agentid,signid,signinfo,createtime,creator,md5,type,tracecode, status) VALUES (?,?,?,?,?,?,?,?,?,?)", taskid, agentid, baitid, data, createtime, creator, md5, types, tracecode, status).Values(&maps)
	if err != nil {
		logs.Error("[InsertBaitPolicy] insert bait policy error,%s", err)
	}
}

func OffSignPolicy(taskid string, status int, offtime int64) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update server_sign set status=?, offlinetime=? where taskid=?", status, offtime, taskid).Values(&maps)
	if err != nil {
		logs.Error("[OffSignPolicy] update sign policy error,%s", err)
	}
}

func UpdateSignPolicy(status int, taskid string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update server_sign set status=? where taskid=?", status, taskid).Values(&maps)
	if err != nil {
		logs.Error("[UpdateHoneySign] update sign policy error,%s", err)
	}
}

func UpdateBaitPolicy(taskid string, offlinetime int64, status int) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update server_bait set offlinetime=?,status=? where taskid=?", offlinetime, status, taskid).Values(&maps)
	if err != nil {
		logs.Error("[UpdateBaitPolicy] insert bait policy error,%s", err)
	}
}

// 写入透明转发策略
func InsertTransPolicy(taskid string, agentid string, forwardport int, honeypotport int, honeyserveragentid string, createtime int64, creator string, status int, agenttype string, path string, honeytypeid string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into fowards (taskid,agentid,forwardport,honeypotport,serverid,createtime,creator,type,path,honeytypeid,status) VALUES (?,?,?,?,?,?,?,?,?,?,?)", taskid, agentid, forwardport, honeypotport, honeyserveragentid, createtime, creator, agenttype, path, honeytypeid, status).Values(&maps)
	if err != nil {
		logs.Error("[InsertTransPolicy] insert bait policy error,%s", err)
	}
}

//
func UpdateHoneyTransPolicyForwardStatus(taskid string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update honeyfowards set forwardstatus=2 where taskid=?", taskid).Values(&maps)
	if err != nil {
		logs.Error("[UpdateHoneyTransPolicyForwardStatus] insert bait policy error,%s", err)
	}
}

// 写入蜜罐流量转发策略
func InsertHoneyTransPolicy(taskid string, agentid string, forwardport int, honeypotport int, honeypotid string, createtime int64, creator string, status int, agenttype string, path string, serverid string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into honeyfowards (taskid,agentid,forwardport,honeypotport,honeypotid,createtime,creator,type,path,serverid,status) VALUES (?,?,?,?,?,?,?,?,?,?,?)", taskid, agentid, forwardport, honeypotport, honeypotid, createtime, creator, agenttype, path, serverid, status).Values(&maps)
	if err != nil {
		logs.Error("[InsertHoneyTransPolicy] insert bait policy error,%s", err)
	}
}

func SelectHoneyServerSSHKey(agentid string) (string, string) {
	o := orm.NewOrm()
	var maps []orm.Params
	sshkey := ""
	errmsg := ""
	_, err := o.Raw("SELECT serversshkey FROM honeyserverconfig WHERE honeyserverid =?", agentid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyServerSSHKey] select event list error,%s", err)
		errmsg = "系统错误"
		sshkey = ""
	} else {
		if len(maps) == 0 {
			sshkey = ""
			errmsg = "蜜网ssh 未插入"
		} else {
			sshkey = util.Strval(maps[0]["serversshkey"])
		}
	}
	return sshkey, errmsg
}

func SelectHoneyPotType(honeytypeid string) string {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT honeypottype from honeypotstype where typeid=?", honeytypeid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyPotType] select event list error,%s", err)
		return ""
	}
	var honeyType []TblHoneyPotsType
	re, _ := json.Marshal(maps)
	err1 := json.Unmarshal(re, &honeyType)
	if err1 != nil {
		logs.Error("[SelectHoneyPotType] oplist unmarshal error,%s", err)
		return ""
	}
	honeytype := honeyType[0].HoneyPotType
	return honeytype
}

func SelectHoneyPotTypeById(honeytypeid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * from honeypotstype where typeid=?", honeytypeid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyPotType] select event list error,%s", err)
		return maps
	}
	return maps
}

func SelectHoneyPotTypeByHoneyId(honeypotid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT a.*,b.honeypottype,b.softpath FROM honeypots a LEFT JOIN honeypotstype b ON a.honeytypeid = b.typeid WHERE a.honeypotid=?", honeypotid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyPotTypeByHoneyId] select event list error,%s", err)
		return maps
	}
	return maps
}

func SelectForwardByTaskId(taskid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("select T1.serverip AS honeyserverip, T3.serverip,T0.`forwardport`,T0.`honeypotport`,T0.`status`,T0.`taskid` from `fowards` T0 LEFT JOIN servers T3 ON T0.agentid = T3.agentid LEFT JOIN honeypotservers T1 ON T0.serverid = T1.agentid where T0.`taskid`=?", taskid).Values(&maps)
	if err != nil {
		logs.Error("[SelectForwardByTaskId] select event list error,%s", err)
		return maps
	}
	return maps
}

func SelectHoneyForwardByTaskId(taskid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("select T1.honeypotid,T1.honeyname,T1.podname,T1.honeynamespce,T0.taskid,T0.forwardport,T1.honeyip,T1.honeyport,T3.serverip from honeyfowards T0 LEFT JOIN honeypots T1 ON T0.honeypotid = T1.honeypotid LEFT JOIN honeypotservers T3 ON T1.agentid = T3.agentid where T0.taskid=?", taskid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyForwardByTaskId] select event list error,%s", err)
		return maps
	}
	return maps
}

func SelectHoneyForwardByHoneyId(honeypotid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * from honeyfowards where honeypotid=?", honeypotid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyForwardByHoneyId] select event list error,%s", err)
		return maps
	}
	return maps
}

func SelectForwardByConditions(agentid string, honeytypeid string, forwardport int, listenport int, honeyserveragentid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * from fowards where agentid=? and honeytypeid=? and forwardport=? and honeypotport=? and serverid=? and status=1", agentid, honeytypeid, forwardport, listenport, honeyserveragentid).Values(&maps)
	if err != nil {
		logs.Error("[SelectForwardByConditions] select event list error,%s", err)
		return maps
	}
	return maps
}

func SelectHoneyForwardByConditions(honeypotid string, honeypotport int, agentid string, listenport int) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * from honeyfowards where honeypotid=? and honeypotport=? and forwardport=? and agentid=? and status=1", honeypotid, honeypotport, listenport, agentid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyForwardByConditions] select event list error,%s", err)
		return maps
	}
	return maps
}

func SelectHoneyTransByHoneyId(honeypotid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	//
	_, err := o.Raw("SELECT T0.honeyip,T0.honeyport,T0.serverid,T1.serverip,T1.agentid,T2.honeypottype,T2.softpath FROM honeypots T0 LEFT JOIN honeypotservers T1 ON T0.agentid = T1.agentid LEFT JOIN honeypotstype T2 ON T0.honeytypeid = T2.typeid where T0.honeypotid = ?", honeypotid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyTransByHoneyId] select event list error,%s", err)
		return maps
	}
	return maps
}

func SelectAllHoneyPotsType() []comm.HoneyPotType {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT honeypottype,typeid as honeyTypeId from honeypotstype").Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyPotType] select event list error,%s", err)
		return nil
	}
	var honeyType []comm.HoneyPotType
	re, _ := json.Marshal(maps)
	err1 := json.Unmarshal(re, &honeyType)
	if err1 != nil {
		logs.Error("[SelectHoneyPotType] oplist unmarshal error,%s", err)
		return nil
	}
	return honeyType
}

func SelectAllBaitType() []comm.BaitType {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT baittype from baittype").Values(&maps)
	if err != nil {
		logs.Error("[SelectAllBaitType] select event list error,%s", err)
		return nil
	}
	var baitType []comm.BaitType
	re, _ := json.Marshal(maps)
	err1 := json.Unmarshal(re, &baitType)
	if err1 != nil {
		logs.Error("[SelectAllBaitType] oplist unmarshal error,%s", err)
		return nil
	}
	return baitType
}

func SelectAllHoneyBaitType() []comm.BaitType {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT baittype from baittype WHERE baittype != 'history'").Values(&maps)
	if err != nil {
		logs.Error("[SelectAllHoneyBaitType] select event list error,%s", err)
		return nil
	}
	var baitType []comm.BaitType
	re, _ := json.Marshal(maps)
	err1 := json.Unmarshal(re, &baitType)
	if err1 != nil {
		logs.Error("[SelectAllHoneyBaitType] oplist unmarshal error,%s", err)
		return nil
	}
	return baitType
}

func SelectAllSignType() []comm.SignType {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT signtype from signtype").Values(&maps)
	if err != nil {
		logs.Error("[SelectAllSignType] select event list error,%s", err)
		return nil
	}
	var signType []comm.SignType
	re, _ := json.Marshal(maps)
	err1 := json.Unmarshal(re, &signType)
	if err1 != nil {
		logs.Error("[SelectAllSignType] oplist unmarshal error,%s", err)
		return nil
	}
	return signType
}

func SelectAllSysType() []comm.SysType {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT sysid,systype from systemtype").Values(&maps)
	if err != nil {
		logs.Error("[SelectAllSysType] select event list error,%s", err)
		return nil
	}
	var sysType []comm.SysType
	re, _ := json.Marshal(maps)
	err1 := json.Unmarshal(re, &sysType)
	if err1 != nil {
		logs.Error("[SelectAllSysType] oplist unmarshal error,%s", err)
		return nil
	}
	return sysType
}

func SelectHoneyPotIP(honeypotid string) string {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT honeyip from honeypots where honeypotid=?", honeypotid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyPotIP] select event list error,%s", err)
		return ""
	}
	var honeyIP []TblHoneyPots
	re, _ := json.Marshal(maps)
	err1 := json.Unmarshal(re, &honeyIP)
	if err1 != nil {
		logs.Error("[SelectHoneyPotIP] oplist unmarshal error,%s", err)
		return ""
	}
	honeyip := honeyIP[0].HoneyIP
	return honeyip
}
func SelectHoneyServerInfo(honeypotid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * from honeypotservers T0 LEFT JOIN honeypots T1 ON T0.agentid = T1.agentid WHERE T1.honeypotid=?", honeypotid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyServerInfo] select event list error,%s", err)
		return maps
	}
	return maps
}

// 更新透明转发策略
func UpdateTransPolicy(taskid string, offlinetime int64, status int) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update fowards set status=? ,offlinetime=? where taskid=?", status, offlinetime, taskid).Values(&maps)
	if err != nil {
		logs.Error("[UpdateTransPolicy]  error,%s", err)
	}
}

// 更新蜜罐流量转发策略
func UpdateHoneyTransPolicy(taskid string, offlinetime int64, status int) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update honeyfowards set status=? ,offlinetime=? where taskid=?",status, offlinetime, taskid).Values(&maps)
	if err != nil {
		logs.Error("[UpdateHoneyTransPolicy] insert bait policy error,%s", err)
	}
}

// 根据任务ID查询蜜罐信息 - 透明转发策略
func SelectHoneyPotInfoByTaskId(taskid string) string {
	o := orm.NewOrm()
	var maps []orm.Params
	//_, err := o.Raw("SELECT T0.`status`,T0.`forwardport`,T0.`honeypotport`,T0.`taskid`,T0.`agentid`,T0.`type`,T0.`path`,T1.`honeyip`,T2.`honeypottype` from `fowards` T0 inner join `honeypots` T1 on T1.`honeypotid` = T0.`honeypotid` inner join `honeypotstype`T2 on T2.`typeid` = T1.`honeytypeid` where T0.`taskid`=?", taskid).Values(&maps)
	_, err := o.Raw("SELECT T0.`status`,T0.`forwardport`,T0.`honeypotport`,T0.`taskid`,T0.`agentid`,T0.`type`,T0.`path`,T1.serverip,T2.`honeypottype`,T2.`typeid`  from `fowards` T0 left join honeypotservers  T1 on T1.agentid = T0.agentid inner join `honeypotstype`T2 on T2.`typeid` = T0.`honeytypeid` where T0.`taskid`=?", taskid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyPotInfoByTaskId] select event list error,%s", err)

	}
	var transPolicy []TransPolicy
	var transPolicyJson TransPolicyJson
	re, _ := json.Marshal(maps)
	err1 := json.Unmarshal(re, &transPolicy)
	if err1 != nil {
		logs.Error("[SelectHoneyPotInfoByTaskId] oplist unmarshal error,%s", err)

	}
	transPolicyJson.TaskId = transPolicy[0].TaskId
	transPolicyJson.AgentId = transPolicy[0].AgentId
	transPolicyJson.HoneyIP = transPolicy[0].Serverip
	transPolicyJson.ServerType = transPolicy[0].ServerType
	transPolicyJson.Status, _ = strconv.Atoi(transPolicy[0].Status)
	transPolicyJson.ListenPort, _ = strconv.Atoi(transPolicy[0].ListenPort)
	transPolicyJson.HoneyPort, _ = strconv.Atoi(transPolicy[0].HoneyPort)
	transPolicyJson.Type = "UN_EDGE"
	transpolicy, _ := json.Marshal(transPolicyJson)
	return string(transpolicy)

}

// 根据任务ID查询蜜罐信息 - 蜜罐流量转发策略
func SelectHoneyPotInfoByTaskIdHoneyTrans(taskid string) string {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT T0.`status`,T0.`forwardport`,T0.`honeypotport`,T0.`taskid`,T0.`agentid`,T1.`honeyip`,T2.`honeypottype` from `honeyfowards` T0 left join `honeypots` T1 on T1.`honeypotid` = T0.`honeypotid` left join `honeypotstype`T2 on T2.`typeid` = T1.`honeytypeid` where T0.`taskid`=?", taskid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyPotInfoByTaskIdHoneyTrans] select event list error,%s", err)

	}
	var transPolicy []TransPolicy
	var transPolicyJson TransPolicyJson
	re, _ := json.Marshal(maps)
	err1 := json.Unmarshal(re, &transPolicy)
	if err1 != nil {
		logs.Error("[SelectHoneyPotInfoByTaskIdHoneyTrans] oplist unmarshal error,%s", err)

	}
	transPolicyJson.TaskId = transPolicy[0].TaskId
	transPolicyJson.AgentId = transPolicy[0].AgentId
	transPolicyJson.HoneyIP = transPolicy[0].Serverip
	transPolicyJson.ServerType = transPolicy[0].ServerType
	transPolicyJson.Status = 1
	transPolicyJson.ListenPort, _ = strconv.Atoi(transPolicy[0].ListenPort)
	transPolicyJson.HoneyPort, _ = strconv.Atoi(transPolicy[0].HoneyPort)
	transPolicyJson.Type = comm.AgentTypeUNRelay

	transpolicy, _ := json.Marshal(transPolicyJson)
	return string(transpolicy)
}

func SelectHoneyPotTrans(honeyport int, honeypotid string, agentid string, creator string, forwardport int, honeypottypeid string, honeyip string, status int, createStartTime string, createEndTime string, offlineStartTime string, offlineEndTime string, pageSize int, pageNum int) (map[string]interface{}, string, int) {
	var data map[string]interface{}
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		log.Printf("[SelectHoneyPotTrans]open mysql fail %v\n", err1)
		logs.Error("[SelectHoneyPotTrans]open mysql fail %s", err1)
		return data, comm.DBConnectError, comm.ErrorCode
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqltotal := "select COUNT(1) from honeyfowards T0 LEFT JOIN honeypots T1 ON T0.honeypotid = T1.honeypotid LEFT JOIN honeypotstype T2 ON T1.honeytypeid = T2.typeid where 1=1"
	sqlstr := "select T0.taskid,T0.forwardport,T2.honeypottype,T1.honeyip,T1.honeyport,T0.createtime,T0.offlinetime,T0.creator,T0.`status`,T3.serverip from honeyfowards T0 LEFT JOIN honeypots T1 ON T0.honeypotid = T1.honeypotid LEFT JOIN honeypotstype T2 ON T1.honeytypeid = T2.typeid LEFT JOIN honeypotservers T3 ON T1.agentid = T3.agentid where 1=1"
	var condition string
	var argsList []interface{}
	if honeypotid != "" {
		condition += " and T0.honeypotid = ?"
		argsList = append(argsList, honeypotid)
	}
	if honeyport != 0 {
		condition += " and T0.serverid = ?"
		argsList = append(argsList, honeyport)
	}
	if agentid != "" {
		condition += " and T1.honeyport = ?"
		argsList = append(argsList, agentid)
	}
	if forwardport != 0 {
		condition += " and T0.forwardport = ?"
		argsList = append(argsList, forwardport)
	}
	if honeypottypeid != "" {
		condition += " and T1.honeytypeid = ?"
		argsList = append(argsList, honeypottypeid)
	}
	if creator != "" {
		condition += " and T0.creator = ?"
		argsList = append(argsList, creator)
	}
	if honeyip != "" {
		condition += " and T1.honeyip = ?"
		argsList = append(argsList, honeyip)
	}
	if status != 0 {
		condition += " and T0.`status` = ?"
		argsList = append(argsList, status)
	}
	if createStartTime != "" && createEndTime != "" {
		condition += " and T0.`createtime` BETWEEN ? and ?"
		argsList = append(argsList, createStartTime)
		argsList = append(argsList, createEndTime)
	}
	if offlineStartTime != "" && offlineEndTime != "" {
		condition += " and T0.`offlinetime` BETWEEN ? and ?"
		argsList = append(argsList, offlineStartTime)
		argsList = append(argsList, offlineEndTime)
	}
	var total int
	sqltotal = sqltotal + condition
	fmt.Println("sqltotal:", sqltotal)
	err := DbCon.QueryRow(sqltotal, argsList...).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneyPotTrans] select total error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " ORDER BY T0.`createtime` desc limit ?,?"
	argsList = append(argsList, offset)
	argsList = append(argsList, pagesize)
	sqlstr = sqlstr + condition
	fmt.Println("sqlstr:", sqlstr)
	rows, err := DbCon.Query(sqlstr, argsList...)
	if err != nil {
		logs.Error("[SelectHoneyPotTrans] select list error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectHoneyPotTrans] rows.Columns() error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	list, count, err := util.GetHoneyPotTransMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectHoneyPotTrans] Unmarshal error,%s", err)
			return data, comm.DataUnmarshalError, comm.ErrorCode
		}
	} else {
		lists := []map[string]interface{}{}
		data = map[string]interface{}{"list": lists, "total": 0}
		return data, comm.DataSelectSuccess, comm.SuccessCode
	}

	return data, comm.DataSelectSuccess, comm.SuccessCode
}

// 使用taskid更新诱饵策略状态（主要是agent执行失败、无响应，变更状态）
func UpdateBaitPolicyStatusByTaskid(status int, taskid string) {
	o := orm.NewOrm()
	//var maps []orm.Params
	_, err := o.Raw("update server_bait set status=? where taskid=?", status, taskid).Exec()
	if err != nil {
		fmt.Println("[UpdateBaitPolicyStatusByTaskid] update bait policy error,%s", err)
		logs.Error("[UpdateBaitPolicyStatusByTaskid] update bait policy error,%s", err)
	}
}

// 使用taskid更新透明转发策略状态（主要是agent执行失败，变更状态）
func UpdateTransPolicyStatusByTaskid(status int, taskid string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update fowards set status=? where taskid=?", status, taskid).Values(&maps)
	fmt.Println("[UpdateTransPolicyStatusByTaskid]")
	if err != nil {
		logs.Error("[UpdateTransPolicyStatusByTaskid] insert bait policy error,%s", err)
	}
}

// 使用taskid更新蜜罐流量转发策略状态（主要是agent执行失败、无响应，变更状态）
func UpdateHoneyTransPolicyStatusByTaskid(status int, taskid string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update honeyfowards set status=? where taskid=?", status, taskid).Values(&maps)
	if err != nil {
		logs.Error("[UpdateHoneyTransPolicyStatusByTaskid] insert bait policy error,%s", err)
	}
}

// 查询密签列表
func SelectSignPolicy(agentId string, signId string, signtype string, creator string, createStartTime string, createEndTime string, offlineStartTime string, offlineEndTime string, status int, pageSize int, pageNum int, signinfo string) (map[string]interface{}, string, int) {
	var data map[string]interface{}
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		log.Printf("[SelectSignPolicy]open mysql fail %v\n", err1)
		logs.Error("[SelectSignPolicy]open mysql fail %s", err1)
		return data, comm.DBConnectError, comm.ErrorCode
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqltotal := "select count(1) from `server_sign` T0 left join `signs` T1 on T1.`signid`=T0.`signid` where 1=1"
	sqlstr := "select T0.`tracecode` ,T0.signinfo as data, T0.`taskid`,T1.`signname`,T1.`signinfo`,T0.`type`,T0.`creator`,T0.`createtime`,T0.`offlinetime`,T0.`status` from `server_sign` T0 left join `signs` T1 on T1.`signid`=T0.`signid` where 1=1"
	var condition string
	var argsList []interface{}
	if agentId != "" {
		condition += " and T0.`agentid` = ?"
		argsList = append(argsList, agentId)
	}
	if creator != "" {
		condition += " and T0.`creator` = ?"
		argsList = append(argsList, creator)
	}
	if signinfo != "" {
		condition += " and T0.`signinfo` = ?"
		argsList = append(argsList, signinfo)
	}
	if status != 0 {
		condition += " and T0.`status` = ?"
		argsList = append(argsList, status)
	}
	if createStartTime != "" && createEndTime != "" {
		condition += " and T0.`createtime` BETWEEN ? and ?"
		argsList = append(argsList, createStartTime)
		argsList = append(argsList, createEndTime)
	}
	if offlineStartTime != "" && offlineEndTime != "" {
		condition += " and T0.`offlinetime` BETWEEN ? and ?"
		argsList = append(argsList, offlineStartTime)
		argsList = append(argsList, offlineEndTime)
	}
	if signId != "" {
		condition += " and T1.`signid` = ?"
		argsList = append(argsList, signId)
	}
	if signtype != "" {
		condition += " and T0.`type` = ?"
		argsList = append(argsList, signId)
	}
	var total int
	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, argsList...).Scan(&total)
	if err != nil {
		logs.Error("[SelectSignPolicy] select total error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " ORDER BY T0.`createtime` desc limit ?,?"
	argsList = append(argsList, offset)
	argsList = append(argsList, pagesize)
	sqlstr = sqlstr + condition
	//fmt.Println(sqlstr)
	rows, err := DbCon.Query(sqlstr, argsList...)
	if err != nil {
		logs.Error("[SelectSignPolicy] select list error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectSignPolicy] rows.Columns() error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	list, count, err := util.GetSignPolicyMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectSignPolicy] Unmarshal error,%s", err)
			return data, comm.DataUnmarshalError, comm.ErrorCode
		}
	} else {
		lists := []map[string]interface{}{}
		data = map[string]interface{}{"list": lists, "total": 0}
		return data, comm.DataSelectSuccess, comm.SuccessCode
	}

	return data, comm.DataSelectSuccess, comm.SuccessCode
}

// 查询诱饵策略
func SelectBaitPolicy(agentId string, baitId string, creator string, createStartTime string, createEndTime string, offlineStartTime string, offlineEndTime string, status int, pageSize int, pageNum int, baitinfo string) (map[string]interface{}, string, int) {
	var data map[string]interface{}
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		log.Printf("[SelectBaitPolicy]open mysql fail %v\n", err1)
		logs.Error("[SelectBaitPolicy]open mysql fail %s", err1)
		return data, comm.DBConnectError, comm.ErrorCode
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqltotal := "select count(1) from `server_bait` T0 left join `baits` T1 on T1.`baitid`=T0.`baitid` where 1=1"
	sqlstr := "select T0.`type`,T0.`data`,T1.`baitinfo`,T0.`taskid`,T1.`baitname`,T0.`creator`,T0.`createtime`,T0.`offlinetime`,T0.`status` from `server_bait` T0 left join `baits` T1 on T1.`baitid`=T0.`baitid` where 1=1"
	var condition string
	var argsList []interface{}
	if agentId != "" {
		condition += " and T0.`agentid` = ?"
		argsList = append(argsList, agentId)
	} else {
		return data, "请选择agentId", comm.ErrorCode
	}
	if creator != "" {
		condition += " and T0.`creator` = ?"
		argsList = append(argsList, creator)
	}
	if baitinfo != "" {
		condition += " and T0.`baitinfo` = ?"
		argsList = append(argsList, baitinfo)
	}
	if status != 0 {
		condition += " and T0.`status` = ?"
		argsList = append(argsList, status)
	}
	if createStartTime != "" && createEndTime != "" {
		condition += " and T0.`createtime` BETWEEN ? and ?"
		argsList = append(argsList, createStartTime)
		argsList = append(argsList, createEndTime)
	}
	if offlineStartTime != "" && offlineEndTime != "" {
		condition += " and T0.`offlinetime` BETWEEN ? and ?"
		argsList = append(argsList, offlineStartTime)
		argsList = append(argsList, offlineEndTime)
	}
	if baitId != "" {
		condition += " and T1.`baitid` = ?"
		argsList = append(argsList, baitId)
	}

	var total int
	sqltotal = sqltotal + condition

	err := DbCon.QueryRow(sqltotal, argsList...).Scan(&total)
	if err != nil {
		logs.Error("[SelectBaitPolicy] select total error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " ORDER BY T0.`createtime` desc limit ?,?"
	argsList = append(argsList, offset)
	argsList = append(argsList, pagesize)
	sqlstr = sqlstr + condition
	rows, err := DbCon.Query(sqlstr, argsList...)
	if err != nil {
		logs.Error("[SelectBaitPolicy] select list error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectBaitPolicy] rows.Columns() error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list, count, err := util.GetBaitPolicyMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectBaitPolicy] Unmarshal error,%s", err)
			return data, comm.DataUnmarshalError, comm.ErrorCode
		}
	} else {
		lists := []map[string]interface{}{}
		data = map[string]interface{}{"list": lists, "total": 0}
		return data, comm.DataSelectSuccess, comm.SuccessCode
	}

	return data, comm.DataSelectSuccess, comm.SuccessCode
}

// 查询透明转发策略
func SelectTransPolicy(serverIp string, agentId string, forwardPort int, honeyPort int, createStartTime string, createEndTime string, offlineStartTime string, offlineEndTime string, creator string, status int, honeytypeid string, pageSize int, pageNum int) (map[string]interface{}, string, int) {
	var data map[string]interface{}
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectTransPolicy] open mysql fail,%s", err1)
		return data, comm.DBConnectError, comm.ErrorCode
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqltotal := "select count(1) from `fowards` T0 left join `honeypotservers` T1 on T1.`agentid` = T0.`serverid` left join `honeypotstype` T2 on T2.`typeid` = T0.`honeytypeid` LEFT JOIN servers T3 ON T0.agentid = T3.agentid where 1=1"
	sqlstr := "select T1.serverip AS honeyserverip, T3.serverip,T0.`forwardport`,T0.`honeypotport`,T0.`createtime`,T0.`offlinetime`,T0.`creator`,T0.`status`,T2.`honeypottype`,T0.`taskid` from `fowards` T0 left join `honeypotstype` T2 on T2.`typeid` = T0.`honeytypeid` LEFT JOIN servers T3 ON T0.agentid = T3.agentid LEFT JOIN honeypotservers T1 ON T0.serverid = T1.agentid  where 1=1 "
	var condition string
	var argsList []interface{}
	if agentId != "" {
		condition += " and T0.`agentid` = ?"
		argsList = append(argsList, agentId)
	}
	if serverIp != "" {
		condition += " and T3.`serverip` = ?"
		argsList = append(argsList, serverIp)
	}
	if forwardPort > 0 {
		condition += " and T0.`forwardport` = ?"
		argsList = append(argsList, forwardPort)
	}
	if honeyPort > 0 {
		condition += " and T0.`honeypotport` = ?"
		argsList = append(argsList, honeyPort)
	}
	if status != 0 {
		condition += " and T0.`status` = ?"
		argsList = append(argsList, status)
	}
	if createStartTime != "" && createEndTime != "" {
		condition += " and T0.`createtime` BETWEEN ? and ?"
		argsList = append(argsList, createStartTime)
		argsList = append(argsList, createEndTime)
	}
	if offlineStartTime != "" && offlineEndTime != "" {
		condition += " and T0.`offlinetime` BETWEEN ? and ?"
		argsList = append(argsList, offlineStartTime)
		argsList = append(argsList, offlineEndTime)
	}
	if creator != "" {
		condition += " and T0.`creator` = ?"
		argsList = append(argsList, creator)
	}
	if honeytypeid != "" {
		condition += " and T2.`typeid` = ?"
		argsList = append(argsList, honeytypeid)
	}
	var total int
	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, argsList...).Scan(&total)
	if err != nil {
		logs.Error("[SelectTransPolicy] select total error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode

	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " ORDER BY T0.`createtime` desc limit ?,?"
	argsList = append(argsList, offset)
	argsList = append(argsList, pagesize)
	sqlstr = sqlstr + condition
	rows, err := DbCon.Query(sqlstr, argsList...)
	if err != nil {
		logs.Error("[SelectTransPolicy] select list error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("[2] %s\v", err)
		logs.Error("[SelectTransPolicy] rows.Columns() error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	list, count, err := util.GetTransPolicyMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			log.Printf("[3] %s\v", err)
			logs.Error("[SelectTransPolicy] Unmarshal error,%s", err)
			return data, comm.DataUnmarshalError, comm.ErrorCode
		}
	} else {
		lists := []map[string]interface{}{}
		data = map[string]interface{}{"list": lists, "total": 0}
		return data, comm.DataSelectSuccess, comm.SuccessCode
	}

	return data, comm.DataSelectSuccess, comm.SuccessCode
}

// 查询所有存在代理转发的蜜罐服务器IP、对应的监听端口、蜜罐流量转发策略状态为1、蜜罐状态为1
func SelectHoneyPotsTransInfoByHoneyTypeId(typeId string) (map[string]interface{}, string, int) {
	var data map[string]interface{}
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectHoneyPotsTransInfoByHoneyTypeId] open mysql fail,%s", err1)
		return data, comm.DBConnectError, comm.ErrorCode
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqltotal := "select count(1) from `honeypotservers` T1 left join `honeyfowards` T2 on T2.`serverid` = T1.`serverid` left join `honeypots` T3 on T2.`honeypotid` = T2.`honeypotid` left join `honeypotstype` T0 on T3.`honeytypeid`= T0.`typeid` where T1.`status`=1 and T2.`status`=1"
	sqlstr := "select T1.`servername`,T1.`serverip`,T2.`forwardport`,T1.`serverid` from `honeypotservers` T1 left join `honeyfowards` T2 on T2.`serverid` = T1.`serverid` left join `honeypots` T3 on T2.`honeypotid` = T2.`honeypotid` left join `honeypotstype` T0 on T3.`honeytypeid`= T0.`typeid` where T1.`status`=1 and T2.`status`=1"
	var condition string
	var argsList []interface{}
	if typeId != "" {
		condition += " and T3.`honeytypeid`=?"
		argsList = append(argsList, typeId)
	} else {
		return data, "请选择服务类型", comm.ErrorCode
	}
	var total int
	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, argsList...).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneyPotsTransInfoByHoneyTypeId] select total error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode

	}
	sqlstr = sqlstr + condition
	fmt.Println(sqlstr)
	rows, err := DbCon.Query(sqlstr, argsList...)
	if err != nil {
		logs.Error("[SelectHoneyPotsTransInfoByHoneyTypeId] select list error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectHoneyPotsTransInfoByHoneyTypeId] rows.Columns() error,%s", err)
		return data, comm.DBSelectError, comm.ErrorCode
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	list, count, err := util.GetHoneyPotsInfoMysqlJson(rows, columns, total, values, scanArgs)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectHoneyPotsTransInfoByHoneyTypeId] Unmarshal error,%s", err)
			return data, comm.DataUnmarshalError, comm.ErrorCode
		}
	} else {
		lists := []map[string]interface{}{}
		data = map[string]interface{}{"list": lists, "total": 0}
		return data, comm.DataSelectSuccess, comm.SuccessCode
	}

	return data, comm.DataSelectSuccess, comm.SuccessCode
}

//func SelectHoneyPotsTransInfo(typeId string, serverId string) (map[string]interface{}, string, int) {
//	var data map[string]interface{}
//	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
//	if err1 != nil {
//		logs.Error("[SelectHoneyPotsTransInfo] open mysql fail,%s", err1)
//		return data, comm.DBConnectError, comm.ErrorCode
//	}
//	DbCon := sqlCon
//	defer sqlCon.Close()
//	sqltotal := "select count(1) from `honeypotservers` T1 inner join `honeyfowards` T2 on T2.`serverid` = T1.`serverid` inner join `honeypots` T3 on T2.`honeypotid` = T2.`honeypotid` inner join `honeypotstype` T0 on T3.`honeytypeid`=T0.`typeid` where T1.`status`=1 and T2.`status`=1"
//	sqlstr := "select T1.`servername`,T2.`forwardport`,T2.`taskid` from `honeypotservers` T1 inner join `honeyfowards` T2 on T2.`serverid` = T1.`serverid` inner join `honeypots` T3 on T2.`honeypotid` = T2.`honeypotid` inner join `honeypotstype` T0 on T3.`honeytypeid`= T0.`typeid` where T1.`status`=1 and T2.`status`=1 and T2.`forwardstatus` = 1"
//	var condition string
//	var argsList []interface{}
//	if typeId != "" && serverId != "" {
//		condition += " and T3.`honeytypeid` = ? and T2.`serverid`=?"
//		argsList = append(argsList, typeId)
//		argsList = append(argsList, serverId)
//	} else {
//		return data, "请选择集群和类型", comm.ErrorCode
//	}
//	var total int
//	sqltotal = sqltotal + condition
//	err := DbCon.QueryRow(sqltotal, argsList...).Scan(&total)
//	if err != nil {
//		logs.Error("[SelectHoneyPotsTransInfo] select total error,%s", err)
//		return data, comm.DBSelectError, comm.ErrorCode
//
//	}
//	sqlstr = sqlstr + condition
//	fmt.Println(sqlstr)
//	rows, err := DbCon.Query(sqlstr, argsList...)
//	if err != nil {
//		logs.Error("[SelectHoneyPotsTransInfo] select list error,%s", err)
//		return data, comm.DBSelectError, comm.ErrorCode
//	}
//	columns, err := rows.Columns()
//	if err != nil {
//		log.Printf("[2] %s\v", err)
//		logs.Error("[SelectHoneyPotsTransInfo] rows.Columns() error,%s", err)
//		return data, comm.DBSelectError, comm.ErrorCode
//	}
//
//	values := make([]sql.RawBytes, len(columns))
//	scanArgs := make([]interface{}, len(values))
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//	list, count, err := util.GetHoneyPotsInfoMysqlJson(rows, columns, total, values, scanArgs)
//	if count > 0 {
//		err = json.Unmarshal([]byte(list), &data)
//		if err != nil {
//			log.Printf("[3] %s\v", err)
//			logs.Error("[SelectHoneyPotsTransInfo] Unmarshal error,%s", err)
//			return data, comm.DataUnmarshalError, comm.ErrorCode
//		}
//	}
//
//	return data, comm.DataSelectSuccess, comm.SuccessCode
//}

func InsertAttackLog(proxytype string, srchost string, srcport int, attackip string, honeypotid string, honeypotport int, honeytypeid string, attacktime int64) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("INSERT INTO attacklog(srchost,srcport,attackip,honeypotid,honeypotport,honeytypeid,attacktime,country,province,proxytype) VALUES (?,?,?,?,?,?,?,'局域网','局域网',?)", srchost, srcport, attackip, honeypotid, honeypotport, honeytypeid, attacktime,proxytype).Values(&maps)
	if err != nil {
		logs.Error("[InsertAttackLog] insert config error,%s", err)
		msg = "日志数据插入失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode

}

// 查询攻击日志列表
func SelectAttackLogList(serverid string, honeytypeid string, srchost string, attackip string, startTime string, endTime string, pageSize int, pageNum int, honeyip string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectAttackLogList]open mysql fail %s", err1)
		msg = "数据库连接失败"
		return data, msg, comm.ErrorCode
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqltotal := "select count(1) from `attacklog` T0 left join `honeypotservers` T1 on T1.`serverid` = T0.`serverid` left join `honeypotstype` T2 on T2.`typeid` = T0.`honeytypeid` LEFT JOIN honeypots T3 ON T0.honeypotid= T3.honeypotid where 1=1 and T0.honeypotid != '' "
	sqlstr := "select T0.`id`,T0.`srchost`,T0.`attackip`,T0.`attacktime`,T1.`servername`,T1.`serverip`,T2.`honeypottype`,T0.`country`,T0.`province`,T3.honeyip from `attacklog` T0 left join `honeypotservers` T1 on T1.`serverid` = T0.`serverid` left join `honeypotstype` T2 on T2.`typeid` = T0.`honeytypeid` LEFT JOIN honeypots T3 ON T0.honeypotid= T3.honeypotid where 1=1 and T0.honeypotid != '' "

	var condition string
	var argsList []interface{}
	if honeytypeid != "" {
		condition += " and T2.`typeid` = ?"
		argsList = append(argsList, honeytypeid)
	}
	if serverid != "" {
		condition += " and T1.`serverid` = ?"
		argsList = append(argsList, serverid)
	}
	if srchost != "" {
		condition += " and T0.`srchost` = ?"
		argsList = append(argsList, srchost)
	}
	if attackip != "" {
		condition += " and T0.`attackip` = ?"
		argsList = append(argsList, attackip)
	}
	if honeyip != "" {
		condition += " and T3.`honeyip` = ?"
		argsList = append(argsList, honeyip)
	}
	if startTime != "" && endTime != "" {
		condition += " and T0.`attacktime` BETWEEN ? and ?"
		argsList = append(argsList, startTime)
		argsList = append(argsList, endTime)
	}

	var total int
	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, argsList...).Scan(&total)
	if err != nil {
		logs.Error("[SelectAttackLogList] select total error,%s", err)
		msg = "数据库查询失败"
		return data, msg, comm.ErrorCode
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " ORDER BY T0.`attacktime` desc limit ?,?"
	argsList = append(argsList, offset)
	argsList = append(argsList, pagesize)
	sqlstr = sqlstr + condition
	rows, err := DbCon.Query(sqlstr, argsList...)
	if err != nil {
		logs.Error("[SelectAttackLogList] select list error,%s", err)
		msg = "数据库查询失败"
		return data, msg, comm.ErrorCode
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectAttackLogList] rows.Columns() error,%s", err)
		msg = "数据库查询失败"
		return data, msg, comm.ErrorCode
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list, count, err := util.GetAttackLogListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectAttackLogList] Unmarshal list error,%s", err)
			msg = "数据格式转换失败"
			return data, msg, comm.ErrorCode
		}
	} else {
		lists := []map[string]interface{}{}
		data = map[string]interface{}{"list": lists, "total": 0}
		return data, comm.DataSelectSuccess, comm.SuccessCode
	}

	return data, msg, comm.SuccessCode

}

// 查询攻击日志详情
func SelectAttackLogDetail(id int, honeytypeid string, srchost string, attackip string, startTime string, endTime string, honeypotport int, honeyIP string, pageSize int, pageNum int, srcport int, eventdetail string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectAttackLogDetail]open mysql fail %s", err1)
		msg = "数据库连接失败"
		return data, msg, comm.ErrorCode
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqltotal := "select count(1) from `attacklog` T0 left join `honeypotservers` T1 on T1.`serverid` = T0.`serverid` left join `honeypotstype` T2 on T2.`typeid` = T0.`honeytypeid` where 1=1"
	sqlstr := "select T0.`srchost`,T0.`srcport`,T0.`attackip`,T0.`attacktime`,T0.`honeypotport`,T0.`logdata`,T0.`eventdetail`,T1.`serverip`,T2.`honeypottype`, T3.honeyip from `attacklog` T0 left join `honeypotservers` T1 on T1.`serverid` = T0.`serverid` left join `honeypotstype` T2 on T2.`typeid` = T0.`honeytypeid` LEFT JOIN honeypots T3 ON T0.honeypotid = T3.honeypotid where 1=1"

	var condition string
	var argsList []interface{}
	if id > 0 {
		condition += " and T0.`id` = ?"
		argsList = append(argsList, id)
	} else {
		return data, "请选择正确的事件", comm.ErrorCode
	}
	if honeytypeid != "" {
		condition += " and T2.`typeid` = ?"
		argsList = append(argsList, honeytypeid)
	}
	if srchost != "" {
		condition += " and T0.`srchost` = ?"
		argsList = append(argsList, srchost)
	}

	if attackip != "" {
		condition += " and T0.`attackip` = ?"
		argsList = append(argsList, attackip)
	}
	if honeyIP != "" {
		condition += " and T1.`serverip` = ?"
		argsList = append(argsList, honeyIP)
	}
	if honeypotport > 0 {
		condition += " and T0.`honeypotport` = ?"
		argsList = append(argsList, honeypotport)
	}

	if startTime != "" && endTime != "" {
		condition += " and T0.`attacktime` BETWEEN ? and ?"
		argsList = append(argsList, startTime)
		argsList = append(argsList, endTime)
	}
	if srcport > 0 {
		condition += " and T0.`srcport` = ?"
		argsList = append(argsList, srcport)
	}
	if eventdetail != "" {
		condition += " and T0.`eventdetail` = ?"
		argsList = append(argsList, eventdetail)
	}

	var total int
	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, argsList...).Scan(&total)
	if err != nil {
		logs.Error("[SelectAttackLogDetail] select total error,%s", err)
		msg = "数据库查询失败"
		return data, msg, comm.ErrorCode
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " ORDER BY T0.`attacktime` desc limit ?,?"
	argsList = append(argsList, offset)
	argsList = append(argsList, pagesize)
	sqlstr = sqlstr + condition
	rows, err := DbCon.Query(sqlstr, argsList...)
	if err != nil {
		logs.Error("[SelectAttackLogDetail] select list error,%s", err)
		msg = "数据库查询失败"
		return data, msg, comm.ErrorCode
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectAttackLogDetail] rows.Columns() error,%s", err)
		msg = "数据库查询失败"
		return data, msg, comm.ErrorCode
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list, count, err := util.GetAttackLogDetailMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)

	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectAttackLogDetail] Unmarshal list error,%s", err)
			msg = "数据库连接失败"
			return data, msg, comm.ErrorCode
		}
	} else {
		lists := []map[string]interface{}{}
		data = map[string]interface{}{"list": lists, "total": 0}
		return data, comm.DataSelectSuccess, comm.SuccessCode
	}

	return data, msg, comm.SuccessCode

}

// 增加config
func InsertConfig(confname string, confvalue string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into config (confname,confvalue) VALUES (?,?)", confname, confvalue).Values(&maps)
	if err != nil {
		logs.Error("[InsertConfig] insert config error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	}

	return data, msg, comm.SuccessCode
}

// 查询config
func SelectConfig(id int, confname string, pageSize int, pageNum int) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectConfig]open mysql fail %s", err1)
		msg = "数据库连接失败"
		return data, msg, comm.ErrorCode
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqltotal := "select count(1) from config where 1=1"
	sqlstr := "select confname,confvalue,id from config where 1=1"
	var condition string
	var argsList []interface{}
	if confname != "" {
		condition += " and confname = ?"
		argsList = append(argsList, confname)
	}
	if id > 0 {
		condition += " and id = ?"
		argsList = append(argsList, id)
	}

	var total int
	sqltotal = sqltotal + condition
	log.Println(sqltotal)
	err := DbCon.QueryRow(sqltotal, argsList...).Scan(&total)
	if err != nil {
		logs.Error("[SelectConfig] select total error,%s", err)
		msg = "数据库查询失败"
		return data, msg, comm.ErrorCode
	}
	//totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	//offset := (pagenum - 1) * pagesize
	//condition += " order by id desc limit ?,?"
	//argsList = append(argsList, offset)
	//argsList = append(argsList, pagesize)
	sqlstr = sqlstr + condition
	rows, err := DbCon.Query(sqlstr, argsList...)
	if err != nil {
		logs.Error("[SelectConfig] select list error,%s", err)
		msg = "数据库查询失败"
		return data, msg, comm.ErrorCode
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectConfig] rows.Columns() error,%s", err)
		msg = "数据库查询失败"
		return data, msg, comm.ErrorCode
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	list, count, err := util.GetConfigJson(rows, columns, total, values, scanArgs, 0, 0, 0)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectConfig] Unmarshal list error,%s", err)
			msg = "数据解析失败"
			return data, msg, comm.ErrorCode
		}
	} else {
		lists := []map[string]interface{}{}
		data = map[string]interface{}{"list": lists, "total": 0}
		return data, comm.DataSelectSuccess, comm.SuccessCode
	}

	return data, msg, comm.SuccessCode

}

// 删除config
func DeleteConfig(id int, confname string) (map[string]interface{}, string, int) {
	msg := "成功"
	o := orm.NewOrm()
	var maps []orm.Params
	var data map[string]interface{}
	_, err := o.Raw("DELETE FROM config where id=? and confname=?", id, confname).Values(&maps)
	if err != nil {
		logs.Error("[DeleteConfig] delete config error,%s", err)
		msg = "数据删除失败"
		return data, msg, comm.ErrorCode
	}

	return data, msg, comm.SuccessCode
}

// 查询未变更上线状态的诱饵策略
func SelectCreateBaitPolicyNotFinish() ([]util.BaitPolicyByStatus, error) {
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectCreateBaitPolicyNotFinish]open mysql fail %s", err1)
		return nil, err1
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "select taskid,agentid,status,createtime,data,md5 from server_bait where status = 0 or status is null"

	fmt.Println(sqlstr)
	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectCreateBaitPolicyNotFinish] select list error,%s", err)
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectCreateBaitPolicyNotFinish] rows.Columns() error,%s", err)
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list, err := util.GetBaitPolicyCheckCreateStatusMysqlJson(rows)

	return list, err
}

// 查询未变更下线状态的诱饵策略
func SelectOfflineBaitPolicyNotFinish() ([]util.BaitPolicyByStatus, error) {
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectOfflineBaitPolicyNotFinish]open mysql fail %s", err1)
		return nil, err1
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "select taskid,agentid,status,offlinetime,address,md5 from server_bait where status = 1 and offlinetime is not null"

	log.Println(sqlstr)
	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectOfflineBaitPolicyNotFinish] select list error,%s", err)
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectOfflineBaitPolicyNotFinish] rows.Columns() error,%s", err)
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list, err := util.GetBaitPolicyCheckOfflineStatusMysqlJson(rows)

	return list, err
}

// 查询未变更上线状态的透明转发策略
func SelectCreateTransPolicyNotFinish() ([]util.TransPolicyByStatus, error) {
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectCreateTransPolicyNotFinish]open mysql fail %s", err1)
		return nil, err1
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "select T0.`taskid`,T0.`agentid`,T0.`forwardport`,T0.`honeypotport`,T0.`createtime`,T0.`status`,T0.`type`,T0.`path`,T1.`serverip` as `honeyIP`,T2.`honeypottype` from `fowards` T0 inner join `honeypotservers` T1 on T1.`serverid` = T0.`serverid` inner join `honeypotstype` T2 on T2.`typeid` = T0.`honeytypeid` where (T0.`status` = 0 or T0.`status` is null) and T0.`createtime` is not null "

	// fmt.Println(sqlstr)
	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectCreateTransPolicyNotFinish] select list error,%s", err)
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectCreateTransPolicyNotFinish] rows.Columns() error,%s", err)
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list, err := util.GetTransPolicyCheckCreateStatusMysqlJson(rows)

	return list, err
}

// 查询未变更下线状态的透明转发策略
func SelectOfflineTransPolicyNotFinish() ([]util.TransPolicyByStatus, error) {
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectOfflineTransPolicyNotFinish]open mysql fail %s", err1)
		return nil, err1
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "select T0.`taskid`,T0.`agentid`,T0.`forwardport`,T0.`honeypotport`,T0.`offlinetime`,T0.`status`,T0.`type`,T0.`path`,T1.`serverip`,T2.`honeypottype` from `fowards` T0 inner join `honeypotservers` T1 on T1.`serverid` = T0.`serverid` inner join `honeypotstype` T2 on T2.`typeid` = T0.`honeytypeid` where T0.`offlinetime` is not null and T0.`status` = 1"

	// log.Println(sqlstr)
	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectOfflineTransPolicyNotFinish] select list error,%s", err)
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectOfflineTransPolicyNotFinish] rows.Columns() error,%s", err)
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list, err := util.GetTransPolicyCheckOfflineStatusMysqlJson(rows)

	return list, err
}

// 查询未变更上线状态的蜜罐流量转发策略
func SelectCreateHoneyTransPolicyNotFinish() ([]util.TransPolicyByStatus, error) {
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectCreateHoneyTransPolicyNotFinish]open mysql fail %s", err1)
		return nil, err1
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "select T0.`taskid`,T0.`agentid`,T0.`forwardport`,T0.`honeypotport`,T0.`createtime`,T0.`status`,T0.`type`,T0.`path`,T1.`honeyip`,T2.`honeypottype` from `honeyfowards` T0 inner join `honeypots` T1 on T1.`honeypotid` = T0.`honeypotid` inner join `honeypotstype` T2 on T2.`typeid` = T1.`honeytypeid` where T0.`status` = 0 or T0.`status` is null and T0.`createtime` is not null "

	// fmt.Println(sqlstr)
	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectCreateHoneyTransPolicyNotFinish] select list error,%s", err)
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectCreateHoneyTransPolicyNotFinish] rows.Columns() error,%s", err)
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list, err := util.GetTransPolicyCheckCreateStatusMysqlJson(rows)

	return list, err
}

// 查询未变更下线状态的蜜罐流量转发策略
func SelectOfflineHoneyTransPolicyNotFinish() ([]util.TransPolicyByStatus, error) {
	sqlCon, err1 := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err1 != nil {
		logs.Error("[SelectOfflineHoneyTransPolicyNotFinish]open mysql fail %s", err1)
		return nil, err1
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "select T0.`taskid`,T0.`agentid`,T0.`forwardport`,T0.`honeypotport`,T0.`offlinetime`,T0.`status`,T0.`type`,T0.`path`,T1.`honeyip`,T2.`honeypottype` from `honeyfowards` T0 inner join `honeypots` T1 on T1.`honeypotid` = T0.`honeypotid` inner join `honeypotstype` T2 on T2.`typeid` = T1.`honeytypeid` where T0.`offlinetime` is not null and T0.`status` = 1"

	// log.Println(sqlstr)
	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectOfflineHoneyTransPolicyNotFinish] select list error,%s", err)
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectOfflineHoneyTransPolicyNotFinish] rows.Columns() error,%s", err)
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list, err := util.GetTransPolicyCheckOfflineStatusMysqlJson(rows)

	return list, err
}
