package models

import (
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type TokenTraceLog struct {
	Id              int64  `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"` //蜜签跟踪攻击ID
	TokenTraceLogId string `gorm:"unique;not null; size:128" json:"TokenTraceLogId"`               //蜜签ID
	BaitName        string `gorm:"not null;size:128" json:"BaitName"`                              //蜜签跟踪码
	FileName        string `gorm:"not null;size:128" json:"FileName"`                              //蜜签跟踪码
	TokenType       string `gorm:"not null;size:128" json:"TokenType"`                             //蜜签跟踪码
	TraceCode       string `gorm:"not null;size:128" json:"TraceCode"`                             //蜜签跟踪码
	OpenIP          string `gorm:"not null;size:128" json:"OpenIP"`                                //攻击IP
	UserAgent       string `gorm:"not null;size:256" json:"UserAgent"`                             //攻击UA
	Location        string `gorm:"not null;size:128" json:"Location"`                              //攻击位置
	OpenTime        int64  `gorm:"not null;size:64" json:"OpenTime"`                               //攻击时间

}

func (tokenTraceLog *TokenTraceLog) CreateTokenTraceLog() error {
	if err := db.Create(tokenTraceLog).Error; err != nil {
		return err
	}
	return nil
}

func (tokenTraceLog *TokenTraceLog) GetTokenTraceLog(queryMap map[string]interface{}) (*[]TokenTraceLog, int64, error) {
	var ret []TokenTraceLog
	var total int64
	sql := fmt.Sprintf("SELECT * FROM token_trace_logs ")
	sqlTotal := fmt.Sprintf("SELECT count(1) FROM token_trace_logs ")
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
			conditionSql = fmt.Sprintf(" %s %s open_time > %s", conditionSql, condition, val)
		}
		if key == "EndTime" {
			conditionSql = fmt.Sprintf(" %s %s open_time < %s", conditionSql, condition, val)
		}
		if key == "Payload" {
			conditionSql = fmt.Sprintf(" %s %s concat(bait_name, file_name, open_ip) like '%s'", conditionSql, condition, val)
		}
	}

	pageSize := int(queryMap["PageSize"].(float64))

	pageNumber := int(queryMap["PageNumber"].(float64))

	t := fmt.Sprintf("order by open_time DESC limit %d offset %d ", pageSize, (pageNumber-1)*pageSize)
	sql = strings.Join([]string{sql, conditionSql, t}, " ")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, total, err
	}

	sqlTotal = strings.Join([]string{sqlTotal, conditionSql}, " ")
	if err := db.Raw(sqlTotal).Scan(&total).Error; err != nil {
		return nil, total, err
	}
	return &ret, total, nil
}
