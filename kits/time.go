package kits

import (
	"fmt"
	"time"
)

// 按照 年-月-日_时:分:秒 格式形式返回时间字符串
func GetCurrentTimeString() string {
	return time.Now().Format("2006-01-02_15:04:05")
}

// 统计给定时间戳到当前的时间消耗, 单位:秒(s), float64
func StatCostTime(st time.Time) float64 {
	return time.Since(st).Seconds()
}

// 统计给定时间戳到当前的时间消耗, 单位:秒(s), string
func StatCostTimeString(st time.Time, nameStr string) string {
	return fmt.Sprintf("%s cost time %3.2f", nameStr, StatCostTime(st))
}

// 按照 年-月-日 时:分:秒 格式形式返回时间字符串
func GetTimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 获取当前月开始时间格式化字符串
func GetCurrentMonthStartTimeString() string {
	return GetTimeString(GetCurrentMonthStartTime())
}

// 获取当前月结束时间格式化字符串
func GetCurrentMonthEndTimeString() string {
	return GetTimeString(GetCurrentMonthEndTime())
}

// 获取当前月开始时间时间戳
func GetCurrentMonthStartTimeTs() int64 {
	return GetCurrentMonthStartTime().Unix()
}

// 获取当前月结束时间时间戳
func GetCurrentMonthEndTimeTs() int64 {
	return GetCurrentMonthEndTime().Unix()
}

// 获取当前月开始时间
func GetCurrentMonthStartTime() time.Time {
	return GetMonthStartTimeForTime(time.Now())
}

// 获取当前月结束时间
func GetCurrentMonthEndTime() time.Time {
	return GetMonthEndTimeForTime(time.Now())
}

// 获取指定时间所在月份起始时间
func GetMonthStartTimeForTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// 获取指定时间所在月份结束时间 -> 月最后一天的 23:59:59
func GetMonthEndTimeForTime(t time.Time) time.Time {
	return GetMonthStartTimeForTime(t.AddDate(0, 1, 0)).Add(-1 * time.Second)
}
