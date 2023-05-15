// package time 和时间相关的公共工具方法
package time

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	"github.com/smiecj/go_common/errorcode"
)

const (
	normalFormat    = "2006-01-02 15:04:05"
	dateFormat      = "2006-01-02"
	monthFormat     = "2006-01"
	minuteFormat    = "15:04"
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

	monthEngToIntMap = map[string]int{
		"January":   1,
		"February":  2,
		"March":     3,
		"April":     4,
		"May":       5,
		"June":      6,
		"July":      7,
		"August":    8,
		"September": 9,
		"October":   10,
		"November":  11,
		"December":  12,
	}
)

// 获取当前时间戳
func CurrentTimestamp() string {
	return time.Now().Format(normalFormat)
}

// 获取当前日期
func CurrentDate() string {
	return time.Now().Format(dateFormat)
}

// 获取当前分钟
func CurrentMinute() string {
	return time.Now().Format(minuteFormat)
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
func CurrentTimeAfterDuration(dur time.Duration) string {
	return time.Now().Add(dur).Format(normalFormat)
}

// 获取指定时间戳、指定偏移时间后的时间戳
func TimestampAfterDuration(startTimestamp string, dur time.Duration) (targetTimestamp string, err error) {
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
func TimestampByUnixtimeStr(unixtime string) (targetTimestamp string, err error) {
	unixtimeInt, err := strconv.Atoi(unixtime)
	if nil != err {
		err = errorcode.BuildErrorWithMsg(errorcode.ParseTimeFailed, err.Error())
	}
	targetTime := time.Unix(int64(unixtimeInt), 0)
	targetTimestamp = targetTime.Format(normalFormat)
	return
}

// unix 时间戳 （秒） 格式转 date 格式
func DateByUnixtimeStr(unixtime string) (dateStr string, err error) {
	unixtimeInt, err := strconv.Atoi(unixtime)
	if nil != err {
		err = errorcode.BuildErrorWithMsg(errorcode.ParseTimeFailed, err.Error())
	}
	targetTime := time.Unix(int64(unixtimeInt), 0)
	dateStr = targetTime.Format(dateFormat)
	return
}

// unix 时间戳（毫秒）格式转 normal 格式
func TimestampByUnixMill(unixMill int) string {
	targetTime := time.UnixMilli(int64(unixMill))
	targetTimestamp := targetTime.Format(normalFormat)
	return targetTimestamp
}

// 获取指定日期的星期数
func WeekDayByDate(date string) (time.Weekday, error) {
	t, err := time.Parse(dateFormat, date)
	if nil != err {
		return 0, err
	}
	return t.Weekday(), nil
}

// 获取指定日期的年份数
func YearByDate(date string) (int, error) {
	t, err := time.Parse(dateFormat, date)
	if nil != err {
		return 0, err
	}
	return t.Year(), nil
}

// 获取指定日期的星期数（中文）
func WeekDayStringByDate(date string) (string, error) {
	t, err := time.Parse(dateFormat, date)
	if nil != err {
		return "", err
	}
	return weekDayToChineseMap[t.Weekday()], nil
}

// 获取当前日期指定天数之前的日期
func DateBeforeDay(day int) string {
	targetTime := time.Now().Add(-time.Duration(day*24) * time.Hour)
	return targetTime.Format(dateFormat)
}

// 获取当前日期指定天数之后的短日期
func ShortDateAfterDay(day int) string {
	targetTime := time.Now().Add(time.Duration(day*24) * time.Hour)
	return targetTime.Format(shortDateFormat)
}

// 获取当前日期指定天数之后的日期
func DateAfterDay(day int) string {
	targetTime := time.Now().Add(time.Duration(day*24) * time.Hour)
	return targetTime.Format(dateFormat)
}

// 获取当前日期指定天数之后的月份
func MonthAfterDay(day int) string {
	targetTime := time.Now().Add(time.Duration(day*24) * time.Hour)
	return targetTime.Format(monthFormat)
}

// 获取指定天数之前的日期（阴历）
func LunarDateBeforeDay(day int) string {
	targetTime := time.Now().Add(-time.Duration(day*24) * time.Hour)
	lunarObj := calendar.BySolar(int64(targetTime.Year()), int64(targetTime.Month()), int64(targetTime.Day()),
		00, 00, 00)
	return fmt.Sprintf("%.4d-%.2d-%.2d", lunarObj.Lunar.GetYear(), lunarObj.Lunar.GetMonth(),
		lunarObj.Lunar.GetDay())
}

// 获取指定天数之后的短日期（阴历）
func ShortLunarDateAfterDay(day int) string {
	targetTime := time.Now().Add(time.Duration(day*24) * time.Hour)
	lunarObj := calendar.BySolar(int64(targetTime.Year()), int64(targetTime.Month()), int64(targetTime.Day()),
		00, 00, 00)
	return fmt.Sprintf("%.2d-%.2d", lunarObj.Lunar.GetMonth(), lunarObj.Lunar.GetDay())
}

// 获取本周的最后一天（周日）日期
func ThisWeekLastDate() string {
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
func CurrentDateZeroTimestmapMill() int {
	currentTime := time.Now()
	zeroTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(),
		0, 0, 0, 0, currentTime.Location())
	return int(zeroTime.UnixMilli())
}

// "1 February 2021" to "2021-02-01"
func TimestampByGoFormat(timestamp string) (string, error) {
	timestampSplitArr := strings.Split(timestamp, " ")
	if len(timestampSplitArr) != 3 {
		return "", errorcode.BuildError(errorcode.ParseTimeFailed)
	}

	date, month, year := timestampSplitArr[0], timestampSplitArr[1], timestampSplitArr[2]
	var monthInt int
	var ok bool
	if monthInt, ok = monthEngToIntMap[month]; !ok {
		return "", errorcode.BuildError(errorcode.ParseTimeFailed)
	}
	dateInt, _ := strconv.Atoi(date)
	return fmt.Sprintf("%s-%02d-%02d", year, monthInt, dateInt), nil
}
