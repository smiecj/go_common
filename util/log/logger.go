package log

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	prefixed "github.com/smiecj/logrus-prefixed-formatter"
)

type LogLevel int

var (
	globalFormatter log.Formatter
	globalLogger    Logger
)

// 自定义formatter, 参考: https://cloud.tencent.com/developer/article/1830707
// 使用方式: log.Info("[method name] log msg: %s", msg)
const (
	timeFormat = "2006-01-02 15:04:05"
)

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

func init() {
	// TextFormatter: https://stackoverflow.com/a/48972299
	globalFormatter = &prefixed.TextFormatter{
		DisableColors:   true,
		TimestampFormat: timeFormat,
		FullTimestamp:   true,
		ForceFormatting: true,
	}
	// globalFormatter = &log.TextFormatter{}
	// globalFormatter = &defaultFormatter{}

	// global logger
	newLogger := log.New()
	newLogger.SetFormatter(globalFormatter)
	newLogger.SetLevel(log.InfoLevel)
	globalLogger = &customLogger{
		logger: newLogger,
	}
}

type defaultFormatter struct {
}

func (m *defaultFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(timeFormat)
	newLog := fmt.Sprintf("[%s] %s %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

type prefixFormatter struct {
	prefix          string
	parentFormatter log.Formatter
}

func (f *prefixFormatter) Format(entry *log.Entry) ([]byte, error) {
	entry.Data["prefix"] = f.prefix

	return f.parentFormatter.Format(entry)
}

func SetLevel(level LogLevel) {
	globalLogger.SetLevel(level)
}

func Debug(format string, args ...interface{}) {
	globalLogger.Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	globalLogger.Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	globalLogger.Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	globalLogger.Error(format, args...)
}

/********* 自定义logger *********/
type Logger interface {
	SetLevel(level LogLevel)
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type customLogger struct {
	logger *log.Logger
}

func PrefixLogger(prefix string) Logger {
	logger := log.New()
	logger.SetFormatter(&prefixFormatter{
		prefix:          prefix,
		parentFormatter: globalFormatter,
	})
	return &customLogger{
		logger: logger,
	}
}

func (cl *customLogger) SetLevel(level LogLevel) {
	logLevel := log.InfoLevel
	switch level {
	case LevelDebug:
		logLevel = log.DebugLevel
	case LevelInfo:
		logLevel = log.InfoLevel
	case LevelWarn:
		logLevel = log.WarnLevel
	case LevelError:
		logLevel = log.ErrorLevel
	}
	cl.logger.SetLevel(logLevel)
}

func (cl *customLogger) Debug(format string, args ...interface{}) {
	cl.logger.Debugf(format, args...)
}

func (cl *customLogger) Info(format string, args ...interface{}) {
	cl.logger.Infof(format, args...)
}

func (cl *customLogger) Warn(format string, args ...interface{}) {
	cl.logger.Warnf(format, args...)
}

func (cl *customLogger) Error(format string, args ...interface{}) {
	cl.logger.Errorf(format, args...)
}
