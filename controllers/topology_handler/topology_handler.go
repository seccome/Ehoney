package topology_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

var (
	LastAttackLog comm.AttackLog
	LastNodes     []comm.TopologyNode
	LastLines     []comm.TopologyLine
	NewAttack     bool
)

var WsUpGrader = websocket.Upgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,
	HandshakeTimeout:  5 * time.Second,
	// CheckOrigin: 处理跨域问题，线上环境慎用
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Setup() {
	go topology()
	go attackMonitJob()
}

type Message struct {
	Type    int         `json:"type"`
	Content interface{} `json:"content"`
}

var (
	subscribe   = make(chan Client, 10)
	unsubscribe = make(chan Client, 10)
	message     = make(chan Message, 10)
	lastNodes   []comm.TopologyNode
	lastLines   []comm.TopologyLine
	newAttack   bool
	NodeSum     int
)

type Client struct {
	Name string
	Conn *websocket.Conn
}

func attackMonitJob() {

	for {
		time.Sleep(time.Second * 10)
		refreshTopology()
	}

}

/*
-- 节点
	1. 10分钟内的透明代理查到攻击IP  为攻击节点
	2. 所有部署的透明代理 为透明代理节点  以及 透明代理和协议代理的Green连线信息
	3. 所有部署的协议代理 为协议代理节点 以及 透明代理和蜜罐的连线信息
   4. 所有部署的蜜罐 为蜜罐节点
-- 连接
	1. 透明代理上报的日志 Red连接 攻击节点和透明代理节点
   2. 协议代理日志联合透明代理日志 Red连接 透明代理和协议代理节点
   3  协议代理日志 查询 Red连接 协议代理和蜜罐的连接
*/

func refreshTopology() {

	attackPodMap := map[string]comm.TopologyLine{}

	lastNodes, lastLines = QueryTopology()

	potNum := 0

	if lastNodes != nil {
		for _, node := range lastNodes {
			if node.NodeType == "POD" {
				potNum = potNum + 1
			}
		}
	} else {
		lastNodes = make([]comm.TopologyNode, 0)
	}

	if lastLines != nil {
		for _, line := range lastLines {
			if line.Status == "RED" && strings.Index(line.Target, "POD") > -1 {
				attackPodMap[line.Target] = line
			}
		}
	} else {
		lastLines = make([]comm.TopologyLine, 0)
	}

	attackedPotNum := len(attackPodMap)

	var data = struct {
		Nodes          []comm.TopologyNode `json:"nodes"`
		Lines          []comm.TopologyLine `json:"lines"`
		PotNum         int                 `json:"potNum"`
		AttackedPotNum int                 `json:"attackedPotNum"`
	}{
		lastNodes,
		lastLines,
		potNum,
		attackedPotNum,
	}
	var msg Message
	msg.Type = 3
	msg.Content = data
	message <- msg
}

func topology() {
	clients := make(map[string]Client)
	for {
		select {
		case msg := <-message:
			for _, client := range clients {
				data, err := json.Marshal(msg)
				if err != nil {
					zap.L().Error(err.Error())
					return
				}
				if client.Conn.WriteMessage(websocket.TextMessage, data) != nil {
					zap.L().Info("Fail to write message")
				}
			}
		// 有Client加入
		case client := <-subscribe:
			clients[client.Name] = client
			var msg Message
			msg.Type = 1
			msg.Content = fmt.Sprintf("connect success")
			message <- msg
			refreshTopology()

		// 有Client退出
		case client := <-unsubscribe:
			if _, ok := clients[client.Name]; !ok {
				zap.L().Info("the client had leaved, client's name:" + client.Name)
				break
			}
			delete(clients, client.Name)
			var msg Message
			msg.Type = 2
			msg.Content = fmt.Sprintf("disconnect")
			message <- msg
		}
	}
}

func TopologyMapHandle(context *gin.Context) {
	// 检验http头中upgrader属性，若为websocket，则将http协议升级为websocket协议
	ws, err := WsUpGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		zap.L().Error("websocket upgrade err")
		return
	}
	defer ws.Close()
	if err != nil {
		zap.L().Error(err.Error())
	}
	if _, ok := err.(websocket.HandshakeError); ok {
		zap.L().Error("Not a websocket connection")
		return
	} else if err != nil {
		fmt.Println("Cannot setup WebSocket connection:", err)
		zap.L().Error(err.Error())
		return
	}

	var client Client
	client.Conn = ws
	client.Name = context.Request.Host
	subscribe <- client

	// 当函数返回时，将该用户加入退出通道，并断开用户连接
	defer func() {
		unsubscribe <- client
		client.Conn.Close()
	}()

	// 由于WebSocket一旦连接，便可以保持长时间通讯，则该接口函数可以一直运行下去，直到连接断开
	for {
		// 读取消息。如果连接断开，则会返回错误
		_, msgStr, err := client.Conn.ReadMessage()

		// 如果返回错误，就退出循环
		if err != nil {
			break
		}
		// 如果没有错误，则把用户发送的信息放入message通道中
		var msg Message
		msg.Type = 0
		msg.Content = string(msgStr)
		message <- msg
	}
	return
}

func QueryTopology() ([]comm.TopologyNode, []comm.TopologyLine) {
	nodes, lines := QueryTopologyMap()

	if nodes == nil {
		return nodes, lines
	}

	return nodes, lines
}

/*
	1. 10分钟内的透明代理查到攻击IP  为攻击节点
	2. 所有部署的透明代理 为透明代理节点  以及 透明代理和协议代理的Green连线信息
	3. 所有部署的协议代理 为协议代理节点 以及 透明代理和蜜罐的连线信息
    4. 所有部署的蜜罐 为蜜罐节点
*/
func QueryTopologyMap() ([]comm.TopologyNode, []comm.TopologyLine) {

	var topologyNodes []comm.TopologyNode
	var protocolAttackNodes []comm.TopologyNode
	var topologyNodeLines = map[string]comm.TopologyLine{}
	var topologyNodeLineArray []comm.TopologyLine
	var topologyNodeGreenLines = map[string]comm.TopologyLine{}
	var topologyNodeRedLines = map[string]comm.TopologyLine{}
	var protocolProxies models.ProtocolProxy
	var honeypots models.Honeypot
	var transparentEvents models.TransparentEvent

	protocolProxyNodes, err := protocolProxies.QueryProtocolProxyNode()

	if err == nil {
		for _, protocolProxyNode := range protocolProxyNodes {
			topologyNodes = append(topologyNodes, protocolProxyNode)
		}
	}

	honeypotNodes, err := honeypots.GetHoneypotNodes()

	if err == nil {
		for _, honeypotNode := range honeypotNodes {
			topologyNodes = append(topologyNodes, honeypotNode)
		}
	}

	// 10分钟之内有攻击的节点 才算入
	attackNodesForAgent, err := transparentEvents.GetTransparentEventNodes()

	if err == nil {
		for _, attackNode := range attackNodesForAgent {
			topologyNodes = append(topologyNodes, attackNode)
		}
	}

	attackedAgentNodes, err := transparentEvents.GetAttackedAgentNodes()

	if err == nil && attackNodesForAgent != nil && len(attackNodesForAgent) > 0 {
		for _, attackNode := range attackedAgentNodes {
			topologyNodes = append(topologyNodes, attackNode)
		}
	}

	protocolAttackOriginNodes, err := protocolProxies.QueryProtocolAttackNode()

	// 去除为透明代理的节点
	if err == nil {
		for _, attackNode := range protocolAttackOriginNodes {
			if attackedAgentNodes != nil {
				for _, transparentProxyNode := range attackedAgentNodes {
					if !strings.Contains(transparentProxyNode.Ip, attackNode.Ip) {
						topologyNodes = append(topologyNodes, attackNode)
						protocolAttackNodes = append(protocolAttackNodes, attackNode)
					}
				}
			} else {
				topologyNodes = append(topologyNodes, attackNode)
				protocolAttackNodes = append(protocolAttackNodes, attackNode)
			}
		}
	}
	// 查询攻击者 至 Agent 节点之间的 连线
	buildAttack2AgentRedLines(attackNodesForAgent, attackedAgentNodes, topologyNodeRedLines)
	// 查询攻击者 至 协议代理 节点之间的 连线
	buildAttack2ProtocolRedLines(protocolAttackNodes, topologyNodeRedLines)
	// 查询被攻击的Agent 至 协议代理 节点之间的 连线
	QueryAttackedAgent2ProtocolRedLines(attackedAgentNodes, protocolProxyNodes, topologyNodeRedLines)
	// 查询所有协议代理 至 蜜罐 节点之间的 连线
	QueryProtocol2PodGreenLines(protocolProxyNodes, honeypotNodes, topologyNodeGreenLines)
	// 查询被攻击的协议代理 至 蜜罐 节点之间的 连线
	QueryProtocol2PodRedLines(topologyNodeRedLines)

	for key, line := range topologyNodeGreenLines {
		topologyNodeLines[key] = line
	}

	for key, line := range topologyNodeRedLines {
		topologyNodeLines[key] = line
	}

	for _, line := range topologyNodeLines {
		topologyNodeLineArray = append(topologyNodeLineArray, line)
	}

	return topologyNodes, topologyNodeLineArray
}

func buildIpParam(nodes []comm.TopologyNode) string {

	var ipParams string
	isFirst := true
	for index, node := range nodes {
		if index > 0 && !isFirst && node.Ip != "" {
			ipParams = ipParams + ",'" + node.Ip + "'"
		} else {
			ipParams = "'" + node.Ip + "'"
			isFirst = false
		}
	}

	return ipParams
}

func buildSplitIpParam(nodes []comm.TopologyNode) string {
	var ipList []string
	for _, node := range nodes {
		if strings.Contains(node.Ip, ",") {
			ips := strings.Split(node.Ip, ",")
			for _, ip := range ips {
				ipList = append(ipList, ip)
			}
		} else {
			ipList = append(ipList, node.Ip)
		}
	}
	var ipParams string
	isFirst := true
	for index, ip := range ipList {
		if index > 0 && !isFirst && ip != "" {
			ipParams = ipParams + ",'" + ip + "'"
		} else {
			ipParams = "'" + ip + "'"
			isFirst = false
		}
	}
	return ipParams
}

/*
	1. 查attack 节点到 透明代理的线 为红线
	2. 查透明代理到协议代理 策略表 为 绿线
	3. 查协议代理到 蜜罐的 策略表 为绿线
    4. 查协透明代理的连接事件 为 透明代理到协议代理的红线
    5. 查协议代理攻击日志表 为协议代理到蜜罐的红线
*/

func buildAttack2AgentRedLines(attackNodes, agentNodes []comm.TopologyNode, redLines map[string]comm.TopologyLine) {
	var transparentEvent models.TransparentEvent
	attackIpParams := buildIpParam(attackNodes)
	agentIpParams := buildIpParam(agentNodes)
	lineArray, _ := transparentEvent.QueryAttack2AgentLines(attackIpParams, agentIpParams)
	if len(lineArray) == 0 {
		return
	}
	for _, line := range lineArray {
		lineKey := fmt.Sprintf("%s-%s", line.Source, line.Target)
		redLines[lineKey] = line
	}
}

func buildAttack2ProtocolRedLines(protocolAttackNodes []comm.TopologyNode, redLines map[string]comm.TopologyLine) {
	var protocolEvents models.ProtocolEvent
	attackIpParams := buildIpParam(protocolAttackNodes)
	lineArray, _ := protocolEvents.QueryAttack2ProtocolLines(attackIpParams)
	if len(lineArray) == 0 {
		return
	}
	for _, line := range lineArray {
		lineKey := fmt.Sprintf("%s-%s", line.Source, line.Target)
		redLines[lineKey] = line
	}
}

func QueryAttackedAgent2ProtocolRedLines(attackedAgentNodes, protocolNodes []comm.TopologyNode, redLines map[string]comm.TopologyLine) {
	var protocolEvents models.ProtocolEvent

	if attackedAgentNodes == nil || len(attackedAgentNodes) == 0 {
		return
	}

	agentIpParams := buildSplitIpParam(attackedAgentNodes)

	lineArray, _ := protocolEvents.QueryAttackedAgent2ProtocolRedLines(agentIpParams)

	if len(lineArray) == 0 {
		return
	}
	for _, line := range lineArray {
		lineKey := fmt.Sprintf("%s-%s", line.Source, line.Target)
		redLines[lineKey] = line
	}
}

func QueryProtocol2PodGreenLines(protocolNodes, podNodes []comm.TopologyNode, greenLines map[string]comm.TopologyLine) {
	var protocolProxy models.ProtocolProxy
	lineArray, _ := protocolProxy.QueryProtocol2PodGreenLines()
	if len(lineArray) == 0 {
		return
	}
	for _, line := range lineArray {
		lineKey := fmt.Sprintf("%s-%s", line.Source, line.Target)
		greenLines[lineKey] = line
	}
}

func QueryProtocol2PodRedLines(redLines map[string]comm.TopologyLine) {
	var protocolEvent models.ProtocolEvent
	lineArray, _ := protocolEvent.QueryProtocol2PodRedLines()
	if len(lineArray) == 0 {
		return
	}
	for _, line := range lineArray {
		if line.Target != "-POD" {
			lineKey := fmt.Sprintf("%s-%s", line.Source, line.Target)
			redLines[lineKey] = line
		}
	}
}
