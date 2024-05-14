package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/taylormonacelli/goldbug"
	"github.com/taylormonacelli/underbeing"
	optmod "github.com/taylormonacelli/underbeing/options"
	"github.com/taylormonacelli/underbeing/version"
)

func main() {
	opts := optmod.ParseOptions()

	if opts.Version {
		buildInfo := version.GetBuildInfo()
		fmt.Println(buildInfo)
		os.Exit(0)
	}

	if opts.Verbose || opts.LogFormat != "" {
		if opts.LogFormat == "json" {
			goldbug.SetDefaultLoggerJson(slog.LevelDebug)
		} else {
			goldbug.SetDefaultLoggerText(slog.LevelDebug)
		}
	}

	code := underbeing.Main(opts)
	os.Exit(code)
}
