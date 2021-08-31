package models

import (
	"decept-defense/controllers/comm"
	"fmt"
)

type ProtocolEvent struct {
	ID           int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`                    //透明代理攻击事件ID
	AttackType   string `json:"AttackType" form:"AttackType" gorm:"not null;size:256" binding:"required"` //攻击类型
	AgentID      string `json:"AgentID" form:"AgentID" gorm:"not null;size:256"`                          //agentID
	AttackIP     string `json:"AttackIP" form:"AttackIP" gorm:"not null;size:256"`                        //攻击IP
	AttackPort   int32  `json:"AttackPort" form:"AttackPort" gorm:"not null"`                             //攻击端口
	ProxyIP      string `json:"ProxyIP" form:"ProxyIP" gorm:"not null;size:256"`                          //代理IP
	ProxyPort    int32  `json:"ProxyPort" form:"ProxyPort" gorm:"not null"`                               //代理端口
	ProtocolType string `json:"ProtocolType" form:"ProtocolType" gorm:"not null"`                         //协议类型
	DestIP       string `json:"DestIP" form:"DestIP" gorm:"not null;size:256"`                            //目标IP
	DestPort     int32  `json:"DestPort" form:"DestPort" gorm:"not null"`                                 //目标端口
	AttackDetail string `json:"AttackDetail" form:"AttackDetail" gorm:"not null"`                         //攻击详情
	EventTime    string `json:"EventTime" form:"EventTime" gorm:"not null"`                               //创建时间
}

func (event *ProtocolEvent) QueryProtocol2PodRedLines() ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	sql := fmt.Sprintf("SELECT concat(hs.server_ip, \"-RELAY\") AS Source, concat(pe.dest_ip, \"-POD\") AS Target,  \"RED\" Status FROM protocol_events pe INNER JOIN honeypot_servers hs ON pe.agent_id = hs.agent_id WHERE pe.dest_ip != \"\" AND hs.`status` = 3 AND TIMESTAMPDIFF(HOUR, pe.event_time, NOW()) < 6 GROUP BY hs.server_ip, pe.dest_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (event *ProtocolEvent) QueryAttack2ProbeLines(attackIpParams, probeIpParams string) ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	if attackIpParams == "" || probeIpParams == "" {
		return ret, nil
	}
	sql := fmt.Sprintf("SELECT concat(attack_ip, \"-HACK\") AS Source, concat(proxy_ip, \"-RELAY\") AS Target,  \"RED\" Status FROM protocol_events pe WHERE attack_ip in (%s) AND proxy_ip IN (%s) AND TIMESTAMPDIFF(HOUR, pe.event_time, NOW()) < 6 GROUP BY attack_ip, proxy_ip", attackIpParams, probeIpParams)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (event *ProtocolEvent) QueryTransparent2ProbeRedLines() ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	sql := fmt.Sprintf("SELECT concat(pe.attack_ip, \"-EDGE\") AS Source, concat(hs.server_ip, \"-RELAY\") AS Target,  \"RED\" Status  FROM protocol_events pe INNER JOIN honeypot_servers hs ON pe.agent_id = hs.agent_id WHERE TIMESTAMPDIFF(HOUR, pe.event_time, NOW()) < 6 GROUP BY  pe.attack_ip, hs.server_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (event *ProtocolEvent) CreateEvent() error {
	result := db.Create(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (event *ProtocolEvent) Query() error {
	result := db.Create(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (event *ProtocolEvent) GetAttackStatisticsByProtocol() (*[]comm.AttackStatistics, error) {
	var ret []comm.AttackStatistics
	sql := fmt.Sprintf("select protocol_type as Data, count(*) from protocol_events group by protocol_type")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
