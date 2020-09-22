package toolbox

import (
	"fmt"
	"log"
	"log/syslog"
	"path"
	"runtime"
	"strings"

	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"

	"github.com/sirupsen/logrus"
)

// This will be used for the default go log messages
var defaultLogger = logrus.New()

// LogParameters contains the default log parameters
type LogParameters struct {
	Level string `kong:"help='Logging level',default='debug',enum='debug,info,warning,error'"`
	Type  string `kong:"help='Log type',default='plain',enum='syslog,plain'"`
}

// LogLevel returns the log level as a value
func logLevelFromString(level string) logrus.Level {
	if level == "" {
		return logrus.DebugLevel
	}
	switch strings.ToLower(level)[0] {
	case 'd':
		return logrus.DebugLevel
	case 'i':
		return logrus.InfoLevel
	case 'w':
		return logrus.WarnLevel
	case 'e':
		return logrus.ErrorLevel
	default:
		return logrus.WarnLevel
	}
}

const (
	plainLogs = "plain"
	sysLogs   = "syslog"
)

// InitLogs configures logging for a service. This will turn on syslog
// logs  if they are enabled or just a plain text stderr log. If there
// are any errors while running
func InitLogs(service string, params LogParameters) {
	logrus.SetLevel(logLevelFromString(params.Level))
	w := defaultLogger.Writer()
	// Note: This writer will never be closed so technically we're leaking
	// a handle here
	log.SetOutput(w)
	if params.Type == sysLogs {
		// Enable syslog
		hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, service)
		if err != nil {
			logrus.WithError(err).Fatal("Could not attach syslog")
			return
		}
		logrus.AddHook(hook)
	}
	if params.Type == plainLogs {
		formatter := &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "15:04:05.000",
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				file = fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
				function = ""
				return
			},
		}
		logrus.SetFormatter(formatter)
		defaultLogger.Formatter = formatter
	}

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.SetReportCaller(true)
	}
}
