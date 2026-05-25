package probe

import (
	"net"
	"time"
)

type TCPProbe struct{}

func (p *TCPProbe) Type() string {
	return "tcp"
}

func (p *TCPProbe) Probe(target string, timeout time.Duration) Result {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, timeout)
	latency := float64(time.Since(start).Milliseconds())
	if err != nil {
		return Result{Status: "down", Latency: latency, Message: err.Error()}
	}
	conn.Close()
	return Result{Status: "up", Latency: latency, Message: "connected"}
}
