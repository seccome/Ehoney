package test

import (
	"decept-defense/models"
	"decept-defense/pkg/app"
	"decept-defense/pkg/util"
	"decept-defense/router"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"strconv"
	"testing"
)

func TestSignature(t *testing.T) {
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
	signatureName := "test"
	signatureType := "file"

	header := map[string]string{
		"Authorization" : "Bearer " + token.(string),
		"SignatureType" : signatureType,
		"TokenName" : signatureName,
	}
	w = NewFileUploadRequest(r, "/api/v1/token", header, "file", util.WorkingPath() + "/test/resource/test.txt")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, util.FileExists(path.Join(util.WorkingPath(), viper.GetString("Project.UploadPath"), "token", signatureName, "test.txt")))

	header = map[string]string{
		"Authorization" : "Bearer " + token.(string),
	}
	w = GETRequest(r, "/api/v1/token/" + signatureName, header)
	assert.Equal(t, http.StatusOK, w.Code)
	var signature models.Token
	json.Unmarshal(w.Body.Bytes(), &response)
	data, _ := json.Marshal(response.Data)
	json.Unmarshal(data, &signature)
	assert.Equal(t, "akita", signature.Creator)
	assert.Equal(t, signatureName, signature.TokenName)
	assert.Equal(t, signatureType, signature.TokenType)
	signatureID := signature.ID
	w = DELETERequest(r, path.Join("/api/v1/token/", strconv.FormatInt(signatureID, 10)), header)
	fmt.Println(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, false, util.FileExists(path.Join(util.WorkingPath(), viper.GetString("Project.UploadPath"), "token", signatureName, "test.txt")))
	var user models.User
	err := user.RevokeAccountByName("akita")
	assert.Nil(t, err)
	teardown()
}
