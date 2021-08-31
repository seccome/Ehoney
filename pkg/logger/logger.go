package logger

import (
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
)

func SetUp()  {
	logPath := path.Join(util.WorkingPath(), "logs", "decept-defense.log")
	switch configs.GetSetting().Log.Level {
	case "info":
		SetLogs(zap.InfoLevel, configs.GetSetting().Log.LogFormatJson, logPath)
	case "debug":
		SetLogs(zap.DebugLevel, configs.GetSetting().Log.LogFormatJson, logPath)
	case "error":
		SetLogs(zap.ErrorLevel, configs.GetSetting().Log.LogFormatJson, logPath)
	case "fatal":
		SetLogs(zap.FatalLevel, configs.GetSetting().Log.LogFormatJson, logPath)
	case "warn":
		SetLogs(zap.WarnLevel, configs.GetSetting().Log.LogFormatJson, logPath)
	default:
		SetLogs(zap.DebugLevel, configs.GetSetting().Log.LogFormatJson, logPath)
	}
}

func SetLogs(logLevel zapcore.Level, logFormat, fileName string) {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        configs.GetSetting().Log.TimeKey,
		LevelKey:       configs.GetSetting().Log.LevelKey,
		NameKey:        configs.GetSetting().Log.NameKey,
		CallerKey:      configs.GetSetting().Log.CallerKey,
		MessageKey:     configs.GetSetting().Log.MessageKey,
		StacktraceKey:  configs.GetSetting().Log.StackTraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志输出格式
	var encoder zapcore.Encoder
	switch logFormat {
	case configs.GetSetting().Log.LogFormatJson:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 添加日志切割归档功能
	hook := lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    configs.GetSetting().Log.MaxSize,
		MaxBackups: configs.GetSetting().Log.MaxBackups,
		MaxAge:     configs.GetSetting().Log.MaxAge,
		Compress:   true,
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(&hook)),
		zap.NewAtomicLevelAt(logLevel), // 日志级别
	)
	caller := zap.AddCaller()
	development := zap.Development()
	logger := zap.New(core, caller, development)
	zap.ReplaceGlobals(logger)
}