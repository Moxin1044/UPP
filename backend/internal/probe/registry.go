package probe

var registry = map[string]Probe{}

func init() {
	Register(&HTTPProbe{})
	Register(&HTTPSProbe{})
	Register(&TCPProbe{})
	Register(&PingProbe{})
}

func Register(p Probe) {
	registry[p.Type()] = p
}

func Get(probeType string) (Probe, bool) {
	p, ok := registry[probeType]
	return p, ok
}
