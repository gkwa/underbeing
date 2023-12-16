package options

import "flag"

type Options struct {
	Verbose    bool
	LogFormat  string
	GithubUser string
}

// ParseOptions parses the command-line options and returns an Options struct
func ParseOptions() *Options {
	opts := &Options{}

	flag.BoolVar(&opts.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&opts.Verbose, "v", false, "Enable verbose output (shorthand)")
	flag.StringVar(&opts.LogFormat, "log-format", "", "Log format (text or json)")

	flag.StringVar(&opts.GithubUser, "github-user", "", "GitHub username (overrides GITHUB_USER environment variable)")

	flag.Parse()

	return opts
}
