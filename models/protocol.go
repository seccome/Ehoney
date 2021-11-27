package models

import (
	"decept-defense/controllers/comm"
	"fmt"
	"strings"
)

type Protocols struct {
	ID           int64           `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`                               //协议ID
	ProtocolType string          `json:"ProtocolType" form:"ProtocolType" gorm:"unique;not null;size:128" binding:"required"` //协议类型
	CreateTime   string          `json:"CreateTime" form:"CreateTime" gorm:"not null"`                                        //创建时间
	DeployPath   string          `json:"DeployPath" form:"DeployPath" gorm:"not null;size:256"`                               //部署路径
	LocalPath    string          `json:"LocalPath" form:"LocalPath" gorm:"not null;size:256"`                                 //本地路径
	FileName     string          `json:"FileName" form:"FileName" gorm:"not null;size:256"`                                   //文件名称
	Creator      string          `json:"Creator" form:"Creator" gorm:"not null;size:256"`                                     //创建用户
	Status       comm.TaskStatus `json:"Status" form:"Status" gorm:"not null"`                                                //状态
	TaskID       string          `json:"TaskID" form:"TaskID" gorm:"not null"`                                                //任务ID
	MinPort      int32           `json:"MinPort" form:"MinPort" gorm:"null, default:1" binding:"required"`
	MaxPort      int32           `json:"MaxPort" form:"MaxPort" gorm:"null, default:65535" binding:"required"`
	DefaultFlag  bool            `json:"DefaultFlag" form:"DefaultFlag" gorm:"null, default:false"` //默认属性
}

var DefaultProtocol = []Protocols{
	{ProtocolType: "httpproxy", CreateTime: "2021-8-3 19:19:30", DeployPath: "/home/ehoney_proxy", LocalPath: "/fake/path", FileName: "httpproxy", Creator: "default", Status: comm.SUCCESS, TaskID: "/fake/task/id", DefaultFlag: true, MinPort: 18080, MaxPort: 18082},
	{ProtocolType: "mysqlproxy", CreateTime: "2021-8-3 19:19:30", DeployPath: "/home/ehoney_proxy", LocalPath: "/fake/path", FileName: "mysqlproxy", Creator: "default", Status: comm.SUCCESS, TaskID: "/fake/task/id", DefaultFlag: true, MinPort: 13306, MaxPort: 13308},
	{ProtocolType: "redisproxy", CreateTime: "2021-8-3 19:19:30", DeployPath: "/home/ehoney_proxy", LocalPath: "/fake/path", FileName: "redisproxy", Creator: "default", Status: comm.SUCCESS, TaskID: "/fake/task/id", DefaultFlag: true, MinPort: 16379, MaxPort: 16381},
	{ProtocolType: "sshproxy", CreateTime: "2021-8-3 19:19:30", DeployPath: "/home/ehoney_proxy", LocalPath: "/fake/path", FileName: "sshproxy", Creator: "default", Status: comm.SUCCESS, TaskID: "/fake/task/id", DefaultFlag: true, MinPort: 10020, MaxPort: 10022},
	{ProtocolType: "telnetproxy", CreateTime: "2021-8-3 19:19:30", DeployPath: "/home/ehoney_proxy", LocalPath: "/fake/path", FileName: "telnetproxy", Creator: "default", Status: comm.SUCCESS, TaskID: "/fake/task/id", DefaultFlag: true, MinPort: 10023, MaxPort: 10025},
	{ProtocolType: "smbproxy", CreateTime: "2021-8-3 19:19:30", DeployPath: "/home/ehoney_proxy", LocalPath: "/fake/path", FileName: "smbproxy", Creator: "default", Status: comm.SUCCESS, TaskID: "/fake/task/id", DefaultFlag: true, MinPort: 10445, MaxPort: 10447},
	{ProtocolType: "ftpproxy", CreateTime: "2021-8-3 19:19:30", DeployPath: "/home/ehoney_proxy", LocalPath: "/fake/path", FileName: "ftpproxy", Creator: "default", Status: comm.SUCCESS, TaskID: "/fake/task/id", DefaultFlag: true, MinPort: 10021, MaxPort: 10023},
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

func (protocol *Protocols) GetProtocolByID(id int64) (*Protocols, error) {
	var ret Protocols
	if err := db.Where("id = ?", id).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (protocol *Protocols) CheckProtocol(agentID string) *Protocols {
	var ret Protocols
	if err := db.First(&ret, "agent_id = ?", agentID).Error; err != nil {
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

func (protocol *Protocols) GetProtocol(payload *comm.SelectPayload) (*[]comm.ProtocolSelectResultPayload, int64, error) {
	var ret []comm.ProtocolSelectResultPayload
	var count int64
	var p = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select id, creator, status, create_time, protocol_type, deploy_path, default_flag, min_port, max_port from protocols where CONCAT(id, creator, create_time, protocol_type, deploy_path, min_port, max_port) LIKE '%s' order by create_time DESC", p)
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

func (protocol *Protocols) GetProtocolByTaskID(taskID string) (*Protocols, error) {
	var ret Protocols
	if err := db.Where("task_id = ?", taskID).First(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (protocol *Protocols) GetProtocolTypeList() (*[]string, error) {
	var ret []string
	if err := db.Model(protocol).Select("protocol_type").Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (protocol *Protocols) DeleteProtocolByID(id int64) error {
	if err := db.Delete(&Protocols{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (protocol *Protocols) UpdateProtocolStatusByTaskID(status int, taskID string) error {
	if err := db.Model(protocol).Where("task_id = ?", taskID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (protocol *Protocols) UpdateProtocolPortRange(min, max int32, ID int64) error {
	if err := db.Model(protocol).Where("id = ?", ID).Update("min_port", min).Update("max_port", max).Error; err != nil {
		return err
	}
	return nil
}
