package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	"github.com/kr/pretty"
)

const (
	EPOLLET        = 1 << 31
	MaxEpollEvents = 1 << 6
)

// var CRFL = []byte{'\r', '\n', '\r', '\n'}

var config LBConfig
var mm sync.Mutex
var num = 0

func handle(clifd int, config *LBConfig, sched Scheduler) {
	defer syscall.Close(clifd)
	buf := make([]byte, 1024)
	rbuf := make([]byte, 0, 1024)
	for {
		n, _ := syscall.Read(clifd, buf[:])
		rbuf = append(rbuf, buf[:n]...)
		if bytes.Equal(rbuf[len(rbuf)-4:], append(CRFL, CRFL...)) {
			break
		}
	}
	httpparser := NewHttpParser(rbuf, false)
	srvid := sched.Schedule(config, httpparser)

	httpparser.headers = append(httpparser.headers,
		Header{name: []byte("Host"), value: []byte(config.Servers[srvid])[:len(config.Servers[srvid])-3]})
	httpparser.headers = append(httpparser.headers,
		Header{name: []byte("Connection"), value: []byte("close")})
	t := time.Now()
	srvconn, err := net.Dial("tcp", config.Servers[srvid])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srvconn.Close()
	// rbuf = rbuf[:0]
	rbuf = make([]byte, 0, len(rbuf))
	httpparser.Serealize(&rbuf)
	srvconn.Write(rbuf)
	for {
		n, err := srvconn.Read(buf[:])
		if err != nil || n <= 0 {
			break
		}
		syscall.Write(clifd, buf[:n])
	}
	dur := time.Since(t)
	defer sched.Finish(config, httpparser, srvid, int64(dur))
}

func main() {
	// sched := &UrlBase{}
	// sched := &RoundRobin{}
	// sched := &RoundRobinBad{}
	sched := &LeastLoad{}
	// sched := &LeastTime{}
	// sched := &WeightedRoundRobin{}
	// sched := &LeastLoadBad{}
	config, err := loadLBConfig("config.json")
	if err != nil {
		fmt.Println(err)
	}
	pretty.Println(config)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err = syscall.SetNonblock(fd, true); err != nil {
		fmt.Println("setnonblock1: ", err)
		os.Exit(1)
	}
	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		fmt.Println("setnonblock1: ", err)
		os.Exit(1)
	}
	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		fmt.Println("setnonblock1: ", err)
		os.Exit(1)
	}
	addr := syscall.SockaddrInet4{Port: config.Port}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	syscall.Bind(fd, &addr)
	syscall.Listen(fd, config.Listen)
	epfd, e := syscall.EpollCreate1(0)
	if e != nil {
		fmt.Println("epoll_create1: ", e)
		os.Exit(1)
	}
	defer syscall.Close(epfd)
	var event syscall.EpollEvent
	var events [MaxEpollEvents]syscall.EpollEvent
	fmt.Println(fd)
	event.Events = syscall.EPOLLIN
	event.Fd = int32(fd)
	if e = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, &event); e != nil {
		fmt.Println("epoll_ctl: ", e)
		os.Exit(1)
	}
	for {
		nevents, e := syscall.EpollWait(epfd, events[:], -1)
		if e != nil {
			fmt.Println("epoll_wait: ", e)
			break
		}
		for ev := 0; ev < nevents; ev++ {
			if int(events[ev].Fd) == fd {
				connFd, _, err := syscall.Accept(fd)
				if err != nil {
					fmt.Println("accept: ", err)
					continue
				}
				syscall.SetNonblock(fd, true)
				event.Events = syscall.EPOLLIN | EPOLLET
				event.Fd = int32(connFd)
				if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, connFd, &event); err != nil {
					fmt.Print("epoll_ctl: ", connFd, err)
					os.Exit(1)
				}
			} else {
				go handle(int(events[ev].Fd), config, sched)
				if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_DEL, int(events[ev].Fd), &event); err != nil {
					fmt.Print("epoll_ctl: ", int(events[ev].Fd), err)
					os.Exit(1)
				}
			}
		}
	}
}
