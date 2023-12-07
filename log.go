package onion_log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gitlab.yc345.tv/backend/onion-log/formatter"
	"gitlab.yc345.tv/backend/onion-log/hook"
	"gitlab.yc345.tv/backend/onion-log/logger"
)

type Log struct {
	Level  logrus.Level
	Logger *logrus.Logger
}

// New 实例化 Log
//     level: debug info warn error panic
//     goEnv: development stage production
func New(level string, goEnv string, hooks ...logrus.Hook) *Log {
	parseLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parseLevel = logrus.InfoLevel
		fmt.Printf("err: log level is err:(%s) set level info", err)
	}

	l := logrus.New()
	l.Formatter = &formatter.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999",
	}
	l.Level = parseLevel
	l.Out = os.Stdout

	if len(hooks) > 0 {
		for _, v := range hooks {
			l.AddHook(v)
		}
	}

	// 记录行号 // 并发环境会有性能问题; 只在非生产环境使用
	if goEnv != "production" {
		l.AddHook(hook.NewLineHook(logrus.ErrorLevel, logrus.PanicLevel))
	}

	return &Log{
		Level:  parseLevel,
		Logger: l,
	}
}

// With 携带基础信息
func (l *Log) With(baseContentInfo *logger.BaseContentInfo) *logger.Logger {
	if baseContentInfo == nil {
		baseContentInfo = &logger.BaseContentInfo{}
	}

	// if baseContentInfo.TraceID == "" {
	//	baseContentInfo.InitTraceID()
	// }

	return &logger.Logger{
		Content: baseContentInfo,
		Entry: l.Logger.WithFields(logrus.Fields{
			"uid":     baseContentInfo.UID,
			"traceId": baseContentInfo.TraceID,
			"appFrom": baseContentInfo.From,
			"spanId":  baseContentInfo.SpanID,
			"type":    "content",
		}),
	}
}
