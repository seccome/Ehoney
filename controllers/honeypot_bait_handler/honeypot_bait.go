package honeypot_bait_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/internal/cluster"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

type HoneypotBaitCreatePayload struct {
	BaitId     string `form:"BaitId" binding:"required"`
	DeployPath string `form:"DeployPath"`
	HoneypotId string `form:"HoneypotId" binding:"required"`
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
// @Param Authorization header string true "Insert your access token_builder" default(Bearer <Add access token_builder here>)
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
	var honeypotBait models.BaitTask
	var bait models.Bait
	var honeypot models.Honeypot

	var payload HoneypotBaitCreatePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}

	honeypotBait.BaitTaskId = util.GenerateId()
	honeypotBait.CreateTime = util.GetCurrentIntTime()
	r, err := bait.GetBaitById(payload.BaitId)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorBaitNotExist, nil)
		return
	}

	honeypotBait.BaitId = r.BaitId

	s, err := honeypot.GetHoneypotByID(payload.HoneypotId)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorHoneypotServerNotExist, nil)
		return
	}

	if r.BaitType == "FILE" {
		if payload.DeployPath == "" {
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		honeypotBait.HoneypotId = payload.HoneypotId
		honeypotBait.BaitName = r.BaitName
		honeypotBait.BaitType = r.BaitType
		honeypotBait.DeployPath = payload.DeployPath
		honeypotBait.LocalPath = r.UploadPath
		//TODO should check dest dir is valid fixme
		if err = cluster.CopyToPod(s.PodName, s.HoneypotName, r.UploadPath, path.Join(payload.DeployPath, path.Base(r.UploadPath))); err != nil {
			appG.Response(http.StatusOK, app.ErrorHoneypotK8SCP, nil)
			return
		}
	} else if r.BaitType == "HISTORY" {
		appG.Response(http.StatusOK, app.ErrorHoneypotHistoryBait, nil)
		return
	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	if err := honeypotBait.CreateBaitTask(); err != nil {
		if r.BaitType == "FILE" {
			_ = cluster.RemoveFromPod(s.PodName, s.HoneypotName, path.Join(payload.DeployPath, path.Base(r.UploadPath)))
		}
		appG.Response(http.StatusOK, app.ErrorHoneypotBaitCreate, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

func GetHoneypotBait(c *gin.Context) {
	appG := app.Gin{C: c}
	var queryMap map[string]interface{}
	var baitTasks models.BaitTask
	err := c.ShouldBindJSON(&queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := baitTasks.QueryBaitTaskPage(queryMap)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

func DeleteHoneypotBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypotBait models.BaitTask
	var honeypot models.Honeypot
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := honeypotBait.GetBaitTaskById(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	s, err := honeypot.GetHoneypotByID(r.HoneypotId)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if r.BaitType == "FILE" {
		err = cluster.RemoveFromPod(s.PodName, s.HoneypotName, path.Join(r.DeployPath, path.Base(r.LocalPath)))
		if err != nil {
			appG.Response(http.StatusOK, app.ErrorHoneypotBaitWithdraw, nil)
			return
		}
		_ = os.RemoveAll(filepath.Dir(filepath.Dir(r.LocalPath)))
	} else {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	if err := honeypotBait.DeleteBaitTaskById(id); err != nil {
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
// @Param Authorization header string true "Insert your access token_builder" default(Bearer <Add access token_builder here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token_builder/id [get]
func DownloadHoneypotBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypotBait models.BaitTask
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := honeypotBait.GetBaitTaskById(id)
	if err != nil {
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
	c.Header("Content-Disposition", "attachment; filename="+path.Base(r.LocalPath))
	c.Header("Content-Type", "application/text/plain")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	_, _ = c.Writer.Write(content)
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// DownloadHoneyBaitByID 下载蜜罐诱饵
// @Summary 下载蜜罐诱饵
// @Description 下载蜜罐诱饵
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token_builder" default(Bearer <Add access token_builder here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/bait/probe/:id [get]
func DownloadHoneyBaitByID(c *gin.Context) {
	appG := app.Gin{C: c}
	var honeypotBaits models.Bait
	id := com.StrTo(c.Param("id")).String()
	if id == "" {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	p, err := honeypotBaits.GetBaitById(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	URL := "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/bait/" + p.BaitName + "/" + path.Base(p.LocalPath)

	appG.Response(http.StatusOK, app.SUCCESS, URL)
}
