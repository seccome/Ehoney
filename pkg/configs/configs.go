package configs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
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
	EnableCors        bool
	ScriptPath        string
	UploadPath        string
	TokenTraceApiPath string
	EngineVersion     string
	WebHook           string
	TraceHost         string
	Extranet          string
}

type HarborSetting struct {
	HarborURL     string
	HarborProject string
	User          string
	Password      string
	APIVersion    string
}

func SetUp() {
	_ = viper.BindEnv("WorkingDir")
	viper.AddConfigPath("configs")
	viper.SetConfigType("toml")
	viper.SetConfigName("configs")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	_ = viper.Unmarshal(&setting)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		_ = viper.Unmarshal(&setting)
	})
}

func UpdateIPConfig(hostIp string) {
	viper.Set("database.dbhost", hostIp)
	viper.Set("server.apphost", hostIp)
	setting.Redis.RedisHost = hostIp
	setting.Database.DBHost = hostIp
	setting.Server.AppHost = hostIp
	if setting.App.Extranet == "http://localhost:8082" {
		viper.Set("app.extranet", fmt.Sprintf("http://%s:8082", hostIp))
		setting.App.Extranet = fmt.Sprintf("http://%s:8082", hostIp)
	}
	time.Sleep(100)
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
