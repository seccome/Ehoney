package test

import (
	"bytes"
	"decept-defense/models"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/logger"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

//setup
//@description setup
func setup()  {
	//logger module
	logger.SetUp()
	//config module
	configs.SetUp()
	//db module
	models.SetUp()
}

//teardown
//@description teardown
func teardown()  {

}

//POSTRequest
//@description send post request
func POSTRequest(r http.Handler, path string, header map[string]string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, path,  bytes.NewBuffer(body))
	req.Header.Set("Content-TaskType", "application/json")
	for k, v := range header{
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

//GETRequest
//@description send get request
func GETRequest(r http.Handler, path string, header map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	for k, v := range header{
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

//DELETERequest
//@description send delete request
func DELETERequest(r http.Handler, path string, header map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodDelete, path, nil)
	for k, v := range header{
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}


// NewFileUploadRequest Creates a new file upload http request with optional extra params
func NewFileUploadRequest(r http.Handler, url string, header map[string]string, paramName, path string) *httptest.ResponseRecorder {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil
	}
	io.Copy(part, file)
	for key, val := range header {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil
	}
	req, _  := http.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-TaskType", writer.FormDataContentType())
	for k, v := range header{
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}