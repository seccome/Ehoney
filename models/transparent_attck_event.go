package models

import (
	"decept-defense/controllers/comm"
	"fmt"
)

type TransparentEvent struct {
	ID                       int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`                    //透明代理攻击事件ID
	AttackType               string `json:"AttackType" form:"AttackType" gorm:"not null;size:256" binding:"required"` //攻击类型
	AgentID                  string `json:"AgentID" form:"AgentID" gorm:"not null;size:256"`                          //agentID
	AttackIP                 string `json:"AttackIP" form:"AttackIP" gorm:"not null;size:256"`                        //攻击IP
	AttackPort               int32  `json:"AttackPort" form:"AttackPort" gorm:"not null"`                             //攻击端口
	ProxyIP                  string `json:"ProxyIP" form:"ProxyIP" gorm:"not null;size:256"`                          //代理IP
	ProxyPort                int32  `json:"ProxyPort" form:"ProxyPort" gorm:"not null"`                               //代理端口
	Transparent2ProtocolPort int32  `json:"Transparent2ProtocolPort" form:"Transparent2ProtocolPort" gorm:"not null"` //透明代理转发到协议代理的内部端口
	DestIP                   string `json:"DestIP" form:"DestIP" gorm:"not null;size:256"`                            //目标IP
	DestPort                 int32  `json:"DestPort" form:"DestPort" gorm:"not null"`                                 //目标端口
	EventTime                string `json:"EventTime" form:"EventTime" gorm:"not null"`                               //创建时间
}

func (event *TransparentEvent) QueryAttack2ProbeLines(attackIpParams, probeIpParams string) ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	if attackIpParams == "" || probeIpParams == "" {
		return ret, nil
	}
	sql := fmt.Sprintf("SELECT  concat(attack_ip, \"-HACK\") AS Source, concat(proxy_ip, \"-EDGE\") AS Target,  \"RED\" Status FROM transparent_events WHERE attack_ip in (%s) AND proxy_ip IN  (%s) GROUP BY attack_ip, proxy_ip", attackIpParams, probeIpParams)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (event *TransparentEvent) GetTransparentEventNodes() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode

	sql := fmt.Sprintf("SELECT concat(attack_ip, \"-HACK\") AS Id,  \"HACK\" NodeType, attack_ip AS Ip, attack_ip AS HostName FROM transparent_events WHERE TIMESTAMPDIFF(HOUR, event_time, NOW()) < 6 GROUP BY attack_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}

	return ret, nil
}

func (event *TransparentEvent) CreateEvent() error {
	result := db.Create(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (event *TransparentEvent) GetAttackStatisticsByIP() (*[]comm.AttackStatistics, error) {
	var ret []comm.AttackStatistics
	sql := fmt.Sprintf("select attack_ip as Data, count(*) from transparent_events  group by attack_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
func (event *TransparentEvent) GetAttackStatisticsByProbeIP() (*[]comm.AttackStatistics, error) {
	var ret []comm.AttackStatistics
	sql := fmt.Sprintf("select proxy_ip as Data, count(*) from transparent_events  group by proxy_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
