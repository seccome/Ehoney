package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"strings"
)

type VirusRecord struct {
	ID            int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`
	VirusName     string `json:"VirusName" form:"VirusName" gorm:"not null;size:256" binding:"required"`
	HoneypotIP    string `json:"HoneypotIP" form:"HoneypotIP" gorm:"not null;size:256" binding:"required"`
	CreateTime    string `json:"CreateTime" form:"CreateTime" gorm:"not null"`
	VirusFilePath string `json:"VirusFilePath" form:"VirusFilePath" gorm:"not null" binding:"required"`
}

func (virusRecord *VirusRecord) CreateVirusRecord() error {
	result := db.Create(virusRecord)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (virusRecord *VirusRecord) DeleteVirusRecordByID(id int64) error {
	if err := db.Delete(&VirusRecord{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (virusRecord *VirusRecord) GetVirusRecord(payload *comm.SelectVirusPayload) (*[]VirusRecord, error) {
	var virusRecords []VirusRecord
	if util.CheckInjectionData(payload.VirusName) || util.CheckInjectionData(payload.HoneypotName) || util.CheckInjectionData(payload.VirusFilePath) {
		return nil, nil
	}

	virusName := strings.Join([]string{"%", payload.VirusName, "%"}, "")
	honeypotIP := strings.Join([]string{"%", payload.HoneypotName, "%"}, "")
	virusFilePath := strings.Join([]string{"%", payload.VirusFilePath, "%"}, "")
	if payload.StartTimestamp != 0 && payload.EndTimestamp != 0 {
		if err := db.Limit(payload.PageSize).Offset((payload.PageNumber-1)*payload.PageSize).Where("virus_name LIKE ? AND honeypot_ip LIKE ? AND virus_file_path LIKE ? AND create_time BETWEEN ? AND ?", virusName, honeypotIP, virusFilePath, util.Sec2TimeStr(payload.StartTimestamp, ""), util.Sec2TimeStr(payload.EndTimestamp, "")).Find(&virusRecords).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Limit(payload.PageSize).Offset((payload.PageNumber-1)*payload.PageSize).Where("virus_name LIKE ? AND honeypot_ip LIKE ? AND virus_file_path LIKE ?", virusName, honeypotIP, virusFilePath).Find(&virusRecords).Error; err != nil {
			return nil, err
		}
	}
	return &virusRecords, nil
}
