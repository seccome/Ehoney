package honeypot_bait_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/internal/cluster"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type  HoneypotBaitCreatePayload struct {
	BaitID       int64                   `form:"BaitID" binding:"required"`
	DeployPath   string                  `form:"DeployPath"`
	HoneypotID   int64                   `form:"HoneypotID" binding:"required"`
}

// CreateHoneypotBait 创建蜜罐诱饵
// @Summary 创建蜜罐诱饵
// @Description 创建蜜罐诱饵
// @Tags 蜜罐管理
// @Produce application/json
// @Accept multipart/form-data
// @Param BaitID body HoneypotBaitCreatePayload true "BaitID"
// @Param DeployPath body HoneypotBaitCreatePayload false "DeployPath"
// @Param HoneypotID body HoneypotBaitCreatePayload true "HoneypotID"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 5001 {string} json "{"code":5001,"msg":"蜜罐服务器不存在、请检测蜜罐服务状态","data":{}}"
// @Failure 3005 {string} json "{"code":3005,"msg":"蜜罐诱饵创建异常","data":{}}"
// @Failure 3011 {string} json "{"code":3011,"msg":"诱饵不存在","data":{}}"
// @Failure 3012 {string} json "{"code":3012,"msg":"K8S拷贝异常","data":{}}"
// @Failure 3013 {string} json "{"code":3013,"msg":"暂不支持下发history类型诱饵到蜜罐","data":{}}"
// @Router /api/v1/bait/honeypot [post]
func CreateHoneypotBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypotBait models.HoneypotBaits
	var bait models.Baits
	var honeypot models.Honeypot
	var payload HoneypotBaitCreatePayload
	err :=  c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist{
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	//TODO .() will panic, best to use value, ok to accept
	honeypotBait.Creator = *(currentUser.(*string))
	honeypotBait.CreateTime = util.GetCurrentTime()

	r, err :=  bait.GetBaitByID(payload.BaitID)
	if err != nil{
		appG.Response(http.StatusOK, app.ErrorBaitNotExist, nil)
		return
	}
	s, err :=  honeypot.GetHoneypotByID(payload.HoneypotID)
	if err != nil{
		appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, nil)
		return
	}

	if r.BaitType == "FILE"{
		if payload.DeployPath == ""{
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		honeypotBait.HoneypotID = payload.HoneypotID
		honeypotBait.BaitName = r.BaitName
		honeypotBait.BaitType = r.BaitType
		honeypotBait.DeployPath = payload.DeployPath
		honeypotBait.LocalPath = r.UploadPath
		//TODO should check dest dir is valid fixme
		if err =  cluster.CopyToPod(s.PodName, s.HoneypotName, r.UploadPath, path.Join(payload.DeployPath, path.Base(r.UploadPath))); err != nil{
			appG.Response(http.StatusOK, app.ErrorHoneypotK8SCP, nil)
			return
		}
	}else if r.BaitType == "HISTORY"{
		appG.Response(http.StatusOK, app.ErrorHoneypotHistoryBait, nil)
		return
	}else{
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	if err := honeypotBait.CreateHoneypotBait(); err != nil{
		if r.BaitType == "FILE"{
			cluster.RemoveFromPod(s.PodName, s.HoneypotName, path.Join(payload.DeployPath, path.Base(r.UploadPath)))
		}
		appG.Response(http.StatusOK, app.ErrorHoneypotBaitCreate, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetHoneypotBait 查找蜜罐诱饵
// @Summary 查找蜜罐诱饵
// @Description 查找蜜罐诱饵
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param BaitType body comm.ServerBaitSelectPayload false "BaitType"
// @Param BaitName body comm.ServerBaitSelectPayload false "BaitName"
// @Param Status body comm.ServerBaitSelectPayload false "Status"
// @Param StartTimestamp body comm.ServerBaitSelectPayload false "StartTimestamp"
// @Param EndTimestamp body comm.ServerBaitSelectPayload false "EndTimestamp"
// @Param PageNumber body comm.ServerBaitSelectPayload true "PageNumber"
// @Param PageSize body comm.ServerBaitSelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.ServerBaitSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/bait/honeypot/set [post]
func GetHoneypotBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.ServerBaitSelectPayload
	var honeypotBaits models.HoneypotBaits
	err := c.ShouldBindJSON(&payload)
	if err != nil{
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := honeypotBaits.GetHoneypotBait(&payload)
	if err != nil{
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// DeleteHoneypotBaitByID 删除蜜罐诱饵接口
// @Summary 删除蜜罐诱饵接口
// @Description 删除蜜罐诱饵接口
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/bait/honeypot/:id [delete]
func DeleteHoneypotBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypotBait models.HoneypotBaits
	var honeypot models.Honeypot
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := honeypotBait.GetHoneypotBaitByID(id)
	if err != nil{
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	s, err :=  honeypot.GetHoneypotByID(r.HoneypotID)
	if err != nil{
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if r.BaitType == "FILE"{
		err = cluster.RemoveFromPod(s.PodName, s.HoneypotName, path.Join(r.DeployPath, path.Base(r.LocalPath)))
		if err != nil{
			appG.Response(http.StatusOK, app.ErrorHoneypotBaitWithdraw, nil)
			return
		}
		os.RemoveAll(filepath.Dir(filepath.Dir(r.LocalPath)))
	}else{
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if err := honeypotBait.DeleteHoneypotBaitByID(id); err != nil{
		appG.Response(http.StatusOK, app.ErrorHoneypotBaitDelete, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// DownloadHoneypotBaitByID 下载蜜罐诱饵接口
// @Summary 下载蜜罐诱饵接口
// @Description 下载蜜罐诱饵接口
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token/id [get]
func DownloadHoneypotBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var honeypotBait models.HoneypotBaits
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := honeypotBait.GetHoneypotBaitByID(id)
	if err != nil{
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	file, err := os.Open(r.LocalPath)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	defer file.Close()
	content, err := ioutil.ReadAll(file)
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=" + path.Base(r.LocalPath))
	c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Writer.Write([]byte(content))

	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// DownloadHoneyBaitByID 下载蜜罐诱饵
// @Summary 下载蜜罐诱饵
// @Description 下载蜜罐诱饵
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/bait/probe/:id [get]
func DownloadHoneyBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var honeypotBiat models.HoneypotBaits
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	p, err := honeypotBiat.GetHoneypotBaitByID(id)
	if err != nil{
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	var URL string = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath  + "/bait/" + p.BaitName + "/" + path.Base(p.LocalPath)

	appG.Response(http.StatusOK, app.SUCCESS, URL)
}
