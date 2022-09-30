package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type Protocols struct {
	Id           int64           `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"Id"`                       //协议ID
	ProtocolId   string          `json:"ProtocolId" form:"ProtocolId" gorm:"not null;" binding:"required"`            //镜像地址
	ProtocolType string          `json:"ProtocolType" form:"ProtocolType" gorm:"not null;size:32" binding:"required"` //协议类型
	LocalPath    string          `json:"LocalPath" form:"LocalPath" gorm:"not null;size:256"`                         //本地路径
	FileName     string          `json:"FileName" form:"FileName" gorm:"not null;size:256"`                           //文件名称
	MinPort      int32           `json:"MinPort" form:"MinPort" gorm:"null, default:1" binding:"required"`
	MaxPort      int32           `json:"MaxPort" form:"MaxPort" gorm:"null, default:65535" binding:"required"`
	CreateTime   int64           `gorm:"not null" json:"CreateTime"`
	DefaultFlag  bool            `json:"DefaultFlag" form:"DefaultFlag" gorm:"null, default:false"` //默认属性
	Status       comm.TaskStatus `json:"Status" form:"Status" gorm:"not null"`                      //状态

}

var DefaultProtocol = []Protocols{
	{ProtocolType: "httpproxy", ProtocolId: util.GenerateId(), CreateTime: util.GetCurrentIntTime(), LocalPath: "/tool/protocol", FileName: "httpproxy", Status: comm.SUCCESS, DefaultFlag: true, MinPort: 18080, MaxPort: 18082},
	{ProtocolType: "mysqlproxy", ProtocolId: util.GenerateId(), CreateTime: util.GetCurrentIntTime(), LocalPath: "/tool/protocol", FileName: "mysqlproxy", Status: comm.SUCCESS, DefaultFlag: true, MinPort: 13306, MaxPort: 13308},
	{ProtocolType: "redisproxy", ProtocolId: util.GenerateId(), CreateTime: util.GetCurrentIntTime(), LocalPath: "/tool/protocol", FileName: "redisproxy", Status: comm.SUCCESS, DefaultFlag: true, MinPort: 16379, MaxPort: 16381},
	{ProtocolType: "sshproxy", ProtocolId: util.GenerateId(), CreateTime: util.GetCurrentIntTime(), LocalPath: "/tool/protocol", FileName: "sshproxy", Status: comm.SUCCESS, DefaultFlag: true, MinPort: 10020, MaxPort: 10022},
	{ProtocolType: "telnetproxy", ProtocolId: util.GenerateId(), CreateTime: util.GetCurrentIntTime(), LocalPath: "/tool/protocol", FileName: "telnetproxy", Status: comm.SUCCESS, DefaultFlag: true, MinPort: 10023, MaxPort: 10025},
	{ProtocolType: "smbproxy", ProtocolId: util.GenerateId(), CreateTime: util.GetCurrentIntTime(), LocalPath: "/tool/protocol", FileName: "smbproxy", Status: comm.SUCCESS, DefaultFlag: true, MinPort: 10445, MaxPort: 10447},
	{ProtocolType: "ftpproxy", ProtocolId: util.GenerateId(), CreateTime: util.GetCurrentIntTime(), LocalPath: "/tool/protocol", FileName: "ftpproxy", Status: comm.SUCCESS, DefaultFlag: true, MinPort: 10021, MaxPort: 10023},
}

func (protocol *Protocols) CreateProtocol() error {
	result := db.Create(protocol)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (protocol *Protocols) CreateDefaultProtocol() error {
	for _, d := range DefaultProtocol {
		p, _ := protocol.GetProtocolByType(d.ProtocolType)
		if p != nil {
			continue
		}
		err := d.CreateProtocol()
		if err != nil {
			continue
		}
	}
	return nil
}

func (protocol *Protocols) GetProtocolByID(protocolId string) (*Protocols, error) {
	var ret Protocols
	if err := db.Where("protocol_id = ?", protocolId).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (protocol *Protocols) CheckProtocol(agentId string) *Protocols {
	var ret Protocols
	if err := db.First(&ret, "agent_id = ?", agentId).Error; err != nil {
		return nil
	}
	return &ret
}

func (protocol *Protocols) GetProtocolByType(protocolType string) (*Protocols, error) {
	var ret Protocols
	if err := db.Where("protocol_type = ?", protocolType).First(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (protocol *Protocols) GetProtocol(payload *comm.SelectPayload) (*[]Protocols, int64, error) {
	var ret []Protocols
	var count int64
	var p = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select * from protocols where CONCAT(id, create_time, protocol_type, local_path, min_port, max_port) LIKE '%s' order by create_time DESC", p)
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

func (protocol *Protocols) GetProtocolTypeList() (*[]string, error) {
	var ret []string
	if err := db.Model(protocol).Select("protocol_type").Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (protocol *Protocols) DeleteProtocolByID(id string) error {
	if err := db.Where("protocol_id= ?", id).Delete(&Protocols{}).Error; err != nil {
		return err
	}
	return nil
}

func (protocol *Protocols) UpdateProtocolPortRange(min, max int32, ID string) error {
	if err := db.Model(protocol).Where("protocol_id = ?", ID).Update("min_port", min).Update("max_port", max).Error; err != nil {
		return err
	}
	return nil
}
