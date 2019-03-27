package main

import (
	cmd "github.com/cnative/servicebuilder/cmd"
)

var (
	version = "unknown"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	cmd.Execute(version, commit, date)
}
