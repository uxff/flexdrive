package utils

import (
	"time"
)

type JsonTime time.Time

func NewJsonTimeNow() JsonTime {
	return JsonTime(time.Now())
}

// 重写 time.Time 的 MarshalJSON
func (j JsonTime) MarshalJSON() ([]byte, error) {
	if j.IsEmptyTime() {
		return []byte(`"` + ZeroTimeStr + `"`), nil
	}
	return []byte(`"` + j.String() + `"`), nil
}

// 重写 time.Time 的 UnmarshalJSON
func (j *JsonTime) UnmarshalJSON(b []byte) error {
	now, err := time.ParseInLocation(`"`+DefaultTimeFmt+`"`, string(b), time.Local)
	*j = JsonTime(now)
	return err
}

// 重写 time.Time 的 String
func (j JsonTime) String() string {
	str := time.Time(j).Format(DefaultTimeFmt)
	if str == StdZeroTimeStr {
		return ZeroTimeStr
	}
	return str
}

// 判断当前时间是不是"0"值
func (j JsonTime) IsEmptyTime() bool {
	return j.String() == ZeroTimeStr
}

// 转换为标准库里time.Time类型
func (j JsonTime) ToStdTime() time.Time {
	return time.Time(j)
}
