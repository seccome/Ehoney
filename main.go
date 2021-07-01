package main

import (

	//httpSwagger "github.com/swaggo/http-swagger"
	"database/sql"
	"decept-defense/controllers"
	"decept-defense/models"
	"decept-defense/models/honeycluster"
	"decept-defense/models/redisCenter"
	"decept-defense/models/util"
	"decept-defense/models/util/comhttp"
	"decept-defense/models/util/comm"
	"decept-defense/models/util/k3s"
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/gorilla/mux"
	"github.com/lestrrat/go-file-rotatelogs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func makeRouter() *mux.Router {
	r := mux.NewRouter()
	/*登录认证*/
	r.PathPrefix("/upload/").Handler(http.StripPrefix("/upload/", http.FileServer(http.Dir("./upload/"))))

	/*大屏*/
	r.HandleFunc("/deceptdefense/api/gettopattackmap", GetTopAttackMap)
	r.HandleFunc("/deceptdefense/api/gettopattackips", GetTopAttackIps)
	r.HandleFunc("/deceptdefense/api/gettopsourceips", GetTopSourceIps)
	r.HandleFunc("/deceptdefense/api/gettopareas", GetTopAreas)
	r.HandleFunc("/deceptdefense/api/gettopattacktypes", GetTopAttackTypes)

	/*管理员登录*/
	r.HandleFunc("/deceptdefense/api/login", AdminLogin)
	r.HandleFunc("/deceptdefense/api/logout", AdminLogout)

	/*镜像列表更新*/
	r.HandleFunc("/deceptdefense/api/refreshimages", RefreshImages)

	r.HandleFunc("/deceptdefense/api/agent/download", Download)

	//add route for 思瀚
	r.HandleFunc("/deceptdefense/api/dbportmap", controllers.Map)

	/*容器列表更新*/
	r.HandleFunc("/deceptdefense/api/refreshpods", RefreshPods)

	/*病毒信息插入*/
	r.HandleFunc("/deceptdefense/api/addclamavresult", InsertClamavResultHandler)

	/*日志插入*/
	r.HandleFunc("/deceptdefense/api/insertFalcoLog", InsertFalcoLogHandler)

	r.HandleFunc("/deceptdefense/api/insertAttackLog", InsertAttackLogHandler)

	/*SSH Key 插入*/
	r.HandleFunc("/deceptdefense/api/insertsshkey", InsertSSHKeyHandler)

	/*服务器心跳*/
	r.HandleFunc("/deceptdefense/api/getapplocationssignmsg", ApplicationSignMsgHandler)

	/*心态检测*/
	r.HandleFunc("/health", healthHandler)

	//r.HandleFunc("/test",testHandle)

	api := r.PathPrefix("/deceptdefense/api").Subrouter()
	/*应用集群管理*/
	api.HandleFunc("/getapplicationclusters", GetApplicationClusters)
	api.HandleFunc("/getapplicationlists", GetApplicationLists)
	api.HandleFunc("/applicationclusteradd", CreateApplicationClusters)

	/*应用集群诱饵策略管理*/
	api.HandleFunc("/applicationbaitpolicyadd", CreateApplicationBaitPolicyHandler)
	api.HandleFunc("/applicationbaitpolicydelete", DeleteApplicationBaitPolicyHandler)
	api.HandleFunc("/downloadapplicationbaitpolicy", DownloadApplicationBaitPolicyHandler)

	/*应用集群密签策略管理*/
	api.HandleFunc("/signpolicyselectagentid", SelectSignPolicyHandler)
	api.HandleFunc("/applicationsignpolicyadd", CreateApplicationSignPolicyHandler)
	api.HandleFunc("/applicationsignpolicydelete", DeleteApplicationSignPolicyHandler)
	api.HandleFunc("/downloadapplicationsign", DownloadApplicationSignPolicyHandler)


	/*策略管理*/
	/*诱饵策略管理*/
	api.HandleFunc("/baitPolicyDelete", DeleteBaitPolicyHandler)
	api.HandleFunc("/baitPolicySelectAgentId", SelectBaitPolicyHandler)
	api.HandleFunc("/getAllBaitType", SelectAllBaitTypeHandler)
	api.HandleFunc("/getallhoneybaittype", SelectAllHoneyBaitTypeHandler)


	/*透明转发策略管理*/
	api.HandleFunc("/transPolicyAdd", CreateTransPolicyHandler)
	api.HandleFunc("/transPolicyDelete", DeleteTransPolicyHandler)
	api.HandleFunc("/transPolicySelectAgentId", SelectTransPolicyHandler)
	api.HandleFunc("/transpolicytest", TestTransPolicyHandler)

	/*协议转发策略管理*/
	api.HandleFunc("/honeytranspolicyadd", CreateHoneyTransPolicyHandler)
	api.HandleFunc("/honeytranspolicydelete", DeleteHoneyTransPolicyHandler)
	api.HandleFunc("/honeytranspolicyselectagentid", SelectHoneyTransPolicyHandler)
	api.HandleFunc("/honeytranspolicytest", TestHoneyTransPolicyHandler)

	/*攻击事件管理*/
	api.HandleFunc("/attackLogList", SelectAttackLogListHandler)
	api.HandleFunc("/attackLogDetail", SelectAttackLogDetailHandler)

	/*告警管理*/
	api.HandleFunc("/addConfig", AddConfigHandler)
	api.HandleFunc("/selectConfig", SelectConfigHandler)
	api.HandleFunc("/deleteConfig", DeleteConfigHandler)

	/*蜜罐集群管理*/
	api.HandleFunc("/addhoneyserver", CreateHoneyServer)

	/*蜜罐容器管理*/
	//api.HandleFunc("/podadd", AddPod)
	api.HandleFunc("/podadd", AddPodv2)
	api.HandleFunc("/deletepod", DeletePod)
	api.HandleFunc("/getpodimage", GetPodImage)
	api.HandleFunc("/podstatuscheck", CheckPodNetStatus) //新建蜜罐状态监测
	api.HandleFunc("/testpodstatus", TestCheckPodStatus) //蜜罐网络情况探测

	/*蜜罐集群管理*/
	api.HandleFunc("/gethoneyclusters", GetHoneyClusters)

	/*蜜罐管理*/
	api.HandleFunc("/gethoneylist", GetHoneyList)
	api.HandleFunc("/gethoneyinfos", GetHoneyInfos)
	api.HandleFunc("/gethoneytransinfos", GetHoneyTransInfos)
	api.HandleFunc("/gethoneypotstype", GetHoneyPotsType)
	api.HandleFunc("/gethoneylistfortrans", GetHoneyListForTrans)
	api.HandleFunc("/gethoneyforwardlistfortrans", GetHoneyForwardListForTrans)
	api.HandleFunc("/gethoneytransportslist", GetHoneyTransPortsListForTrans)

	/*操作系统类型*/
	api.HandleFunc("/getsystype", SelectAllSysTypeHandler)

	/*蜜罐密签管理*/
	//api.HandleFunc("/gethoneysigns",)
	api.HandleFunc("/honeysignpolicyadd", CreateHoneySignPolicyHandler)
	api.HandleFunc("/honeysignpolicydelete", DeleteHoneySignPolicyHandler)
	api.HandleFunc("/gethoneysigns", GetHoneySigns)
	api.HandleFunc("/gethoneysignmsg", HoneySignMsgHandler)
	api.HandleFunc("/downloadhoneysign", DownloadHoneySignPolicyHandler)

	/*蜜罐诱饵管理*/
	api.HandleFunc("/honeybaitpolicyadd", CreateHoneyBaitPolicyHandler)
	api.HandleFunc("/honeybaitpolicydetele", DeleteHoneyBaitPolicyHandler)
	api.HandleFunc("/gethoneybaits", GetHoneyBaits)
	api.HandleFunc("/downloadhoneybaits", DownloadHoneyBaits)

	/*病毒扫描结果*/
	api.HandleFunc("/getclamavresult", GetClamavResultHandler)

	/*诱饵管理*/
	api.HandleFunc("/createbaits", CreateBaitHandler)
	api.HandleFunc("/getbaits", GetBaits)
	api.HandleFunc("/getbaitsbytype", GetBaitsByType)
	api.HandleFunc("/baitdelete", DeleteBaitHandler)

	/*密签管理*/
	api.HandleFunc("/createsign", CreateSignHandler)
	api.HandleFunc("/deletesign", DeleteSignHandler)
	api.HandleFunc("/getsigns", GetSigns)
	api.HandleFunc("/getsigntype", GetSignType)
	api.HandleFunc("/getsignsbytype", GetSignsByType)

	/*协议转发模块管理*/
	api.HandleFunc("/createprotocoltype", CreateProtocolType)
	api.HandleFunc("/getprotocoltype", GetProtocolType)
	api.HandleFunc("/deleteprotocoltype", DeleteProtocolType)

	/*镜像列表模块*/
	api.HandleFunc("/gethoneyimagelists", GetHoneyImageLists)
	api.HandleFunc("/honeyimageedit", UpdateHoneyImage)

	/**
	1、文件上传，生产。有了。通过redis策略下发
	2、协议模块删除，文件删除没写，通过redis策略下发，文件删除
	*/

	/*配置管理*/
	api.HandleFunc("/deleteConfig", DeleteConfigHandler)
	api.HandleFunc("/addharborconfig", AddHarborConfig)
	api.HandleFunc("/getharborinfo", GetHarborInfo)
	api.HandleFunc("/deleteharborinfo", DeleteHarborInfo)

	/*redis 初始化*/
	api.HandleFunc("/addredisconfig", AddRedisConfigHandler)
	api.HandleFunc("/getredisinfo", GetRedisInfo)
	api.HandleFunc("/deleteredisinfo", DeleteRedisInfo)

	/*密签追踪ip*/
	api.HandleFunc("/addtraceinfo", AddTraceHostHandler)
	api.HandleFunc("/gettraceinfo", GetTraceHostHandler)

	api.Use(BeforeAction)

	return r
}

var (
	daemon   = flag.Bool("daemon", false, "run as daemon")
	Addr     = flag.String("addr", ":8082", "default web")
	MysqlCon = flag.String("mysql", "", "mysql server connect")

	DbCon *sql.DB
)

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}

	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	return string(path), nil
}

func initiLogger() {
	appPath, err := GetCurrentPath()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	dirName := fmt.Sprintf("%slog/", appPath)
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		fmt.Println(dirName, err)
		log.Fatal(dirName, err)
	}

	fileName := fmt.Sprintf("%s/deceptdefense.%s", dirName, "%Y%m%d%H%M")

	writer, err := rotatelogs.New(fileName,
		rotatelogs.WithLinkName(""),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		fmt.Println(err)
		log.Fatalf("[initiLogger] Failed to Initialize Log File :%v", err)
	}

	log.SetFlags(log.Lshortfile | log.Llongfile | log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
	log.SetOutput(writer)
	err1 := util.InitLogger() //初始化调用logger
	if err1 != nil {
		fmt.Println("[MAIN] logger init error!")
	}

	//获取数据库配置
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	err_mysql := orm.RegisterDataBase("default", "mysql", dbuser+":"+dbpassword+"@tcp("+dbhost+":"+dbport+")/"+dbname+"?charset=utf8&parseTime=true&loc=Local")
	if err_mysql != nil {
		logs.Error("[MAIN] mysql connect error:", err_mysql)
	}
	orm.Debug = true

	// 初始化tor列表
	err2 := models.SetupExitNodeMap()
	if err2 != nil {
		logs.Error(err2)
	}

	tracehost := beego.AppConfig.String("tracehost")
	if tracehost != "" {
		createtime := time.Now().Unix()
		traceid := util.GetStrMd5("ehoney")
		honeycluster.InsertTraceInfo(traceid, tracehost, createtime)
	}

	// Redis消费者
	go redisCenter.RedisPubConsumerBaitPolicyResponse()  //业务服务器诱饵、密签
	go redisCenter.RedisPubConsumerTransPolicyResponse() //透明、协议转发
	go redisCenter.RedisPubConsumerTransEventResponse()

	// 根据心跳注册主机
	go redisCenter.RedisPubConsumerServerRegResponse()

	go k3s.FreshPods()
	return
}

func parseStartupParams(configs string) {
	if configs == "" {
		return
	}
	paramsArray := strings.Split(configs, ";")
	for _, param := range paramsArray {
		log.Println("start pares param: " + param)

		keyValArray := strings.Split(param, ":")

		if len(keyValArray) == 2 {
			log.Println("set config ", keyValArray[0], " : ", keyValArray[1])
			beego.AppConfig.Set(keyValArray[0], keyValArray[1])
		}
	}
}

func main() {

	configs := flag.String("CONFIGS", "xxxx", "configs")

	flag.Parse()

	log.Printf("find startup parameter %s", *configs)

	parseStartupParams(*configs)

	initiLogger()

	if len(*MysqlCon) <= 0 {
		*MysqlCon = dbuser + ":" + dbpassword + "@tcp(" + dbhost + ")/" + dbname + "?charset=utf8&loc=Asia%2FShanghai"
	}

	sqlCon, err := sql.Open("mysql", *MysqlCon)
	if err != nil {
		log.Printf("open mysql fail %v\n", err)
		return
	}

	DbCon = sqlCon
	defer sqlCon.Close()

	log.Printf("start web server\n")
	DbCon.SetConnMaxLifetime(300 * time.Second)
	DbCon.SetMaxOpenConns(120)
	DbCon.SetMaxIdleConns(12)

	r := makeRouter()
	http.Handle("/", r)
	http.ListenAndServe(*Addr, r)
}

func BeforeAction(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			comhttp.SendJSONResponse(w, comm.Response{Code: comm.SuccessCode, Data: nil, Message: ""})
		} else {
			islogin := false
			name := sess.Start(w, r).GetString("name")
			if name != "" {
				islogin = true
			}
			if !islogin {
				comhttp.SendJSONResponse(w, comm.Response{Code: comm.LoginoutCode, Data: nil, Message: "no login"})
			} else {
				h.ServeHTTP(w, r)
			}
		}
	})
}
