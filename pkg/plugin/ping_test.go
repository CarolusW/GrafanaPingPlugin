package plugin

import (
	"testing"
)

func TestParseLatency(t *testing.T) {
	cases := map[string]float64{
		"64 bytes from 8.8.8.8: icmp_seq=1 ttl=118 time=12.3 ms":     12.3,
		"64 bytes from 8.8.8.8: icmp_seq=1 ttl=118 time=1.234 ms":    1.234,
		"Reply from 8.8.8.8: bytes=32 time=13ms TTL=118":             13,
		"64 bytes from 1.1.1.1: icmp_seq=0 ttl=57 time=0.078 ms":    0.078,
		"64 bytes from 8.8.8.8: icmp_seq=1 ttl=118 time<1 ms":       1,
		"64 bytes from example.com (1.2.3.4): time=99.999 ms ttl=64": 99.999,
	}
	for out, want := range cases {
		got, err := parseLatency(out)
		if err != nil {
			t.Errorf("parseLatency(%q) error: %v", out, err)
			continue
		}
		if got != want {
			t.Errorf("parseLatency(%q) = %v, want %v", out, got, want)
		}
	}

	if _, err := parseLatency("ping: unknown host"); err == nil {
		t.Error("expected error for output without latency")
	}
}
