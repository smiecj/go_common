// package ticker 定时调度器
package ticker

import (
	"context"
	"sync"
	"time"

	"github.com/smiecj/go_common/errorcode"
	"github.com/smiecj/go_common/util/log"
	timeutil "github.com/smiecj/go_common/util/time"
)

type TickerStatus string

const (
	StatusInit    TickerStatus = "init"
	StatusRunning TickerStatus = "running"
	StatusStop    TickerStatus = "stop"

	fixedHourTickerCheckDuration = 5 * time.Minute
)

// 调度器定义
type Ticker interface {
	Status() TickerStatus
	Start() error
	Stop() error
	Error() <-chan error
}

// 按小时调度器，每天只调度一次，在指定小时的时候调度
// 需要获取 errorChan 中的数据，否则可能会导致阻塞
type fixHourTicker struct {
	errorChan chan error
	ticker    *time.Ticker

	tickLoopFunc func()
	tickOnce     sync.Once
	closeOnce    sync.Once

	conf   *tickerConf
	status TickerStatus
	// 状态信息，记录上次调度的日期
	lastExecuteDate string
	todayHasRun     bool
}

// ticker 具体配置
type tickerConf struct {
	name          string
	hour          int
	f             func() error
	ctx           context.Context
	cancel        context.CancelFunc
	isIgnoreError bool
	// timeout time.Duration

}

type tickerConfFunc func(*tickerConf) error

// 启动 fixed hour ticker
func (ticker *fixHourTicker) Start() error {
	ticker.tickOnce.Do(func() {
		go ticker.tickLoopFunc()
		ticker.status = StatusRunning
	})
	return nil
}

// 停止 fixed hour ticker
func (ticker *fixHourTicker) Stop() error {
	ticker.closeOnce.Do(func() {
		ticker.conf.cancel()
		ticker.status = StatusStop
	})
	return nil
}

// 获取状态 fixed hour ticker
func (ticker *fixHourTicker) Status() TickerStatus {
	return ticker.status
}

// 获取 error chan
func (ticker *fixHourTicker) Error() <-chan error {
	return ticker.errorChan
}

// 后续: 可参考 https://github.com/mileusna/crontab 实现更完整的 crontab 定时调度器
func NewFixHourTicker(confFuncArr ...tickerConfFunc) Ticker {
	conf := getTickerConf()
	for _, currentConfFunc := range confFuncArr {
		// 暂时不会有致命错误，先不处理
		_ = currentConfFunc(conf)
	}
	hourTicker := getFixHourTicker(conf)

	// 开启协程 定时调度
	hourTicker.tickLoopFunc = func() {
		for {
			select {
			case <-hourTicker.ticker.C:
				if time.Now().Hour() == conf.hour && !hourTicker.todayHasRun {
					log.Info("[FixHourTicker.tick] start")
					hourTicker.todayHasRun = true
					hourTicker.lastExecuteDate = timeutil.GetCurrentDate()
					// 后续: 支持选择 同步 or 异步，目前是同步
					jobFinishChan := make(chan struct{})
					go func() {
						defer func() {
							if err := recover(); nil != err {
								log.Warn("[FixHourTicker.tick] job exec throw err: %s", err)
							}
							close(jobFinishChan)
						}()
						e := conf.f()
						if nil != e && !conf.isIgnoreError {
							log.Warn("[FixHourTicker.tick] put error to chan")
							hourTicker.errorChan <- e
						}
					}()
					<-jobFinishChan
					log.Info("[FixHourTicker.tick] end")
				} else if timeutil.GetCurrentDate() != hourTicker.lastExecuteDate {
					hourTicker.todayHasRun = false
				}
			case <-conf.ctx.Done():
				hourTicker.ticker.Stop()
				close(hourTicker.errorChan)
			}
		}
	}

	return hourTicker
}

// 设置定时调度的小时数，限制范围: 0~23
func SetHour(hour int) tickerConfFunc {
	return func(conf *tickerConf) error {
		if hour >= 0 && hour < 24 {
			conf.hour = hour
		} else {
			return errorcode.BuildError(errorcode.ServiceError)
		}
		return nil
	}
}

// 设置定时调度的方法
func SetFunc(f func() error) tickerConfFunc {
	return func(conf *tickerConf) error {
		conf.f = f
		return nil
	}
}

// 设置定时调度的 context
func SetContext(ctx context.Context) tickerConfFunc {
	return func(conf *tickerConf) error {
		conf.ctx, conf.cancel = context.WithCancel(ctx)
		return nil
	}
}

// 设置忽略 error，这样调用方不再需要处理 Error() 返回的 chan
// 后续: error chan 使用懒加载的方式，节省空间
func SetIsIgnoreError(isIgnoreError bool) tickerConfFunc {
	return func(conf *tickerConf) error {
		conf.isIgnoreError = isIgnoreError
		return nil
	}
}

// 获取带有一些初始化配置的 定时调度配置
func getTickerConf() *tickerConf {
	conf := new(tickerConf)
	conf.ctx = context.Background()
	return conf
}

func getFixHourTicker(conf *tickerConf) *fixHourTicker {
	hourTicker := fixHourTicker{}
	hourTicker.ticker = time.NewTicker(fixedHourTickerCheckDuration)
	hourTicker.errorChan = make(chan error)
	hourTicker.conf = conf
	hourTicker.status = StatusInit
	hourTicker.todayHasRun = false
	return &hourTicker
}
