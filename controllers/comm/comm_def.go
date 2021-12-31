package comm

type SelectVirusPayload struct {
	StartTimestamp int64  `json:"StartTimestamp"`
	EndTimestamp   int64  `json:"EndTimestamp" `
	VirusName      string `json:"VirusName"`
	HoneypotName   string `json:"HoneypotName" `
	PageNumber     int    `json:"PageNumber" binding:"required"`
	PageSize       int    `json:"PageSize" binding:"required"`
	VirusFilePath  string `json:"VirusFilePath" `
}

type HoneypotSelectResultPayload struct {
	ID           int64      `json:"ID"`           //蜜罐ID
	ServerType   string     `json:"ServerType"`   //蜜罐类型
	HoneypotName string     `json:"HoneypotName"` //蜜罐名称
	HoneypotIP   string     `json:"HoneypotIP"`   //蜜罐IP
	ServerIP     string     `json:"ServerIP"`     //蜜网IP
	CreateTime   string     `json:"CreateTime"`   //创建时间
	Status       TaskStatus `json:"Status"`       //蜜罐状态
	Creator      string     `json:"Creator"`      //创建用户
}

type ProbeSelectResultPayload struct {
	ID            int64      `json:"ID"`            //探针ID
	ProbeIP       string     `json:"ProbeIP"`       //探针IP
	HostName      string     `json:"HostName"`      //探针名称
	HeartbeatTime string     `json:"HeartbeatTime"` //心跳时间
	CreateTime    string     `json:"CreateTime"`    //注册时间
	SystemType    string     `json:"SystemType"`    //系统类型
	Status        TaskStatus `json:"Status"`        //状态
}

type TransProxySelectPayload struct {
	ProxyPort          int32  `json:"ProxyPort"`
	ProtocolType       string `json:"ProtocolType"`
	ProbeIP            string `json:"ProbeIP"`
	HoneypotServerIP   string `json:"HoneypotServerIP"`
	HoneypotServerPort int32  `json:"HoneypotServerPort"`
	StartTimestamp     int64  `json:"StartTimestamp"`
	EndTimestamp       int64  `json:"EndTimestamp"`
	PageNumber         int    `json:"PageNumber" binding:"required"`
	PageSize           int    `json:"PageSize" binding:"required"`
}

type FileBaitSelectResultPayload struct {
	ID         int64  `json:"ID"`         //文件诱饵ID
	BaitName   string `json:"BaitName"`   //诱饵名称
	FileName   string `json:"FileName"`   //文件名称
	BaitType   string `json:"BaitType"`   //诱饵类型
	BaitData   string `json:"BaitData"`   //HISTORY诱饵数据
	Creator    string `json:"Creator"`    //创建用户
	CreateTime string `json:"CreateTime"` //创建时间
}

type HistoryBaitSelectResultPayload struct {
	ID         int64  `json:"ID"`
	BaitName   string `json:"BaitName"`
	BaitType   string `json:"BaitType"`
	Creator    string `json:"Creator"`
	BaitData   string `json:"BaitData"`
	CreateTime string `json:"CreateTime"`
}

type TokenSelectResultPayload struct {
	ID          int64  `json:"ID"`          //密签ID
	FileName    string `json:"FileName"`    //文件名称
	TokenName   string `json:"TokenName"`   //蜜签名称
	TokenType   string `json:"TokenType"`   //蜜签类型
	Creator     string `json:"Creator"`     //创建用户
	CreateTime  string `json:"CreateTime"`  //创建时间
	DefaultFlag bool   `json:"DefaultFlag"` //默认属性
}

type ProtocolProxySelectResultPayload struct {
	ID           int64  `json:"ID"`           //协议代理ID
	HoneypotIP   string `json:"HoneypotIP"`   //蜜罐IP
	HoneypotName string `json:"HoneypotName"` //蜜罐名称
	ProxyPort    int32  `json:"ProxyPort"`    //代理端口
	ServerIP     string `json:"ServerIP"`     //蜜网IP
	ServerType   string `json:"ServerType"`   //蜜罐服务
	ServerPort   string `json:"ServerPort"`   //蜜罐端口
	Creator      string `json:"Creator"`      //创建用户
	CreateTime   string `json:"CreateTime"`   //创建时间
	Status       int    `json:"Status"`       //状态
	MinPort      int32  `json:"MinPort"`      //协议MinPort
	MaxPort      int32  `json:"MaxPort"`      //协议MaxPort
}

type TransparentProxySelectResultPayload struct {
	ID           int64  `json:"ID"`           //透明代理ID
	ProbeIP      string `json:"ProbeIP"`      //探针IP
	ProxyPort    int32  `json:"ProxyPort"`    //代理端口
	Creator      string `json:"Creator"`      //创建用户
	CreateTime   string `json:"CreateTime"`   //创建时间
	ProtocolPort int32  `json:"ProtocolPort"` //协议代理端口
	ProtocolType string `json:"ProtocolType"` //协议类型
	Status       int    `json:"Status"`       //状态
}

type ServerBaitSelectPayload struct {
	ServerID   int64  `json:"ServerID" binding:"required"`
	PageNumber int    `json:"PageNumber" binding:"required"`
	PageSize   int    `json:"PageSize" binding:"required"`
	Payload    string `json:"Payload"`
}

type ServerTokenSelectPayload struct {
	ServerID   int64  `json:"ServerID" binding:"required"`   //服务器ID
	PageNumber int    `json:"PageNumber" binding:"required"` //页面number
	PageSize   int    `json:"PageSize" binding:"required"`   //页面size
	Payload    string `json:"Payload"`                       //查找payload
}

type ServerBaitSelectResultPayload struct {
	ID         int64      `json:"ID"`
	BaitName   string     `json:"BaitName"`
	BaitType   string     `json:"BaitType"`
	Creator    string     `json:"Creator"`
	DeployPath string     `json:"DeployPath"`
	CreateTime string     `json:"CreateTime"`
	Status     TaskStatus `json:"Status"`
}

type ServerTokenSelectResultPayload struct {
	ID         int64      `json:"ID"`         //蜜罐密签ID
	TokenName  string     `json:"TokenName"`  //密签名称
	TokenType  string     `json:"TokenType"`  //密签类型
	Creator    string     `json:"Creator"`    //创建用户
	CreateTime string     `json:"CreateTime"` //创建时间
	DeployPath string     `json:"DeployPath"` //部署路径
	Status     TaskStatus `json:"Status"`     //状态
}

type AttackSelectResultPayload struct {
	//ID             int64  `json:"ID"`             //攻击日志ID
	AttackIP       string            `json:"AttackIP"`       //攻击IP
	ProbeIP        string            `json:"ProbeIP"`        //探针IP
	JumpIP         string            `json:"JumpIP"`         //跳转IP
	HoneypotIP     string            `json:"HoneypotIP"`     //蜜罐IP
	ProtocolType   string            `json:"ProtocolType"`   //协议类型
	AttackTime     string            `json:"AttackTime"`     //攻击时间
	AttackLocation string            `json:"AttackLocation"` //攻击位置
	AttackDetail   string            `json:"AttackDetail"`   //攻击详情
	CounterInfo    map[string]string `json:"CounterInfo"`    //反制详情
}

type TokenTraceSelectResultPayload struct {
	ID        int64  `json:"ID"`        //攻击日志ID
	TokenType string `json:"TokenType"` //密签类型
	TokenName string `json:"TokenName"` //密签名称
	OpenTime  string `json:"OpenTime"`  //攻击时间
	OpenIP    string `json:"OpenIP"`    //攻击IP
	UserAgent string `json:"UserAgent"` //用户UA
	Location  string `json:"Location"`  //攻击位置
}

type FalcoSelectResultPayload struct {
	ID           int64  `json:"ID"`           //falco攻击日志ID
	HoneypotName string `json:"HoneypotName"` //蜜罐名称
	Event        string `json:"Event"`        //falco事件
	Time         string `json:"Time"`         //发生时间
	Output       string `json:"Output"`       //输出
	Level        string `json:"Level"`        //事件等级
	FileFlag     bool   `json:"FileFlag"`     //文件标记
	DownloadPath string `json:"DownloadPath"` //文件下载链接
}

type ImageUpdatePayload struct {
	ImagePort int32  `form:"ImagePort" binding:"required"` //镜像端口
	ImageType string `form:"ImageType" binding:"required"` //协议类型
}

type SelectPayload struct {
	Payload    string `json:"Payload"`                       //查找payload
	PageNumber int    `json:"PageNumber" binding:"required"` //页number
	PageSize   int    `json:"PageSize" binding:"required"`   //页size
}

type AttackEventSelectPayload struct {
	SelectPayload
	StartTime    string `json:"StartTime"`    //开始时间
	EndTime      string `json:"EndTime"`      //结束时间
	AttackIP     string `json:"AttackIP"`     //攻击IP
	JumpIP       string `json:"JumpIP"`       //跳转IP
	ProbeIP      string `json:"ProbeIP"`      //探针IP
	HoneypotIP   string `json:"HoneypotIP"`   //蜜罐IP
	ProtocolType string `json:"ProtocolType"` //协议类型
}

type FalcoEventSelectPayload struct {
	SelectPayload
	StartTime string `json:"StartTime"` //开始时间
	EndTime   string `json:"EndTime"`   //结束时间
}

type BatchSelectPayload struct {
	Ids []int64 `json:"Ids"`
}

type TokenTraceSelectPayload struct {
	SelectPayload
	ServerType string `json:"ServerType"` //服务类型 蜜罐（honeypot）or 探针（probe）
	AttackIP   string `json:"AttackIP"`   //攻击IP
	StartTime  string `json:"StartTime"`  //开始时间
	EndTime    string `json:"EndTime"`    //结束时间
}

type AttackTraceSelectPayload struct {
	SelectPayload
	Type         string `json:"Type"`         //攻击类型
	AttackIP     string `json:"AttackIP"`     //攻击IP
	HoneypotIP   string `json:"HoneypotIP"`   //蜜罐IP
	ProtocolType string `json:"ProtocolType"` //协议类型
	StartTime    string `json:"StartTime"`    //开始时间
	EndTime      string `json:"EndTime"`      //结束时间
}

type BaitSelectPayload struct {
	SelectPayload
	BaitType string `json:"BaitType"` //诱饵类型
}

type HoneypotSelectPayload struct {
	SelectPayload
	ProtocolType string `json:"ProtocolType"` //协议类型
}

type SelectTransparentProxyPayload struct {
	SelectPayload
	ProtocolProxyID int64 `json:"ProtocolProxyID"` //协议代理ID
	Status          int64 `json:"Status"`          //: 2 下线 3 上线
}

type AttackStatistics struct {
	Data  string `json:"Data"`  //攻击信息
	Count int64  `json:"Count"` //出现次数
}

type ProtocolSelectResultPayload struct {
	ID           int64  `json:"ID"`           //协议ID
	ProtocolType string `json:"ProtocolType"` //协议类型
	DeployPath   string `json:"DeployPath"`   //部署路径
	Status       int    `json:"Status"`       //状态
	Creator      string `json:"Creator"`      //创建用户
	CreateTime   string `json:"CreateTime"`   //创建时间
	MinPort      int32  `json:"MinPort"`      //端口low
	MaxPort      int32  `json:"MaxPort"`      //端口high
	DefaultFlag  bool   `json:"DefaultFlag"`  //默认属性
}

type ProtocolHoneypotSelectResultPayload struct {
	ID                     int64
	ProtocolHoneypotIpPort string
}

// attack event

type AttackType string

const (
	ProtocolAttackEvent    AttackType = "PROTOCOL_ATTACK_EVENT"
	TransparentAttackEvent AttackType = "TRANSPARENT_ATTACK_EVENT"
)

type TransparentEvent struct {
	AttackType               AttackType `json:"AttackType"`               //事件类型
	AgentID                  string     `json:"AgentID"`                  //agentID
	AttackIP                 string     `json:"AttackIP"`                 //攻击IP
	AttackPort               int32      `json:"AttackPort"`               //攻击Port
	ProxyIP                  string     `json:"ProxyIP"`                  //透明代理IP
	ProxyPort                int32      `json:"ProxyPort"`                //透明代理端口
	Transparent2ProtocolPort int32      `json:"Transparent2ProtocolPort"` //透明代理转发到协议代理的内部端口
	DestIP                   string     `json:"DestIP"`                   //目标IP
	DestPort                 int32      `json:"DestPort"`                 //目标端口
	EventTime                string     `json:"EventTime"`                //攻击事件发生时间
}

// redis task

type TaskPayload struct {
	TaskID       string       `json:"TaskID"`
	Status       TaskStatus   `json:"Status"`
	AgentID      string       `json:"AgentID"`
	TaskType     TaskType     `json:"TaskType"`
	OperatorType OperatorType `json:"OperatorType"`
}

type FileTaskPayload struct {
	TaskPayload
	FileMD5           string            `json:"FileMD5"`
	CommandParameters map[string]string `json:"CommandParameters"`
	URL               string            `json:"URL"`
	ScriptName        string            `json:"ScriptName"`
}

type TokenFileTaskPayload struct {
	FileTaskPayload
	TokenType string `json:"TokenType"`
}

type BaitFileTaskPayload struct {
	FileTaskPayload
	BaitType string `json:"BaitType"`
}

type HistoryBaitDeployTaskPayload struct {
	TaskPayload
	BaitType string `json:"BaitType"`
	BaitData string `json:"BaitData"`
}

type ProtocolProxyTaskPayload struct {
	TaskPayload
	ProxyPort    int32  `json:"ProxyPort"`
	HoneypotPort int32  `json:"HoneypotPort"`
	HoneypotIP   string `json:"HoneypotIP"`
	DeployPath   string `json:"DeployPath"`
}

type TransparentProxyTaskPayload struct {
	TaskPayload
	ProxyPort          int32  `json:"ProxyPort"`
	HoneypotServerPort int32  `json:"HoneypotServerPort"`
	HoneypotServerIP   string `json:"HoneypotServerIP"`
	ProbeIP            string `json:"ProbeIP"`
}

// TODO this is a bad design, put two contents into one struct、should use above two struct
// anyway、 design for agent

type ProxyTaskPayload struct {
	TaskPayload
	ProxyPort          int32  `json:"ProxyPort"`
	HoneypotPort       int32  `json:"HoneypotPort"`
	HoneypotIP         string `json:"HoneypotIP"`
	DeployPath         string `json:"DeployPath"`
	HoneypotServerPort int32  `json:"HoneypotServerPort"`
	HoneypotServerIP   string `json:"HoneypotServerIP"`
	ProbeIP            string `json:"ProbeIP"`
}

type SelectResultPayload struct {
	Count int64
	List  interface{}
}

type TaskType string
type OperatorType string
type TaskStatus int

const (
	PROTOCOL         TaskType = "PROTOCOL"
	TOKEN            TaskType = "TOKEN"
	BAIT             TaskType = "BAIT"
	ProtocolProxy    TaskType = "PROTOCOL_PROXY"
	TransparentProxy TaskType = "TRANSPARENT_PROXY"
	Heartbeat        TaskType = "HEARTBEAT"
)

const (
	DEPLOY   OperatorType = "DEPLOY"
	WITHDRAW OperatorType = "WITHDRAW"
)

const (
	IDLE    TaskStatus = -1 //初始
	RUNNING TaskStatus = 1  //下发中
	FAILED  TaskStatus = 2  //异常
	SUCCESS TaskStatus = 3  //成功
)

var BaitType = []string{"FILE", "HISTORY"}
var TokenType = []string{"FILE", "BrowserPDF", "EXE", "WPS"}

type TopologyNode struct {
	Id       string `json:"id"`
	Ip       string `json:"ip"`
	HostName string `json:"hostName"`
	NodeType string `json:"nodeType"`
}

type TopologyLine struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Status string `json:"status"`
}

type AttackLog struct {
	attackIp   string
	edgeIp     string
	honeyIp    string
	attackTime int32
	relayIp    string
	relayPort  int32
}

type TraceSourceResultPayload struct {
	ID           int64  `json:"ID"`           //ID
	Type         string `json:"Type"`         //攻击类型
	Time         string `json:"Time"`         //攻击时间
	HoneypotIP   string `json:"HoneypotIP"`   //蜜罐IP
	AttackIP     string `json:"AttackIP"`     //攻击IP
	ProtocolType string `json:"ProtocolType"` //协议类型
	Log          string `json:"Log"`          //日志
	Detail       string `json:"Detail"`       //详情
}
