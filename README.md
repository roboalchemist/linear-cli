# linear-cli

A command-line interface for the [Linear API](https://developers.linear.app/), built for both humans and AI agents. Forked from [linctl](https://github.com/dorkitude/linear-cli) by [Kyle Wild](https://github.com/dorkitude).

## Features

- **12 entity types** with full CRUDL: issues, projects, cycles, labels, documents, initiatives, views, milestones, status updates, relations, attachments, comments
- **Quick actions**: `issue start` (In Progress + assign), `issue done`, `issue triage`, `issue archive`
- **Hierarchy navigation**: `initiative projects`, `project issues`
- **Raw GraphQL**: `graphql` command for arbitrary API queries
- **3 output formats**: colored table (default), plaintext (`-p`), JSON (`-j`)
- **Env var auth**: `LINEAR_API_KEY` for CI/CD pipelines
- **Time filtering**: `--newer-than 2_weeks_ago` on all list commands (default: 6 months)

## Installation

### Homebrew
```bash
brew tap roboalchemist/linear-cli
brew install linear-cli
```

### From Source
```bash
git clone https://github.com/roboalchemist/linear-cli.git
cd linear-cli
make build && make install
```

## Quick Start

```bash
# Authenticate
linear-cli auth

# Or use env var (CI/CD)
export LINEAR_API_KEY="lin_api_..."

# List your issues
linear-cli issue list --assignee me

# Start working on one
linear-cli issue start ENG-123

# Mark it done
linear-cli issue done ENG-123

# Run arbitrary GraphQL
linear-cli gql 'query { viewer { id name email } }'
```

## Commands

### Issue Management
```bash
linear-cli issue list [flags]              # List issues (aliases: ls)
linear-cli issue search "query" [flags]    # Full-text search
linear-cli issue get ISSUE-ID              # Get details (aliases: show)
linear-cli issue create [flags]            # Create issue (aliases: new)
linear-cli issue update ISSUE-ID [flags]   # Update issue (aliases: edit)
linear-cli issue assign ISSUE-ID           # Assign to yourself
linear-cli issue start ISSUE-ID            # Set In Progress + assign to me
linear-cli issue done ISSUE-ID             # Mark as Done
linear-cli issue archive ISSUE-ID          # Archive (soft delete)
linear-cli issue triage TEAM-KEY           # List untriaged/backlog issues
linear-cli issue activity ISSUE-ID         # Show activity timeline

# Issue list flags
  -a, --assignee string     Filter by assignee (email or 'me')
  -s, --state string        Filter by state name
  -t, --team string         Filter by team key
  -r, --priority int        Filter by priority (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)
  -L, --label strings       Filter/set labels by name (repeatable)
  -l, --limit int           Max results (default 50)
  -o, --sort string         Sort: linear (default), created, updated
  -n, --newer-than string   Time filter (default: 6_months_ago, use 'all_time' for all)
  -c, --include-completed   Include completed/canceled issues
      --view string         Execute a custom view by ID (overrides other filters)

# Issue create flags
      --title string        Issue title (required)
  -d, --description string  Description
  -t, --team string         Team key (required)
      --priority int        Priority 0-4 (default 3)
  -m, --assign-me           Assign to yourself
      --parent string       Parent issue identifier
      --project string      Project ID
      --milestone string    Milestone ID or name (requires --project)
  -L, --label strings       Label names (repeatable, case-insensitive)

# Issue update flags
      --title string        New title
  -d, --description string  New description
  -a, --assignee string     Assignee (email, name, 'me', or 'unassigned')
  -s, --state string        State name (e.g., 'Todo', 'In Progress', 'Done')
      --priority int        Priority (0-4)
      --due-date string     Due date (YYYY-MM-DD, or empty to remove)
      --milestone string    Milestone ID or name (or 'none' to unset)
      --parent string       Parent issue (or 'none' to unset)
```

### Comments (under issue)
```bash
linear-cli issue comment list ISSUE-ID     # List comments (aliases: ls)
linear-cli issue comment create ISSUE-ID   # Add comment (aliases: add, new)
linear-cli issue comment update COMMENT-ID # Edit comment (aliases: edit)
linear-cli issue comment delete COMMENT-ID # Delete comment (aliases: rm)

# Create/update flags
  -b, --body string         Comment body (required)
```

### Issue Relations
```bash
linear-cli issue relation list ISSUE-ID                          # List relations
linear-cli issue relation add ISSUE-ID --type blocks --target ID # Add relation
linear-cli issue relation remove ISSUE-ID --type blocks --target ID
linear-cli issue relation update RELATION-ID --type related

# Types: blocks, blocked-by, related, duplicate, parent, sub-issue
```

### Issue Attachments
```bash
linear-cli issue attachment list ISSUE-ID
linear-cli issue attachment create ISSUE-ID --url URL --title TITLE
linear-cli issue attachment link ISSUE-ID --url URL    # Smart link (auto-detects GitHub PRs, etc.)
linear-cli issue attachment update ATTACHMENT-ID --title TITLE
linear-cli issue attachment delete ATTACHMENT-ID
```

### Projects
```bash
linear-cli project list [flags]            # List projects
linear-cli project get PROJECT-ID          # Get details
linear-cli project create [flags]          # Create project
linear-cli project update PROJECT-ID       # Update project
linear-cli project archive PROJECT-ID      # Archive project
linear-cli project delete PROJECT-ID       # Permanently delete
linear-cli project issues PROJECT-ID       # List issues in project
linear-cli project add-team PROJECT-ID KEY # Add team(s)
linear-cli project remove-team PROJECT-ID KEY

# Create flags
      --name string         Project name (required)
  -d, --description string  Description
      --team-ids strings    Team UUIDs to associate
      --state string        State: planned, started, paused, completed, canceled
      --start-date string   Start date (YYYY-MM-DD)
      --target-date string  Target date (YYYY-MM-DD)
```

### Milestones (under project)
```bash
linear-cli project milestone list PROJECT-ID
linear-cli project milestone get MILESTONE-ID
linear-cli project milestone create PROJECT-ID --name NAME
linear-cli project milestone update MILESTONE-ID --name NAME
linear-cli project milestone delete MILESTONE-ID
```

### Project Status Updates
```bash
linear-cli project status list PROJECT-ID
linear-cli project status get UPDATE-ID
linear-cli project status create PROJECT-ID --body TEXT [--health onTrack|atRisk|offTrack]
linear-cli project status update UPDATE-ID --body TEXT
linear-cli project status delete UPDATE-ID
```

### Cycles (Sprints)
```bash
linear-cli cycle list [--team KEY] [--active]
linear-cli cycle get CYCLE-ID
linear-cli cycle create --team-id UUID --starts YYYY-MM-DD --ends YYYY-MM-DD [--name NAME]
linear-cli cycle update CYCLE-ID [--name NAME] [--starts DATE] [--ends DATE]
linear-cli cycle archive CYCLE-ID
```

### Labels
```bash
linear-cli label list [--team KEY]
linear-cli label create --name NAME [--color HEX]
linear-cli label update LABEL-ID [--name NAME] [--color HEX] [--description TEXT]
linear-cli label delete LABEL-ID
```

### Teams
```bash
linear-cli team list
linear-cli team get TEAM-KEY
linear-cli team members TEAM-KEY
linear-cli team states TEAM-KEY            # Show workflow states (helps discover --state values)
```

### Initiatives
```bash
linear-cli initiative list [--status Active] [--include-completed]
linear-cli initiative get INITIATIVE-ID
linear-cli initiative create --name NAME [--status Planned|Active|Completed]
linear-cli initiative update INITIATIVE-ID [--name NAME] [--status STATUS]
linear-cli initiative delete INITIATIVE-ID
linear-cli initiative projects INITIATIVE-ID   # List projects under initiative
```

### Documents
```bash
linear-cli document list [--project ID] [--team KEY]
linear-cli document get DOC-ID
linear-cli document search "query"
linear-cli document create --title TITLE [--content MD] [--project ID]
linear-cli document update DOC-ID [--title TITLE] [--content MD]
linear-cli document delete DOC-ID
```

### Custom Views
```bash
linear-cli view list
linear-cli view get VIEW-ID
linear-cli view run VIEW-ID                # Execute saved filters
linear-cli view create --name NAME [--model issue|project] [--filter-json JSON]
linear-cli view update VIEW-ID [--name NAME]
linear-cli view delete VIEW-ID

# Also available as: linear-cli issue list --view VIEW-ID
```

### Users
```bash
linear-cli user list [--active]
linear-cli user get EMAIL
linear-cli user me
linear-cli whoami                          # Shortcut for user me
```

### Raw GraphQL
```bash
linear-cli graphql 'query { viewer { id name } }'     # aliases: gql, gl
linear-cli gql 'query($id: String!) { issue(id: $id) { title } }' \
  -v '{"id": "UUID"}'
```

### Authentication
```bash
linear-cli auth login                      # Interactive login
linear-cli auth status                     # Check auth (shows source: env var or config)
linear-cli auth logout                     # Clear stored credentials

# Environment variable override (useful for CI/CD)
export LINEAR_API_KEY="lin_api_..."         # Primary
export LINCTL_API_KEY="lin_api_..."         # Legacy alias
# Precedence: LINEAR_API_KEY > LINCTL_API_KEY > config file
```

## Global Flags

```
-p, --plaintext   Plaintext output (tab-separated, no colors)
-j, --json        JSON output (for scripting/agents)
-h, --help        Help for any command
-v, --version     Show version
    --config      Config file (default: ~/.linear-cli.yaml)
```

## Default Filters

List commands default to showing items from the **last 6 months** and **exclude completed/canceled** items. Override with:
- `--newer-than all_time` to see everything
- `--include-completed` to include done/canceled items
- `--include-archived` (on issue search) to include archived items

### Time expressions
```
1_day_ago, 2_weeks_ago, 3_months_ago, 1_year_ago, all_time, 2025-07-01
```

## Authentication

Credentials are stored in `~/.linear-cli-auth.json` (0600 permissions).

1. Get a Personal API Key from [Linear Settings > API](https://linear.app/settings/api)
2. Run `linear-cli auth` and paste your key

For CI/CD, set `LINEAR_API_KEY` environment variable instead.

## Testing

```bash
make test           # Run smoke tests (requires valid auth)
make test-verbose   # With bash tracing
```

Smoke tests exercise all read-only commands across all 3 output formats. Write operations are tested manually against a test workspace.

## Development

```bash
go run main.go <command> [flags]   # Run without building
make build                         # Build binary
make fmt                           # Format code
make lint                          # Lint (requires golangci-lint)
make everything                    # Build + fmt + lint + test + install
```

## License

MIT - see [LICENSE](LICENSE)

## Links

- [Linear API Documentation](https://developers.linear.app/)
- [GraphQL Schema Explorer](https://studio.apollographql.com/public/Linear-API/variant/current/schema/reference)
- [GitHub Repository](https://github.com/roboalchemist/linear-cli)
