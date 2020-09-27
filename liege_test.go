package main

import (
	"flag"
	"os"
	"testing"
)

func Test_parseArgs(t *testing.T) {
	type env struct {
		root    string
		port    string
		verbose string
	}
	type want struct {
		root     string
		port     uint
		verbose  bool
		err      bool
		exitCode int
	}
	tests := []struct {
		name string
		args []string
		env  env
		want want
	}{
		{name: "ok/cli-min", args: []string{"l", "data"}, env: env{}, want: want{root: "data", port: 3000, verbose: false}},
		{name: "ok/cli-all", args: []string{"l", "-p=3001", "-v", "data"}, env: env{}, want: want{root: "data", port: 3001, verbose: true}},
		{name: "ok/env-all", args: []string{"l"}, env: env{root: "data", port: "3001", verbose: "1"}, want: want{root: "data", port: 3001, verbose: true}},
		{name: "ok/cli-env", args: []string{"l", "-p=3002", "-v=false", "./"}, env: env{root: "data", port: "3001", verbose: "1"}, want: want{root: "./", port: 3002, verbose: false}},
		{name: "err/port-number", args: []string{"l", "-p=99999", "data"}, env: env{}, want: want{err: true, exitCode: 2}},
		{name: "err/root-missing", args: []string{"l"}, env: env{}, want: want{err: true, exitCode: 3}},
		{name: "err/root-not-found", args: []string{"l", "nowhere"}, env: env{}, want: want{err: true, exitCode: 4}},
		{name: "err/root-not-dir", args: []string{"l", "README.md"}, env: env{}, want: want{err: true, exitCode: 5}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set command-line arguments and environment variables
			os.Args = test.args
			_ = os.Setenv(rootEnvVar, test.env.root)
			_ = os.Setenv(portEnvVar, test.env.port)
			_ = os.Setenv(verboseEnvVar, test.env.verbose)
			// Reset flags configuration
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			// Run
			root, port, verbose, err, exitCode := parseArgs()
			// Assert
			if test.want.err != (err != nil) {
				t.Fatalf("want error %v, got %v", test.want.err, err != nil)
			}
			if exitCode != test.want.exitCode {
				t.Errorf("want exitCode = %v, got %v", test.want.exitCode, exitCode)
			}
			if root != test.want.root {
				t.Errorf("want root = %v, got %v", test.want.root, root)
			}
			if port != test.want.port {
				t.Errorf("want port = %v, got %v", test.want.port, port)
			}
			if verbose != test.want.verbose {
				t.Errorf("want verbose = %v, got %v", test.want.verbose, verbose)
			}
		})
	}
}

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
