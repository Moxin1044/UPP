package probe

import (
	"fmt"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type PingProbe struct{}

func (p *PingProbe) Type() string {
	return "ping"
}

func (p *PingProbe) Probe(target string, timeout time.Duration) Result {
	pinger, err := probing.NewPinger(target)
	if err != nil {
		return Result{Status: "down", Latency: 0, Message: err.Error()}
	}
	pinger.Count = 3
	pinger.Timeout = timeout
	err = pinger.Run()
	if err != nil {
		return Result{Status: "down", Latency: 0, Message: err.Error()}
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		return Result{Status: "down", Latency: float64(stats.AvgRtt.Milliseconds()), Message: "no response"}
	}
	return Result{Status: "up", Latency: float64(stats.AvgRtt.Milliseconds()), Message: fmt.Sprintf("%d/%d packets", stats.PacketsRecv, stats.PacketsSent)}
}
