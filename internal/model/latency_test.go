package model

import (
	"testing"
	"time"
)

func TestLatency_Compute(t *testing.T) {
	tests := []struct {
		name   string
		global Latency
		local  Latency
		want   int
	}{
		{"-1/-1", Latency{-1, -1}, Latency{-1, -1}, 0},
		{"-1/0", Latency{-1, -1}, Latency{0, 0}, 0},
		{"-1/n", Latency{-1, -1}, Latency{5, 5}, 0},
		{"0/-1", Latency{0, 0}, Latency{-1, -1}, 0},
		{"0/0", Latency{0, 0}, Latency{0, 0}, 0},
		{"0/n", Latency{0, 0}, Latency{5, 5}, 5},
		{"n/-1", Latency{4, 4}, Latency{-1, -1}, 4},
		{"n/0", Latency{4, 4}, Latency{0, 0}, 0},
		{"n/n", Latency{4, 4}, Latency{5, 5}, 5},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.local.Compute(test.global)
			expected := time.Duration(test.want) * time.Millisecond
			if actual != time.Duration(test.want)*time.Millisecond {
				t.Errorf("want Compute() = %d, got %d", expected, actual)
			}
		})
	}
}
