package probe

import "time"

type HTTPSProbe struct{}

func (p *HTTPSProbe) Type() string {
	return "https"
}

func (p *HTTPSProbe) Probe(target string, timeout time.Duration) Result {
	// HTTPS probe is the same as HTTP probe, just with TLS verification
	return (&HTTPProbe{}).Probe(target, timeout)
}
