package pretty

import (
	"strconv"
	"time"
)

func Duration(dur time.Duration) string {
	var out string
	if amount := dur.Nanoseconds(); amount < 1_000 {
		out = strconv.Itoa(int(amount)) + "ns"
	} else if amount := dur.Microseconds(); amount < 1_000 {
		out = strconv.Itoa(int(amount)) + "µs"
	} else if amount := dur.Milliseconds(); amount < 1_000 {
		out = strconv.Itoa(int(amount)) + "ms"
	} else {
		out = strconv.Itoa(int(dur.Seconds())) + "s"
	}
	return out
}
