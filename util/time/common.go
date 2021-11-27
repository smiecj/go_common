// package time 和时间相关的公共工具方法
package time

import (
	"time"

	"github.com/smiecj/go_common/errorcode"
)

const (
	normalFormat = "2006-01-02 15:04:05"
	dateFormat   = "2006-01-02"
)

// 获取当前时间戳
func GetCurrentTimestamp() string {
	return time.Now().Format(normalFormat)
}

// 将指定时间戳和当前时间戳 做时间差值对比，当前时间 - 传入时间
// 若传入的时间格式不合法，返回 error
func CompareTimestampWithNow(timestamp string) (dur time.Duration, err error) {
	inputTime, err := time.Parse(normalFormat, timestamp)
	if nil != err {
		err = errorcode.BuildErrorWithMsg(errorcode.ParseTimeFailed, err.Error())
	} else {
		dur = time.Now().Sub(inputTime)
	}
	return
}

// 获取当前时间 指定偏移时间后的时间戳
func GetCurrentTimeAfterDuration(dur time.Duration) string {
	return time.Now().Add(dur).Format(normalFormat)
}

// 获取指定时间戳、指定偏移时间后的时间戳
func GetTimestampAfterDuration(startTimestamp string, dur time.Duration) (targetTimestamp string, err error) {
	startTime, err := time.Parse(normalFormat, startTimestamp)
	if nil != err {
		err = errorcode.BuildErrorWithMsg(errorcode.ParseTimeFailed, err.Error())
	} else {
		targetTime := startTime.Add(dur)
		targetTimestamp = targetTime.Format(normalFormat)
	}
	return
}
