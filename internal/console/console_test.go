package console

import (
	"flag"
	"os"
	"testing"
)

func Test_Parse(t *testing.T) {
	type env struct {
		root string
		port string
	}
	tests := []struct {
		name    string
		args    []string
		env     env
		want    Args
		wantErr bool
	}{
		{name: "ok/cli-min", args: []string{"l", ".."}, env: env{}, want: Args{Root: "..", Port: 3000}},
		{name: "ok/cli-all", args: []string{"l", "-p=3001", ".."}, env: env{}, want: Args{Root: "..", Port: 3001}},
		{name: "ok/env-all", args: []string{"l"}, env: env{root: "..", port: "3001"}, want: Args{Root: "..", Port: 3001}},
		{name: "ok/cli-env", args: []string{"l", "-p=3002", "./"}, env: env{root: "..", port: "3001"}, want: Args{Root: "./", Port: 3002}},
		{name: "err/port-number", args: []string{"l", "-p=99999", ".."}, env: env{}, want: Args{}, wantErr: true},
		{name: "err/root-missing", args: []string{"l"}, env: env{}, want: Args{}, wantErr: true},
		{name: "err/root-not-found", args: []string{"l", "nowhere"}, env: env{}, want: Args{}, wantErr: true},
		{name: "err/root-not-dir", args: []string{"l", "cli.go"}, env: env{}, want: Args{}, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set command-line arguments and environment variables
			os.Args = test.args
			_ = os.Setenv(RootEnvVar, test.env.root)
			_ = os.Setenv(PortEnvVar, test.env.port)
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
		})
	}
}
