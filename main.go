package main

import "github.com/YaleUniversity/deco/cmd"

var (
	version = "v0.0.0"
	date    = "unset"
	commit  = "unset"
)

func main() {
	cmd.Version = &cmd.CmdVersion{
		AppVersion: version,
		BuildTime:  date,
		GitCommit:  commit,
	}

	cmd.Execute()
}
