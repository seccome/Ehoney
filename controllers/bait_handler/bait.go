package bait_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"github.com/astaxie/beego/validation"
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

// CreateBait 创建诱饵
// @Summary 创建诱饵
// @Description 创建诱饵
// @Tags 诱捕管理
// @Produce application/json
// @Accept multipart/form-data
// @Param BaitType body BaitCreatePayload true "BaitType"
// @Param BaitName body BaitCreatePayload true "BaitName"
// @Param BaitData body BaitCreatePayload false "BaitData"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 3009 {string} json "{"code":3009,"msg":"诱饵名称重复","data":{}}"
// @Router /api/v1/bait [post]
func CreateBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var bait models.Baits
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
	currentUser, exist := c.Get("currentUser")
	if !exist {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	//TODO .() will panic, best to use value, ok to accept
	bait.Creator = *(currentUser.(*string))
	bait.CreateTime = util.GetCurrentTime()
	if payload.BaitType == "FILE" {
		file, err := c.FormFile("file")
		if err != nil {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		bait.FileName = file.Filename
		savePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "bait", payload.BaitName)
		if err := util.CreateDir(savePath); err != nil {
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
		bait.UploadPath = path.Join(savePath, file.Filename)
		if err := c.SaveUploadedFile(file, path.Join(savePath, file.Filename)); err != nil {
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
	} else if payload.BaitType == "HISTORY" {
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

	if err := bait.CreateBait(); err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetBait 查找诱饵
// @Summary 查找诱饵
// @Description 查找诱饵
// @Tags 诱捕管理
// @Produce application/json
// @Accept application/json
// @Param Payload body comm.BaitSelectPayload true "Payload"
// @Param PageNumber body comm.BaitSelectPayload true "PageNumber"
// @Param PageSize body comm.BaitSelectPayload true "PageSize"
// @Param BaitType body comm.BaitSelectPayload false "BaitType"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.FileBaitSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/bait/set [post]
func GetBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.BaitSelectPayload
	var record models.Baits
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

//func GetBaitByName(c *gin.Context) {
//	appG := app.Gin{C: c}
//	name := c.Param("name")
//	if len(name) == 0{
//		appG.Response(http.StatusOK, app.InvalidParams, nil)
//		return
//	}
//	var bait models.Baits
//	record, err :=  bait.GetBaitByName(name)
//	if err != nil{
//		appG.Response(http.StatusOK, app.ErrorBaitNotExist, nil)
//		return
//	}
//	appG.Response(http.StatusOK, app.SUCCESS, record)
//}

// GetBaitByType 查找支持诱饵类型
// @Summary 查找支持诱饵类型
// @Description 查找支持诱饵类型
// @Tags 诱捕管理
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{["FILE", "HISTORY"]}}"
// @Router /api/v1/bait/type [get]
func GetBaitByType(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.SUCCESS, comm.BaitType)
}

// DeleteBaitByID 删除诱饵
// @Summary 删除诱饵
// @Description 删除诱饵
// @Tags 诱捕管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/bait/:id [delete]
func DeleteBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var bait models.Baits
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := bait.GetBaitByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	if err := bait.DeleteBaitByID(id); err != nil {
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

// DownloadBaitByID 下载文件诱饵
// @Summary 下载文件诱饵
// @Description 下载文件诱饵
// @Tags 诱捕管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/bait/:id [get]
func DownloadBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var bait models.Baits
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := bait.GetBaitByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	if r.BaitType == "FILE" {
		var URL = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/bait/" + r.BaitName + "/" + r.FileName
		appG.Response(http.StatusOK, app.SUCCESS, URL)
	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, "不支持下载此类型诱饵")
	}

}
