package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type TransparentProxy struct {
	ID              int64         `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"` //透明代理ID
	CreateTime      string        `gorm:"not null"`                                                       //创建时间
	Creator         string        `gorm:"not null;size:256"`                                              //创建用户
	ProxyPort       int32         `gorm:"not null;" json:"ProxyPort" form:"ProxyPort" binding:"required"` //代理端口
	ProtocolProxyID int64         `binding:"required" json:"ProtocolProxyID"`                             //协议代理ID
	ProtocolProxy   ProtocolProxy `gorm:"ForeignKey:ProtocolProxyID;constraint:OnDelete:CASCADE"`
	ServerID        int64         `binding:"required" json:"ProbeID"` //服务ID
	Servers         Probes        `gorm:"ForeignKey:ServerID;constraint:OnDelete:CASCADE"`
	TaskID          string
	AgentID         string          `gorm:"not null;size:256" json:"AgentID"`
	Status          comm.TaskStatus `json:"Status" form:"Status"` //状态
}

func (proxy *TransparentProxy) QueryProbe2ProtocolGreenLines() ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	sql := fmt.Sprintf("SELECT concat(p.server_ip, \"-EDGE\") AS Source, concat(hs.server_ip, \"-RELAY\") AS Target,  \"GREEN\" Status FROM transparent_proxies tp LEFT JOIN probes p on p.id = tp.server_id LEFT JOIN protocol_proxies pp ON tp.protocol_proxy_id = pp.id LEFT JOIN honeypots h ON pp.honeypot_id = h.id LEFT JOIN honeypot_servers hs on h.servers_id = hs.id GROUP BY p.server_ip, hs.server_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (proxy *TransparentProxy) GetTransparentProxyNodes() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode
	sql := fmt.Sprintf("SELECT concat(p.server_ip, \"-EDGE\") AS Id,  \"EDGE\" NodeType, p.server_ip AS Ip, p.host_name AS HostName FROM transparent_proxies tp LEFT JOIN probes p ON tp.server_id = p.id WHERE tp.status = 3 GROUP BY  p.server_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (proxy *TransparentProxy) CreateTransparentProxy() error {
	result := db.Create(proxy)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (proxy *TransparentProxy) GetTransparentByAgent(agentId string) (*[]TransparentProxy, error) {
	var ret []TransparentProxy
	if err := db.Take(&ret, "agent_id = ? and status = ?", agentId, 3).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *TransparentProxy) GetTransparentProxyByID(id int64) (*TransparentProxy, error) {
	var ret TransparentProxy
	if err := db.Take(&ret, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *TransparentProxy) GetTransparentProxy(payload *comm.SelectTransparentProxyPayload) (*[]comm.TransparentProxySelectResultPayload, int64, error) {
	var ret []comm.TransparentProxySelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}
	baseSql := "select h.id, h.proxy_port, h2.server_ip as ProbeIP, h.create_time, h.creator, h.status, h3.proxy_port as ProtocolPort, h4.protocol_type  from transparent_proxies h, probes h2, protocol_proxies h3, protocols h4 where "
	var p = "%" + payload.Payload + "%"

	if payload.Status != 0 {
		baseSql = fmt.Sprintf("%s h.status = %d ", baseSql, payload.Status)
	} else {
		baseSql = fmt.Sprintf("%s h.status != 0 ", baseSql)
	}

	var sql = ""
	if payload.ProtocolProxyID != 0 {
		sql = fmt.Sprintf("%s and h3.id = %d AND  protocol_proxy_id = %d AND h4.id = h3.protocol_id AND h.server_id = h2.id AND CONCAT(h.id, h.proxy_port, h2.server_ip, h.create_time, h.creator) LIKE '%s' order by h.create_time DESC", baseSql, payload.ProtocolProxyID, payload.ProtocolProxyID, p)
	} else {
		sql = fmt.Sprintf("%s and h.server_id = h2.id  AND h3.id = h.protocol_proxy_id AND h4.id = h3.protocol_id AND CONCAT(h.id, h.proxy_port, h2.server_ip, h.create_time, h.creator) LIKE '%s' order by h.create_time DESC", baseSql, p)
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

func (proxy *TransparentProxy) CheckTransparentByID(id int64) *TransparentProxy {
	var ret TransparentProxy
	if err := db.Take(&ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return &ret
}

func (proxy *TransparentProxy) DeleteTransparentProxyByID(id int64) error {
	if err := db.Delete(&TransparentProxy{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (proxy *TransparentProxy) GetTransparentProxyByProxyPort(port int32, probeID int64) (*TransparentProxy, error) {
	var ret TransparentProxy
	if err := db.Take(&ret, "proxy_port = ? and server_id = ?", port, probeID).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *TransparentProxy) GetTransparentProxyByTaskID(taskID string) (*TransparentProxy, error) {
	var ret TransparentProxy
	if err := db.Take(&ret, "task_id = ?", taskID).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *TransparentProxy) UpdateTransparentProxyStatusByTaskID(status comm.TaskStatus, taskID string) error {
	if err := db.Model(proxy).Where("task_id = ?", taskID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (proxy *TransparentProxy) UpdateTransparentProxyStatusByID(status int64, ID int64) error {
	if err := db.Model(proxy).Where("id = ?", ID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (proxy *TransparentProxy) GetTransparentProxies() (*[]TransparentProxy, error) {
	var ret []TransparentProxy
	if err := db.Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
