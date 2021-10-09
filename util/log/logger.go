package log

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// 自定义formatter, 参考: https://cloud.tencent.com/developer/article/1830707
// 使用方式: log.Info("[method name] log msg: %s", msg)
const (
	timeFormat = "2006-01-02 15:04:05"
)

type MyFormatter struct {
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(timeFormat)
	var newLog string
	newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func init() {
	log.SetFormatter(&MyFormatter{})
	// some other formatter: &log.JSONFormatter{}
}

func Debug(format string, args ...interface{}) {
	log.Debugf(format, args)
}

func Info(format string, args ...interface{}) {
	log.Infof(format, args)
}

func Warn(format string, args ...interface{}) {
	log.Warnf(format, args)
}

func Error(format string, args ...interface{}) {
	log.Errorf(format, args)
}
