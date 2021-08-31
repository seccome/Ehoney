package virus_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateVirusRecord(c *gin.Context) {
	appG := app.Gin{C: c}
	var record models.VirusRecord
	err := c.ShouldBindJSON(&record)
	if err != nil{
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	record.CreateTime = util.GetCurrentTime()
	err = record.CreateVirusRecord()
	if err != nil{
		appG.Response(http.StatusOK, app.ErrorCreateVirusRecord, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func SelectVirusRecord(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.SelectVirusPayload
	var record models.VirusRecord
	err := c.ShouldBindJSON(&payload)
	if err != nil{
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, err := record.GetVirusRecord(&payload)
	if err != nil{
		appG.Response(http.StatusOK, app.ErrorSelectVirusRecord, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, data)
}