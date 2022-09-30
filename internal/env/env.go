package env

import (
	"decept-defense/pkg/util"
	"fmt"
	"go.uber.org/zap"
)

func Setup() {
	coverProjectK3sConfig()
	chmodDirs()
}

func coverProjectK3sConfig() {
	err := util.CoverProjectK3sConfig(util.ProjectDir())
	if err != nil {
		zap.L().Error("环境配置设置失败")
		zap.L().Error(err.Error())
	}
}

func chmodDirs() {
	tokensDir := fmt.Sprintf("%s/tool/file_token_trace/linux", util.ProjectDir())
	err := util.ChmodDir(tokensDir)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Dir %s 授权失败", tokensDir))
		zap.L().Error(err.Error())
	}
}
