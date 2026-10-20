package helper

import (
	"log/slog"
	"time"
)

func UnixToTaipeiTime(seconds int64, nanos int32) string {
	slog.Debug("Time: ", "seconds", seconds, "nanoseconds", nanos)
	t := time.Unix(seconds, int64(nanos))

	loc, _ := time.LoadLocation("Asia/Taipei")
	return t.In(loc).Format("2006-01-02 15:04:05")
}
