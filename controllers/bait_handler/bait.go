package bait_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/token_builder"
	"decept-defense/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"os"
	"path"
	"strconv"
)

type BaitCreatePayload struct {
	BaitType string `form:"BaitType" json:"BaitType" binding:"required"` //诱饵类型
	BaitName string `form:"BaitName" json:"BaitName" binding:"required"` //诱饵名称
	BaitData string `form:"BaitData" json:"BaitData"`                    //HISTORY诱饵数据
}

func CreateBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var bait models.Bait
	var payload BaitCreatePayload
	err := c.ShouldBind(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	_, err = bait.GetBaitByName(payload.BaitName)
	if err == nil {
		appG.Response(http.StatusOK, app.ErrorDuplicateBaitName, nil)
		return
	}

	if payload.BaitType == "FILE" || payload.BaitType == "EXE" {
		file, err := c.FormFile("file")
		if err != nil {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		bait.FileName = file.Filename
		savePath := path.Join(util.WorkingPath(), "upload", "bait", payload.BaitName)
		if err := util.CreateDir(savePath); err != nil {
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
		bait.UploadPath = path.Join(savePath, file.Filename)
		if err := c.SaveUploadedFile(file, bait.UploadPath); err != nil {
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
	} else if payload.BaitType == "HISTORY" || payload.BaitType == "WPS" {
		data := c.PostForm("BaitData")
		if data == "" {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		bait.BaitData = data
	} else {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	bait.BaitType = payload.BaitType
	bait.BaitName = payload.BaitName
	bait.LocalPath = bait.UploadPath
	if err := bait.CreateBait(); err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func SignTokenForFileBait(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	var baits models.Bait

	bait, err := baits.GetBaitById(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if !bait.TokenAble {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	newLocalPath, tokenTraceUrl, err := token_builder.TokenBaitFile(*bait)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	bait.UploadPath = newLocalPath
	bait.TokenTraceUrl = tokenTraceUrl
	err = bait.UpdateForToken()
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, err)
}

func GetBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.BaitSelectPayload
	var record models.Bait
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	var result comm.SelectResultPayload
	//TODO different bait type
	data, count, err := record.GetBaitsRecord(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	result = comm.SelectResultPayload{Count: count, List: data}
	appG.Response(http.StatusOK, app.SUCCESS, result)
}

func GetBaitByType(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.SUCCESS, comm.BaitType)
}

func DeleteBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	var bait models.Bait

	r, err := bait.GetBaitById(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	if err := bait.DeleteBaitById(id); err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	//TODO can ignore this error
	if util.FileExists(r.UploadPath) {
		err = os.Remove(r.UploadPath)
		if err != nil {
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func DownloadBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).String()
	var bait models.Bait
	r, err := bait.GetBaitById(id)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorBaitNotExist, nil)
		return
	}
	if r.BaitType == "FILE" {
		var URL = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/bait/" + r.BaitName + "/" + r.FileName
		appG.Response(http.StatusOK, app.SUCCESS, URL)
	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, "不支持下载此类型诱饵")
	}

}
