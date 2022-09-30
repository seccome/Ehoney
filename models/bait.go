package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

/**
诱饵包含了之前的蜜签类型 主要考虑蜜钱基本都是文件、可以统一管理; 蜜签从一种类似诱饵类型降格为一种属性; 文件型诱饵的一种属性
文件诱饵加签之后文件路径不变 进行替换 获得 TokenTraceUrl 值


*/

type Bait struct {
	Id            int64  `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`
	BaitId        string `json:"BaitId" form:"BaitId" gorm:"unique;not null;size:64" binding:"required"`
	BaitType      string `json:"BaitType" form:"BaitType" gorm:"not null;size:64" binding:"required"`
	BaitName      string `json:"BaitName" form:"BaitName" gorm:"unique;size:128" binding:"required"`
	BaitData      string `json:"BaitData" form:"BaitData" gorm:"size:2048"`
	UploadPath    string `json:"UploadPath" form:"UploadPath" gorm:"size:256" `
	LocalPath     string `json:"LocalPath" form:"LocalPath" gorm:"size:256" `
	FileName      string `json:"FileName" form:"FileName" gorm:"size:256"`
	TokenAble     bool   `json:"TokenAble" form:"TokenAble"  gorm:"size:3" form:"TokenAble"`                  // 根据类型判断是否可加签
	TokenTraceUrl string `json:"TokenTraceUrl" form:"TokenTraceUrl" gorm:"size:256"  binding:"TokenTraceUrl"` // 文件型诱饵加签后 获得 TokenTraceUrl 值
	CreateTime    int64  `json:"CreateTime" form:"CreateTime" gorm:"size:64"`
}

var TokenAbleBaitTypeArray = [2]string{"WPS"}
var TokenAbleFileTypeArray = [5]string{".pdf", ".docx", ".xlsx", ".pptx", ".exe"}

func (bait *Bait) IsTokenAble() bool {
	for _, bt := range TokenAbleBaitTypeArray {
		if bait.BaitType == bt {
			return true
		}
	}
	for _, ft := range TokenAbleFileTypeArray {
		if strings.HasSuffix(bait.FileName, ft) {
			return true
		}
	}
	return false
}

func (bait *Bait) CreateBait() error {
	bait.BaitId = util.GenerateId()
	bait.CreateTime = util.GetCurrentIntTime()
	bait.TokenAble = bait.IsTokenAble()
	if err := db.Create(bait).Error; err != nil {
		return err
	}
	return nil
}

func (bait *Bait) GetBaitsRecord(payload *comm.BaitSelectPayload) (*[]Bait, int64, error) {
	var ret []Bait
	var count int64
	if util.CheckInjectionData(payload.Payload) || util.CheckInjectionData(payload.BaitType) {
		return nil, 0, nil
	}
	var p string = "%" + payload.Payload + "%"
	var sql = ""

	if payload.BaitType != "" {
		sql = fmt.Sprintf("select * from baits where CONCAT(bait_type, bait_name, bait_data, file_name, create_time) LIKE '%s' AND bait_type = '%s' order by create_time DESC", p, payload.BaitType)
	} else {
		sql = fmt.Sprintf("select * from baits where CONCAT(bait_type, bait_name, bait_data, file_name, create_time) LIKE '%s' order by create_time DESC", p)
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

func (bait *Bait) UpdateForToken() error {
	if err := db.Model(bait).Where("bait_id = ?", bait.BaitId).Update("local_path", bait.LocalPath).Update("token_trace_url", bait.TokenTraceUrl).Error; err != nil {
		return err
	}
	return nil
}

func (bait *Bait) GetBaitById(baitId string) (*Bait, error) {
	var ret Bait
	if err := db.Where("bait_id = ?", baitId).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (bait *Bait) GetBaitByName(name string) (*Bait, error) {
	var ret Bait
	if err := db.Where("bait_name = ?", name).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (bait *Bait) DeleteBaitById(baitId string) error {
	if err := db.Where("bait_id= ?", baitId).Delete(&Bait{}).Error; err != nil {
		return err
	}
	return nil
}
