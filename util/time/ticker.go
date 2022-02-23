package time

import (
	"context"
	"sync"
	"time"

	"github.com/smiecj/go_common/util/log"
)

type TickerStatus string

const (
	StatusInit    TickerStatus = "init"
	StatusRunning TickerStatus = "running"
	StatusStop    TickerStatus = "stop"
)

// 调度器定义
type Ticker interface {
	Status() string
	Start() error
	Stop() error
}

// 按小时调度器，每天只调度一次，在指定小时的时候调度
// 需要获取 errorChan 中的数据，否则可能会导致阻塞
type fixHourTicker struct {
	ErrorChan chan error
	ticker    *time.Ticker

	tickFunc  func()
	tickOnce  sync.Once
	closeOnce sync.Once

	conf   *fixHourTickerConf
	status TickerStatus
	// 状态信息，记录上次调度的日期
	lastExecuteDate string
	todayHasRun     bool
}

// 按小时调度器具体配置
type fixHourTickerConf struct {
	hour   int
	f      func() error
	ctx    context.Context
	cancel context.CancelFunc
	// timeout time.Duration
}

type fixHourTickerConfFunc func(*fixHourTickerConf)

// 启动 fixed hour ticker
func (ticker *fixHourTicker) Start() error {
	ticker.tickOnce.Do(func() {
		go ticker.tickFunc()
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

// 当前: 实现指定小时调度的ticker
// 后续: 可参考 https://github.com/mileusna/crontab 实现更完整的 crontab 定时调度器
func NewFixHourTicker(confFuncArr ...fixHourTickerConfFunc) *fixHourTicker {
	conf := new(fixHourTickerConf)
	for _, currentConfFunc := range confFuncArr {
		currentConfFunc(conf)
	}
	hourTicker := getFixHourTicker(conf)

	// 开启协程 定时调度
	hourTicker.tickFunc = func() {
		for {
			select {
			case <-hourTicker.ticker.C:
				if time.Now().Hour() == conf.hour && !hourTicker.todayHasRun {
					log.Info("[FixHourTicker.tick] start")
					hourTicker.todayHasRun = true
					hourTicker.lastExecuteDate = GetCurrentDate()
					e := conf.f()
					if nil != e {
						log.Info("[FixHourTicker.tick] put error to chan")
						hourTicker.ErrorChan <- e
					}
					log.Info("[FixHourTicker.tick] end")
				} else if GetCurrentDate() != hourTicker.lastExecuteDate {
					hourTicker.todayHasRun = false
				}
			case <-conf.ctx.Done():
				hourTicker.ticker.Stop()
			}
		}
	}

	return hourTicker
}

// 设置定时调度的小时数，限制范围: 0~23
func SetHour(hour int) fixHourTickerConfFunc {
	return func(conf *fixHourTickerConf) {
		if hour >= 0 && hour < 24 {
			conf.hour = hour
		}
	}
}

// 设置定时调度的方法
func SetFunc(f func() error) fixHourTickerConfFunc {
	return func(conf *fixHourTickerConf) {
		conf.f = f
	}
}

// 设置定时调度的 context
func SetContext(ctx context.Context) fixHourTickerConfFunc {
	return func(conf *fixHourTickerConf) {
		conf.ctx, conf.cancel = context.WithCancel(ctx)
	}
}

func getFixHourTicker(conf *fixHourTickerConf) *fixHourTicker {
	hourTicker := fixHourTicker{}
	hourTicker.ticker = time.NewTicker(30 * time.Minute)
	hourTicker.ErrorChan = make(chan error)
	hourTicker.conf = conf
	hourTicker.status = StatusInit
	hourTicker.todayHasRun = false
	return &hourTicker
}
