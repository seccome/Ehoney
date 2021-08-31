package main

import (
	"decept-defense/controllers/topology_handler"
	"decept-defense/internal/cluster"
	"decept-defense/internal/cron"
	"decept-defense/internal/harbor"
	message "decept-defense/internal/message_client"
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
	//logger module
	logger.SetUp()
	//db module
	models.SetUp()
	//message_client module
	message.Setup()
	//cluster module
	cluster.Setup()
	//harbor module
	harbor.Setup()
	//cron module
	cron.SetUp()
	//topology module
	topology_handler.Setup()
}

// @title 欺骗防御接口文档
// @version 2.0
// @description

// @contact.name akita.tian
// @contact.url
// @contact.email chronoyq@163.com

// @license.name MIT
// @license.url

// @host 127.0.0.1:8082
// @BasePath

func main() {
	//config module
	configs.SetUp()
	config := pflag.String("CONFIGS", "", "one click install parameter")
	pflag.Parse()
	configs.UpdateIPConfig(*config)
	gin.SetMode(configs.GetSetting().Server.RunMode)
	SetUp()
	router := router.MakeRoute()
	server := &http.Server{Addr: ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort), Handler: router}
	zap.L().Fatal(server.ListenAndServe().Error())
}
