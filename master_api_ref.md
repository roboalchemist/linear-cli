# Linear API Master Reference

> Comprehensive reference for Linear's GraphQL API endpoints and features for the `linear-cli` CLI tool

## Table of Contents
- [Authentication](#authentication)
- [Core GraphQL Operations](#core-graphql-operations)
- [Issues](#issues)
- [Projects](#projects)
- [Teams](#teams)
- [Users](#users)
- [Comments](#comments)
- [Attachments](#attachments)
- [Webhooks](#webhooks)
- [Pagination & Filtering](#pagination--filtering)
- [Rate Limiting](#rate-limiting)
- [CLI Command Mapping](#cli-command-mapping)

## Authentication

### Methods
1. **Personal API Keys** (recommended for CLI)
   - Header: `Authorization: <API_KEY>`
   - Created in "Security & access" settings
   
2. **OAuth 2.0** (for web applications)
   - Header: `Authorization: Bearer <ACCESS_TOKEN>`
   - Requires app registration and token exchange

### Scopes
- `read` - Read access to all resources
- `write` - Write access to all resources
- `issues:create` - Create issues
- `issues:write` - Modify issues
- `teams:read` - Read team information
- `admin` - Administrative access

## Core GraphQL Operations

### Endpoint
```
https://api.linear.app/graphql
```

### Introspection
```graphql
query IntrospectionQuery {
  __schema {
    types {
      name
      description
    }
  }
}
```

## Issues

### List Issues
```graphql
query Issues(
  $filter: IssueFilter,
  $orderBy: PaginationOrderBy,
  $first: Int,
  $after: String
) {
  issues(
    filter: $filter,
    orderBy: $orderBy,
    first: $first,
    after: $after
  ) {
    nodes {
      id
      identifier
      title
      description
      priority
      state {
        name
        type
      }
      assignee {
        name
        email
      }
      team {
        name
        key
      }
      createdAt
      updatedAt
      dueDate
      estimate
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

### Get Single Issue
```graphql
query Issue($id: String!) {
  issue(id: $id) {
    id
    identifier
    title
    description
    priority
    state {
      name
      type
    }
    assignee {
      name
      email
    }
    team {
      name
      key
    }
    labels {
      nodes {
        name
        color
      }
    }
    comments {
      nodes {
        body
        user {
          name
        }
        createdAt
      }
    }
    attachments {
      nodes {
        title
        url
      }
    }
    createdAt
    updatedAt
    dueDate
    estimate
  }
}
```

### Create Issue
```graphql
mutation IssueCreate($input: IssueCreateInput!) {
  issueCreate(input: $input) {
    success
    issue {
      id
      identifier
      title
    }
  }
}
```

### Update Issue
```graphql
mutation IssueUpdate($id: String!, $input: IssueUpdateInput!) {
  issueUpdate(id: $id, input: $input) {
    success
    issue {
      id
      identifier
      title
    }
  }
}
```

### Archive Issue
```graphql
mutation IssueArchive($id: String!) {
  issueArchive(id: $id) {
    success
  }
}
```

## Projects

### List Projects
```graphql
query Projects($filter: ProjectFilter, $first: Int, $after: String) {
  projects(filter: $filter, first: $first, after: $after) {
    nodes {
      id
      name
      description
      state
      progress
      startDate
      targetDate
      lead {
        name
      }
      members {
        nodes {
          name
        }
      }
      teams {
        nodes {
          name
          key
        }
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

### Create Project
```graphql
mutation ProjectCreate($input: ProjectCreateInput!) {
  projectCreate(input: $input) {
    success
    project {
      id
      name
    }
  }
}
```

### Update Project
```graphql
mutation ProjectUpdate($id: String!, $input: ProjectUpdateInput!) {
  projectUpdate(id: $id, input: $input) {
    success
    project {
      id
      name
    }
  }
}
```

## Teams

### List Teams
```graphql
query Teams($filter: TeamFilter, $first: Int, $after: String) {
  teams(filter: $filter, first: $first, after: $after) {
    nodes {
      id
      key
      name
      description
      private
      issueCount
      members {
        nodes {
          name
          email
        }
      }
      states {
        nodes {
          name
          type
          color
        }
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

### Get Team
```graphql
query Team($id: String!) {
  team(id: $id) {
    id
    key
    name
    description
    private
    issueCount
    members {
      nodes {
        id
        name
        email
        isMe
      }
    }
    issues(first: 50) {
      nodes {
        identifier
        title
        state {
          name
        }
        assignee {
          name
        }
      }
    }
  }
}
```

## Users

### Current User (Viewer)
```graphql
query Viewer {
  viewer {
    id
    name
    email
    avatarUrl
    isMe
    teams {
      nodes {
        name
        key
      }
    }
    assignedIssues(first: 50) {
      nodes {
        identifier
        title
        state {
          name
        }
      }
    }
  }
}
```

### List Users
```graphql
query Users($filter: UserFilter, $first: Int, $after: String) {
  users(filter: $filter, first: $first, after: $after) {
    nodes {
      id
      name
      email
      avatarUrl
      active
      admin
      guest
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

## Comments

### List Comments for Issue
```graphql
query Comments($issueId: String!, $first: Int, $after: String) {
  issue(id: $issueId) {
    comments(first: $first, after: $after) {
      nodes {
        id
        body
        user {
          name
          email
        }
        createdAt
        updatedAt
      }
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
}
```

### Create Comment
```graphql
mutation CommentCreate($input: CommentCreateInput!) {
  commentCreate(input: $input) {
    success
    comment {
      id
      body
    }
  }
}
```

## Attachments

### Create Attachment
```graphql
mutation AttachmentCreate($input: AttachmentCreateInput!) {
  attachmentCreate(input: $input) {
    success
    attachment {
      id
      title
      url
    }
  }
}
```

## Webhooks

### List Webhooks
```graphql
query Webhooks($first: Int, $after: String) {
  webhooks(first: $first, after: $after) {
    nodes {
      id
      url
      label
      enabled
      resourceTypes
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
```

### Create Webhook
```graphql
mutation WebhookCreate($input: WebhookCreateInput!) {
  webhookCreate(input: $input) {
    success
    webhook {
      id
      url
      label
    }
  }
}
```

## Pagination & Filtering

### Pagination Arguments
- `first: Int` - Number of items to fetch
- `after: String` - Cursor for pagination
- `orderBy: PaginationOrderBy` - Sort order (createdAt, updatedAt)

### Common Filters

#### Issue Filters
```graphql
{
  team: { id: { eq: "TEAM_ID" } }
  assignee: { id: { eq: "USER_ID" } }
  state: { name: { eq: "In Progress" } }
  priority: { eq: 1 }  # 0=None, 1=Urgent, 2=High, 3=Normal, 4=Low
  createdAt: { gte: "2024-01-01T00:00:00Z" }
}
```

#### String Comparators
- `eq` - equals
- `neq` - not equals
- `contains` - contains substring
- `containsIgnoreCase` - case-insensitive contains
- `startsWith` - starts with
- `endsWith` - ends with

## Rate Limiting

### Limits
- **Personal API Keys**: 5,000 requests per hour
- **OAuth Apps**: 15,000 requests per hour
- **Complexity Limit**: 1,000,000 per request

### Headers
- `X-RateLimit-Limit` - Request limit
- `X-RateLimit-Remaining` - Remaining requests
- `X-RateLimit-Reset` - Reset timestamp

## CLI Command Mapping

### Issue Commands
```bash
# List issues
linear-cli issue list --assignee me --state "In Progress"
linear-cli issue ls -a me -s "In Progress"

# Get specific issue
linear-cli issue get LIN-123
linear-cli issue show LIN-123

# Create issue
linear-cli issue create --title "Bug fix" --team TEAM_KEY
linear-cli issue new -t "Bug fix" --team TEAM_KEY

# Update issue
linear-cli issue update LIN-123 --assignee user@example.com
linear-cli issue edit LIN-123 -a user@example.com

# Archive issue
linear-cli issue archive LIN-123
```

### Project Commands
```bash
# List projects
linear-cli project list --team TEAM_KEY
linear-cli project ls -t TEAM_KEY

# Get project
linear-cli project get PROJECT_ID
linear-cli project show PROJECT_ID

# Create project
linear-cli project create --name "New Feature" --team TEAM_KEY
```

### Team Commands
```bash
# List teams
linear-cli team list
linear-cli team ls

# Get team info
linear-cli team get TEAM_KEY
linear-cli team show TEAM_KEY

# List team members
linear-cli team members TEAM_KEY
```

### User Commands
```bash
# Show current user
linear-cli user me
linear-cli whoami

# List users
linear-cli user list
linear-cli user ls

# Show user info
linear-cli user get user@example.com
```

### Comment Commands
```bash
# List comments
linear-cli comment list LIN-123
linear-cli comment ls LIN-123

# Add comment
linear-cli comment create LIN-123 --body "Comment text"
linear-cli comment add LIN-123 -b "Comment text"
```

### Auth Commands
```bash
# Authenticate
linear-cli auth
linear-cli auth login

# Show current auth status
linear-cli auth status
linear-cli auth whoami

# Logout
linear-cli auth logout
```

### Global Flags
- `--plaintext, -p` - Plain text output (non-interactive)
- `--json, -j` - JSON output
- `--help, -h` - Show help
- `--version, -v` - Show version

### Output Formats
1. **Table** (default) - Formatted table with colors
2. **Plaintext** (`-p`) - Simple text output for scripts
3. **JSON** (`-j`) - Structured data for parsing

---

**Note**: All commands support both long form (`issue list`) and short form (`issue ls`) for better UX. The `--plaintext` flag ensures output is suitable for automation and other CLI tools.