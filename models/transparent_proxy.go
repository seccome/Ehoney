package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type TransparentProxy struct {
	Id                 int64           `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`                                  //透明代理ID //创建时间
	TransparentProxyId string          `gorm:"unique;not null; size:64" json:"TransparentProxyId" form:"TransparentProxyId" binding:"required"` //代理端口
	ProtocolProxyId    string          `gorm:"index;not null;" json:"ProtocolProxyId" form:"ProtocolProxyId" binding:"required"`                //协议代理ID
	ProtocolType       string          `json:"ProtocolType" form:"ProtocolType" gorm:"not null;size:32" binding:"required"`                     //协议类型
	ProxyPort          int32           `gorm:"not null;" json:"ProxyPort" form:"ProxyPort" binding:"required"`                                  //代理端口
	DestIp             string          `gorm:"not null;" json:"DestIp" form:"DestIp" binding:"required"`                                        //代理端口
	DestPort           int32           `gorm:"not null;" json:"DestPort" form:"DestPort" binding:"required"`                                    //代理端口
	AgentId            string          `gorm:"not null;size:64" json:"AgentId"`
	AgentToken         string          `gorm:"not null;size:128" json:"AgentToken" form:"AgentToken" `
	AgentIp            string          `gorm:"not null;size:256" json:"AgentIp" form:"AgentIp" `
	AgentHost          string          `gorm:"not null;size:128" json:"AgentHost" form:"AgentHost"`
	Status             comm.TaskStatus `gorm:"not null;" json:"Status" form:"Status"`
	CreateTime         int64           `gorm:"not null;" json:"CreateTime" form:"CreateTime"`
}

func (proxy *TransparentProxy) QueryTransparent2ProtocolGreenLines() ([]comm.TopologyLine, error) {
	var ret []comm.TopologyLine
	sql := fmt.Sprintf("SELECT concat(p.server_ip, \"-EDGE\") AS Source, concat(hs.server_ip, \"-RELAY\") AS Target,  \"GREEN\" Status FROM transparent_proxies tp LEFT JOIN probes p on p.id = tp.server_id LEFT JOIN protocol_proxies pp ON tp.protocol_proxy_id = pp.id LEFT JOIN honeypots h ON pp.honeypot_id = h.id LEFT JOIN honeypot_servers hs on h.servers_id = hs.id GROUP BY p.server_ip, hs.server_ip")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (proxy *TransparentProxy) GetTransparentProxyNodes() ([]comm.TopologyNode, error) {
	var ret []comm.TopologyNode
	sql := fmt.Sprintf("SELECT concat(agent_ip, \"-EDGE\") AS Id,  \"EDGE\" NodeType, agent_ip AS Ip, agent_host AS HostName FROM transparent_proxies WHERE status = 3 GROUP BY HostName ")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (proxy *TransparentProxy) CreateTransparentProxy() error {
	proxy.TransparentProxyId = util.GenerateId()
	proxy.CreateTime = util.GetCurrentIntTime()
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

func (proxy *TransparentProxy) GetTransparentProxyByID(transparentProxyId string) (*TransparentProxy, error) {
	var ret TransparentProxy
	if err := db.Take(&ret, "transparent_proxy_id = ?", transparentProxyId).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *TransparentProxy) GetTransparentProxy(queryMap map[string]interface{}) (*[]TransparentProxy, int64, error) {
	var ret []TransparentProxy
	//

	var total int64
	sql := fmt.Sprintf("SELECT * FROM transparent_proxies ")
	sqlTotal := fmt.Sprintf("SELECT count(1) FROM transparent_proxies ")

	conditionFlag := false
	conditionSql := ""
	for key, val := range queryMap {

		if key == "PageSize" || key == "PageNumber" {
			continue
		}
		if val == "" {
			continue
		}
		if util.CheckInjectionData(val.(string)) {
			return nil, 0, nil
		}
		condition := "where"
		if !conditionFlag {
			conditionFlag = true
		} else {
			condition = "and"
		}
		if key == "AgentId" {
			conditionSql = fmt.Sprintf(" %s %s agent_id = '%s'", conditionSql, condition, val)
		}
		if key == "Payload" {
			conditionSql = fmt.Sprintf(" %s %s concat(proxy_port, dest_ip, dest_port) like '%s'", conditionSql, condition, "%"+val.(string)+"%")
		}
	}

	pageSize := int(queryMap["PageSize"].(float64))

	pageNumber := int(queryMap["PageNumber"].(float64))

	t := fmt.Sprintf("order by create_time DESC limit %d offset %d ", pageSize, (pageNumber-1)*pageSize)
	sql = strings.Join([]string{sql, conditionSql, t}, " ")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, total, err
	}

	sqlTotal = strings.Join([]string{sqlTotal, conditionSql}, " ")
	if err := db.Raw(sqlTotal).Scan(&total).Error; err != nil {
		return nil, total, err
	}
	return &ret, total, nil
}

func (proxy *TransparentProxy) CheckTransparentById(transparentProxyId string) *TransparentProxy {
	var ret TransparentProxy
	if err := db.Take(&ret, "transparent_proxy_id = ?", transparentProxyId).Error; err != nil {
		return nil
	}
	return &ret
}

func (proxy *TransparentProxy) DeleteTransparentProxyByID(id string) error {
	if err := db.Where("transparent_proxy_id = ?", id).Delete(&TransparentProxy{}).Error; err != nil {
		return err
	}
	return nil
}

func (proxy *TransparentProxy) GetTransparentProxyByProxyPort(port int32, agentId string) (*TransparentProxy, error) {
	var ret TransparentProxy
	if err := db.Take(&ret, "proxy_port = ? and agent_id = ?", port, agentId).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (proxy *TransparentProxy) UpdateTransparentProxyStatusByID(status comm.TaskStatus, transparentProxyId string) error {
	if err := db.Model(proxy).Where("transparent_proxy_id = ?", transparentProxyId).Update("status", status).Error; err != nil {
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
