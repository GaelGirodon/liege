package console

import (
	"flag"
	"gaelgirodon.fr/liege/internal/model"
	"os"
	"strings"
	"testing"
)

func Test_Parse(t *testing.T) {
	files := []string{"../../go.mod", "../../go.sum"}
	type env struct {
		root    string
		port    string
		cert    string
		key     string
		latency string
	}
	tests := []struct {
		name    string
		args    []string
		env     env
		want    model.Config
		wantErr bool
	}{
		{name: "ok/cli-min", args: []string{"l", ".."}, env: env{},
			want: model.Config{Root: "..", Port: 3000, Latency: model.Latency{Min: 0, Max: 0}}},
		{name: "ok/cli-all", args: []string{"l", "-p=3001", "-c=" + files[0], "-k=" + files[1], "-l=5", ".."}, env: env{},
			want: model.Config{Root: "..", Port: 3001, Cert: files[0], Key: files[1], Latency: model.Latency{Min: 5, Max: 5}}},
		{name: "ok/env-all", args: []string{"l"}, env: env{root: "..", port: "3001", cert: files[0], key: files[1], latency: "4"},
			want: model.Config{Root: "..", Port: 3001, Cert: files[0], Key: files[1], Latency: model.Latency{Min: 4, Max: 4}}},
		{name: "ok/cli-env", args: []string{"l", "-p=3002", "-c=" + files[0], "-k=" + files[1], "-l=5-6", "./"},
			env:  env{root: "..", port: "3001", cert: files[1], key: files[0], latency: "4"},
			want: model.Config{Root: "./", Port: 3002, Cert: files[0], Key: files[1], Latency: model.Latency{Min: 5, Max: 6}}},
		{name: "err/root-missing", args: []string{"l"}, env: env{}, want: model.Config{}, wantErr: true},
		{name: "err/root-not-found", args: []string{"l", "nowhere"}, env: env{}, want: model.Config{}, wantErr: true},
		{name: "err/root-not-dir", args: []string{"l", "cli.go"}, env: env{}, want: model.Config{}, wantErr: true},
		{name: "err/port-number", args: []string{"l", "-p=99999", ".."}, env: env{}, want: model.Config{}, wantErr: true},
		{name: "err/cert-no-key", args: []string{"l", "-c=" + files[0], ".."}, env: env{}, want: model.Config{}, wantErr: true},
		{name: "err/key-no-cert", args: []string{"l", "-k=" + files[1], ".."}, env: env{}, want: model.Config{}, wantErr: true},
		{name: "err/bad-cert", args: []string{"l", "-c=bad", "-k=" + files[1], ".."}, env: env{}, want: model.Config{}, wantErr: true},
		{name: "err/bad-key", args: []string{"l", "-c=" + files[0], "-k=bad", ".."}, env: env{}, want: model.Config{}, wantErr: true},
		{name: "err/latency", args: []string{"l", "-l=999999", ".."}, env: env{}, want: model.Config{}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set command-line arguments and environment variables
			os.Args = test.args
			_ = os.Setenv(RootEnvVar, test.env.root)
			_ = os.Setenv(PortEnvVar, test.env.port)
			_ = os.Setenv(LatencyEnvVar, test.env.latency)
			// Reset flags configuration
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			// Run
			args, err := ParseArgs()
			// Assert
			if test.wantErr != (err != nil) {
				t.Fatalf("want error %v, got %v (%v)", test.wantErr, err != nil, err)
			}
			if err != nil {
				return
			}
			if args.Root != test.want.Root {
				t.Errorf("want root = %v, got %v", test.want.Root, args.Root)
			}
			if args.Port != test.want.Port {
				t.Errorf("want port = %v, got %v", test.want.Port, args.Port)
			}
			if args.Latency != test.want.Latency {
				t.Errorf("want latency = %v, got %v", test.want.Latency, args.Latency)
			}
		})
	}
}

func Test_ValidateRootDirPath(t *testing.T) {
	tests := []struct {
		name    string
		root    string
		wantErr string
	}{
		{name: "ok", root: "../../internal", wantErr: ""},
		{name: "err/empty", root: "", wantErr: "required"},
		{name: "err/not-found", root: "unknown", wantErr: "exist"},
		{name: "err/not-dir", root: "console.go", wantErr: "directory"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateRootDirPath(test.root)
			if (len(test.wantErr) != 0) != (err != nil) ||
				err != nil && !strings.Contains(err.Error(), test.wantErr) {
				t.Errorf("want error '%v', got '%v'", test.wantErr, err)
			}
		})
	}
}
