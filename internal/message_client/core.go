package message_client

import (
	"context"
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

// TODO change channel name
//   1.agent-heart-beat-channel to server_heartbeat
//   2.agent-task-response-channel
//   3.agent-task-request-channel
var (
	ctx         = context.Background()
	RedisClient *redis.Client
	subChannels = []string{
		"agent-heart-beat-channel",
		"agent-task-response-channel",
	}
)

func Setup() {
	if err := createRedisClient(); err != nil {
		zap.L().Fatal("redis 连接异常: " + err.Error())
	}
	go subscribeChannel()
}

func createRedisClient() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     configs.GetSetting().Redis.RedisHost + ":" + strconv.Itoa(configs.GetSetting().Redis.RedisPort),
		Password: configs.GetSetting().Redis.RedisPassword,
		DB:       0,
	})
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func AppendList(key string, value string) error {
	err := RedisClient.RPush(ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetList(key string) ([]string, error) {
	len, _ := RedisClient.LLen(ctx, key).Result()
	ret, _ := RedisClient.LRange(ctx, key, 0, len).Result()
	return ret, nil
}

func PublishMessage(channel string, msg string) error {
	var err error
	channel = "agent-task-request-channel"
	zap.L().Debug(fmt.Sprintf("Publish message %s to channel %s", msg, channel))
	fmt.Printf("Publish message %s to channel %s", msg, channel)
	fmt.Printf("Publish message %s to channel %s", util.Base64Encode(msg), channel)

	err = RedisClient.Publish(ctx, channel, util.Base64Encode(msg)).Err()
	if err != nil {
		fmt.Printf("try publish message to channel %s error[%s]\n", channel, err.Error())
		return err
	}
	return nil
}

func subscribeChannel() {
	pubSub := RedisClient.Subscribe(ctx, "agent-heart-beat-channel", "agent-task-response-channel")
	for msg := range pubSub.Channel() {
		decoded := util.Base64Decode(msg.Payload)
		taskType := gjson.Get(decoded, "TaskType").String()
		operationType := gjson.Get(decoded, "OperatorType").String()
		agentID := gjson.Get(decoded, "AgentId").String()
		status := int(gjson.Get(decoded, "Status").Int())
		taskID := gjson.Get(decoded, "TaskID").String()
		attackType := gjson.Get(decoded, "AttackType").String()
		if taskType == string(comm.TOKEN) {
			if operationType == (string(comm.DEPLOY)) {
				var token models.ProbeToken
				token.UpdateProbeTokenStatusByTaskID(status, taskID)
			}
			continue

		} else if taskType == string(comm.BAIT) {
			if operationType == (string(comm.DEPLOY)) {
				var bait models.ProbeBaits
				bait.UpdateProbeBaitStatusByTaskID(status, taskID)
			}
			continue
		} else if taskType == string(comm.Heartbeat) {
			ips := gjson.Get(decoded, "IPs").String()
			hostName := gjson.Get(decoded, "HostName").String()
			serverType := strings.ToLower(gjson.Get(decoded, "Type").String())
			if strings.ToLower(serverType) == "edge" {
				system := gjson.Get(decoded, "System").String()
				var server models.Probes
				server.ServerIP = ips
				server.HostName = hostName
				server.HeartbeatTime = util.GetCurrentTime()
				server.CreateTime = util.GetCurrentTime()
				server.AgentID = agentID
				server.SystemType = system
				if err := server.CreateServer(); err != nil {
					zap.L().Error("heartbeat create probe error: " + err.Error())
				}
				continue
			} else if strings.ToLower(serverType) == "relay" {
				var server models.HoneypotServers
				server.ServerIP = ips
				server.HostName = hostName
				server.HeartbeatTime = util.GetCurrentTime()
				server.CreateTime = util.GetCurrentTime()
				server.AgentID = agentID
				if err := server.CreateServer(); err != nil {
					zap.L().Error("heartbeat create honeypotServer error: " + err.Error())
				}
				continue
			}
		} else if taskType == string(comm.PROTOCOL) {
			if operationType == (string(comm.DEPLOY)) {
				var protocol models.Protocols
				protocol.UpdateProtocolStatusByTaskID(status, taskID)
			}
			continue
		} else if taskType == string(comm.TransparentProxy) {
			if operationType == (string(comm.DEPLOY)) {
				var transparentProxy models.TransparentProxy
				if status == 3 {
					transparentProxy.UpdateTransparentProxyStatusByTaskID(comm.SUCCESS, taskID)
				} else {
					transparentProxy.UpdateTransparentProxyStatusByTaskID(comm.FAILED, taskID)
				}
			} else if operationType == (string(comm.WITHDRAW)) {
				var transparentProxy models.TransparentProxy
				if status == 3 {
					transparentProxy.UpdateTransparentProxyStatusByTaskID(comm.FAILED, taskID)
				} else {
					transparentProxy.UpdateTransparentProxyStatusByTaskID(comm.SUCCESS, taskID)
				}
			}
			continue
		} else if taskType == string(comm.ProtocolProxy) {
			if operationType == (string(comm.DEPLOY)) {
				var protocolProxy models.ProtocolProxy
				if status == 3 {
					protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.SUCCESS, taskID)
				} else {
					protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.FAILED, taskID)
				}
			} else if operationType == (string(comm.WITHDRAW)) {
				var protocolProxy models.ProtocolProxy
				if status == 3 {
					protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.FAILED, taskID)
				} else {
					protocolProxy.UpdateProtocolProxyStatusByTaskID(comm.SUCCESS, taskID)
				}
			}
			continue
		} else if taskType != "" {
			zap.L().Debug("不支持的task类型" + taskType)
		}

		if attackType == string(comm.ProtocolAttackEvent) {
			//TODO get attack event by http
		} else if attackType == string(comm.TransparentAttackEvent) {
			AttackType := gjson.Get(decoded, "AttackType").String()
			AgentID := gjson.Get(decoded, "AgentID").String()
			AttackIP := gjson.Get(decoded, "AttackIP").String()
			AttackPort := int32(gjson.Get(decoded, "AttackPort").Int())
			ProxyIP := gjson.Get(decoded, "ProxyIP").String()
			ProxyPort := int32(gjson.Get(decoded, "ProxyPort").Int())
			Transparent2ProtocolPort := int32(gjson.Get(decoded, "Transparent2ProtocolPort").Int())
			DestIP := gjson.Get(decoded, "DestIP").String()
			DestPort := int32(gjson.Get(decoded, "DestPort").Int())
			EventTime := gjson.Get(decoded, "EventTime").String()
			if util.IsLocalIP(AttackIP) {
				continue
			}
			var transparentAttack = models.TransparentEvent{
				AttackType:               AttackType,
				AgentID:                  AgentID,
				AttackIP:                 AttackIP,
				AttackPort:               AttackPort,
				ProxyIP:                  ProxyIP,
				ProxyPort:                ProxyPort,
				Transparent2ProtocolPort: Transparent2ProtocolPort,
				DestIP:                   DestIP,
				DestPort:                 DestPort,
				EventTime:                EventTime,
			}
			if err := transparentAttack.CreateEvent(); err != nil {
				zap.L().Error("透明代理事件创建异常")
				continue
			}
			waning := `攻击类型: ` + transparentAttack.AttackType + `\n\n > AgentID:  ` + transparentAttack.AgentID + `\n\n > 攻击IP:  ` + transparentAttack.AttackIP + `\n\n > 攻击端口:  ` + strconv.Itoa(int(transparentAttack.AttackPort)) + `\n\n > 代理IP:  ` + transparentAttack.ProxyIP + `\n\n > 代理端口:  ` + strconv.Itoa(int(transparentAttack.ProxyPort)) + `\n\n > 蜜网IP:  ` + transparentAttack.DestIP + `\n\n > 蜜网端口:  ` + strconv.Itoa(int(transparentAttack.DestPort)) + `\n\n > 连接端口:  ` + strconv.Itoa(int(transparentAttack.Transparent2ProtocolPort)) + `\n\n > 创建时间:  ` + transparentAttack.EventTime + ``
			util.SendDingMsg("欺骗防御告警", "透明代理告警", waning)
			continue

		} else {
			zap.L().Debug("不支持的攻击类型类型")
		}
	}
}
