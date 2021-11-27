package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type TokenTraceLog struct {
	ID        int64  `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"` //蜜签跟踪攻击ID
	OpenTime  string `gorm:"not null；size:256" json:"opentime"`                              //攻击时间
	TraceCode string `gorm:"not null；size:256" json:"tracecode"`                             //蜜签跟踪码
	OpenIP    string `gorm:"not null；size:256" json:"OpenIP"`                                //攻击IP
	UserAgent string `gorm:"not null；size:256" json:"useragent"`                             //攻击UA
	Location  string `gorm:"not null；size:256" json:"Location"`                              //攻击位置
}

func (token *TokenTraceLog) GetTokenTraceLog(payload comm.TokenTraceSelectPayload) (*[]comm.TokenTraceSelectResultPayload, int64, error) {
	var ret []comm.TokenTraceSelectResultPayload
	var honeypotRet []comm.TokenTraceSelectResultPayload
	var probeRet []comm.TokenTraceSelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) || util.CheckInjectionData(payload.AttackIP) {
		return nil, 0, nil
	}
	var p = "%" + payload.Payload + "%"
	var attackIP = "%" + payload.AttackIP + "%"

	var sql = fmt.Sprintf("select h2.id, h2.open_time, h2.open_ip, h2.user_agent, h2.location, h.token_name, h.token_type from honeypot_tokens h, token_trace_logs h2  where h.trace_code = h2.trace_code  AND CONCAT(h2.open_time, h2.open_ip, h2.user_agent, h.token_name) LIKE '%s' and h2.open_ip LIKE '%s'", p, attackIP)
	if payload.StartTime != "" && payload.EndTime != "" {
		sql = strings.Join([]string{sql, fmt.Sprintf("AND h2.open_time between '%s' and '%s'", payload.StartTime, payload.EndTime)}, " ")
	}
	sql = strings.Join([]string{sql, "order by open_time DESC"}, " ")
	if err := db.Raw(sql).Scan(&honeypotRet).Error; err != nil {
		return nil, 0, err
	}

	sql = fmt.Sprintf("select h2.id, h2.open_time, h2.open_ip, h2.user_agent, h2.location, h.token_name, h.token_type from probe_tokens h, token_trace_logs h2  where h.trace_code = h2.trace_code  AND CONCAT(h2.open_time, h2.open_ip, h2.user_agent, h.token_name) LIKE '%s' and h2.open_ip LIKE '%s'", p, attackIP)
	if payload.StartTime != "" && payload.EndTime != "" {
		sql = strings.Join([]string{sql, fmt.Sprintf("AND h2.open_time between '%s' and '%s'", payload.StartTime, payload.EndTime)}, " ")
	}
	sql = strings.Join([]string{sql, "order by open_time DESC"}, " ")
	if err := db.Raw(sql).Scan(&probeRet).Error; err != nil {
		return nil, 0, err
	}

	if payload.ServerType == "honeypot" {
		for _, i := range honeypotRet {
			ret = append(ret, i)
		}
	} else if payload.ServerType == "probe" {
		for _, i := range probeRet {
			ret = append(ret, i)
		}
	} else {
		for _, i := range honeypotRet {
			ret = append(ret, i)
		}
		for _, i := range probeRet {
			ret = append(ret, i)
		}
	}

	for index, data := range ret {
		d, _ := util.GetLocationByIP(data.OpenIP)
		if d.City == "-" || d.Country_long == "-" {
			ret[index].Location = "LAN"
		} else {
			ret[index].Location = d.City + "-" + d.Country_long
		}
	}
	count = int64(len(ret))

	var start int = payload.PageSize * (payload.PageNumber - 1)
	var end int = payload.PageSize*(payload.PageNumber-1) + payload.PageSize
	if payload.PageSize*(payload.PageNumber-1) > len(ret) {
		start = len(ret)
	}
	if payload.PageSize*(payload.PageNumber-1)+payload.PageSize > len(ret) {
		end = len(ret)
	}
	ret = ret[start:end]
	return &ret, count, nil
}
