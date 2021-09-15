package nicehash

import (
	"time"
)

type Time int64

func (t Time) Time() time.Time {
	return time.Unix(0, int64(t)/1000000)
}
