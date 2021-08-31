package images_handler

import (
    "decept-defense/controllers/comm"
    "decept-defense/models"
    "decept-defense/pkg/app"
    "github.com/astaxie/beego/validation"
    "github.com/gin-gonic/gin"
    "github.com/unknwon/com"
    "go.uber.org/zap"
    "net/http"
)

// GetImages 获得镜像列表
// @Summary 获得镜像列表
// @Description 获得镜像列表
// @Tags 镜像管理
// @Produce application/json
// @Accept application/json
// @Param Payload body comm.SelectPayload false "Payload"
// @Param PageNumber body comm.SelectPayload true "PageNumber"
// @Param PageSize body comm.SelectPayload true "PageSize"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.Images
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/images/set [post]
func GetImages(c *gin.Context) {
    appG := app.Gin{C: c}
    var payload comm.SelectPayload
    var record models.Images
    err := c.ShouldBindJSON(&payload)
    if err != nil{
        zap.L().Error(err.Error())
        appG.Response(http.StatusOK, app.InvalidParams, err.Error())
        return
    }
    data, count, err := record.GetImage(&payload)
    if err != nil{
        zap.L().Error(err.Error())
        appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
        return
    }
    appG.Response(http.StatusOK, app.SUCCESS, comm.SelectResultPayload{Count: count, List: data})
}

// GetPodImages 获得POD镜像列表
// @Summary 获得POD镜像列表
// @Description 获得POD镜像列表
// @Tags 蜜罐管理
// @Produce application/json
// @Accept application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":[]}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/images/pod [get]
func GetPodImages(c *gin.Context) {
    appG := app.Gin{C: c}
    var record models.Images
    data, err := record.GetPodImageList()
    if err != nil{
        zap.L().Error(err.Error())
        appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
        return
    }
    appG.Response(http.StatusOK, app.SUCCESS, data)
}

// UpdateImage 更新镜像信息
// @Summary 更新镜像信息
// @Description 更新镜像信息
// @Tags 镜像管理
// @Produce application/json
// @Accept application/json
// @Param ImagePort body comm.ImageUpdatePayload true "ImagePort"
// @Param ImageType body comm.ImageUpdatePayload true "ImagePort"
// @Param int query int true "int valid" minimum(1)
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/images/:id [put]
func UpdateImage(c *gin.Context)  {
    appG := app.Gin{C: c}
    var payload comm.ImageUpdatePayload
    var image models.Images
    err := c.ShouldBindJSON(&payload)
    if err != nil{
        zap.L().Error(err.Error())
        appG.Response(http.StatusOK, app.InvalidParams, err.Error())
        return
    }
    valid := validation.Validation{}
    id := com.StrTo(c.Param("id")).MustInt64()
    valid.Min(id, 1, "id").Message("ID必须大于0")
    if valid.HasErrors() {
        appG.Response(http.StatusOK, app.InvalidParams, nil)
        return
    }
    if err = image.UpdateImageByID(id, payload); err != nil{
        zap.L().Error(err.Error())
        appG.Response(http.StatusOK, app.ErrorDatabase, nil)
        return
    }
    appG.Response(http.StatusOK, app.SUCCESS, nil)
}
