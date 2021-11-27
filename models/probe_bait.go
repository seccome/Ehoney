package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type ProbeBaits struct {
	ID         int64           `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`
	CreateTime string          `gorm:"not null"`
	Creator    string          `gorm:"not null;size:256" json:"Creator"`
	DeployPath string          `gorm:"not null;size:256" json:"DeployPath"`
	LocalPath  string          `gorm:"not null;size:256" json:"LocalPath"`
	BaitName   string          `gorm:"not null;size:256" json:"BaitName"`
	BaitType   string          `gorm:"not null;size:256" json:"BaitType"`
	ServerID   int64           `gorm:"not null" json:"ServerID"`
	Servers    Probes          `gorm:"ForeignKey:ServerID;constraint:OnDelete:CASCADE"`
	Status     comm.TaskStatus `gorm:"not null" json:"Status"`
	TaskID     string          `gorm:"not null" json:"TaskID"`
}

func (probeBaits *ProbeBaits) CreateProbeBait() error {
	if err := db.Create(probeBaits).Error; err != nil {
		return err
	}
	return nil
}

func (probeBaits *ProbeBaits) GetProbeBaitByID(id int64) (*ProbeBaits, error) {
	var ret ProbeBaits
	if err := db.Take(&ret, id).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (probeBaits *ProbeBaits) GetProbeBait(payload *comm.ServerBaitSelectPayload) (*[]comm.ServerBaitSelectResultPayload, int64, error) {
	var ret []comm.ServerBaitSelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}
	var p string = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select h.id, h.bait_name, h.bait_type, h.deploy_path, h.creator, h.status, h.create_time from probe_baits h where h.server_id = %d   AND CONCAT(h.id, h.bait_name, h.bait_type, h.deploy_path, h.creator, h.create_time) LIKE '%s' order by h.create_time DESC", payload.ServerID, p)
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

func (probeBaits *ProbeBaits) DeleteProbeBaitByID(id int64) error {
	if err := db.Delete(&ProbeBaits{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (probeBaits *ProbeBaits) UpdateProbeBaitStatusByTaskID(status int, taskID string) error {
	if err := db.Model(probeBaits).Where("task_id = ?", taskID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}
