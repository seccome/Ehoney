package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"encoding/json"
	"fmt"
)

type AttackEvent struct {
	ID            int64  `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`
	CreateTime    string `gorm:"not null"`
	AttackIP      string `gorm:"unique" form:"AttackIP" json:"AttackIP" gorm:"unique;size:256" binding:"required"`
	HoneypotIP    string `gorm:"unique" form:"HoneypotIP" json:"HoneypotIP" gorm:"unique;size:256" binding:"required"`
	AttackAddress string `gorm:"unique" form:"AttackAddress" json:"AttackAddress" gorm:"unique;size:256" binding:"required"`
	ProtocolType  string `gorm:"unique" form:"ProtocolType" json:"ProtocolType" gorm:"unique;size:256"`
	ProbeIP       string `gorm:"unique" form:"ProbeIP" json:"ProbeIP" gorm:"unique;size:256" binding:"required"`
}

func (event *AttackEvent) CreateAttackEvent() error {
	if err := db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

type AttackSelectResultPayload struct {
	ID             int64  `json:"ID"`             //攻击日志ID
	AttackIP       string `json:"AttackIP"`       //攻击IP
	ProbeIP        string `json:"ProbeIP"`        //探针IP
	JumpIP         string `json:"JumpIP"`         //跳转IP
	HoneypotIP     string `json:"HoneypotIP"`     //蜜罐IP
	ProtocolType   string `json:"ProtocolType"`   //协议类型
	AttackTime     string `json:"AttackTime"`     //攻击时间
	AttackLocation string `json:"AttackLocation"` //攻击位置
}

func (event *AttackEvent) GetAttackEvent(payload comm.AttackEventSelectPayload) (*[]comm.AttackSelectResultPayload, error) {
	var ret []comm.AttackSelectResultPayload
	var ret1 []comm.AttackSelectResultPayload

	if util.CheckInjectionData(payload.ProtocolType) || util.CheckInjectionData(payload.AttackIP) || util.CheckInjectionData(payload.JumpIP) || util.CheckInjectionData(payload.ProbeIP) || util.CheckInjectionData(payload.HoneypotIP) || util.CheckInjectionData(payload.ProtocolType) {
		return nil, nil
	}

	var p = "%" + payload.Payload + "%"
	var attackIP = "%" + payload.AttackIP + "%"
	var jumpIP = "%" + payload.JumpIP + "%"
	var probeIP = "%" + payload.ProbeIP + "%"
	var honeypotIP = "%" + payload.HoneypotIP + "%"
	var protocolType = "%" + payload.ProtocolType + "%"

	sql := fmt.Sprintf("select h.attack_ip, h.proxy_ip as ProbeIP, h2.proxy_ip as JumpIP, h2.dest_ip as HoneypotIP, h2.protocol_type, h2.event_time as AttackTime, h2.attack_detail from transparent_events h, protocol_events h2 where LOCATE(h2.attack_ip, h.proxy_ip, 1) > 0 AND LOCATE(h.dest_ip, h2.proxy_ip, 1) > 0 AND h.transparent2_protocol_port = h2.attack_port AND CONCAT(h.attack_ip, h.proxy_ip, h2.proxy_ip, h2.dest_ip, h2.protocol_type, h2.event_time) LIKE '%s' and h.attack_ip LIKE '%s' and h.proxy_ip LIKE '%s' and  h2.proxy_ip LIKE '%s' and h2.dest_ip LIKE '%s' and h2.protocol_type LIKE '%s' order by h2.event_time DESC", p, attackIP, probeIP, jumpIP, honeypotIP, protocolType)
	if err := db.Raw(sql).Scan(&ret1).Error; err != nil {
		return nil, err
	}

	var ret2 []comm.AttackSelectResultPayload
	sql = fmt.Sprintf("select attack_ip, dest_ip as HoneypotIP, proxy_ip as JumpIP, protocol_type, event_time as AttackTime, attack_detail from  protocol_events  where  CONCAT(attack_ip, dest_ip, protocol_type, event_time) LIKE '%s' and attack_ip LIKE '%s' and proxy_ip LIKE '%s'  and dest_ip LIKE '%s' and protocol_type LIKE '%s' order by event_time DESC", p, attackIP, jumpIP, honeypotIP, protocolType)
	if err := db.Raw(sql).Scan(&ret2).Error; err != nil {
		return nil, err
	}
	for _, i := range ret1 {
		ret = append(ret, i)
	}

	for _, i := range ret2 {
		ret = append(ret, i)
	}

	for index, data := range ret {
		d, _ := util.GetLocationByIP(data.AttackIP)
		if d.City == "-" || d.Country_long == "-" {
			ret[index].AttackLocation = "LAN"
		} else {
			ret[index].AttackLocation = d.City + "-" + d.Country_long
		}

		ret[index].CounterInfo = findCounterMapByIP(ret[index].AttackIP)
	}
	return &ret, nil
}

func findCounterMapByIP(attackIP string) map[string]string {
	CounterEventMap := make(map[string]string)
	if attackIP == "" {
		return CounterEventMap
	}
	var counterEvents CounterEvent
	data, err := counterEvents.GetCounterEventsByAttackIp(attackIP)

	if data == nil || err != nil {
		return CounterEventMap
	}

	for _, d := range *data {
		m := StructToMapViaJson(d)
		for key, value := range m {
			if value != "" && value != "{}" {
				CounterEventMap[key] = value
			}
		}
	}
	return CounterEventMap
}

func StructToMapViaJson(data interface{}) map[string]string {
	m := make(map[string]string)
	j, _ := json.Marshal(data)
	json.Unmarshal(j, &m)
	return m
}

func (event *AttackEvent) GetAttackEventForSource(payload comm.AttackTraceSelectPayload) (*[]comm.TraceSourceResultPayload, error) {
	var result []comm.TraceSourceResultPayload

	var ret []comm.AttackSelectResultPayload
	var ret1 []comm.AttackSelectResultPayload
	var ret2 []comm.AttackSelectResultPayload
	if util.CheckInjectionData(payload.ProtocolType) || util.CheckInjectionData(payload.AttackIP) || util.CheckInjectionData(payload.HoneypotIP) || util.CheckInjectionData(payload.Payload) || util.CheckInjectionData(payload.StartTime) || util.CheckInjectionData(payload.EndTime) {
		return nil, nil
	}
	protocolType := "%" + payload.ProtocolType + "%"
	attackIP := "%" + payload.AttackIP + "%"
	honeypotIP := "%" + payload.HoneypotIP + "%"
	selectPayload := "%" + payload.Payload + "%"

	if payload.StartTime != "" && payload.EndTime != "" {
		sql := fmt.Sprintf("select h.attack_ip, h.proxy_ip as ProbeIP, h2.proxy_ip as JumpIP, h2.dest_ip as HoneypotIP, h2.protocol_type, h2.event_time as AttackTime, h2.attack_detail from transparent_events h, protocol_events h2 where LOCATE(h2.attack_ip, h.proxy_ip, 1) > 0 AND LOCATE(h.dest_ip, h2.proxy_ip, 1) > 0 AND h.transparent2_protocol_port = h2.attack_port AND h.attack_ip LIKE '%s' AND h2.dest_ip LIKE '%s' AND h2.protocol_type LIKE '%s' AND CONCAT(h.attack_ip, h.proxy_ip, h2.proxy_ip, h2.dest_ip, h2.protocol_type, h2.event_time, h2.attack_detail) LIKE '%s' AND h2.event_time betweent '%s' and '%s' order by h2.event_time DESC", attackIP, honeypotIP, protocolType, selectPayload, payload.StartTime, payload.EndTime)
		if err := db.Raw(sql).Scan(&ret1).Error; err != nil {
			return nil, err
		}

		sql = fmt.Sprintf("select attack_ip, dest_ip as HoneypotIP, protocol_type, event_time as AttackTime, attack_detail from  protocol_events  where attack_ip LIKE '%s' AND dest_ip LIKE '%s' AND protocol_type LIKE '%s' AND CONCAT(attack_ip, dest_ip, protocol_type, event_time, attack_detail) LIKE '%s' ANT event_time between '%s' and '%s' order by event_time DESC", attackIP, honeypotIP, protocolType, selectPayload, payload.StartTime, payload.EndTime)
		if err := db.Raw(sql).Scan(&ret2).Error; err != nil {
			return nil, err
		}
		for _, i := range ret1 {
			ret = append(ret, i)
		}

		for _, i := range ret2 {
			ret = append(ret, i)
		}
	} else {
		sql := fmt.Sprintf("select h.attack_ip, h.proxy_ip as ProbeIP, h2.proxy_ip as JumpIP, h2.dest_ip as HoneypotIP, h2.protocol_type, h2.event_time as AttackTime, h2.attack_detail from transparent_events h, protocol_events h2 where LOCATE(h2.attack_ip, h.proxy_ip, 1) > 0 AND LOCATE(h.dest_ip, h2.proxy_ip, 1) > 0 AND h.transparent2_protocol_port = h2.attack_port AND h.attack_ip LIKE '%s' AND h2.dest_ip LIKE '%s' AND h2.protocol_type LIKE '%s' AND CONCAT(h.attack_ip, h.proxy_ip, h2.proxy_ip, h2.dest_ip, h2.protocol_type, h2.event_time, h2.attack_detail) LIKE '%s' order by h2.event_time DESC", attackIP, honeypotIP, protocolType, selectPayload)
		if err := db.Raw(sql).Scan(&ret1).Error; err != nil {
			return nil, err
		}

		sql = fmt.Sprintf("select attack_ip, dest_ip as HoneypotIP, protocol_type, event_time as AttackTime, attack_detail from  protocol_events  where attack_ip LIKE '%s' AND dest_ip LIKE '%s' AND protocol_type LIKE '%s' AND CONCAT(attack_ip, dest_ip, protocol_type, event_time, attack_detail) LIKE '%s' order by event_time DESC", attackIP, honeypotIP, protocolType, selectPayload)
		if err := db.Raw(sql).Scan(&ret2).Error; err != nil {
			return nil, err
		}
		for _, i := range ret1 {
			ret = append(ret, i)
		}

		for _, i := range ret2 {
			ret = append(ret, i)
		}
	}

	for _, i := range ret {
		result = append(result, comm.TraceSourceResultPayload{Time: i.AttackTime, HoneypotIP: i.HoneypotIP, AttackIP: i.AttackIP, ProtocolType: i.ProtocolType, Detail: i.AttackDetail, Log: i.AttackDetail})
	}
	return &result, nil
}
