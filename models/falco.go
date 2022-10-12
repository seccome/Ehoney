package models

import (
	"decept-defense/pkg/util"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

/**





{"code":400,"msg":"请求参数错误","data":"json: cannot unmarshal string into Go struct field FalcoAttackEvent.time of type int64"}{"output":"23:50:54.719348051: Error wreating files in container (user=root user_loginuid=-1 command=sshd -D -e k8s.ns=default k8s.pod=ssssssh-6948f886dc-b9p44 container=027dae4c3693 file=/proc/self/oom_score_adj container_id=027dae4c3693 container_name=k8s_ssssssh_ssssssh-6948f886dc-b9p44_default_397f6b51-5cb1-47c5-870d-d3f60576e960_0 image=ehoney/ssh) k8s.ns=default k8s.pod=ssssssh-6948f886dc-b9p44 container=027dae4c3693","priority":"Error","rule":"Create files below container any dir","time":"2022-09-25T15:50:54.719348051Z", "output_fields": {"container.id":"027dae4c3693","container.image.repository":"ehoney/ssh","container.name":"k8s_ssssssh_ssssssh-6948f886dc-b9p44_default_397f6b51-5cb1-47c5-870d-d3f60576e960_0","evt.time":1664121054719348051,"fd.name":"/proc/self/oom_score_adj","k8s.ns.name":"default","k8s.pod.name":"ssssssh-6948f886dc-b9p44","proc.cmdline":"sshd -D -e","user.loginuid":-1,"user.name":"root"}}


*/

type FalcoAttackEvent struct {
	Id                 int64        `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`
	FalcoAttackEventId string       `gorm:"not null; size:32" json:"falcoAttackEventId"` // 容器ID
	Output             string       `gorm:"not null; size:1024"     json:"output"`
	Priority           string       `gorm:"not null; size:32"     json:"priority"`
	Rule               string       `gorm:"not null; size:256"     json:"rule"`
	Time               string       `gorm:"not null; size:32"     json:"time"`
	FileFlag           bool         `gorm:"not null; default: false"     json:"FileFlag"`
	HoneypotName       string       `gorm:"not null; size:64"     json:"honeypotName"`
	DownloadPath       string       `gorm:"size:128"         json:"downloadPath"`
	OutputFields       OutputFields `gorm:"embedded"     json:"output_fields"`
	CreateTime         int64        `form:"CreateTime" json:"createTime"`
}

type OutputFields struct {
	ContainerId   string `gorm:"" json:"container.id"`               // 容器ID
	Repository    string `gorm:"" json:"container.image.repository"` // 镜像仓库
	ContainerName string `gorm:"" json:"container.name"`             // 容器名称
	EventTime     int64  `gorm:"" json:"evt.time"`                   // 事件时间
	FilePath      string `gorm:"" json:"fd.name"`                    // 上传的文件名称
	Namespace     string `gorm:"" json:"k8s.ns.name"`                // 命名空间
	PodName       string `gorm:"" json:"k8s.pod.name"`               // pod名称
	Cmdline       string `gorm:"" json:"proc.cmdline"`               // 命令行
	LoginUID      int    `gorm:"" json:"user.loginuid"`              // 登录UID
	UserName      string `gorm:"" json:"user.name"`                  // 用户名
	ProcTTY       int    `gorm:"" json:"proc.tty"`                   // TTY
	ProcessPName  string `gorm:"" json:"proc.pname"`                 // 父进程名称
	ProcessName   string `gorm:"" json:"proc.name"`                  // 进程名称
	Connection    string `gorm:"" json:"connection"`                 // 进程名称
}

func (event *FalcoAttackEvent) CreateFalcoEvent() error {
	event.FalcoAttackEventId = util.GenerateId()
	event.CreateTime = util.GetCurrentIntTime()
	if err := db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (event *FalcoAttackEvent) GetFalcoEvent(queryMap map[string]interface{}) (*[]FalcoAttackEvent, int64, error) {
	var ret []FalcoAttackEvent
	var total int64
	sql := fmt.Sprintf("SELECT * FROM falco_attack_events ")
	sqlTotal := fmt.Sprintf("SELECT count(1) FROM falco_attack_events ")
	conditionFlag := false
	conditionSql := ""

	if queryMap["AttackIp"] != nil && queryMap["AttackIp"].(string) != "" && queryMap["Payload"] == "" {
		queryMap["Payload"] = queryMap["AttackIp"]
	}

	if queryMap["AgentIp"] != nil && queryMap["AgentIp"].(string) != "" && queryMap["Payload"] == "" {
		queryMap["Payload"] = queryMap["AgentIp"]
	}

	for key, val := range queryMap {

		if key == "PageSize" || key == "PageNumber" || key == "AttackIp" || key == "AgentIp" {
			continue
		}
		if val == "" {
			continue
		}
		if util.CheckInjectionData(val.(string)) {
			return nil, 0, nil
		}
		condition := "where"
		if !conditionFlag {
			conditionFlag = true
		} else {
			condition = "and"
		}
		if key == "StartTime" {
			conditionSql = fmt.Sprintf(" %s %s create_time > %s", conditionSql, condition, val)
		}
		if key == "EndTime" {
			conditionSql = fmt.Sprintf(" %s %s create_time < %s", conditionSql, condition, val)
		}
		// TODO 目前 attackIp 和 AgentIp 无法处理 当前处置是为防止干扰

		if key == "Payload" {
			conditionSql = fmt.Sprintf(" %s %s output like '%s' or cmdline like '%s'", conditionSql, condition, "%"+val.(string)+"%", "%"+val.(string)+"%")
		}
		if key == "HoneypotName" {
			conditionSql = fmt.Sprintf(" %s %s honeypot_name like '%s'", conditionSql, condition, "%"+val.(string)+"%")
		}
		if key == "ProtocolType" {
			val = strings.ReplaceAll(val.(string), "proxy", "")
			conditionSql = fmt.Sprintf(" %s %s repository like '%s'", conditionSql, condition, "%"+val.(string)+"%")
		}
	}
	pageSize := int(queryMap["PageSize"].(float64))
	pageNumber := int(queryMap["PageNumber"].(float64))
	t := fmt.Sprintf("order by create_time DESC limit %d offset %d ", pageSize, (pageNumber-1)*pageSize)
	sql = strings.Join([]string{sql, conditionSql, t}, " ")
	zap.L().Info(sql)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, total, err
	}

	sqlTotal = strings.Join([]string{sqlTotal, conditionSql}, " ")
	if err := db.Raw(sqlTotal).Scan(&total).Error; err != nil {
		return nil, total, err
	}
	return &ret, total, nil
}

func (event *FalcoAttackEvent) GetFalcoEventByID(falcoAttackEventId string) (*FalcoAttackEvent, error) {
	var ret FalcoAttackEvent
	if err := db.Where("falco_attack_event_id = ?", falcoAttackEventId).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (event *FalcoAttackEvent) GetFalcoEvents() (*[]FalcoAttackEvent, error) {
	var ret []FalcoAttackEvent
	if err := db.Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
