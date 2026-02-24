package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/crit/gif2vid/internal/app"
	"github.com/crit/gif2vid/internal/config"
)

func main() {
	fs := flag.NewFlagSet("gif2vid", flag.ExitOnError)
	cfg := config.AddFlags(fs)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <inputs...>\n", os.Args[0])
		fs.PrintDefaults()
	}
	_ = fs.Parse(os.Args[1:])
	if err := cfg.Finalize(fs.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fs.Usage()
		os.Exit(2)
	}

	ctx := context.Background()
	if err := app.Run(ctx, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "gif2vid: %v\n", err)
		os.Exit(1)
	}
}
