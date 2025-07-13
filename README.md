# 🚀 linctl - Linear CLI Tool

A comprehensive command-line interface for Linear's API, built with Go and Cobra.

## ✨ Features

- 🔐 **Authentication**: Personal API Key and OAuth 2.0 support
- 📋 **Issue Management**: Create, list, update, and archive issues
- 👥 **Team Management**: List teams, members, and team details
- 🚀 **Project Tracking**: Manage projects and milestones
- 💬 **Comments**: Add and view issue comments
- 📎 **Attachments**: Handle file uploads and attachments
- 🔗 **Webhooks**: Configure and manage webhooks
- 🎨 **Multiple Output Formats**: Table, plaintext, and JSON output
- ⚡ **Performance**: Fast and lightweight CLI tool

## 🛠️ Installation

### Homebrew (macOS/Linux)
```bash
brew install linctl
```

### apt (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install linctl
```

### From Source
```bash
git clone https://github.com/dorkitude/linctl.git
cd linctl
go build -o linctl
sudo mv linctl /usr/local/bin/
```

## 🚀 Quick Start

### 1. Authentication
```bash
# Interactive authentication
linctl auth

# Check authentication status
linctl auth status

# Show current user
linctl whoami
```

### 2. Issue Management
```bash
# List all issues
linctl issue list

# List issues assigned to you
linctl issue list --assignee me

# List issues in a specific state
linctl issue list --state "In Progress"

# Get issue details
linctl issue get LIN-123

# Create a new issue
linctl issue create --title "Bug fix" --team ENG
```

### 3. Team Management
```bash
# List all teams
linctl team list

# Get team details
linctl team get ENG

# List team members
linctl team members ENG
```

## 📖 Command Reference

### Global Flags
- `--plaintext, -p`: Plain text output (non-interactive)
- `--json, -j`: JSON output for scripting
- `--help, -h`: Show help
- `--version, -v`: Show version

### Authentication Commands
```bash
linctl auth               # Interactive authentication
linctl auth login         # Same as above
linctl auth status        # Check authentication status
linctl auth logout        # Clear stored credentials
linctl whoami            # Show current user
```

### Issue Commands
```bash
# List issues with filters
linctl issue list [flags]
linctl issue ls [flags]     # Short alias

# Flags:
  -a, --assignee string    Filter by assignee (email or 'me')
  -s, --state string       Filter by state name
  -t, --team string        Filter by team key
  -r, --priority int       Filter by priority (0-4)
  -l, --limit int          Maximum results (default 50)

# Get issue details
linctl issue get <issue-id>
linctl issue show <issue-id>  # Alias

# Create issue
linctl issue create [flags]
linctl issue new [flags]      # Alias

# Update issue
linctl issue update <issue-id> [flags]
linctl issue edit <issue-id> [flags]    # Alias

# Archive issue
linctl issue archive <issue-id>
```

### Team Commands
```bash
# List teams
linctl team list
linctl team ls              # Alias

# Get team details
linctl team get <team-key>
linctl team show <team-key> # Alias

# List team members
linctl team members <team-key>
```

### Project Commands
```bash
# List projects
linctl project list
linctl project ls           # Alias

# Get project details
linctl project get <project-id>
linctl project show <project-id>  # Alias

# Create project
linctl project create [flags]
```

### User Commands
```bash
# List users
linctl user list
linctl user ls              # Alias

# Get user details
linctl user get <email>
linctl user show <email>    # Alias

# Show current user
linctl user me
```

### Comment Commands
```bash
# List comments for issue
linctl comment list <issue-id>
linctl comment ls <issue-id>    # Alias

# Add comment to issue
linctl comment create <issue-id> --body "Comment text"
linctl comment add <issue-id> -b "Comment text"  # Alias
```

## 🎨 Output Formats

### Table Format (Default)
```bash
linctl issue list
```
```
ID       Title                State        Assignee    Team  Priority
LIN-123  Fix authentication   In Progress  john@co.com ENG   High
LIN-124  Update documentation Done         jane@co.com DOC   Normal
```

### Plaintext Format
```bash
linctl issue list --plaintext
```
```
ID	Title	State	Assignee	Team	Priority
LIN-123	Fix authentication	In Progress	john@co.com	ENG	High
LIN-124	Update documentation	Done	jane@co.com	DOC	Normal
```

### JSON Format
```bash
linctl issue list --json
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

Configuration is stored in `~/.linctl.yaml`:

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

Authentication credentials are stored securely in `~/.linctl-auth.json`.

## 🔒 Authentication

### Personal API Key (Recommended)
1. Go to [Linear Settings > API](https://linear.app/settings/api)
2. Create a new Personal API Key
3. Run `linctl auth` and paste your key

### OAuth 2.0 (Future)
OAuth 2.0 flow will be available in a future release for building applications.

## 🤖 Scripting & Automation

Use `--plaintext` or `--json` flags for scripting:

```bash
#!/bin/bash

# Get all urgent issues in JSON format
urgent_issues=$(linctl issue list --priority 1 --json)

# Parse with jq
echo "$urgent_issues" | jq '.[] | select(.assignee == "me") | .id'

# Plaintext output for simple parsing
linctl issue list --assignee me --plaintext | cut -f1 | tail -n +2
```

## 🐛 Troubleshooting

### Authentication Issues
```bash
# Check authentication status
linctl auth status

# Re-authenticate
linctl auth logout
linctl auth
```

### API Rate Limits
Linear has the following rate limits:
- Personal API Keys: 5,000 requests/hour
- OAuth Apps: 15,000 requests/hour

### Common Errors
- `Not authenticated`: Run `linctl auth` first
- `Team not found`: Use team key (e.g., "ENG") not display name
- `Invalid priority`: Use numbers 0-4 (0=None, 1=Urgent, 2=High, 3=Normal, 4=Low)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Linear API Documentation](https://developers.linear.app/)
- [GitHub Repository](https://github.com/dorkitude/linctl)
- [Issue Tracker](https://github.com/dorkitude/linctl/issues)

---

**Built with ❤️ using Go, Cobra, and the Linear API**