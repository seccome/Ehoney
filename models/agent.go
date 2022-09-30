package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type Agent struct {
	Id            int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`         //探针ID
	AgentId       string `json:"AgentId" form:"AgentId" gorm:"unique;not null;size:128"`        //代理ID
	AgentToken    string `json:"AgentToken" form:"AgentToken" gorm:"not null;size:256"`         //代理Token
	AgentIp       string `json:"AgentIp" form:"AgentIp" gorm:"not null;size:256"`               //服务IP
	SystemType    string `json:"SystemType" form:"SystemType" gorm:"not null;size:256"`         //系统类型
	SubnetMask    string `json:"SubnetMask" form:"SubnetMask" gorm:"not null;size:256"`         //子网掩码
	HostName      string `json:"HostName" form:"HostName" gorm:"not null;size:256"`             //探针名称
	Status        int    `json:"Status" form:"Status" gorm:"not null;size:3;default:3"`         //状态
	LostStatus    int    `json:"LostStatus" form:"LostStatus" gorm:"not null;size:3;default:3"` //失陷状态
	CreateTime    int64  `json:"CreateTime" form:"CreateTime" gorm:"not null"`                  //注册时间
	HeartbeatTime int64  `json:"HeartbeatTime" form:"HeartbeatTime" gorm:"not null"`            //心跳时间
}

func (agent *Agent) CreateAgent() error {
	result := db.Create(agent)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (agent *Agent) UpdateAgentHeartBeat() error {
	db.Model(agent).Where("agent_token = ?", agent.AgentToken).
		Update("heartbeat_time", agent.HeartbeatTime).
		Update("status", agent.Status)
	return nil
}

func (agent *Agent) QueryAgentPage(payload *comm.SelectPayload) (*[]Agent, int64, error) {
	var agents []Agent
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}
	p := "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select agent_id, agent_token, agent_ip, system_type, subnet_mask, host_name, lost_status, create_time, heartbeat_time, if(heartbeat_time < %d,2,3) as status from agents where CONCAT(agent_ip, host_name, create_time, heartbeat_time, system_type) LIKE '%s' order by create_time DESC", util.GetCurrentIntTime()-(3*60), p)
	if err := db.Raw(sql).Scan(&agents).Error; err != nil {
		return nil, 0, err
	}
	count = (int64)(len(agents))
	t := fmt.Sprintf("limit %d offset %d", payload.PageSize, (payload.PageNumber-1)*payload.PageSize)
	sql = strings.Join([]string{sql, t}, " ")
	fmt.Printf(sql)
	if err := db.Raw(sql).Scan(&agents).Error; err != nil {
		return nil, 0, err
	}
	for _, a := range agents {
		if a.HeartbeatTime < util.GetCurrentIntTime()-(3*60) {
			a.Status = 2
		}
	}
	return &agents, count, nil
}

func (agent *Agent) GetAgentByAgentId(agentId string) *Agent {
	var ret Agent
	if err := db.Take(&ret, "agent_id = ?", agentId).Error; err != nil {
		return nil
	}
	return &ret
}

func (agent *Agent) GetAgentByAgentToken(agentToken string) *Agent {
	var ret Agent
	if err := db.Take(&ret, "agent_token = ?", agentToken).Error; err != nil {
		return nil
	}
	return &ret
}

func (agent *Agent) RefreshAgentStatus() error {
	db.Model(agent).Where("heartbeat_time > ?", util.GetCurrentIntTime(), 120).Updates(Agent{Status: 2})
	db.Model(agent).Where("heartbeat_time < ?", util.GetCurrentIntTime(), 120).Updates(Agent{Status: 3})
	return nil
}
