package models

import (
	"decept-defense/pkg/util"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type AttackEvent struct {
	TransparentEventId string `form:"TransparentEventId" json:"TransparentEventId"`
	ProtocolEventId    string `form:"ProtocolEventId" json:"ProtocolEventId"`
	AttackIp           string `form:"AttackIp" json:"AttackIp"`
	AttackLocation     string `form:"AttackLocation" json:"AttackLocation"`
	AgentIp            string `form:"AgentIp" json:"AgentIp"`
	AgentPort          int32  `json:"AgentPort" form:"AgentPort"` //代理端口
	ProtocolType       string `form:"ProtocolType" json:"ProtocolType"`
	ProtocolPort       int32  `json:"ProtocolPort" form:"ProtocolPort"` //协议端口
	HoneypotIp         string `form:"HoneypotIp" json:"HoneypotIp"`
	HoneypotPort       int32  `form:"HoneypotPort" json:"HoneypotPort"` //代理端口
	AttackDetail       string `form:"AttackDetail" json:"AttackDetail"` //代理端口
	CreateTime         int64  `form:"CreateTime" json:"CreateTime"`
}

type AttackSelectResultPayload struct {
	AttackIp       string `json:"AttackIp"`       //攻击IP
	AgentIp        string `json:"AgentIp"`        //探针IP
	HoneypotIp     string `json:"HoneypotIp"`     //蜜罐IP
	ProtocolType   string `json:"ProtocolType"`   //协议类型
	AttackTime     string `json:"AttackTime"`     //攻击时间
	AttackLocation string `json:"AttackLocation"` //攻击位置
}

func (event *AttackEvent) GetAttackEvent(queryMap map[string]interface{}) (*[]AttackEvent, int64, error) {
	var ret []AttackEvent
	var total int64
	sql := fmt.Sprintf("SELECT te.transparent_event_id as TransparentEventId, pe.protocol_event_id as ProtocolEventId, ifnull(te.create_time, pe.create_time) As CreateTime,  ifnull(te.attack_ip, pe.attack_ip) as AttackIp, te.proxy_ip as AgentIp, te.proxy_port as AgentPort, pe.proxy_port as ProtocolPort, pe.protocol_type as ProtocolType,  pe.attack_detail as AttackDetail, pe.dest_ip as HoneypotIp, pe.dest_port as HoneypotPort FROM protocol_events pe left join transparent_events te on pe.attack_port = te.out_port left join protocol_proxies pp on pe.protocol_proxy_id = pp.protocol_proxy_id left join honeypots h on pp.honeypot_id = h.honeypot_id ")
	sqlTotal := fmt.Sprintf("SELECT count(1) FROM protocol_events pe left join transparent_events te on pe.attack_port = te.out_port left join protocol_proxies pp on pe.protocol_proxy_id = pp.protocol_proxy_id left join honeypots h on pp.honeypot_id = h.honeypot_id ")

	conditionFlag := false
	conditionSql := ""
	for key, val := range queryMap {

		if key == "PageSize" || key == "PageNumber" {
			continue
		}
		if val == "" {
			continue
		}
		if util.CheckInjectionData(val.(string)) {
			return nil, 0, nil
		}
		condition := "where"
		if !conditionFlag {
			conditionFlag = true
		} else {
			condition = "and"
		}
		if key == "StartTime" {
			conditionSql = fmt.Sprintf(" %s %s te.create_time > %s or pe.create_time > %s", conditionSql, condition, val, val)
		}
		if key == "EndTime" {
			conditionSql = fmt.Sprintf(" %s %s te.create_time > %s or pe.create_time > %s", conditionSql, condition, val, val)
		}
		if key == "AttackIp" {
			conditionSql = fmt.Sprintf(" %s %s te.attack_ip = '%s' or pe.attack_ip = '%s' ", conditionSql, condition, val, val)
		}
		if key == "AgentIp" {
			conditionSql = fmt.Sprintf(" %s %s te.proxy_ip like '%s'", conditionSql, condition, "%"+val.(string)+"%")
		}
		if key == "HoneypotName" {
			conditionSql = fmt.Sprintf(" %s %s h.honeypot_name like '%s'", conditionSql, condition, "%"+val.(string)+"%")
		}
		if key == "Payload" {
			conditionSql = fmt.Sprintf(" %s %s pe.attack_detail like '%s'", conditionSql, condition, "%"+val.(string)+"%")
		}
		if key == "ProtocolType" {
			val = strings.ReplaceAll(val.(string), "proxy", "")
			conditionSql = fmt.Sprintf(" %s %s pe.protocol_type like '%s'", conditionSql, condition, "%"+val.(string)+"%")
		}
	}

	pageSize := int(queryMap["PageSize"].(float64))

	pageNumber := int(queryMap["PageNumber"].(float64))

	t := fmt.Sprintf("order by CreateTime DESC limit %d, %d ", (pageNumber-1)*pageSize, pageSize)
	sql = strings.Join([]string{sql, conditionSql, t}, " ")
	zap.L().Info(sql)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, total, err
	}

	sqlTotal = strings.Join([]string{sqlTotal, conditionSql}, " ")
	zap.L().Info(sqlTotal)
	if err := db.Raw(sqlTotal).Scan(&total).Error; err != nil {
		return nil, total, err
	}

	for index, data := range ret {
		d, _ := util.GetLocationByIP(data.AttackIp)
		if d.City == "-" || d.Country_long == "-" {
			ret[index].AttackLocation = "LAN"
		} else {
			ret[index].AttackLocation = d.City + "-" + d.Country_long
		}
	}
	return &ret, total, nil
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
