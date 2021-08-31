package cron

import (
	"decept-defense/controllers/honeypot_handler"
	"decept-defense/controllers/probe_handler"
	"decept-defense/controllers/protocol_proxy_handler"
	"fmt"
	"time"
)

//SetUp cron task
func SetUp(){
	RefreshServerStatus()
	RefreshProxyStatus()
}

func RefreshServerStatus(){
	ticker := time.NewTicker(time.Second * 5)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				honeypot_handler.RefreshHoneypotStatus()
				probe_handler.RefreshProbeStatus()
				fmt.Println("refresh status, tick at", t)
			}
		}
	}()
}

func RefreshProxyStatus(){
	ticker := time.NewTicker(time.Hour * 1)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				protocol_proxy_handler.UpdateProtocolProxyStatus()
				//trans_proxy_handler.UpdateTransparentProxyStatus()
				fmt.Println("refresh proxy status, tick at", t)
			}
		}
	}()
}


