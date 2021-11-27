package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type HoneypotBaits struct {
	ID         int64    `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`
	CreateTime string   `gorm:"not null"`
	Creator    string   `gorm:"not null;size:256" json:"Creator"`
	LocalPath  string   `gorm:"not null;size:256" json:"LocalPath"`
	DeployPath string   `gorm:"not null;size:256" json:"DeployPath"`
	BaitName   string   `gorm:"not null;size:256" json:"BaitName"`
	BaitType   string   `gorm:"not null;size:256" json:"BaitType"`
	BaitData   string   `gorm:"not null;size:256" json:"BaitData"`
	HoneypotID int64    `gorm:"not null" json:"HoneypotID"`
	Honeypot   Honeypot `gorm:"ForeignKey:HoneypotID;constraint:OnDelete:CASCADE"`
}

func (honeypotBaits *HoneypotBaits) CreateHoneypotBait() error {
	if err := db.Create(honeypotBaits).Error; err != nil {
		return err
	}
	return nil
}

func (honeypotBaits *HoneypotBaits) GetHoneypotBaitByID(id int64) (*HoneypotBaits, error) {
	var ret HoneypotBaits
	if err := db.Take(&ret, id).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (honeypotBaits *HoneypotBaits) GetHoneypotBait(payload *comm.ServerBaitSelectPayload) (*[]comm.ServerBaitSelectResultPayload, int64, error) {
	var ret []comm.ServerBaitSelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}

	var p string = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select h.id, h.bait_name, h.bait_type, h.creator, h.deploy_path, h.create_time from honeypot_baits h where h.honeypot_id = %d  AND CONCAT(h.id, h.bait_name, h.bait_type, h.creator,  h.deploy_path, h.create_time) LIKE '%s' order by h.create_time DESC", payload.ServerID, p)
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

func (honeypotBaits *HoneypotBaits) DeleteHoneypotBaitByID(id int64) error {
	if err := db.Delete(&HoneypotBaits{}, id).Error; err != nil {
		return err
	}
	return nil
}
