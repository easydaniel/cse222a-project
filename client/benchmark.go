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

func request(url string, reqCnt int, conType int, message chan Status) {
	defer wg.Done()

	for idx := 0; idx < reqCnt; idx++ {
		qstr := qs(genParam(conType))
		endpoint := fmt.Sprintf("%s?%s", url, qstr)
		start := time.Now()
		resp, err := http.Get(endpoint)
		latency := time.Since(start)
		code := -1
		if err != nil {
			fmt.Println(err)
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

func main() {
	var numReq int
	flag.IntVar(&numReq, "n", 1000, "# of requests")
	var numCon int
	flag.IntVar(&numCon, "c", 10, "# of connections")
	var conType int
	flag.IntVar(&conType, "t", 0, "type of connection: 1: short, 2: long, 3: mix")
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Printf("Usage: benchmark [options] <url>\n")
		return
	}
	message := make(chan Status)
	status := make(map[int][]float32)
	numEach := numReq / numCon
	for total := numReq; total > 0; total -= numEach {
		wg.Add(1)
		go request(flag.Arg(0), min(total, numEach), conType, message)
	}

	go func() {
		for result := range message {
			status[result.code] = append(status[result.code], (float32(result.latency) / 1000000))
		}
	}()
	wg.Wait()
	close(message)
	// fmt.Fprintf(os.Stderr, "Total: %d, Num: %d, Con: %d\n", len(status[200]), numReq, numCon)
	failRate := float32(len(status[-1])) / float32(numReq)
	avgLatency := float32(0)
	for _, latency := range status[200] {
		avgLatency = avgLatency + latency
	}
	fmt.Printf("%f %f\n", failRate, avgLatency/float32(len(status[200])))
}
