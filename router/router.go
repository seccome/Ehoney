package router

import (
	"decept-defense/controllers/agent_bait_handler"
	"decept-defense/controllers/agent_handler"
	"decept-defense/controllers/attack_handler"
	"decept-defense/controllers/bait_handler"
	"decept-defense/controllers/extranet_handler"
	"decept-defense/controllers/harbor_handler"
	"decept-defense/controllers/heartbeat_handler"
	"decept-defense/controllers/honeypot_bait_handler"
	"decept-defense/controllers/honeypot_handler"
	"decept-defense/controllers/images_handler"
	"decept-defense/controllers/protocol_handler"
	"decept-defense/controllers/protocol_proxy_handler"
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
	"html/template"
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
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})
	r.LoadHTMLGlob("front/*.html")
	r.Static("/decept-defense/static", "front/static")
	r.GET("/decept-defense/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/decept-defense/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/decept-defense/threaten-perception/*id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/decept-defense/honeypots/*id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/decept-defense/probes-list/*id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/decept-defense/proxy-manage/*id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/decept-defense/trap-manage/*id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/decept-defense/system-config", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/decept-defense/datav", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api := r.Group("/api")
	{
		//create counter info
		api.GET("/info", attack_handler.CreateCountEvent)
		public := api.Group("/public")
		{
			public.GET("/health", heartbeat_handler.Heartbeat)
			public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
			public.POST("/login", user_handler.Login)
			public.POST("/signup", user_handler.SignUp)
			public.POST("/attack/protocol", attack_handler.CreateProtocolAttackEvent)
			public.POST("/attack/transparent", attack_handler.CreateTransparentEventEvent)
			public.POST("/attack/falco", attack_handler.CreateFalcoAttackEvent)
			public.GET("/topology/map", topology_handler.TopologyMapHandle)
			public.GET("/extranet", extranet_handler.GetExtranetConfig)
			public.GET("/engine/linux/version", agent_handler.QueryLinuxEngineVersion)
			public.GET("/engine/linux", agent_handler.DownloadLinuxEngine)
			public.POST("/engine/linux/upload", agent_handler.UploadLinuxEngine)
			public.GET("/transparentList", trans_proxy_handler.GetTransparentByAgent)
			public.PUT("/webhook", webhook_handler.UpdateWebHookConfig)
			public.POST("/agent/heartbeat", agent_handler.AgentHeartBeat)
			public.GET("/agent/task/:id", agent_handler.LoadAgentTask)
			public.POST("/agent/task/callback", agent_handler.AgentTaskCallBack)
			public.GET("/token/trace-alert", token_trace_handler.TraceMsgReceive)
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
			//get bait
			private.POST("/bait/set", bait_handler.GetBait)
			private.PUT("/bait/token/:id", bait_handler.SignTokenForFileBait)
			//get bait type

			private.GET("/bait/type", bait_handler.GetBaitByType)
			//add bait
			private.POST("/bait", bait_handler.CreateBait)
			//delete bait by ID
			private.DELETE("/bait/:id", bait_handler.DeleteBaitByID)
			//download bait by ID
			private.GET("/bait/:id", bait_handler.DownloadBaitByID)
			//create honeypot bait
			private.POST("/honeypot/bait-task", honeypot_bait_handler.CreateHoneypotBait)
			//delete honeypot bait
			private.DELETE("/honeypot/bait-task/:id", honeypot_bait_handler.DeleteHoneypotBaitByID)

			//download honeypot bait by ID
			private.GET("/honeypot/bait-task/:id", honeypot_bait_handler.DownloadHoneyBaitByID)
			//create probe bait
			private.POST("/agent/bait-task", agent_bait_handler.CreateProbeBait)

			//get bait
			private.POST("/honeypot/bait-task/set", honeypot_bait_handler.GetHoneypotBait)

			//get bait
			private.POST("/agent/bait-task/set", agent_bait_handler.GetProbeBait)

			//delete probe bait
			private.DELETE("/bait/probe/:id", agent_bait_handler.DeleteProbeBaitByID)
			//download probe bait by ID
			private.GET("/bait/probe/:id", agent_bait_handler.DownloadProbeBaitByID)
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
			private.PUT("/proxy/transparent/:id", trans_proxy_handler.UpdateTransparentProxyStatus)
			//get image list
			private.POST("/images/set", images_handler.GetImages)
			//change image list
			private.PUT("/images/:id", images_handler.UpdateImage)
			//change image list
			private.POST("/images", images_handler.CreateImage)
			//get pod image
			private.GET("/images/pod", images_handler.GetPodImages)
			//get agent record
			private.POST("/agent/set", agent_handler.AgentPage)
			//download linux agent
			private.GET("/agent/linux", agent_handler.DownloadLinuxAgent)
			private.POST("/agent/transparent/set", trans_proxy_handler.GetTransparentProxy)

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
