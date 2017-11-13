package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr      = flag.String("listen", "127.0.0.1:8080", "address to listen on for HTTP requests.")
	hostsFile = flag.String("hosts", "./hosts.yaml", "file contains hosts to be tested")
	interval  = flag.Uint("interval", 10, "test interval")
)

func main() {
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	go perf()
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func perf() {
	f, err := ioutil.ReadFile(*hostsFile)
	if err != nil {
		log.Fatalf("failed to open hosts file, %s\n", err)
	}

	var hosts Hosts
	if err := yaml.Unmarshal(f, &hosts); err != nil {
		log.Fatalf("failed to unmarshal hosts file, %s", err)
	}

	testInterval := time.Duration(*interval) * time.Second
	for _, h := range hosts.Hosts {
		go perfHost(h, testInterval)
	}
}

func perfHost(h *Host, interval time.Duration) {
	t := time.NewTicker(interval)
	timeout := interval / 2
	for {
		select {
		case <-t.C:
			if h.HTTPEndpoint != "" {
				go perfHTTP(h.Name, h.HTTPEndpoint, timeout)
			}
		}
	}
}

func perfHTTP(name string, endpoint string, timeout time.Duration) {
	log.Printf("perf http: %s, %s", name, endpoint)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("invalid http url, %s, %s\n", endpoint, err)
		return
	}

	var t0, t1, t2, t3, t4 time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { t1 = time.Now() },
		ConnectStart: func(_, _ string) {
			if t1.IsZero() {
				// connecting to IP
				t1 = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				log.Printf("unable to connect to host %v: %v\n", addr, err)
			}
			t2 = time.Now()

			//log.Printf("\n%s%s\n", color.GreenString("Connected to "), color.CyanString(addr))
		},
		GotConn:              func(_ httptrace.GotConnInfo) { t3 = time.Now() },
		GotFirstResponseByte: func() { t4 = time.Now() },
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(httptrace.WithClientTrace(ctx, trace))

	tr := &http.Transport{}
	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to read response: %s, %v", endpoint, err)
		return
	}

	// readbody

	resp.Body.Close()
	//t5 := time.Now() // after read body
	TCPConnectTime.WithLabelValues(name).Set(float64(t1.Sub(t0)))

}
