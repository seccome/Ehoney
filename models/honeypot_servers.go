package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
)

type HoneypotServers struct {
	ID            int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`  //蜜网ID
	AgentID       string `json:"AgentID" form:"AgentID" gorm:"unique;not null;size:128"` //代理ID
	CreateTime    string `json:"CreateTime" form:"CreateTime" gorm:"not null"`           //注册时间
	HeartbeatTime string `json:"HeartbeatTime" form:"HeartbeatTime" gorm:"not null"`     //心跳时间
	ServerIP      string `json:"ServerIP" form:"ServerIP" gorm:"not null;size:256"`      //蜜网IP
	HostName      string `json:"HostName" form:"HostName" gorm:"not null;size:256"`      //主机名称
	Status        comm.TaskStatus    `json:"Status" form:"Status" gorm:"not null;size:256"`          //状态
}

func (server *HoneypotServers) CreateServer() error {
	if server.GetServerByAgentID(server.AgentID) != nil {
		db.Model(server).Where("agent_id = ?", server.AgentID).Updates(map[string]interface{}{"heartbeat_time": server.HeartbeatTime, "server_ip": server.ServerIP, "host_name": server.HostName, "Status": comm.SUCCESS})
	} else {
		result := db.Create(server)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (server *HoneypotServers) GetServerByAgentID(agentID string) *HoneypotServers {
	var ret HoneypotServers
	if err := db.Take(&ret, "agent_id = ?", agentID).Error; err != nil {
		return nil
	}
	return &ret
}

func (server *HoneypotServers) GetServerStatusByID(ID int64) (*HoneypotServers, error) {
	var ret HoneypotServers
	if err := db.Where("id = ? AND status = ?", ID, comm.SUCCESS).First(&ret).Error; err != nil {
		return &ret, err
	}
	return &ret, nil
}

func (server *HoneypotServers) GetServerByID(ID int64) (*HoneypotServers, error) {
	var ret HoneypotServers
	if err := db.Take(&ret, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (server *HoneypotServers) GetFirstHoneypotServer() (*HoneypotServers, error) {
	var honeypot HoneypotServers
	err := db.Model(server).Take(&honeypot).Error
	if err != nil {
		return nil, err
	}
	return &honeypot, nil
}

func (server *HoneypotServers) RefreshServerStatus() error {
	db.Model(server).Where("TIMESTAMPDIFF(second, heartbeat_time, ?) > ?", util.GetCurrentTime(),120).Updates(HoneypotServers{Status: comm.FAILED})
	db.Model(server).Where("TIMESTAMPDIFF(second, heartbeat_time, ?) < ?", util.GetCurrentTime(),120).Updates(HoneypotServers{Status: comm.SUCCESS})
	return nil
}

