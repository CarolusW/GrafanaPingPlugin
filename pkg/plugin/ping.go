package plugin

import (
	"context"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"os/exec"
)

var timeRe = regexp.MustCompile(`time[=<]([\d.]+)\s*ms`)

// Ping runs a single ICMP echo request against host using the system ping
// binary and returns the round-trip time in milliseconds.
func Ping(host string, timeout time.Duration) (float64, error) {
	if host == "" {
		return 0, fmt.Errorf("host is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		args := []string{"-n", "1", "-w", strconv.Itoa(int(timeout.Milliseconds())), host}
		cmd = exec.CommandContext(ctx, "ping", args...)
	default:
		args := []string{"-c", "1", "-W", strconv.Itoa(int(mathMax(1, int(timeout.Seconds())))), host}
		cmd = exec.CommandContext(ctx, "ping", args...)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("ping failed: %v: %s", err, truncate(string(out), 200))
	}

	return parseLatency(string(out))
}

func parseLatency(output string) (float64, error) {
	m := timeRe.FindStringSubmatch(output)
	if len(m) < 2 {
		return 0, fmt.Errorf("could not parse ping output: %s", truncate(output, 200))
	}

	ms, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse latency %q: %v", m[1], err)
	}
	return ms, nil
}

func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "..."
	}
	return s
}

func mathMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
