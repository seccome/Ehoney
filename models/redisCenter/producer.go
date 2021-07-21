package redisCenter

import (
	"decept-defense/models/honeycluster"
	"decept-defense/models/policyCenter"
	"decept-defense/models/util"
	"decept-defense/models/util/comm"
	"github.com/astaxie/beego"
	"strconv"
	"strings"

	//"decept-defense/models/util"
	//"decept-defense/models/util/comm"
	"fmt"
	//"github.com/go-redis/redis"
	"github.com/garyburd/redigo/redis"
	"github.com/tidwall/gjson"
	_ "github.com/tidwall/gjson"

	"time"
)

type redisServer struct {
	redisHost string
	redisAuth string
}

var redisurl = beego.AppConfig.String("redisurl")
var redisport = beego.AppConfig.String("redisport")
var redispwd = beego.AppConfig.String("redispwd")

//NewRedis ...
func NewRedis(redisHost, redisAuth string) *redisServer {
	return &redisServer{redisHost, redisAuth}
}

func (rs *redisServer) NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {

			conn, err := redis.Dial("tcp", rs.redisHost)
			if err != nil {
				return nil, err
			}
			if _, err := conn.Do("AUTH", rs.redisAuth); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func (rs *redisServer) BaitListen(pool *redis.Pool, key string) {
	conn := pool.Get()
	defer conn.Close()
	psc := redis.PubSubConn{conn}
	psc.Subscribe(key)
	for {
		time.Now().Unix()
		switch v := psc.Receive().(type) {
		//有消息Publish到指定Channel时
		case redis.Message:
			data := string(v.Data)
			data = util.Base64Decode(data)
			status := gjson.Get(data, "Status")
			tasktype := gjson.Get(data, "Type")
			taskid := gjson.Get(data, "TaskId")
			if status.Raw == "1" {
				if tasktype.Str == comm.BaitHis || tasktype.Str == comm.BaitFile {
					status := 1
					policyCenter.UpdateBaitPolicyStatusByTaskid(status, taskid.Str)
				} else if tasktype.Str == comm.BaitUNFile {
					status := 5
					offlinetime := time.Now().Unix()
					policyCenter.UpdateBaitPolicy(taskid.Str, offlinetime, status)
				} else if tasktype.Str == comm.SignFile {
					status := 1
					policyCenter.UpdateSignPolicy(status, taskid.Str)
				} else if tasktype.Str == comm.SignUNFile {
					offlinetime := time.Now().Unix()
					status := 5
					policyCenter.OffSignPolicy(taskid.Str, status, offlinetime)
				}
			}
			if status.Raw == "-1" {
				if tasktype.Str == comm.BaitHis || tasktype.Str == comm.BaitFile {
					status := 2
					policyCenter.UpdateBaitPolicyStatusByTaskid(status, taskid.Str)
				} else if tasktype.Str == comm.BaitUNFile {
					status := 6
					offlinetime := time.Now().Unix()
					policyCenter.UpdateBaitPolicy(taskid.Str, offlinetime, status)
				} else if tasktype.Str == comm.SignFile {
					status := 2
					policyCenter.UpdateSignPolicy(status, taskid.Str)
				} else if tasktype.Str == comm.SignUNFile {
					offlinetime := time.Now().Unix()
					status := 6
					policyCenter.OffSignPolicy(taskid.Str, status, offlinetime)
				}
			}
		case redis.Subscription: //Subscribe一个Channel时
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			//RedisPubConsumerBaitPolicyResponse()
			return
		}
	}
}

func (rs *redisServer) TransListen(pool *redis.Pool, key string) {
	conn := pool.Get()
	defer conn.Close()
	psc := redis.PubSubConn{conn}
	psc.Subscribe(key)
	for {
		time.Now().Unix()
		switch v := psc.Receive().(type) {
		case redis.Message: //有消息Publish到指定Channel时
			data := util.Base64Decode(string(v.Data))
			status := gjson.Get(data, "Status")
			taskid := gjson.Get(data, "TaskId")
			tasktype := gjson.Get(data, "Type")
			if status.Num == 1 {
				if tasktype.Str == comm.AgentTypeEdge {
					policyCenter.UpdateTransPolicyStatusByTaskid(1, taskid.Str)
				} else if tasktype.Str == comm.AgentTypeUNEdge {
					offlinetime := time.Now().Unix()
					policyCenter.UpdateTransPolicy(taskid.Str, offlinetime, 5)
				} else if tasktype.Str == comm.AgentTypeRelay {
					policyCenter.UpdateHoneyTransPolicyStatusByTaskid(1, taskid.Str)
				} else if tasktype.Str == comm.AgentTypeUNRelay {
					offlinetime := time.Now().Unix()
					policyCenter.UpdateHoneyTransPolicy(taskid.Str, offlinetime, 5)
				}
			} else if status.Num == -1 {
				if tasktype.Str == comm.AgentTypeEdge {
					policyCenter.UpdateTransPolicyStatusByTaskid(2, taskid.Str)
				} else if tasktype.Str == comm.AgentTypeUNEdge {
					offlinetime := time.Now().Unix()
					policyCenter.UpdateTransPolicy(taskid.Str, offlinetime, 6)
				} else if tasktype.Str == comm.AgentTypeRelay {
					policyCenter.UpdateHoneyTransPolicyStatusByTaskid(2, taskid.Str)
				} else if tasktype.Str == comm.AgentTypeUNRelay {
					offlinetime := time.Now().Unix()
					policyCenter.UpdateHoneyTransPolicy(taskid.Str, offlinetime, 6)
				}
			}
		case redis.Subscription: //Subscribe一个Channel时
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			//RedisPubConsumerTransPolicyResponse()
			return
		}
	}
}

func (rs *redisServer) SendMsg(pool *redis.Pool, key, value string) {
	//i := 0
	conn := pool.Get()
	defer conn.Close()
	//for {
	fmt.Println("send:", value)
	value = util.Base64Encode(value)
	//logs.Error("redis value:" + value)
	pubtext := value
	res, err := redis.Values(conn.Do("publish", key, pubtext))
	if err != nil {

	}
	fmt.Println(res)
	//time.Sleep(60000000 * time.Microsecond)
	//i++

	//}
}

func (rs *redisServer) ServerHeartBeatListen(pool *redis.Pool, key string) {
	conn := pool.Get()
	defer conn.Close()
	psc := redis.PubSubConn{conn}
	psc.Subscribe(key)
	for {
		time.Now().Unix()
		switch v := psc.Receive().(type) {
		case redis.Message:
			data := util.Base64Decode(string(v.Data))

			agentid := gjson.Get(data, "AgentId")
			status := gjson.Get(data, "Status")
			ips := gjson.Get(data, "IPs")
			servername := gjson.Get(data, "HostName")
			servertype := gjson.Get(data, "Type")
			sys := gjson.Get(data, "Sys")
			var system string
			if sys.Str == "" {
				system = "Linux"
			} else {
				system = sys.Str
			}
			timenow := time.Now().Unix()
			status1 := 0
			if status.Str == "running" {
				status1 = 1
			}
			if strings.ToLower(servertype.Str) == "edge" {
				//fmt.Println(fmt.Println("EDGE data：", data))
				honeycluster.ServerHeartBeatAct(agentid.Str, system, status1, ips.Str, servername.Str, timenow)
			} else if strings.ToLower(servertype.Str) == "relay" {
				//fmt.Println(fmt.Println("RELAY data：", data))
				honeycluster.HoneyServerHeartBeatAct(agentid.Str, status1, ips.Str, servername.Str, timenow)
			}
		case redis.Subscription: //Subscribe一个Channel时
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			//RedisPubConsumerServerRegResponse()
			return
		}
	}
}

func (rs *redisServer) TransEventListen(pool *redis.Pool, key string) {
	conn := pool.Get()
	defer conn.Close()
	psc := redis.PubSubConn{conn}
	psc.Subscribe(key)
	for {
		time.Now().Unix()
		switch v := psc.Receive().(type) {
		case redis.Message:
			data := util.Base64Decode(string(v.Data))
			fmt.Println(data)

			agentid := gjson.Get(data, "AgentId")
			sourceaddr := gjson.Get(data, "SourceAddr")
			bindport := gjson.Get(data, "BindPort")
			destaddr := gjson.Get(data, "DestAddr")
			destport := gjson.Get(data, "DestPort")
			eventtime := gjson.Get(data, "EventTime")
			proxytype := gjson.Get(data, "ProxyType")
			exportport := gjson.Get(data, "ExportPort")
			serverip := ""
			honeypotid := ""
			honeypotport := ""
			honeytypeid := ""
			serverinfo := honeycluster.SelectApplicationByAgentID(agentid.Str)
			if len(serverinfo) > 0 {
				serverip = util.Strval(serverinfo[0]["serverip"])
			}
			dport, _ := strconv.Atoi(destport.Raw)
			honeyinfo := honeycluster.SelectHoneyInfoByTransInfo(destaddr.Str, dport)
			if len(honeyinfo) > 0 {
				honeypotid = util.Strval(honeyinfo[0]["honeypotid"])
				honeypotport = util.Strval(honeyinfo[0]["honeypotport"])
				honeytypeid = util.Strval(honeyinfo[0]["honeytypeid"])
			}
			hport, _ := strconv.Atoi(honeypotport)
			bport, _ := strconv.Atoi(bindport.Raw)
			eport, _ := strconv.Atoi(exportport.Raw)
			policyCenter.InsertAttackLog(proxytype.Str, serverip, bport, util.GetIp2(sourceaddr.Str), honeypotid, hport, honeytypeid, eventtime.Int(), eport)

		case redis.Subscription: //Subscribe一个Channel时
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			//RedisPubConsumerTransEventResponse()
			return
		}
	}
}

//func WritePolicyToRedis(value string)  {
//	// 建立连接
//	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
//	if err != nil {
//		fmt.Println("redis.Dial err=", err)
//		return
//	}
//	// 通过go向redis写入数据
//	_, err = conn.Do("Set", "agentId", value)
//	if err != nil {
//		fmt.Println("set err=", err)
//		return
//	}
//	// 关闭连接
//	defer conn.Close()
//	// 读取数据 获取名字
//	r, err := redis.String(conn.Do("Get", "agentId"))
//	if err != nil {
//		fmt.Println("set err=", err)
//		return
//	}
//	fmt.Println("Manipulate success, the name is", r)
//
//}

func RedisSubProducerBaitPolicy(msg string) {

	redis := NewRedis(redisurl+":"+redisport, redispwd)
	pool := redis.NewPool()
	//logs.Error("RedisSubProducerBaitPolicy:", msg)
	go redis.SendMsg(pool, "deception-strategy-sub-channel", msg)
	//select {}
}

func RedisSubProducerTransPolicy(msg string) {

	redis := NewRedis(redisurl+":"+redisport, redispwd)
	pool := redis.NewPool()
	go redis.SendMsg(pool, "proxy-strategy-sub-channel", msg)
	select {}
}

func RedisPubConsumerBaitPolicyResponse() {
	redis := NewRedis(redisurl+":"+redisport, redispwd)
	pool := redis.NewPool()
	go redis.BaitListen(pool, "deception-strategy-pub-channel")
	select {}
}

func RedisPubConsumerTransPolicyResponse() {
	redis := NewRedis(redisurl+":"+redisport, redispwd)
	pool := redis.NewPool()
	go redis.TransListen(pool, "proxy-strategy-pub-channel")
	select {}
}

func RedisPubConsumerServerRegResponse() {
	redis := NewRedis(redisurl+":"+redisport, redispwd)
	pool := redis.NewPool()
	go redis.ServerHeartBeatListen(pool, "agent-heart-beat-channel")
	select {}
}

func RedisPubConsumerTransEventResponse() {
	redis := NewRedis(redisurl+":"+redisport, redispwd)
	pool := redis.NewPool()
	go redis.TransEventListen(pool, "proxy-event-pub-channel")
	select {}
}
