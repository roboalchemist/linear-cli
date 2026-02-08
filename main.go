package main

import (
	_ "embed"

	"github.com/dorkitude/linear-cli/cmd"
)

//go:embed README.md
var readmeContents string

func main() {
	cmd.SetReadmeContents(readmeContents)
	cmd.Execute()
}
