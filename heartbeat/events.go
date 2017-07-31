package heartbeat

import (
	"time"
)

type HeartBeatStartingEvent struct {
	Timeout time.Duration
}

type HeartBeatTickedEvent struct{}
