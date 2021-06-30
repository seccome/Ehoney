package routers

import (
	"decept-defense/controllers"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/baitPolicyAdd", &controllers.BaitPolicy{},"post:Add")
	//beego.Router("/baitPolicyDelete", &controllers.BaitPolicy{},"post:Delete")
	//beego.Router("/baitPolicySelectAgentId", &controllers.BaitPolicy{},"post:Select")
	//beego.Router("/transPolicyAdd", &controllers.TransparentTransponderPolicy{},"post:Add")
	//beego.Router("/transPolicyDelete", &controllers.TransparentTransponderPolicy{},"post:Delete")
	//beego.Router("/transPolicySelectAgentId", &controllers.TransparentTransponderPolicy{},"post:SelectAgentId")
	//beego.Router("/attackLogList", &controllers.AttackLog{},"post:SelectList")
	//beego.Router("/attackLogDetail", &controllers.AttackLog{},"post:SelectDetail")
	//beego.Router("/addConfig", &controllers.Config{},"post:AddConf")
	//beego.Router("/selectConfig", &controllers.Config{},"post:SelectConf")
	//beego.Router("/deleteConfig", &controllers.Config{},"post:DeleteConf")
	beego.Router("/map", &controllers.WebSocket{}, "get:Map")
}
