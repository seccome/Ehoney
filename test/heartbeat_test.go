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

//TestHeartbeat
//@description heartbeat API unit test

func TestHeartbeat(t *testing.T) {
	setup()
	r := router.MakeRoute()
	body := []byte(`{"username":"akita", "password":"secret"}`)
	w := POSTRequest(r,"/api/public/signup", nil, body)
	assert.Equal(t, http.StatusOK, w.Code)
	body = []byte(`{"username":"akita", "password":"secret"}`)
	w = POSTRequest(r,"/api/public/login", nil, body)
	assert.Equal(t, http.StatusOK, w.Code)
	var response app.Response
	json.Unmarshal(w.Body.Bytes(), &response)
	token := response.Data.(map[string]interface{})["token"]
	header := map[string]string{
		"Authorization" : "Bearer " + token.(string),
	}
	w = GETRequest(r, "/api/v1/health", header)
	assert.Equal(t, http.StatusOK, w.Code)
	var user models.User
	err := user.RevokeAccountByName("akita")
	assert.Nil(t, err)
	teardown()
}
