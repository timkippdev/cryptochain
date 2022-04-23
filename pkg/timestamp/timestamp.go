package timestamp

import (
	"strconv"
	"time"
)

type Timestamp int64

func FromTime(t time.Time) Timestamp {
	return Timestamp(t.UTC().UnixNano() / int64(time.Millisecond))
}

func Now() Timestamp {
	return Timestamp(time.Now().UTC().UnixNano() / int64(time.Millisecond))
}

func (t Timestamp) String() string {
	return strconv.FormatInt(int64(t), 10)
}

func (t Timestamp) Add(v int64) Timestamp {
	return Timestamp(int64(t) + v)
}
