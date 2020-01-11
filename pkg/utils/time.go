package utils

import (
	"time"
)

const (
	DefaultTimeFmt = "2006-01-02 15:04:05"
	StdZeroTimeStr = "0001-01-01 00:00:00" // 标准库里time.Time{} jsonMashal以后的字符串
	ZeroTimeStr    = "0000-00-00 00:00:00" // 我们自己期望使用的时间 "0" 值字符串
)

// time.Time 转字符串
func TimeToString(t time.Time) string {
	return t.Format(DefaultTimeFmt)
}

// 字符串转 time.Time
func StringToTime(str string) time.Time {
	local, err := time.LoadLocation("Local")
	if err != nil {
		panic(err)
	}
	t, err := time.ParseInLocation(DefaultTimeFmt, str, local)
	if err != nil {
		panic(err)
	}
	return t
}

// 判断当前时间是不是"0"值
func IsEmptyTime(t time.Time) bool {
	return t == time.Time{}
}
