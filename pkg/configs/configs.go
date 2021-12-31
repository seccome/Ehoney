package configs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/url"
	"strings"
)

var setting Config

type Config struct {
	Log      LogSetting
	Server   ServerSetting
	Database DatabaseSetting
	Redis    RedisSetting
	App      AppSetting
	Harbor   HarborSetting
}

type LogSetting struct {
	Level            string
	LogFormatJson    string
	LogFormatConsole string
	TimeKey          string
	LevelKey         string
	NameKey          string
	CallerKey        string
	StackTraceKey    string
	MessageKey       string
	MaxSize          int
	MaxBackups       int
	MaxAge           int
}

type ServerSetting struct {
	RunMode  string
	HttpPort int
	AppHost  string
}

type DatabaseSetting struct {
	DBType     string
	DBUser     string
	DBPassword string
	DBPort     string
	DBHost     string
	DBName     string
}

type RedisSetting struct {
	RedisHost     string
	RedisPort     int
	RedisPassword string
}

type AppSetting struct {
	EnableCors         bool
	ProtocolDeployPath string
	TaskChannel        string
	ScriptPath         string
	UploadPath         string
	TokenTraceAddress  string
	TokenTraceApiPath  string
	EngineVersion      string
	WebHook            string
	TraceHost          string
	Extranet           string
}

type HarborSetting struct {
	HarborURL     string
	HarborProject string
	User          string
	Password      string
	APIVersion    string
}

func SetUp() {
	viper.BindEnv("WorkingDir")
	viper.AddConfigPath("configs")
	viper.SetConfigType("toml")
	viper.SetConfigName("configs")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.Unmarshal(&setting)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.Unmarshal(&setting)
	})
}

func UpdateIPConfig(ipConfig string) {
	data := strings.Split(ipConfig, ";")
	fmt.Println("update ip config IP value : " + ipConfig)
	zap.L().Debug("update ip config IP value : " + ipConfig)
	for _, d := range data {
		if d == "" {
			continue
		}
		p := strings.Split(d, ":")
		if p[0] == "host" {
			fmt.Println("start config IP value, current value: " + setting.Server.AppHost)
			zap.L().Debug("start config IP value, current value: " + setting.Server.AppHost)
			data, _ := url.Parse(setting.App.TokenTraceAddress)
			traceAddress := setting.App.TokenTraceAddress
			if data != nil {
				traceAddress = data.Scheme + "://" + p[1] + ":" + data.Port()
			}

			viper.Set("redis.redishost", p[1])
			viper.Set("database.dbhost", p[1])
			viper.Set("server.apphost", p[1])
			viper.Set("app.extranet", p[1])
			viper.Set("app.tokentraceaddress", traceAddress)
			setting.Redis.RedisHost = p[1]
			setting.Database.DBHost = p[1]
			setting.Server.AppHost = p[1]
			setting.App.TokenTraceAddress = traceAddress
			setting.App.Extranet = p[1]
			err := viper.WriteConfig()
			if err != nil {
				zap.L().Error("write config err " + err.Error())
				fmt.Println("write config err " + err.Error())
				fmt.Println("try again")
				viper.WriteConfig()
				zap.L().Error("write config err " + err.Error())
				zap.L().Error("请手动修改相关配置为主机IP地址")
				fmt.Println("配置文件修改异常、请手动修改相关配置为主机IP地址")
			}
			zap.L().Debug("finish config IP value, current value: " + setting.Server.AppHost)
			fmt.Println("finish config IP value, current value: " + setting.Server.AppHost)
		}
	}
}

func GetSetting() *Config {
	return &setting
}

func ProjectName() string {
	return "decept-defense"
}

func ProjectLogFile() string {
	return fmt.Sprintf("./logs/%s_.log", ProjectName())
}
