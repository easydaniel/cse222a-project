package main

import (
	"math"
	"sync"
)

type LeastTime struct {
	m    sync.Mutex
	id   int
	time []int64
}

func (sched *LeastTime) Schedule(config *LBConfig, httpparser *HttpParser) int {
	sched.m.Lock()
	defer sched.m.Unlock()
	if sched.time == nil {
		sched.time = make([]int64, len(config.Servers))
	}
	mn, mnid := int64(math.MaxInt64), 0
	for id, p := range sched.time {
		if p < mn {
			mn = p
			mnid = id
		}
	}
	return mnid
}
func (sched *LeastTime) Finish(config *LBConfig, httpparser *HttpParser, srvid int, dur int64) {
	sched.m.Lock()
	defer sched.m.Unlock()
	sched.time[srvid] = sched.time[srvid] + dur
}
