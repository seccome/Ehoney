package router

import (
	"decept-defense/controllers/agent_handler"
	"decept-defense/controllers/attack_handler"
	"decept-defense/controllers/bait_handler"
	"decept-defense/controllers/extranet_handler"
	"decept-defense/controllers/harbor_handler"
	"decept-defense/controllers/heartbeat_handler"
	"decept-defense/controllers/honeypot_bait_handler"
	"decept-defense/controllers/honeypot_handler"
	"decept-defense/controllers/honeypot_token_handler"
	"decept-defense/controllers/images_handler"
	"decept-defense/controllers/probe_bait_handler"
	"decept-defense/controllers/probe_handler"
	"decept-defense/controllers/probe_token_handler"
	"decept-defense/controllers/protocol_handler"
	"decept-defense/controllers/protocol_proxy_handler"
	"decept-defense/controllers/token_handler"
	"decept-defense/controllers/token_trace_handler"
	"decept-defense/controllers/topology_handler"
	_ "decept-defense/controllers/topology_handler"
	"decept-defense/controllers/trans_proxy_handler"
	"decept-defense/controllers/user_handler"
	"decept-defense/controllers/virus_handler"
	"decept-defense/controllers/webhook_handler"
	_ "decept-defense/docs"
	"decept-defense/middleware/cors"
	"decept-defense/middleware/jwt"
	"decept-defense/pkg/configs"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func MakeRoute() *gin.Engine {
	r := gin.New()
	r.MaxMultipartMemory = 32 << 20
	//enable Cors
	if configs.GetSetting().App.EnableCors {
		r.Use(cors.Cors())
	}
	//enable recover middleware
	r.Use(gin.Recovery())

	r.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(zap.L(), true))

	r.StaticFS("/upload/", http.Dir("upload"))
	r.StaticFS("/agent/", http.Dir("agent"))

	api := r.Group("/api/")
	{
		//create counter info
		api.GET("/info", attack_handler.CreateCountEvent)

		public := api.Group("/public")
		{
			public.GET("/health", heartbeat_handler.Heartbeat)
			public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
			public.POST("/login", user_handler.Login)
			public.POST("/signup", user_handler.SignUp)
			//create attack event
			public.POST("/attack/protocol", attack_handler.CreateProtocolAttackEvent)
			public.POST("/attack/falco", attack_handler.CreateFalcoAttackEvent)

			//insert ssh key
			public.POST("/protocol/key", protocol_handler.CreateSSHKey)
			public.GET("/topology/map", topology_handler.TopologyMapHandle)

			public.GET("/extranet", extranet_handler.GetExtranetConfig)

			public.GET("/engine/linux/version", agent_handler.QueryLinuxEngineVersion)
			public.GET("/engine/linux", agent_handler.DownloadLinuxEngine)
			public.POST("/engine/linux/upload", agent_handler.UploadLinuxEngine)
			public.GET("/transparentList", trans_proxy_handler.GetTransparentByAgent)
			public.PUT("/webhook", webhook_handler.UpdateWebHookConfig)

		}
		private := api.Group("/v1/")
		private.Use(jwt.JWT())
		{
			//use
			private.PUT("/user/password", user_handler.ChangePassword)

			//attack
			private.POST("/attack", attack_handler.GetAttackList)
			private.POST("/attack/falco", attack_handler.GetFalcoAttackList)
			private.GET("/attack/falco/:id", attack_handler.GetFalcoAttackDetail)
			private.POST("/attack/token", attack_handler.GetTokenTraceLog)
			private.POST("/attack/trace", attack_handler.GetAttackSource)

			// screen display
			private.GET("/attack/display/attackIP", attack_handler.GetAttackIPStatistics)
			private.GET("/attack/display/probeIP", attack_handler.GetProbeIPStatistics)
			private.GET("/attack/display/protocol", attack_handler.GetAttackProtocolStatistics)
			private.GET("/attack/display/location", attack_handler.GetAttackLocationStatistics)

			//get token
			private.POST("/token/set", token_handler.GetToken)
			//get token type
			private.GET("/token/type", token_handler.GetTokenByType)
			//add token
			private.POST("/token", token_handler.CreateToken)
			//delete token by ID
			private.DELETE("/token/:id", token_handler.DeleteTokenByID)
			//download token by ID
			private.GET("/token/:id", token_handler.DownloadTokenByID)
			//get token name list
			private.GET("/token/name/set", token_handler.GetTokenNameList)

			//create honeypot token
			private.POST("/token/honeypot", honeypot_token_handler.CreateHoneypotTokenNew)
			//delete honeypot token
			private.DELETE("/token/honeypot/:id", honeypot_token_handler.DeleteHoneypotTokenByID)
			//get honeypot token
			private.POST("/token/honeypot/set", honeypot_token_handler.GetHoneypotToken)
			//download honeypot token by ID
			private.GET("/token/honeypot/:id", honeypot_token_handler.DownloadHoneypotTokenByID)

			//create probe token
			private.POST("/token/probe", probe_token_handler.CreateProbeTokenNew)
			//delete probe token
			private.DELETE("/token/probe/:id", probe_token_handler.DeleteProbeTokenByID)
			//get probe token
			private.POST("/token/probe/set", probe_token_handler.GetProbeToken)
			//download probe token by ID
			private.GET("/token/probe/:id", probe_token_handler.DownloadProbeTokenByID)

			//get bait
			private.POST("/bait/set", bait_handler.GetBait)
			//get bait type
			private.GET("/bait/type", bait_handler.GetBaitByType)
			//add bait
			private.POST("/bait", bait_handler.CreateBait)
			//delete bait by ID
			private.DELETE("/bait/:id", bait_handler.DeleteBaitByID)
			//download bait by ID
			private.GET("/bait/:id", bait_handler.DownloadBaitByID)

			//create honeypot bait
			private.POST("/bait/honeypot", honeypot_bait_handler.CreateHoneypotBait)
			//delete honeypot bait
			private.DELETE("/bait/honeypot/:id", honeypot_bait_handler.DeleteHoneypotBaitByID)
			//get honeypot bait
			private.POST("/bait/honeypot/set", honeypot_bait_handler.GetHoneypotBait)
			//download honeypot bait by ID
			private.GET("/bait/honeypot/:id", honeypot_bait_handler.DownloadHoneyBaitByID)

			//create probe bait
			private.POST("/bait/probe", probe_bait_handler.CreateProbeBait)
			//delete probe bait
			private.DELETE("/bait/probe/:id", probe_bait_handler.DeleteProbeBaitByID)
			//get probe bait
			private.POST("/bait/probe/set", probe_bait_handler.GetProbeBait)
			//download probe bait by ID
			private.GET("/bait/probe/:id", probe_bait_handler.DownloadProbeBaitByID)

			//create virus record
			private.POST("/virus", virus_handler.CreateVirusRecord)
			//select virus record
			private.POST("/virus/set", virus_handler.SelectVirusRecord)

			//get honeypot record
			private.POST("/honeypot/set", honeypot_handler.GetHoneypots)
			//create honeypot
			private.POST("/honeypot", honeypot_handler.CreateHoneypot)
			//delete honeypot
			private.DELETE("/honeypot/:id", honeypot_handler.DeleteHoneypot)
			//get honeypot detail
			private.GET("/honeypot/:id", honeypot_handler.GetHoneypotDetail)
			//get protocol proxy honeypot
			private.GET("/honeypot/protocol", honeypot_handler.GetProtocolProxyHoneypots)

			//get probe record
			private.POST("/probe/set", probe_handler.GetProbes)

			//get protocol
			private.POST("/protocol/set", protocol_handler.GetProtocol)
			//create protocol
			private.POST("/protocol", protocol_handler.CreateProtocol)
			//delete protocol
			private.DELETE("/protocol/:id", protocol_handler.DeleteProtocol)
			//update  protocol port range
			private.POST("/protocol/port/:id", protocol_handler.UpdateProtocolPortRange)

			private.GET("/protocol/type", protocol_handler.GetProtocolType)

			//get protocol proxy
			private.POST("/proxy/protocol/set", protocol_proxy_handler.GetProtocolProxy)
			//create protocol proxy
			private.POST("/proxy/protocol", protocol_proxy_handler.CreateProtocolProxy)
			//delete protocol proxy
			private.DELETE("/proxy/protocol/:id", protocol_proxy_handler.DeleteProtocolProxy)
			//online protocol proxy
			private.POST("/proxy/protocol/online/:id", protocol_proxy_handler.OnlineProtocolProxy)
			//offline protocol proxy
			private.POST("/proxy/protocol/offline/:id", protocol_proxy_handler.OfflineProtocolProxy)
			//test protocol proxy
			private.GET("/proxy/protocol/test/:id", protocol_proxy_handler.TestProtocolProxy)

			//get transparent proxy
			private.POST("/proxy/transparent/set", trans_proxy_handler.GetTransparentProxy)
			//create transparent proxy
			private.POST("/proxy/transparent", trans_proxy_handler.CreateTransparentProxy)
			//delete transparent proxy
			private.DELETE("/proxy/transparent/:id", trans_proxy_handler.DeleteTransparentProxy)
			//online transparent proxy
			private.POST("/proxy/transparent/online/:id", trans_proxy_handler.OnlineTransparentProxy)

			//batch online transparent proxy
			private.POST("/proxy/transparent/online/batch", trans_proxy_handler.BatchOnlineTransparentProxy)

			//offline transparent proxy
			private.POST("/proxy/transparent/offline/:id", trans_proxy_handler.OfflineTransparentProxy)

			//batch offline transparent proxy
			private.POST("/proxy/transparent/offline/batch", trans_proxy_handler.BatchOfflineTransparentProxy)
			//test transparent proxy
			private.GET("/proxy/transparent/test/:id", trans_proxy_handler.TestTransparentProxy)

			private.PUT("/proxy/transparent/test/:id", trans_proxy_handler.UpdateTransparentProxyStatus)

			//get image list
			private.POST("/images/set", images_handler.GetImages)
			//change image list
			private.PUT("/images/:id", images_handler.UpdateImage)
			//get pod image
			private.GET("/images/pod", images_handler.GetPodImages)

			//download linux agent
			private.GET("/agent/linux", agent_handler.DownloadLinuxAgent)
			//download window agent
			private.GET("/agent/windows", agent_handler.DownloadWindowsAgent)

			//set harbor info
			private.PUT("/harbor", harbor_handler.UpdateHarborConfig)
			private.GET("/harbor", harbor_handler.GetHarborConfig)
			private.POST("/harbor/health", harbor_handler.TestHarborConnection)
			private.PUT("/webhook", webhook_handler.UpdateWebHookConfig)
			private.GET("/webhook", webhook_handler.GetWebHookConfig)
			private.GET("/token/trace", token_trace_handler.GetTraceHostConfig)
			private.PUT("/token/trace", token_trace_handler.UpdateTraceHostConfig)
			private.PUT("/extranet", extranet_handler.UpdateExtranetConfig)
			private.GET("/extranet", extranet_handler.GetExtranetConfig)
		}
	}
	return r
}
