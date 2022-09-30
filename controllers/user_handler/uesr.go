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
	OldPassword    string `form:"OldPassword" json:"OldPassword" binding:"required"`
	NewPassword    string `form:"NewPassword" json:"NewPassword" binding:"required"`
	RepeatPassword string `form:"RepeatPassword" json:"RepeatPassword" binding:"required"`
}

func Login(c *gin.Context) {
	appG := app.Gin{C: c}
	var payload LoginPayload
	var user models.User
	err := c.ShouldBindJSON(&payload)
	if err != nil {
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
		"name":  payload.Username,
	})
}

func SignUp(c *gin.Context) {
	appG := app.Gin{C: c}
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	err = user.HashPassword(user.Password)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.INTERNAlERROR, err.Error())
		return
	}

	p, _ := user.GetUserByName(user.Username)
	if p != nil {
		appG.Response(http.StatusOK, app.ErrorDuplicateUser, nil)
		return
	}
	err = user.CreateUserRecord()
	if err != nil {
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
		"name":  user.Username,
	})
}

func ChangePassword(c *gin.Context) {
	appG := app.Gin{C: c}
	var user models.User
	var payload PasswordChangePayload
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		zap.L().Error(err.Error())
		appG.Response(http.StatusOK, app.InvalidParams, err.Error())
		return
	}
	currentUser, exist := c.Get("currentUser")
	if !exist {
		appG.Response(http.StatusOK, app.INTERNAlERROR, nil)
		return
	}
	u := *(currentUser.(*string))
	err = user.CheckPassword(u, payload.OldPassword)
	if err != nil {
		appG.Response(http.StatusOK, app.ErrorPasswordCheck, "原始密码错误")
		return
	}

	if payload.NewPassword != payload.RepeatPassword {
		appG.Response(http.StatusOK, app.INTERNAlERROR, "新密码两次输入不一致")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.GetPassword(u)), []byte(payload.NewPassword))
	if err == nil {
		appG.Response(http.StatusOK, app.INTERNAlERROR, "新密码与旧密码一致")
		return
	}

	user.HashPassword(payload.NewPassword)
	user.UpdatePassword(u, user.Password)
	appG.Response(http.StatusOK, app.SUCCESS, nil)
}
