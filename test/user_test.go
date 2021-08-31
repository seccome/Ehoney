package test

import (
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/router"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

//TestSignup
//@description signup API unit test

func TestSignup(t *testing.T) {
	setup()
	r := router.MakeRoute()
	body := []byte(`{"username":"", "password":"secret"}`)
	w := POSTRequest(r,"/api/public/signup", nil, body)
	assert.Equal(t, http.StatusOK, w.Code)
	rsp, _ := json.Marshal(app.Response{Code: app.InvalidParams, Msg: app.GetMsg(app.InvalidParams), Data: nil})
	assert.Equal(t, w.Body.Bytes(), rsp)
	body = []byte(`{"username":"akita"}`)
	w = POSTRequest(r,"/api/public/signup", nil, body)
	assert.Equal(t, http.StatusOK, w.Code)
	rsp, _ = json.Marshal(app.Response{Code: app.InvalidParams, Msg: app.GetMsg(app.InvalidParams), Data: nil})
	assert.Equal(t, w.Body.Bytes(), rsp)
	body = []byte(`{"username":"akita", "password":"secret"}`)
	w = POSTRequest(r,"/api/public/signup", nil, body)
	assert.Equal(t, http.StatusOK, w.Code)
	rsp, _ = json.Marshal(app.Response{Code: app.SUCCESS, Msg: app.GetMsg(app.SUCCESS), Data: nil})
	assert.Equal(t, w.Body.Bytes(), rsp)
	var user models.User
	err := user.RevokeAccountByName("akita")
	assert.Nil(t, err)
	teardown()
}

//TestLogin
//@description login API unit test

func TestLogin(t *testing.T) {
	setup()
	r := router.MakeRoute()
	body := []byte(`{"username":"akita", "password":"secret"}`)
	w := POSTRequest(r,"/api/public/signup", nil, body)
	assert.Equal(t, http.StatusOK, w.Code)
	rsp, _ := json.Marshal(app.Response{Code: app.SUCCESS, Msg: app.GetMsg(app.SUCCESS), Data: nil})
	assert.Equal(t, w.Body.Bytes(), rsp)
	body = []byte(`{"username":"", "password":"secret"}`)
	w = POSTRequest(r,"/api/public/login", nil, body)
	assert.Equal(t, http.StatusOK, w.Code)
	rsp, _ = json.Marshal(app.Response{Code: app.InvalidParams, Msg: app.GetMsg(app.InvalidParams), Data: nil})
	assert.Equal(t, w.Body.Bytes(), rsp)
	body = []byte(`{"username":"akita", "password":"secret_error"}`)
	w = POSTRequest(r,"/api/public/login", nil, body)
	assert.Equal(t, http.StatusOK, w.Code)
	rsp, _ = json.Marshal(app.Response{Code: app.ErrorPasswordCheck, Msg: app.GetMsg(app.ErrorPasswordCheck), Data: nil})
	assert.Equal(t, w.Body.Bytes(), rsp)
	body = []byte(`{"username":"akita", "password":"secret"}`)
	w = POSTRequest(r,"/api/public/login", nil, body)
	assert.Equal(t, http.StatusOK, w.Code)
	var user models.User
	err := user.RevokeAccountByName("akita")
	assert.Nil(t, err)
	teardown()
}
