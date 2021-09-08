package user_handler

import (
	"decept-defense/models"
	"decept-defense/pkg/app"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginPayload struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type PasswordChangePayload struct {
	OldPassword string `form:"OldPassword" json:"OldPassword" binding:"required"`
	NewPassword string `form:"NewPassword" json:"NewPassword" binding:"required"`
	RepeatPassword string `form:"RepeatPassword" json:"RepeatPassword" binding:"required"`
}



// Login 登录
// @Summary 登录
// @Description 登录
// @Tags 登录管理
// @Produce application/json
// @Accept application/json
// @Param username body models.User true "username"
// @Param password body models.User true "password"
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{"token":""}}"
// @Failure 401 {string} json "{"code":2001,"msg":"账号或是密码错误","data":{}}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/public/login [post]
func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload LoginPayload
	var user  models.User
	err := c.ShouldBindJSON(&payload)
	if err != nil{
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	err = user.CheckPassword(payload.Username, payload.Password)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.ErrorPasswordCheck, err.Error())
		return
	}
	token, err := app.GenerateToken(payload.Username, payload.Password)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, map[string]string{
		"token": token,
		"name" : payload.Username,
	})
}

// SignUp 注册
// @Summary 注册
// @Description 注册
// @Tags 登录管理
// @Produce application/json
// @Accept application/json
// @Param username body models.User true "username"
// @Param password body models.User true "password"
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{"token":""}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/public/signup [post]
func SignUp(c *gin.Context)  {
	appG := app.Gin{C: c}
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil{
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	err = user.HashPassword(user.Password)
	if err != nil{
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}

	p, _ :=  user.GetUserByName(user.Username)
	if p != nil{
		appG.Response(http.StatusOK, app.ErrorDuplicateUser, nil)
		return
	}
	err = user.CreateUserRecord()
	if err != nil{
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	token, err := app.GenerateToken(user.Username, user.Password)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}
	appG.Response(http.StatusOK, app.SUCCESS, map[string]string{
		"token": token,
		"name" : user.Username,
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改密码
// @Tags 登录管理
// @Produce application/json
// @Accept application/json
// @Param NewPassword body PasswordChangePayload true "NewPassword"
// @Param RepeatPassword body PasswordChangePayload true "RepeatPassword"
// @Param OldPassword body PasswordChangePayload true "OldPassword"
// @Success 200 {string} json "{"code":200,"msg":"ok","data":{}}"
// @Failure 500 {string} json "{"code":500,"msg":"内部异常","data":{}}"
// @Router /api/v1/user/password [PUT]
func ChangePassword(c *gin.Context)  {
	appG := app.Gin{C: c}
	var user models.User
	var payload PasswordChangePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil{
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist{
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	u := *(currentUser.(*string))
	err = user.CheckPassword(u, payload.OldPassword)
	if err != nil{
		appG.Response(http.StatusOK, app.ErrorPasswordCheck, "原始密码错误")
		return
	}

	if payload.NewPassword != payload.RepeatPassword{
		appG.Response(http.StatusOK, app.INTERNAlERROR, "新密码两次输入不一致")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.GetPassword(u)), []byte(payload.NewPassword))
	if err == nil{
		appG.Response(http.StatusOK, app.INTERNAlERROR, "新密码与旧密码一致")
		return
	}

	user.HashPassword(payload.NewPassword)
	user.UpdatePassword(u, user.Password)
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}




