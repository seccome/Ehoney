package agent_client

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

/**
1. 维护一张以AgentToken为Key 代理、诱饵组成的集合为Value的Map



*/

var (
	TransparentProxyMap = make(map[string]map[string]comm.AgentTransparentProxyTask)
	BaitMap             = make(map[string]map[string]comm.AgentBaitTask)
	LoadedTaskMap       = make(map[string]comm.TaskType)
)

func RegisterTransparentProxy(transparentProxy models.TransparentProxy, operatorType comm.OperatorType) error {
	transparentProxyTask := buildAgentTransparentProxyTask(transparentProxy, operatorType)
	if _, ok := TransparentProxyMap[transparentProxy.AgentToken]; !ok {
		TransparentProxyMap[transparentProxy.AgentToken] = make(map[string]comm.AgentTransparentProxyTask)
	}
	TransparentProxyMap[transparentProxyTask.AgentToken][transparentProxyTask.TaskId] = transparentProxyTask
	return nil
}

func buildAgentTransparentProxyTask(transparentProxy models.TransparentProxy, operatorType comm.OperatorType) comm.AgentTransparentProxyTask {
	var agentTransparentProxyTask comm.AgentTransparentProxyTask
	agentTransparentProxyTask.AgentToken = transparentProxy.AgentToken
	agentTransparentProxyTask.TaskId = transparentProxy.TransparentProxyId
	agentTransparentProxyTask.TaskType = comm.TRANSPARENT_PROXY_TASK
	agentTransparentProxyTask.OperatorType = operatorType
	agentTransparentProxyTask.ProxyPort = transparentProxy.ProxyPort
	agentTransparentProxyTask.DestIp = transparentProxy.DestIp
	agentTransparentProxyTask.DestPort = transparentProxy.DestPort
	return agentTransparentProxyTask
}

func RegisterBaitTask(baitTask models.BaitTask, operatorType comm.OperatorType) error {
	agentBaitTask := buildAgentBaitTask(baitTask, operatorType)

	if _, ok := BaitMap[agentBaitTask.AgentToken]; !ok {
		BaitMap[agentBaitTask.AgentToken] = make(map[string]comm.AgentBaitTask)
	}

	BaitMap[agentBaitTask.AgentToken][agentBaitTask.TaskId] = agentBaitTask
	return nil
}
func buildAgentBaitTask(baitTask models.BaitTask, operatorType comm.OperatorType) comm.AgentBaitTask {
	var agentBaitTask comm.AgentBaitTask
	agentBaitTask.AgentToken = baitTask.AgentToken
	agentBaitTask.TaskId = baitTask.BaitTaskId
	agentBaitTask.OperatorType = operatorType
	agentBaitTask.BaitTaskId = baitTask.BaitTaskId
	agentBaitTask.TaskType = comm.BAIT_TASK
	agentBaitTask.BaitType = baitTask.BaitType
	agentBaitTask.DeployPath = baitTask.DeployPath
	agentBaitTask.BaitName = baitTask.BaitName
	agentBaitTask.URL = baitTask.URL
	agentBaitTask.BaitData = baitTask.BaitData
	agentBaitTask.FileMD5 = baitTask.FileMD5
	agentBaitTask.ScriptName = baitTask.ScriptName
	agentBaitTask.CommandParameters = baitTask.CommandParameters
	return agentBaitTask
}

// FinishTask 如果状态正确
func FinishTask(taskId string) comm.TaskType {
	taskType, ok := LoadedTaskMap[taskId]
	if ok {
		zap.L().Info(fmt.Sprintf("Finish Task Id: {%s} Type: {%s}", taskId, taskType))
		delete(LoadedTaskMap, taskId)
	} else {
		zap.L().Info(fmt.Sprintf("Task Id: {%s} Unregistered", taskId))
	}
	return taskType
}

func HasTask(agentToken string) bool {
	fmt.Printf(string(len(BaitMap)))
	if _, ok := TransparentProxyMap[agentToken]; ok {
		return true
	}
	if _, ok := BaitMap[agentToken]; ok {
		return true
	}
	return false
}

func LoadAgentTask(agentToken string) string {

	var agentTask interface{}

	if transparentProxiesMap, ok := TransparentProxyMap[agentToken]; ok {
		for _, task := range transparentProxiesMap {
			LoadedTaskMap[task.TaskId] = comm.TRANSPARENT_PROXY_TASK
			delete(TransparentProxyMap[agentToken], task.TaskId)
			agentTask = task
		}
	}
	if baits, ok := BaitMap[agentToken]; ok {
		for _, task := range baits {
			LoadedTaskMap[task.TaskId] = comm.BAIT_TASK
			delete(BaitMap[agentToken], task.TaskId)
			agentTask = task
		}
	}

	bytes, err := json.Marshal(agentTask)

	if err != nil {
		return ""
	}
	return string(bytes)
}
