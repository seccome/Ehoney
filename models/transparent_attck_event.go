package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
)

type TransparentEvent struct {
	Id                 int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`                       //透明代理攻击事件ID
	TransparentEventId string `gorm:"index;not null; size:64" json:"TransparentEventId" form:"TransparentEventId"` //透明代理攻击事件ID
	TransparentProxyId string `gorm:"index;not null; size:64" json:"TransparentProxyId" form:"TransparentProxyId"` //透明代理攻击事件ID
	AgentToken         string `json:"AgentToken" form:"AgentToken" gorm:"not null;size:128"`                       //代理Token
	AttackType         string `json:"AttackType" form:"AttackType" gorm:"not null;size:64"`                        //攻击类型
	AttackIp           string `json:"AttackIP" form:"AttackIP" gorm:"not null;size:256"`                           //攻击IP
	AttackPort         int32  `json:"AttackPort" form:"AttackPort" gorm:"index;not null;"`                         //攻击端口
	ProxyIp            string `json:"ProxyIP" form:"ProxyIP" gorm:"index;not null;size:64"`                        //代理IP
	ProxyPort          int32  `json:"ProxyPort" form:"ProxyPort" gorm:"index;not null;"`                           //代理端口
	DestIp             string `json:"DestIp" form:"DestIp" gorm:"not null;size:64"`                                //目标IP
	DestPort           int32  `json:"DestPort" form:"DestPort"`                                                    //目标端口
	OutPort            int32  `json:"OutPort" form:"OutPort"`
	AttackLocation     string `json:"AttackLocation" form:"AttackLocation" gorm:"not null;size:128"`
	CreateTime         int64  `json:"CreateTime" form:"CreateTime"` //创建时间
}

func (event *TransparentEvent) QueryAttack2AgentLines(attackIpParams, probeIpParams string) ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	if attackIpParams == "" || probeIpParams == "" {
		return ret, nil
	}
	sql := fmt.Sprintf("SELECT concat(attack_ip, \"-HACK\") AS Source, concat(proxy_ip, \"-EDGE\") AS Target,  \"RED\" Status FROM transparent_events WHERE attack_ip in (%s) AND proxy_ip IN (%s) GROUP BY attack_ip, proxy_ip", attackIpParams, probeIpParams)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (event *TransparentEvent) GetTransparentEventNodes() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode

	sql := fmt.Sprintf("SELECT concat(attack_ip, \"-HACK\") AS Id,  \"HACK\" NodeType, attack_ip AS Ip, attack_ip AS HostName FROM transparent_events WHERE create_time > %d GROUP BY attack_ip", util.GetCurrentIntTime()-(10*60))
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}

	return ret, nil
}

func (event *TransparentEvent) GetAttackedAgentNodes() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode

	sql := fmt.Sprintf("SELECT concat(a.agent_ip, \"-EDGE\") AS Id,  \"EDGE\" NodeType, a.agent_ip AS Ip, a.host_name AS HostName FROM transparent_events te left join agents a on te.agent_token = a.agent_token GROUP BY a.agent_ip ")

	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}

	return ret, nil
}

func (event *TransparentEvent) CreateEvent() error {
	event.TransparentEventId = util.GenerateId()
	event.CreateTime = util.GetCurrentIntTime()
	result := db.Create(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (event *TransparentEvent) GetAttackStatisticsByIP() (*[]comm.AttackStatistics, error) {
	var ret []comm.AttackStatistics
	sql := fmt.Sprintf("select attack_ip as Data, count(*) as Count from transparent_events group by attack_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
func (event *TransparentEvent) GetAttackStatisticsByProbeIP() (*[]comm.AttackStatistics, error) {
	var ret []comm.AttackStatistics
	sql := fmt.Sprintf("select proxy_ip as Data, count(1) as Count from transparent_events group by proxy_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
