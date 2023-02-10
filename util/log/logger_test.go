package log

import "testing"

func TestGlobalLogger(t *testing.T) {
	Info("global logger info")
	Warn("global logger warn")
	Error("global logger error")
	Info("[class] [func] do something.")
}

func TestCustomLogger(t *testing.T) {
	prefixLogger := PrefixLogger("this is prefix")
	prefixLogger.Info("after prefix")
}
