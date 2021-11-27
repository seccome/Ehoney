package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type ProtocolProxy struct {
	ID         int64           `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"` //协议代理ID
	CreateTime string          `gorm:"not null"`                                                       //创建时间
	Creator    string          `gorm:"not null;size:256"`                                              //创建用户
	AgentID    string          `gorm:"not null;size:256" json:"AgentID"`
	ProxyPort  int32           `gorm:"not null;unique" json:"ProxyPort" form:"ProxyPort" binding:"required"` //代理端口
	HoneypotID int64           `binding:"required" json:"HoneypotID"`                                        //蜜罐ID
	Honeypot   Honeypot        `gorm:"ForeignKey:HoneypotID"`
	ProtocolID int64           //协议ID
	Protocols  Protocols       `gorm:"ForeignKey:ProtocolID"`
	TaskID     string          //任务ID
	Status     comm.TaskStatus `json:"Status" form:"Status"` //状态
}

func (proxy *ProtocolProxy) QueryProtocol2PodGreenLines() ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	sql := fmt.Sprintf("SELECT concat(hs.server_ip, \"-RELAY\") AS Source, concat(h.honeypot_ip, \"-POD\") AS Target, \"GREEN\" Status FROM protocol_proxies pp LEFT JOIN honeypots h ON pp.honeypot_id = h.id LEFT JOIN honeypot_servers hs ON h.servers_id = hs.id GROUP BY hs.server_ip, h.honeypot_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (proxy *ProtocolProxy) QueryProtocolProxyNode() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode
	sql := fmt.Sprintf("SELECT concat(h.server_ip, \"-RELAY\") AS Id, \"RELAY\" NodeType, h.server_ip AS Ip, h.host_name AS HostName FROM protocol_proxies pp INNER JOIN honeypot_servers h ON pp.agent_id = h.agent_id WHERE pp.status = 3 GROUP BY  h.server_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (proxy *ProtocolProxy) QueryProtocolAttackNode() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode
	sql := fmt.Sprintf("SELECT concat(attack_ip, \"-HACK\") AS Id,  \"HACK\" NodeType, attack_ip AS Ip, attack_ip AS HostName FROM protocol_events pe WHERE pe.attack_ip NOT IN (SELECT server_ip FROM probes GROUP BY server_ip) AND TIMESTAMPDIFF(HOUR, pe.event_time, NOW()) < 6 GROUP BY attack_ip")
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

func (proxy *ProtocolProxy) GetProtocolProxyByID(id int64) (*ProtocolProxy, error) {
	var ret ProtocolProxy
	if err := db.Take(&ret, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *ProtocolProxy) GetProtocolProxy(payload *comm.SelectPayload) (*[]comm.ProtocolProxySelectResultPayload, int64, error) {
	var ret []comm.ProtocolProxySelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}
	var p string = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select h.id, h2.server_ip, h3.server_type, h3.honeypot_name, h3.server_port, h.proxy_port, h3.honeypot_ip, h.create_time, h.creator, h.status, h4.min_port, h4.max_port from protocol_proxies h, honeypot_servers h2, honeypots h3, protocols h4 where h.honeypot_id = h3.id AND h3.servers_id = h2.id AND h.protocol_id = h4.id AND CONCAT(h.id, h2.server_ip, h3.server_type, h3.honeypot_name, h3.server_port, h.proxy_port, h3.honeypot_ip, h.create_time, h.creator, h.status) LIKE '%s' order by h.create_time DESC", p)
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

func (proxy *ProtocolProxy) DeleteProtocolProxyByID(id int64) error {
	if err := db.Delete(&ProtocolProxy{}, id).Error; err != nil {
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

func (proxy *ProtocolProxy) GetProtocolProxyByTaskID(taskID string) (*ProtocolProxy, error) {
	var ret ProtocolProxy
	if err := db.Take(&ret, "task_id = ?", taskID).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *ProtocolProxy) GetProtocolProxyByHoneypotID(honeypotID int64) (*ProtocolProxy, error) {
	var ret ProtocolProxy
	if err := db.Take(&ret, "honeypot_id = ?", honeypotID).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *ProtocolProxy) UpdateProtocolProxyStatusByTaskID(status comm.TaskStatus, taskID string) error {
	if err := db.Model(proxy).Where("task_id = ?", taskID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (proxy *ProtocolProxy) GetProtocolProxies() (*[]ProtocolProxy, error) {
	var ret []ProtocolProxy
	if err := db.Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
