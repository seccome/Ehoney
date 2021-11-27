package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type Baits struct {
	ID         int64  `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`
	CreateTime string `gorm:"not null"`
	Creator    string `gorm:"not null;size:256"`
	BaitType   string `gorm:"not null;size:256" form:"BaitType" binding:"required"`
	UploadPath string `gorm:"size:256"`
	FileName   string `gorm:"size:256"`
	BaitName   string `gorm:"unique"form:"BaitName" gorm:"unique;size:256" binding:"required"`
	BaitData   string `form:"BaitData"`
}

func (bait *Baits) CreateBait() error {
	if err := db.Create(bait).Error; err != nil {
		return err
	}
	return nil
}

func (bait *Baits) GetBaitsRecord(payload *comm.BaitSelectPayload) (*[]comm.FileBaitSelectResultPayload, int64, error) {
	var ret []comm.FileBaitSelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) || util.CheckInjectionData(payload.BaitType) {
		return nil, 0, nil
	}
	var p string = "%" + payload.Payload + "%"
	var sql = ""

	if payload.BaitType != "" {
		sql = fmt.Sprintf("select id, bait_type, bait_name, file_name, bait_data, create_time, creator from baits where CONCAT(bait_type, bait_name, bait_data, file_name, create_time, creator) LIKE '%s' AND bait_type = '%s' order by create_time DESC", p, payload.BaitType)
	} else {
		sql = fmt.Sprintf("select id, bait_type, bait_name, file_name, bait_data, create_time, creator from baits where CONCAT(bait_type, bait_name, bait_data, file_name, create_time, creator) LIKE '%s' order by create_time DESC", p)
	}

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

func (bait *Baits) GetBaitByID(id int64) (*Baits, error) {
	var ret Baits
	if err := db.Take(&ret, id).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (bait *Baits) GetBaitByName(name string) (*Baits, error) {
	var ret Baits
	if err := db.Where("bait_name = ?", name).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (bait *Baits) DeleteBaitByID(id int64) error {
	if err := db.Delete(&Baits{}, id).Error; err != nil {
		return err
	}
	return nil
}
