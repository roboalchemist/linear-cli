# linear-cli

CLI for the [Linear API](https://developers.linear.app/) — issues, projects, cycles, labels, documents, initiatives, views, teams, and raw GraphQL. Built for both humans and AI agents.

**Version**: 0.3.0 | **Binary**: `linear-cli` | **Config**: `~/.linear-cli.yaml` | **Auth**: `~/.linear-cli-auth.json`

## Setup

```bash
brew tap roboalchemist/tap && brew install linear-cli
linear-cli auth                    # Interactive (paste Personal API Key)
export LINEAR_API_KEY="lin_api_…"  # Or env var for CI/agents
linear-cli auth status             # Verify
```

## Output Formats

| Flag | Format | Use case |
|------|--------|----------|
| (none) | Colored table | Human terminal |
| `-p` / `--plaintext` | Markdown/plain | Piping, agents |
| `-j` / `--json` | JSON | Programmatic parsing |

**Always use `--json` for programmatic parsing.** Use `--plaintext` for readable agent output.

## Default Filters

List commands default to **last 6 months** and **exclude completed/canceled**. Override:
- `--newer-than all_time` — show everything
- `--include-completed` — include done/canceled
- `--include-archived` — include archived (search only)

Time expressions: `1_day_ago`, `2_weeks_ago`, `3_months_ago`, `1_year_ago`, `all_time`, `2025-07-01`

<examples>
<example>
Task: List my in-progress issues

```bash
linear-cli issue list --assignee me --state "In Progress"
```

Output:
```
TITLE                       STATE         ASSIGNEE             TEAM   CREATED      URL
CRUD test (updated)         In Progress   joe@schlesinger.io   ROB    2026-02-09   https://linear.app/roboalchemist/issue/ROB-27/crud-test-updated

✓ 1 issues
```
</example>

<example>
Task: Get full details of an issue in JSON for parsing

```bash
linear-cli issue get ROB-27 --json
```

Output (truncated):
```json
{
  "id": "7582dfde-ca3c-48bb-98c0-96e786dee3f4",
  "identifier": "ROB-27",
  "title": "CRUD test (updated)",
  "state": { "name": "In Progress", "type": "started" },
  "assignee": { "email": "joe@schlesinger.io" },
  "team": { "key": "ROB" },
  "url": "https://linear.app/roboalchemist/issue/ROB-27/crud-test-updated"
}
```
</example>

<example>
Task: Create a bug issue, assign to me, with labels

```bash
linear-cli issue create --title "Login broken on Safari" --team ROB --priority 2 --assign-me --label Bug --description "Users can't log in on Safari 17"
```

Output:
```
✅ Created issue ROB-34: Login broken on Safari
https://linear.app/roboalchemist/issue/ROB-34/login-broken-on-safari
```
</example>

<example>
Task: Start working on an issue (sets In Progress + assigns to me)

```bash
linear-cli issue start ROB-34
```

Output:
```
✅ Started ROB-34: Login broken on Safari
   State: In Progress | Assigned: joe@schlesinger.io
```
</example>

<example>
Task: Run a raw GraphQL query to get viewer info

```bash
linear-cli gql 'query { viewer { id name email } }'
```

Output:
```json
{
  "viewer": {
    "id": "71732280-872f-4c10-8973-41ebdff055c2",
    "name": "joe@schlesinger.io",
    "email": "joe@schlesinger.io"
  }
}
```
</example>
</examples>

## Quick Reference

### Issues (most used)

```bash
linear-cli issue list [flags]              # List issues (alias: ls)
linear-cli issue search "query" [flags]    # Full-text search (alias: find)
linear-cli issue get ISSUE-ID              # Get details (alias: show)
linear-cli issue create [flags]            # Create (alias: new)
linear-cli issue update ISSUE-ID [flags]   # Update (alias: edit)
linear-cli issue start ISSUE-ID            # Set In Progress + assign to me
linear-cli issue done ISSUE-ID             # Mark as Done
linear-cli issue assign ISSUE-ID           # Assign to yourself
linear-cli issue archive ISSUE-ID          # Soft delete (alias: delete, rm)
linear-cli issue triage TEAM-KEY           # List untriaged/backlog issues
linear-cli issue activity ISSUE-ID         # Activity timeline
```

### Issue List Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--assignee` | `-a` | Filter by email or `me` |
| `--state` | `-s` | Filter by state name |
| `--team` | `-t` | Filter by team key |
| `--priority` | `-r` | 0=None, 1=Urgent, 2=High, 3=Normal, 4=Low |
| `--limit` | `-l` | Max results (default 50) |
| `--sort` | `-o` | `linear` (default), `created`, `updated` |
| `--newer-than` | `-n` | Time filter (default: `6_months_ago`) |
| `--include-completed` | `-c` | Include done/canceled |
| `--view` | | Execute a custom view by ID |

### Issue Create Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--title` | | Title (required) |
| `--team` | `-t` | Team key (required) |
| `--description` | `-d` | Description |
| `--priority` | | 0-4 (default 3=Normal) |
| `--assign-me` | `-m` | Assign to yourself |
| `--label` | `-L` | Label name (repeatable) |
| `--parent` | | Parent issue ID |
| `--project` | | Project ID |
| `--milestone` | | Milestone (requires `--project`) |

### Issue Update Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--title` | | New title |
| `--description` | `-d` | New description |
| `--assignee` | `-a` | Email, name, `me`, or `unassigned` |
| `--state` | `-s` | State name (e.g., `Todo`, `In Progress`, `Done`) |
| `--priority` | | 0-4 |
| `--due-date` | | `YYYY-MM-DD` or empty to remove |
| `--milestone` | | Milestone or `none` to unset |
| `--parent` | | Parent issue or `none` to unset |

### Comments, Relations, Attachments

```bash
# Comments
linear-cli issue comment list ISSUE-ID
linear-cli issue comment create ISSUE-ID --body "text"
linear-cli issue comment update COMMENT-ID --body "new text"
linear-cli issue comment delete COMMENT-ID

# Relations (types: blocks, blocked-by, related, duplicate, parent, sub-issue)
linear-cli issue relation list ISSUE-ID
linear-cli issue relation add ISSUE-ID --type blocks --target OTHER-ID
linear-cli issue relation remove ISSUE-ID --type blocks --target OTHER-ID

# Attachments
linear-cli issue attachment list ISSUE-ID
linear-cli issue attachment create ISSUE-ID --url URL --title TITLE
linear-cli issue attachment link ISSUE-ID --url URL    # Smart link (auto-detects GitHub PRs)
```

### Projects

```bash
linear-cli project list [flags]
linear-cli project get PROJECT-ID
linear-cli project create --name NAME [--team-ids UUID] [--state planned|started|paused|completed|canceled]
linear-cli project update PROJECT-ID [flags]
linear-cli project issues PROJECT-ID           # List issues in project
linear-cli project archive PROJECT-ID
linear-cli project delete PROJECT-ID           # Permanent delete
linear-cli project add-team PROJECT-ID KEY
linear-cli project remove-team PROJECT-ID KEY
```

### Milestones & Status Updates

```bash
# Milestones (under project)
linear-cli project milestone list PROJECT-ID
linear-cli project milestone create PROJECT-ID --name NAME
linear-cli project milestone update MILESTONE-ID --name NAME
linear-cli project milestone delete MILESTONE-ID

# Status updates (health: onTrack, atRisk, offTrack)
linear-cli project status list PROJECT-ID
linear-cli project status create PROJECT-ID --body "text" --health onTrack
linear-cli project status update UPDATE-ID --body "text"
```

### Teams & Users

```bash
linear-cli team list
linear-cli team get TEAM-KEY
linear-cli team members TEAM-KEY
linear-cli team states TEAM-KEY      # Discover valid --state values!

linear-cli user list
linear-cli user me
linear-cli whoami                    # Shortcut for user me
```

### Other Entities

```bash
# Labels
linear-cli label list [--team KEY]
linear-cli label create --name NAME [--color HEX] [--team TEAM-ID]

# Cycles
linear-cli cycle list [--team KEY] [--active]
linear-cli cycle get CYCLE-ID

# Documents
linear-cli document list [--project ID] [--team KEY]
linear-cli document search "query"
linear-cli document create --title TITLE [--content MD]

# Initiatives
linear-cli initiative list [--status Active]
linear-cli initiative get INITIATIVE-ID
linear-cli initiative projects INITIATIVE-ID

# Views
linear-cli view list
linear-cli view run VIEW-ID          # Execute saved filter, returns matching issues
linear-cli view create --name NAME [--model issue|project] [--filter-json JSON]
```

### GraphQL

```bash
linear-cli gql 'query { viewer { id name } }'                    # Aliases: graphql, gql, gl
linear-cli gql 'query($id: String!) { issue(id: $id) { title } }' -v '{"id":"UUID"}'
```

### Auth & Utility

```bash
linear-cli auth status               # Check auth (shows source: env var or config)
linear-cli auth rate-limit            # Show API rate limit status
linear-cli auth login                 # Interactive login
linear-cli auth logout                # Clear credentials
linear-cli docs                       # Show full embedded documentation
linear-cli skill print                # Print embedded Claude Code skill
linear-cli skill add                  # Install skill to ~/.claude/skills/
linear-cli completion zsh             # Generate shell completions
```

## Agent Best Practices

1. **Always use `--json`** for programmatic parsing, pipe to `jq`
2. **Use `--plaintext`** when you need readable structured output without ANSI colors
3. **Prefix comments with "🤖 says:"** so users know they're automated
4. **Use `team states TEAM-KEY`** to discover valid `--state` values before filtering
5. **Auth via env var** `LINEAR_API_KEY` — no interactive login needed
6. **Rate limits**: 5000 requests/hour, 3M complexity points/hour — check with `auth rate-limit`
7. **IDs**: Issue identifiers like `ROB-27` work everywhere. UUIDs needed for projects, milestones, etc.
8. **Env vars**: `LINEAR_API_KEY` (primary), `LINCTL_API_KEY` (legacy). Precedence: env > config file

## Gotchas

- **Priority numbering is inverted**: 1=Urgent (highest), 4=Low (lowest), 0=None
- **Default 6-month filter**: `issue list` hides older issues unless `--newer-than all_time`
- **State names are case-sensitive**: Use `"In Progress"` not `"in progress"` — discover with `team states`
- **`issue archive` is aliased to `issue delete`/`issue rm`** — it's a soft delete (restorable in UI)
- **`project delete` is permanent** — unlike issue archive
- **Milestone requires `--project`** on issue create
- **Label flag is `-L` (uppercase)** not `-l` (which is `--limit`)

See [reference/commands.md](reference/commands.md) for complete command details.
