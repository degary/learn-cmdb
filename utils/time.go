package utils

import (
	"time"
)

func StrToTime(ti string) *time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation("2006-01-02", ti, loc)
	return &t
}
