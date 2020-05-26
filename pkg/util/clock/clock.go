package clock

import "time"

func UnixMillisecondFromTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
