package server

import "testing"

func Test_parseFileName(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantName string
		wantExt  string
		wantCode int
	}{
		{"empty", "", "", "", 200},
		{"no-ext", "test", "test", "", 200},
		{"ext", "test.json", "test", ".json", 200},
		{"code", "test__500", "test", "", 500},
		{"ext-code", "test__500.txt", "test", ".txt", 500},
		{"bad-code", "test__999", "test__999", "", 200},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, ext, code := parseFileName(test.filename)
			if name != test.wantName {
				t.Errorf("want name = %v, got %v", test.wantName, name)
			}
			if ext != test.wantExt {
				t.Errorf("want ext = %v, got %v", test.wantExt, ext)
			}
			if code != test.wantCode {
				t.Errorf("want code = %v, got %v", test.wantCode, code)
			}
		})
	}
}
