package helper

import (
	"time"
)

func UnixToTaipeiTime(seconds int64, nanos int32) string {
	t := time.Unix(seconds, int64(nanos))

	loc, _ := time.LoadLocation("Asia/Taipei")
	return t.In(loc).Format("2006-01-02 15:04:05")
}
