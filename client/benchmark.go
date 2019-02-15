package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

// Status of network response
type Status struct {
	latency time.Duration
	code    int
}

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

func request(url string, reqCnt int, params map[string]string, message chan Status) {

	qstr := qs(params)
	endpoint := fmt.Sprintf("%s?%s", url, qstr)
	for idx := 0; idx < reqCnt; idx++ {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var numReq int
	flag.IntVar(&numReq, "n", 10000, "# of requests")
	var numCon int
	flag.IntVar(&numCon, "c", 10, "# of connections")
	flag.Parse()
	message := make(chan Status)
	status := make(map[int][]float32)
	numEach := numReq / numCon
	for total := numReq; total > 0; total -= numEach {
		go request("http://allen.nctu.me:8388/", min(total, numEach), make(map[string]string), message)
	}
	for idx := 0; idx < numReq; idx++ {
		result := <-message
		status[result.code] = append(status[result.code], (float32(result.latency) / 1000000))
	}
	close(message)
	// fmt.Fprintf(os.Stderr, "Total: %d, Num: %d, Con: %d\n", len(status[200]), numReq, numCon)
	failRate := float32(len(status[-1])) / float32(numReq)
	avgLatency := float32(0)
	for _, latency := range status[200] {
		avgLatency = avgLatency + latency
	}
	fmt.Printf("%f %f\n", failRate, avgLatency/float32(len(status[200])))
}
