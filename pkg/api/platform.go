package api

import (
	"context"
	"time"
)

// ============================================================================
// Platform types (agent sessions/skills, integrations, external users,
// organization, initiative updates, semantic search, notification
// subscriptions).
//
// NOTE: ExternalUser, Issues and Projects already exist in queries.go and are
// reused where they fit. The types below are new and must never duplicate an
// existing name.
// ============================================================================

// AgentSession represents an agent (coding) session.
type AgentSession struct {
	ID         string     `json:"id"`
	Summary    string     `json:"summary"`
	Status     string     `json:"status"`
	SlugID     string     `json:"slugId"`
	URL        *string    `json:"url"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	StartedAt  *time.Time `json:"startedAt"`
	EndedAt    *time.Time `json:"endedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
	AppUser    *User      `json:"appUser"`
	Creator    *User      `json:"creator"`
	Issue      *Issue     `json:"issue"`
}

// AgentSessions is a paginated collection of agent sessions.
type AgentSessions struct {
	Nodes    []AgentSession `json:"nodes"`
	PageInfo PageInfo       `json:"pageInfo"`
}

// AgentSkill represents a workspace agent skill.
type AgentSkill struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	SlugID           string     `json:"slugId"`
	Description      *string    `json:"description"`
	TeamID           *string    `json:"teamId"`
	Shared           bool       `json:"shared"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	ArchivedAt       *time.Time `json:"archivedAt"`
	LastUsedAt       *time.Time `json:"lastUsedAt"`
	RecentUsageCount float64    `json:"recentUsageCount"`
	Owner            *User      `json:"owner"`
	Creator          *User      `json:"creator"`
}

// AgentSkills is a paginated collection of agent skills.
type AgentSkills struct {
	Nodes    []AgentSkill `json:"nodes"`
	PageInfo PageInfo     `json:"pageInfo"`
}

// Integration represents a workspace integration (Slack, GitHub, Jira, ...).
type Integration struct {
	ID           string        `json:"id"`
	Service      string        `json:"service"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	ArchivedAt   *time.Time    `json:"archivedAt"`
	Creator      *User         `json:"creator"`
	Team         *Team         `json:"team"`
	Organization *Organization `json:"organization"`
}

// Integrations is a paginated collection of integrations.
type Integrations struct {
	Nodes    []Integration `json:"nodes"`
	PageInfo PageInfo      `json:"pageInfo"`
}

// ExternalUserRecord is the wider representation of an ExternalUser used for
// listing. The pre-existing api.ExternalUser type only carries id/name/email,
// so CreatedAt/LastSeen/displayName need this dedicated type.
type ExternalUserRecord struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	DisplayName string     `json:"displayName"`
	Email       *string    `json:"email"`
	AvatarURL   *string    `json:"avatarUrl"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	LastSeen    *time.Time `json:"lastSeen"`
	ArchivedAt  *time.Time `json:"archivedAt"`
}

// ExternalUserRecords is a paginated collection of external users.
type ExternalUserRecords struct {
	Nodes    []ExternalUserRecord `json:"nodes"`
	PageInfo PageInfo             `json:"pageInfo"`
}

// Organization represents the authenticated user's workspace.
type Organization struct {
	ID                   string     `json:"id"`
	Name                 string     `json:"name"`
	URLKey               string     `json:"urlKey"`
	LogoURL              *string    `json:"logoUrl"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	ArchivedAt           *time.Time `json:"archivedAt"`
	TrialEndsAt          *time.Time `json:"trialEndsAt"`
	UserCount            int        `json:"userCount"`
	CreatedIssueCount    int        `json:"createdIssueCount"`
	CustomerCount        int        `json:"customerCount"`
	PreviousURLKeys      []string   `json:"previousUrlKeys"`
	ReleaseChannel       string     `json:"releaseChannel"`
	SAMLEnabled          bool       `json:"samlEnabled"`
	SCIMEnabled          bool       `json:"scimEnabled"`
	FiscalYearStartMonth float64    `json:"fiscalYearStartMonth"`
}

// InitiativeUpdate represents a status update posted to an initiative.
type InitiativeUpdate struct {
	ID           string      `json:"id"`
	Body         string      `json:"body"`
	Health       string      `json:"health"`
	SlugID       string      `json:"slugId"`
	IsStale      bool        `json:"isStale"`
	CommentCount int         `json:"commentCount"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	EditedAt     *time.Time  `json:"editedAt"`
	ArchivedAt   *time.Time  `json:"archivedAt"`
	Initiative   *Initiative `json:"initiative"`
}

// InitiativeUpdates is a paginated collection of initiative updates.
type InitiativeUpdates struct {
	Nodes    []InitiativeUpdate `json:"nodes"`
	PageInfo PageInfo           `json:"pageInfo"`
}

// ProjectSearchResult is a single hit from searchProjects. The GraphQL node
// type is ProjectSearchResult (not Project), so this dedicated type is used.
type ProjectSearchResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	Progress    float64   `json:"progress"`
	Health      string    `json:"health"`
	URL         string    `json:"url"`
	Icon        *string   `json:"icon"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	TargetDate  *string   `json:"targetDate"`
	Lead        *User     `json:"lead"`
	Teams       *Teams    `json:"teams"`
	Status      *struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"status"`
}

// ProjectSearchResults is a paginated collection of project search hits.
type ProjectSearchResults struct {
	Nodes    []ProjectSearchResult `json:"nodes"`
	PageInfo PageInfo              `json:"pageInfo"`
}

// SemanticSearchResult is one node of a semanticSearch response.
type SemanticSearchResult struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Issue      *Issue      `json:"issue"`
	Project    *Project    `json:"project"`
	Initiative *Initiative `json:"initiative"`
	Document   *Document   `json:"document"`
}

// SemanticSearchPayload is the result of semanticSearch.
type SemanticSearchPayload struct {
	Results []SemanticSearchResult `json:"results"`
}

// NotificationSubscription represents a single notification subscription.
// The GraphQL type is an interface, but all selected fields live on the
// interface itself, so no inline fragments are required.
type NotificationSubscription struct {
	ID                          string      `json:"id"`
	Active                      bool        `json:"active"`
	IncludeSubInitiativeUpdates bool        `json:"includeSubInitiativeUpdates"`
	CreatedAt                   time.Time   `json:"createdAt"`
	UpdatedAt                   time.Time   `json:"updatedAt"`
	ArchivedAt                  *time.Time  `json:"archivedAt"`
	Subscriber                  *User       `json:"subscriber"`
	Team                        *Team       `json:"team"`
	Project                     *Project    `json:"project"`
	Cycle                       *Cycle      `json:"cycle"`
	Initiative                  *Initiative `json:"initiative"`
	Label                       *Label      `json:"label"`
	User                        *User       `json:"user"`
}

// NotificationSubscriptions is a paginated collection of notification
// subscriptions.
type NotificationSubscriptions struct {
	Nodes    []NotificationSubscription `json:"nodes"`
	PageInfo PageInfo                   `json:"pageInfo"`
}

const agentSessionFields = `
	id
	summary
	status
	slugId
	url
	createdAt
	updatedAt
	startedAt
	endedAt
	archivedAt
	appUser {
		id
		name
	}
	creator {
		id
		name
	}
	issue {
		id
		identifier
		title
		state {
			id
			name
			type
		}
	}
`

// GetAgentSessions returns agent sessions, most recent first by default.
func (c *Client) GetAgentSessions(ctx context.Context, first int, after string, orderBy string, includeArchived bool) (*AgentSessions, error) {
	query := `
		query AgentSessions($first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			agentSessions(first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {` + agentSessionFields + `}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first":           first,
		"includeArchived": includeArchived,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		AgentSessions AgentSessions `json:"agentSessions"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.AgentSessions, nil
}

// GetAgentSession returns a single agent session by ID.
func (c *Client) GetAgentSession(ctx context.Context, id string) (*AgentSession, error) {
	query := `
		query AgentSession($id: String!) {
			agentSession(id: $id) {
				` + agentSessionFields + `
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		AgentSession AgentSession `json:"agentSession"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.AgentSession, nil
}

// GetAgentSkills returns all agent skills.
func (c *Client) GetAgentSkills(ctx context.Context, first int, after string, orderBy string) (*AgentSkills, error) {
	query := `
		query AgentSkills($first: Int, $after: String, $orderBy: PaginationOrderBy) {
			agentSkills(first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					title
					slugId
					description
					teamId
					shared
					createdAt
					updatedAt
					archivedAt
					lastUsedAt
					recentUsageCount
					owner {
						id
						name
					}
					creator {
						id
						name
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		AgentSkills AgentSkills `json:"agentSkills"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.AgentSkills, nil
}

// GetIntegrations returns all integrations for the workspace.
func (c *Client) GetIntegrations(ctx context.Context, first int, after string, orderBy string) (*Integrations, error) {
	query := `
		query Integrations($first: Int, $after: String, $orderBy: PaginationOrderBy) {
			integrations(first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					service
					createdAt
					updatedAt
					archivedAt
					creator {
						id
						name
					}
					team {
						id
						key
						name
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Integrations Integrations `json:"integrations"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.Integrations, nil
}

// GetExternalUsers returns all external users for the organization.
func (c *Client) GetExternalUsers(ctx context.Context, first int, after string, orderBy string) (*ExternalUserRecords, error) {
	query := `
		query ExternalUsers($first: Int, $after: String, $orderBy: PaginationOrderBy) {
			externalUsers(first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					name
					displayName
					email
					avatarUrl
					createdAt
					updatedAt
					lastSeen
					archivedAt
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		ExternalUsers ExternalUserRecords `json:"externalUsers"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.ExternalUsers, nil
}

// GetOrganization returns the authenticated user's workspace.
func (c *Client) GetOrganization(ctx context.Context) (*Organization, error) {
	query := `
		query Organization {
			organization {
				id
				name
				urlKey
				logoUrl
				createdAt
				updatedAt
				archivedAt
				trialEndsAt
				userCount
				createdIssueCount
				customerCount
				previousUrlKeys
				releaseChannel
				samlEnabled
				scimEnabled
				fiscalYearStartMonth
			}
		}
	`

	var response struct {
		Organization Organization `json:"organization"`
	}

	if err := c.Execute(ctx, query, nil, &response); err != nil {
		return nil, err
	}
	return &response.Organization, nil
}

// GetInitiativeUpdates returns initiative status updates, optionally scoped to
// a single initiative.
func (c *Client) GetInitiativeUpdates(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*InitiativeUpdates, error) {
	query := `
		query InitiativeUpdates($filter: InitiativeUpdateFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			initiativeUpdates(filter: $filter, first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					body
					health
					slugId
					isStale
					commentCount
					createdAt
					updatedAt
					editedAt
					archivedAt
					initiative {
						id
						name
						url
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first": first,
	}
	if filter != nil {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		InitiativeUpdates InitiativeUpdates `json:"initiativeUpdates"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.InitiativeUpdates, nil
}

// SearchProjects searches projects by text query using full-text and vector
// search. Results are ranked by relevance.
func (c *Client) SearchProjects(ctx context.Context, term string, first int, after string, orderBy string, includeArchived bool) (*ProjectSearchResults, error) {
	query := `
		query ProjectSearch($term: String!, $first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			searchProjects(term: $term, first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {
					id
					name
					description
					state
					progress
					health
					url
					icon
					createdAt
					updatedAt
					targetDate
					lead {
						id
						name
					}
					teams {
						nodes {
							id
							key
							name
						}
					}
					status {
						name
						type
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"term":            term,
		"first":           first,
		"includeArchived": includeArchived,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		SearchProjects ProjectSearchResults `json:"searchProjects"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.SearchProjects, nil
}

// SemanticSearch searches issues, projects, initiatives and documents using
// natural language.
func (c *Client) SemanticSearch(ctx context.Context, query string, types []string, maxResults int, includeArchived bool) (*SemanticSearchPayload, error) {
	gqlQuery := `
		query SemanticSearch($query: String!, $maxResults: Int, $types: [SemanticSearchResultType!], $includeArchived: Boolean) {
			semanticSearch(query: $query, maxResults: $maxResults, types: $types, includeArchived: $includeArchived) {
				results {
					id
					type
					issue {
						id
						identifier
						title
						state {
							name
						}
					}
					project {
						id
						name
						state
					}
					initiative {
						id
						name
					}
					document {
						id
						title
						url
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"query":           query,
		"includeArchived": includeArchived,
	}
	if maxResults > 0 {
		variables["maxResults"] = maxResults
	}
	if len(types) > 0 {
		variables["types"] = types
	}

	var response struct {
		SemanticSearch SemanticSearchPayload `json:"semanticSearch"`
	}

	if err := c.Execute(ctx, gqlQuery, variables, &response); err != nil {
		return nil, err
	}
	return &response.SemanticSearch, nil
}

// GetNotificationSubscriptions returns the authenticated user's notification
// subscriptions.
func (c *Client) GetNotificationSubscriptions(ctx context.Context, first int, after string, orderBy string, includeArchived bool) (*NotificationSubscriptions, error) {
	query := `
		query NotificationSubscriptions($first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			notificationSubscriptions(first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {
					id
					active
					includeSubInitiativeUpdates
					createdAt
					updatedAt
					archivedAt
					subscriber {
						id
						name
					}
					team {
						id
						key
						name
					}
					project {
						id
						name
					}
					cycle {
						id
						number
						name
					}
					initiative {
						id
						name
					}
					label {
						id
						name
					}
					user {
						id
						name
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first":           first,
		"includeArchived": includeArchived,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		NotificationSubscriptions NotificationSubscriptions `json:"notificationSubscriptions"`
	}

	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.NotificationSubscriptions, nil
}
