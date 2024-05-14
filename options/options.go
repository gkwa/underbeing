package options

import "flag"

type Options struct {
	Verbose    bool
	LogFormat  string
	GithubUser string
	RepoName   string
	Version    bool
}

func ParseOptions() *Options {
	opts := &Options{}

	flag.BoolVar(&opts.Verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&opts.Verbose, "v", false, "Enable verbose output (shorthand)")

	flag.BoolVar(&opts.Version, "version", false, "Display version")

	flag.StringVar(&opts.LogFormat, "log-format", "", "Log format (text or json)")

	flag.StringVar(&opts.GithubUser, "github-user", "", "GitHub username (overrides GITHUB_USER environment variable)")
	flag.StringVar(&opts.RepoName, "repo-name", "", "GitHub repository name (overrides the current directory name)")

	flag.Parse()

	return opts
}
