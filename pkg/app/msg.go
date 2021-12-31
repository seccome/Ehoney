package app

var MsgFlags = map[int]string{
	SUCCESS:       "ok",
	INTERNAlERROR: "内部异常",
	InvalidParams: "请求参数错误",
	ErrorAuth:     "鉴权错误",

	//OK Message
	OKHeartBeat:       "心跳正常",
	ErrorFileUpload:   "文件上传失败",
	ErrorDirCreate:    "目录创建失败",
	ErrorFileCompress: "文件压缩失败",
	ErrorUUID:         "UUID生成失败",
	ErrorDatabase:     "数据库异常",
	ErrorRedis:        "Redis异常",
	ErrorConnectTest:  "网络探测异常",

	//INTERNAlERROR Message
	ErrorPasswordCheck:         "账号或是密码错误",
	ErrorAuthCheckTokenTimeout: "Token已超时",
	ErrorAuthToken:             "Token生成失败",
	ErrorPasswdHash:            "密码加密失败",
	ErrorCreateUser:            "用户创建失败",
	ErrorAuthCheckTokenFail:    "Token鉴权错误",
	ErrorDuplicateUser:         "用户名称重复",

	ErrorCreateToken:            "密签创建错误",
	ErrorDeleteToken:            "密签删除错误",
	ErrorGetToken:               "密签查找失败",
	ErrorDeleteBait:             "诱饵文件删除失败",
	ErrorHoneypotBaitCreate:     "蜜罐诱饵创建异常",
	ErrorHoneypotBaitDeploy:     "蜜罐诱饵部署失败",
	ErrorHoneypotBaitDelete:     "蜜罐诱饵删除异常",
	ErrorHoneypotBaitWithdraw:   "蜜罐诱饵撤销失败",
	ErrorDuplicateBaitName:      "诱饵名称重复",
	ErrorDuplicateTokenName:     "密签名称重复",
	ErrorBaitNotExist:           "诱饵不存在",
	ErrorHoneypotK8SCP:          "K8S拷贝异常",
	ErrorHoneypotHistoryBait:    "暂不支持下发history类型诱饵到蜜罐",
	ErrorProbeBaitCreate:        "探针诱饵创建失败",
	ErrorTokenNotExist:          "密签不存在",
	ErrorHoneypotTokenCreate:    "蜜罐密签创建失败",
	ErrorFileTokenType:          "不支持的文件密签类型、仅支持pdf、docx、xlsx、pptx",
	ErrorDoFileTokenTrace:       "文件加密签失败",
	ErrorDoBrowserPDFTokenTrace: "浏览器PDF密签创建失败",

	ErrorCreateVirusRecord: "木马文件记录创建失败",
	ErrorSelectVirusRecord: "木马记录查询失败",

	ErrorHoneypotServerNotExist:     "蜜罐服务器异常、请检测蜜罐服务器状态",
	ErrorHoneypotCreate:             "K8S创建镜像失败",
	ErrorHoneypotIpAddress:          "蜜罐IP获取错误",
	ErrorHoneypotRecordCreate:       "蜜罐数据插入失败",
	ErrorHoneypotDelete:             "K8S删除蜜罐错误",
	ErrorHoneypotNameExist:          "蜜罐名称存在",
	ErrorProbeNotExist:              "探针服务器不在线",
	ErrorHoneypotPodExist:           "容器名称存在请先删除容器或是修改容器名字",
	ErrorHoneypotNotExist:           "蜜罐不存在",
	ErrorProbeServerNotExist:        "探针服务器异常、请检测探针服务器状态",
	ErrorHoneypotProtocolProxyExist: "存在关联此蜜罐的协议代理、请先删除协议代理",

	ErrorProtocolCreate: "协议创建失败",
	ErrorProtocolGet:    "协议获取失败",
	ErrorProtocolDel:    "协议删除失败",
	ErrorProtocolDup:    "协议名称重复",

	ErrorProtocolProxyCreate:         "协议转发创建失败",
	ErrorProxyPortDup:                "代理端口重复",
	ErrorProtocolProxyFail:           "协议代理下发失败",
	ErrorProtocolProxyOfflineFail:    "协议代理离线失败",
	ErrorProtocolProxyOnlineFail:     "协议代理上线失败",
	ErrorTransparentProxyFail:        "透明代理下发失败",
	ErrorTransparentProxyOfflineFail: "透明代理离线失败, 请确认探针是否在线",
	ErrorTransparentProxyOnlineFail:  "透明代理上线失败, 请确认探针是否在线",
	ErrorProtocolNotExist:            "对应协议不存在",
	ErrorProtocolProxyNotExist:       "协议代理不存在",
	ErrorTransparentProxyCreate:      "透明转发创建失败",
	ErrorTransparentProxyNotExist:    "透明代理不存在",
	ErrorProtocolUpdate:              "协议更新异常",
	ErrorProtocolPortRange:           "端口范围异常",
	ErrorTransparentProxyPortRange:   "透明代理端口不在设置端口范围、请先修改端口范围",
	ErrorImagesNotConfig:             "镜像的端口或是服务未配置、请先进行配置",
	ErrorProtocolOffline:             "协议代理已离线",
	ErrorProtocolOnline:              "协议代理已上线",
	ErrorTransparentOffline:          "透明代理已离线",
	ErrorTransparentOnline:           "透明代理已上线",
	EngineVersionSame:                "Engine版本已存在",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[INTERNAlERROR]
}
