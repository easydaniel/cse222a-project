package main

import (
	"math"
	"sync"
)

type Status struct {
	lat int64
	cnt int
}

type UrlBase struct {
	m       sync.Mutex
	id      int
	uri2srv map[string]map[int]*Status
}

func (sched *UrlBase) Schedule(config *LBConfig, httpparser *HttpParser) int {
	sched.m.Lock()
	defer sched.m.Unlock()
	if sched.uri2srv == nil {
		sched.uri2srv = make(map[string]map[int]*Status)
	}
	url := string(httpparser.url)
	if _, ok := sched.uri2srv[url]; !ok {
		sched.uri2srv[url] = make(map[int]*Status, len(config.Servers))
		for i := 0; i < len(config.Servers); i++ {
			sched.uri2srv[url][i] = &Status{cnt: 1}
		}
	}
	srvid, mn := 0, math.MaxFloat64
	for i := 0; i < len(config.Servers); i++ {
		if mn > float64(sched.uri2srv[url][i].lat)/float64(sched.uri2srv[url][i].cnt) {
			mn = float64(sched.uri2srv[url][i].lat) / float64(sched.uri2srv[url][i].cnt)
			srvid = i
		}
	}
	return srvid
}
func (sched *UrlBase) Finish(config *LBConfig, httpparser *HttpParser, srvid int, dur int64) {
	sched.m.Lock()
	defer sched.m.Unlock()
	url := string(httpparser.url)
	sched.uri2srv[url][srvid].cnt++
	sched.uri2srv[url][srvid].lat = sched.uri2srv[url][srvid].lat + dur
}
