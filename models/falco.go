package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"encoding/json"
	"fmt"
	"strings"
)

type FalcoAttackEvent struct {
	ID           int64        `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`
	Output       string       `gorm:"not null"     json:"output"`
	Priority     string       `gorm:"not null"     json:"priority"`
	Rule         string       `gorm:"not null"     json:"rule"`
	Time         string       `gorm:"not null"     json:"time"`
	FileFlag     bool         `gorm:"not null; default: false"     json:"FileFlag"`
	DownloadPath string       `gorm:"null"         json:"DownloadPath"`
	OutputFields OutputFields `gorm:"embedded"     json:"output_fields"`
}

type OutputFields struct {
	ContainerID   string `gorm:"" json:"container.id"`               // 容器ID
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
	if err := db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (event *FalcoAttackEvent) GetFalcoEvent(payload comm.FalcoEventSelectPayload) (*[]comm.FalcoSelectResultPayload, int64, error) {
	var ret []comm.FalcoSelectResultPayload
	var count int64

	if util.CheckInjectionData(payload.Payload) || util.CheckInjectionData(payload.StartTime) || util.CheckInjectionData(payload.EndTime) {
		return nil, 0, nil
	}
	var p string = "%" + payload.Payload + "%"
	var sql = ""
	if payload.StartTime == "" && payload.EndTime == "" {
		sql = fmt.Sprintf("select h.id, h.pod_name as HoneypotName, rule as event, time, output, priority as level, file_flag, download_path from falco_attack_events h, honeypots h2  where h.pod_name = h2.pod_name AND TIMESTAMPDIFF(second, h2.create_time, h.time) > 60 AND CONCAT(h.pod_name, rule, output, priority) LIKE '%s' order by time DESC", p)
	} else {
		sql = fmt.Sprintf("select h.id, h.pod_name as HoneypotName, rule as event, time, output, priority as level, file_flag, download_path from falco_attack_events h, honeypots h2  where h.pod_name = h2.pod_name AND TIMESTAMPDIFF(second, h2.create_time, h.time) > 60 AND CONCAT(pod_name, rule, output, priority) LIKE '%s' AND time between '%s' and '%s' order by time DESC", p, payload.StartTime, payload.EndTime)
	}
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, 0, err
	}
	count = (int64)(len(ret))
	t := fmt.Sprintf("limit %d offset %d", payload.PageSize, (payload.PageNumber-1)*payload.PageSize)
	sql = strings.Join([]string{sql, t}, " ")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, 0, err
	}
	return &ret, count, nil
}

func (event *FalcoAttackEvent) GetFalcoEventForTraceSource(payload comm.AttackTraceSelectPayload) (*[]comm.TraceSourceResultPayload, error) {
	var ret []comm.TraceSourceResultPayload
	var result []comm.TraceSourceResultPayload
	if util.CheckInjectionData(payload.Payload) || util.CheckInjectionData(payload.StartTime) || util.CheckInjectionData(payload.EndTime) {
		return nil, nil
	}
	selectPayload := "%" + payload.Payload + "%"

	var sql = ""
	if payload.StartTime != "" && payload.EndTime != "" {
		sql = fmt.Sprintf("select h.id, time, output as Log, h.pod_name  from falco_attack_events h, honeypots h2  where h.pod_name = h2.pod_name AND TIMESTAMPDIFF(second, h2.create_time, h.time) > 3 AND CONCAT(output) LIKE '%s' AND time between '%s' and '%s' order by time DESC", selectPayload, payload.StartTime, payload.EndTime)
	} else {
		sql = fmt.Sprintf("select h.id, time, output as Log, h.pod_name  from falco_attack_events h, honeypots h2  where h.pod_name = h2.pod_name AND TIMESTAMPDIFF(second, h2.create_time, h.time) > 3 AND CONCAT(output) LIKE '%s' order by time DESC", selectPayload)
	}

	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}

	var falcoIDDetailMap map[int64]OutputFields
	falcoIDDetailMap = make(map[int64]OutputFields)
	var honeypotPodIpMap map[string]string
	honeypotPodIpMap = make(map[string]string)
	x, _ := (&Honeypot{}).GetHoneypots()
	y, _ := event.GetFalcoEvents()
	for _, data := range *y {
		falcoIDDetailMap[data.ID] = data.OutputFields
	}
	for _, data := range *x {
		honeypotPodIpMap[data.PodName] = data.HoneypotIP
	}

	for index, data := range ret {
		detail, _ := json.Marshal(falcoIDDetailMap[data.ID])
		ret[index].Detail = string(detail)
		ret[index].ProtocolType = "falco"
		if !strings.Contains(ret[index].ProtocolType, payload.ProtocolType) {
			continue
		}
		_, ok := honeypotPodIpMap[falcoIDDetailMap[data.ID].PodName]
		if ok {
			ret[index].HoneypotIP = honeypotPodIpMap[falcoIDDetailMap[data.ID].PodName]
		} else {
			ret[index].HoneypotIP = ""
		}
		if !strings.Contains(ret[index].HoneypotIP, payload.HoneypotIP) {
			continue
		}
		result = append(result, ret[index])
	}
	return &result, nil
}

func (event *FalcoAttackEvent) GetFalcoEventByID(id int64) (*FalcoAttackEvent, error) {
	var ret FalcoAttackEvent
	if err := db.Take(&ret, id).Error; err != nil {
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
