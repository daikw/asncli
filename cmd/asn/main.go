package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/michalvavra/asncli/internal/asana"
	"github.com/michalvavra/asncli/internal/auth"
	"github.com/michalvavra/asncli/internal/cli"
	"github.com/michalvavra/asncli/internal/cli/cmd"
	"github.com/michalvavra/asncli/internal/config"
)

func main() {
	exitCode := run(os.Args, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}

func run(args []string, stdout, stderr *os.File) int {
	root := cmd.Root{}
	parser, err := kong.New(&root,
		kong.Name("asn"),
		kong.Description("Asana CLI"),
		kong.UsageOnError(),
		kong.BindTo(context.Background(), (*context.Context)(nil)),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
	)
	if err != nil {
		fmt.Fprintf(stderr, "asn: %s\n", err.Error())
		return 1
	}
	ctx, err := parser.Parse(args[1:])
	if err != nil {
		fmt.Fprintf(stderr, "asn: %s\nTry 'asn --help'\n", err.Error())
		return 2
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(stderr, "asn: failed to load config: %s\n", err.Error())
		return 1
	}

	rootCtx := cli.Context{
		Stdout: stdout,
		Stderr: stderr,
		JSON:   root.JSON,
		TokenSource: auth.NewTokenSource(auth.TokenSourceOptions{
			Service: auth.DefaultService,
			User:    auth.DefaultUser,
		}),
		Store: auth.NewKeyringStore(),
		ClientFactory: func(source auth.TokenSource) *asana.Client {
			return asana.NewClient(source, nil)
		},
		Renderer: cli.NewRenderer(stdout, stderr, root.JSON),
		Config:   cfg,
	}

	if err := ctx.Run(&rootCtx); err != nil {
		fmt.Fprintf(stderr, "asn: %s\n", err.Error())
		return 1
	}
	return 0
}
