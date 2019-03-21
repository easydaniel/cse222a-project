package main

import (
	"math"
	"sync"
)

type LeastLoadBad struct {
	m       sync.Mutex
	id      int
	pending []int
}

func (sched *LeastLoadBad) Schedule(config *LBConfig, httpparser *HttpParser) int {
	sched.m.Lock()
	defer sched.m.Unlock()
	if sched.pending == nil {
		sched.pending = make([]int, len(config.Servers))
	}
	httpparser.url = []byte("/main/time?period=500")
	mn, mnid := math.MaxInt32, 0
	for id, p := range sched.pending {
		if p < mn {
			mn = p
			mnid = id
		}
	}
	sched.pending[mnid]++
	return mnid
}
func (sched *LeastLoadBad) Finish(config *LBConfig, httpparser *HttpParser, srvid int, dur int64) {
	sched.m.Lock()
	defer sched.m.Unlock()
	sched.pending[srvid]--
}
