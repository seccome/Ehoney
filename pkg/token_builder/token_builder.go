package token_builder

import (
	"decept-defense/models"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path"
	"strings"
)

type CreateTokenPayload struct {
	TokenType  string `json:"TokenType" form:"TokenType" binding:"required"`
	TokenName  string `json:"TokenName" form:"TokenName" binding:"required"`
	TokenData  string `json:"TokenData" form:"TokenData"`
	UploadPath string `json:"UploadPath" form:"UploadPath"`
	TraceUrl   string `json:"UploadPath" form:"UploadPath"`
	TraceCode  string `json:"TraceCode" form:"TraceCode"`
	LocalPath  string `json:"UploadPath" form:"LocalPath"`
}

func TokenBaitFile(bait models.Bait) (string, string, error) {

	zap.L().Info(fmt.Sprintf("开始注入文件[%s] 类型[%s]蜜签", bait.BaitType, bait.UploadPath))
	if bait.BaitType != "WPS" && !util.CheckPathIsExist(bait.UploadPath) {
		zap.L().Error("待加签文件不存在")
		return "", "", errors.New("source file is not exist")
	}

	var toolPath = path.Join(util.WorkingPath(), "tool", "file_token_trace", "linux", "TraceFile")

	if !util.CheckPathIsExist(toolPath) {
		zap.L().Error("加签工具不存在")
		return "", "", errors.New("trace file is not exist")
	}
	bait.TokenTraceUrl = strings.Join([]string{configs.GetSetting().App.Extranet, configs.GetSetting().App.TokenTraceApiPath}, "/") + "?tracecode=" + bait.BaitId
	bait.LocalPath = path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "bait", bait.BaitId, bait.FileName, "/")

	zap.L().Info(fmt.Sprintf("开始生产加签后文件夹[%s]", bait.LocalPath))
	zap.L().Info(fmt.Sprintf("开始生产加签后文件夹Dir:[%s]", path.Dir(bait.LocalPath)))
	zap.L().Info(fmt.Sprintf("TokenTraceUrl: [%s]", bait.TokenTraceUrl))

	err := util.CreateDir(path.Dir(bait.LocalPath) + "/")
	if err != nil {
		zap.L().Error("密签文件夹创建失败: " + bait.LocalPath)
		zap.L().Error(err.Error())
	}
	var cmd *exec.Cmd

	// 命令组装
	switch bait.BaitType {
	case "WPS":
		cmd = exec.Command(toolPath, "-u", bait.TokenTraceUrl, "-o", fmt.Sprintf("%s/%s.docx", bait.LocalPath, bait.BaitName), "-w", bait.BaitData, "-t", "wps")
	case "FILE":
		cmd = exec.Command(toolPath, "-u", bait.TokenTraceUrl, "-o", bait.LocalPath, "-i", bait.UploadPath, "-t", "office")
	case "EXE":
		cmd = exec.Command(toolPath, "-u", bait.TokenTraceUrl, "-o", bait.LocalPath, "-i", bait.UploadPath, "-t", "exe")
	default:
		zap.L().Error("无法处理的蜜签类型: " + bait.BaitType)
		return "", "", errors.New("无法处理的蜜签类型: " + bait.BaitType)
	}
	cmd.Dir = path.Dir(toolPath)
	zap.L().Info("cmd : " + cmd.String())
	_, err = cmd.CombinedOutput()
	if err != nil {
		zap.L().Error("文件密签加签失败")
		zap.L().Error(err.Error())
		fmt.Println("文件密签加签失败:" + err.Error())
		_ = os.RemoveAll(path.Dir(bait.LocalPath))
		return "", "", err
	}
	return bait.LocalPath, bait.TokenTraceUrl, nil
}
