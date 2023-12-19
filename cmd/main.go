package main

import (
	"log/slog"
	"os"

	"github.com/taylormonacelli/goldbug"
	"github.com/taylormonacelli/underbeing"
	optmod "github.com/taylormonacelli/underbeing/options"
)

func main() {
	opts := optmod.ParseOptions()

	slog.Info("hi")
	slog.Error("example fail")

	if opts.Verbose || opts.LogFormat != "" {
		if opts.LogFormat == "json" {
			goldbug.SetDefaultLoggerJson(slog.LevelDebug)
		} else {
			goldbug.SetDefaultLoggerText(slog.LevelDebug)
		}
	}

	slog.Info("hi")
	slog.Error("example fail")

	code := underbeing.Main(opts)
	os.Exit(code)
}
