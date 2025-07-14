package main

import (
	_ "embed"
	
	"github.com/dorkitude/linctl/cmd"
)

//go:embed README.md
var readmeContents string

func main() {
	cmd.SetReadmeContents(readmeContents)
	cmd.Execute()
}
