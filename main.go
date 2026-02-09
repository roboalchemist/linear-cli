package main

import (
	_ "embed"

	"github.com/roboalchemist/linear-cli/cmd"
)

//go:embed README.md
var readmeContents string

func main() {
	cmd.SetReadmeContents(readmeContents)
	cmd.Execute()
}
