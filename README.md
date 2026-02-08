# 🚀 linear-cli - Linear CLI Tool

> **Fork notice:** `linear-cli` is a fork of [`linctl`](https://github.com/dorkitude/linctl), the original Linear CLI tool. This fork adds initiative management, attachment CRUDL, issue activity timelines, and renames the binary for clarity.

A comprehensive command-line interface for Linear's API, built with agents in mind (but nice for humans too).

## ✨ Features

- 🔐 **Authentication**: Personal API Key support
- 📋 **Issue Management**: Create, list, view, update, assign, and manage issues with full details
  - Sub-issue hierarchy with parent/child relationships
  - Git branch integration showing linked branches
  - Cycle (sprint) and project associations
  - Attachments and recent comments preview
  - Due dates, snoozed status, and completion tracking
  - Full-text search via `linear-cli issue search`
- 👥 **Team Management**: View teams, get team details, and list team members
- 🚀 **Project Tracking**: Comprehensive project information
  - Progress visualization with issue statistics
  - Team and member associations
  - Milestone management (create, list, update, delete)
  - Recent issues preview
  - Timeline tracking (created, updated, completed dates)
- 👤 **User Management**: List all users, view user details, and current user info
- 📄 **Document Management**: List, view, search, create, update, and delete documents
  - Full-text search across all documents
  - Project and team associations
  - Full markdown content display
- 🔍 **Custom Views**: Run saved filters to query matching issues or projects
  - List, create, update, and delete custom views
  - Execute views to see matching issues or projects
  - Support for shared and team-scoped views
- 🎯 **Initiative Management**: CRUDL for high-level strategic objectives
  - List, create, update, delete initiatives
  - View linked projects and sub-initiatives
  - Filter by status (Planned, Active, Completed)
- 📎 **Attachment Management**: Full CRUDL for issue attachments
  - Create, list, update, delete attachments
  - Smart URL linking (auto-detects GitHub PRs, Slack threads, etc.)
- 📊 **Issue Activity Timeline**: Chronological view of all issue changes
  - State, assignee, priority, project, cycle changes
  - Label additions/removals, linked attachments, comments
- 💬 **Comments**: List and create comments on issues with time-aware formatting
- 🔗 **Webhooks**: Configure and manage webhooks
- 🎨 **Multiple Output Formats**: Table, plaintext, and JSON output
- ⚡ **Performance**: Fast and lightweight CLI tool
- 🔄 **Flexible Sorting**: Sort lists by Linear's default order, creation date, or update date
- 📅 **Time-based Filtering**: Filter lists by creation date with intuitive time expressions
- 📚 **Built-in Documentation**: Access full documentation with `linear-cli docs`
- 🧪 **Smoke Testing**: Automated smoke tests for all read-only commands

## 🛠️ Installation

### Homebrew (macOS/Linux)
```bash
brew tap dorkitude/linear-cli
brew install linear-cli
linear-cli docs      # Render the README.md
```

### From Source
```bash
git clone https://github.com/dorkitude/linear-cli.git
cd linear-cli
make deps        # Install dependencies
make build       # Build the binary
make install     # Install to /usr/local/bin (requires sudo)
linear-cli docs      # Render the README.md
```

### For Development
```bash
git clone https://github.com/dorkitude/linear-cli.git
cd linear-cli
make deps        # Install dependencies
go run main.go   # Run directly without building
make dev         # Or build and run in development mode
make test        # Run all tests
make lint        # Run linter
make fmt         # Format code
linear-cli docs      # Render the README.md
```

## Important: Default Filters

**By default, `issue list`, `issue search`, and `project list` commands only show items created in the last 6 months!**
 
This improves performance and prevents overwhelming data loads. To see older items:
 - Use `--newer-than 1_year_ago` for items from the last year
 - Use `--newer-than all_time` to see ALL items ever created
 - See the [Time-based Filtering](#-time-based-filtering) section for details

**By default, `issue list` and `issue search` also filter out canceled and completed items. To see all items, use the `--include-completed` flag.**
- Need archived matches? Add `--include-archived` when using `issue search`.


## 🚀 Quick Start

> **IMPORTANT**  Agents like Claude Code, Cursor, and Gemini should use the `--json` flag on all read operations.

### 1. Authentication
```bash
# Interactive authentication
linear-cli auth

# Check authentication status
linear-cli auth status

# Show current user
linear-cli whoami

# View full documentation
linear-cli docs | less
```

### 2. Issue Management
```bash
# List all issues
linear-cli issue list

# List issues assigned to you
linear-cli issue list --assignee me

# List issues in a specific state
linear-cli issue list --state "In Progress"

# List issues sorted by update date
linear-cli issue list --sort updated

# Search issues using Linear's full-text index (shares the same filters as list)
linear-cli issue search "login bug" --team ENG
linear-cli issue search "customer:" --include-completed --include-archived

# List recent issues (last 2 weeks instead of default 6 months)
linear-cli issue list --newer-than 2_weeks_ago

# List ALL issues ever created (override 6-month default)
linear-cli issue list --newer-than all_time

# List today's issues
linear-cli issue list --newer-than 1_day_ago

# Get issue details (now includes git branch, cycle, project, attachments, and comments)
linear-cli issue get LIN-123

# Create a new issue
linear-cli issue create --title "Bug fix" --team ENG

# Assign issue to yourself
linear-cli issue assign LIN-123

# Update issue fields
linear-cli issue update LIN-123 --title "New title"
linear-cli issue update LIN-123 --description "Updated description"
linear-cli issue update LIN-123 --assignee john.doe@company.com
linear-cli issue update LIN-123 --assignee me  # Assign to yourself
linear-cli issue update LIN-123 --assignee unassigned  # Remove assignee
linear-cli issue update LIN-123 --state "In Progress"
linear-cli issue update LIN-123 --priority 1  # 0=None, 1=Urgent, 2=High, 3=Normal, 4=Low
linear-cli issue update LIN-123 --due-date "2024-12-31"
linear-cli issue update LIN-123 --due-date ""  # Remove due date

# Update multiple fields at once
linear-cli issue update LIN-123 --title "Critical Bug" --assignee me --priority 1

# Set or remove parent issue
linear-cli issue update LIN-123 --parent LIN-100
linear-cli issue update LIN-123 --parent none
linear-cli issue create --title "Sub-task" --team ENG --parent LIN-100
```

### Issue Relations
```bash
# List all relationships for an issue
linear-cli issue relation list LIN-123

# Add relations
linear-cli issue relation add LIN-123 --type blocks --target LIN-456
linear-cli issue relation add LIN-123 --type blocked-by --target LIN-456
linear-cli issue relation add LIN-123 --type related --target LIN-456
linear-cli issue relation add LIN-123 --type duplicate --target LIN-456

# Set parent/sub-issue
linear-cli issue relation add LIN-123 --type parent --target LIN-100
linear-cli issue relation add LIN-123 --type sub-issue --target LIN-789

# Remove relations
linear-cli issue relation remove LIN-123 --type blocks --target LIN-456
linear-cli issue relation remove LIN-123 --type parent --target LIN-100

# Update a relation's type (use relation ID from --json output)
linear-cli issue relation update RELATION-UUID --type related
```

### Issue Attachments
```bash
# List attachments for an issue
linear-cli issue attachment list LIN-123

# Create a manual attachment
linear-cli issue attachment create LIN-123 --url "https://example.com/spec" --title "Spec Doc"

# Smart link (auto-detects GitHub PR, Slack, Notion, etc.)
linear-cli issue attachment link LIN-123 --url "https://github.com/org/repo/pull/42"

# Update an attachment
linear-cli issue attachment update ATTACHMENT-ID --title "New Title"

# Delete an attachment
linear-cli issue attachment delete ATTACHMENT-ID
```

### Issue Activity Timeline
```bash
# Show full activity timeline for an issue
linear-cli issue activity LIN-123

# Show more history entries
linear-cli issue activity LIN-123 --limit 100

# JSON output for scripting
linear-cli issue activity LIN-123 --json
```

### 3. Project Management
```bash
# List all projects (shows IDs)
linear-cli project list

# Filter projects by team
linear-cli project list --team ENG

# List projects created in the last month (instead of default 6 months)
linear-cli project list --newer-than 1_month_ago

# List ALL projects regardless of age
linear-cli project list --newer-than all_time

# Get project details (use ID from list command)
linear-cli project get 65a77a62-ec5e-491e-b1d9-84aebee01b33

# Add teams to a project (required for cross-team issue assignment)
linear-cli project add-team PROJECT-ID ENG
linear-cli project add-team PROJECT-ID ENG DESIGN OPS

# Remove teams from a project
linear-cli project remove-team PROJECT-ID ENG
```

### Milestone Management (within Projects)
```bash
# List milestones for a project
linear-cli project milestone list PROJECT-ID

# Get milestone details (including issues)
linear-cli project milestone get MILESTONE-ID

# Create a milestone
linear-cli project milestone create PROJECT-ID --name "Beta Release"
linear-cli project milestone create PROJECT-ID --name "GA" --target-date "2025-06-01" --description "General availability"

# Update a milestone
linear-cli project milestone update MILESTONE-ID --name "New Name"
linear-cli project milestone update MILESTONE-ID --target-date "2025-07-01"

# Delete a milestone
linear-cli project milestone delete MILESTONE-ID

# Assign a milestone to an issue (by name or ID)
linear-cli issue update LIN-123 --milestone "Beta Release"
linear-cli issue update LIN-123 --milestone none  # Remove milestone

# Create an issue with a milestone
linear-cli issue create --title "Fix bug" --team ENG --project PROJECT-ID --milestone "Beta Release"
```

### Project Status Updates
```bash
# List status updates for a project
linear-cli project status list PROJECT-ID

# Get a specific status update
linear-cli project status get UPDATE-ID

# Create a status update (health: onTrack, atRisk, offTrack)
linear-cli project status create PROJECT-ID --body "Sprint on track" --health onTrack
linear-cli project status create PROJECT-ID --body "Blocked on API" --health atRisk

# Update a status update
linear-cli project status update UPDATE-ID --body "Updated text"
linear-cli project status update UPDATE-ID --health offTrack

# Archive (soft-delete) a status update
linear-cli project status delete UPDATE-ID
```

### 4. Team Management
```bash
# List all teams
linear-cli team list

# Get team details
linear-cli team get ENG

# List team members
linear-cli team members ENG
```

### 5. User Management
```bash
# List all users
linear-cli user list

# Show only active users
linear-cli user list --active

# Get user details by email
linear-cli user get john@example.com

# Show your own profile
linear-cli user me
```

### 6. Documents
```bash
# List all documents
linear-cli document list

# List documents for a project
linear-cli document list --project PROJECT-ID

# Search documents
linear-cli document search "onboarding spec"

# View a document with full content
linear-cli document get DOC-ID

# Create a document
linear-cli document create --title "API Spec" --content "# Overview" --project PROJECT-ID

# Update a document
linear-cli document update DOC-ID --title "Updated Title"

# Delete a document
linear-cli document delete DOC-ID
```

### 7. Comments
```bash
# List comments on an issue
linear-cli comment list LIN-123

# Add a comment to an issue
linear-cli comment create LIN-123 --body "Fixed the authentication bug"
```

## 📖 Command Reference

### Global Flags
- `--plaintext, -p`: Plain text output (non-interactive)
- `--json, -j`: JSON output for scripting
- `--help, -h`: Show help
- `--version, -v`: Show version

### Authentication Commands
```bash
linear-cli auth               # Interactive authentication
linear-cli auth login         # Same as above
linear-cli auth status        # Check authentication status
linear-cli auth logout        # Clear stored credentials
linear-cli whoami            # Show current user
```

### Issue Commands
```bash
# List issues with filters
linear-cli issue list [flags]
linear-cli issue ls [flags]     # Short alias

# Flags:
  -a, --assignee string     Filter by assignee (email or 'me')
  -c, --include-completed   Include completed and canceled issues
  -s, --state string       Filter by state name
  -t, --team string        Filter by team key
  -r, --priority int       Filter by priority (0-4, default: -1)
  -l, --limit int          Maximum results (default 50)
  -o, --sort string        Sort order: linear (default), created, updated
  -n, --newer-than string  Show items created after this time (default: 6_months_ago, use 'all_time' for no filter)

# Get issue details (shows parent and sub-issues)
linear-cli issue get <issue-id>
linear-cli issue show <issue-id>  # Alias

# Create issue
linear-cli issue create [flags]
linear-cli issue new [flags]      # Alias
# Flags:
  --title string           Issue title (required)
  -d, --description string Issue description
  -t, --team string        Team key (required)
  --priority int       Priority 0-4 (default 3)
  -m, --assign-me          Assign to yourself
  --parent string          Parent issue identifier

# Assign issue to yourself
linear-cli issue assign <issue-id>

# Update issue
linear-cli issue update <issue-id> [flags]
linear-cli issue edit <issue-id> [flags]    # Alias
# Flags:
  --title string           New title
  -d, --description string New description
  -a, --assignee string    Assignee (email, name, 'me', or 'unassigned')
  -s, --state string       State name (e.g., 'Todo', 'In Progress', 'Done')
  --priority int           Priority (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)
  --due-date string        Due date (YYYY-MM-DD format, or empty to remove)
  --milestone string       Milestone ID or name (or 'none' to unset)
  --parent string          Parent issue identifier (or 'none' to unset)

# Issue Relations
linear-cli issue relation list <issue-id>
linear-cli issue relation ls <issue-id>       # Alias
linear-cli issue relation add <issue-id> [flags]
linear-cli issue relation create <issue-id> [flags]   # Alias
linear-cli issue relation remove <issue-id> [flags]
linear-cli issue relation rm <issue-id> [flags]       # Alias
linear-cli issue relation update <relation-id> [flags]
linear-cli issue relation edit <relation-id> [flags]  # Alias
# Flags (add/remove):
  --type string            Relation type: blocks, blocked-by, related, duplicate, parent, sub-issue
  --target string          Target issue identifier
# Flags (update):
  --type string            New type: blocks, related, duplicate

# Archive issue (coming soon)
linear-cli issue archive <issue-id>
```

### Team Commands
```bash
# List all teams with issue counts
linear-cli team list
linear-cli team ls              # Alias
# Flags:
  -l, --limit int          Maximum results (default 50)
  -o, --sort string        Sort order: linear (default), created, updated

# Get team details
linear-cli team get <team-key>
linear-cli team show <team-key> # Alias

# Examples:
linear-cli team get ENG         # Shows Engineering team details
linear-cli team get DESIGN      # Shows Design team details

# List team members with roles and status
linear-cli team members <team-key>

# Examples:
linear-cli team members ENG     # Lists all Engineering team members
```

### Project Commands
```bash
# List projects
linear-cli project list [flags]
linear-cli project ls [flags]     # Alias
# Flags:
  -t, --team string        Filter by team key
  -s, --state string       Filter by state (planned, started, paused, completed, canceled)
  -l, --limit int          Maximum results (default 50)
  -o, --sort string        Sort order: linear (default), created, updated
  -n, --newer-than string  Show items created after this time (default: 6_months_ago)
  -c, --include-completed  Include completed and canceled projects

# Get project details
linear-cli project get <project-id>
linear-cli project show <project-id>  # Alias

# Add teams to a project (required for cross-team issue assignment)
linear-cli project add-team <project-id> <team-key> [team-key...]

# Remove teams from a project
linear-cli project remove-team <project-id> <team-key> [team-key...]

# Examples:
linear-cli project add-team PROJECT-ID ENG DESIGN   # Add ENG and DESIGN teams
linear-cli project remove-team PROJECT-ID OPS        # Remove OPS team

# Create project (coming soon)
linear-cli project create [flags]
```

### Milestone Commands (under project)
```bash
# List milestones for a project
linear-cli project milestone list <project-id> [flags]
linear-cli project milestone ls <project-id> [flags]  # Alias
# Flags:
  -l, --limit int          Maximum results (default 50)

# Get milestone details
linear-cli project milestone get <milestone-id>
linear-cli project milestone show <milestone-id>  # Alias

# Create milestone
linear-cli project milestone create <project-id> [flags]
linear-cli project milestone new <project-id> [flags]    # Alias
# Flags:
  --name string            Milestone name (required)
  -d, --description string Milestone description
  --target-date string     Target date (YYYY-MM-DD)

# Update milestone
linear-cli project milestone update <milestone-id> [flags]
linear-cli project milestone edit <milestone-id> [flags]  # Alias
# Flags:
  --name string            New name
  -d, --description string New description
  --target-date string     New target date (YYYY-MM-DD, or empty to remove)

# Delete milestone
linear-cli project milestone delete <milestone-id>
```

### Project Status Update Commands
```bash
# List status updates for a project
linear-cli project status list <project-id> [flags]
linear-cli project status ls <project-id> [flags]  # Alias
# Flags:
  -l, --limit int          Maximum results (default 20)

# Get status update details
linear-cli project status get <update-id>
linear-cli project status show <update-id>  # Alias

# Create status update
linear-cli project status create <project-id> [flags]
linear-cli project status new <project-id> [flags]    # Alias
# Flags:
  -b, --body string        Status update body text (required)
  --health string           Project health: onTrack, atRisk, offTrack

# Update status update
linear-cli project status update <update-id> [flags]
linear-cli project status edit <update-id> [flags]  # Alias
# Flags:
  -b, --body string        New body text
  --health string           New health: onTrack, atRisk, offTrack

# Archive (soft-delete) status update
linear-cli project status delete <update-id>
linear-cli project status archive <update-id>  # Alias
linear-cli project status rm <update-id>       # Alias
```

### User Commands
```bash
# List all users in workspace
linear-cli user list [flags]
linear-cli user ls [flags]      # Alias
# Flags:
  -a, --active             Show only active users
  -l, --limit int          Maximum results (default 50)
  -o, --sort string        Sort order: linear (default), created, updated

# Examples:
linear-cli user list            # List all users
linear-cli user list --active   # List only active users

# Get user details by email
linear-cli user get <email>
linear-cli user show <email>    # Alias

# Examples:
linear-cli user get john@example.com
linear-cli user get jane.doe@company.com

# Show current authenticated user
linear-cli user me              # Shows your profile with admin status
```

### Document Commands
```bash
# List documents
linear-cli document list [flags]
linear-cli document ls [flags]     # Alias
# Flags:
  --project string           Filter by project ID
  -t, --team string          Filter by team key
  -l, --limit int            Maximum results (default 50)
  -o, --sort string          Sort order: linear (default), created, updated
  -n, --newer-than string    Show documents created after this time (default: 6_months_ago)

# Get document details (shows full content)
linear-cli document get <document-id>
linear-cli document show <document-id>  # Alias

# Search documents
linear-cli document search <query> [flags]
linear-cli document find <query> [flags]  # Alias
# Flags:
  -t, --team string          Filter by team ID
  -l, --limit int            Maximum results (default 50)
  -o, --sort string          Sort order: linear (default), created, updated
  --include-comments         Include document comments in search

# Create document
linear-cli document create [flags]
linear-cli document new [flags]    # Alias
# Flags:
  --title string             Document title (required)
  --content string           Document content (markdown)
  --project string           Project ID to associate with
  -t, --team string          Team key to associate with
  --icon string              Document icon (emoji)
  --color string             Document icon color (hex)

# Update document
linear-cli document update <document-id> [flags]
linear-cli document edit <document-id> [flags]  # Alias
# Flags:
  --title string             New title
  --content string           New content (markdown)
  --icon string              New icon (emoji)
  --color string             New icon color (hex)

# Delete document
linear-cli document delete <document-id>
```

### Comment Commands
```bash
# List all comments for an issue
linear-cli comment list <issue-id> [flags]
linear-cli comment ls <issue-id> [flags]    # Alias
# Flags:
  -l, --limit int          Maximum results (default 50)
  -o, --sort string        Sort order: linear (default), created, updated

# Examples:
linear-cli comment list LIN-123      # Shows all comments with timestamps
linear-cli comment list LIN-456 -l 10 # Show latest 10 comments

# Add comment to issue
linear-cli comment create <issue-id> --body "Comment text"
linear-cli comment add <issue-id> -b "Comment text"    # Alias
linear-cli comment new <issue-id> -b "Comment text"    # Alias

# Examples:
linear-cli comment create LIN-123 --body "I've started working on this"
linear-cli comment add LIN-123 -b "Fixed in commit abc123"
linear-cli comment create LIN-456 --body "@john please review this PR"
```

### 8. Initiatives
```bash
# List all initiatives
linear-cli initiative list

# List by status
linear-cli initiative list --status Active

# Include completed
linear-cli initiative list --include-completed

# Get initiative details (with linked projects and sub-initiatives)
linear-cli initiative get INITIATIVE-ID

# Create an initiative
linear-cli initiative create --name "Q1 Goals" --description "Company objectives for Q1"
linear-cli initiative create --name "Mobile Launch" --status Active --target-date "2025-06-01"

# Update an initiative
linear-cli initiative update INITIATIVE-ID --status Completed
linear-cli initiative update INITIATIVE-ID --name "Q2 Goals" --target-date "2025-09-01"

# Delete an initiative
linear-cli initiative delete INITIATIVE-ID
```

### 9. Custom Views (Saved Filters)
```bash
# List all custom views
linear-cli view list
linear-cli view list --shared          # Only shared views
linear-cli view list --model issue     # Only issue views
linear-cli view list --team ENG        # Only views for a team

# Get view details (shows filter configuration)
linear-cli view get VIEW-ID

# Run a view — execute its saved filters and see matching items
linear-cli view run VIEW-ID
linear-cli view run VIEW-ID --limit 100
linear-cli view run VIEW-ID --json

# Create a custom view
linear-cli view create --name "My Bugs" --model issue
linear-cli view create --name "Active Projects" --model project --shared
linear-cli view create --name "Urgent Issues" --filter-json '{"priority":{"eq":1}}'
linear-cli view create --name "Team WIP" --team ENG --filter-json '{"state":{"type":{"eq":"started"}}}'

# Update a view
linear-cli view update VIEW-ID --name "Renamed View"
linear-cli view update VIEW-ID --shared
linear-cli view update VIEW-ID --filter-json '{"priority":{"in":[1,2]}}'

# Delete a view
linear-cli view delete VIEW-ID
```

### Custom View Commands
```bash
# List custom views
linear-cli view list [flags]
linear-cli view ls [flags]     # Alias
# Flags:
  -l, --limit int          Maximum results (default 50)
  --shared                 Show only shared views
  -m, --model string       Filter by model type (issue, project)
  -t, --team string        Filter by team key

# Get view details
linear-cli view get <view-id>
linear-cli view show <view-id>  # Alias

# Run (execute) a view
linear-cli view run <view-id> [flags]
linear-cli view exec <view-id> [flags]  # Alias
# Flags:
  -l, --limit int          Maximum results (default 50)

# Create view
linear-cli view create [flags]
linear-cli view new [flags]    # Alias
# Flags:
  --name string            View name (required)
  -d, --description string View description
  -m, --model string       Model type: issue (default), project
  -t, --team string        Team key
  --shared                 Make the view shared
  --filter-json string     Raw JSON filter (IssueFilter or ProjectFilter schema)

# Update view
linear-cli view update <view-id> [flags]
linear-cli view edit <view-id> [flags]  # Alias
# Flags:
  --name string            New name
  -d, --description string New description
  --shared                 Set shared status
  --filter-json string     New raw JSON filter

# Delete view
linear-cli view delete <view-id>
linear-cli view rm <view-id>   # Alias
```

## 🎨 Output Formats

### Table Format (Default)
```bash
linear-cli issue list
```
```
ID       Title                State        Assignee    Team  Priority
LIN-123  Fix authentication   In Progress  john@co.com ENG   High
LIN-124  Update documentation Done         jane@co.com DOC   Normal
```

### Plaintext Format
```bash
linear-cli issue list --plaintext
```
```
# Issues
## BUG: Fix login button alignment
- **ID**: FAK-123
- **State**: In Progress
- **Assignee**: Jane Doe
- **Team**: WEB
- **Created**: 2025-07-12
- **URL**: https://linear.app/example/issue/FAK-123/bug-fix-login-button-alignment
- **Description**: The login button on the main page is misaligned on mobile devices.

Steps to reproduce:
1. Open the website on a mobile browser.
2. Navigate to the login page.
3. Observe the button alignment.

## FEAT: Add dark mode support
- **ID**: FAK-124
- **State**: Todo
- **Assignee**: John Smith
- **Team**: APP
- **Created**: 2025-07-11
- **URL**: https://linear.app/example/issue/FAK-124/feat-add-dark-mode-support
- **Description**: Implement a dark mode theme for the entire application to improve user experience in low-light environments.
```

### JSON Format
```bash
linear-cli issue list --json
```
```json
[
  {
    "id": "LIN-123",
    "title": "Fix authentication",
    "state": "In Progress",
    "assignee": "john@co.com",
    "team": "ENG",
    "priority": "High"
  }
]
```

## ⚙️ Configuration

Configuration is stored in `~/.linear-cli.yaml`:

```yaml
# Default output format
output: table

# Default pagination limit
limit: 50

# API settings
api:
  timeout: 30s
  retries: 3
```

Authentication credentials are stored securely in `~/.linear-cli-auth.json`.

## 🔒 Authentication

### Personal API Key (Recommended)
1. Go to [Linear Settings > API](https://linear.app/settings/api)
2. Create a new Personal API Key
3. Run `linear-cli auth` and paste your key

## 📅 Time-based Filtering

**⚠️ Default Behavior**: To improve performance and prevent overwhelming data loads, list commands **only show items created in the last 6 months by default**. This is especially important for large workspaces.

### Using the --newer-than Flag

The `--newer-than` (or `-n`) flag is available on `issue list` and `project list` commands:

```bash
# Default behavior (last 6 months)
linear-cli issue list

# Show items from a specific time period
linear-cli issue list --newer-than 2_weeks_ago
linear-cli project list --newer-than 1_month_ago

# Show ALL items regardless of age
linear-cli issue list --newer-than all_time
```

### Supported Time Formats

1. **Relative time expressions**: `N_units_ago`
   - Units: `minutes`, `hours`, `days`, `weeks`, `months`, `years`
   - Examples: `30_minutes_ago`, `2_hours_ago`, `3_days_ago`, `1_week_ago`, `6_months_ago`

2. **Special values**:
   - `all_time` - Shows all items without any date filter
   - ISO dates - `2025-07-01` or `2025-07-01T15:30:00Z`

3. **Default value**: `6_months_ago` (when flag is not specified)

### Quick Reference

| Time Expression | Description | Example Command |
|----------------|-------------|-----------------|
| *(no flag)* | Last 6 months (default) | `linear-cli issue list` |
| `1_day_ago` | Last 24 hours | `linear-cli issue list --newer-than 1_day_ago` |
| `1_week_ago` | Last 7 days | `linear-cli issue list --newer-than 1_week_ago` |
| `2_weeks_ago` | Last 14 days | `linear-cli issue list --newer-than 2_weeks_ago` |
| `1_month_ago` | Last month | `linear-cli issue list --newer-than 1_month_ago` |
| `3_months_ago` | Last quarter | `linear-cli issue list --newer-than 3_months_ago` |
| `6_months_ago` | Last 6 months | `linear-cli issue list --newer-than 6_months_ago` |
| `1_year_ago` | Last year | `linear-cli issue list --newer-than 1_year_ago` |
| `all_time` | No date filter | `linear-cli issue list --newer-than all_time` |
| `2025-07-01` | Since specific date | `linear-cli issue list --newer-than 2025-07-01` |

### Common Use Cases

```bash
# Recent activity - issues from last week
linear-cli issue list --newer-than 1_week_ago

# Sprint planning - issues from current month
linear-cli issue list --newer-than 1_month_ago --state "Todo"

# Quarterly review - all projects from last 3 months
linear-cli project list --newer-than 3_months_ago

# Historical analysis - ALL issues ever created
linear-cli issue list --newer-than all_time --sort created

# Today's issues
linear-cli issue list --newer-than 1_day_ago

# Combine with other filters
linear-cli issue list --newer-than 2_weeks_ago --assignee me --sort updated
```

## 🔄 Sorting Options

All list commands support sorting with the `--sort` or `-o` flag:

- **linear** (default): Linear's built-in sorting order (respects manual ordering in the UI)
- **created**: Sort by creation date (newest first)
- **updated**: Sort by last update date (most recently updated first)

### Examples
```bash
# Get recently updated issues
linear-cli issue list --sort updated

# Get oldest projects first
linear-cli project list --sort created

# Get recently joined users
linear-cli user list --sort created --active

# Get latest comments on an issue
linear-cli comment list LIN-123 --sort created

# Combine sorting with filters
linear-cli issue list --assignee me --state "In Progress" --sort updated

# Combine time filtering with sorting
linear-cli issue list --newer-than 1_week_ago --sort updated

# Get all projects sorted by creation date
linear-cli project list --newer-than all_time --sort created
```

### Performance Tips

- The 6-month default filter significantly improves performance for large workspaces
- Use specific time ranges when possible instead of `all_time`
- Combine time filtering with other filters (assignee, state, team) for faster results

## 🧪 Testing

linear-cli includes comprehensive unit and integration tests to ensure reliability.

### Running Tests
```bash
# Run all tests  (currently just a smoke test)
make test
```

### Integration Testing
Integration tests require a Linear API key. Create a `.env.test` file:
```bash
cp .env.test.example .env.test
# Edit .env.test and add your LINEAR_TEST_API_KEY
```

Or set it as an environment variable:
```bash
export LINEAR_TEST_API_KEY="your-test-api-key"
make test-integration
```

⚠️ **Note**: Integration tests are read-only and safe to run with production API keys.

### Test Structure
- `tests/unit/` - Unit tests with mocked API responses
- `tests/integration/` - End-to-end tests with real Linear API
- `tests/testutils/` - Shared test utilities and helpers

See [tests/README.md](tests/README.md) for detailed testing documentation.

## 🤖 Scripting & Automation

Use `--plaintext` or `--json` flags for scripting:

```bash
#!/bin/bash

# Get all urgent issues in JSON format
urgent_issues=$(linear-cli issue list --priority 1 --json)

# Parse with jq
echo "$urgent_issues" | jq '.[] | select(.assignee == "me") | .id'

# Plaintext output for simple parsing
linear-cli issue list --assignee me --plaintext | cut -f1 | tail -n +2

# Get issue count for different time periods
echo "Last week: $(linear-cli issue list --newer-than 1_week_ago --json | jq '. | length')"
echo "Last month: $(linear-cli issue list --newer-than 1_month_ago --json | jq '. | length')"
echo "All time: $(linear-cli issue list --newer-than all_time --json | jq '. | length')"

# Create and assign issue in one command
linear-cli issue create --title "Fix bug" --team ENG --assign-me --json

# Get all projects for a team
linear-cli project list --team ENG --json | jq '.[] | {name, progress}'

# List all admin users
linear-cli user list --json | jq '.[] | select(.admin == true) | {name, email}'

# Get team member count
linear-cli team members ENG --json | jq '. | length'

# Export issue comments
linear-cli comment list LIN-123 --json > issue-comments.json
```

## 📡 Real-World Examples

### Team Workflows
```bash
# Find which team a user belongs to
for team in $(linear-cli team list --json | jq -r '.[].key'); do
  echo "Checking team: $team"
  linear-cli team members $team --json | jq '.[] | select(.email == "john@example.com")'
done

# List all private teams
linear-cli team list --json | jq '.[] | select(.private == true) | {key, name}'

# Get teams with more than 50 issues
linear-cli team list --json | jq '.[] | select(.issueCount > 50) | {key, name, issueCount}'
```

### User Management
```bash
# Find inactive users
linear-cli user list --json | jq '.[] | select(.active == false) | {name, email}'

# Check if you're an admin
linear-cli user me --json | jq '.admin'

# List users who are admins but not the current user
linear-cli user list --json | jq '.[] | select(.admin == true and .isMe == false) | .email'
```

### Issue Comments
```bash
# Add a comment mentioning the issue is blocked
linear-cli comment create LIN-123 --body "Blocked by LIN-456. Waiting for API changes."

# Get all comments by a specific user
linear-cli comment list LIN-123 --json | jq '.[] | select(.user.email == "john@example.com") | .body'

# Count comments per issue
for issue in LIN-123 LIN-124 LIN-125; do
  count=$(linear-cli comment list $issue --json | jq '. | length')
  echo "$issue: $count comments"
done
```

### Project Tracking
```bash
# List projects nearing completion (>80% progress)
linear-cli project list --json | jq '.[] | select(.progress > 0.8) | {name, progress}'

# Get all paused projects
linear-cli project list --state paused

# Show project timeline
linear-cli project get PROJECT-ID --json | jq '{name, startDate, targetDate, progress}'
```

### Daily Standup Helper
```bash
#!/bin/bash
# Show my recent activity
echo "=== My Issues ==="
linear-cli issue list --assignee me --limit 10

echo -e "\n=== Recent Comments ==="
for issue in $(linear-cli issue list --assignee me --json | jq -r '.[].identifier'); do
  echo "Comments on $issue:"
  linear-cli comment list $issue --limit 3
done
```

## 🐛 Troubleshooting

### Authentication Issues
```bash
# Check authentication status
linear-cli auth status

# Re-authenticate
linear-cli auth logout
linear-cli auth
```

### API Rate Limits
Linear has the following rate limits:
- Personal API Keys: 5,000 requests/hour

### Common Errors
- `Not authenticated`: Run `linear-cli auth` first
- `Team not found`: Use team key (e.g., "ENG") not display name
- `Invalid priority`: Use numbers 0-4 (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)

### Time Filtering Issues
- **Missing old issues?** Remember that list commands default to showing only the last 6 months
  - Solution: Use `--newer-than all_time` to see all issues
- **Invalid time expression?** Check the format: `N_units_ago` (e.g., `3_weeks_ago`)
  - Valid units: `minutes`, `hours`, `days`, `weeks`, `months`, `years`
- **Performance issues?** Avoid using `all_time` on large workspaces
  - Solution: Use specific time ranges like `--newer-than 1_year_ago`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

See CONTRIBUTING.md for a detailed release checklist and the Homebrew tap auto-bump workflow.

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Linear API Documentation](https://developers.linear.app/)
- [GitHub Repository](https://github.com/dorkitude/linear-cli)
- [Issue Tracker](https://github.com/dorkitude/linear-cli/issues)
- [Original Project (linctl)](https://github.com/dorkitude/linctl) - the upstream project this was forked from

---

**Forked from [linctl](https://github.com/dorkitude/linctl). Built with Go, Cobra, and the Linear API.**
