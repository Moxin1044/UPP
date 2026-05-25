package probe

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type HTTPProbe struct{}

func (p *HTTPProbe) Type() string {
	return "http"
}

func (p *HTTPProbe) Probe(target string, timeout time.Duration) Result {
	start := time.Now()
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}
	resp, err := client.Get(target)
	latency := float64(time.Since(start).Milliseconds())
	if err != nil {
		return Result{Status: "down", Latency: latency, Message: err.Error()}
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return Result{Status: "down", Latency: latency, Message: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}
	return Result{Status: "up", Latency: latency, Message: fmt.Sprintf("HTTP %d", resp.StatusCode)}
}
