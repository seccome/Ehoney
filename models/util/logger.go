package util

import (
	"github.com/astaxie/beego/logs"
	"log"
)

func InitLogger() (err error) {
	logs.EnableFuncCallDepth(true)
	_ = logs.SetLogger("console")
	//_ = logs.SetLogger(logs.AdapterMultiFile,`{"filename":"./logs/app.log","level":6,"maxlines":1000,"separate":["error"]}`)
	error := logs.SetLogger(logs.AdapterFile, `{"filename":"logs/error.log","level":6,"maxlines":1000}`)
	if error != nil {
		log.Println("[ERROR]init Logger error:", error)
	}
	//logs.Async()
	logs.SetLogFuncCall(true)
	return
}
