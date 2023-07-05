package Dateutil

import (
	"time"
)

// DateUtils 是日期工具类
type DateUtils struct{}

// Format 格式化日期为指定格式的字符串
func (du *DateUtils) Format(date time.Time, format string) string {
	return date.Format(format)
}

// GetCurrentTimestamp 获取当前时间戳（秒）
func (du *DateUtils) GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetHistoricalDate 计算距离今日指定天数的历史日期
func (du *DateUtils) GetHistoricalDate(days int) time.Time {
	today := time.Now()
	historicalDate := today.AddDate(0, 0, -days)
	return historicalDate
}

// GetFutureDate 计算距离今日指定天数的未来日期
func (du *DateUtils) GetFutureDate(days int) time.Time {
	today := time.Now()
	futureDate := today.AddDate(0, 0, days)
	return futureDate
}

//大小比较
func (du *DateUtils) CompareDates(date1, date2 time.Time) int {
	if date1.Before(date2) {
		return -1
	} else if date1.After(date2) {
		return 1
	}
	return 0
}
