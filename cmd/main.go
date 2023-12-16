package main

import (
	"log/slog"
	"os"

	"github.com/taylormonacelli/goldbug"
	"github.com/taylormonacelli/underbeing"
	"github.com/taylormonacelli/underbeing/options"
)

func main() {
	opts := options.ParseOptions()

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
