package token_handler

import (
	"decept-defense/controllers/comm"
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path"
	"strconv"
)

type TokenCreatePayload struct {
	TokenType string `json:"TokenType" form:"TokenType" binding:"required"`
	TokenName string `json:"TokenName" form:"TokenName" binding:"required"`
	TokenData string `json:"TokenData" form:"TokenData"`
}

// CreateToken 创建密签
// @Summary 创建密签
// @Description 创建密签
// @Tags 诱捕管理
// @Produce application/json
// @Accept multipart/form-data
// @Param TokenType body TokenCreatePayload true "TokenType"
// @Param TokenName body TokenCreatePayload true "TokenName"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Failure 1006 {string} json "{"code":1006,"msg":"数据库异常","data":{}}"
// @Failure 3010 {string} json "{"code":3010,"msg":"密签名称重复","data":{}}"
// @Router /api/v1/token [post]
func CreateToken(c *gin.Context) {
	appG := app.Gin{C: c}
	var token models.Token
	var payload TokenCreatePayload
	err := c.ShouldBind(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	_, err = token.GetTokenByName(payload.TokenName)
	if err == nil {
		zap.L().Error("密签名称重复")
		appG.Response(http.StatusOK, app.ErrorDuplicateTokenName, nil)
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist {
		zap.L().Error("当前用户获取错误")
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	value, ok := currentUser.(*string)
	if !ok {
		zap.L().Error("当前用户名称解析错误")
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	token.Creator = *(value)
	token.CreateTime = util.GetCurrentTime()
	if payload.TokenType == "FILE" || payload.TokenType == "EXE" {
		file, err := c.FormFile("file")
		if err != nil {
			zap.L().Error("参数错误、文件类型密签未上传密签文件")
			appG.Response(http.StatusOK, app.InvalidParams, nil)
			return
		}
		fileExt := path.Ext(file.Filename)
		if payload.TokenType == "FILE" && fileExt != ".pdf" && fileExt != ".docx" && fileExt != ".xlsx" && fileExt != ".pptx" {
			zap.L().Error("不支持的文件密签文件类型")
			appG.Response(http.StatusOK, app.ErrorFileTokenType, nil)
			return
		}
		if payload.TokenType == "EXE" && fileExt != ".exe" {
			zap.L().Error("不支持的EXE密签文件类型")
			appG.Response(http.StatusOK, app.ErrorFileTokenType, nil)
			return
		}
		token.FileName = file.Filename
		savePath := path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "token", payload.TokenName)
		if err := util.CreateDir(savePath); err != nil {
			zap.L().Error("创建密签路径异常")
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
		token.UploadPath = path.Join(savePath, file.Filename)
		if err := c.SaveUploadedFile(file, path.Join(savePath, file.Filename)); err != nil {
			zap.L().Error("保存密签文件异常")
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
	} else if payload.TokenType == "BrowserPDF" || payload.TokenType == "WPS" {
		token.TokenData = payload.TokenData
	} else {
		zap.L().Error("未支持的密签类型")
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	token.TokenType = payload.TokenType
	token.TokenName = payload.TokenName

	if err := token.CreateToken(); err != nil {
		zap.L().Error("创建密签失败、数据库异常")
		appG.Response(http.StatusOK, app.ErrorDatabase, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// GetToken 查找密签
// @Summary 查找密签
// @Description 查找密签
// @Tags 诱捕管理
// @Produce application/json
// @Accept application/json
// @Param Payload body comm.SelectPayload false "Payload"
// @Param PageNumber body comm.SelectPayload true "PageNumber"
// @Param PageSize body comm.SelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} comm.TokenSelectResultPayload
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token/set [post]
func GetToken(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload comm.SelectPayload
	var record models.Token
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	data, count, err := record.GetToken(&payload)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// GetTokenNameList 查看密签名称列表
// @Summary 查看密签名称列表
// @Description 查看密签名称列表
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}"
// @Router /api/v1/token/name/set [get]
func GetTokenNameList(c *gin.Context) {
	appG := app.Gin{C: c}
	var sign models.Token
	record, err := sign.GetTokenNameList()
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, record)
}

// GetTokenByType 查看支持密签类型
// @Summary 查看支持密签类型
// @Description 查看支持密签类型
// @Tags 诱捕管理
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{["FILE", "HISTORY"]}}"
// @Router /api/v1/token/type [get]
func GetTokenByType(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, app.SUCCESS, comm.TokenType)
}

// DeleteTokenByID 删除密签
// @Summary 删除密签
// @Description 删除密签
// @Tags 诱捕管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token/id [delete]
func DeleteTokenByID(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var sign models.Token
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := sign.GetTokenByID(id)
	if err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, "密签不存在")
		return
	}

	if err := sign.DeleteTokenByID(id); err != nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	//TODO could ignore this error
	if util.FileExists(r.UploadPath) {
		err = os.Remove(r.UploadPath)
		if err != nil {
			appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
			return
		}
	}
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}

// DownloadTokenByID 下载密签
// @Summary 下载密签
// @Description 下载密签
// @Tags 诱捕管理
// @Produce application/json
// @Accept application/json
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/token/id [get]
func DownloadTokenByID(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt64()
	var token models.Token
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		appG.Response(http.StatusOK, app.InvalidParams, nil)
		return
	}
	r, err := token.GetTokenByID(id)
	if err != nil {
		zap.L().Error("token不存在")
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}

	var URL string = "http:" + "//" + configs.GetSetting().Server.AppHost + ":" + strconv.Itoa(configs.GetSetting().Server.HttpPort) + "/" + configs.GetSetting().App.UploadPath + "/token/" + r.TokenName + "/" + r.FileName

	appG.Response(http.StatusOK, app.SUCCESS, URL)
}
