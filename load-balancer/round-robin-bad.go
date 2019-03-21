package main

import "sync"

type RoundRobinBad struct {
	m  sync.Mutex
	id int
}

func (sched *RoundRobinBad) Schedule(config *LBConfig, httpparser *HttpParser) int {
	sched.m.Lock()
	defer sched.m.Unlock()
	srvid := sched.id
	sched.id = (sched.id + 1) % len(config.Servers)
	if srvid < len(config.Servers)/2 {
		httpparser.url = []byte("url1")
	} else {
		httpparser.url = []byte("url2")
	}
	return srvid
}
func (sched *RoundRobinBad) Finish(config *LBConfig, httpparser *HttpParser, srvid int, dur int64) {}
