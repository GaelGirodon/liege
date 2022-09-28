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
		want   Latency
	}{
		{"-1/-1", Latency{-1, -1}, Latency{-1, -1}, Latency{0, 0}},
		{"-1/0", Latency{-1, -1}, Latency{0, 0}, Latency{0, 0}},
		{"-1/n", Latency{-1, -1}, Latency{5, 5}, Latency{0, 0}},
		{"0/-1", Latency{0, 0}, Latency{-1, -1}, Latency{0, 0}},
		{"0/0", Latency{0, 0}, Latency{0, 0}, Latency{0, 0}},
		{"0/n", Latency{0, 0}, Latency{5, 5}, Latency{5, 5}},
		{"n/-1", Latency{4, 4}, Latency{-1, -1}, Latency{4, 4}},
		{"n/0", Latency{4, 4}, Latency{0, 0}, Latency{0, 0}},
		{"n/n", Latency{4, 4}, Latency{5, 5}, Latency{5, 5}},
		{"0/rand", Latency{0, 0}, Latency{5, 10}, Latency{5, 10}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.local.Compute(test.global)
			wantMin := time.Duration(test.want.Min) * time.Millisecond
			wantMax := time.Duration(test.want.Max) * time.Millisecond
			if actual < wantMin || actual > wantMax {
				t.Errorf("want %d <= Compute() <= %d, got %d", wantMin, wantMax, actual)
			}
		})
	}
}
