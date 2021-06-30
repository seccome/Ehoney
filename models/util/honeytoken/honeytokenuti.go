package honeytoken

import (
	"github.com/astaxie/beego/logs"
	"os/exec"
)

func DoFileSignTrace(orginalfilename string, newfilename string, uploaddir string, outputdir string, tracecode string) error {
	//log.Println("sign trace cmd :", "python", "tool/TraceFile.py", "-i", uploaddir+"/"+orginalfilename, "-o", outputdir+"/"+newfilename, "-b", tracecode)
	cmd := exec.Command("tool/TraceFile", "-i", uploaddir+"/"+orginalfilename, "-o", outputdir+"/"+newfilename, "-b", tracecode)
	out, err := cmd.CombinedOutput()
	logs.Error("cmd out:", string(out))
	if err != nil {
		logs.Error("[文件跟踪]: 文件%s操作,执行命令报错%s", orginalfilename, string(out))
		return err
	}
	return nil
}
