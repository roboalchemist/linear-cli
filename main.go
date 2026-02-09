package main

import (
	_ "embed"

	"github.com/roboalchemist/linear-cli/cmd"
)

//go:embed README.md
var readmeContents string

//go:embed .claude/skills/linear-cli/SKILL.md
var skillContents string

//go:embed .claude/skills/linear-cli/reference/commands.md
var skillRefContents string

func main() {
	cmd.SetReadmeContents(readmeContents)
	cmd.SetSkillContents(skillContents, skillRefContents)
	cmd.Execute()
}
