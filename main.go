package main

import (
	"decept-defense/controllers/topology_handler"
	"decept-defense/internal/cluster"
	"decept-defense/internal/cron"
	"decept-defense/internal/env"
	"decept-defense/models"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/logger"
	"decept-defense/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func SetUp() {
	logger.SetUp()

	models.SetUp()

	env.Setup()

	cluster.Setup()

	cron.SetUp()

	topology_handler.Setup()
}

// @title 欺骗防御接口文档
// @version 2.0
// @description

// @license.name MIT
// @license.url

// @host 127.0.0.1:8082
// @BasePath

func main() {
	configs.SetUp()
	serverIp := pflag.String("ip", "", "set server ip")
	pflag.Parse()
	if serverIp != nil && *serverIp != "" {
		configs.UpdateIPConfig(*serverIp)
	}
	gin.SetMode(configs.GetSetting().Server.RunMode)
	SetUp()
	route := router.MakeRoute()
	server := &http.Server{Addr: ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort), Handler: route}
	zap.L().Fatal(server.ListenAndServe().Error())
}
