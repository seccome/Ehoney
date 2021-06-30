package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

type WebSocket struct {
	beego.Controller
}

var (
	Dbhost        = beego.AppConfig.String("dbhost")
	Dbport        = beego.AppConfig.String("dbport")
	Dbuser        = beego.AppConfig.String("dbuser")
	Dbpassword    = beego.AppConfig.String("dbpassword")
	Dbname        = beego.AppConfig.String("dbname")
	LastAttackLog AttackLog
	LastNodes     []TopologyNode
	LastLines     []TopologyLine
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

func init() {
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
	lastNodes   []TopologyNode
	lastLines   []TopologyLine
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

func refreshTopology() {

	attackPodMap := map[string]TopologyLine{}

	lastNodes, lastLines, newAttack = QueryTopology()

	potNum := 0

	if lastNodes != nil {
		for _, node := range lastNodes {
			if node.NodeType == "POT" {
				potNum = potNum + 1
			}
		}
	} else {
		lastNodes = make([]TopologyNode, 0)
	}

	if lastLines != nil {
		for _, line := range lastLines {
			if line.Status == "RED" && strings.Index(line.Target, "POT") > -1 {
				attackPodMap[line.Target] = line
			}
		}
	} else {
		lastLines = make([]TopologyLine, 0)
	}

	attackedPotNum := len(attackPodMap)

	var data = struct {
		Nodes          []TopologyNode `json:"nodes"`
		Lines          []TopologyLine `json:"lines"`
		PotNum         int            `json:"potNum"`
		AttackedPotNum int            `json:"attackedPotNum"`
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
		// 哪个case可以执行，则转入到该case。若都不可执行，则堵塞。
		select {
		// 消息通道中有消息则执行，否则堵塞
		case msg := <-message:
			for _, client := range clients {
				data, err := json.Marshal(msg)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(data)
				if client.Conn.WriteMessage(websocket.TextMessage, data) != nil {
					fmt.Println("Fail to write message")
				}
			}
		// 有用户加入
		case client := <-subscribe:
			clients[client.Name] = client // 将用户加入映射
			// 将用户加入消息放入消息通道
			var msg Message
			msg.Type = 1
			msg.Content = fmt.Sprintf("connect success")
			// 此处要设置有缓冲的通道。因为这是goroutine自己从通道中发送并接受数据。
			// 若是无缓冲的通道，该goroutine发送数据到通道后就被锁定，需要数据被接受后才能解锁，而恰恰接受数据的又只能是它自己
			message <- msg
			refreshTopology()

		// 有用户退出
		case client := <-unsubscribe:
			str := fmt.Sprintf("broadcaster-----------%s leave\n", client.Conn)
			fmt.Println(str)

			// 如果该用户已经被删除

			if _, ok := clients[client.Name]; !ok {
				beego.Info("the client had leaved, client's name:" + client.Name)
				break
			}

			delete(clients, client.Name) // 将用户从映射中删除

			// 将用户退出消息放入消息通道
			var msg Message
			msg.Type = 2
			msg.Content = fmt.Sprintf("disconnect")
			message <- msg
		}
	}
}

func Map(w http.ResponseWriter, r *http.Request) {
	// 检验http头中upgrader属性，若为websocket，则将http协议升级为websocket协议
	// conn, err := wsUpGrader.Upgrade(w, r, nil)

	conn, err := WsUpGrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	if _, ok := err.(websocket.HandshakeError); ok {
		fmt.Println("Not a websocket connection")
		return
	} else if err != nil {
		fmt.Println("Cannot setup WebSocket connection:", err)
		return
	}

	var client Client
	client.Conn = conn
	client.Name = r.Host
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
}

// 查询被攻击pod 的完整路径
func QueryTopology() ([]TopologyNode, []TopologyLine, bool) {
	var lineMap map[string]TopologyLine
	var lineArray []TopologyLine
	var newNodes []TopologyNode

	nodes := QueryTopologyNodes()

	if nodes == nil {
		return newNodes, lineArray, NewAttack
	}

	attackLog := QueryLatestAttackLog()

	NewAttack = true
	if attackLog.attackip != "" && (LastAttackLog.attacktime == 0 || LastAttackLog.attacktime < attackLog.attacktime-5 || NodeSum != len(nodes)) { // 第一次或 比上次日志延迟至少5毫秒之后
		LastAttackLog = attackLog
		mark := false
		var tmpNodes []TopologyNode

		for _, node := range nodes {
			if node.Ip == attackLog.attackip || strings.Index(node.Ip, attackLog.attackip) > -1 {
				node.NodeType = "HACK"
				mark = true
			}
			tmpNodes = append(tmpNodes, node)
		}

		nodes = tmpNodes

		if !mark {
			hackNode := TopologyNode{
				Id:       fmt.Sprintf("%s-%s", attackLog.attackip, "HACK"),
				Ip:       attackLog.attackip,
				HostName: attackLog.attackip,
				NodeType: "HACK",
			}
			nodes = append(nodes, hackNode)
		}
	} else {
		NewAttack = false
	}

	lineMap = QueryTopologyLines(attackLog, nodes, "GREEN")

	if NewAttack || attackLog.attackip == "" {
		dyeingRedLines(nodes, attackLog, lineMap, "RELAY")
	} else {
		return LastNodes, LastLines, NewAttack
	}

	// 查找Attack线
	var attackLine TopologyLine
	for _, line1 := range lineMap {
		if line1.Status == "RED" && strings.Index(line1.Source, "HACK") > -1 {
			attackLine = line1
			log.Printf("find attack line %v\n", attackLine)
		}
	}

	for _, line := range lineMap {

		if strings.Index(line.Source, "-") == -1 && strings.Index(line.Target, "HACK") > -1 {
			continue
		}

		if line.Source == line.Target {
			continue
		}

		lineArray = append(lineArray, line)
	}

	for _, node := range nodes {
		isIsolated := true
		for _, line := range lineArray {
			if line.Source == node.Id || line.Target == node.Id {
				isIsolated = false
			}
		}
		if !isIsolated {
			newNodes = append(newNodes, node)
		}

	}
	LastNodes = newNodes
	LastLines = lineArray
	NodeSum = len(nodes)
	return newNodes, lineArray, NewAttack
}

func dyeingRedLines(nodes []TopologyNode, attackLog AttackLog, lineMap map[string]TopologyLine, scanType string) {

	if scanType == "RELAY" {
		for _, node := range nodes {
			if node.NodeType != "RELAY" {
				continue
			}
			// 有连线的情况
			for key, line := range lineMap {
				attackPodId := fmt.Sprintf("%s-%s", attackLog.attackPotIp, "POT")
				if line.Source == node.Id && line.Target == attackPodId {
					line.Status = "RED"
					lineMap[key] = line
				}
			}
		}
		dyeingRedLines(nodes, attackLog, lineMap, "EDGE")
	} else if scanType == "EDGE" {
		// 有连线的情况

		redLineMap := map[string]TopologyLine{}

		for key, line := range lineMap {
			if line.Status == "RED" {
				redLineMap[key] = line
			}
		}

		for key, line := range lineMap {
			headMatch := false
			tailMatch := false
			for redKey, redLine := range redLineMap {
				if key == redKey {
					continue
				}

				if line.Source == redLine.Target {
					headMatch = true
				}
				if line.Target == redLine.Source {
					tailMatch = true
				}
			}

			if headMatch && tailMatch {
				line.Status = "RED"
				lineMap[key] = line
			}
		}
	}

}

func QueryLatestAttackLog() AttackLog {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		fmt.Errorf("[SelectTopAttackTypes]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var attackLog AttackLog
	defer sqlCon.Close()

	queryLatestAttackLogSql := "SELECT a.attackip, h.honeyip, a.attacktime, a.srchost as byIp FROM attacklog a LEFT JOIN honeypots h on a.honeypotid = h.honeypotid WHERE a.attackip IS NOT NULL AND a.honeypotid  IS NOT NULL AND a.proxytype IS NOT NULL AND a.proxytype != \"falco\" ORDER BY a.attacktime DESC LIMIT 1"
	rows, err := DbCon.Query(queryLatestAttackLogSql)
	if err != nil {
		fmt.Errorf("[SelectTopAttackTypes] select list error,%s", err)
	}

	if rows != nil {
		for rows.Next() {
			err = rows.Scan(&attackLog.attackip, &attackLog.attackPotIp, &attackLog.attacktime, &attackLog.byIp)
			if err != nil {
				fmt.Errorf("attackLog row init err %s", err)
			}
		}
	}

	if strings.Index(attackLog.byIp, ",") > -1 {
		ips := strings.Split(attackLog.byIp, ",")
		attackLog.byIp = ips[0]
	}

	return attackLog
}

func QueryTopologyNodes() []TopologyNode {
	sqlCon, err := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err != nil {
		fmt.Errorf("[SelectTopAttackTypes]open mysql fail %s", err)
	}
	DbCon := sqlCon

	defer sqlCon.Close()

	allHoneyPotSelectSql := "SELECT * FROM (SELECT concat(s.serverip, \"-EDGE\") as id, s.serverip as ip, s.servername as hostName, \"EDGE\" nodeType FROM servers s LEFT JOIN fowards f  ON f.agentid = s.agentid WHERE f.status != 2 UNION ALL SELECT concat( s.serverip, \"-RELAY\") as id,s.serverip AS ip, s.servername as hostName, \"RELAY\" nodeType FROM honeyfowards hf INNER JOIN honeypotservers s ON hf.agentid = s.agentid WHERE hf.status = 1 UNION ALL SELECT concat(honeyip, \"-POT\") as id, honeyip AS ip, honeyname as hostName, \"POT\" nodeType  FROM honeypots WHERE status = 1 AND honeytypeid !=\"\") as tmp GROUP BY id, ip, hostname, nodeType"
	rows, err := DbCon.Query(allHoneyPotSelectSql)

	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		fmt.Errorf("[SelectTopAttackTypes] select list error,%s", err)
	}
	var nodes []TopologyNode

	if rows != nil {
		for rows.Next() {
			var node TopologyNode
			err = rows.Scan(&node.Id, &node.Ip, &node.HostName, &node.NodeType)
			nodes = append(nodes, node)
		}
	}

	//for _, node := range nodes {
	//
	//
	//}

	return nodes
}

// 1. 从攻击目的蜜罐 进行溯源  产生新的 AttackNode 以及 路径Line
// 2. 从存量的Node 进行连线 产生 存量的 路径Line
func QueryTopologyLines(attackLog AttackLog, topologyNodes []TopologyNode, lineType string) map[string]TopologyLine {

	lineMap := map[string]TopologyLine{}

	for _, topologyNode := range topologyNodes {
		if topologyNode.Ip == attackLog.byIp || strings.Index(topologyNode.Ip, attackLog.byIp) > -1 {
			attackId := fmt.Sprintf("%s-%s", attackLog.attackip, "HACK")
			line := TopologyLine{
				Source: attackId,
				Target: topologyNode.Id,
				Status: "RED",
			}
			lineKey := fmt.Sprintf("%s-%s", attackId, topologyNode.Id)
			lineMap[lineKey] = line
		}
	}

	for _, topologyNode := range topologyNodes {
		findLinkedParentNode(topologyNode, lineMap, lineType)
	}

	return lineMap
}

/**
1. 蜜罐到协议代理 未有连线用虚线代替
2. 独立节点不包含
*/
func findLinkedParentNode(node TopologyNode, lineMap map[string]TopologyLine, lineType string) {

	if node.NodeType != "POT" && node.NodeType != "RELAY" {
		return
	}

	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		fmt.Errorf("[SelectTopAttackTypes]open mysql fail %s", err1)
	}
	defer sqlCon.Close()
	var parentNodes []TopologyNode
	var rows *sql.Rows
	var err error
	if node.NodeType == "POT" {
		findRelayNodesSql := fmt.Sprintf("SELECT concat(hps.serverip, \"-RELAY\") as id, hps.serverip AS ip, hps.servername as hostName, \"RELAY\" nodeType FROM honeyfowards hf INNER JOIN honeypotservers hps ON hf.agentid =  hps.agentid INNER JOIN honeypots hp ON hf.honeypotid = hp.honeypotid WHERE hf.status = 1 AND hp.honeyip = '%s'", node.Ip)
		rows, err = sqlCon.Query(findRelayNodesSql)
	} else if node.NodeType == "RELAY" {
		findEdgeNodesSql := fmt.Sprintf("SELECT concat(s.serverip, \"-EDGE\") as id, s.serverip,  s.servername as hostName, \"EDGE\" nodeType FROM fowards f  INNER JOIN servers s ON f.agentid = s.agentid WHERE f.status = 1 GROUP BY s.serverip,  s.servername")
		rows, err = sqlCon.Query(findEdgeNodesSql)
	} // 默认只连接到协议代理ip(协议代理只有一个)

	if rows == nil || err != nil {
		return
	}

	for rows.Next() {
		var parentNode TopologyNode
		err := rows.Scan(&parentNode.Id, &parentNode.Ip, &parentNode.HostName, &parentNode.NodeType)
		if err != nil {
			continue
		}

		line := TopologyLine{
			Source: parentNode.Id,
			Target: node.Id,
			Status: lineType,
		}
		parentNodes = append(parentNodes, parentNode)
		lineKey := fmt.Sprintf("%s-%s", parentNode.Id, node.Id)
		lineMap[lineKey] = line
	}

	if parentNodes != nil && len(parentNodes) > 0 {
		for _, parentNode := range parentNodes {
			findLinkedParentNode(parentNode, lineMap, lineType)
		}
	} else {
		if node.NodeType == "POT" {
			findRelayNodeSql := fmt.Sprintf(" SELECT concat(serverip, \"-RELAY\") as id, serverip, servername as hostName, \"POT\" nodeType FROM honeypotservers WHERE status = 1")
			rows, err = sqlCon.Query(findRelayNodeSql)
			if rows == nil || err != nil {
				return
			}
			for rows.Next() {
				var parentNode TopologyNode
				err := rows.Scan(&parentNode.Id, &parentNode.Ip, &parentNode.HostName, &parentNode.NodeType)
				if err != nil {
					continue
				}

				line := TopologyLine{
					Source: parentNode.Id,
					Target: node.Id,
					Status: "VIRTUAL",
				}
				parentNodes = append(parentNodes, parentNode)
				lineKey := fmt.Sprintf("%s-%s", parentNode.Id, node.Id)
				lineMap[lineKey] = line
			}
		}
	}
}

// TODO 需要确定AttackIp 和 哪个透明转发 进行关联

type TopologyNode struct {
	Id       string `json:"id"`
	Ip       string `json:"ip"`
	HostName string `json:"hostName"`
	NodeType string `json:"nodeType"`
}

type TopologyLine struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Status string `json:"status"`
}

type AttackLog struct {
	attackip    string
	attackPotIp string
	attacktime  int32
	byIp        string
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n+1:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
