# linear-cli Command Reference

**Version**: 0.2.2 | **Binary**: `/opt/homebrew/bin/linear-cli`

## Table of Contents

- [Global Flags](#global-flags)
- [Issue Commands](#issue-commands)
- [Project Commands](#project-commands)
- [Cycle Commands](#cycle-commands)
- [Label Commands](#label-commands)
- [Team Commands](#team-commands)
- [User Commands](#user-commands)
- [Document Commands](#document-commands)
- [Initiative Commands](#initiative-commands)
- [View Commands](#view-commands)
- [GraphQL](#graphql)
- [Auth Commands](#auth-commands)
- [Utility Commands](#utility-commands)

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--config` | | Config file path (default: `~/.linear-cli.yaml`) |
| `--json` | `-j` | JSON output |
| `--plaintext` | `-p` | Plaintext output (no colors, tab-separated) |
| `--help` | `-h` | Help for any command |
| `--version` | `-v` | Show version |

## Issue Commands

### `issue list` (alias: `ls`)

List issues with optional filtering.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--assignee` | `-a` | | Email or `me` |
| `--state` | `-s` | | State name |
| `--team` | `-t` | | Team key |
| `--priority` | `-r` | -1 | 0=None, 1=Urgent, 2=High, 3=Normal, 4=Low |
| `--limit` | `-l` | 50 | Max results |
| `--sort` | `-o` | `linear` | `linear`, `created`, `updated` |
| `--newer-than` | `-n` | `6_months_ago` | Time filter |
| `--include-completed` | `-c` | false | Include done/canceled |
| `--view` | | | Execute custom view by ID |

### `issue search` (alias: `find`)

Full-text search across issues.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--assignee` | `-a` | | Email or `me` |
| `--state` | `-s` | | State name |
| `--team` | `-t` | | Team key |
| `--priority` | `-r` | -1 | Priority filter |
| `--limit` | `-l` | 50 | Max results |
| `--sort` | `-o` | `linear` | Sort order |
| `--newer-than` | `-n` | `6_months_ago` | Time filter |
| `--include-completed` | `-c` | false | Include done/canceled |
| `--include-archived` | | false | Include archived |

### `issue get` (alias: `show`)

Get detailed issue information. Shows description, state, assignee, team, priority, git branch, URL, attachments, relations, and recent comments.

```bash
linear-cli issue get ROB-27
linear-cli issue get ROB-27 --json    # Full JSON with all fields
```

### `issue create` (alias: `new`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--title` | | (required) | Issue title |
| `--team` | `-t` | (required) | Team key |
| `--description` | `-d` | | Description |
| `--priority` | | 3 | 0-4 |
| `--assign-me` | `-m` | false | Assign to self |
| `--label` | `-L` | | Label name (repeatable) |
| `--parent` | | | Parent issue ID |
| `--project` | | | Project ID |
| `--milestone` | | | Milestone ID or name (requires `--project`) |

### `issue update` (alias: `edit`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--title` | | | New title |
| `--description` | `-d` | | New description |
| `--assignee` | `-a` | | Email, name, `me`, or `unassigned` |
| `--state` | `-s` | | State name |
| `--priority` | | -1 | Priority 0-4 |
| `--due-date` | | | `YYYY-MM-DD` or empty to remove |
| `--milestone` | | | Milestone ID/name or `none` |
| `--parent` | | | Parent ID or `none` |

### `issue start`

Sets issue to "In Progress" and assigns to current user.

```bash
linear-cli issue start ROB-25
```

### `issue done`

Sets issue to the "Done" (completed) state.

```bash
linear-cli issue done ROB-25
```

### `issue assign`

Assigns issue to current user.

```bash
linear-cli issue assign ROB-25
```

### `issue archive` (aliases: `delete`, `rm`)

Soft delete — archived issues can be restored in Linear UI.

### `issue triage`

List issues in Triage or Backlog state for a team.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--limit` | `-l` | 50 | Max results |

```bash
linear-cli issue triage ROB
```

### `issue activity`

Show chronological activity timeline: state changes, assignee changes, priority changes, attachments, relations, comments.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--limit` | `-l` | 50 | History entries |

### `issue comment list` (alias: `ls`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--limit` | `-l` | 50 | Max comments |
| `--sort` | `-o` | `linear` | Sort order |

### `issue comment create` (aliases: `add`, `new`)

| Flag | Short | Description |
|------|-------|-------------|
| `--body` | `-b` | Comment body (required) |

### `issue comment update` / `issue comment delete`

Update or delete a comment by COMMENT-ID.

### `issue relation list`

List all relationships for an issue.

### `issue relation add`

| Flag | Description |
|------|-------------|
| `--type` | `blocks`, `blocked-by`, `related`, `duplicate`, `parent`, `sub-issue` |
| `--target` | Target issue ID |

### `issue relation remove`

Same flags as `add`.

### `issue relation update`

Update a relation's type by RELATION-ID.

### `issue attachment list`

List attachments for an issue.

### `issue attachment create`

| Flag | Description |
|------|-------------|
| `--url` | URL (required) |
| `--title` | Title |

### `issue attachment link`

Smart link — auto-detects GitHub PRs, etc.

| Flag | Description |
|------|-------------|
| `--url` | URL (required) |

### `issue attachment update` / `issue attachment delete`

By ATTACHMENT-ID.

## Project Commands

### `project list`

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--include-completed` | `-c` | false | Include completed |
| `--newer-than` | `-n` | `6_months_ago` | Time filter |

### `project get`

Get project details by PROJECT-ID.

### `project create` (alias: `new`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--name` | | (required) | Project name |
| `--description` | `-d` | | Description |
| `--team-ids` | | | Team UUIDs |
| `--state` | | `planned` | `planned`, `started`, `paused`, `completed`, `canceled` |
| `--start-date` | | | `YYYY-MM-DD` |
| `--target-date` | | | `YYYY-MM-DD` |

### `project update` (alias: `edit`)

Same flags as create (all optional), by PROJECT-ID.

### `project issues` (alias: `issue`)

List issues in a project.

| Flag | Short | Default |
|------|-------|---------|
| `--limit` | `-l` | 50 |

### `project archive` / `project delete`

Archive (soft) or delete (permanent!) a project.

### `project add-team` / `project remove-team`

Add/remove team associations by team KEY.

### `project milestone list` (alias: `ls`)

| Flag | Short | Default |
|------|-------|---------|
| `--limit` | `-l` | 50 |

### `project milestone create` / `update` / `delete` / `get`

CRUD for milestones. Create requires `--name`.

### `project status list` (alias: `ls`)

| Flag | Short | Default |
|------|-------|---------|
| `--limit` | `-l` | 20 |

### `project status create` (aliases: `new`, `add`)

| Flag | Short | Description |
|------|-------|-------------|
| `--body` | `-b` | Status text (required) |
| `--health` | | `onTrack`, `atRisk`, `offTrack` |

### `project status update` / `project status delete` / `project status get`

By UPDATE-ID.

## Cycle Commands

### `cycle list`

| Flag | Description |
|------|-------------|
| `--team` | Team key |
| `--active` | Show only active cycle |

### `cycle get`

By CYCLE-ID. Shows cycle details with issues.

### `cycle create`

| Flag | Description |
|------|-------------|
| `--team-id` | Team UUID (required) |
| `--starts` | Start date YYYY-MM-DD (required) |
| `--ends` | End date YYYY-MM-DD (required) |
| `--name` | Optional name |

### `cycle update`

By CYCLE-ID. Optional `--name`, `--starts`, `--ends`.

### `cycle archive`

By CYCLE-ID.

## Label Commands

### `label list`

| Flag | Description |
|------|-------------|
| `--team` | Team key |

### `label create`

| Flag | Description |
|------|-------------|
| `--name` | Label name (required) |
| `--color` | Hex color (e.g., `#e11d48`) |
| `--team` | Team ID |

### `label update`

By LABEL-ID. Optional `--name`, `--color`, `--description`.

### `label delete`

By LABEL-ID.

## Team Commands

### `team list`

Lists all teams with key, name, description, privacy, issue count.

### `team get`

By TEAM-KEY (e.g., `ROB`).

### `team members`

List members of a team by TEAM-KEY.

### `team states`

**Important**: Use this to discover valid `--state` values for issue filtering.

```bash
linear-cli team states ROB
# Output:
# NAME          TYPE        COLOR
# In Review     started     #0f783c
# Duplicate     canceled    #95a2b3
# Canceled      canceled    #95a2b3
# Todo          unstarted   #e2e2e2
# Backlog       backlog     #bec2c8
# Done          completed   #5e6ad2
# In Progress   started     #f2c94c
```

## User Commands

### `user list`

List all users. Shows name, email, role, status.

### `user get`

By email address.

### `user me`

Show current user (same as `whoami`).

## Document Commands

### `document list`

| Flag | Description |
|------|-------------|
| `--project` | Project ID |
| `--team` | Team key |

### `document get` / `document search` / `document create` / `document update` / `document delete`

Standard CRUD. Create requires `--title`, optional `--content` (markdown).

## Initiative Commands

### `initiative list`

| Flag | Description |
|------|-------------|
| `--status` | `Planned`, `Active`, `Completed` |
| `--include-completed` | Include completed |

### `initiative get` / `initiative create` / `initiative update` / `initiative delete`

Standard CRUD. Create requires `--name`.

### `initiative projects`

List projects under an initiative.

## View Commands

### `view list`

List all custom views.

### `view run` (alias: `exec`)

Execute a saved view filter and return matching issues/projects.

| Flag | Short | Default |
|------|-------|---------|
| `--limit` | `-l` | 50 |

### `view create` (alias: `new`)

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--name` | | (required) | View name |
| `--model` | `-m` | `issue` | `issue` or `project` |
| `--filter-json` | | | Raw JSON filter |
| `--description` | `-d` | | Description |
| `--shared` | | false | Make shared |
| `--team` | `-t` | | Team key |

### `view get` / `view update` / `view delete`

By VIEW-ID.

## GraphQL

```bash
linear-cli graphql QUERY [-v VARIABLES]    # Aliases: gql, gl
```

Runs arbitrary GraphQL against `https://api.linear.app/graphql`. Always returns JSON.

| Flag | Short | Description |
|------|-------|-------------|
| `--variables` | `-v` | JSON string of variables |

## Auth Commands

### `auth login`

Interactive — prompts for Personal API Key.

### `auth status`

Shows authenticated user and auth source (env var or config file).

### `auth rate-limit`

Shows current request and complexity limits.

### `auth logout`

Clears stored credentials.

## Utility Commands

### `whoami`

Shortcut for `user me`.

### `docs`

Outputs full embedded documentation in markdown.

### `completion`

Generate shell completion scripts: `bash`, `zsh`, `fish`, `powershell`.
