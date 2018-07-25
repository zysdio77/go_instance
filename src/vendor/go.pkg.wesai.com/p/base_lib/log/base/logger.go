package base

import "go.pkg.wesai.com/p/base_lib/log/base/field"

// LoggerType 代表日志记录器类型。
type LoggerType string

const (
	// LOGRUS 代表日志代码包logrus。
	LOGRUS LoggerType = "logrus"
	// ZAP 代表日志代码包zap。
	//ZAP LoggerType = "zap"
)

// MyLogger 代表日志记录器接口。
type MyLogger interface {
	Name() string

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})

	// 增加需记录的额外字段。
	WithFields(fields ...field.Field) MyLogger
}
