package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Status of network response
type Status struct {
	latency time.Duration
	code    int
}

type Request struct {
	url    string
	params map[string]string
}

var wg sync.WaitGroup

func qs(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	var s string
	for k, v := range params {
		s = fmt.Sprintf("%s%s=%s&", s, k, v)
	}
	return s[:len(s)-1]
}

func genParam(conType int) map[string]string {
	param := make(map[string]string)
	p := 0
	if conType == 3 {
		r := rand.Intn(2)
		if r == 1 {
			p = 500
		} else {
			p = 100
		}
	} else if conType == 2 {
		p = 500
	} else if conType == 1 {
		p = 100
	}
	param["period"] = strconv.Itoa(p)
	return param
}

func genEndpoint(conType int) string {
	s := ""
	return s
	if conType == 3 {
		r := rand.Intn(2)
		if r == 1 {
			s = "/ep1"
		} else {
			s = "/ep2"
		}
	} else if conType == 2 {
		s = "/ep2"
	} else if conType == 1 {
		s = "/ep1"
	}
	return s
}

func request(url string, reqCnt int, conType int, message chan Status) {
	defer wg.Done()

	for idx := 0; idx < reqCnt; idx++ {
		qstr := qs(genParam(conType))
		ep := genEndpoint(conType)
		endpoint := fmt.Sprintf("%s%s?%s", url, ep, qstr)
		start := time.Now()
		resp, err := http.Get(endpoint)
		latency := time.Since(start)
		code := -1
		if err != nil {
			// fmt.Println(err)
		} else {
			code = resp.StatusCode
			resp.Body.Close()
		}
		message <- Status{latency, code}
	}
}
func consumer(reqCnt int, works chan Request, message chan Status) {
	defer wg.Done()

	for idx := 0; idx < reqCnt; idx++ {
		req := <-works
		qstr := qs(req.params)
		endpoint := fmt.Sprintf("%s?%s", req.url, qstr)
		// fmt.Println(endpoint)
		start := time.Now()
		resp, err := http.Get(endpoint)
		latency := time.Since(start)
		code := -1
		if err != nil {
		//	fmt.Println(err)
		} else {
			code = resp.StatusCode
			resp.Body.Close()
		}
		message <- Status{latency, code}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func reqlatency(numCon, numReq, conType int) {

	message := make(chan Status, numCon*numReq)
	status := make(map[int][]float32)
	url := flag.Arg(0)
	start := time.Now()
	for idx := 0; idx < numCon; idx++ {
		wg.Add(1)
		go request(url, numReq, conType, message)
	}

	go func() {
		for result := range message {
			status[result.code] = append(status[result.code], (float32(result.latency) / 1000000))
		}
	}()
	wg.Wait()
	duration := time.Since(start)
	close(message)
	// fmt.Fprintf(os.Stderr, "Total: %d, Num: %d, Con: %d\n", len(status[200]), numReq, numCon)
	failRate := float32(len(status[-1])) / (float32(numReq) * float32(numCon))
	throughput := (float32(numReq) * float32(numCon)) / (float32(duration) / 1000000000)
	avgLatency := float32(0)
	for _, latency := range status[200] {
		avgLatency = avgLatency + latency
	}
	fmt.Printf("%f %f %f\n", failRate, avgLatency/float32(len(status[200])), throughput)
}

func orderedReq(numCon, numReq, srvCnt int) {
	message := make(chan Status)
	works := make(chan Request, numCon*numReq)
	status := make(map[int][]float32)
	url := flag.Arg(0)
	start := time.Now()
	for idx := 0; idx < numCon*numReq; idx++ {
		params := make(map[string]string)
/*
		if idx%srvCnt == 0 {
			params["period"] = "500"
		} else {
			params["period"] = "100"
		}
*/
		params["period_s"] = "100"
		params["period_not_s"] = "500"
		if idx % 2 == 4 {
			params["server"] = "52.41.69.101"
		} else {
			params["server"] = "34.212.92.108"
		}
		works <- Request{url, params}
	}
	for idx := 0; idx < numCon; idx++ {
		wg.Add(1)
		go consumer(numReq, works, message)
	}

	go func() {
		for result := range message {
			status[result.code] = append(status[result.code], (float32(result.latency) / 1000000))
		}
	}()
	wg.Wait()
	duration := time.Since(start)
	close(message)
	// fmt.Fprintf(os.Stderr, "Total: %d, Num: %d, Con: %d\n", len(status[200]), numReq, numCon)
	failRate := float32(len(status[-1])) / (float32(numReq) * float32(numCon))
	throughput := (float32(numReq) * float32(numCon)) / (float32(duration) / 1000000000)
	avgLatency := float32(0)
	for _, latency := range status[200] {
		avgLatency = avgLatency + latency
	}
	fmt.Printf("%f %f %f\n", failRate, avgLatency/float32(len(status[200])), throughput)
}

func main() {
	var numReq int
	flag.IntVar(&numReq, "n", 30, "# of requests for each connection")
	var numCon int
	flag.IntVar(&numCon, "c", 10, "# of connections")
	var conType int
	flag.IntVar(&conType, "t", 0, "type of connection: 1: short, 2: long, 3: mix")
	var srvCnt int
	flag.IntVar(&srvCnt, "s", 1, "server count")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: benchmark [options] <url>\n")
		return
	}
	// reqlatency(numCon, numReq, conType)
	orderedReq(numCon, numReq, srvCnt)
}
