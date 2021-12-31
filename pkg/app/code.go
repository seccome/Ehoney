package app

const (
	SUCCESS       = 200
	INTERNAlERROR = 500
	InvalidParams = 400
	ErrorAuth     = 401

	OKHeartBeat       = 1001
	ErrorFileUpload   = 1002
	ErrorDirCreate    = 1003
	ErrorFileCompress = 1004
	ErrorUUID         = 1005
	ErrorDatabase     = 1006
	ErrorRedis        = 1007
	ErrorConnectTest  = 1008

	ErrorPasswordCheck         = 2001
	ErrorAuthCheckTokenTimeout = 2002
	ErrorAuthToken             = 2003
	ErrorPasswdHash            = 2005
	ErrorCreateUser            = 2006
	ErrorAuthCheckTokenFail    = 2007
	ErrorDuplicateUser         = 2008

	ErrorCreateToken            = 3001
	ErrorDeleteToken            = 3002
	ErrorGetToken               = 3003
	ErrorDeleteBait             = 3004
	ErrorHoneypotBaitCreate     = 3005
	ErrorHoneypotBaitDeploy     = 3006
	ErrorHoneypotBaitDelete     = 3007
	ErrorHoneypotBaitWithdraw   = 3008
	ErrorDuplicateBaitName      = 3009
	ErrorDuplicateTokenName     = 3010
	ErrorBaitNotExist           = 3011
	ErrorHoneypotK8SCP          = 3012
	ErrorHoneypotHistoryBait    = 3013
	ErrorProbeBaitCreate        = 3014
	ErrorTokenNotExist          = 3015
	ErrorHoneypotTokenCreate    = 3016
	ErrorFileTokenType          = 3017
	ErrorDoFileTokenTrace       = 3018
	ErrorDoBrowserPDFTokenTrace = 3019

	ErrorCreateVirusRecord = 4001
	ErrorSelectVirusRecord = 4002

	ErrorHoneypotServerNotExist     = 5001
	ErrorHoneypotCreate             = 5002
	ErrorHoneypotIpAddress          = 5003
	ErrorHoneypotRecordCreate       = 5004
	ErrorHoneypotDelete             = 5005
	ErrorHoneypotNameExist          = 5006
	ErrorProbeNotExist              = 5007
	ErrorHoneypotPodExist           = 5008
	ErrorHoneypotNotExist           = 5009
	ErrorProbeServerNotExist        = 5010
	ErrorHoneypotProtocolProxyExist = 5011

	ErrorProtocolCreate              = 6001
	ErrorProtocolGet                 = 6002
	ErrorProtocolDel                 = 6003
	ErrorProtocolDup                 = 6004
	ErrorProtocolProxyCreate         = 6005
	ErrorProxyPortDup                = 6006
	ErrorProtocolProxyFail           = 6007
	ErrorProtocolProxyOfflineFail    = 6008
	ErrorProtocolProxyOnlineFail     = 6009
	ErrorTransparentProxyFail        = 6010
	ErrorTransparentProxyOfflineFail = 6011
	ErrorTransparentProxyOnlineFail  = 6012
	ErrorProtocolNotExist            = 6013
	ErrorProtocolProxyNotExist       = 6014
	ErrorTransparentProxyCreate      = 6015
	ErrorTransparentProxyNotExist    = 6016
	ErrorProtocolUpdate              = 6017
	ErrorProtocolPortRange           = 6018
	ErrorTransparentProxyPortRange   = 6019

	ErrorImagesNotConfig = 7001

	ErrorProtocolOffline    = 8001
	ErrorProtocolOnline     = 8002
	ErrorTransparentOffline = 8003
	ErrorTransparentOnline  = 8004
	EngineVersionSame       = 8005
)
