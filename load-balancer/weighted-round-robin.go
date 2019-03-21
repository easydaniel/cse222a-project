package main

import "sync"

type WeightedRoundRobin struct {
	m    sync.Mutex
	pool []int
	id   int
}

func (sched *WeightedRoundRobin) Schedule(config *LBConfig, httpparser *HttpParser) int {
	sched.m.Lock()
	if sched.pool == nil {
		for i := 0; i < len(config.Servers); i++ {
			for j := 0; j < config.Weights[i]; j++ {
				sched.pool = append(sched.pool, i)
			}
		}
	}
	defer sched.m.Unlock()
	srvid := sched.pool[sched.id]
	sched.id = (sched.id + 1) % len(sched.pool)
	return srvid
}
func (sched *WeightedRoundRobin) Finish(config *LBConfig, httpparser *HttpParser, srvid int, dur int64) {
}
