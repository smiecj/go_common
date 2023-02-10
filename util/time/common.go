// package time 和时间相关的公共工具方法
package time

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	"github.com/smiecj/go_common/errorcode"
)

const (
	normalFormat    = "2006-01-02 15:04:05"
	dateFormat      = "2006-01-02"
	monthFormat     = "2006-01"
	shortDateFormat = "01-02"
)

var (
	weekDayToChineseMap = map[time.Weekday]string{
		time.Monday:    "一",
		time.Tuesday:   "二",
		time.Wednesday: "三",
		time.Thursday:  "四",
		time.Friday:    "五",
		time.Saturday:  "六",
		time.Sunday:    "日",
	}
)

// 获取当前时间戳
func GetCurrentTimestamp() string {
	return time.Now().Format(normalFormat)
}

// 获取当前日期
func GetCurrentDate() string {
	return time.Now().Format(dateFormat)
}

// 将指定时间戳和当前时间戳 做时间差值对比，当前时间 - 传入时间
// 若传入的时间格式不合法，返回 error
func CompareTimestampWithNow(timestamp string) (dur time.Duration, err error) {
	inputTime, err := time.Parse(normalFormat, timestamp)
	if nil != err {
		err = errorcode.BuildErrorWithMsg(errorcode.ParseTimeFailed, err.Error())
	} else {
		dur = time.Since(inputTime)
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

// unix 时间戳（秒）格式转 normal 格式
func GetTimestampByUnixtime(unixtime string) (targetTimestamp string, err error) {
	unixtimeInt, err := strconv.Atoi(unixtime)
	if nil != err {
		err = errorcode.BuildErrorWithMsg(errorcode.ParseTimeFailed, err.Error())
	}
	targetTime := time.Unix(int64(unixtimeInt), 0)
	targetTimestamp = targetTime.Format(normalFormat)
	return
}

// 获取指定日期的星期数
func GetWeekDayByDate(date string) (time.Weekday, error) {
	t, err := time.Parse(dateFormat, date)
	if nil != err {
		return 0, err
	}
	return t.Weekday(), nil
}

// 获取指定日期的年份数
func GetYearByDate(date string) (int, error) {
	t, err := time.Parse(dateFormat, date)
	if nil != err {
		return 0, err
	}
	return t.Year(), nil
}

// 获取指定日期的星期数（中文）
func GetWeekDayStringByDate(date string) (string, error) {
	t, err := time.Parse(dateFormat, date)
	if nil != err {
		return "", err
	}
	return weekDayToChineseMap[t.Weekday()], nil
}

// 获取当前日期指定天数之前的日期
func GetDateBeforeDay(day int) string {
	targetTime := time.Now().Add(-time.Duration(day*24) * time.Hour)
	return targetTime.Format(dateFormat)
}

// 获取当前日期指定天数之后的短日期
func GetShortDateAfterDay(day int) string {
	targetTime := time.Now().Add(time.Duration(day*24) * time.Hour)
	return targetTime.Format(shortDateFormat)
}

// 获取当前日期指定天数之后的日期
func GetDateAfterDay(day int) string {
	targetTime := time.Now().Add(time.Duration(day*24) * time.Hour)
	return targetTime.Format(dateFormat)
}

// 获取当前日期指定天数之后的月份
func GetMonthAfterDay(day int) string {
	targetTime := time.Now().Add(time.Duration(day*24) * time.Hour)
	return targetTime.Format(monthFormat)
}

// 获取指定天数之前的日期（阴历）
func GetLunarDateBeforeDay(day int) string {
	targetTime := time.Now().Add(-time.Duration(day*24) * time.Hour)
	lunarObj := calendar.BySolar(int64(targetTime.Year()), int64(targetTime.Month()), int64(targetTime.Day()),
		00, 00, 00)
	return fmt.Sprintf("%.4d-%.2d-%.2d", lunarObj.Lunar.GetYear(), lunarObj.Lunar.GetMonth(),
		lunarObj.Lunar.GetDay())
}

// 获取指定天数之后的短日期（阴历）
func GetShortLunarDateAfterDay(day int) string {
	targetTime := time.Now().Add(time.Duration(day*24) * time.Hour)
	lunarObj := calendar.BySolar(int64(targetTime.Year()), int64(targetTime.Month()), int64(targetTime.Day()),
		00, 00, 00)
	return fmt.Sprintf("%.2d-%.2d", lunarObj.Lunar.GetMonth(), lunarObj.Lunar.GetDay())
}

// 获取本周的最后一天（周日）日期
func GetThisWeekLastDate() string {
	currentTime := time.Now()
	currentWeekDate := currentTime.Weekday()
	var lastDateDiff int
	// 特殊逻辑: time 库中一周默认从 周日开始
	if time.Sunday == currentWeekDate {
		lastDateDiff = 0
	} else {
		lastDateDiff = int(time.Saturday-currentWeekDate) + 1
	}
	return currentTime.Add(time.Hour * 24 * time.Duration(lastDateDiff)).Format(dateFormat)
}

// 获取当天的0点的时间戳, 毫秒格式
func GetCurrentDateZeroTimestmapMill() int {
	currentTime := time.Now()
	zeroTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
		0, 0, 0, 0, currentTime.Location())
	return int(zeroTime.UnixMilli())
}
