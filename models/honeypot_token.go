package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type HoneypotToken struct {
	ID         int64           `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"` //蜜罐密签ID
	CreateTime string          `gorm:"not null"`                                                       //创建时间
	Creator    string          `gorm:"not null;size:256"`                                              //创建用户
	DeployPath string          `gorm:"not null;size:256"`                                              //部署路径
	LocalPath  string          `gorm:"not null;size:256"`                                              //本地路径
	TokenType  string          `gorm:"not null;size:256"`                                              //密签类型
	TokenName  string          `gorm:"not null;size:256"`                                              //密签名称
	HoneypotID int64           `gorm:"not null" json:"HoneypotID"`                                     //蜜罐ID
	Honeypot   Honeypot        `gorm:"ForeignKey:HoneypotID;constraint:OnDelete:CASCADE"`
	Status     comm.TaskStatus `gorm:"not null" json:"Status"`    //状态
	TraceCode  string          `gorm:"not null" json:"TraceCode"` //跟踪码
}

func (honeypotToken *HoneypotToken) CreateHoneypotToken() error {
	if err := db.Create(honeypotToken).Error; err != nil {
		return err
	}
	return nil
}

func (honeypotToken *HoneypotToken) GetHoneypotTokenByID(id int64) (*HoneypotToken, error) {
	var ret HoneypotToken
	if err := db.Take(&ret, id).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (honeypotToken *HoneypotToken) GetHoneypotToken(payload *comm.ServerTokenSelectPayload) (*[]comm.ServerTokenSelectResultPayload, int64, error) {
	var ret []comm.ServerTokenSelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}

	var p string = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select h.id, h.token_type, h.token_name, h.deploy_path, h.creator, h.create_time from honeypot_tokens h where h.honeypot_id = %d  AND CONCAT(h.id, h.token_type, h.token_name,  h.deploy_path, h.creator, h.create_time) LIKE '%s' order by h.create_time DESC", payload.ServerID, p)
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

func (honeypotToken *HoneypotToken) DeleteHoneypotTokenByID(id int64) error {
	if err := db.Delete(&HoneypotToken{}, id).Error; err != nil {
		return err
	}
	return nil
}
