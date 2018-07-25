package log

import (
	"testing"

	"go.pkg.wesai.com/p/base_lib/log/base"
)

func TestLogger(t *testing.T) {

	// default logger
	logger := DLogger()
	if logger == nil {
		t.Fatal("The default logger is invalid!")
	}
	if logger.Name() != "logrus" {
		t.Fatalf("Expect '%s', but the actual is '%s'\n", "logrus", logger.Name())
	}
	t.Logf("The default logger: %#v\n", logger)

	// logrus logger
	logrusLogger := Logger(base.LOGRUS)
	if logrusLogger == nil {
		t.Fatal("The logrus logger is invalid!")
	}
	if logrusLogger.Name() != "logrus" {
		t.Fatalf("Expect '%s', but the actual is '%s'\n", "logrus", logrusLogger.Name())
	}
	t.Logf("The logrus logger: %#v\n", logrusLogger)

	// zap logger
	zapLogger := Logger(base.ZAP)
	if logrusLogger == nil {
		t.Fatal("The zap logger is invalid!")
	}
	if zapLogger.Name() != "zap" {
		t.Fatalf("Expect '%s', but the actual is '%s'\n", "zap", zapLogger.Name())
	}
	t.Logf("The zap logger: %#v\n", zapLogger)
}
