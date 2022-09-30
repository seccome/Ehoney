package images_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
)

func GetImages(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.SelectPayload
	var record models.Images
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	data, count, err := record.GetImage(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

func CreateImage(c *gin.Context) {
	appG := app.Gin{C: c}
	var record models.Images
	err := c.ShouldBindJSON(&record)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	record.ImageId = util.GenerateId()
	record.CreateTime = util.GetCurrentIntTime()
	err = record.CreateImage()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetPodImages(c *gin.Context) {
	appG := app.Gin{C: c}
	var record models.Images
	data, err := record.GetPodImageList()
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, data)
}

func UpdateImage(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.ImageUpdatePayload
	var image models.Images
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	if err = image.UpdateImageByID(id, payload); err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}
