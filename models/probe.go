package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type Probes struct {
	ID            int64           `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`   //探针ID
	AgentID       string          `json:"AgentID" form:"AgentID" gorm:"unique;not null;size:128"`  //代理ID
	CreateTime    string          `json:"CreateTime" form:"CreateTime" gorm:"not null"`            //注册时间
	HeartbeatTime string          `json:"HeartbeatTime" form:"HeartbeatTime" gorm:"not null"`      //心跳时间
	ServerIP      string          `json:"ServerIP" form:"ServerIP" gorm:"not null;size:256"`       //服务IP
	SystemType    string          `json:"SystemType" form:"SystemType" gorm:"not null;size:256"`   //系统类型
	HostName      string          `json:"HostName" form:"HostName" gorm:"not null;size:256"`       //探针名称
	Status        comm.TaskStatus `json:"Status" form:"Status" gorm:"not null;size:256;default:3"` //状态
}

func (server *Probes) CreateServer() error {
	if server.GetServerByAgentID(server.AgentID) != nil {
		db.Model(server).Where("agent_id = ?", server.AgentID).Updates(map[string]interface{}{"HeartbeatTime": server.HeartbeatTime, "ServerIP": server.ServerIP, "HostName": server.HostName, "Status": comm.SUCCESS})
	} else {
		result := db.Create(server)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (server *Probes) GetProbe(payload *comm.SelectPayload) (*[]comm.ProbeSelectResultPayload, int64, error) {
	var ret []comm.ProbeSelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}

	var p string = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select id, server_ip as ProbeIP, host_name, create_time, heartbeat_time, status, system_type from probes where CONCAT(server_ip, host_name, create_time, heartbeat_time, system_type) LIKE '%s' order by create_time DESC", p)
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

func (server *Probes) GetServerStatusByAgentID(agentID string) *Probes {
	var ret Probes
	if err := db.Where("agent_id = ? AND status = ?", agentID, comm.SUCCESS).Find(server).Scan(&ret).Error; err != nil {
		return nil
	}
	return &ret
}

func (server *Probes) GetServerByAgentID(agentID string) *Probes {
	var ret Probes
	if err := db.Take(&ret, "agent_id = ?", agentID).Error; err != nil {
		return nil
	}
	return &ret
}

func (server *Probes) GetServerStatusByID(ID int64) (*Probes, error) {
	var ret Probes
	if err := db.Take(&ret, "id = ? AND status = ?", ID, comm.SUCCESS).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (server *Probes) GetServerStatusByID2(ID int64) (*Probes, error) {
	var ret Probes
	if err := db.Take(&ret, "id = ?", ID, comm.SUCCESS).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (server *Probes) GetServerByID(ID int64) (*Probes, error) {
	var ret Probes
	if err := db.Where("id = ?", ID).Find(server).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (server *Probes) RefreshServerStatus() error {

	db.Model(server).Where("TIMESTAMPDIFF(second, heartbeat_time, ?) > ?", util.GetCurrentTime(), 120).Updates(Probes{Status: comm.FAILED})
	db.Model(server).Where("TIMESTAMPDIFF(second, heartbeat_time, ?) < ?", util.GetCurrentTime(), 120).Updates(Probes{Status: comm.SUCCESS})
	return nil
}
