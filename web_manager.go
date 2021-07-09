package main

import (
	"context"
	"database/sql"
	"decept-defense/models/clamavcenter"
	"decept-defense/models/datavcenter"
	"decept-defense/models/honeycluster"
	"decept-defense/models/policyCenter"
	"decept-defense/models/redisCenter"
	"decept-defense/models/util"
	"decept-defense/models/util/comhttp"
	"decept-defense/models/util/comm"
	"decept-defense/models/util/honeytoken"
	"decept-defense/models/util/k3s"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/kataras/go-sessions/v3"
	"io"
	"io/ioutil"
	apiV1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"net"

	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var typeRegistry = make(map[string]reflect.Type)
var (
	dbhost     = beego.AppConfig.String("dbhost")
	dbport     = beego.AppConfig.String("dbport")
	dbuser     = beego.AppConfig.String("dbuser")
	dbpassword = beego.AppConfig.String("dbpassword")
	dbname     = beego.AppConfig.String("dbname")
	sess       = sessions.New(sessions.Config{
		Cookie:                      "ehoney-session",
		Expires:                     time.Hour * 2,
		DisableSubdomainPersistence: false,
	})
)

var sessionCookieName = "ehoney-session"

func healthHandler(w http.ResponseWriter, r *http.Request) {
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: "success", Message: "心跳正常"})
	return
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	var adminuser comm.Admin
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &adminuser); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	result := honeycluster.CheckAdminLogin(adminuser.UserName, adminuser.Password)
	if len(result) > 0 {
		s := sess.Start(w, r)
		s.Set("name", adminuser.UserName)
		if result[0]["status"] == "1" {
			honeycluster.UpdateAdminLoginStatus(adminuser.UserName, adminuser.Password)
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: 1, Message: "登录成功"})
			return
		} else {
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: "登录成功"})
			return
		}
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "登录失败"})
		return
	}

}

func AdminLogout(w http.ResponseWriter, r *http.Request) {

	sess.Destroy(w, r)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: nil, Message: "成功"})
	return
}

// 拓扑图
func TopologyMap(w http.ResponseWriter, r *http.Request) {
	datas := datavcenter.GetTopAttackMap()

	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetTopAttackMap(w http.ResponseWriter, r *http.Request) {
	datas := datavcenter.GetTopAttackMap()
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetTopAttackIps(w http.ResponseWriter, r *http.Request) {
	datas := datavcenter.SelectTopAttackIps()
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetTopAttackTypes(w http.ResponseWriter, r *http.Request) {
	datas := datavcenter.SelectTopAttackTypes()
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetTopSourceIps(w http.ResponseWriter, r *http.Request) {
	datas := datavcenter.SelectTopSourceIps()
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetTopAreas(w http.ResponseWriter, r *http.Request) {
	datas := datavcenter.SelectTopAreas()
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func CreateHoneyServer(w http.ResponseWriter, r *http.Request) {
	var appclusters comm.ApplicationsClusters
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &appclusters); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	serverid := util.GetUUID()
	data, msg, code := honeycluster.InsertApplicationCluster(appclusters.ServerName, appclusters.ServerIp, serverid, appclusters.Agentid)
	comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: data, Message: msg})
	return
}

func GetApplicationClusters(w http.ResponseWriter, r *http.Request) {
	var appclusters comm.ApplicationsClusters
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &appclusters); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectApplicationClusters(appclusters.ServerIp, appclusters.ServerName, appclusters.VpcName, appclusters.Status, appclusters.PageSize, appclusters.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetApplicationLists(w http.ResponseWriter, r *http.Request) {
	var appclusters comm.ApplicationsClusters
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &appclusters); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	nowtime := time.Now().Unix()
	nowtimestr := util.Strval(nowtime - 300)
	datas := honeycluster.SelectApplicationLists(nowtimestr)
	if len(datas) > 0 {
		for i, value := range datas {
			datas[i]["servername"] = util.Strval(value["servername"]) + "(" + util.Strval(value["serverip"]) + " 转发数：" + util.Strval(value["servercount"]) + ")"
		}
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func CreateApplicationClusters(w http.ResponseWriter, r *http.Request) {
	var appclusters comm.ApplicationsClusters
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &appclusters); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	ecsId := util.GetUUID()
	datas, msg, _ := honeycluster.InsertApplication(appclusters.ServerName, appclusters.ServerIp, ecsId, 1, appclusters.VpcName, "admin", appclusters.Agentid)
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: datas, Message: msg})
}

/*
撤回诱饵策略
*/
func DeleteApplicationBaitPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var bait comm.BaitJson
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &bait); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	sysproto := ""
	if r.TLS == nil {
		sysproto = "http://"
	} else {
		sysproto = "https://"
	}
	datas := honeycluster.SelectApplicationBaitsById(bait.TaskId)
	baitfilepath := datas[0]["data"].(string)
	baitfilename := datas[0]["baitfilename"].(string)
	baitfilepath = baitfilepath + "/" + baitfilename

	baitname := datas[0]["baitname"].(string) + "-uninstall"
	util.ModifyUninstallFile(baitfilepath, baitname)
	baituploadpath := "upload/" + baitname
	currentpath := util.GetCurrentPathString()
	util.FileTarZip(baituploadpath, baituploadpath+".tar.gz")
	filemd5 := util.GetFileMd5(currentpath + "/" + baituploadpath + ".tar.gz")
	apphost := beego.AppConfig.String("apphost")
	appport := beego.AppConfig.String("appport")
	apphost = apphost + ":" + appport
	baiturl := sysproto + apphost + "/" + baituploadpath + ".tar.gz"
	baiturl = util.Base64Encode(baiturl)
	DeleteAppBaitPolicyHandler(bait.TaskId, datas[0]["agentid"].(string), datas[0]["baitid"].(string), comm.BaitUNFile, baiturl, bait.Address, 4, filemd5)
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: "成功"})
}

func DownloadApplicationBaitPolicyHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	taskid := query.Get("taskid")

	if taskid == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	datas := honeycluster.SelectApplicationBaitsById(taskid)
	if len(datas) > 0 {
		baittype := util.Strval(datas[0]["type"])
		if baittype == "file" {
			baitname := datas[0]["baitname"].(string)
			baitfilename := datas[0]["baitfilename"].(string)
			filePath := util.GetCurrentPathString() + "/upload/" + baitname + "/" + baitfilename
			file, err := os.Open(filePath)
			if err != nil {
				logs.Error("sign download fail: %v , %s", err, filePath)
				return
			}

			defer file.Close()
			fileHeader := make([]byte, 65536)
			file.Read(fileHeader)
			fileStat, _ := file.Stat()
			w.Header().Set("Content-Disposition", "attachment; filename="+baitfilename)
			w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
			w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
			file.Seek(0, 0)
			io.Copy(w, file)
			return
		} else if baittype == "history" {
			baitcontent := util.Strval(datas[0]["data"])
			baitcontent = util.Base64Decode(baitcontent)
			var content = []byte(baitcontent)
			w.Header().Set("Content-Disposition", "attachment; filename= history")
			w.Header().Add("Content-Type", "application/octet-stream")
			w.Write(content)
			return
		}
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: datas, Message: "没有该诱饵下发任务"})
		return
	}
}

/**
新增诱饵
*/
func CreateApplicationBaitPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var bait comm.BaitJson
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &bait); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	if strings.Contains(bait.Address, "./") || strings.Contains(bait.Address, ".\\") {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.DirUseError})
		return
	}
	sysproto := ""
	if r.TLS == nil {
		sysproto = "http://"
	} else {
		sysproto = "https://"
	}
	if bait.Type == "file" {
		datas := honeycluster.SelectBaitsById(bait.BaitId)
		baitfilepath := "upload/" + datas[0]["baitname"].(string)
		util.ModifyFile(baitfilepath, bait.Address, datas[0]["baitname"].(string), datas[0]["baitinfo"].(string))
		currentpath := util.GetCurrentPathString()
		util.FileTarZip(baitfilepath, baitfilepath+".tar.gz")
		filemd5 := util.GetFileMd5(currentpath + "/" + baitfilepath + ".tar.gz")
		apphost := beego.AppConfig.String("apphost")
		appport := beego.AppConfig.String("appport")
		apphost = apphost + ":" + appport
		baiturl := sysproto + apphost + "/" + baitfilepath + ".tar.gz"
		baiturl = util.Base64Encode(baiturl)
		CreateBaitPolicyHandler(bait.AgentId, bait.BaitId, comm.BaitFile, baiturl, bait.Address, 3, filemd5)
	} else if bait.Type == "history" {
		baitdata := util.Base64Encode(bait.Data)
		CreateBaitPolicyHandler(bait.AgentId, "", comm.BaitHis, baitdata, "", 1, "")
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: "成功"})
	return
}

func CreateProtocolTypeHandler(agentID string, protocolType string, protocolData string, status int, fileMD5 string) {
	taskID := util.GetUUID()

	var protocolTypePolicyJson comm.ProtocolTypePolicyJson
	protocolTypePolicyJson.Status = status
	protocolTypePolicyJson.Md5 = fileMD5
	protocolTypePolicyJson.TaskId = taskID
	protocolTypePolicyJson.Type = protocolType
	protocolTypePolicyJson.Data = protocolData
	protocolTypePolicyJson.AgentId = agentID

	protocolPolicy, err3 := json.Marshal(protocolTypePolicyJson)
	if err3 != nil {
		return
	}
	redisCenter.RedisSubProducerBaitPolicy(string(protocolPolicy))
	return
}

/**
调用诱饵新增
*/
// 新增诱饵策略，insert数据库，下发策略到Redis,返回taskid
func CreateBaitPolicyHandler(agentid string, baitid string, baittype string, baitdata string, baitpath string, baitstatus int, filemd5 string) {

	taskid := util.GetUUID()

	// 前端json转换成策略结构体
	var baitPolicyJson comm.BaitPolicyJson
	baitPolicyJson.Status = baitstatus
	baitPolicyJson.Md5 = filemd5
	baitPolicyJson.TaskId = taskid
	baitPolicyJson.Type = baittype
	baitPolicyJson.Data = baitdata
	baitPolicyJson.AgentId = agentid
	baitPolicy, err3 := json.Marshal(baitPolicyJson)
	if err3 != nil {
		logs.Error("[CreateBaitPolicyHandler] 策略结构体转换失败 %v", err3)
		return
	}
	if baittype == comm.BaitFile || baittype == comm.BaitHis {
		if baittype == comm.BaitFile {
			baitdata = baitpath
			baittype = "file"
		} else if baittype == comm.BaitHis {
			baittype = "history"
		}
		createtime := time.Now().Unix()
		policyCenter.InsertBaitPolicy(taskid, agentid, baitid, "", createtime, comm.Creator, baitstatus, baitdata, filemd5, baittype)
	}

	redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
	return
}

// 删除协议转发 不包含下线功能
func RemoveHoneyTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var transJson comm.TransOfflineJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &transJson); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	policyCenter.DeleteHoneyTransPolicy(transJson.TaskId)
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.DataSelectSuccess})
	return
}

// 删除透明转发策略
func RemoveTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	// 前端json转换为结构体
	var transJson comm.TransOfflineJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}
	if err := json.Unmarshal([]byte(body), &transJson); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	policyCenter.DeleteTransPolicy(transJson.TaskId)
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.DataSelectSuccess})
	return
}

func DeleteProtocolTypePolicyHandler(agentID string, protocolType string, status int, protocolURL string, fileMD5 string) {

	// 前端json转换成策略结构体
	var protocolTypePolicyJson comm.ProtocolTypePolicyJson
	protocolTypePolicyJson.Status = status
	protocolTypePolicyJson.Md5 = fileMD5
	protocolTypePolicyJson.Type = protocolType
	protocolTypePolicyJson.AgentId = agentID
	protocolTypePolicyJson.Data = protocolURL

	baitPolicy, err3 := json.Marshal(protocolTypePolicyJson)
	if err3 != nil {
		logs.Error("[DeleteProtocolTypePolicyHandler] 策略结构体转换失败 %v", err3)
		return
	}

	redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
	return
}

/**
调用诱饵撤销
*/
// 新增诱饵策略，insert数据库，下发策略到Redis,返回taskid
func DeleteAppBaitPolicyHandler(taskid string, agentid string, baitid string, baittype string, baitdata string, baitpath string, baitstatus int, filemd5 string) {

	// 前端json转换成策略结构体
	var baitPolicyJson comm.BaitPolicyJson
	baitPolicyJson.Status = baitstatus
	baitPolicyJson.Md5 = filemd5
	baitPolicyJson.TaskId = taskid
	baitPolicyJson.Type = baittype
	baitPolicyJson.Data = baitdata
	baitPolicyJson.AgentId = agentid

	baitPolicy, err3 := json.Marshal(baitPolicyJson)
	if err3 != nil {
		logs.Error("[CreateBaitPolicyHandler] 策略结构体转换失败 %v", err3)
		return
	}
	offlinetime := time.Now().Unix()
	policyCenter.UpdateBaitPolicy(taskid, offlinetime, baitstatus)

	redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
	return
}

/**
应用服务器密签策略列表
*/
func SelectSignPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var selectSign comm.SelectSignPolicyJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &selectSign); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	datas, msg, code := policyCenter.SelectSignPolicy(selectSign.AgentId, selectSign.SignId, selectSign.SignType, selectSign.Creator, selectSign.CreateStartTime, selectSign.CreateEndTime, selectSign.OfflineStartTime, selectSign.OfflineEndTime, selectSign.Status, selectSign.PageSize, selectSign.PageNum, selectSign.SignInfo)
	comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: datas, Message: msg})
}

/**
applications新增密签
*/
func CreateApplicationSignPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var sign comm.SignJson
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &sign); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	if strings.Contains(sign.Address, "./") || strings.Contains(sign.Address, ".\\") {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.DirUseError})
		return
	}
	sysproto := ""
	signname := ""
	traceurl := ""
	taskid := util.GetUUID()
	if r.TLS == nil {
		sysproto = "http://"
	} else {
		sysproto = "https://"
	}
	if sign.Type == "file" {
		datas := honeycluster.SelectSignsById(sign.SignId)
		if len(datas) == 0 {
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "密签不存在!"})
			return
		}
		signname = datas[0]["signname"].(string)
		signsourcepath := "upload/honeytoken/" + signname
		signfilepath := "upload/honeytoken/" + taskid + "/" + signname
		err := util.MakeDir(signfilepath)
		if err != nil {
			logs.Error("[CreateHoneySignPolicyHandler] signfilepath make err:", err)
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "密签目录创建失败"})
			return
		}
		tracecode := util.GetUUID()
		traceinfo := honeycluster.GetTraceHost()
		if len(traceinfo) == 0 {
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "请到系统设置配置密签追踪url"})
			return
		}
		tracehost := util.Strval(traceinfo[0]["tracehost"])
		if strings.Contains(tracehost, "http") {
			traceurl = fmt.Sprintf("%s/api/msgreceive?tracecode=%s", tracehost, tracecode)
		} else {
			traceurl = fmt.Sprintf("http://%s/api/msgreceive?tracecode=%s", tracehost, tracecode)
		}
		signinfoname := datas[0]["signinfo"].(string)
		var Suffix = regexp.MustCompile(`.(ppt|pptx|doc|docx|pdf|xls|xlsx)$`)
		if len(Suffix.FindString(signinfoname)) > 0 {

			err := honeytoken.DoFileSignTrace(signinfoname, signinfoname, signsourcepath, signfilepath, traceurl)
			if err != nil {
				error := util.CopyDir(signsourcepath, signfilepath)
				if error != nil {
					logs.Error("[CreateApplicationSignPolicyHandler] util.CopyDir err:", error)
				}
			}
		}
		util.ModifyFile(signfilepath, sign.Address, datas[0]["signname"].(string), datas[0]["signinfo"].(string))
		currentpath := util.GetCurrentPathString()
		util.FileTarZip(signfilepath, signfilepath+".tar.gz")

		filemd5 := util.GetFileMd5(currentpath + "/" + signfilepath + ".tar.gz")
		apphost := beego.AppConfig.String("apphost")
		appport := beego.AppConfig.String("appport")
		apphost = apphost + ":" + appport
		baiturl := sysproto + apphost + "/" + signfilepath + ".tar.gz"
		baiturl = util.Base64Encode(baiturl)
		go CreateSignPolicyHandler(taskid, sign.AgentId, sign.SignId, comm.SignFile, baiturl, sign.Address+"/"+datas[0]["signinfo"].(string), 3, filemd5, tracecode)
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: "成功"})
}

/**
探针密签管理-下载
*/
func DownloadApplicationSignPolicyHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	taskid := query.Get("taskid")

	if taskid == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	datas := honeycluster.SelectApplicationSignsById(taskid)
	signname := datas[0]["signname"].(string)
	signfilename := datas[0]["signfilename"].(string)
	filePath := util.GetCurrentPathString() + "/upload/honeytoken/" + taskid + "/" + signname + "/" + signfilename
	file, err := os.Open(filePath)
	if err != nil {
		logs.Error("sign download fail: %v , %s", err, filePath)
		return
	}
	defer file.Close()
	fileHeader := make([]byte, 1024)
	file.Read(fileHeader)
	fileStat, _ := file.Stat()
	w.Header().Set("Content-Disposition", "attachment; filename="+signfilename)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	file.Seek(0, 0)
	io.Copy(w, file)
	return

}

func DeleteApplicationSignPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var signs comm.SignJson
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &signs); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	sysproto := ""
	if r.TLS == nil {
		sysproto = "http://"
	} else {
		sysproto = "https://"
	}
	datas := honeycluster.SelectApplicationSignsById(signs.TaskId)
	signfilepath := datas[0]["signinfo"].(string)
	signname := datas[0]["signname"].(string) + "-uninstall"

	util.ModifySignUninstallFile(signfilepath, signname)
	signuploadpath := "upload/honeysign/" + signname
	currentpath := util.GetCurrentPathString()
	util.FileTarZip(signuploadpath, signuploadpath+".tar.gz")
	filemd5 := util.GetFileMd5(currentpath + "/" + signuploadpath + ".tar.gz")
	apphost := beego.AppConfig.String("apphost")
	appport := beego.AppConfig.String("appport")
	apphost = apphost + ":" + appport
	baiturl := sysproto + apphost + "/" + signuploadpath + ".tar.gz"
	baiturl = util.Base64Encode(baiturl)
	go DeleteSignPolicyHandler(datas[0]["agentid"].(string), signs.TaskId, comm.SignUNFile, baiturl, signs.Address, 1, filemd5, "")
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: "成功"})

}

func ApplicationSignMsgHandler(w http.ResponseWriter, r *http.Request) {
	var fileList struct {
		Tracecode string `json:"tracecode"`
		Map       string `json:"map"`
		PageNum   int    `json:"pageNum"`
		PageSize  int    `json:"pageSize"`
		StartTime string `json:"starttime"`
		EndTime   string `json:"endtime"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &fileList); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.ApplicationSelectSignMsg(fileList.Tracecode, fileList.Map, fileList.StartTime, fileList.EndTime, fileList.PageSize, fileList.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: datas, Message: "成功"})
	return
}

func DeleteSignPolicyHandler(agentid string, taskid string, signtype string, signdata string, signpath string, signstatus int, filemd5 string, tracecode string) {

	// 前端json转换成策略结构体
	var baitPolicyJson comm.BaitPolicyJson
	baitPolicyJson.Status = signstatus
	baitPolicyJson.Md5 = filemd5
	baitPolicyJson.TaskId = taskid
	baitPolicyJson.Type = signtype
	baitPolicyJson.Data = signdata
	baitPolicyJson.AgentId = agentid
	baitPolicy, err3 := json.Marshal(baitPolicyJson)
	if err3 != nil {
		logs.Error("[DeleteSignPolicyHandler] 策略结构体转换失败 %v", err3)
		return
	}
	if strings.ToLower(signtype) == "un_file" {
		signdata = signpath
		signtype = "file"
		offlinetime := time.Now().Unix()
		status := 4
		policyCenter.OffSignPolicy(taskid, status, offlinetime)
	}
	redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
	return
}

func CreateSignPolicyHandler(taskid string, agentid string, signid string, signtype string, signdata string, signpath string, signstatus int, filemd5 string, tracecode string) {

	// 前端json转换成策略结构体
	var baitPolicyJson comm.BaitPolicyJson
	baitPolicyJson.Status = signstatus
	baitPolicyJson.Md5 = filemd5
	baitPolicyJson.TaskId = taskid
	baitPolicyJson.Type = signtype
	baitPolicyJson.Data = signdata
	baitPolicyJson.AgentId = agentid
	baitPolicy, err3 := json.Marshal(baitPolicyJson)
	if err3 != nil {
		logs.Error("[CreateSignPolicyHandler] 策略结构体转换失败 %v", err3)
		return
	}

	if signtype == comm.SignFile {
		signdata = signpath
		signtype = "file"
		createtime := time.Now().Unix()
		policyCenter.InsertSignPolicy(taskid, agentid, signid, "", createtime, comm.Creator, signstatus, signdata, filemd5, signtype, tracecode)
	}
	redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
	return
}

// 查询所有诱饵类型
func SelectAllBaitTypeHandler(w http.ResponseWriter, r *http.Request) {
	results := policyCenter.SelectAllBaitType()
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: results, Message: comm.DataSelectSuccess})
	return
}

func SelectAllHoneyBaitTypeHandler(w http.ResponseWriter, r *http.Request) {
	results := policyCenter.SelectAllHoneyBaitType()
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: results, Message: comm.DataSelectSuccess})
	return
}

// 查询所有操作类型
func SelectAllSysTypeHandler(w http.ResponseWriter, r *http.Request) {
	results := policyCenter.SelectAllSysType()
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: results, Message: comm.DataSelectSuccess})
	return
}

// 删除诱饵策略，update数据库，下发策略到Redis
func DeleteBaitPolicyHandler(w http.ResponseWriter, r *http.Request) {
	// 前端json转换为结构体
	var bait comm.BaitPolicyJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}
	if err := json.Unmarshal([]byte(body), &bait); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	if bait.Status == 1 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "请确认诱饵策略的状态"})

		return
	}
	// 前端json转换成策略结构体
	var baitPolicyJson comm.BaitPolicyJson

	err1 := json.Unmarshal([]byte(body), &baitPolicyJson)
	if err1 != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})

		return
	}
	currentTime := time.Now()
	//offlineTime := currentTime.Format("2006-01-02 15:04:05")
	policyCenter.UpdateBaitPolicy(bait.TaskId, currentTime.Unix(), bait.Status)
	//策略结构体转换成策略json
	baitPolicy, err2 := json.Marshal(baitPolicyJson)
	if err2 != nil {
		logs.Error(fmt.Sprintf("[DeleteBaitPolicyHandler] 策略结构体转换失败 %v", err1))
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.DataSelectSuccess})
	return
}

// 新增透明转发策略，insert数据库，下发策略到Redis
func CreateTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var transJson comm.TransparentTransponderJson
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &transJson); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	if transJson.ListenPort < 1025 || transJson.ListenPort > 65535 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.PortUseError})
		return
	}
	var transPolicyJson comm.TransparentTransponderPolicyJson
	taskid := util.GetUUID()
	forwardport := transJson.ForwardPort
	honeyforward := policyCenter.SelectHoneyForwardByHoneyId(transJson.HoneyPotId)
	honeyType := policyCenter.SelectHoneyPotTypeByHoneyId(transJson.HoneyPotId)
	serverInfo := policyCenter.SelectHoneyServerInfo(transJson.HoneyPotId)

	if len(serverInfo) == 0 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("蜜罐服务器获取异常")})
		return
	}
	if len(honeyforward) == 0 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("蜜罐端口获取异常")})
		return
	}
	if len(honeyType) == 0 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("蜜罐类型获取异常")})
		return
	}

	honeyserveragentid := honeyforward[0]["agentid"].(string)
	honeytypeid := honeyType[0]["honeytypeid"].(string)
	servertype := util.Strval(honeyType[0]["honeypottype"])

	forwardinfo := policyCenter.SelectForwardByConditions(transJson.AgentId, honeytypeid, forwardport, transJson.ListenPort, honeyserveragentid)
	if len(forwardinfo) > 0 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "透明转发已经创建"})
		return
	} else {
		transPolicyJson.ServerType = servertype
		transPolicyJson.HoneyIP = serverInfo[0]["serverip"].(string)
		transPolicyJson.AgentId = transJson.AgentId
		transPolicyJson.TaskId = taskid
		transPolicyJson.HoneyPort = forwardport
		transPolicyJson.ListenPort = transJson.ListenPort
		transPolicyJson.Status = 1
		transPolicyJson.Type = comm.AgentTypeEdge
		transPolicyJson.Path = util.Strval(honeyType[0]["softpath"])
		policyCenter.InsertTransPolicy(taskid, transJson.AgentId, transJson.ListenPort, forwardport, honeyserveragentid, time.Now().Unix(), comm.Creator, 3, comm.AgentTypeEdge, transPolicyJson.Path, honeytypeid)
		policyCenter.UpdateHoneyTransPolicyForwardStatus(taskid)
		// 策略结构体转换成策略json
		transPolicy, err2 := json.Marshal(transPolicyJson)
		if err2 != nil {
			logs.Error(fmt.Sprintf("[CreateTransPolicyHandler] Failed to marshal data: %v", err2))
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "透明代理数据转换异常"})
			return
		}
		go redisCenter.RedisSubProducerTransPolicy(string(transPolicy))
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.DataSelectSuccess})
		return
	}

}

// 删除透明转发策略，insert数据库，下发策略到Redis
func DeleteTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	// 前端json转换为结构体
	var transJson comm.TransOfflineJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}
	if err := json.Unmarshal([]byte(body), &transJson); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	currentTime := time.Now().Unix()
	policyCenter.UpdateTransPolicy(transJson.TaskId, currentTime, 2)
	policy := policyCenter.SelectHoneyPotInfoByTaskId(transJson.TaskId)
	go redisCenter.RedisSubProducerTransPolicy(policy)
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.DataSelectSuccess})
	return
}

// 查看指定agent的所有诱饵策略
func SelectBaitPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var selectBait comm.SelectBaitPolicyJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &selectBait); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	datas, msg, code := policyCenter.SelectBaitPolicy(selectBait.AgentId, selectBait.BaitId, selectBait.Creator, selectBait.CreateStartTime, selectBait.CreateEndTime, selectBait.OfflineStartTime, selectBait.OfflineEndTime, selectBait.Status, selectBait.PageSize, selectBait.PageNum, selectBait.BaitInfo)
	comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: datas, Message: msg})

	return
}

// 查看指定agent的所有透明转发策略
func SelectTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var agentId comm.SelectTransPolicyJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}
	if err := json.Unmarshal([]byte(body), &agentId); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	datas, msg, code := policyCenter.SelectTransPolicy(agentId.ServerIP, agentId.AgentId, agentId.ForwardPort, agentId.HoneyPotPort, agentId.CreateStartTime, agentId.CreateEndTime, agentId.OfflineStartTime, agentId.OfflineEndTime, agentId.Creator, agentId.Status, agentId.HoneyTypeId, agentId.PageSize, agentId.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: datas, Message: msg})

	return
}

/**
透明转发网络探测
*/
func TestTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var transJson comm.TransTestJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &transJson); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	transInfo := policyCenter.SelectForwardByTaskId(transJson.TaskId)
	if len(transInfo) > 0 {
		iscon := false
		ips := util.Strval(transInfo[0]["serverip"])
		iplist := strings.Split(ips, ",")
		forwardport := util.Strval(transInfo[0]["forwardport"])
		if len(iplist) > 0 {
			for _, ip := range iplist {
				if util.NetConnectTest(ip, forwardport) {
					iscon = true
					break
				}
			}
		}
		if iscon {
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.NetWorkSuccess})
			return
		} else {
			currentTime := time.Now().Unix()
			policyCenter.UpdateTransPolicy(transJson.TaskId, currentTime, 2)
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.NetWorkfail})
			return
		}
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "未查到该taskid的透明转发信息"})
		return
	}
}

// 创建蜜罐转发
func CreateHoneyTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var honeytrans comm.HoneyTrans
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeytrans); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	if honeytrans.ListenPort < 1025 || honeytrans.ListenPort > 65535 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.PortUseError})
		return
	}
	var transPolicyJson comm.TransparentTransponderPolicyJson
	taskid := util.GetUUID()
	honeytraninfo := policyCenter.SelectHoneyTransByHoneyId(honeytrans.HoneyPotId)
	if len(honeytraninfo) > 0 {
		honeyport, _ := strconv.Atoi(util.Strval(honeytraninfo[0]["honeyport"]))
		agentId := util.Strval(honeytraninfo[0]["agentid"])

		honeyforwardinfo := policyCenter.SelectHoneyForwardByConditions(honeytrans.HoneyPotId, honeyport, agentId, honeytrans.ListenPort)
		if len(honeyforwardinfo) > 0 {
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "协议转发已经创建"})
			return
		} else {
			sysproto := ""
			transPolicyJson.ServerType = util.Strval(honeytraninfo[0]["honeypottype"])
			transPolicyJson.HoneyIP = util.Strval(honeytraninfo[0]["honeyip"])
			transPolicyJson.AgentId = agentId
			transPolicyJson.TaskId = taskid
			transPolicyJson.HoneyPort = honeyport
			transPolicyJson.ListenPort = honeytrans.ListenPort
			transPolicyJson.Status = 1
			transPolicyJson.Type = comm.AgentTypeRelay
			transPolicyJson.Path = util.Strval(honeytraninfo[0]["softpath"])
			if r.TLS == nil {
				sysproto = "http://"
			} else {
				sysproto = "https://"
			}
			apphost := beego.AppConfig.String("apphost")
			appport := beego.AppConfig.String("appport")
			apphost = apphost + ":" + appport
			transPolicyJson.SecCenter = sysproto + apphost

			createtime := time.Now().Unix()
			policyCenter.InsertHoneyTransPolicy(taskid, transPolicyJson.AgentId, transPolicyJson.ListenPort, honeyport, honeytrans.HoneyPotId, createtime, "admin", 3, comm.AgentTypeRelay, transPolicyJson.Path, util.Strval(honeytraninfo[0]["serverid"]))
			transPolicy, err2 := json.Marshal(transPolicyJson)
			if err2 != nil {
				logs.Error(fmt.Sprintf("[CreateHoneyTransPolicyHandler] Failed to marshal data: %v", err2))
				return
			}
			if strings.ToLower(transPolicyJson.ServerType) == "ssh" {
				sshkey, errmsg := policyCenter.SelectHoneyServerSSHKey(transPolicyJson.AgentId)
				if errmsg != "" {
					comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: nil, Message: errmsg})
					return
				}
				if sshkey == "" {
					comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "蜜网的SSH Key为空，请重新安装agent"})
					return
				} else {
					sshkeycheckcmdstr := []string{"/bin/sh", "-c", "find /root/.ssh"}
					checkerr := k3s.ExecPodCmd(honeytrans.HoneyPotId, sshkeycheckcmdstr)
					if checkerr != nil {
						cmdstr := []string{"mkdir", "/root/.ssh"}
						makeerr := k3s.ExecPodCmd(honeytrans.HoneyPotId, cmdstr)
						if makeerr != nil {
							logs.Error("[CreateHoneyTransPolicyHandler] makedir /root/.ssh err:", makeerr)
							comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: "makedirsshkeyerr", Message: "ssh 蜜罐创建失败"})
							return
						} else {
							cmdstr = []string{"/bin/sh", "-c", "echo " + sshkey + " |  base64 -d > /root/.ssh/authorized_keys"}
							echosshkeyerr := k3s.ExecPodCmd(honeytrans.HoneyPotId, cmdstr)
							if echosshkeyerr != nil {
								logs.Error("[CreateHoneyTransPolicyHandler] write sshkey err:", echosshkeyerr)
								comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: "echosshkeyerr", Message: "ssh 蜜罐创建失败"})
								return
							}
						}
					}
				}
			}
			go redisCenter.RedisSubProducerTransPolicy(string(transPolicy))
			comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: nil, Message: "成功"})
			return
		}
	}
}

// 下线协议转发
func DeleteHoneyTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var transJson comm.TransOfflineJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &transJson); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	currentTime := time.Now().Unix()
	policyCenter.UpdateHoneyTransPolicy(transJson.TaskId, currentTime, 2)
	policy := policyCenter.SelectHoneyPotInfoByTaskIdHoneyTrans(transJson.TaskId)
	go redisCenter.RedisSubProducerTransPolicy(policy)
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.DataSelectSuccess})
	return
}

func TestHoneyTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var transJson comm.TransTestJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &transJson); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	transInfo := policyCenter.SelectHoneyForwardByTaskId(transJson.TaskId)
	if len(transInfo) > 0 {
		honeypotname := util.Strval(transInfo[0]["podname"])
		honeynamespce := util.Strval(transInfo[0]["honeynamespce"])
		_, _, error, _, errorcode := k3s.CheckPodStatus(honeypotname, honeynamespce)
		if errorcode == 5 {
			iscon := false
			ips := util.Strval(transInfo[0]["serverip"])
			iplist := strings.Split(ips, ",")
			forwardport := util.Strval(transInfo[0]["forwardport"])
			if len(iplist) > 0 {
				for _, ip := range iplist {
					if util.NetConnectTest(ip, forwardport) {
						iscon = true
						break
					}
				}
			}
			if iscon {
				comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.NetWorkSuccess})
				return
			} else {
				offlinetime := time.Now().Unix()
				policyCenter.UpdateHoneyTransPolicy(transJson.TaskId, offlinetime, 2)
				comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.NetWorkfail})
				return
			}
		} else {
			offlinetime := time.Now().Unix()
			honeycluster.UpdatePod("", util.Strval(transInfo[0]["honeypotid"]), offlinetime)
			policyCenter.UpdateHoneyTransPolicy(transJson.TaskId, offlinetime, 2)
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: util.Strval(transInfo[0]["honeyname"]) + " 状态异常： " + error.Error()})
			return
		}
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "未查到该taskid的协议转发信息"})
		return
	}
}

// 查询蜜罐转发列表
func SelectHoneyTransPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var honeytrans comm.SelectHoneypotTransPolicyJson
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}
	if err := json.Unmarshal([]byte(body), &honeytrans); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas, msg, code := policyCenter.SelectHoneyPotTrans(honeytrans.HoneyPotPort, honeytrans.HoneyPotId, honeytrans.AgentId, honeytrans.Creator, honeytrans.ForwardPort, honeytrans.HoneyTypeId, honeytrans.HoneyPotIp, honeytrans.Status, honeytrans.CreateStartTime, honeytrans.CreateEndTime, honeytrans.OfflineStartTime, honeytrans.OfflineEndTime, honeytrans.PageSize, honeytrans.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: datas, Message: msg})
	return

}

// 查询所有蜜罐类型
func GetHoneyPotsType(w http.ResponseWriter, r *http.Request) {
	var agentId comm.SelectTransPolicyJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}
	if err := json.Unmarshal([]byte(body), &agentId); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	results := policyCenter.SelectAllHoneyPotsType()
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: results, Message: comm.DataSelectSuccess})
	return
}

// 查询所有存在代理转发的蜜罐服务器IP、对应的监听端口、蜜罐流量转发策略状态为1、蜜罐状态为1
func GetHoneyTransInfos(w http.ResponseWriter, r *http.Request) {
	var typeId comm.HoneyPotsType
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}
	if err := json.Unmarshal([]byte(body), &typeId); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	datas, msg, code := policyCenter.SelectHoneyPotsTransInfoByHoneyTypeId(typeId.HoneyTypeId)
	comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: datas, Message: msg})

	return
}

// 查询攻击事件列表
func SelectAttackLogListHandler(w http.ResponseWriter, r *http.Request) {
	var attackLogList comm.AttackLogListJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &attackLogList); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas, msg, code := policyCenter.SelectAttackLogList(attackLogList.ServerName, attackLogList.HoneyTypeId, attackLogList.SrcHost, attackLogList.AttackIP, attackLogList.StartTime, attackLogList.EndTime, attackLogList.PageSize, attackLogList.PageNum, attackLogList.HoneyIP)
	comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: datas, Message: msg})
	return
}

// 查询攻击事件详情
func SelectAttackLogDetailHandler(w http.ResponseWriter, r *http.Request) {
	var attackLogDetail comm.AttackLogDetailJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}
	if err := json.Unmarshal([]byte(body), &attackLogDetail); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas, msg, code := policyCenter.SelectAttackLogDetail(attackLogDetail.Id, attackLogDetail.HoneyTypeId, attackLogDetail.SrcHost, attackLogDetail.AttackIP, attackLogDetail.StartTime, attackLogDetail.EndTime, attackLogDetail.HoneyPotPort, attackLogDetail.HoneyIP, attackLogDetail.PageSize, attackLogDetail.PageNum, attackLogDetail.SrcPort, attackLogDetail.EventDetail)
	comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: datas, Message: msg})
	return
}

// 接收falco的日志
func InsertFalcoLogHandler(w http.ResponseWriter, r *http.Request) {
	var falcoLog struct {
		Output       string            `json:"output"`   // HoneyIP
		Priority     string            `json:"priority"` // HoneyPotPort
		Rule         string            `json:"rule"`     // serverIP
		Time         string            `json:"time"`     // serverPort
		OutputFields comm.OutputFields `json:"output_fields"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}

	if err := json.Unmarshal([]byte(body), &falcoLog); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	if falcoLog.OutputFields.ContainerId == "" || falcoLog.OutputFields.ContainerId == "host" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertFalcoLogHandler] Invalid log")})
		return
	}
	if falcoLog.OutputFields.FileName != "" {
		kubeconfig := beego.AppConfig.String("kubeconfig")
		kubeconfig = util.GetCurrentPathString() + "/" + kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logs.Error("[InsertFalcoLogHandler] clientcmd.BuildConfigFromFlags ERR:", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			logs.Error("[InsertFalcoLogHandler] kubernetes.NewForConfig ERR:", err)
		}
		reader, outStream := io.Pipe()
		srcPath := falcoLog.OutputFields.FileName
		destPath := beego.AppConfig.String("clamavpath")
		cmdArr := []string{"tar", "cf", "-", srcPath}

		req := clientset.CoreV1().RESTClient().Get().
			Resource("pods").
			Name(falcoLog.OutputFields.K8sPodName).
			Namespace("default").
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Command: cmdArr,
				Stdin:   true,
				Stdout:  true,
				Stderr:  true,
				TTY:     false,
			}, scheme.ParameterCodec)
		exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())

		go func() {
			defer outStream.Close()
			err = exec.Stream(remotecommand.StreamOptions{
				Stdin:  os.Stdin,
				Stdout: outStream,
				Stderr: os.Stderr,
				Tty:    false,
			})
		}()
		prefix := util.GetPrefix(srcPath)
		prefix = path.Clean(prefix)
		destPath = path.Join(destPath, path.Base(prefix))
		err = util.UntarAll(reader, destPath, prefix)
		if err != nil {
			logs.Error("kubectl cp file err:", err)
		}

	}
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.DBInsertSuccess})
	return

	//db1, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	//if err != nil {
	//	log.Printf("open mysql fail %v\n", err)
	//	comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertFalcoLogHandler] open mysql fail %v", err)})
	//	return
	//}
	//defer db1.Close()
	//
	////DbCon = sqlCon
	////defer sqlCon.Close()
	//
	////log.Printf("start web server\n")
	//db1.SetConnMaxLifetime(300 * time.Second)
	//db1.SetMaxOpenConns(120)
	//db1.SetMaxIdleConns(12)
	//
	//honeyPotid0 := ""
	//honeytypeid := ""
	//honeypotport := ""
	//honeyip := ""
	//serverid := ""
	//if len(maps) > 0 {
	//	honeyPotid0 = util.Strval(maps[0]["honeypotid"])
	//	honeytypeid = util.Strval(maps[0]["honeytypeid"])
	//	honeypotport = util.Strval(maps[0]["honeyport"])
	//	honeyip = util.Strval(maps[0]["honeyip"])
	//	serverid = util.Strval(maps[0]["serverid"])
	//}
	//country := "局域网"
	//province := "局域网"
	//if honeyip != "" {
	//	data := models.GetGeoDataForAliYun(honeyip)
	//	if data.Country != "" {
	//		if data.CountryCode == "CN" {
	//			country = data.Country
	//			province = data.Province
	//			//site = data.Province + "-" + data.City
	//		} else if data.Country == "局域网" {
	//			country = data.Country
	//			province = "局域网"
	//		} else {
	//			if data.City != "" && data.Province != "" {
	//				country = data.Country
	//				province = data.Province
	//				//site = data.Country + data.Province + "-" + data.City
	//			} else {
	//				datas := models.GetGeoData(honeyip)
	//				log.Println(datas)
	//				if datas.CountryCode == "unk" {
	//					country = ""
	//					province = ""
	//				} else {
	//					country = data.Country
	//					province = datas.RegionName
	//					//site = data.Country + "-" + datas.RegionName + "-" + datas.City
	//				}
	//			}
	//		}
	//	} else {
	//		datas := models.GetGeoData(honeyip)
	//		if datas.CountryCode == "unk" {
	//			country = ""
	//			province = ""
	//		} else if datas.CountryCode == "CN" {
	//			country = datas.CountryName
	//			province = datas.RegionName
	//			//site = datas.RegionName + "-" + datas.City
	//		} else {
	//			country = datas.CountryName
	//			province = datas.RegionName
	//			//site = datas.CountryName + "-" + datas.RegionName + "-" + datas.City
	//		}
	//
	//	}
	//}

	//attacktime := time.Now().Unix()
	//eventDetail := "以" + falcoLog.OutputFields.UserName + "权限执行" + falcoLog.OutputFields.Command
	//res, err3 := db1.Exec("insert INTO attacklog(serverid,honeypotid,honeytypeid,honeypotport,attackip,country, province,attacktime,eventdetail,proxytype,sourcetype,logdata) values(?,?,?,?,?,?,?,?,?,?,?,?)", serverid, honeyPotid0, honeytypeid, honeypotport, honeyip, country, province, attacktime, eventDetail, "falco", 2, body)
	//
	//if err3 != nil {
	//	log.Printf("insert mysql fail: ", err3)
	//	comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertFalcoLogHandler] insert mysql fail: %v", err3)})
	//	return
	//}
	//
	//_, err = res.RowsAffected()
	//if err != nil {
	//	log.Printf("insert failed, err: ", err)
	//	comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertFalcoLogHandler] insert failed, err:%v", err)})
	//	return
	//}
}

//comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.DBInsertSuccess})
//return

//}

// 添加钉钉告警
func AddConfigHandler(w http.ResponseWriter, r *http.Request) {
	var config comm.ConfJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &config); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	if config.ConfName != "DingDing" && config.ConfName != "Message" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf(" 类型错误")})
		return
	} else {
		result, msg, code := policyCenter.InsertConfig(config.ConfName, config.ConfValue)
		comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: result, Message: msg})

		return
	}

}

// 查询钉钉告警
func SelectConfigHandler(w http.ResponseWriter, r *http.Request) {
	var config comm.ConfJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}
	if err := json.Unmarshal([]byte(body), &config); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	if config.ConfName != "DingDing" && config.ConfName != "Message" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf(" 类型错误")})
		return

	} else {
		result, msg, code := policyCenter.SelectConfig(config.Id, config.ConfName, config.PageSize, config.PageNum)
		comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: result, Message: msg})

		return
	}
}

func CreateProtocolType(w http.ResponseWriter, r *http.Request) {
	protocolFileName, protocolName, _ := ProtocolFileUpload(r)
	if protocolFileName == "" || protocolName == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "协议模块上传信息不完整"})
		return
	}
	typeid := util.GetStrMd5(protocolName)
	softPath := "/home/ehoney_proxy/" + protocolFileName
	createTime := time.Now().Unix()
	data, msg, msgCode := honeycluster.InsertProtocol(protocolName, typeid, softPath, createTime)

	util.ModifyFile("upload/protocol/"+protocolName, "/home/sys_admin/", protocolName, protocolFileName)
	currentpath := util.GetCurrentPathString()
	protocolPath := "upload/protocol/" + protocolName
	util.FileTarZip(protocolPath, protocolPath+".tar.gz")
	fileMD5 := util.GetFileMd5(currentpath + "/" + protocolPath + ".tar.gz")
	sysProto := ""
	if r.TLS == nil {
		sysProto = "http://"
	} else {
		sysProto = "https://"
	}
	apphost := beego.AppConfig.String("apphost")
	appport := beego.AppConfig.String("appport")
	apphost = apphost + ":" + appport
	protocolURL := sysProto + apphost + "/" + protocolPath + ".tar.gz"
	protocolURL = util.Base64Encode(protocolURL)

	agentID := honeycluster.SelectHoneypotServer()[0]["agentid"].(string)

	CreateProtocolTypeHandler(agentID, "FILE", protocolURL, 1, fileMD5)
	comhttp.SendJSONResponse(w, comm.Response{Code: msgCode, Data: data, Message: msg})
	return
}

func GetProtocolType(w http.ResponseWriter, r *http.Request) {
	var protocoltype comm.ProtocolType
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &protocoltype); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectProtocol(protocoltype.PageSize, protocoltype.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func DeleteProtocolType(w http.ResponseWriter, r *http.Request) {
	var protocoltype comm.ProtocolType
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &protocoltype); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	sysproto := ""
	if r.TLS == nil {
		sysproto = "http://"
	} else {
		sysproto = "https://"
	}

	protocolFilePath := honeycluster.SelectHoneypotTypeByID(protocoltype.TypeId)[0]["softpath"].(string)
	protocolName := protocolFilePath[strings.LastIndex(protocolFilePath, "/")+1:]
	protocolName = protocolName + "-uninstall"
	util.ModifyUninstallFile(protocolFilePath, protocolName)
	protocolUploadPath := "upload/protocol/" + protocolName
	currentpath := util.GetCurrentPathString()
	util.FileTarZip(protocolUploadPath, protocolUploadPath+".tar.gz")
	fileMD5 := util.GetFileMd5(currentpath + "/" + protocolUploadPath + ".tar.gz")
	apphost := beego.AppConfig.String("apphost")
	appport := beego.AppConfig.String("appport")
	apphost = apphost + ":" + appport
	protocolURL := sysproto + apphost + "/" + protocolUploadPath + ".tar.gz"
	protocolURL = util.Base64Encode(protocolURL)
	agentID := honeycluster.SelectHoneypotServer()[0]["agentid"].(string)
	data, msg, msgCode := honeycluster.DeleteProtocol(protocoltype.TypeId)
	DeleteProtocolTypePolicyHandler(agentID, "UN_FILE", 2, protocolURL, fileMD5)
	comhttp.SendJSONResponse(w, comm.Response{Code: msgCode, Data: data, Message: msg})
	return
}

func GetHoneyImageLists(w http.ResponseWriter, r *http.Request) {
	var honeyimage comm.HoneyImage
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeyimage); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectHoneyImage(honeyimage.ImagesName, honeyimage.ImageType, honeyimage.ImageOS, honeyimage.PageSize, honeyimage.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func UpdateHoneyImage(w http.ResponseWriter, r *http.Request) {
	var honeyimage comm.HoneyImage
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeyimage); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	if honeyimage.ID < 6 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "无法修改默认镜像模板"})
		return
	}
	if honeyimage.ImagePort < 0 || honeyimage.ImagePort > 65535 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "端口范围为0-65535"})
		return
	}
	err := honeycluster.UpdateHoneyImageById(honeyimage.ID, honeyimage.ImageType, honeyimage.ImageOS, honeyimage.ImagePort)
	if err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "镜像信息更新失败"})
		return
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: "镜像信息更新成功"})
	return
}

func DeleteConfigHandler(w http.ResponseWriter, r *http.Request) {
	var config comm.ConfJson
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &config); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	if config.ConfName != "DingDing" && config.ConfName != "Message" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf(" 类型错误")})
		return

	} else {
		result, msg, code := policyCenter.DeleteConfig(config.Id, config.ConfName)
		comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: result, Message: msg})

		return
	}
}

func AddHarborConfig(w http.ResponseWriter, r *http.Request) {
	var harborinfo struct {
		HarborUrl   string `json:"harborUrl"`
		UserName    string `json:"userName"`
		Password    string `json:"password"`
		ProjectName string `json:"projectName"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &harborinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	harborurl := util.GetHost(harborinfo.HarborUrl)
	if harborurl == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "harborurl 请输入标准格式：http://ip/"})
		return
	}
	harborInfo := honeycluster.SelectHarborInfo()
	if len(harborInfo) > 0 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "harbor配置只支持单条记录"})
		return
	}
	harborip := ""
	harborport := ""
	harborhost := util.GetHost(harborinfo.HarborUrl)
	if harborhost != "" {
		harborhostinfo := strings.Split(harborhost, ":")
		if len(harborhostinfo) == 1 {
			harborip = harborhostinfo[0]
			harborport = "80"
		} else {
			harborip = harborhostinfo[0]
			harborport = harborhostinfo[1]
		}
	} else {
		err := fmt.Errorf("%s", "harhost解析异常")
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: err.Error()})
		return
	}
	if !util.NetConnectTest(harborip, harborport) {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "harbor url网络异常"})
		return
	}

	createtime := time.Now().Unix()
	harborid := util.GetStrMd5("ehoney")
	data, msg, msgcode := honeycluster.InsertHarborInfo(harborid, harborinfo.HarborUrl, harborinfo.UserName, harborinfo.Password, harborinfo.ProjectName, createtime)
	go GetPodImageList()
	comhttp.SendJSONResponse(w, comm.Response{Code: msgcode, Data: data, Message: msg})
	return
}

func GetHarborInfo(w http.ResponseWriter, r *http.Request) {
	datas := honeycluster.SelectHarborInfo()
	if len(datas) > 0 {
		datas[0]["password"] = "***"
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func DeleteHarborInfo(w http.ResponseWriter, r *http.Request) {
	var harborinfo comm.HarborInfo
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &harborinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	data, msg, msgCode := honeycluster.DeleteHarborInfo(harborinfo.HarborId)
	comhttp.SendJSONResponse(w, comm.Response{Code: msgCode, Data: data, Message: msg})
	return
}

func AddRedisConfigHandler(w http.ResponseWriter, r *http.Request) {
	var redisinfo struct {
		RedisIP   string `json:"redisip"`
		RedisPort string `json:"redisport"`
		Password  string `json:"password"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &redisinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	redisaddress := net.ParseIP(redisinfo.RedisIP)
	if redisaddress == nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "redis 地址错误"})
		return
	} else if !util.NetWorkStatus(redisaddress.String()) {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "redis 网络连接异常"})
		return
	}
	createtime := time.Now().Unix()
	redisid := util.GetStrMd5("ehoney")
	data, msg, msgcode := honeycluster.InsertRedisInfo(redisid, redisinfo.RedisIP, redisinfo.RedisPort, redisinfo.Password, createtime)
	comhttp.SendJSONResponse(w, comm.Response{Code: msgcode, Data: data, Message: msg})
	return
}

func GetRedisInfo(w http.ResponseWriter, r *http.Request) {
	datas := honeycluster.SelectRedisInfo()
	if len(datas) > 0 {
		datas[0]["password"] = "***"
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func DeleteRedisInfo(w http.ResponseWriter, r *http.Request) {
	var redisinfo struct {
		RedisId string `json:"redisid"`
	}
	body := comhttp.GetPostBody(w, r)

	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &redisinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	data, msg, msgCode := honeycluster.DeleteRedisInfo(redisinfo.RedisId)
	comhttp.SendJSONResponse(w, comm.Response{Code: msgCode, Data: data, Message: msg})
	return
}

func AddTraceHostHandler(w http.ResponseWriter, r *http.Request) {
	var traceinfo struct {
		TraceHost string `json:"tracehost"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &traceinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	createtime := time.Now().Unix()
	traceid := util.GetStrMd5("ehoney")
	data, msg, msgcode := honeycluster.InsertTraceInfo(traceid, traceinfo.TraceHost, createtime)
	comhttp.SendJSONResponse(w, comm.Response{Code: msgcode, Data: data, Message: msg})
	return
}

func GetTraceHostHandler(w http.ResponseWriter, r *http.Request) {
	datas := honeycluster.SelectTraceInfo()
	tracehost := ""
	if len(datas) > 0 {
		tracehost = util.Strval(datas[0]["tracehost"])
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: tracehost, Message: "成功"})
	return
}

// 接收欺骗防御的攻击日志
func InsertAttackLogHandler(w http.ResponseWriter, r *http.Request) {
	var attackLog struct {
		DstHost    string      `json:"DstHost"` // HoneyIP
		DstPort    int64       `json:"DstPort"` // HoneyPotPort
		SrcHost    string      `json:"SrcHost"` // serverIP
		SrcPort    int64       `json:"SrcPort"` // serverPort
		AttackHost string      `json:"AttackHost"`
		AttackPort int64       `json:"AttackPort"`
		LocalTime  string      `json:"LocalTime"` // AttackTime
		LogType    string      `json:"LogType"`
		NodeId     string      `json:"NodeId"`
		LogData    interface{} `json:"LogData"`
		SourceType int         `json:"SourceType"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})

		return
	}

	if err := json.Unmarshal([]byte(body), &attackLog); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	logdata, err := json.Marshal(attackLog.LogData)
	if err != nil {
		logs.Error("[InsertAttackLogHandler] logdata json.Marshal err:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertAttackLogHandler] log data marshal error: %v", err)})
		return
	}
	logdatas := string(logdata)
	logdatas = strings.ReplaceAll(logdatas, "\\r\\n", "\\\\r\\\\n")
	eventDetail := comhttp.InterfaceToString(attackLog.LogType, logdata)
	eventDetail = strings.ReplaceAll(eventDetail, "\\n", "")
	sqlCon, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")

	if err != nil {
		logs.Error("open mysql fail %v\n", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertAttackLogHandler] open mysql fail %v", err)})
		return
	}

	country := ""
	province := ""

	if util.IsLocalIP(attackLog.AttackHost){
		country = "局域网"
		province = "局域网"
	}else{
		result, err := util.GetLocationByIP(attackLog.AttackHost)
		if err != nil{
			logs.Error("get location error: %v\n", err)
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("get location error %v", err)})
			return
		}
		country = result.Country_long
		province = result.City
	}

	honeymaps := honeycluster.SelectHoneyPotByIp(attackLog.DstHost)
	if len(honeymaps) > 0 {
		srchost := honeymaps[0]["serverip"]
		serverid := honeymaps[0]["serverid"]
		honeyportport := honeymaps[0]["honeyport"]

		DbCon = sqlCon
		defer sqlCon.Close()
		DbCon.SetConnMaxLifetime(300 * time.Second)
		DbCon.SetMaxOpenConns(120)
		DbCon.SetMaxIdleConns(12)
		sqlstr := "select honeytypeid,serverid from honeypots where honeyip = ? and serverid = ?"
		var argsList []interface{}
		argsList = append(argsList, attackLog.DstHost)
		argsList = append(argsList, attackLog.NodeId)
		rows, err0 := DbCon.Query(sqlstr, argsList...)
		if err0 != nil {
			logs.Error("[InsertAttackLogHandler] %s\v", err0)
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertAttackLogHandler] scan failed, err:%v", err0)})
			return
		}

		columns, errr := rows.Columns()
		if errr != nil {
			logs.Error("[InsertAttackLogHandler] %s\v", errr)
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertAttackLogHandler] scan failed, err:%v", errr)})
			return
		}

		values := make([]sql.RawBytes, len(columns))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		honeyTypeIds := []string{}

		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				logs.Error("select mysql fail: ", err)
				comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertAttackLogHandler] insert mysql fail: %v", err)})
				return
			}

			var value string
			for _, col := range values {
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}
				honeyTypeIds = append(honeyTypeIds, value)
			}
		}
		db1, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
		if err != nil {

			logs.Error("open mysql fail %v\n", err)
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertAttackLogHandler] open mysql fail %v", err)})
			return
		}
		defer db1.Close()

		db1.SetConnMaxLifetime(300 * time.Second)
		db1.SetMaxOpenConns(120)
		db1.SetMaxIdleConns(12)

		honeyTypeId := ""
		if len(honeyTypeIds) > 0 {
			honeyTypeId = honeyTypeIds[0]
		}
		var attacktime int64
		if err == nil {
			attacktime = time.Now().Unix()
		} else {
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertFalcoLogHandler] insert mysql fail: %v", err)})
			return
		}
		honeypot := honeycluster.SelectHoneyInfoByIp(attackLog.DstHost)
		if len(honeypot) > 0 {
			honeypotid := honeypot[0]["honeypotid"]
			honeytypeid := honeypot[0]["honeytypeid"]
			if honeyTypeId == "" {
				honeyTypeId = honeytypeid.(string)
			}
			res, err3 := db1.Exec("INSERT INTO attacklog(srchost,srcport,serverid,honeypotid,honeypotport,attackip,attackport,attacktime,eventdetail,proxytype,sourcetype,logdata,honeytypeid,country,province) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
				srchost, attackLog.SrcPort, serverid, honeypotid, honeyportport, attackLog.AttackHost, attackLog.AttackPort, attacktime, eventDetail, attackLog.LogType, attackLog.SourceType, logdatas, honeyTypeId, country, province)
			if err3 != nil {
				logs.Error("insert mysql fail: ", err3)
				comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertAttackLogHandler] insert mysql fail: %v", err3)})
				return
			}

			_, err = res.RowsAffected()
			if err != nil {
				logs.Error("insert failed, err: ", err)
				comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: fmt.Sprintf("[InsertAttackLogHandler] insert failed, err:%v", err)})
				return
			}
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: comm.DBInsertSuccess})
			return
		} else {
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.DataUnmarshalError})
			return
		}
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "获取蜜罐异常"})
		return
	}
}

func InsertSSHKeyHandler(w http.ResponseWriter, r *http.Request) {
	var sshkey struct {
		SSHKey  string `json:"ssh_key"`
		AgentID string `json:"agentid"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &sshkey); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	datas, msg, msgcode := honeycluster.InsertSSHInfo(sshkey.SSHKey, sshkey.AgentID)
	comhttp.SendJSONResponse(w, comm.Response{Code: msgcode, Data: datas, Message: msg})
	return

}

// 创建容器pod
func AddPod(w http.ResponseWriter, r *http.Request) {

	var podinfo struct {
		Honeytypeid   string `json:"honeytypeid"`
		Honeyserverid string `json:"honeyserverid"`
		Honeysysid    string `json:"honeysysid"`
		Image         string `json:"image"`
		Name          string `json:"name"`
		ContainerPort int    `json:"containerport"`
	}
	honeypotServergAnentid := ""
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}

	if err := json.Unmarshal([]byte(body), &podinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + kubeconfig
	namespace := "default"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logs.Error("[AddPod] Error:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: err.Error(), Message: "失败"})
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		logs.Error("[AddPod] Error:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: err.Error(), Message: "失败"})
	}

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": podinfo.Name,
			},
			"spec": map[string]interface{}{
				"replicas": 1,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "demo",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "demo",
						},
					},

					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  podinfo.Name,
								"image": podinfo.Image,
								"ports": []map[string]interface{}{
									{
										"name":          podinfo.Name,
										"protocol":      "TCP",
										"containerPort": podinfo.ContainerPort,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := client.Resource(deploymentRes).Namespace(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logs.Error("[AddPod] Error:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: err.Error(), Message: err.Error()})
	} else {
		pod := k3s.GetPodinfo(podinfo.Name)
		createtime := time.Now().Unix()
		honeypotserver := honeycluster.SelectHoneypotServer()
		if len(honeypotserver) > 0 {
			honeypotServergAnentid = util.Strval(honeypotserver[0]["agentid"])
		}
		honeycluster.InsertPodinfo(podinfo.Name, podinfo.Honeytypeid, util.GetUUID(), pod.Status.PodIP, podinfo.ContainerPort, createtime, "admin", honeypotServergAnentid, podinfo.Honeysysid)
		comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: result.GetClusterName(), Message: "成功"})
		time.Sleep(3 * time.Second)
		go k3s.FreshPods()
	}
	return
}

func AddPodv2(w http.ResponseWriter, r *http.Request) {

	var podinfo struct {
		Honeyserverid string `json:"honeyserverid"`
		ImageID       int    `json:"imageid"`
		Name          string `json:"name"`
	}
	honeypotServergAnentid := ""
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}

	if err := json.Unmarshal([]byte(body), &podinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	imageinfo := honeycluster.SelectHoneyImageById(podinfo.ImageID)
	sysid := ""
	typeid := ""
	imageport := 0
	imageaddr := ""
	if len(imageinfo) > 0 && imageinfo[0]["imageos"] != nil && imageinfo[0]["imagetype"] != nil && imageinfo[0]["imageport"] != nil {
		sysid = util.Strval(imageinfo[0]["imageos"])
		typeid = util.Strval(imageinfo[0]["imagetype"])
		imageport, _ = strconv.Atoi(util.Strval(imageinfo[0]["imageport"]))
		imageaddr = util.Strval(imageinfo[0]["imageaddress"])
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "镜像属性信息异常"})
		return
	}

	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + kubeconfig
	namespace := "default"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logs.Error("[AddPod] Error:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: err.Error(), Message: "失败"})
		return
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		logs.Error("[AddPod] Error:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: err.Error(), Message: "失败"})
		return
	}

	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": podinfo.Name,
			},
			"spec": map[string]interface{}{
				"replicas": 1,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "demo",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "demo",
						},
					},

					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  podinfo.Name,
								"image": imageaddr,
								"ports": []map[string]interface{}{
									{
										"name":          podinfo.Name,
										"protocol":      "TCP",
										"containerPort": imageport,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := client.Resource(deploymentRes).Namespace(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logs.Error("[AddPod] Error:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: err.Error(), Message: err.Error()})
		return
	} else {
		pod := k3s.GetPodinfo(podinfo.Name)
		createtime := time.Now().Unix()
		honeypotserver := honeycluster.SelectHoneypotServer()
		if len(honeypotserver) > 0 {
			honeypotServergAnentid = util.Strval(honeypotserver[0]["agentid"])
		}
		honeycluster.InsertPodinfo(podinfo.Name, typeid, util.GetUUID(), pod.Status.PodIP, imageport, createtime, "admin", honeypotServergAnentid, sysid)
		time.Sleep(4 * time.Second)
		comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: result.GetClusterName(), Message: "成功"})
		return
	}
	return
}

// 删除容器pod
func DeletePod(w http.ResponseWriter, r *http.Request) {
	var podinfo struct {
		PodId string `json:"podid"`
		Name  string `json:"name"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}

	if err := json.Unmarshal([]byte(body), &podinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + kubeconfig
	namespace := "default"

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logs.Error("[DeletePod] Error:", err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		logs.Error("[DeletePod] Error:", err)
	}
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := client.Resource(deploymentRes).Namespace(namespace).Delete(context.TODO(), podinfo.Name, deleteOptions, ""); err != nil {
		logs.Error("[DeletePod] Error:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "下线失败:" + err.Error()})
	} else {
		offlinetime := time.Now().Unix()
		honeycluster.UpdatePod(podinfo.Name, podinfo.PodId, offlinetime)
		comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: nil, Message: "成功"})
	}
	return
}

func CheckPodStatus(w http.ResponseWriter, r *http.Request) {
	var podinfo struct {
		PodId string `json:"podid"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}

	if err := json.Unmarshal([]byte(body), &podinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	pod := honeycluster.SelectHoneyInfoById(podinfo.PodId)
	if len(pod) > 0 {
		podname := util.Strval(pod[0]["honeyname"])
		podinfo := k3s.GetPodinfo(podname)
		if podinfo.Name != "" {
			if len(podinfo.Status.ContainerStatuses) > 0 {
				if podinfo.Status.ContainerStatuses[0].State.Running != nil {
					comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: "成功"})
					return
				} else {
					err := fmt.Errorf("%s", podinfo.Status.ContainerStatuses[0].State.Waiting.Reason)
					comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: err.Error()})
					return
				}
			}
		}
	}

}

/**
容器状态校验
*/
func CheckPodNetStatus(w http.ResponseWriter, r *http.Request) {
	var podinfo struct {
		//PodId string `json:"podid"`
		Name string `json:"name"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}

	if err := json.Unmarshal([]byte(body), &podinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	errmsg := ""
	_, _, err, createstatus, errcode := k3s.IsPodRunningV3(podinfo.Name, apiV1.NamespaceDefault)
	errcodestr := fmt.Sprintf("%d", errcode)
	if err != nil {
		logs.Error(err.Error() + " " + errcodestr)
		errmsg = err.Error()

	}
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: createstatus, Message: errmsg})
	return
}

/**
容器状态校验
*/
func TestCheckPodStatus(w http.ResponseWriter, r *http.Request) {
	var podinfo struct {
		PodId string `json:"honeypotid"`
		//Name string `json:"name"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}

	if err := json.Unmarshal([]byte(body), &podinfo); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	pod := honeycluster.SelectHoneyInfoById(podinfo.PodId)
	if len(pod) == 0 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.PodNotFound})
		return
	}

	honeyname := util.Strval(pod[0]["honeyname"])
	podname := util.Strval(pod[0]["podname"])
	_, _, _, _, errcode := k3s.CheckPodStatus(podname, apiV1.NamespaceDefault)
	if errcode == 0 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: 0, Message: "容器服务配置异常"})
		return
	} else if errcode == 2 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: 0, Message: "容器" + honeyname + "正在创建中"})
		return
	} else if errcode == 3 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: 0, Message: "容器" + honeyname + "创建异常：" + " 状态 ErrImagePull 镜像拉取失败"})
		return
	} else if errcode == 4 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: 0, Message: "容器" + honeyname + "创建异常：" + " 状态 CrashLoopBackOff 容器创建失败"})
		return
	} else if errcode == 5 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: 0, Message: "容器" + honeyname + "正常运行中"})
		return
	}
}

/**
容器镜像列表定时拉取
*/
func GetPodImage(w http.ResponseWriter, r *http.Request) {
	var podimage comm.PodImageResult
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &podimage); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.GetPodImage(podimage.PodName, podimage.PageSize, podimage.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetPodImageList() error {
	harborinfo := honeycluster.SelectHarborInfo()
	if len(harborinfo) > 0 {
		harborip := ""
		harborport := ""
		harborurl := util.Strval(harborinfo[0]["harborurl"])
		harborprojectname := harborinfo[0]["projectname"]
		harborproject := util.Strval(harborprojectname)
		harborhost := util.GetHost(harborurl)
		if harborhost != "" {
			harborhostinfo := strings.Split(harborhost, ":")
			if len(harborhostinfo) == 1 {
				harborip = harborhostinfo[0]
				harborport = "80"
			} else {
				harborip = harborhostinfo[0]
				harborport = harborhostinfo[1]
			}
		} else {
			err := fmt.Errorf("%s", "harhost解析异常")
			return err
		}
		if util.NetConnectTest(harborip, harborport) {
			harborUname := util.Strval(harborinfo[0]["username"])
			harborPwd := util.Strval(harborinfo[0]["password"])
			hearders := make(map[string]string)
			authorizationstr := util.Base64Encode(harborUname + ":" + harborPwd)
			hearders["authorization"] = "Basic " + authorizationstr
			resp := comhttp.GetResponseByHttp("GET", harborurl+"/api/v2.0/projects/"+harborproject+"/repositories?page=1&page_size=20", "", hearders)
			if resp != nil && resp.StatusCode == 200 {
				body, harborerr := ioutil.ReadAll(resp.Body)
				if harborerr != nil {
					logs.Error("requesr harborurl err:", harborerr)
					return harborerr
				}
				respbody := string(body)
				dec := json.NewDecoder(strings.NewReader(respbody))
				var pods []comm.PodName
				err := dec.Decode(&pods)
				if err != nil {
					logs.Error("Decode pods info:", err)
					return err
				} else {
					for _, pod := range pods {
						imagename := strings.Replace(pod.Name, harborproject, "", -1)
						imageresp := comhttp.GetResponseByHttp("GET", harborurl+"/api/v2.0/projects/"+harborproject+"/repositories/"+imagename+"/artifacts?with_tag=true&with_scan_overview=true&with_label=true&page_size=20&page=1", "", hearders)
						imagebody, _ := ioutil.ReadAll(imageresp.Body)
						imagerespbody := string(imagebody)
						imagedec := json.NewDecoder(strings.NewReader(imagerespbody))
						var podimagetags []comm.PodImageInfo
						err := imagedec.Decode(&podimagetags)
						if err != nil {
							logs.Error("imagedec Decode Error:", err)
							return err
						} else {
							for _, podimg := range podimagetags {
								for _, podtag := range podimg.Tags {
									podurl := harborhost + "/" + pod.Name + ":" + podtag.Name
									honeycluster.InsertPodImage(pod.Name, podurl)
								}
							}
						}
					}
				}
			}
			time.Sleep(5 * time.Minute)
			go GetPodImageList()
		} else {
			err := fmt.Errorf("%s", "harbor 网络异常")
			return err
		}
	} else {
		err := fmt.Errorf("%s", "harbor记录为空")
		return err
	}
	return nil
}

/*
上传诱饵
*/
func CreateBaitHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	sysid := query.Get("sysid")
	baitname := query.Get("baitname")
	filemd5, baitinfo, _ := FileUpload(r, baitname)
	createtime := time.Now().Unix()
	baitid := util.GetUUID()
	creator := "admin"
	honeycluster.InsertBait(baitname, createtime, creator, baitid, sysid, baitinfo, filemd5)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: nil, Message: "上传成功"})
}

/**
删除诱饵
*/
func DeleteBaitHandler(w http.ResponseWriter, r *http.Request) {
	var baits comm.Baits
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &baits); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	baitid := baits.BaitId
	baitinfo := honeycluster.SelectBaitsById(baitid)
	baitname := baitinfo[0]["baitname"]
	baitpath := "upload/" + util.Strval(baitname)
	result, msg, code := honeycluster.DeleteBaitById(baitid)
	err := os.RemoveAll(baitpath)
	if err != nil {
		logs.Error("delete baitpath err:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: result, Message: msg})
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: result, Message: msg})
	}
	return
}

/**
创建密签
*/
func CreateSignHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	signname := query.Get("signname")
	signinfo, err, isfalse := SignFileUpload(r, signname)
	if err != nil {
		logs.Error("CreateSignHandler error:", err)
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "上传失败"})
		return
	}
	if isfalse {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: "上传失败", Message: "上传类型只支持doc,docx,ppt,pptx,xlsx,xls,pdf"})
		return
	}
	createtime := time.Now().Unix()
	signid := util.GetUUID()
	creator := "admin"
	honeycluster.InsertSign("file", signname, createtime, creator, signid, "", signinfo)
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: "上传成功"})
	return
}

/**
删除密签
*/
func DeleteSignHandler(w http.ResponseWriter, r *http.Request) {
	var signs comm.SignJson
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &signs); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	signid := signs.SignId
	signinfo := honeycluster.SelectSignsById(signid)
	if len(signinfo) > 0 {
		signname := signinfo[0]["signname"]
		signpath := "upload/honeytoken/" + util.Strval(signname)
		result, msg, code := honeycluster.DeleteSignById(signid)
		err := os.RemoveAll(signpath)
		if err != nil {
			logs.Error("delete signpath err:", err)
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: result, Message: msg})
		} else {
			comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: result, Message: msg})
		}
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "no this sign"})
	}

	return
}

/**
密签列表

*/
func GetSigns(w http.ResponseWriter, r *http.Request) {
	var signs comm.SignJson
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &signs); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectSigns(signs.SignName, signs.SignId, signs.Type, signs.Creator, signs.StartTime, signs.EndTime, signs.PageSize, signs.PageNum)

	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

/**
密签列表

*/
func GetSignType(w http.ResponseWriter, r *http.Request) {
	results := policyCenter.SelectAllSignType()
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: results, Message: comm.DataSelectSuccess})
	return
}

/**
诱饵管理
*/
func GetBaits(w http.ResponseWriter, r *http.Request) {
	var baits comm.Baits
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &baits); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectBaits(baits.BaitId, baits.BaitType, baits.BaitSysType, baits.Creator, baits.StartTime, baits.EndTime, baits.PageSize, baits.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

/**
根据诱饵类型诱饵管理
*/
func GetBaitsByType(w http.ResponseWriter, r *http.Request) {
	var baits comm.Baits
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &baits); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectBaitsByType(baits.BaitId, baits.BaitType, baits.BaitSysType, baits.Creator, baits.PageSize, baits.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

/**
根据密签类型密签管理
*/
func GetSignsByType(w http.ResponseWriter, r *http.Request) {
	var signs comm.SignJson
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &signs); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectSignsByType(signs.SignId, signs.Type, signs.Creator)
	comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: datas, Message: "成功"})
	return
}

/**
获取蜜罐集群列表
*/
func GetHoneyClusters(w http.ResponseWriter, r *http.Request) {
	var honeyClusters comm.HoneyCluster
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeyClusters); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectHoneyClusters(honeyClusters.ClusterIp, honeyClusters.ClusterName, honeyClusters.ClusterStats, honeyClusters.PageSize, honeyClusters.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return

}

func GetHoneyList(w http.ResponseWriter, r *http.Request) {
	var honeypots comm.HoneyPots
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeypots); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectHoneyList(honeypots.ServerId, honeypots.HoneyTypeId, honeypots.HoneyIp, honeypots.HoneyName, honeypots.StartTime, honeypots.EndTime, honeypots.Status, honeypots.Creator, honeypots.PageSize, honeypots.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetHoneyListForTrans(w http.ResponseWriter, r *http.Request) {
	var honeypots comm.HoneyPots
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeypots); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectHoneyListForTrans(honeypots.HoneyTypeId)
	if len(datas) > 0 {
		for i, value := range datas {
			datas[i]["honeyname"] = util.Strval(value["honeyname"]) + "(" + util.Strval(value["honeyip"]) + ":" + util.Strval(value["honeyport"]) + ")"
			delete(datas[i], "honeyip")
			delete(datas[i], "honeyport")
		}
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetHoneyForwardListForTrans(w http.ResponseWriter, r *http.Request) {
	var honeypots comm.HoneyPots
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeypots); err != nil {
		logs.Error(fmt.Sprintf("[GetHoneyListForTrans] Failed to read the request body: %v", err))
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectHoneyForwardListForTrans()
	if len(datas) > 0 {
		for i, value := range datas {
			datas[i]["honeyname"] = util.Strval(value["honeyname"]) + "(" + util.Strval(value["honeyip"]) + ":" + util.Strval(value["honeyport"]) + ")"
			delete(datas[i], "honeyip")
			delete(datas[i], "honeyport")
		}
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func GetHoneyTransPortsListForTrans(w http.ResponseWriter, r *http.Request) {
	var honeypots comm.HoneyPots
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeypots); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectHoneyTransPortsListForTrans(honeypots.HoneypotId)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

/**
蜜罐管理
*/
func GetHoneyInfos(w http.ResponseWriter, r *http.Request) {
	var honeypots comm.HoneyPots
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeypots); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectHoneyInfos(honeypots.ServerId, honeypots.SysId, honeypots.HoneyTypeId, honeypots.HoneyIp, honeypots.HoneyName, honeypots.StartTime, honeypots.EndTime, honeypots.Status, honeypots.Creator, honeypots.PageSize, honeypots.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return

}

/**
蜜罐诱饵管理
*/
func CreateHoneyBaitPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var honeybait comm.HoneyBaits
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeybait); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	if strings.Contains(honeybait.Address, "./") || strings.Contains(honeybait.Address, ".\\") {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.DirUseError})
		return
	}
	honeypotid := honeybait.HoneypotId
	honeypotinfo := honeycluster.SelectHoneyInfoById(honeypotid)
	honeybaitdespath := honeybait.Address
	honeybaitdespath = strings.ReplaceAll(honeybaitdespath, "'", "''")
	baitdespathcheckcmdstr := []string{"/bin/sh", "-c", "find '" + honeybaitdespath + "'"}
	checkerr := k3s.ExecPodCmd(honeypotid, baitdespathcheckcmdstr)
	if checkerr != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: honeybaitdespath + " 目录不存在"})
		return
	}
	baitid := honeybait.BaitId
	baitfilepath := ""
	baitname := ""
	baiturl := ""
	sysproto := ""
	createbaitstatus := 0
	baitdespath := ""
	if r.TLS == nil {
		sysproto = "http://"
	} else {
		sysproto = "https://"
	}
	if honeybait.BaitType == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: nil, Message: "BaitType 为空"})
	}
	if honeybait.BaitId == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: nil, Message: "BaitId 为空"})
	}
	if honeybait.BaitType == "file" {
		datas := honeycluster.SelectBaitsById(baitid)
		baitname = datas[0]["baitname"].(string)
		baitfilepath = "upload/" + datas[0]["baitname"].(string)
		baitdespath = honeybait.Address + "/" + datas[0]["baitinfo"].(string)
		util.ModifyFile(baitfilepath, honeybait.Address, datas[0]["baitname"].(string), datas[0]["baitinfo"].(string))

		util.FileTarZip(baitfilepath, baitfilepath+".tar.gz")
		apphost := beego.AppConfig.String("apphost")
		appport := beego.AppConfig.String("appport")
		apphost = apphost + ":" + appport
		baiturl = sysproto + apphost + "/" + baitfilepath + ".tar.gz"

	}
	kubeconfig := beego.AppConfig.String("kubeconfig")
	kubeconfig = util.GetCurrentPathString() + "/" + kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logs.Error("clientcmd.BuildConfigFromFlags ERR:", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logs.Error("kubernetes.NewForConfig ERR:", err)
	}
	var cmdlist [][]string
	cmd1 := []string{"wget", "-P", "/tmp/", baiturl}
	cmd2 := []string{"mkdir", "/tmp/" + baitname}
	cmd3 := []string{"tar", "-zxvf", "/tmp/" + baitname + ".tar.gz", "-C", "/tmp/" + baitname}
	cmd4 := []string{"/bin/sh", "/tmp/" + baitname + "/install.sh"}
	cmd5 := []string{"rm", "-rf", "/tmp/" + baitname + ".tar.gz"}
	cmd6 := []string{"rm", "-rf", "/tmp/" + baitname}
	cmdlist = append(cmdlist, cmd1)
	cmdlist = append(cmdlist, cmd2)
	cmdlist = append(cmdlist, cmd3)
	cmdlist = append(cmdlist, cmd4)
	cmdlist = append(cmdlist, cmd5)
	cmdlist = append(cmdlist, cmd6)

	for index, value := range cmdlist {
		fmt.Println(index)
		req := clientset.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(util.Strval(honeypotinfo[0]["podname"])).
			Namespace("default").
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Command: value,
				Stdin:   true,
				Stdout:  true,
				Stderr:  true,
				TTY:     false,
			}, scheme.ParameterCodec)
		exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
		if err != nil {
			logs.Error("CreateHoneyBaitPolicyHandler cmd err:", err)
		}

		screen := struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}

		if err = exec.Stream(remotecommand.StreamOptions{
			Stdin:  screen,
			Stdout: screen,
			Stderr: screen,
			Tty:    false,
		}); err != nil {
			logs.Error("[CreateHoneyBaitPolicyHandler] exec Error:", err)
		}

	}
	baitdespath = strings.ReplaceAll(baitdespath, "'", "''")
	cmdstr := []string{"/bin/sh", "-c", "find '" + baitdespath + "'"}
	err = k3s.ExecPodCmd(honeypotid, cmdstr)
	if err != nil {
		logs.Error("[CreateHoneyBaitPolicyHandler] err:", err)
		createbaitstatus = 2
	} else {
		createbaitstatus = 1
	}

	createtime := time.Now().Unix()
	taskid := util.GetUUID()
	datas, msg, _ := honeycluster.InsertHoneyBait(createbaitstatus, taskid, honeybait.BaitId, honeybait.Address, honeybait.HoneypotId, createtime, "admin")
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: msg})
	return

}

/**
蜜罐诱饵管理-撤回
*/
func DeleteHoneyBaitPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var honeybait comm.HoneyBaits
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeybait); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	baitfilepath := ""
	taskid := honeybait.Taskid
	deletebaitstatus := 0
	honeypotid := ""
	datas := honeycluster.SelectHoneyBaitsById(taskid)
	if len(datas) != 0 {
		honeypotid = datas[0]["honeypotid"].(string)
		podname := datas[0]["podname"].(string)
		baitfilepath = datas[0]["baitinfo"].(string)
		baitfilename := datas[0]["baitfilename"].(string)
		baitfilepath = baitfilepath + "/" + baitfilename
		baitname := datas[0]["baitname"].(string) + "-uninstall"
		util.ModifyUninstallFile(baitfilepath, baitname)
		baituploadpath := "upload/" + baitname
		util.FileTarZip(baituploadpath, baituploadpath+".tar.gz")
		sysproto := ""
		if r.TLS == nil {
			sysproto = "http://"
		} else {
			sysproto = "https://"
		}
		apphost := beego.AppConfig.String("apphost")
		appport := beego.AppConfig.String("appport")
		apphost = apphost + ":" + appport
		baiturl := sysproto + apphost + "/" + baituploadpath + ".tar.gz"

		kubeconfig := beego.AppConfig.String("kubeconfig")
		kubeconfig = util.GetCurrentPathString() + "/" + kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logs.Error("clientcmd.BuildConfigFromFlags Error:", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			logs.Error("kubernetes.NewForConfig Error:", err)
		}
		uninstallpath := "/tmp/uninstall" + taskid

		var cmdlist [][]string
		cmd1 := []string{"mkdir", uninstallpath}
		cmd2 := []string{"wget", "-O", uninstallpath + "/" + "uninstall.tar.gz", baiturl}
		cmd3 := []string{"tar", "-zxvf", uninstallpath + "/" + "uninstall.tar.gz", "-C", uninstallpath}
		cmd4 := []string{"/bin/sh", uninstallpath + "/install.sh"}
		cmd5 := []string{"rm", "-rf", uninstallpath}
		cmdlist = append(cmdlist, cmd1)
		cmdlist = append(cmdlist, cmd2)
		cmdlist = append(cmdlist, cmd3)
		cmdlist = append(cmdlist, cmd4)
		cmdlist = append(cmdlist, cmd5)

		for index, value := range cmdlist {
			fmt.Println(index)
			req := clientset.CoreV1().RESTClient().Post().
				Resource("pods").
				Name(podname).
				Namespace("default").
				SubResource("exec").
				VersionedParams(&corev1.PodExecOptions{
					Command: value,
					Stdin:   true,
					Stdout:  true,
					Stderr:  true,
					TTY:     false,
				}, scheme.ParameterCodec)
			exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
			if err != nil {
				logs.Error("cmd err:", err)
			}
			screen := struct {
				io.Reader
				io.Writer
			}{os.Stdin, os.Stdout}

			if err = exec.Stream(remotecommand.StreamOptions{
				Stdin:  screen,
				Stdout: screen,
				Stderr: screen,
				Tty:    false,
			}); err != nil {
				logs.Error("[DeleteHoneyBaitPolicyHandler] exec Error:", err)
			}

		}
		baitfilepath = strings.ReplaceAll(baitfilepath, "'", "''")
		cmdstr := []string{"/bin/sh", "-c", "find '" + baitfilepath + "'"}
		err = k3s.ExecPodCmd(honeypotid, cmdstr)
		if err != nil {
			deletebaitstatus = 5
		} else {
			deletebaitstatus = 6
		}

		offtime := time.Now().Unix()
		honeycluster.UpdateHoneyBait(deletebaitstatus, honeybait.Taskid, offtime)
	}
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: nil, Message: "成功"})
	return

}

/**
蜜罐密签管理
*/
func CreateHoneySignPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var honeysign comm.HoneySigns
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeysign); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	if strings.Contains(honeysign.Address, "./") || strings.Contains(honeysign.Address, ".\\") {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.DirUseError})
		return
	}

	honeypotid := honeysign.HoneypotId
	honeypotinfo := honeycluster.SelectHoneyInfoById(honeypotid)
	honeysigndespath := honeysign.Address
	honeysigndespath = strings.ReplaceAll(honeysigndespath, "'", "''")
	signdespathcheckcmdstr := []string{"/bin/sh", "-c", "find '" + honeysigndespath + "'"}
	checkerr := k3s.ExecPodCmd(honeypotid, signdespathcheckcmdstr)
	if checkerr != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: honeysigndespath + " 目录不存在"})
		return
	}
	taskid := util.GetUUID()
	signid := honeysign.SignId
	signfilepath := ""
	signname := ""
	signfilename := ""
	signurl := ""
	sysproto := ""
	createsignstatus := 0
	signdespath := ""
	traceurl := ""
	if r.TLS == nil {
		sysproto = "http://"
	} else {
		sysproto = "https://"
	}
	tracecode := util.GetUUID()
	traceinfo := honeycluster.GetTraceHost()
	if len(traceinfo) == 0 {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "请到系统设置配置密签追踪url"})
		return
	}
	tracehost := util.Strval(traceinfo[0]["tracehost"])
	if strings.Contains(tracehost, "http") {
		traceurl = fmt.Sprintf("%s/api/msgreceive?tracecode=%s", tracehost, tracecode)
	} else {
		traceurl = fmt.Sprintf("http://%s/api/msgreceive?tracecode=%s", tracehost, tracecode)
	}
	if honeysign.SignType != "" {
		if honeysign.SignType == "file" {
			datas := honeycluster.SelectSignsById(signid)
			if len(datas) == 0 {
				comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "密签不存在!"})
				return
			}
			signname = datas[0]["signname"].(string)
			signfilename = datas[0]["signinfo"].(string)
			signsourcepath := "upload/honeytoken/" + signname
			signfilepath = "upload/honeytoken/" + taskid + "/" + signname
			err := util.MakeDir(signfilepath)
			if err != nil {
				logs.Error("[CreateHoneySignPolicyHandler] signfilepath make err:", err)
				comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "密签目录创建失败"})
				return
			}
			var Suffix = regexp.MustCompile(`.(ppt|pptx|doc|docx|pdf|xls|xlsx)$`)
			if len(Suffix.FindString(signfilename)) > 0 {
				err := honeytoken.DoFileSignTrace(signfilename, signfilename, signsourcepath, signfilepath, traceurl)
				if err != nil {
					error := util.CopyDir(signsourcepath, signfilepath)
					if error != nil {
						logs.Error("[CreateHoneySignPolicyHandler] util.CopyDir err:", error)
					}
				}
			}
			signdespath = honeysign.Address + "/" + datas[0]["signinfo"].(string)
			util.ModifyFile(signfilepath, honeysign.Address, datas[0]["signname"].(string), datas[0]["signinfo"].(string))
			util.FileTarZip(signfilepath, signfilepath+".tar.gz")
			apphost := beego.AppConfig.String("apphost")
			appport := beego.AppConfig.String("appport")
			apphost = apphost + ":" + appport
			signurl = sysproto + apphost + "/" + signfilepath + ".tar.gz"
			log.Println("honeysignurl:", signurl)
		}

		kubeconfig := beego.AppConfig.String("kubeconfig")
		kubeconfig = util.GetCurrentPathString() + "/" + kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logs.Error("[CreateHoneySignPolicyHandler]  clientcmd.BuildConfigFromFlags ERR:", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			logs.Error("[CreateHoneySignPolicyHandler]  kubernetes.NewForConfig ERR:", err)
		}

		var cmdlist [][]string
		cmd1 := []string{"wget", "-P", "/tmp/", signurl}
		cmd2 := []string{"mkdir", "/tmp/" + signname}
		cmd3 := []string{"tar", "-zxvf", "/tmp/" + signname + ".tar.gz", "-C", "/tmp/" + signname}
		cmd4 := []string{"/bin/sh", "/tmp/" + signname + "/install.sh"}
		cmd5 := []string{"rm", "-rf", "/tmp/" + signname + ".tar.gz"}
		cmd6 := []string{"rm", "-rf", "/tmp/" + signname}
		cmdlist = append(cmdlist, cmd1)
		cmdlist = append(cmdlist, cmd2)
		cmdlist = append(cmdlist, cmd3)
		cmdlist = append(cmdlist, cmd4)
		cmdlist = append(cmdlist, cmd5)
		cmdlist = append(cmdlist, cmd6)

		for index, value := range cmdlist {
			fmt.Println(index)
			req := clientset.CoreV1().RESTClient().Post().
				Resource("pods").
				Name(util.Strval(honeypotinfo[0]["podname"])).
				Namespace("default").
				SubResource("exec").
				VersionedParams(&corev1.PodExecOptions{
					Command: value,
					Stdin:   true,
					Stdout:  true,
					Stderr:  true,
					TTY:     false,
				}, scheme.ParameterCodec)
			exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
			if err != nil {
				logs.Error("cmd err:", err)
			}

			screen := struct {
				io.Reader
				io.Writer
			}{os.Stdin, os.Stdout}

			if err = exec.Stream(remotecommand.StreamOptions{
				Stdin:  screen,
				Stdout: screen,
				Stderr: screen,
				Tty:    false,
			}); err != nil {
				logs.Error("[CreateHoneySignPolicyHandler] exec Error:", err)
			}

		}
		signdespath = strings.ReplaceAll(signdespath, "'", "''")
		cmdstr := []string{"/bin/sh", "-c", "find '" + signdespath + "'"}
		err = k3s.ExecPodCmd(honeypotid, cmdstr)
		if err != nil {
			logs.Error("[CreateHoneySignPolicyHandler] err:", err)
			createsignstatus = 2
		} else {
			createsignstatus = 1
		}
		createtime := time.Now().Unix()

		datas, msg, _ := honeycluster.InsertHonetSign(createsignstatus, taskid, honeysign.SignId, signdespath, honeysign.HoneypotId, createtime, "admin", tracecode)
		comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: msg})
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: "SignType is needed!"})
	}
	return
}

/**
蜜罐密签管理-下载
*/
func DownloadHoneySignPolicyHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	taskid := query.Get("taskid")

	if taskid == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	datas := honeycluster.SelectHoneySignById(taskid)
	if len(datas) > 0 {
		signname := datas[0]["signname"].(string)
		signfilename := datas[0]["signfilename"].(string)
		filePath := util.GetCurrentPathString() + "/upload/honeytoken/" + taskid + "/" + signname + "/" + signfilename
		file, err := os.Open(filePath)
		if err != nil {
			logs.Error("sign download fail: %v , %s", err, filePath)
			return
		}
		defer file.Close()
		fileHeader := make([]byte, 1024)
		file.Read(fileHeader)
		fileStat, _ := file.Stat()
		w.Header().Set("Content-Disposition", "attachment; filename="+signfilename)
		w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
		w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
		file.Seek(0, 0)
		io.Copy(w, file)
		return
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: datas, Message: "没有该诱饵下发任务"})
		return
	}
}

/**
蜜罐密签管理-撤回
*/
func DeleteHoneySignPolicyHandler(w http.ResponseWriter, r *http.Request) {

	var honeysign comm.HoneySigns
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeysign); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	signfilepath := ""
	deletesignstatus := 0
	honeypotid := ""
	taskid := honeysign.Taskid
	datas := honeycluster.SelectHoneySignById(taskid)
	if len(datas) > 0 {
		honeypotid = datas[0]["honeypotid"].(string)
		podname := datas[0]["podname"].(string)
		signfilepath = datas[0]["signinfo"].(string)
		signfilename := datas[0]["signfilename"].(string)
		signfilepath = signfilepath + "/" + signfilename
		signname := datas[0]["signname"].(string) + "-uninstall"
		util.ModifyUninstallFile(signfilepath, signname)
		signuploadpath := "upload/" + signname
		util.FileTarZip(signuploadpath, signuploadpath+".tar.gz")
		sysproto := ""
		if r.TLS == nil {
			sysproto = "http://"
		} else {
			sysproto = "https://"
		}
		apphost := beego.AppConfig.String("apphost")
		appport := beego.AppConfig.String("appport")
		apphost = apphost + ":" + appport
		signurl := sysproto + apphost + "/" + signuploadpath + ".tar.gz"

		kubeconfig := beego.AppConfig.String("kubeconfig")
		kubeconfig = util.GetCurrentPathString() + "/" + kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logs.Error("clientcmd.BuildConfigFromFlags Error:", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			logs.Error("kubernetes.NewForConfig Error:", err)
		}
		uninstallpath := "/tmp/uninstall" + taskid
		var cmdlist [][]string
		cmd1 := []string{"mkdir", uninstallpath}
		cmd2 := []string{"wget", "-O", uninstallpath + "/" + "uninstall.tar.gz", signurl}
		cmd3 := []string{"tar", "-zxvf", uninstallpath + "/" + "uninstall.tar.gz", "-C", uninstallpath}
		cmd4 := []string{"/bin/sh", uninstallpath + "/install.sh"}
		cmd5 := []string{"rm", "-rf", uninstallpath}
		cmdlist = append(cmdlist, cmd1)
		cmdlist = append(cmdlist, cmd2)
		cmdlist = append(cmdlist, cmd3)
		cmdlist = append(cmdlist, cmd4)
		cmdlist = append(cmdlist, cmd5)

		for index, value := range cmdlist {
			fmt.Println(index)
			req := clientset.CoreV1().RESTClient().Post().
				Resource("pods").
				Name(podname).
				Namespace("default").
				SubResource("exec").
				VersionedParams(&corev1.PodExecOptions{
					Command: value,
					Stdin:   true,
					Stdout:  true,
					Stderr:  true,
					TTY:     false,
				}, scheme.ParameterCodec)
			exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
			if err != nil {
				logs.Error("DeleteHoneySignPolicyHandler cmd err:", err)
			}

			screen := struct {
				io.Reader
				io.Writer
			}{os.Stdin, os.Stdout}

			if err = exec.Stream(remotecommand.StreamOptions{
				Stdin:  screen,
				Stdout: screen,
				Stderr: screen,
				Tty:    false,
			}); err != nil {
				logs.Error("[DeleteHoneySignPolicyHandler] exec Error:", err)
			}

		}
		signfilepath = strings.ReplaceAll(signfilepath, "'", "''")
		cmdstr := []string{"/bin/sh", "-c", "find '" + signfilepath + "'"}
		err = k3s.ExecPodCmd(honeypotid, cmdstr)
		if err != nil {
			deletesignstatus = 5
		} else {
			deletesignstatus = 6
		}

		offtime := time.Now().Unix()
		honeycluster.UpdateHoneySign(deletesignstatus, honeysign.Taskid, offtime)
		comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
		return
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: datas, Message: "不存在该任务"})
		return
	}

}

/**
蜜罐密签列表管理
*/
func GetHoneySigns(w http.ResponseWriter, r *http.Request) {
	var honeysigns comm.HoneySigns
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeysigns); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectHoneySigns(honeysigns.HoneypotId, honeysigns.SignType, "admin", honeysigns.Status, honeysigns.StartTime, honeysigns.EndTime, honeysigns.PageSize, honeysigns.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func HoneySignMsgHandler(w http.ResponseWriter, r *http.Request) {
	var fileList struct {
		Tracecode string `json:"tracecode"`
		Map       string `json:"map"`
		PageNum   int    `json:"pageNum"`
		PageSize  int    `json:"pageSize"`
		StartTime string `json:"starttime"`
		EndTime   string `json:"endtime"`
	}
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &fileList); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.HoneySelectSignMsg(fileList.Tracecode, fileList.Map, fileList.StartTime, fileList.EndTime, fileList.PageSize, fileList.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return

}

/**
蜜罐诱饵管理
*/
func GetHoneyBaits(w http.ResponseWriter, r *http.Request) {
	var honeybaits comm.HoneyBaits
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &honeybaits); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	datas := honeycluster.SelectHoneyBaits(honeybaits.HoneypotId, honeybaits.BaitType, "admin", honeybaits.Status, honeybaits.StartTime, honeybaits.EndTime, honeybaits.PageSize, honeybaits.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func DownloadHoneyBaits(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	taskid := query.Get("taskid")

	if taskid == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	datas := honeycluster.SelectHoneyBaitsById(taskid)
	if len(datas) > 0 {
		baitname := datas[0]["baitname"].(string)
		baitfilename := datas[0]["baitfilename"].(string)
		filePath := util.GetCurrentPathString() + "/upload/" + baitname + "/" + baitfilename
		file, err := os.Open(filePath)
		if err != nil {
			logs.Error("sign download fail: %v , %s", err, filePath)
			return
		}
		defer file.Close()
		fileHeader := make([]byte, 1024)
		file.Read(fileHeader)
		fileStat, _ := file.Stat()
		w.Header().Set("Content-Disposition", "attachment; filename="+baitfilename)
		w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
		w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
		file.Seek(0, 0)
		io.Copy(w, file)
		return
	} else {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: datas, Message: "没有该诱饵下发任务"})
		return
	}
}

func InsertClamavResultHandler(w http.ResponseWriter, r *http.Request) {
	var clamavdata comm.ClamavData
	body := comhttp.GetPostBody(w, r)
	body = r.PostForm.Encode()
	body, _ = url.QueryUnescape(body)
	if util.HasSuffix(body, "=") {
		body = strings.TrimRight(body, "=")
	}
	honeypotip := ""
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &clamavdata); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}
	createtime := time.Now().Unix()

	if r.RemoteAddr != "" {
		honeypotip = util.GetIP(r.RemoteAddr)
	}
	datas, msg, code := clamavcenter.InsertClamavData(clamavdata.Filename, clamavdata.Virname, createtime, honeypotip)
	comhttp.SendJSONResponse(w, comm.Response{Code: code, Data: datas, Message: msg})
	return
}

func GetClamavResultHandler(w http.ResponseWriter, r *http.Request) {
	var clamavdata comm.ClamavData
	body := comhttp.GetPostBody(w, r)
	if body == "" {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyNullMsg})
		return
	}
	if err := json.Unmarshal([]byte(body), &clamavdata); err != nil {
		comhttp.SendJSONResponse(w, comm.Response{Code: comm.ErrorCode, Data: nil, Message: comm.BodyUnmarshalEorrMsg})
		return
	}

	datas := clamavcenter.SelectClamavData(clamavdata.Filename, clamavdata.Virname, clamavdata.HoneyPotip, clamavdata.CreateStartTime, clamavdata.CreateEndTime, clamavdata.PageSize, clamavdata.PageNum)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: datas, Message: "成功"})
	return
}

func testHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("xxxx:", *Addr)
	comhttp.SendJSONResponse(w, comm.Response{Code: 0, Data: nil, Message: "心跳正常"})
	return
}

//监控数据库中哪些上线诱饵策略未收到agent的返回信息，一分钟下发一次，三分钟内下发三次左右，三次之后均无响应，则置状态为-2
func MonitorForAgentByCreateBaitPolicy() {
	for {
		list, _ := policyCenter.SelectCreateBaitPolicyNotFinish()
		fmt.Println(list)
		if len(list) > 0 {
			for i := 0; i < len(list); i++ {
				fmt.Println(fmt.Sprintf(" taskid status is %v", list[i].Status))
				if list[i].Status == 0 {
					localTimes := time.Now().Unix()
					lastTimes, _ := strconv.ParseInt(list[i].CreateTime, 10, 64)
					sub := (localTimes - lastTimes) / 60
					if sub >= 0 && sub <= 3 {
						var baitPolicyJson comm.BaitPolicyJson
						baitPolicyJson.AgentId = list[i].AgentId
						baitPolicyJson.TaskId = list[i].TaskId
						baitPolicyJson.Data = list[i].Data
						baitPolicyJson.Md5 = list[i].Md5
						baitPolicyJson.Status = comm.CreateStatus
						baitPolicyJson.Type = list[i].Type

						baitPolicy, err := json.Marshal(baitPolicyJson)
						if err != nil {
							continue
						} else {
							redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
							fmt.Println(string(baitPolicy))
						}

					} else {
						policyCenter.UpdateBaitPolicyStatusByTaskid(comm.AgentOutTimeStatus, list[i].TaskId)
						fmt.Println(fmt.Sprintf("taskid %v status -2", list[i].TaskId))
					}
				}

			}
		}
		time.Sleep(time.Minute)
	}
}

//监控数据库中哪些下线诱饵策略未收到agent的返回信息，一分钟下发一次，三分钟内下发三次左右，三次之后均无响应，则置状态为-2
func MonitorForAgentByOfflineBaitPolicy() {
	go func() {
		for {
			list, _ := policyCenter.SelectOfflineBaitPolicyNotFinish()
			fmt.Println(list)
			if len(list) > 0 {
				for i := 0; i < len(list); i++ {
					fmt.Println(fmt.Sprintf(" taskid status is %v", list[i].Status))
					if list[i].Status == 1 {
						localTimes := time.Now().Unix()
						lastTimes, _ := strconv.ParseInt(list[i].OfflineTime, 10, 64)
						sub := (localTimes - lastTimes) / 60
						if sub >= 0 && sub <= 3 {
							var baitPolicyJson comm.BaitPolicyJson
							baitPolicyJson.AgentId = list[i].AgentId
							baitPolicyJson.TaskId = list[i].TaskId
							baitPolicyJson.Data = list[i].Data
							baitPolicyJson.Md5 = list[i].Md5
							baitPolicyJson.Status = comm.OfflineStatus
							baitPolicyJson.Type = list[i].Type

							baitPolicy, err := json.Marshal(baitPolicyJson)
							if err != nil {
								continue
							} else {
								redisCenter.RedisSubProducerBaitPolicy(string(baitPolicy))
								fmt.Println(string(baitPolicy))
							}

						} else {
							policyCenter.UpdateBaitPolicyStatusByTaskid(comm.AgentOutTimeStatus, list[i].TaskId)
							fmt.Println(fmt.Sprintf("taskid %v status -2", list[i].TaskId))
						}
					}

				}
			}
			time.Sleep(time.Minute)
		}
	}()

}

//监控数据库中哪些上线透明转发策略未收到agent的返回信息，一分钟下发一次，三分钟内下发三次左右，三次之后均无响应，则置状态为-2
func MonitorForAgentByCreateTransPolicy() {
	for {
		list, _ := policyCenter.SelectCreateTransPolicyNotFinish()
		fmt.Println(list)
		if len(list) > 0 {
			for i := 0; i < len(list); i++ {
				fmt.Println(fmt.Sprintf(" taskid status is %v", list[i].Status))
				if list[i].Status == 0 {
					localTimes := time.Now().Unix()
					lastTimes, _ := strconv.ParseInt(list[i].CreateTime, 10, 64)
					sub := (localTimes - lastTimes) / 60
					if sub >= 0 && sub <= 3 {
						var transPolicyJson comm.TransparentTransponderPolicyJson
						transPolicyJson.AgentId = list[i].AgentId
						transPolicyJson.TaskId = list[i].TaskId
						transPolicyJson.ListenPort = list[i].ForwardPort
						transPolicyJson.ServerType = list[i].HoneyPotType
						transPolicyJson.HoneyIP = list[i].HoneyIP
						transPolicyJson.HoneyPort = list[i].HoneyPotPort
						transPolicyJson.Status = comm.CreateStatus
						transPolicyJson.Type = comm.AgentTypeEdge
						transPolicyJson.Path = ""

						transPolicy, err := json.Marshal(transPolicyJson)
						if err != nil {
							continue
						} else {
							redisCenter.RedisSubProducerTransPolicy(string(transPolicy))
							fmt.Println(string(transPolicy))
						}

					} else {
						policyCenter.UpdateTransPolicyStatusByTaskid(comm.AgentOutTimeStatus, list[i].TaskId)
						fmt.Println(fmt.Sprintf("taskid %v status -2", list[i].TaskId))
					}
				}

			}
		}
		time.Sleep(time.Minute)
	}
}

//监控数据库中哪些下线透明转发策略未收到agent的返回信息，一分钟下发一次，三分钟内下发三次左右，三次之后均无响应，则置状态为-2
func MonitorForAgentByOfflineTransPolicy() {
	for {
		list, _ := policyCenter.SelectOfflineTransPolicyNotFinish()
		fmt.Println(list)
		if len(list) > 0 {
			for i := 0; i < len(list); i++ {
				fmt.Println(fmt.Sprintf(" taskid status is %v", list[i].Status))
				if list[i].Status == 1 {
					localTimes := time.Now().Unix()
					lastTimes, _ := strconv.ParseInt(list[i].OfflineTime, 10, 64)
					sub := (localTimes - lastTimes) / 60
					if sub >= 0 && sub <= 3 {
						var transPolicyJson comm.TransparentTransponderPolicyJson
						transPolicyJson.AgentId = list[i].AgentId
						transPolicyJson.TaskId = list[i].TaskId
						transPolicyJson.ListenPort = list[i].ForwardPort
						transPolicyJson.ServerType = list[i].HoneyPotType
						transPolicyJson.HoneyIP = list[i].HoneyIP
						transPolicyJson.HoneyPort = list[i].HoneyPotPort
						transPolicyJson.Status = comm.OfflineStatus

						transPolicy, err := json.Marshal(transPolicyJson)
						if err != nil {
							continue
						} else {
							redisCenter.RedisSubProducerTransPolicy(string(transPolicy))
							fmt.Println(string(transPolicy))
						}

					} else {
						policyCenter.UpdateTransPolicyStatusByTaskid(comm.AgentOutTimeStatus, list[i].TaskId)
						fmt.Println(fmt.Sprintf("taskid %v status -2", list[i].TaskId))
					}
				}

			}
		}
		time.Sleep(time.Minute)
	}
}

//监控数据库中哪些上线蜜罐流量转发策略未收到agent的返回信息，一分钟下发一次，三分钟内下发三次左右，三次之后均无响应，则置状态为-2
func MonitorForAgentByCreateHoneyTransPolicy() {
	for {
		list, _ := policyCenter.SelectCreateHoneyTransPolicyNotFinish()
		fmt.Println(list)
		if len(list) > 0 {
			for i := 0; i < len(list); i++ {
				fmt.Println(fmt.Sprintf(" taskid status is %v", list[i].Status))
				if list[i].Status == 0 {
					localTimes := time.Now().Unix()
					lastTimes, _ := strconv.ParseInt(list[i].CreateTime, 10, 64)
					sub := (localTimes - lastTimes) / 60
					if sub >= 0 && sub <= 3 {
						var transPolicyJson comm.TransparentTransponderPolicyJson
						transPolicyJson.AgentId = list[i].AgentId
						transPolicyJson.TaskId = list[i].TaskId
						transPolicyJson.ListenPort = list[i].ForwardPort
						transPolicyJson.ServerType = list[i].HoneyPotType
						transPolicyJson.HoneyIP = list[i].HoneyIP
						transPolicyJson.HoneyPort = list[i].HoneyPotPort
						transPolicyJson.Status = comm.CreateStatus
						transPolicyJson.Type = comm.AgentTypeEdge
						transPolicyJson.Path = ""

						transPolicy, err := json.Marshal(transPolicyJson)
						if err != nil {
							continue
						} else {
							redisCenter.RedisSubProducerTransPolicy(string(transPolicy))
							fmt.Println(string(transPolicy))
						}

					} else {
						policyCenter.UpdateHoneyTransPolicyStatusByTaskid(comm.AgentOutTimeStatus, list[i].TaskId)
						fmt.Println(fmt.Sprintf("taskid %v status -2", list[i].TaskId))
					}
				}

			}
		}
		time.Sleep(time.Minute)
	}
}

//监控数据库中哪些下线蜜罐流量转发策略未收到agent的返回信息，一分钟下发一次，三分钟内下发三次左右，三次之后均无响应，则置状态为-2
func MonitorForAgentByOfflineHoneyTransPolicy() {
	for {
		list, _ := policyCenter.SelectOfflineHoneyTransPolicyNotFinish()
		fmt.Println(list)
		if len(list) > 0 {
			for i := 0; i < len(list); i++ {
				fmt.Println(fmt.Sprintf(" taskid status is %v", list[i].Status))
				if list[i].Status == 1 {
					localTimes := time.Now().Unix()
					lastTimes, _ := strconv.ParseInt(list[i].OfflineTime, 10, 64)
					sub := (localTimes - lastTimes) / 60
					if sub >= 0 && sub <= 3 {
						var transPolicyJson comm.TransparentTransponderPolicyJson
						transPolicyJson.AgentId = list[i].AgentId
						transPolicyJson.TaskId = list[i].TaskId
						transPolicyJson.ListenPort = list[i].ForwardPort
						transPolicyJson.ServerType = list[i].HoneyPotType
						transPolicyJson.HoneyIP = list[i].HoneyIP
						transPolicyJson.HoneyPort = list[i].HoneyPotPort
						transPolicyJson.Status = comm.OfflineStatus

						transPolicy, err := json.Marshal(transPolicyJson)
						if err != nil {
							continue
						} else {
							redisCenter.RedisSubProducerTransPolicy(string(transPolicy))
							fmt.Println(string(transPolicy))
						}

					} else {
						policyCenter.UpdateTransPolicyStatusByTaskid(comm.AgentOutTimeStatus, list[i].TaskId)
						fmt.Println(fmt.Sprintf("taskid %v status -2", list[i].TaskId))
					}
				}

			}
		}
		time.Sleep(time.Minute)
	}
}

func SignFileUpload(r *http.Request, signname string) (string, error, bool) {
	isfalse := false
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")

	if strings.Contains(handler.Filename, "..") || strings.Contains(handler.Filename, "/") {
		logs.Error("Filename invalid: ", handler.Filename)
		isfalse = true
		return "", errors.New("filename invalid"), isfalse
	}
	if err != nil {
		isfalse = true
		return "", err, isfalse
	}
	defer file.Close()
	filename := handler.Filename
	if path.Ext(filename) == ".doc" || path.Ext(filename) == ".docx" || path.Ext(filename) == ".pdf" || path.Ext(filename) == ".ppt" || path.Ext(filename) == ".pptx" || path.Ext(filename) == ".xls" || path.Ext(filename) == ".xlsx" {
		err = os.MkdirAll("upload/honeytoken/", os.ModePerm)
		if err != nil {
			isfalse = true
			return "", err, isfalse
		}
		err = os.MkdirAll("upload/honeytoken/"+signname+"/", os.ModePerm)
		if err != nil {
			isfalse = true
			return "", err, isfalse
		}
		filepath := "upload/honeytoken/" + signname + "/" + filename
		f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			isfalse = true
			return "", err, isfalse
		}
		defer f.Close()
		io.Copy(f, file)
	} else {
		isfalse = true
	}
	return filename, nil, isfalse
}

func ProtocolFileUpload(r *http.Request) (string, string, error) {
	protocolname := ""
	protocolfilename := ""
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if file != nil {
		protocolname = r.FormValue("protocolname")
		if protocolname != "" {
			if strings.Contains(handler.Filename, "..") || strings.Contains(handler.Filename, "/") {
				logs.Error("Filename invalid: ", handler.Filename)
				return "", "", errors.New("filename invalid")
			}
			if err != nil {
				return "", "", err
			}
			defer file.Close()
			protocolfilename = protocolname + "proxy"
			err = os.MkdirAll(path.Join("upload/protocol/", protocolname), os.ModePerm)
			if err != nil {
				return "", "", err
			}
			err = os.MkdirAll(path.Join("upload/protocol/", protocolname), os.ModePerm)
			if err != nil {
				return "", "", err
			}
			filepath := path.Join("upload/protocol/", protocolname, protocolfilename)
			f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, os.ModeAppend|os.ModePerm)
			if err != nil {
				return "", "", err
			}
			defer f.Close()
			io.Copy(f, file)
		}
	}

	return protocolfilename, protocolname, nil
}

func FileUpload(r *http.Request, baitname string) (string, string, error) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")

	if strings.Contains(handler.Filename, "..") || strings.Contains(handler.Filename, "/") {
		logs.Error("Filename invalid: ", handler.Filename)
		return "", "", errors.New("filename invalid")
	}
	if err != nil {
		return "", "", err
	}
	defer file.Close()
	filename := handler.Filename
	err = os.MkdirAll("upload/", os.ModePerm)
	if err != nil {
		return "", "", err
	}
	err = os.MkdirAll("upload/"+baitname+"/", os.ModePerm)
	if err != nil {
		return "", "", err
	}
	filepath := "upload/" + baitname + "/" + filename

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		return "", "", err
	}
	defer f.Close()
	io.Copy(f, file)
	return "", filename, nil
}

func Download(rw http.ResponseWriter, r *http.Request) {
	filePath := util.GetCurrentPathString() + "/agent/decept-agent.tar.gz"
	file, err := os.Open(filePath)
	if err != nil {
		logs.Error("open decept-agent.tar.gz fail: %v , %s", err, filePath)
		return
	}
	defer file.Close()
	fileHeader := make([]byte, 1024)
	file.Read(fileHeader)
	fileStat, _ := file.Stat()
	rw.Header().Set("Content-Disposition", "attachment; filename=decept-agent.tar.gz")
	rw.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	rw.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	file.Seek(0, 0)
	io.Copy(rw, file)
	return
}

func RefreshImages(w http.ResponseWriter, r *http.Request) {
	GetPodImageList()
}

func RefreshPods(w http.ResponseWriter, r *http.Request) {
	go k3s.FreshPods()
}
