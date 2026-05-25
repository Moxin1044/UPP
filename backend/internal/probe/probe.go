package probe

import (
	"time"
)

type Result struct {
	Status  string  `json:"status"` // up, down
	Latency float64 `json:"latency"` // ms
	Message string  `json:"message"`
}

type Probe interface {
	Type() string
	Probe(target string, timeout time.Duration) Result
}
