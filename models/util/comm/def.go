package comm

// 常量定义
const SuccessCode = 0
const ErrorCode = 1001
const LoginoutCode = 401
const DBConnectError = "数据库连接失败"
const DBSelectError = "数据库查询失败"
const DataUnmarshalError = "数据解析失败"
const DataSelectSuccess = "成功"
const DBInsertSuccess = "数据写入成功"
const BodyNullMsg = "请求参数为空"
const BodyUnmarshalEorrMsg = "请求解析失败"
const NetWorkSuccess = "网络连接成功"
const NetWorkfail = "网络异常"
const PodNotFound = "容器记录不存在"
const PortUseError = "建议使用1025-65535范围内的端口号"
const DirUseError = "目录存在恶意字符，请确认！"
const OfflineStatus = 2
const CreateStatus = 1
const ExecFailStatus = -1
const AgentOutTimeStatus = -2
const AgentTypeRelay = "RELAY"      //协议转发类型
const AgentTypeUNRelay = "UN_RELAY" //协议转发撤回类型
const AgentTypeEdge = "EDGE"        //透明转发类型
const AgentTypeUNEdge = "UN_EDGE"   //透明转发撤回类型
const BaitHis = "Bait_HIS"          //history诱饵
const BaitFile = "Bait_FILE"        //文件诱饵
const BaitUNFile = "Bait_UN_FILE"   //撤回文件诱饵
const SignFile = "Sign_FILE"        //文件密签
const SignUNFile = "Sign_UN_FILE"   //撤回文件密签
const Creator = "admin"

type Admin struct {
	UserName string
	Password string
}

type AuthUser struct {
	Code     string `json:"code"`
	Mail     string `json:"mail"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Alias    string `json:"alias"`
	Username string `json:"id"`
	Job      string `json:"job"`
	WorkNum  string `json:"workNum"`
}

type PodImageResult struct {
	PodName         string
	PodImageAddress string
	PageNum         int
	PageSize        int
}

type PodImage struct {
	PodImageInfo []PodName
}

type PodName struct {
	ArtifactCount int    `json:"artifact_count"`
	CreationTime  string `json:"creation_time"`
	Id            int    `json:"id":4`
	Name          string `json:"name"`
	ProjectId     int    `json:"project_id"`
	UpdateTime    string `json:"update_time"`
}

type PodImageInfo struct {
	Digest              string        `json:"digest"`
	Icon                string        `json:"icon"`
	Id                  int           `json:"id"`
	Labels              string        `json:"labels"`
	Manifest_media_type string        `json:"manifest_media_type"`
	Media_type          string        `json:"media_type"`
	Project_id          int           `json:"project_id"`
	Pull_time           string        `json:"pull_time"`
	Push_time           string        `json:"push_time"`
	References          string        `json:"references"`
	Repository_id       int           `json:"repository_id"`
	Size                int           `json:"size"`
	Tags                []PodImageTag `json:"tags"`
	Type                string        `json:"type"`
}

type PodImageTag struct {
	Artifact_id   int    `json:"artifact_id"`
	Id            int    `json:"id"`
	Immutable     bool   `json:"immutable"`
	Name          string `json:"name"`
	Pull_time     string `json:"pull_time"`
	Push_time     string `json:"push_time"`
	Repository_id int    `json:"repository_id"`
	Signed        bool   `json:"signed"`
}

// 响应包结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 查询密签策略传入的json
type SelectBaitPolicyJson struct {
	AgentId          string `json:"agentId"`
	PageSize         int    `json:"pageSize"`
	PageNum          int    `json:"pageNum"`
	BaitId           string `json:"baitId"`
	BaitInfo         string `json:"baitInfo"`
	Creator          string `json:"creator"`
	CreateStartTime  string `json:"createStartTime"`
	CreateEndTime    string `json:"createEndTime"`
	OfflineStartTime string `json:"offlineStartTime"`
	OfflineEndTime   string `json:"offlineEndTime"`
	Status           int    `json:"status"`
}

// 查询诱饵策略传入的json
type SelectSignPolicyJson struct {
	AgentId          string `json:"agentId"`
	PageSize         int    `json:"pageSize"`
	PageNum          int    `json:"pageNum"`
	SignId           string `json:"signid"`
	SignType         string `json:"signtype"`
	Address          string `json:"address"`
	SignInfo         string `json:"signinfo"`
	Creator          string `json:"creator"`
	CreateStartTime  string `json:"createStartTime"`
	CreateEndTime    string `json:"createEndTime"`
	OfflineStartTime string `json:"offlineStartTime"`
	OfflineEndTime   string `json:"offlineEndTime"`
	Status           int    `json:"status"`
}

// 查询透明转发策略传入的json
type SelectTransPolicyJson struct {
	AgentId          string `json:"agentId"`
	PageSize         int    `json:"pageSize"`
	PageNum          int    `json:"pageNum"`
	Creator          string `json:"creator"`
	CreateStartTime  string `json:"createStartTime"`
	CreateEndTime    string `json:"createEndTime"`
	OfflineStartTime string `json:"offlineStartTime"`
	OfflineEndTime   string `json:"offlineEndTime"`
	ForwardPort      int    `json:"forwardPort"`
	HoneyTypeId      string `json:"honeyTypeId"`
	HoneyPotPort     int    `json:"honeyPotPort"`
	Status           int    `json:"status"`
	ServerIP         string `json:"serverIp"`
}

// 查询服务协议转发策略传入的json
type SelectHoneypotTransPolicyJson struct {
	AgentId          string
	PageSize         int
	PageNum          int
	Creator          string
	CreateStartTime  string
	CreateEndTime    string
	OfflineStartTime string
	OfflineEndTime   string
	ForwardPort      int
	HoneyPotId       string
	HoneyTypeId      string
	HoneyPotPort     int
	HoneyPotIp       string
	Status           int
}

// 前端传入查询蜜罐列表的json
type HoneyPotsType struct {
	//ServerId    string
	HoneyTypeId string
}

// 蜜罐类型
type HoneyPotType struct {
	HoneyPotType string `json:"honeyPotType"`
	TypeId       string `json:"honeyTypeId"`
}

// 诱饵类型
type BaitType struct {
	BaitId   string `json:"baitId"`
	BaitType string `json:"baitType"`
}

// 密签类型
type SignType struct {
	SignId   string `json:"signid"`
	SignType string `json:"signtype"`
}

//系统类型
type SysType struct {
	SysId   string `json:"sysId"`
	SysType string `json:"sysType"`
}

type ApplicationsClusters struct {
	ServerName string
	ServerIp   string
	ServerId   string
	Status     int
	Agentid    string
	VpcName    string
	PageSize   int
	PageNum    int
}

// 下发密签策略传入的json
type SignJson struct {
	TaskId    string
	AgentId   string
	SignId    string
	SignName  string
	Creator   string
	Status    int
	Data      string
	SignPath  string
	Address   string
	Md5       string
	Type      string
	Signinfo  *BaitInfo
	StartTime string
	EndTime   string
	PageSize  int
	PageNum   int
}

// 下发诱饵策略传入的json
type BaitJson struct {
	TaskId  string
	AgentId string
	BaitId  string
	//BaitType string
	Creator  string
	Status   int
	Data     string
	BaitPath string
	Address  string
	Md5      string
	Type     string
	Baitinfo *BaitInfo
}

// 诱饵部署详情
type BaitInfo struct {
}

// 蜜罐镜像
type HoneyImage struct {
	ID         int
	ImagesName string
	ImageType  string
	ImageOS    string
	ImagePort  int
	ImagesId   string
	PageSize   int
	PageNum    int
}

// 诱饵策略（下发到Redis）
type BaitPolicyJson struct {
	TaskId  string
	AgentId string
	Data    string
	Md5     string
	Status  int
	Type    string
}

type ProtocolTypePolicyJson struct {
	TaskId  string
	AgentId string
	Data    string
	Md5     string
	Status  int
	Type    string
}

// 上线透明转发策略传入的json
type TransparentTransponderJson struct {
	HoneyPotId   string
	AgentId      string
	ListenPort   int
	HoneyTypeId  string
	ForwardPort  int
	ServerId     string
	HoneyPotPort int
	Creator      string
	Status       int
	Path         string
	TaskId       string
}

// 下线透明转发策略的json
type TransOfflineJson struct {
	AgentId string `json:"agentId"`
	TaskId  string `json:"taskId"`
	Status  int    `json:"status"`
}

type TransTestJson struct {
	AgentId string `json:"agentId"`
	TaskId  string `json:"taskId"`
	Status  int    `json:"status"`
}

// 蜜罐端口转发
type HoneyTrans struct {
	ListenPort   int
	ServerTypeid string
	HoneyPotId   string
	SecCenter    string
}

// 钉钉/短信告警，传入json
type ConfJson struct {
	Id        int
	ConfName  string
	ConfValue string
	PageSize  int
	PageNum   int
}

// 透明转发策略（下发到Redis）
type TransparentTransponderPolicyJson struct {
	TaskId     string
	AgentId    string
	ListenPort int
	ServerType string
	HoneyIP    string
	HoneyPort  int
	Status     int
	Type       string
	Path       string
	SecCenter  string
}

// 攻击日志列表-前端json
type AttackLogListJson struct {
	ServerName  string
	HoneyTypeId string
	SrcHost     string
	AttackIP    string
	HoneyIP     string
	StartTime   string
	EndTime     string
	PageSize    int
	PageNum     int
}

// 攻击日志详情 - 前端json
type AttackLogDetailJson struct {
	Id           int
	SrcHost      string
	SrcPort      int
	HoneyPotPort int
	HoneyTypeId  string
	AttackIP     string
	HoneyIP      string
	StartTime    string
	EndTime      string
	EventDetail  string
	PageSize     int
	PageNum      int
}

// mysql代理log data结构体
type MysqlLogData struct {
	PassWord string
	UserName string
	Sql      string
	SqlType  string
}

type TelnetLogData struct {
	UserName string
	PassWord string
	Command  string
	Session  string
}

// ssh代理log data结构体
type SSHLogData struct {
	PassWord string
	UserName string
	Command  string
	Session  string
}

// redis代理log data结构体
type RedisLogData struct {
	PassWord string
	UserName string
	Command  string
	Session  string
}

// http代理log data结构体
type HttpLogData struct {
	PassWord string
	UserName string
	Command  string
	Session  string
}

type HoneyCluster struct {
	ClusterName  string
	ClusterIp    string
	ClusterStats int
	ClusterId    string
	PageSize     int
	PageNum      int
}

type HoneyPots struct {
	HoneyName   string
	HoneyTypeId string
	HoneypotId  string
	HoneyIp     string
	StartTime   string
	EndTime     string
	Creator     string
	ServerId    string
	SysId       string
	Status      int
	PageNum     int
	PageSize    int
}

type HoneySigns struct {
	Taskid      string
	SignId      string
	SignInfo    string
	SignType    string
	Address     string
	CreateTime  string
	OfflineTime string
	StartTime   string
	EndTime     string
	PageSize    int
	PageNum     int
	Creator     string
	HoneypotId  string
	Status      int
}

type HoneyBaits struct {
	Taskid       string
	Address      string
	BaitName     string
	BaitPath     string
	Agentid      string
	BaitId       string
	BaitType     string
	BaitInfo     string
	HoneypotId   string
	HoneypotName string
	StartTime    string
	EndTime      string
	CreateTime   string
	OfflineTime  string
	Creator      string
	Status       int
	PageSize     int
	PageNum      int
}

type Servers struct {
	ServerName  string
	ServerIp    string
	ServerStats int
	ServerId    string
	AgentId     string
	PageSize    int
	PageNum     int
}

type Signs struct {
	SignName    string
	Creator     string
	Createtime  string
	SignId      string
	SignSysType string
	PageSize    int
	PageNum     int
}

type Baits struct {
	BaitType    string
	BaitName    string
	Sysid       string
	CreateTime  string
	Creator     string
	BaitId      string
	BaitSysType string
	StartTime   string
	EndTime     string
	PageSize    int
	PageNum     int
}

// falco output fields结构体
type OutputFields struct {
	ContainerId string `json:"container.id"` // 容器ID
	FileName    string `json:"fd.name"`      // 命令中的文件名
	Command     string `json:"proc.cmdline"` // 完整的命令，例如touch 1.txt
	CmdName     string `json:"proc.name"`    // 使用的系统命令
	Pname       string `json:"proc.pname"`   // 例如bash
	UserName    string `json:"user.name"`    // 登录的用户名，执行命令的用户
	K8sPodName  string `json:"k8s.pod.name"` // pod name
}

type VpsApplications struct {
	PageSize   int      `json:"pageSize"`
	TotalPage  int      `json:"totalPage"`
	TotalCount int      `json:"totalCount"`
	Data       []VpsApp `json:"data"`
}

type VpsApp struct {
	EcsId    string    `json:"ecsId"`
	EcsName  string    `json:"ecsName"`
	VpcIp    string    `json:"vpcIp"`
	Status   VpsStatus `json:"status"`
	Vpc      Vpc       `json:"vpc"`
	VpsOwner VpsOwner  `json:"owner"`
}

type VpsStatus struct {
	StatusName string `json:"statusName"`
}

type Vpc struct {
	VpcName string `json:"vpcName"`
}

type VpsOwner struct {
	FirstName string `json:"first_name"`
}

type VpcAgent struct {
	Data AgentList `json:"data"`
}

type AgentList struct {
	AllCount int             `json:"AllCount"`
	List     []AgentListData `json:"list"`
}

type AgentListData struct {
	AgentID  string `json:"AgentID"`
	Hostname string `json:"Hostname"`
}

type ClamavData struct {
	Virname         string
	Filename        string
	HoneyPotip      string
	CreateStartTime string
	CreateEndTime   string
	PageSize        int
	PageNum         int
}

type ProtocolType struct {
	TypeId       string
	ProtocolName string
	PageSize     int
	PageNum      int
}

type HarborInfo struct {
	HarborId string
}
