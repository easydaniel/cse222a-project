package main

import "sync"

type RoundRobin struct {
	m  sync.Mutex
	id int
}

func (sched *RoundRobin) Schedule(config *LBConfig, httpparser *HttpParser) int {
	sched.m.Lock()
	defer sched.m.Unlock()
	srvid := sched.id
	sched.id = (sched.id + 1) % len(config.Servers)
	return srvid
}
func (sched *RoundRobin) Finish(config *LBConfig, httpparser *HttpParser, srvid int, dur int64) {}
