package ticker

import (
	"context"
	"testing"
	"time"

	"github.com/smiecj/go_common/util/log"
	"github.com/stretchr/testify/require"
)

var (
	tickerConfArr = []tickerConf{
		{
			hour: time.Now().Hour(),
			f:    func() error { log.Info("hour ticker run"); return nil },
			ctx:  context.Background(),
		},
	}
)

// 测试按小时调度器
// 目前没有具体调度逻辑，只要不报错即算通过
func TestFixHourTicker(t *testing.T) {
	for _, currentTickerConf := range tickerConfArr {
		ticker := NewFixHourTicker(SetHour(currentTickerConf.hour), SetFunc(currentTickerConf.f), SetContext(currentTickerConf.ctx))
		ticker.Start()
		require.Equal(t, StatusRunning, ticker.Status())
		ticker.Stop()
	}
}
