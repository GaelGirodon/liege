package model

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Latency is the response latency to simulate.
type Latency struct {
	// Min is the minimum latency value in ms.
	Min int `json:"min"`
	// Max is the maximum latency value in ms.
	Max int `json:"max"`
}

// latencyPattern is the pattern to validate and parse a latency value.
var latencyPattern = regexp.MustCompile("^(?:-1|([0-9]{1,5})(?:-([0-9]{1,5}))?)$")

// ParseLatency validates, parses and returns a latency value.
func ParseLatency(value string, prefix string) (latency Latency, err error) {
	if !strings.HasPrefix(value, prefix) {
		return Latency{}, errors.New("invalid latency value")
	}
	if match := latencyPattern.FindStringSubmatch(value[len(prefix):]); len(match) == 3 {
		min, _ := strconv.Atoi(match[1])
		max := min
		if len(match[2]) > 0 {
			max, _ = strconv.Atoi(match[2])
		}
		latency = Latency{Min: min, Max: max}
		if latency.IsValid() {
			return
		}
	}
	return Latency{}, errors.New("invalid latency value")
}

// IsValid indicates whether the current latency is valid or not.
func (l Latency) IsValid() bool {
	return l.Min == -1 && l.Max == -1 || // Disabled or undefined
		l.Min >= 0 && l.Max >= l.Min && l.Max < 100000
}

// IsDisabledOrUndefined indicates whether the current latency has a value
// meaning that latency must be disabled (CLI flag) or is undefined (file name).
func (l Latency) IsDisabledOrUndefined() bool {
	return l.Min == -1
}

// Compute computes the duration to wait before sending the response.
func (l Latency) Compute(global Latency) time.Duration {
	lat := global // Take the global value by default
	if global.IsDisabledOrUndefined() {
		lat = Latency{Min: 0, Max: 0} // Latency disabled globally
	} else if !l.IsDisabledOrUndefined() {
		lat = l // Latency from file name overrides global latency
	}
	var duration int
	if lat.Min == lat.Max {
		duration = lat.Min // Fixed value
	} else {
		duration = lat.Min + rand.Intn(lat.Max+1-lat.Min) // Random in range
	}
	return time.Duration(duration) * time.Millisecond
}
