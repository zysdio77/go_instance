package logrus

import (
	"testing"

	"go.pkg.wesai.com/p/base_lib/log/base/field"
)

func TestLogrusLogger(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			switch i := p.(type) {
			case error, string:
				t.Fatalf("Fatal error: %s\n", i)
			default:
				t.Fatalf("Fatal error: %#v\n", i)
			}
		}
	}()
	logger := NewLogger("testing")
	logger = logger.WithFields(
		field.Bool("bool", false),
		field.Int64("int64", 12345678),
		field.Float64("float64", 123.456),
		field.String("string", "logrus"),
		field.Object("object", interface{}("abcd")),
	)
	t.Logf("The tested logger: %s", logger.Name())
	logger.Infof("The tested logger: %s", logger.Name())
	logger.Info("Info log (logrus)")
	logger.Infoln("Infoln log (logrus)")
	logger.Error("Error log (logrus)")
	logger.Errorf("%s log (logrus)", "Errorf")
	logger.Errorln("Errorln log (logrus)")
	logger.Warn("Warn log (logrus)")
	logger.Warnf("%s log (logrus)", "Warnf")
	logger.Warnln("Warnln log (logrus)")

	// They will call os.Exit(1)
	// logger.Fatal("Fatal log (logrus)")
	// logger.Fatalf("%s log (logrus)", "Fatalf")
	// logger.Fatalln("Fatalln log (logrus)")

	// They will cause panic
	// logger.Panic("Panic log (logrus)")
	// logger.Panicf("%s log (logrus)", "Panicf")
	// logger.Panicln("Panicln log (logrus)")
}
