package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"fmt"
)

type ProtocolEvent struct {
	Id              int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`                 //透明代理攻击事件ID
	ProtocolEventId string `gorm:"index;not null; size:64" json:"ProtocolEventId" form:"ProtocolEventId"` //透明代理攻击事件ID
	ProtocolProxyId string `json:"ProtocolProxyId" form:"ProtocolProxyId" gorm:"index;not null"`          //协议id
	ProtocolType    string `json:"ProtocolType" form:"ProtocolType" gorm:"not null"`                      //协议类型
	AttackIp        string `json:"AttackIp" form:"AttackIp" gorm:"not null;size:256"`                     //攻击IP
	AttackPort      int32  `json:"AttackPort" form:"AttackPort" gorm:"not null"`                          //攻击端口
	ProxyIp         string `json:"ProxyIp" form:"ProxyIp" gorm:"not null;size:256"`                       //代理IP
	ProxyPort       int32  `json:"ProxyPort" form:"ProxyPort" gorm:"index;not null"`                      //代理端口
	DestIp          string `json:"DestIp" form:"DestIp" gorm:"not null;size:256"`                         //目标IP
	DestPort        int32  `json:"DestPort" form:"DestPort" gorm:"not null"`                              //目标端口
	AttackDetail    string `json:"AttackDetail" form:"AttackDetail" gorm:"not null"`                      //攻击详情
	EventTime       string `json:"EventTime" form:"EventTime" gorm:"not null"`                            //创建时间
	CreateTime      int64  `json:"CreateTime" form:"CreateTime"`                                          //创建时间
}

func (event *ProtocolEvent) QueryProtocol2PodRedLines() ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	sql := fmt.Sprintf("SELECT concat(\"%s-RELAY\") AS Source, concat(dest_ip, \"-POD\") AS Target,  \"RED\" Status FROM protocol_events WHERE create_time > %d GROUP BY dest_ip", configs.GetSetting().Server.AppHost, util.GetCurrentIntTime()-(10*60))
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (event *ProtocolEvent) QueryAttack2ProtocolLines(attackIpParams string) ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	if attackIpParams == "" {
		return ret, nil
	}
	sql := fmt.Sprintf("SELECT concat(attack_ip, \"-HACK\") AS Source, concat(proxy_ip, \"-RELAY\") AS Target,  \"RED\" Status FROM protocol_events WHERE attack_ip in (%s) AND create_time > %d GROUP BY attack_ip, proxy_ip", attackIpParams, util.GetCurrentIntTime()-(10*60))
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (event *ProtocolEvent) QueryAttackedAgent2ProtocolRedLines(agentIpParams string) ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	sql := fmt.Sprintf("SELECT concat(tp.agent_ip, \"-EDGE\") AS Source, concat(proxy_ip, \"-RELAY\") AS Target,  \"RED\" as Status FROM protocol_events pe left join transparent_proxies tp on pe.protocol_proxy_id = tp.protocol_proxy_id WHERE attack_ip in (%s) AND pe.create_time > %d and tp.agent_ip is not null GROUP BY Source", agentIpParams, util.GetCurrentIntTime()-(10*60))
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (event *ProtocolEvent) CreateEvent() error {
	event.ProtocolEventId = util.GenerateId()
	event.CreateTime = util.GetCurrentIntTime()
	result := db.Create(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (event *ProtocolEvent) GetAttackStatisticsByProtocol() (*[]comm.AttackStatistics, error) {
	var ret []comm.AttackStatistics
	sql := fmt.Sprintf("select protocol_type as Data, count(*) as Count from protocol_events where protocol_type is not null group by protocol_type")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
