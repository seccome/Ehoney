package proxy

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/util"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os/exec"
	"strings"
)

func StartProxyProtocol(protocolProxy models.ProtocolProxy) (int, error) {

	if !portCheck(protocolProxy.ProxyPort) {
		zap.L().Error(fmt.Sprintf("proxy port: %d occupied, or to big abort!", protocolProxy.ProxyPort))
		protocolProxy.Status = comm.FAILED
		return 0, nil
	}

	var startCmd = fmt.Sprintf("%s -backend %s:%d -bind :%d -ppid :%s", protocolProxy.ProtocolPath, protocolProxy.HoneypotIp, protocolProxy.HoneypotPort, protocolProxy.ProxyPort, protocolProxy.ProtocolProxyId)

	var startMode = "bash"

	if util.FileExists(protocolProxy.ProtocolPath) == false {
		return 0, errors.New("file not exist")
	}

	startedProxyPid, err := util.StartProcess(startCmd, startMode, "protocol-proxy-"+protocolProxy.ProtocolProxyName)

	if err != nil {
		zap.L().Error(fmt.Sprintf("start proxy err: %v", err))
	} else {
		zap.L().Info(fmt.Sprintf("protocol proxy start success, process [%v]", startedProxyPid))
	}
	return startedProxyPid.Id, nil
}

func portCheck(port int32) bool {

	if port > 65535 {
		zap.L().Error(fmt.Sprintf("Too large port %d, abort!", port))
	}
	checkStatement := fmt.Sprintf("netstat -ano | grep '::%d '", port)

	zap.L().Debug(checkStatement)

	output, _ := exec.Command("sh", "-c", checkStatement).CombinedOutput()

	zap.L().Debug(fmt.Sprintf("output len: %d; out: %s", len(output), output))

	if len(output) > 0 {
		zap.L().Error(fmt.Sprintf("proxy port: %d occupied!", port))
		return false
	}
	return true
}

func TestLocalPortConnection(proxyPort int32) bool {

	return util.TcpGather(strings.Split("127.0.0.1", ","), fmt.Sprintf("%d", proxyPort))
}
