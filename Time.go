package nicehash

import (
	"time"
)

type Time int64

func (t Time) Time() time.Time {
	msec := t % 1000
	sec := (t - msec) / 1000

	nsec := msec * 1000000

	return time.Unix(int64(sec), int64(nsec))
}
