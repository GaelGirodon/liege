package server

import (
	"gaelgirodon.fr/liege/internal/model"
	"reflect"
	"testing"
)

func Test_parseFileName(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		wantName  string
		wantExt   string
		wantRoute model.Route
		wantErr   bool
	}{
		{"empty", "", "", "",
			model.Route{Code: 200, Latency: model.Latency{Min: -1, Max: -1}}, false},
		{"name", "test", "test", "",
			model.Route{Code: 200, Latency: model.Latency{Min: -1, Max: -1}}, false},
		{"name/empty-opt", "test___", "test", "",
			model.Route{Code: 200, Latency: model.Latency{Min: -1, Max: -1}}, false},
		{"ext", "test.json", "test", ".json",
			model.Route{Code: 200, Latency: model.Latency{Min: -1, Max: -1}}, false},
		{"method", "test__GET", "test", "",
			model.Route{Method: "GET", Code: 200, Latency: model.Latency{Min: -1, Max: -1}}, false},
		{"method/err", "test__ERR", "test", "",
			model.Route{}, true},
		{"query/single/name", "test__qn", "test", "",
			model.Route{QueryParams: []model.QueryParam{{Name: "n"}}, Code: 200, Latency: model.Latency{Min: -1, Max: -1}}, false},
		{"query/single/name-value", "test__qn=v", "test", "",
			model.Route{QueryParams: []model.QueryParam{{Name: "n", Value: "v"}}, Code: 200, Latency: model.Latency{Min: -1, Max: -1}}, false},
		{"query/multiple", "test__qn_qs", "test", "",
			model.Route{QueryParams: []model.QueryParam{{Name: "n"}, {Name: "s"}}, Code: 200, Latency: model.Latency{Min: -1, Max: -1}}, false},
		{"query/err", "test__qa=b=c", "test", "",
			model.Route{}, true},
		{"code", "test__500", "test", "",
			model.Route{Code: 500, Latency: model.Latency{Min: -1, Max: -1}}, false},
		{"code/err", "test__999", "test", "",
			model.Route{}, true},
		{"latency/fixed", "test__l20", "test", "",
			model.Route{Code: 200, Latency: model.Latency{Min: 20, Max: 20}}, false},
		{"latency/random", "test__l10-30", "test", "",
			model.Route{Code: 200, Latency: model.Latency{Min: 10, Max: 30}}, false},
		{"latency/err", "test__l999999", "test", "",
			model.Route{}, true},
		{"all", "test__POST_qn_403_l50.txt", "test", ".txt",
			model.Route{Method: "POST", QueryParams: []model.QueryParam{{Name: "n"}}, Code: 403, Latency: model.Latency{Min: 50, Max: 50}}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, ext, route, err := parseFileName(test.filename)
			if test.wantErr != (err != nil) {
				t.Errorf("want error = %v, got %v (%v)", test.wantErr, err != nil, err)
			}
			if err != nil {
				return
			}
			if name != test.wantName {
				t.Errorf("want name = %v, got %v", test.wantName, name)
			}
			if ext != test.wantExt {
				t.Errorf("want ext = %v, got %v", test.wantExt, ext)
			}
			if route.Method != test.wantRoute.Method {
				t.Errorf("want method = %v, got %v", test.wantRoute.Method, route.Method)
			}
			if test.wantRoute.QueryParams == nil {
				test.wantRoute.QueryParams = []model.QueryParam{}
			}
			if !reflect.DeepEqual(route.QueryParams, test.wantRoute.QueryParams) {
				t.Errorf("want queryParams = %v, got %v", test.wantRoute.QueryParams, route.QueryParams)
			}
			if route.Code != test.wantRoute.Code {
				t.Errorf("want code = %v, got %v", test.wantRoute.Code, route.Code)
			}
			if route.Latency != test.wantRoute.Latency {
				t.Errorf("want latency = %v, got %v", test.wantRoute.Latency, route.Latency)
			}
		})
	}
}
