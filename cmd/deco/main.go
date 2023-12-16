package main

import "github.com/YaleUniversity/deco/cli"

var (
	version = "v0.0.0"
	date    = "unset"
	commit  = "unset"
)

func main() {
	cli.Version = &cli.CmdVersion{
		AppVersion: version,
		BuildTime:  date,
		GitCommit:  commit,
	}

	cli.Execute()
}
