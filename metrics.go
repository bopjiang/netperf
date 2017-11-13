package main

import "github.com/prometheus/client_golang/prometheus"

var (
	TCPConnectTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "netperf",
			Name:      "tcp_conn_time",
			Help:      "connection time(unit, ns)",
		}, []string{"host"})
)

func init() {
	prometheus.MustRegister(TCPConnectTime)
}
