package cron

import (
	"decept-defense/controllers/agent_handler"
	"decept-defense/controllers/honeypot_handler"
	"decept-defense/controllers/protocol_proxy_handler"
	"decept-defense/controllers/trans_proxy_handler"
	"time"
)

//SetUp cron task
func SetUp() {
	RefreshServerStatus()
	RefreshProxyStatus()
}

func RefreshServerStatus() {
	ticker := time.NewTicker(time.Second * 10)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				honeypot_handler.RefreshHoneypotStatus()
				agent_handler.RefreshAgentStatus()
			}
		}
	}()
}

func RefreshProxyStatus() {
	ticker := time.NewTicker(time.Second * 10)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				protocol_proxy_handler.UpdateDeployedProtocolProxyStatus()
				trans_proxy_handler.UpdateTransparentProxiesStatus()
			}
		}
	}()
}
