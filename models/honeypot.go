package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type Honeypot struct {
	ID           int64           `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`
	HoneypotName string          `json:"HoneypotName" form:"HoneypotName" gorm:"not null;unique;size:128" binding:"required"` //蜜罐名称
	PodName      string          `json:"PodName" form:"PodName" gorm:"null;size:256"`
	ImageAddress string          `json:"ImageAddress" form:"ImageAddress" gorm:"not null;size:256" binding:"required"` //镜像地址
	ServersID    int64           `json:"ServersID" form:"ServersID" gorm:"not null"`
	Servers      HoneypotServers `gorm:"ForeignKey:ServersID"`
	CreateTime   string          `json:"CreateTime" form:"CreateTime" gorm:"not null"`
	HoneypotIP   string          `json:"HoneypotIP" form:"HoneypotIP" gorm:"null;size:256"`
	Creator      string          `json:"Creator" form:"Creator" gorm:"not null;size:256"`
	ServerPort   int32           `json:"ServerPort"`
	ServerType   string          `json:"ServerType"`
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
	result := db.Create(honeypot)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (honeypot *Honeypot) DeleteHoneypotByID(id int64) error {
	sql := fmt.Sprintf("DELETE FROM honeypots WHERE `id` = %d", id)
	return db.Exec(sql).Error
}

func (honeypot *Honeypot) GetHoneypot(payload *comm.HoneypotSelectPayload) (*[]comm.HoneypotSelectResultPayload, int64, error) {
	var ret []comm.HoneypotSelectResultPayload
	var count int64

	if util.CheckInjectionData(payload.Payload) || util.CheckInjectionData(payload.ProtocolType) {
		return nil, 0, nil
	}

	var p = "%" + payload.Payload + "%"
	var sql string = ""
	if payload.ProtocolType == "" {
		sql = fmt.Sprintf("select h.id, h.server_type, h.honeypot_name, h.honeypot_ip, h2.server_ip, h.create_time, h.status, h.creator from honeypots h, honeypot_servers h2 where h.servers_id = h2.id AND CONCAT(h.server_type, h.honeypot_name, h.honeypot_ip, h2.server_ip, h.create_time, h.status, h.creator) LIKE '%s' order by h.create_time DESC", p)
	} else {
		sql = fmt.Sprintf("select h.id, h.server_type, h.honeypot_name, h.honeypot_ip, h2.server_ip, h.create_time, h.status, h.creator from honeypots h, honeypot_servers h2 where h.servers_id = h2.id AND CONCAT(h.server_type, h.honeypot_name, h.honeypot_ip, h2.server_ip, h.create_time, h.status, h.creator) LIKE '%s' AND server_type = '%s' order by h.create_time DESC", p, payload.ProtocolType)
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

func (honeypot *Honeypot) GetHoneypotByID(id int64) (*Honeypot, error) {

	var ret Honeypot
	if err := db.Where("id = ? and status = ?", id, comm.SUCCESS).Take(&ret).Error; err != nil {
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

func (honeypot *Honeypot) RefreshServerStatusByPodName(podName string, status comm.TaskStatus, podIP string) error {
	db.Model(honeypot).Where("pod_name = ?", podName).Updates(Honeypot{Status: status, HoneypotIP: podIP})
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
