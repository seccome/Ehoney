package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type ProtocolProxy struct {
	Id                int64           `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`                   //协议代理ID
	ProtocolProxyId   string          `gorm:"not null;" json:"ProtocolProxyId" form:"protocol_proxy_id" binding:"required"`     //代理端口
	ProtocolProxyName string          `gorm:"not null;" json:"ProtocolProxyName" form:"protocol_proxy_name" binding:"required"` //代理端口
	ProtocolId        string          `gorm:"not null;" json:"ProtocolId" form:"protocol_id" binding:"required"`                //代理端口
	ProtocolPath      string          `gorm:"not null;" json:"ProtocolPath" form:"protocol_path" binding:"required"`            //代理端口
	ProtocolType      string          `gorm:"not null;size:32" json:"ProtocolType" form:"ProtocolType"  binding:"required"`     //协议类型
	ProxyPort         int32           `gorm:"not null;" json:"ProxyPort" form:"proxy_port" binding:"required"`                  //代理端口
	HoneypotId        string          `gorm:"not null;" json:"HoneypotId" form:"honeypot_id" binding:"required"`                //代理端口                                    //蜜罐ID
	HoneypotIp        string          `gorm:"not null;" json:"HoneypotIp" form:"honeypot_ip" `                                  //代理端口
	HoneypotPodName   string          `gorm:"not null;" json:"HoneypotPodName" form:"honeypot_pod_name"`                        //代理端口
	HoneypotPort      int32           `gorm:"not null;" json:"HoneypotPort" form:"honeypot_port" binding:"required"`            //代理端口
	ProcessPid        int             `gorm:"not null;" json:"ProcessPid" form:"process_pid"`                                   //代理端口
	Status            comm.TaskStatus `gorm:"not null;" json:"Status" form:"status"`
	CreateTime        int64           `gorm:"not null;" json:"CreateTime" form:"CreateTime"`
}

func (proxy *ProtocolProxy) QueryProtocol2PodGreenLines() ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	sql := fmt.Sprintf("SELECT concat(\"%s-RELAY\") AS Source, concat(honeypot_ip, \"-POD\") AS Target, \"GREEN\" Status FROM protocol_proxies  WHERE status = 3 GROUP BY honeypot_ip", configs.GetSetting().Server.AppHost)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (proxy *ProtocolProxy) QueryProtocolProxyNode() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode
	sql := fmt.Sprintf("SELECT concat(\"%s-RELAY\") AS Id, \"RELAY\" NodeType,  \"%s\" AS Ip, \"HoneyNet\" AS HostName FROM protocol_proxies WHERE status = 3", configs.GetSetting().Server.AppHost, configs.GetSetting().Server.AppHost)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (proxy *ProtocolProxy) QueryProtocolAttackNode() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode
	sql := fmt.Sprintf("SELECT concat(attack_ip, \"-HACK\") AS Id,  \"HACK\" NodeType, attack_ip AS Ip, attack_ip AS HostName FROM protocol_events WHERE attack_ip NOT IN (SELECT agent_ip FROM agents GROUP BY agent_ip) AND create_time > %d GROUP BY attack_ip", util.GetCurrentIntTime()-(10*60))
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (proxy *ProtocolProxy) CreateProtocolProxy() error {
	result := db.Create(proxy)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (proxy *ProtocolProxy) GetProtocolProxyByID(protocolProxyId string) (*ProtocolProxy, error) {
	var ret ProtocolProxy
	if err := db.Take(&ret, "protocol_proxy_id = ?", protocolProxyId).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *ProtocolProxy) QueryProtocolProxyByHoneypot(honeypot Honeypot) (*[]ProtocolProxy, error) {
	var ret []ProtocolProxy
	if err := db.Take(&ret, "honeypot_id = ?", honeypot.HoneypotId).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *ProtocolProxy) GetProtocolProxy(payload *comm.SelectPayload) (*[]ProtocolProxy, int64, error) {
	var ret []ProtocolProxy
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}
	var p string = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select * from protocol_proxies where protocol_proxy_name like '%s'", p)
	t := fmt.Sprintf("limit %d offset %d", payload.PageSize, (payload.PageNumber-1)*payload.PageSize)
	sql = strings.Join([]string{sql, t}, " ")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, 0, err
	}

	sql = fmt.Sprintf("select count(1) from protocol_proxies where protocol_proxy_name like '%s'", p)
	if err := db.Raw(sql).Scan(&count).Error; err != nil {
		return nil, 0, err
	}

	return &ret, count, nil
}

func (proxy *ProtocolProxy) DeleteProtocolProxyByID(id string) error {
	if err := db.Where("protocol_proxy_id= ?", id).Delete(&ProtocolProxy{}).Error; err != nil {
		return err
	}
	return nil
}

func (proxy *ProtocolProxy) GetProtocolProxyByProxyPort(port int32) (*ProtocolProxy, error) {
	var ret ProtocolProxy
	if err := db.Take(&ret, "proxy_port = ?", port).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *ProtocolProxy) GetProtocolProxyByHoneypotID(honeypotID string) (*ProtocolProxy, error) {
	var ret ProtocolProxy
	if err := db.Take(&ret, "honeypot_id = ?", honeypotID).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *ProtocolProxy) UpdateStatus() error {
	if err := db.Model(proxy).Where("protocol_proxy_id = ?", proxy.ProtocolProxyId).Update("status", proxy.Status).Update("process_pid", proxy.ProcessPid).Error; err != nil {
		return err
	}
	return nil
}

func (proxy *ProtocolProxy) UpdateHoneypot() error {
	if err := db.Model(proxy).Where("protocol_proxy_id = ?", proxy.ProtocolProxyId).Update("honeypot_ip", proxy.HoneypotIp).Error; err != nil {
		return err
	}
	return nil
}

func (proxy *ProtocolProxy) GetAllProtocolProxies() (*[]ProtocolProxy, error) {
	var ret []ProtocolProxy
	if err := db.Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *ProtocolProxy) GetNewDeployProtocolProxies() (*[]ProtocolProxy, error) {
	var ret []ProtocolProxy
	if err := db.Take(&ret, "status = ?", 1).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
