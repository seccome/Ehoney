package comhttp

import (
	"bytes"
	"decept-defense/models/util/comm"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func GetPostBody(w http.ResponseWriter, r *http.Request) string {
	r.ParseForm()
	if r.Method != "POST" && r.Method != "OPTIONS" {
		log.Printf("Request Method Is Invalid %s %s %s\n", r.Method, r.RequestURI, r.RemoteAddr)
		return ""
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read the request body: %v\n", err)
		return ""
	}
	r.Body.Close()

	return string(body)
}

func SendJSONResponse(w http.ResponseWriter, data interface{}) {
	worigin := beego.AppConfig.String("worigin")
	w.Header().Set("Access-Control-Allow-Origin", worigin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/octet-stream;text/plain; charset=UTF-8")
	w.Header().Add("Access-Control-Allow-Headers", "Origin,Access-Control-Allow-Origin,Access-Control-Allow-Credentials,Access-Control-Allow-Headers,Content-Type,Accept,X-Requested-With,Date,Content-Length,Connection")
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")

	body, err := json.Marshal(data)
	if err != nil {
		logs.Error("[SendJSONResponse] Failed to encode a JSON response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		logs.Error("[SendJSONResponse] Failed to write the response body: %v\n", err)
		return
	}
}

func HttpPostWithNoResponse(requrl string, data string) error {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Post(requrl,
		"application/octet-stream",
		strings.NewReader(data))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("response code error")
	}

	defer resp.Body.Close()

	return nil
}

func HttpGet(requrl string) string {
	respbody := ""
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(requrl)
	if err != nil {
		log.Println("HttpGet Error:", err)
	}

	if resp.StatusCode != 200 {
		log.Println("HttpGet response code error")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		respbody = string(body)
	}

	defer resp.Body.Close()
	return respbody
}

func GetResponseByHttp(method string, url string, postbody string, headers map[string]string) *http.Response {
	var jsonStr = []byte(postbody)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if headers["authorization"] != "" {
		req.Header.Set("authorization", headers["authorization"])
	}
	req.Header.Set("Content-Type", headers["content-type"])
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp == nil && err != nil {
		logs.Error("doHttp Error: %v\n", err)
	}
	return resp

}

func Str2Map(jsonData string) (result map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(jsonData), &result)
	return result, err
}
func FilteredSQLInject(to_match_str string) bool {
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		panic(err.Error())
		return false
	}
	return re.MatchString(to_match_str)
}

// telnet
func InterfaceToString(logType string, logData []byte) string {
	fmt.Println(logType)
	logs.Info(logType)
	logs.Info(string(logData))
	if logType == "mysql" {
		var logdata comm.MysqlLogData
		if err := json.Unmarshal(logData, &logdata); err != nil {
			return ""
		}
		if logdata.SqlType == "Login" {
			return "登录" + logdata.UserName + "账号"
		} else {
			return "使用" + logdata.UserName + "账号" + "执行" + logdata.Sql
		}
	} else if logType == "ssh" {
		var logdata comm.SSHLogData
		if err := json.Unmarshal(logData, &logdata); err != nil {
			return ""
		}
		logs.Info("ssh logdata %v", logdata)
		command := strings.Split(logdata.Command, ",")
		if "" != command[0] {
			return command[0] + ",use " + logdata.UserName
		}

	} else if logType == "redis" {
		var logdata comm.RedisLogData
		if err := json.Unmarshal(logData, &logdata); err != nil {
			return ""
		}
		logmsg := ""
		if logdata.UserName != "" {
			logmsg = ",use account: " + logdata.UserName
		}
		if logdata.PassWord != "" {
			logmsg = logmsg + ",use password: " + logdata.UserName
		}
		return "执行命令：" + logdata.Command + logmsg
	} else if logType == "http" {
		var logdata comm.HttpLogData
		if err := json.Unmarshal(logData, &logdata); err != nil {
			return ""
		}
		command := strings.Split(logdata.Command, ",")
		return command[0]
	} else if logType == "telnet" {
		var logdata comm.TelnetLogData
		if err := json.Unmarshal(logData, &logdata); err != nil {
			return ""
		}
		command := strings.Split(logdata.Command, ",")
		return command[0]
	} else {
		return ""
	}
	return ""
}
