package main

type Scheduler interface {
	Schedule(*LBConfig, *HttpParser) int
	Finish(*LBConfig, *HttpParser, int, int64)
}
