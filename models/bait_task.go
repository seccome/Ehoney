package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type BaitTask struct {
	Id                int64           `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`
	BaitTaskId        string          `gorm:"not null;size:64" json:"BaitTaskId"`
	BaitId            string          `gorm:"not null;size:64" json:"BaitId"`
	BaitName          string          `gorm:"not null;size:64" json:"BaitName"`
	BaitType          string          `gorm:"not null;size:32" json:"BaitType"`
	OperatorType      string          `gorm:"not null;size:32" json:"OperatorType"`
	TokenId           string          `gorm:"not null;size:32" json:"TokenId"`
	TokenType         string          `gorm:"not null;size:32" json:"TokenType"`
	URL               string          `gorm:"not null;size:256" json:"URL"`
	FileMD5           string          `gorm:"not null;size:64" json:"FileMD5"`
	LocalPath         string          `gorm:"not null;size:256" json:"LocalPath"`
	DeployPath        string          `gorm:"not null;size:256" json:"DeployPath"`
	ScriptName        string          `gorm:"not null;size:256" json:"ScriptName"`
	CommandParameters string          `gorm:"not null;size:256" json:"CommandParameters"`
	BaitData          string          `gorm:"not null;size:2048" json:"BaitData"`
	HoneypotId        string          `gorm:"not null" json:"HoneypotId"`
	AgentToken        string          `gorm:"not null" json:"AgentToken"`
	Status            comm.TaskStatus `gorm:"not null" json:"Status"`
	CreateTime        int64           `gorm:"not null" json:"CreateTime"`
}

func (baitTask *BaitTask) CreateBaitTask() error {
	if err := db.Create(baitTask).Error; err != nil {
		return err
	}
	return nil
}

func (baitTask *BaitTask) GetBaitTaskById(baitTaskId string) (*BaitTask, error) {
	var ret BaitTask
	if err := db.Where("bait_task_id = ?", baitTaskId).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (baitTask *BaitTask) QueryBaitTaskPage(queryMap map[string]interface{}) (*[]BaitTask, int64, error) {
	var ret []BaitTask
	var total int64
	sql := fmt.Sprintf("select * from bait_tasks ")
	sqlTotal := fmt.Sprintf("SELECT count(1) FROM bait_tasks ")

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
		if key == "BaitType" {
			conditionSql = fmt.Sprintf(" %s %s bait_type = '%s'", conditionSql, condition, val)
		}
		if key == "HoneypotId" {
			conditionSql = fmt.Sprintf(" %s %s honeypot_id = '%s'", conditionSql, condition, val)
		}
		if key == "AgentToken" {
			conditionSql = fmt.Sprintf(" %s %s agent_token = '%s'", conditionSql, condition, val)
		}
		if key == "Payload" {
			conditionSql = fmt.Sprintf(" %s %s bait_name like '%s'", conditionSql, condition, "%"+val.(string)+"%")
		}
	}

	pageSize := int(queryMap["PageSize"].(float64))

	pageNumber := int(queryMap["PageNumber"].(float64))

	pageSql := fmt.Sprintf("order by create_time desc limit %d offset %d", pageSize, (pageNumber-1)*pageSize)
	zap.L().Info(sql)

	sqlPage := strings.Join([]string{sql, conditionSql, pageSql}, " ")
	if err := db.Raw(sqlPage).Scan(&ret).Error; err != nil {
		zap.L().Info(sqlPage)
		return nil, total, err
	}

	sqlTotal = strings.Join([]string{sqlTotal, conditionSql}, " ")
	if err := db.Raw(sqlTotal).Scan(&total).Error; err != nil {
		zap.L().Info(sql)
		return nil, total, err
	}

	return &ret, total, nil
}

func (baitTask *BaitTask) QueryBaitTaskPageNew(queryMap map[string]string) (*[]BaitTask, int64, error) {
	var ret []BaitTask
	var count int64

	sql := "select * from bait_tasks "

	conditionFlag := false
	for key, val := range queryMap {
		if util.CheckInjectionData(val) {
			return nil, 0, nil
		}
		condition := "where"
		if !conditionFlag {
			conditionFlag = true
		} else {
			condition = "and"
		}
		if key == "BaitType" {
			sql = fmt.Sprintf(" %s %s bait_type = '%s'", sql, condition, val)
			conditionFlag = true
		}
		if key == "HoneypotId" {
			sql = fmt.Sprintf(" %s %s honeypot_id = '%s'", sql, condition, val)
		}
		if key == "AgentToken" {
			sql = fmt.Sprintf(" %s %s agent_token = '%s'", sql, condition, val)
		}
		if key == "Payload" {
			sql = fmt.Sprintf(" %s %s attack_payload like '%s'", sql, "%"+condition+"%", val)
		}
	}
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, 0, err
	}
	count = (int64)(len(ret))

	pageSize, _ := strconv.Atoi(queryMap["PageSize"])
	PageNumber, _ := strconv.Atoi(queryMap["PageNumber"])

	t := fmt.Sprintf("%s limit %d offset %d", sql, pageSize, (PageNumber-1)*pageSize)
	sql = strings.Join([]string{sql, t}, " ")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, 0, err
	}
	return &ret, count, nil
}

func (baitTask *BaitTask) QueryBaitTasks(payload *comm.BaitTaskQueryPayload) (*[]BaitTask, error) {
	var ret []BaitTask
	if util.CheckInjectionData(payload.Payload) {
		return nil, nil
	}
	p := "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select * from bait_tasks where CONCAT(bait_name, bait_type, deploy_path) LIKE '%s' order by create_time DESC", p)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (baitTask *BaitTask) DeleteBaitTaskById(id string) error {
	if err := db.Where("bait_task_id= ?", id).Delete(&BaitTask{}).Error; err != nil {
		return err
	}
	return nil
}

func (baitTask *BaitTask) UpdateBaitTaskStatusById(status comm.TaskStatus, baitTaskId string) error {
	if err := db.Model(baitTask).Where("bait_task_id = ?", baitTaskId).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
