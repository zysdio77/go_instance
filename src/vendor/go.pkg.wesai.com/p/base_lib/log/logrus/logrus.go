package logrus // import "go.pkg.wesai.com/p/base_lib/log/logrus"

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.pkg.wesai.com/p/base_lib/log/base"
	"go.pkg.wesai.com/p/base_lib/log/base/field"

	"github.com/Sirupsen/logrus"
)

func init() {
}

type logger_logrus struct {
	projectName string
	inner       *logrus.Entry
}

// 新建并返回一个日志记录器。
func NewLogger(projectName string) base.MyLogger {
	return NewLoggerBy(projectName, os.Stderr, logrus.DebugLevel)
}

// 根据指定的参数新建并返回一个日志记录器。
func NewLoggerBy(projectName string, w io.Writer, l logrus.Level) base.MyLogger {
	return &logger_logrus{
		projectName: projectName,
		inner:       initInnerLogger(w, l),
	}
}

func initInnerLogger(w io.Writer, l logrus.Level) *logrus.Entry {
	innerLogger := logrus.New()
	innerLogger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: base.TIMESTAMP_FORMAT,
	}
	innerLogger.Level = l
	innerLogger.Out = w
	return logrus.NewEntry(innerLogger)
}

func (logger *logger_logrus) Name() string {
	return "logrus"
}

func (logger *logger_logrus) Debug(v ...interface{}) {
	if base.DebugEnable(logger.projectName) {
		appendRequiredFields(logger.inner).Debug(v...)
	}
}

func (logger *logger_logrus) Debugf(format string, v ...interface{}) {
	if base.DebugEnable(logger.projectName) {
		appendRequiredFields(logger.inner).Debugf(format, v...)
	}
}

func (logger *logger_logrus) Debugln(v ...interface{}) {
	if base.DebugEnable(logger.projectName) {
		appendRequiredFields(logger.inner).Debug(v...)
	}
}

func (logger *logger_logrus) Error(v ...interface{}) {
	appendRequiredFields(logger.inner).Error(genMsg("", v...))
}

func (logger *logger_logrus) Errorf(format string, v ...interface{}) {
	appendRequiredFields(logger.inner).Errorf(genMsg(format, v...))
}

func (logger *logger_logrus) Errorln(v ...interface{}) {
	appendRequiredFields(logger.inner).Errorln(genMsg("", v...))
}

func (logger *logger_logrus) Fatal(v ...interface{}) {
	appendRequiredFields(logger.inner).Fatal(genMsg("", v...))
}

func (logger *logger_logrus) Fatalf(format string, v ...interface{}) {
	appendRequiredFields(logger.inner).Fatalf(genMsg(format, v...))
}

func (logger *logger_logrus) Fatalln(v ...interface{}) {
	appendRequiredFields(logger.inner).Fatalln(genMsg("", v...))
}

func (logger *logger_logrus) Info(v ...interface{}) {
	appendRequiredFields(logger.inner).Info(genMsg("", v...))
}

func (logger *logger_logrus) Infof(format string, v ...interface{}) {
	appendRequiredFields(logger.inner).Infof(genMsg(format, v...))
}

func (logger *logger_logrus) Infoln(v ...interface{}) {
	appendRequiredFields(logger.inner).Infoln(genMsg("", v...))
}

func (logger *logger_logrus) Panic(v ...interface{}) {
	appendRequiredFields(logger.inner).Panic(genMsg("", v...))
}

func (logger *logger_logrus) Panicf(format string, v ...interface{}) {
	appendRequiredFields(logger.inner).Panicf(genMsg(format, v...))
}

func (logger *logger_logrus) Panicln(v ...interface{}) {
	appendRequiredFields(logger.inner).Panicln(genMsg("", v...))
}

func (logger *logger_logrus) Warn(v ...interface{}) {
	appendRequiredFields(logger.inner).Warning(genMsg("", v...))
}

func (logger *logger_logrus) Warnf(format string, v ...interface{}) {
	appendRequiredFields(logger.inner).Warningf(genMsg(format, v...))
}

func (logger *logger_logrus) Warnln(v ...interface{}) {
	appendRequiredFields(logger.inner).Warningln(genMsg("", v...))
}

func (logger *logger_logrus) WithFields(fields ...field.Field) base.MyLogger {
	fieldsLen := len(fields)
	if fieldsLen == 0 {
		return logger
	}
	logrusFields := make(map[string]interface{}, fieldsLen)
	for _, curfield := range fields {
		logrusFields[curfield.Name()] = curfield.Value()
	}
	return &logger_logrus{
		projectName: logger.projectName,
		inner:       logger.inner.WithFields(logrusFields),
	}
}

// 添加必要的字段。
func appendRequiredFields(logger *logrus.Entry) *logrus.Entry {
	return appendLocation(appendTimestamp(logger))
}

// 添加时间戳。
func appendTimestamp(logger *logrus.Entry) *logrus.Entry {
	return logger.WithField("ts", time.Now().UnixNano())
}

// 添加记录日志的代码位置。
func appendLocation(logger *logrus.Entry) *logrus.Entry {
	funcPath, fileName, line := base.GetInvokerLocation(4)
	return logger.WithField(
		"location", map[string]interface{}{
			"func_path": funcPath,
			"file_name": fileName,
			"line":      line,
		},
	)
}

// 生成日志消息。
func genMsg(format string, v ...interface{}) string {
	if len(format) > 0 {
		return fmt.Sprintf(format, v...)
	} else {
		return fmt.Sprint(v...)
	}
}
