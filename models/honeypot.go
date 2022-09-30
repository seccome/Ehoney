package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type Honeypot struct {
	Id           int64           `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"Id"`
	HoneypotId   string          `json:"HoneypotId" form:"HoneypotId"`
	HoneypotName string          `json:"HoneypotName" form:"HoneypotName" gorm:"not null;unique;size:128" binding:"required"` //蜜罐名称
	PodName      string          `json:"PodName" form:"PodName" gorm:"null;size:256"`
	ImageId      string          `json:"ImageId" form:"ImageId" gorm:"not null;size:256"`                              //镜像地址
	ImageAddress string          `json:"ImageAddress" form:"ImageAddress" gorm:"not null;size:256" binding:"required"` //镜像地址
	HoneypotIp   string          `json:"HoneypotIp" form:"HoneypotIp" gorm:"null;size:256"`
	ServerPort   int32           `json:"ServerPort"`
	ProtocolType string          `json:"ProtocolType"`
	CreateTime   int64           `json:"CreateTime" form:"CreateTime" gorm:"not null"`
	Status       comm.TaskStatus `json:"Status"`
}

func (honeypot *Honeypot) GetHoneypotNodes() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode
	sql := fmt.Sprintf("SELECT concat(honeypot_ip, \"-POD\") AS Id,  \"POD\" NodeType, honeypot_ip AS Ip, honeypot_name AS HostName FROM honeypots WHERE `status` = 3")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (honeypot *Honeypot) CreateHoneypot() error {
	honeypot.CreateTime = util.GetCurrentIntTime()
	honeypot.HoneypotId = util.GenerateId()
	result := db.Create(honeypot)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (honeypot *Honeypot) DeleteHoneypotByID(honeypotId string) error {
	if err := db.Where("honeypot_id= ?", honeypotId).Delete(&Honeypot{}).Error; err != nil {
		return err
	}
	return nil
}

func (honeypot *Honeypot) GetHoneypot(payload *comm.HoneypotSelectPayload) (*[]Honeypot, int64, error) {
	var ret []Honeypot
	var count int64

	if util.CheckInjectionData(payload.Payload) || util.CheckInjectionData(payload.ProtocolType) {
		return nil, 0, nil
	}

	var p = "%" + payload.Payload + "%"
	var sql string = ""
	if payload.ProtocolType == "" {
		sql = fmt.Sprintf("select * from honeypots where CONCAT(protocol_type, honeypot_name, honeypot_ip, create_time, status) LIKE '%s' order by create_time DESC", p)
	} else {
		sql = fmt.Sprintf("select * from honeypots where CONCAT(protocol_type, honeypot_name, honeypot_ip, create_time, status) LIKE '%s' and protocol_type = %s order by create_time DESC", p, payload.ProtocolType)
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

func (honeypot *Honeypot) GetHoneypotByID(honeypotId string) (*Honeypot, error) {

	var ret Honeypot
	if err := db.Where("honeypot_id = ?", honeypotId).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (honeypot *Honeypot) GetPodNameList() ([]string, error) {
	var ret []string
	if err := db.Select("pod_name").Find(honeypot).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (honeypot *Honeypot) GetHoneypotByAddress(ip string) (*Honeypot, error) {
	var ret Honeypot
	if err := db.Where("honeypot_ip = ?", ip).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (honeypot *Honeypot) GetHoneypotByName(name string) error {
	var ret Honeypot
	if err := db.Where("honeypot_name = ?", name).Take(&ret).Error; err != nil {
		return err
	}
	return nil
}

func (honeypot *Honeypot) GetHoneypotByPodName(name string) (*Honeypot, error) {
	var ret Honeypot
	if err := db.Where("pod_name = ?", name).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (honeypot *Honeypot) GetHoneypotCount() (int64, error) {
	var count int64
	if err := db.Model(honeypot).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (honeypot *Honeypot) UpdatePodInfoByPodName(podName string, status comm.TaskStatus, podIP string) error {
	db.Model(honeypot).Where("pod_name = ?", podName).Updates(Honeypot{Status: status, HoneypotIp: podIP})
	return nil
}

func (honeypot *Honeypot) GetProtocolProxyHoneypot() (*[]comm.ProtocolHoneypotSelectResultPayload, error) {
	var ret []comm.ProtocolHoneypotSelectResultPayload
	sql := fmt.Sprintf("select h2.id, concat(h.honeypot_name, '-', h.honeypot_ip, ':', h2.proxy_port) as ProtocolHoneypotIpPort from honeypots h, protocol_proxies h2 where h.id = h2.honeypot_id")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (honeypot *Honeypot) GetHoneypots() (*[]Honeypot, error) {
	var ret []Honeypot
	if err := db.Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
