package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type ProbeToken struct {
	ID         int64           `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"` //探针密签ID
	CreateTime string          `gorm:"not null"`                                                       //创建时间
	Creator    string          `gorm:"not null;size:256" json:"CreateTime"`                            //创建用户
	DeployPath string          `gorm:"not null;size:256" json:"Creator"`                               //部署路径
	LocalPath  string          `gorm:"not null;size:256" json:"DeployPath"`                            //本地路径
	TokenName  string          `gorm:"not null;size:256" json:"TokenName"`
	TokenType  string          `gorm:"not null;size:256" json:"LocalPath"`
	ServerID   int64           `gorm:"not null" json:"ServerID"` //探针服务器ID
	Servers    Probes          `gorm:"ForeignKey:ServerID;constraint:OnDelete:CASCADE"`
	Status     comm.TaskStatus `gorm:"not null" json:"Status"`    //状态
	TaskID     string          `gorm:"not null" json:"TaskID"`    //任务ID
	TraceCode  string          `gorm:"not null" json:"TraceCode"` //跟踪码
}

func (probeTokens *ProbeToken) CreateProbeToken() error {
	if err := db.Create(probeTokens).Error; err != nil {
		return err
	}
	return nil
}

func (probeTokens *ProbeToken) GetProbeTokenByID(id int64) (*ProbeToken, error) {
	var ret ProbeToken
	if err := db.Take(&ret, id).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (probeTokens *ProbeToken) GetProbeToken(payload *comm.ServerTokenSelectPayload) (*[]comm.ServerTokenSelectResultPayload, int64, error) {
	var ret []comm.ServerTokenSelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}
	var p string = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select h.id, h.token_type, h.token_name, h.status,  h.deploy_path, h.creator, h.create_time from probe_tokens h where h.server_id = %d  AND CONCAT(h.id, h.token_type, h.token_name,  h.deploy_path, h.creator, h.create_time) LIKE '%s' order by h.create_time DESC", payload.ServerID, p)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, 0, err
	}
	count = (int64)(len(ret))
	t := fmt.Sprintf("limit %d offset %d", payload.PageSize, (payload.PageNumber-1)*payload.PageSize)
	sql = strings.Join([]string{sql, t}, " ")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, 0, err
	}
	return &ret, count, nil
}

func (probeTokens *ProbeToken) DeleteProbeTokenByID(id int64) error {
	if err := db.Delete(&ProbeToken{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (probeTokens *ProbeToken) UpdateProbeTokenStatusByTaskID(status int, taskID string) error {
	if err := db.Model(probeTokens).Where("task_id = ?", taskID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
