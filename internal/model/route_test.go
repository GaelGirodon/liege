package model

import (
	"testing"
)

func TestRoute_Before(t *testing.T) {
	tests := []struct {
		name       string
		r          Route
		r2         Route
		wantBefore bool
	}{
		{"path/alpha", Route{Path: "a"}, Route{Path: "b"}, true},
		{"path/length", Route{Path: "test"}, Route{Path: "test2"}, true},
		{"method", Route{Method: "GET"}, Route{}, true},
		{"query/count", Route{QueryParams: []QueryParam{{}, {}}}, Route{}, true},
		{"query/specific", Route{QueryParams: []QueryParam{{"n", "v"}}},
			Route{QueryParams: []QueryParam{{"n", ""}}}, true},
		{"filepath", Route{FilePath: ""}, Route{FilePath: "longer"}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before := test.r.Before(test.r2)
			if before != test.wantBefore {
				t.Errorf("want before = %v, got %v", test.wantBefore, before)
			}
		})
	}
}
