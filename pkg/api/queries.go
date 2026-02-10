package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// User represents a Linear user
type User struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	AvatarURL   string     `json:"avatarUrl"`
	DisplayName string     `json:"displayName"`
	IsMe        bool       `json:"isMe"`
	Active      bool       `json:"active"`
	Admin       bool       `json:"admin"`
	CreatedAt   *time.Time `json:"createdAt"`
}

// Team represents a Linear team
type Team struct {
	ID          string  `json:"id"`
	Key         string  `json:"key"`
	Name        string  `json:"name"`
	DisplayName string  `json:"displayName"`
	Description string  `json:"description"`
	Icon        *string `json:"icon"`
	Color       string  `json:"color"`
	Private     bool    `json:"private"`
	IssueCount  int     `json:"issueCount"`
	Timezone    string  `json:"timezone"`
	// Cycle settings
	CyclesEnabled                bool   `json:"cyclesEnabled"`
	CycleStartDay                int    `json:"cycleStartDay"`
	CycleDuration                int    `json:"cycleDuration"`
	CycleCooldownTime            int    `json:"cycleCooldownTime"`
	CycleIssueAutoAssignStarted  bool   `json:"cycleIssueAutoAssignStarted"`
	CycleIssueAutoAssignCompleted bool   `json:"cycleIssueAutoAssignCompleted"`
	CycleLockToActive            bool   `json:"cycleLockToActive"`
	UpcomingCycleCount           int    `json:"upcomingCycleCount"`
	CycleCalenderUrl             string `json:"cycleCalenderUrl"`
	// Triage settings
	TriageEnabled              bool `json:"triageEnabled"`
	RequirePriorityToLeaveTriage bool `json:"requirePriorityToLeaveTriage"`
	// Issue estimation settings
	InheritIssueEstimation       bool    `json:"inheritIssueEstimation"`
	IssueEstimationType          string  `json:"issueEstimationType"`
	IssueEstimationAllowZero     bool    `json:"issueEstimationAllowZero"`
	IssueEstimationExtended      bool    `json:"issueEstimationExtended"`
	DefaultIssueEstimate         float64 `json:"defaultIssueEstimate"`
	SetIssueSortOrderOnStateChange string `json:"setIssueSortOrderOnStateChange"`
	// Auto-close/archive settings
	AutoClosePeriod      *float64 `json:"autoClosePeriod"`
	AutoCloseStateId     *string  `json:"autoCloseStateId"`
	AutoArchivePeriod    float64  `json:"autoArchivePeriod"`
	AutoCloseParentIssues *bool   `json:"autoCloseParentIssues"`
	AutoCloseChildIssues  *bool   `json:"autoCloseChildIssues"`
	// Team hierarchy
	Parent   *Team `json:"parent"`
	// Workflow settings
	InheritWorkflowStatuses bool `json:"inheritWorkflowStatuses"`
	GroupIssueHistory       bool `json:"groupIssueHistory"`
	// Access/membership settings
	JoinByDefault     *bool `json:"joinByDefault"`
	AllMembersCanJoin *bool `json:"allMembersCanJoin"`
	ScimManaged       bool  `json:"scimManaged"`
	ScimGroupName     *string `json:"scimGroupName"`
	// AI settings
	AiThreadSummariesEnabled     bool `json:"aiThreadSummariesEnabled"`
	AiDiscussionSummariesEnabled bool `json:"aiDiscussionSummariesEnabled"`
	// Timestamps
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
	RetiredAt  *time.Time `json:"retiredAt"`
	// Nested objects (for detail queries)
	DefaultIssueState             *WorkflowState `json:"defaultIssueState"`
	TriageIssueState              *WorkflowState `json:"triageIssueState"`
	MarkedAsDuplicateWorkflowState *WorkflowState `json:"markedAsDuplicateWorkflowState"`
	ActiveCycle                   *Cycle         `json:"activeCycle"`
	DefaultTemplateForMembers     *Template      `json:"defaultTemplateForMembers"`
	DefaultTemplateForNonMembers  *Template      `json:"defaultTemplateForNonMembers"`
	DefaultProjectTemplate        *Template      `json:"defaultProjectTemplate"`
}

// Issue represents a Linear issue
type Issue struct {
	ID                  string            `json:"id"`
	Identifier          string            `json:"identifier"`
	Title               string            `json:"title"`
	Description         string            `json:"description"`
	Priority            int               `json:"priority"`
	Estimate            *float64          `json:"estimate"`
	CreatedAt           time.Time         `json:"createdAt"`
	UpdatedAt           time.Time         `json:"updatedAt"`
	DueDate             *string           `json:"dueDate"`
	State               *State            `json:"state"`
	Assignee            *User             `json:"assignee"`
	Team                *Team             `json:"team"`
	Labels              *Labels           `json:"labels"`
	Children            *Issues           `json:"children"`
	Parent              *Issue            `json:"parent"`
	URL                 string            `json:"url"`
	BranchName          string            `json:"branchName"`
	Cycle               *Cycle            `json:"cycle"`
	Project             *Project          `json:"project"`
	ProjectMilestone    *ProjectMilestone `json:"projectMilestone"`
	Attachments         *Attachments      `json:"attachments"`
	Documents           *Documents        `json:"documents"`
	Comments            *Comments         `json:"comments"`
	SnoozedUntilAt      *time.Time        `json:"snoozedUntilAt"`
	SnoozedBy           *User             `json:"snoozedBy"`
	StartedAt           *time.Time        `json:"startedAt"`
	CompletedAt         *time.Time        `json:"completedAt"`
	CanceledAt          *time.Time        `json:"canceledAt"`
	ArchivedAt          *time.Time        `json:"archivedAt"`
	TriagedAt           *time.Time        `json:"triagedAt"`
	CustomerTicketCount int               `json:"customerTicketCount"`
	PreviousIdentifiers []string          `json:"previousIdentifiers"`
	Trashed             *bool             `json:"trashed"`
	// SLA fields
	SLAStartedAt    *time.Time `json:"slaStartedAt"`
	SLAMediumRiskAt *time.Time `json:"slaMediumRiskAt"`
	SLAHighRiskAt   *time.Time `json:"slaHighRiskAt"`
	SLABreachesAt   *time.Time `json:"slaBreachesAt"`
	SLAType         *string    `json:"slaType"`
	// Additional fields
	Number                int              `json:"number"`
	BoardOrder            float64          `json:"boardOrder"`
	SubIssueSortOrder     float64          `json:"subIssueSortOrder"`
	PriorityLabel         string           `json:"priorityLabel"`
	IntegrationSourceType *string          `json:"integrationSourceType"`
	Creator               *User            `json:"creator"`
	Subscribers           *Users           `json:"subscribers"`
	Relations             *IssueRelations  `json:"relations"`
	History               *IssueHistory    `json:"history"`
	Reactions             []Reaction       `json:"reactions"`
	SlackIssueComments    []SlackComment   `json:"slackIssueComments"`
	ExternalUserCreator   *ExternalUser    `json:"externalUserCreator"`
	CustomerTickets       []CustomerTicket `json:"customerTickets"`
}

// State represents an issue state
type State struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Color       string  `json:"color"`
	Description *string `json:"description"`
	Position    float64 `json:"position"`
}

// Project represents a Linear project
type Project struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	State             string             `json:"state"`
	Progress          float64            `json:"progress"`
	StartDate         *string            `json:"startDate"`
	TargetDate        *string            `json:"targetDate"`
	Lead              *User              `json:"lead"`
	Teams             *Teams             `json:"teams"`
	URL               string             `json:"url"`
	Icon              *string            `json:"icon"`
	Color             string             `json:"color"`
	CreatedAt         time.Time          `json:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt"`
	CompletedAt       *time.Time         `json:"completedAt"`
	CanceledAt        *time.Time         `json:"canceledAt"`
	ArchivedAt        *time.Time         `json:"archivedAt"`
	Creator           *User              `json:"creator"`
	Members           *Users             `json:"members"`
	Issues            *Issues            `json:"issues"`
	ProjectMilestones *ProjectMilestones `json:"projectMilestones"`
	// Additional fields
	SlugId               string          `json:"slugId"`
	Content              string          `json:"content"`
	ConvertedFromIssue   *Issue          `json:"convertedFromIssue"`
	LastAppliedTemplate  *Template       `json:"lastAppliedTemplate"`
	ProjectUpdates       *ProjectUpdates `json:"projectUpdates"`
	Documents            *Documents      `json:"documents"`
	Health               string          `json:"health"`
	Scope                float64         `json:"scope"`
	SlackNewIssue        bool            `json:"slackNewIssue"`
	SlackIssueComments   bool            `json:"slackIssueComments"`
	SlackIssueStatuses   bool            `json:"slackIssueStatuses"`
	Priority             int             `json:"priority"`
	PriorityLabel        string          `json:"priorityLabel"`
	StartedAt            *time.Time      `json:"startedAt"`
	AutoArchivedAt       *time.Time      `json:"autoArchivedAt"`
	Trashed              bool            `json:"trashed"`
	HealthUpdatedAt      *time.Time      `json:"healthUpdatedAt"`
	StartDateResolution  string          `json:"startDateResolution"`
	TargetDateResolution string          `json:"targetDateResolution"`
	SortOrder            float64         `json:"sortOrder"`
	PrioritySortOrder    float64         `json:"prioritySortOrder"`
}

// Paginated collections
type Issues struct {
	Nodes    []Issue  `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

type Teams struct {
	Nodes    []Team   `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

type Projects struct {
	Nodes    []Project `json:"nodes"`
	PageInfo PageInfo  `json:"pageInfo"`
}

type Users struct {
	Nodes    []User   `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

type Labels struct {
	Nodes    []Label  `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

type Label struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	Description *string `json:"description"`
	Parent      *Label  `json:"parent"`
}

// Cycle represents a Linear cycle (sprint)
type Cycle struct {
	ID                          string     `json:"id"`
	Number                      int        `json:"number"`
	Name                        string     `json:"name"`
	Description                 *string    `json:"description"`
	StartsAt                    string     `json:"startsAt"`
	EndsAt                      string     `json:"endsAt"`
	Progress                    float64    `json:"progress"`
	CompletedAt                 *time.Time `json:"completedAt"`
	CreatedAt                   time.Time  `json:"createdAt"`
	UpdatedAt                   time.Time  `json:"updatedAt"`
	ArchivedAt                  *time.Time `json:"archivedAt"`
	AutoArchivedAt              *time.Time `json:"autoArchivedAt"`
	IsActive                    bool       `json:"isActive"`
	IsFuture                    bool       `json:"isFuture"`
	IsPast                      bool       `json:"isPast"`
	IsNext                      bool       `json:"isNext"`
	IsPrevious                  bool       `json:"isPrevious"`
	ScopeHistory                []float64  `json:"scopeHistory"`
	CompletedScopeHistory       []float64  `json:"completedScopeHistory"`
	InProgressScopeHistory      []float64  `json:"inProgressScopeHistory"`
	IssueCountHistory           []float64  `json:"issueCountHistory"`
	CompletedIssueCountHistory  []float64  `json:"completedIssueCountHistory"`
	Team                        *Team      `json:"team"`
	Issues                      *Issues    `json:"issues"`
}

type Cycles struct {
	Nodes    []Cycle  `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

// Release represents a Linear release (for document linking)
type Release struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Version *string `json:"version"`
}

// Attachment represents a file attachment or link
type Attachment struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Subtitle  *string                `json:"subtitle"`
	URL       string                 `json:"url"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"createdAt"`
	Creator   *User                  `json:"creator"`

	// Use a map to capture any extra fields Linear might return
	Extra map[string]interface{} `json:"-"`
}

// UnmarshalJSON implements custom unmarshaling to handle unexpected fields from Linear API
func (a *Attachment) UnmarshalJSON(data []byte) error {
	// Create an alias to avoid infinite recursion
	type Alias Attachment
	aux := &struct {
		*Alias
		// Capture extra fields that might come from Linear
		Source     interface{} `json:"source,omitempty"`
		SourceType interface{} `json:"sourceType,omitempty"`
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Store unexpected fields in Extra map if needed
	if aux.Source != nil || aux.SourceType != nil {
		a.Extra = make(map[string]interface{})
		if aux.Source != nil {
			a.Extra["source"] = aux.Source
		}
		if aux.SourceType != nil {
			a.Extra["sourceType"] = aux.SourceType
		}
	}

	return nil
}

// Attachments represents a paginated list of attachments
type Attachments struct {
	Nodes []Attachment `json:"nodes"`
}

// Initiative represents a Linear initiative
type Initiative struct {
	ID                   string       `json:"id"`
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Status               string       `json:"status"`
	SlugId               string       `json:"slugId"`
	Color                string       `json:"color"`
	Icon                 *string      `json:"icon"`
	Content              string       `json:"content"`
	TargetDate           *string      `json:"targetDate"`
	TargetDateResolution string       `json:"targetDateResolution"`
	Owner                *User        `json:"owner"`
	Creator              *User        `json:"creator"`
	URL                  string       `json:"url"`
	Health               string       `json:"health"`
	CreatedAt            time.Time    `json:"createdAt"`
	UpdatedAt            time.Time    `json:"updatedAt"`
	ArchivedAt           *time.Time   `json:"archivedAt"`
	CompletedAt          *time.Time   `json:"completedAt"`
	Projects             *Projects    `json:"projects"`
	ParentInitiative     *Initiative  `json:"parentInitiative"`
	SubInitiatives       *Initiatives `json:"subInitiatives"`
}

// Initiatives represents a paginated list of initiatives
type Initiatives struct {
	Nodes    []Initiative `json:"nodes"`
	PageInfo PageInfo     `json:"pageInfo"`
}

type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

// Additional types for expanded fields
type IssueRelations struct {
	Nodes []IssueRelation `json:"nodes"`
}

type IssueRelation struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Issue        *Issue `json:"issue"`
	RelatedIssue *Issue `json:"relatedIssue"`
}

type IssueHistory struct {
	Nodes []IssueHistoryEntry `json:"nodes"`
}

type IssueHistoryEntry struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Changes         string    `json:"changes"`
	Actor           *User     `json:"actor"`
	FromAssignee    *User     `json:"fromAssignee"`
	ToAssignee      *User     `json:"toAssignee"`
	FromState       *State    `json:"fromState"`
	ToState         *State    `json:"toState"`
	FromPriority    *int      `json:"fromPriority"`
	ToPriority      *int      `json:"toPriority"`
	FromTitle       *string   `json:"fromTitle"`
	ToTitle         *string   `json:"toTitle"`
	FromCycle       *Cycle    `json:"fromCycle"`
	ToCycle         *Cycle    `json:"toCycle"`
	FromProject     *Project  `json:"fromProject"`
	ToProject       *Project  `json:"toProject"`
	AddedLabelIds   []string  `json:"addedLabelIds"`
	RemovedLabelIds []string  `json:"removedLabelIds"`
}

type Reaction struct {
	ID        string    `json:"id"`
	Emoji     string    `json:"emoji"`
	User      *User     `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}

type SlackComment struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type ExternalUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CustomerTicket struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"createdAt"`
	ExternalId string    `json:"externalId"`
}

type Template struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Milestone is the legacy workspace-level milestone (deprecated by Linear)
type Milestone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TargetDate  *string   `json:"targetDate"`
	Projects    *Projects `json:"projects"`
}

// ProjectMilestone represents a milestone within a Linear project
type ProjectMilestone struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	TargetDate  *string    `json:"targetDate"`
	Status      string     `json:"status"`
	Progress    float64    `json:"progress"`
	SortOrder   float64    `json:"sortOrder"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
	Project     *Project   `json:"project"`
	Issues      *Issues    `json:"issues"`
}

// ProjectMilestones represents a paginated list of project milestones
type ProjectMilestones struct {
	Nodes    []ProjectMilestone `json:"nodes"`
	PageInfo PageInfo           `json:"pageInfo"`
}

type Roadmaps struct {
	Nodes []Roadmap `json:"nodes"`
}

type Roadmap struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Creator     *User     `json:"creator"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProjectUpdates struct {
	Nodes []ProjectUpdate `json:"nodes"`
}

type ProjectUpdate struct {
	ID         string     `json:"id"`
	Body       string     `json:"body"`
	User       *User      `json:"user"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	EditedAt   *time.Time `json:"editedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
	Health     string     `json:"health"`
	URL        string     `json:"url"`
	Project    *Project   `json:"project"`
}

type Documents struct {
	Nodes    []Document `json:"nodes"`
	PageInfo PageInfo   `json:"pageInfo"`
}

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Icon      *string   `json:"icon"`
	Color     string    `json:"color"`
	SlugId    string    `json:"slugId"`
	URL       string    `json:"url"`
	Creator   *User     `json:"creator"`
	UpdatedBy *User     `json:"updatedBy"`
	Project   *Project  `json:"project"`
	Team      *Team     `json:"team"`
	Issue      *Issue      `json:"issue"`
	Initiative *Initiative `json:"initiative"`
	Cycle      *Cycle      `json:"cycle"`
	Release    *Release    `json:"release"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

type ProjectLinks struct {
	Nodes []ProjectLink `json:"nodes"`
}

type ProjectLink struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Label     string    `json:"label"`
	Creator   *User     `json:"creator"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CustomView represents a Linear custom view (saved filter)
type CustomView struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Description       *string                `json:"description"`
	Icon              *string                `json:"icon"`
	Color             *string                `json:"color"`
	Shared            bool                   `json:"shared"`
	SlugId            string                 `json:"slugId"`
	ModelName         string                 `json:"modelName"`
	FilterData        map[string]interface{} `json:"filterData"`
	ProjectFilterData map[string]interface{} `json:"projectFilterData"`
	Creator           *User                  `json:"creator"`
	Owner             *User                  `json:"owner"`
	Team              *Team                  `json:"team"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

// CustomViews represents a paginated list of custom views
type CustomViews struct {
	Nodes    []CustomView `json:"nodes"`
	PageInfo PageInfo     `json:"pageInfo"`
}

// GetViewer returns the current authenticated user
func (c *Client) GetViewer(ctx context.Context) (*User, error) {
	query := `
		query Me {
			viewer {
				id
				name
				email
				avatarUrl
				isMe
				active
				admin
			}
		}
	`

	var response struct {
		Viewer User `json:"viewer"`
	}

	err := c.Execute(ctx, query, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Viewer, nil
}

// GetIssues returns a list of issues with optional filtering
func (c *Client) GetIssues(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*Issues, error) {
	query := `
		query Issues($filter: IssueFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			issues(filter: $filter, first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					identifier
					title
					description
					priority
					priorityLabel
					estimate
					createdAt
					updatedAt
					startedAt
					completedAt
					canceledAt
					archivedAt
					triagedAt
					snoozedUntilAt
					dueDate
					branchName
					customerTicketCount
					url
					slaStartedAt
					slaMediumRiskAt
					slaHighRiskAt
					slaBreachesAt
					slaType
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
					}
					creator {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					labels {
						nodes {
							id
							name
							color
						}
					}
					cycle {
						id
						number
						name
					}
					project {
						id
						name
						state
					}
					projectMilestone {
						id
						name
						status
					}
					parent {
						id
						identifier
						title
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
		Issues Issues `json:"issues"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issues, nil
}

// IssueSearch returns issues that match a full-text query
func (c *Client) IssueSearch(ctx context.Context, term string, filter map[string]interface{}, first int, after string, orderBy string, includeArchived bool) (*Issues, error) {
	query := `
		query IssueSearch($term: String!, $filter: IssueFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			searchIssues(term: $term, filter: $filter, first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {
					id
					identifier
					title
					description
					priority
					priorityLabel
					estimate
					createdAt
					updatedAt
					startedAt
					completedAt
					canceledAt
					archivedAt
					triagedAt
					snoozedUntilAt
					dueDate
					branchName
					customerTicketCount
					url
					slaStartedAt
					slaMediumRiskAt
					slaHighRiskAt
					slaBreachesAt
					slaType
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
					}
					creator {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					labels {
						nodes {
							id
							name
							color
						}
					}
					cycle {
						id
						number
						name
					}
					project {
						id
						name
						state
					}
					projectMilestone {
						id
						name
						status
					}
					parent {
						id
						identifier
						title
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
		SearchIssues struct {
			Nodes    []Issue  `json:"nodes"`
			PageInfo PageInfo `json:"pageInfo"`
		} `json:"searchIssues"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &Issues{
		Nodes:    response.SearchIssues.Nodes,
		PageInfo: response.SearchIssues.PageInfo,
	}, nil
}

// GetIssue returns a single issue by ID
func (c *Client) GetIssue(ctx context.Context, id string) (*Issue, error) {
	query := `
		query Issue($id: String!) {
			issue(id: $id) {
				id
				identifier
				number
				title
				description
				priority
				priorityLabel
				estimate
				boardOrder
				subIssueSortOrder
				createdAt
				updatedAt
				startedAt
				dueDate
				url
				branchName
				snoozedUntilAt
				snoozedBy {
					id
					name
					email
				}
				completedAt
				canceledAt
				archivedAt
				triagedAt
				trashed
				customerTicketCount
				previousIdentifiers
				integrationSourceType
				slaStartedAt
				slaMediumRiskAt
				slaHighRiskAt
				slaBreachesAt
				slaType
				state {
					id
					name
					type
					color
					description
					position
				}
				assignee {
					id
					name
					email
					avatarUrl
					displayName
					active
					admin
					createdAt
				}
				creator {
					id
					name
					email
					avatarUrl
					displayName
					active
				}
				team {
					id
					key
					name
					description
					icon
					color
					cyclesEnabled
					cycleStartDay
					cycleDuration
					upcomingCycleCount
				}
				labels {
					nodes {
						id
						name
						color
						description
						parent {
							id
							name
						}
					}
				}
				parent {
					id
					identifier
					title
					state {
						name
						type
					}
				}
				children {
					nodes {
						id
						identifier
						title
						priority
						createdAt
						state {
							name
							type
							color
						}
						assignee {
							name
							email
						}
					}
				}
				cycle {
					id
					number
					name
					description
					startsAt
					endsAt
					progress
					completedAt
					scopeHistory
				}
				project {
					id
					name
					description
					state
					progress
					startDate
					targetDate
					health
					lead {
						name
						email
					}
				}
				projectMilestone {
					id
					name
					targetDate
					status
				}
				attachments(first: 20) {
					nodes {
						id
						title
						subtitle
						url
						metadata
						createdAt
						creator {
							name
							email
						}
					}
				}
				documents(first: 20) {
					nodes {
						id
						title
						icon
						color
						slugId
						url
						createdAt
						updatedAt
						creator {
							name
							email
						}
					}
				}
				comments(first: 10) {
					nodes {
						id
						body
						createdAt
						updatedAt
						editedAt
						user {
							name
							email
							avatarUrl
						}
						parent {
							id
						}
						children {
							nodes {
								id
								body
								user {
									name
								}
							}
						}
					}
				}
				subscribers {
					nodes {
						id
						name
						email
						avatarUrl
					}
				}
				relations {
					nodes {
						id
						type
						relatedIssue {
							id
							identifier
							title
							state {
								name
								type
							}
						}
					}
				}
				history(first: 10) {
					nodes {
						id
						createdAt
						updatedAt
						actor {
							name
							email
						}
						fromAssignee {
							name
						}
						toAssignee {
							name
						}
						fromState {
							name
						}
						toState {
							name
						}
						fromPriority
						toPriority
						fromTitle
						toTitle
						fromCycle {
							name
						}
						toCycle {
							name
						}
						fromProject {
							name
						}
						toProject {
							name
						}
						addedLabelIds
						removedLabelIds
					}
				}
				reactions {
					id
					emoji
					user {
						name
						email
					}
					createdAt
				}
				externalUserCreator {
					id
					name
					email
					avatarUrl
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Issue Issue `json:"issue"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issue, nil
}

// GetTeams returns a list of teams
func (c *Client) GetTeams(ctx context.Context, first int, after string, orderBy string) (*Teams, error) {
	query := `
		query Teams($first: Int, $after: String, $orderBy: PaginationOrderBy) {
			teams(first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					key
					name
					displayName
					description
					icon
					color
					private
					issueCount
					timezone
					cyclesEnabled
					cycleStartDay
					cycleDuration
					cycleCooldownTime
					cycleIssueAutoAssignStarted
					cycleIssueAutoAssignCompleted
					cycleLockToActive
					upcomingCycleCount
					triageEnabled
					requirePriorityToLeaveTriage
					inheritIssueEstimation
					issueEstimationType
					issueEstimationAllowZero
					issueEstimationExtended
					defaultIssueEstimate
					setIssueSortOrderOnStateChange
					autoClosePeriod
					autoCloseStateId
					autoArchivePeriod
					autoCloseParentIssues
					autoCloseChildIssues
					inheritWorkflowStatuses
					groupIssueHistory
					joinByDefault
					allMembersCanJoin
					scimManaged
					scimGroupName
					aiThreadSummariesEnabled
					aiDiscussionSummariesEnabled
					createdAt
					updatedAt
					archivedAt
					retiredAt
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
		Teams Teams `json:"teams"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Teams, nil
}

// GetProjects returns a list of projects
func (c *Client) GetProjects(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*Projects, error) {
	query := `
		query Projects($filter: ProjectFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			projects(filter: $filter, first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					name
					description
					state
					progress
					startDate
					targetDate
					url
					icon
					color
					createdAt
					updatedAt
					completedAt
					canceledAt
					health
					priority
					priorityLabel
					scope
					lead {
						id
						name
						email
					}
					teams {
						nodes {
							id
							key
							name
						}
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
		Projects Projects `json:"projects"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Projects, nil
}

// GetProject returns a single project by ID
func (c *Client) GetProject(ctx context.Context, id string) (*Project, error) {
	query := `
		query Project($id: String!) {
			project(id: $id) {
				id
				slugId
				name
				description
				content
				state
				progress
				health
				healthUpdatedAt
				scope
				startDate
				startDateResolution
				targetDate
				targetDateResolution
				url
				icon
				color
				createdAt
				updatedAt
				completedAt
				canceledAt
				archivedAt
				startedAt
				autoArchivedAt
				trashed
				priority
				priorityLabel
				sortOrder
				prioritySortOrder
				slackNewIssue
				slackIssueComments
				slackIssueStatuses
				lead {
					id
					name
					email
					avatarUrl
					displayName
					active
				}
				creator {
					id
					name
					email
					avatarUrl
					active
				}
				convertedFromIssue {
					id
					identifier
					title
				}
				lastAppliedTemplate {
					id
					name
					description
				}
				teams {
					nodes {
						id
						key
						name
						description
						icon
						color
						cyclesEnabled
					}
				}
				members {
					nodes {
						id
						name
						email
						avatarUrl
						displayName
						active
						admin
					}
				}
				projectMilestones(first: 50) {
				nodes {
					id
					name
					description
					targetDate
					status
					progress
					sortOrder
				}
			}
			issues(first: 50, orderBy: updatedAt) {
					nodes {
						id
						identifier
						number
						title
						description
						priority
						estimate
						createdAt
						updatedAt
						completedAt
						state {
							name
							type
							color
						}
						assignee {
							name
							email
						}
						labels {
							nodes {
								name
								color
							}
						}
					}
				}
				projectUpdates(first: 10) {
					nodes {
						id
						body
						health
						createdAt
						updatedAt
						editedAt
						user {
							name
							email
							avatarUrl
						}
					}
				}
				documents(first: 20) {
					nodes {
						id
						title
						content
						icon
						color
						createdAt
						updatedAt
						creator {
							name
							email
						}
						updatedBy {
							name
							email
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Project Project `json:"project"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Project, nil
}

// UpdateProject updates a project's fields
func (c *Client) UpdateProject(ctx context.Context, id string, input map[string]interface{}) (*Project, error) {
	query := `
		mutation ProjectUpdate($id: String!, $input: ProjectUpdateInput!) {
			projectUpdate(id: $id, input: $input) {
				project {
					id
					name
					description
					state
					progress
					url
					icon
					color
					priority
					priorityLabel
					health
					scope
					startDate
					targetDate
					slackNewIssue
					slackIssueComments
					slackIssueStatuses
					createdAt
					updatedAt
					lead {
						id
						name
						email
					}
					members {
						nodes {
							id
							name
							email
						}
					}
					teams {
						nodes {
							id
							key
							name
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		ProjectUpdate struct {
			Project Project `json:"project"`
		} `json:"projectUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ProjectUpdate.Project, nil
}

// UpdateIssue updates an issue's fields
func (c *Client) UpdateIssue(ctx context.Context, id string, input map[string]interface{}) (*Issue, error) {
	query := `
		mutation UpdateIssue($id: String!, $input: IssueUpdateInput!) {
			issueUpdate(id: $id, input: $input) {
				issue {
					id
					identifier
					title
					description
					priority
					estimate
					createdAt
					updatedAt
					dueDate
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					labels {
						nodes {
							id
							name
							color
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		IssueUpdate struct {
			Issue Issue `json:"issue"`
		} `json:"issueUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.IssueUpdate.Issue, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(ctx context.Context, input map[string]interface{}) (*Issue, error) {
	query := `
		mutation CreateIssue($input: IssueCreateInput!) {
			issueCreate(input: $input) {
				issue {
					id
					identifier
					title
					description
					priority
					estimate
					createdAt
					updatedAt
					dueDate
					state {
						id
						name
						type
						color
					}
					assignee {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					labels {
						nodes {
							id
							name
							color
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		IssueCreate struct {
			Issue Issue `json:"issue"`
		} `json:"issueCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.IssueCreate.Issue, nil
}

// GetTeam returns a single team by key
func (c *Client) GetTeam(ctx context.Context, key string) (*Team, error) {
	query := `
		query Team($key: String!) {
			team(id: $key) {
				id
				key
				name
				displayName
				description
				icon
				color
				private
				issueCount
				timezone
				cyclesEnabled
				cycleStartDay
				cycleDuration
				cycleCooldownTime
				cycleIssueAutoAssignStarted
				cycleIssueAutoAssignCompleted
				cycleLockToActive
				upcomingCycleCount
				cycleCalenderUrl
				triageEnabled
				requirePriorityToLeaveTriage
				inheritIssueEstimation
				issueEstimationType
				issueEstimationAllowZero
				issueEstimationExtended
				defaultIssueEstimate
				setIssueSortOrderOnStateChange
				autoClosePeriod
				autoCloseStateId
				autoArchivePeriod
				autoCloseParentIssues
				autoCloseChildIssues
				inheritWorkflowStatuses
				groupIssueHistory
				joinByDefault
				allMembersCanJoin
				scimManaged
				scimGroupName
				aiThreadSummariesEnabled
				aiDiscussionSummariesEnabled
				createdAt
				updatedAt
				archivedAt
				retiredAt
				parent {
					id
					key
					name
				}
				defaultIssueState {
					id
					name
					type
					color
				}
				triageIssueState {
					id
					name
					type
					color
				}
				markedAsDuplicateWorkflowState {
					id
					name
					type
					color
				}
				activeCycle {
					id
					number
					name
					startsAt
					endsAt
				}
				defaultTemplateForMembers {
					id
					name
				}
				defaultTemplateForNonMembers {
					id
					name
				}
				defaultProjectTemplate {
					id
					name
				}
			}
		}
	`

	variables := map[string]interface{}{
		"key": key,
	}

	var response struct {
		Team Team `json:"team"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Team, nil
}

// Comment represents a Linear comment
type Comment struct {
	ID        string     `json:"id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	EditedAt  *time.Time `json:"editedAt"`
	User      *User      `json:"user"`
	Parent    *Comment   `json:"parent"`
	Children  *Comments  `json:"children"`
}

// Comments represents a paginated list of comments
type Comments struct {
	Nodes    []Comment `json:"nodes"`
	PageInfo PageInfo  `json:"pageInfo"`
}

// Notification represents a Linear notification (inbox item)
type Notification struct {
	ID             string     `json:"id"`
	Type           string     `json:"type"`
	CreatedAt      time.Time  `json:"createdAt"`
	ReadAt         *time.Time `json:"readAt"`
	SnoozedUntilAt *time.Time `json:"snoozedUntilAt"`
	ArchivedAt     *time.Time `json:"archivedAt"`
	Actor          *User      `json:"actor"`
	// IssueNotification fields
	Issue        *Issue `json:"issue"`
	CommentID    string `json:"commentId"`
	ReactionEmoji string `json:"reactionEmoji"`
}

// Notifications represents a paginated list of notifications
type Notifications struct {
	Nodes    []Notification `json:"nodes"`
	PageInfo PageInfo       `json:"pageInfo"`
}

// WorkflowState represents a Linear workflow state
type WorkflowState struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Color       string  `json:"color"`
	Description string  `json:"description"`
	Position    float64 `json:"position"`
}

// GetTeamStates returns workflow states for a team
func (c *Client) GetTeamStates(ctx context.Context, teamKey string) ([]WorkflowState, error) {
	query := `
		query TeamStates($key: String!) {
			team(id: $key) {
				states {
					nodes {
						id
						name
						type
						color
						description
						position
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"key": teamKey,
	}

	var response struct {
		Team struct {
			States struct {
				Nodes []WorkflowState `json:"nodes"`
			} `json:"states"`
		} `json:"team"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return response.Team.States.Nodes, nil
}

// GetTeamMembers returns members of a specific team
func (c *Client) GetTeamMembers(ctx context.Context, teamKey string) (*Users, error) {
	query := `
		query TeamMembers($key: String!) {
			team(id: $key) {
				members {
					nodes {
						id
						name
						email
						avatarUrl
						isMe
						active
						admin
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"key": teamKey,
	}

	var response struct {
		Team struct {
			Members Users `json:"members"`
		} `json:"team"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Team.Members, nil
}

// GetUsers returns a list of all users
func (c *Client) GetUsers(ctx context.Context, first int, after string, orderBy string) (*Users, error) {
	query := `
		query Users($first: Int, $after: String, $orderBy: PaginationOrderBy) {
			users(first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					name
					email
					avatarUrl
					isMe
					active
					admin
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
		Users Users `json:"users"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Users, nil
}

// GetUser returns a specific user by email
func (c *Client) GetUser(ctx context.Context, email string) (*User, error) {
	query := `
		query User($email: String!) {
			user(email: $email) {
				id
				name
				email
				avatarUrl
				isMe
				active
				admin
			}
		}
	`

	variables := map[string]interface{}{
		"email": email,
	}

	var response struct {
		User User `json:"user"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.User, nil
}

// GetIssueComments returns comments for a specific issue
func (c *Client) GetIssueComments(ctx context.Context, issueID string, first int, after string, orderBy string) (*Comments, error) {
	query := `
		query IssueComments($id: String!, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			issue(id: $id) {
				comments(first: $first, after: $after, orderBy: $orderBy) {
					nodes {
						id
						body
						createdAt
						updatedAt
						user {
							id
							name
							email
						}
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    issueID,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Issue struct {
			Comments Comments `json:"comments"`
		} `json:"issue"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issue.Comments, nil
}

// CreateComment creates a new comment on an issue
func (c *Client) CreateComment(ctx context.Context, issueID string, body string) (*Comment, error) {
	query := `
		mutation CreateComment($input: CommentCreateInput!) {
			commentCreate(input: $input) {
				comment {
					id
					body
					createdAt
					updatedAt
					user {
						id
						name
						email
					}
				}
			}
		}
	`

	input := map[string]interface{}{
		"issueId": issueID,
		"body":    body,
	}

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CommentCreate struct {
			Comment Comment `json:"comment"`
		} `json:"commentCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CommentCreate.Comment, nil
}

// GetDocuments returns a list of documents with optional filtering
func (c *Client) GetDocuments(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*Documents, error) {
	query := `
		query Documents($filter: DocumentFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy) {
			documents(filter: $filter, first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					title
					icon
					color
					slugId
					url
					createdAt
					updatedAt
					creator {
						id
						name
						email
					}
					updatedBy {
						id
						name
						email
					}
					project {
						id
						name
					}
					team {
						id
						key
						name
					}
					issue {
						id
						identifier
						title
					}
					initiative {
						id
						name
					}
					cycle {
						id
						number
						name
					}
					release {
						id
						name
						version
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
		Documents Documents `json:"documents"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Documents, nil
}

// GetDocument returns a single document by ID
func (c *Client) GetDocument(ctx context.Context, id string) (*Document, error) {
	query := `
		query Document($id: String!) {
			document(id: $id) {
				id
				title
				content
				icon
				color
				slugId
				url
				createdAt
				updatedAt
				creator {
					id
					name
					email
					avatarUrl
				}
				updatedBy {
					id
					name
					email
				}
				project {
					id
					name
					state
					progress
					url
				}
				team {
					id
					key
					name
				}
				issue {
					id
					identifier
					title
				}
				initiative {
					id
					name
				}
				cycle {
					id
					number
					name
				}
				release {
					id
					name
					version
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Document Document `json:"document"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Document, nil
}

// SearchDocuments returns documents matching a search query
func (c *Client) SearchDocuments(ctx context.Context, term string, first int, after string, orderBy string, teamID string, includeComments bool) (*Documents, error) {
	query := `
		query SearchDocuments($term: String!, $first: Int, $after: String, $orderBy: PaginationOrderBy, $teamId: String, $includeComments: Boolean) {
			searchDocuments(term: $term, first: $first, after: $after, orderBy: $orderBy, teamId: $teamId, includeComments: $includeComments) {
				nodes {
					id
					title
					icon
					color
					slugId
					url
					createdAt
					updatedAt
					creator {
						id
						name
						email
					}
					updatedBy {
						id
						name
						email
					}
					project {
						id
						name
					}
					team {
						id
						key
						name
					}
					issue {
						id
						identifier
						title
					}
					initiative {
						id
						name
					}
					cycle {
						id
						number
						name
					}
					release {
						id
						name
						version
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
		"includeComments": includeComments,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}
	if teamID != "" {
		variables["teamId"] = teamID
	}

	var response struct {
		SearchDocuments struct {
			Nodes    []Document `json:"nodes"`
			PageInfo PageInfo   `json:"pageInfo"`
		} `json:"searchDocuments"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &Documents{
		Nodes:    response.SearchDocuments.Nodes,
		PageInfo: response.SearchDocuments.PageInfo,
	}, nil
}

// CreateDocument creates a new document
func (c *Client) CreateDocument(ctx context.Context, input map[string]interface{}) (*Document, error) {
	query := `
		mutation CreateDocument($input: DocumentCreateInput!) {
			documentCreate(input: $input) {
				document {
					id
					title
					content
					icon
					color
					slugId
					url
					createdAt
					updatedAt
					creator {
						id
						name
						email
					}
					project {
						id
						name
					}
					team {
						id
						key
						name
					}
					issue {
						id
						identifier
						title
					}
					initiative {
						id
						name
					}
					cycle {
						id
						number
						name
					}
					release {
						id
						name
						version
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		DocumentCreate struct {
			Document Document `json:"document"`
		} `json:"documentCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.DocumentCreate.Document, nil
}

// UpdateDocument updates a document
func (c *Client) UpdateDocument(ctx context.Context, id string, input map[string]interface{}) (*Document, error) {
	query := `
		mutation UpdateDocument($id: String!, $input: DocumentUpdateInput!) {
			documentUpdate(id: $id, input: $input) {
				document {
					id
					title
					content
					icon
					color
					slugId
					url
					createdAt
					updatedAt
					creator {
						id
						name
						email
					}
					project {
						id
						name
					}
					team {
						id
						key
						name
					}
					issue {
						id
						identifier
						title
					}
					initiative {
						id
						name
					}
					cycle {
						id
						number
						name
					}
					release {
						id
						name
						version
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		DocumentUpdate struct {
			Document Document `json:"document"`
		} `json:"documentUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.DocumentUpdate.Document, nil
}

// DeleteDocument deletes a document
func (c *Client) DeleteDocument(ctx context.Context, id string) error {
	query := `
		mutation DeleteDocument($id: String!) {
			documentDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		DocumentDelete struct {
			Success bool `json:"success"`
		} `json:"documentDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.DocumentDelete.Success {
		return fmt.Errorf("failed to delete document")
	}

	return nil
}

// GetProjectMilestones returns milestones for a specific project
func (c *Client) GetProjectMilestones(ctx context.Context, projectID string, first int, after string) (*ProjectMilestones, error) {
	query := `
		query ProjectMilestones($filter: ProjectMilestoneFilter, $first: Int, $after: String) {
			projectMilestones(filter: $filter, first: $first, after: $after) {
				nodes {
					id
					name
					description
					targetDate
					status
					progress
					sortOrder
					createdAt
					updatedAt
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
		"filter": map[string]interface{}{
			"project": map[string]interface{}{"id": map[string]interface{}{"eq": projectID}},
		},
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		ProjectMilestones ProjectMilestones `json:"projectMilestones"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ProjectMilestones, nil
}

// GetProjectMilestone returns a single project milestone by ID
func (c *Client) GetProjectMilestone(ctx context.Context, id string) (*ProjectMilestone, error) {
	query := `
		query ProjectMilestone($id: String!) {
			projectMilestone(id: $id) {
				id
				name
				description
				targetDate
				status
				progress
				sortOrder
				createdAt
				updatedAt
				archivedAt
				project {
					id
					name
					state
					progress
				}
				issues(first: 50) {
					nodes {
						id
						identifier
						title
						priority
						createdAt
						state {
							name
							type
							color
						}
						assignee {
							name
							email
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		ProjectMilestone ProjectMilestone `json:"projectMilestone"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ProjectMilestone, nil
}

// CreateProjectMilestone creates a new milestone on a project
func (c *Client) CreateProjectMilestone(ctx context.Context, input map[string]interface{}) (*ProjectMilestone, error) {
	query := `
		mutation ProjectMilestoneCreate($input: ProjectMilestoneCreateInput!) {
			projectMilestoneCreate(input: $input) {
				projectMilestone {
					id
					name
					description
					targetDate
					status
					progress
					sortOrder
					createdAt
					updatedAt
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		ProjectMilestoneCreate struct {
			ProjectMilestone ProjectMilestone `json:"projectMilestone"`
		} `json:"projectMilestoneCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ProjectMilestoneCreate.ProjectMilestone, nil
}

// UpdateProjectMilestone updates an existing project milestone
func (c *Client) UpdateProjectMilestone(ctx context.Context, id string, input map[string]interface{}) (*ProjectMilestone, error) {
	query := `
		mutation ProjectMilestoneUpdate($id: String!, $input: ProjectMilestoneUpdateInput!) {
			projectMilestoneUpdate(id: $id, input: $input) {
				projectMilestone {
					id
					name
					description
					targetDate
					status
					progress
					sortOrder
					createdAt
					updatedAt
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		ProjectMilestoneUpdate struct {
			ProjectMilestone ProjectMilestone `json:"projectMilestone"`
		} `json:"projectMilestoneUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ProjectMilestoneUpdate.ProjectMilestone, nil
}

// GetProjectUpdates returns status updates for a project
func (c *Client) GetProjectUpdates(ctx context.Context, projectID string, first int, after string) (*ProjectUpdates, error) {
	query := `
		query ProjectUpdates($id: String!, $first: Int, $after: String) {
			project(id: $id) {
				projectUpdates(first: $first, after: $after) {
					nodes {
						id
						body
						health
						url
						createdAt
						updatedAt
						editedAt
						archivedAt
						user {
							id
							name
							email
							avatarUrl
						}
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    projectID,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		Project struct {
			ProjectUpdates struct {
				Nodes    []ProjectUpdate `json:"nodes"`
				PageInfo PageInfo        `json:"pageInfo"`
			} `json:"projectUpdates"`
		} `json:"project"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &ProjectUpdates{
		Nodes: response.Project.ProjectUpdates.Nodes,
	}, nil
}

// GetProjectUpdate returns a single project status update by ID
func (c *Client) GetProjectUpdate(ctx context.Context, id string) (*ProjectUpdate, error) {
	query := `
		query ProjectUpdate($id: String!) {
			projectUpdate(id: $id) {
				id
				body
				health
				url
				createdAt
				updatedAt
				editedAt
				archivedAt
				user {
					id
					name
					email
					avatarUrl
				}
				project {
					id
					name
					state
					progress
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		ProjectUpdate ProjectUpdate `json:"projectUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ProjectUpdate, nil
}

// CreateProjectUpdate creates a new status update on a project
func (c *Client) CreateProjectUpdate(ctx context.Context, input map[string]interface{}) (*ProjectUpdate, error) {
	query := `
		mutation ProjectUpdateCreate($input: ProjectUpdateCreateInput!) {
			projectUpdateCreate(input: $input) {
				projectUpdate {
					id
					body
					health
					url
					createdAt
					updatedAt
					user {
						id
						name
						email
					}
					project {
						id
						name
					}
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		ProjectUpdateCreate struct {
			ProjectUpdate ProjectUpdate `json:"projectUpdate"`
			Success       bool          `json:"success"`
		} `json:"projectUpdateCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.ProjectUpdateCreate.Success {
		return nil, fmt.Errorf("failed to create project update")
	}

	return &response.ProjectUpdateCreate.ProjectUpdate, nil
}

// UpdateProjectUpdate updates an existing project status update
func (c *Client) UpdateProjectUpdate(ctx context.Context, id string, input map[string]interface{}) (*ProjectUpdate, error) {
	query := `
		mutation ProjectUpdateUpdate($id: String!, $input: ProjectUpdateUpdateInput!) {
			projectUpdateUpdate(id: $id, input: $input) {
				projectUpdate {
					id
					body
					health
					url
					createdAt
					updatedAt
					editedAt
					user {
						id
						name
						email
					}
					project {
						id
						name
					}
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		ProjectUpdateUpdate struct {
			ProjectUpdate ProjectUpdate `json:"projectUpdate"`
			Success       bool          `json:"success"`
		} `json:"projectUpdateUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.ProjectUpdateUpdate.Success {
		return nil, fmt.Errorf("failed to update project update")
	}

	return &response.ProjectUpdateUpdate.ProjectUpdate, nil
}

// ArchiveProjectUpdate archives a project status update
func (c *Client) ArchiveProjectUpdate(ctx context.Context, id string) error {
	query := `
		mutation ProjectUpdateArchive($id: String!) {
			projectUpdateArchive(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		ProjectUpdateArchive struct {
			Success bool `json:"success"`
		} `json:"projectUpdateArchive"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.ProjectUpdateArchive.Success {
		return fmt.Errorf("failed to archive project update")
	}

	return nil
}

// CreateIssueRelation creates a relation between two issues
func (c *Client) CreateIssueRelation(ctx context.Context, issueID, relatedIssueID, relationType string) (*IssueRelation, error) {
	query := `
		mutation IssueRelationCreate($input: IssueRelationCreateInput!) {
			issueRelationCreate(input: $input) {
				issueRelation {
					id
					type
					issue {
						id
						identifier
						title
					}
					relatedIssue {
						id
						identifier
						title
					}
				}
				success
			}
		}
	`

	input := map[string]interface{}{
		"issueId":        issueID,
		"relatedIssueId": relatedIssueID,
		"type":           relationType,
	}

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		IssueRelationCreate struct {
			IssueRelation IssueRelation `json:"issueRelation"`
			Success       bool          `json:"success"`
		} `json:"issueRelationCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.IssueRelationCreate.Success {
		return nil, fmt.Errorf("failed to create issue relation")
	}

	return &response.IssueRelationCreate.IssueRelation, nil
}

// UpdateIssueRelation updates an existing issue relation
func (c *Client) UpdateIssueRelation(ctx context.Context, id string, input map[string]interface{}) (*IssueRelation, error) {
	query := `
		mutation IssueRelationUpdate($id: String!, $input: IssueRelationUpdateInput!) {
			issueRelationUpdate(id: $id, input: $input) {
				issueRelation {
					id
					type
					issue {
						id
						identifier
						title
					}
					relatedIssue {
						id
						identifier
						title
					}
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		IssueRelationUpdate struct {
			IssueRelation IssueRelation `json:"issueRelation"`
			Success       bool          `json:"success"`
		} `json:"issueRelationUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.IssueRelationUpdate.Success {
		return nil, fmt.Errorf("failed to update issue relation")
	}

	return &response.IssueRelationUpdate.IssueRelation, nil
}

// DeleteIssueRelation deletes an issue relation
func (c *Client) DeleteIssueRelation(ctx context.Context, id string) error {
	query := `
		mutation IssueRelationDelete($id: String!) {
			issueRelationDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		IssueRelationDelete struct {
			Success bool `json:"success"`
		} `json:"issueRelationDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.IssueRelationDelete.Success {
		return fmt.Errorf("failed to delete issue relation")
	}

	return nil
}

// DeleteProjectMilestone deletes a project milestone
func (c *Client) DeleteProjectMilestone(ctx context.Context, id string) error {
	query := `
		mutation ProjectMilestoneDelete($id: String!) {
			projectMilestoneDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		ProjectMilestoneDelete struct {
			Success bool `json:"success"`
		} `json:"projectMilestoneDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.ProjectMilestoneDelete.Success {
		return fmt.Errorf("failed to delete project milestone")
	}

	return nil
}

// GetCustomViews returns a list of custom views
func (c *Client) GetCustomViews(ctx context.Context, filter map[string]interface{}, first int, after string) (*CustomViews, error) {
	query := `
		query CustomViews($filter: CustomViewFilter, $first: Int, $after: String) {
			customViews(filter: $filter, first: $first, after: $after) {
				nodes {
					id
					name
					description
					icon
					color
					shared
					slugId
					modelName
					filterData
					projectFilterData
					creator {
						id
						name
						email
					}
					owner {
						id
						name
						email
					}
					team {
						id
						key
						name
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

	var response struct {
		CustomViews CustomViews `json:"customViews"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CustomViews, nil
}

// GetCustomView returns a single custom view by ID
func (c *Client) GetCustomView(ctx context.Context, id string) (*CustomView, error) {
	query := `
		query CustomView($id: String!) {
			customView(id: $id) {
				id
				name
				description
				icon
				color
				shared
				slugId
				modelName
				filterData
				projectFilterData
				creator {
					id
					name
					email
				}
				owner {
					id
					name
					email
				}
				team {
					id
					key
					name
				}
				createdAt
				updatedAt
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		CustomView CustomView `json:"customView"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CustomView, nil
}

// GetCustomViewIssues returns issues matching a custom view's filters
func (c *Client) GetCustomViewIssues(ctx context.Context, viewID string, first int, after string) (*Issues, error) {
	query := `
		query CustomViewIssues($id: String!, $first: Int, $after: String) {
			customView(id: $id) {
				issues(first: $first, after: $after) {
					nodes {
						id
						identifier
						title
						description
						priority
						estimate
						createdAt
						updatedAt
						dueDate
						url
						state {
							id
							name
							type
							color
						}
						assignee {
							id
							name
							email
						}
						team {
							id
							key
							name
						}
						labels {
							nodes {
								id
								name
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
		}
	`

	variables := map[string]interface{}{
		"id":    viewID,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		CustomView struct {
			Issues Issues `json:"issues"`
		} `json:"customView"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CustomView.Issues, nil
}

// GetCustomViewProjects returns projects matching a custom view's filters
func (c *Client) GetCustomViewProjects(ctx context.Context, viewID string, first int, after string) (*Projects, error) {
	query := `
		query CustomViewProjects($id: String!, $first: Int, $after: String) {
			customView(id: $id) {
				projects(first: $first, after: $after) {
					nodes {
						id
						name
						description
						state
						progress
						startDate
						targetDate
						url
						createdAt
						updatedAt
						lead {
							id
							name
							email
						}
						teams {
							nodes {
								id
								key
								name
							}
						}
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    viewID,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		CustomView struct {
			Projects Projects `json:"projects"`
		} `json:"customView"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CustomView.Projects, nil
}

// CreateCustomView creates a new custom view
func (c *Client) CreateCustomView(ctx context.Context, input map[string]interface{}) (*CustomView, error) {
	query := `
		mutation CustomViewCreate($input: CustomViewCreateInput!) {
			customViewCreate(input: $input) {
				customView {
					id
					name
					description
					icon
					color
					shared
					slugId
					modelName
					filterData
					projectFilterData
					creator {
						id
						name
						email
					}
					owner {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					createdAt
					updatedAt
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CustomViewCreate struct {
			CustomView CustomView `json:"customView"`
			Success    bool       `json:"success"`
		} `json:"customViewCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.CustomViewCreate.Success {
		return nil, fmt.Errorf("failed to create custom view")
	}

	return &response.CustomViewCreate.CustomView, nil
}

// UpdateCustomView updates an existing custom view
func (c *Client) UpdateCustomView(ctx context.Context, id string, input map[string]interface{}) (*CustomView, error) {
	query := `
		mutation CustomViewUpdate($id: String!, $input: CustomViewUpdateInput!) {
			customViewUpdate(id: $id, input: $input) {
				customView {
					id
					name
					description
					icon
					color
					shared
					slugId
					modelName
					filterData
					projectFilterData
					creator {
						id
						name
						email
					}
					owner {
						id
						name
						email
					}
					team {
						id
						key
						name
					}
					createdAt
					updatedAt
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		CustomViewUpdate struct {
			CustomView CustomView `json:"customView"`
			Success    bool       `json:"success"`
		} `json:"customViewUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.CustomViewUpdate.Success {
		return nil, fmt.Errorf("failed to update custom view")
	}

	return &response.CustomViewUpdate.CustomView, nil
}

// DeleteCustomView deletes a custom view
func (c *Client) DeleteCustomView(ctx context.Context, id string) error {
	query := `
		mutation CustomViewDelete($id: String!) {
			customViewDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		CustomViewDelete struct {
			Success bool `json:"success"`
		} `json:"customViewDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.CustomViewDelete.Success {
		return fmt.Errorf("failed to delete custom view")
	}

	return nil
}

// GetInitiatives returns a list of initiatives with optional filtering
func (c *Client) GetInitiatives(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string, includeArchived bool) (*Initiatives, error) {
	query := `
		query Initiatives($filter: InitiativeFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			initiatives(filter: $filter, first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {
					id
					name
					description
					status
					slugId
					color
					icon
					targetDate
					targetDateResolution
					health
					url
					createdAt
					updatedAt
					archivedAt
					completedAt
					owner {
						id
						name
						email
					}
					creator {
						id
						name
						email
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
		Initiatives Initiatives `json:"initiatives"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Initiatives, nil
}

// GetInitiative returns a single initiative by ID
func (c *Client) GetInitiative(ctx context.Context, id string) (*Initiative, error) {
	query := `
		query Initiative($id: String!) {
			initiative(id: $id) {
				id
				name
				description
				status
				slugId
				color
				icon
				content
				targetDate
				targetDateResolution
				health
				url
				createdAt
				updatedAt
				archivedAt
				completedAt
				owner {
					id
					name
					email
					avatarUrl
					displayName
					active
				}
				creator {
					id
					name
					email
					avatarUrl
					active
				}
				projects(first: 50) {
					nodes {
						id
						name
						state
						progress
						url
					}
				}
				parentInitiative {
					id
					name
					status
				}
				subInitiatives(first: 50) {
					nodes {
						id
						name
						status
						health
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Initiative Initiative `json:"initiative"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Initiative, nil
}

// CreateInitiative creates a new initiative
func (c *Client) CreateInitiative(ctx context.Context, input map[string]interface{}) (*Initiative, error) {
	query := `
		mutation InitiativeCreate($input: InitiativeCreateInput!) {
			initiativeCreate(input: $input) {
				initiative {
					id
					name
					description
					status
					color
					icon
					targetDate
					health
					url
					createdAt
					updatedAt
					owner {
						id
						name
						email
					}
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		InitiativeCreate struct {
			Initiative Initiative `json:"initiative"`
			Success    bool       `json:"success"`
		} `json:"initiativeCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.InitiativeCreate.Success {
		return nil, fmt.Errorf("failed to create initiative")
	}

	return &response.InitiativeCreate.Initiative, nil
}

// UpdateInitiative updates an existing initiative
func (c *Client) UpdateInitiative(ctx context.Context, id string, input map[string]interface{}) (*Initiative, error) {
	query := `
		mutation InitiativeUpdate($id: String!, $input: InitiativeUpdateInput!) {
			initiativeUpdate(id: $id, input: $input) {
				initiative {
					id
					name
					description
					status
					color
					icon
					targetDate
					health
					url
					createdAt
					updatedAt
					owner {
						id
						name
						email
					}
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		InitiativeUpdate struct {
			Initiative Initiative `json:"initiative"`
			Success    bool       `json:"success"`
		} `json:"initiativeUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.InitiativeUpdate.Success {
		return nil, fmt.Errorf("failed to update initiative")
	}

	return &response.InitiativeUpdate.Initiative, nil
}

// DeleteInitiative deletes an initiative
func (c *Client) DeleteInitiative(ctx context.Context, id string) error {
	query := `
		mutation InitiativeDelete($id: String!) {
			initiativeDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		InitiativeDelete struct {
			Success bool `json:"success"`
		} `json:"initiativeDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.InitiativeDelete.Success {
		return fmt.Errorf("failed to delete initiative")
	}

	return nil
}

// GetIssueAttachments returns attachments for a specific issue
func (c *Client) GetIssueAttachments(ctx context.Context, issueID string, first int, after string) (*Attachments, error) {
	query := `
		query IssueAttachments($id: String!, $first: Int, $after: String) {
			issue(id: $id) {
				attachments(first: $first, after: $after) {
					nodes {
						id
						title
						subtitle
						url
						metadata
						createdAt
						creator {
							id
							name
							email
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    issueID,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		Issue struct {
			Attachments Attachments `json:"attachments"`
		} `json:"issue"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issue.Attachments, nil
}

// CreateAttachment creates a generic attachment on an issue
func (c *Client) CreateAttachment(ctx context.Context, input map[string]interface{}) (*Attachment, error) {
	query := `
		mutation AttachmentCreate($input: AttachmentCreateInput!) {
			attachmentCreate(input: $input) {
				attachment {
					id
					title
					subtitle
					url
					metadata
					createdAt
					creator {
						id
						name
						email
					}
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		AttachmentCreate struct {
			Attachment Attachment `json:"attachment"`
			Success    bool       `json:"success"`
		} `json:"attachmentCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.AttachmentCreate.Success {
		return nil, fmt.Errorf("failed to create attachment")
	}

	return &response.AttachmentCreate.Attachment, nil
}

// LinkURL creates a smart link attachment on an issue (auto-detects type)
func (c *Client) LinkURL(ctx context.Context, issueID string, url string, title string) (*Attachment, error) {
	query := `
		mutation AttachmentLinkURL($issueId: String!, $url: String!, $title: String) {
			attachmentLinkURL(issueId: $issueId, url: $url, title: $title) {
				attachment {
					id
					title
					subtitle
					url
					metadata
					createdAt
					creator {
						id
						name
						email
					}
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"issueId": issueID,
		"url":     url,
	}
	if title != "" {
		variables["title"] = title
	}

	var response struct {
		AttachmentLinkURL struct {
			Attachment Attachment `json:"attachment"`
			Success    bool       `json:"success"`
		} `json:"attachmentLinkURL"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.AttachmentLinkURL.Success {
		return nil, fmt.Errorf("failed to link URL")
	}

	return &response.AttachmentLinkURL.Attachment, nil
}

// UpdateAttachment updates an existing attachment
func (c *Client) UpdateAttachment(ctx context.Context, id string, input map[string]interface{}) (*Attachment, error) {
	query := `
		mutation AttachmentUpdate($id: String!, $input: AttachmentUpdateInput!) {
			attachmentUpdate(id: $id, input: $input) {
				attachment {
					id
					title
					subtitle
					url
					metadata
					createdAt
					creator {
						id
						name
						email
					}
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		AttachmentUpdate struct {
			Attachment Attachment `json:"attachment"`
			Success    bool       `json:"success"`
		} `json:"attachmentUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.AttachmentUpdate.Success {
		return nil, fmt.Errorf("failed to update attachment")
	}

	return &response.AttachmentUpdate.Attachment, nil
}

// DeleteAttachment deletes an attachment
func (c *Client) DeleteAttachment(ctx context.Context, id string) error {
	query := `
		mutation AttachmentDelete($id: String!) {
			attachmentDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		AttachmentDelete struct {
			Success bool `json:"success"`
		} `json:"attachmentDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.AttachmentDelete.Success {
		return fmt.Errorf("failed to delete attachment")
	}

	return nil
}

// GetIssueActivity returns an issue with expanded history for activity timeline
func (c *Client) GetIssueActivity(ctx context.Context, issueID string, historyFirst int) (*Issue, error) {
	query := `
		query IssueActivity($id: String!, $historyFirst: Int) {
			issue(id: $id) {
				id
				identifier
				title
				url
				state {
					id
					name
					type
					color
				}
				assignee {
					id
					name
					email
				}
				team {
					id
					key
					name
				}
				attachments(first: 50) {
					nodes {
						id
						title
						subtitle
						url
						metadata
						createdAt
						creator {
							name
							email
						}
					}
				}
				relations {
					nodes {
						id
						type
						relatedIssue {
							id
							identifier
							title
							state {
								name
								type
							}
						}
					}
				}
				comments(first: 10) {
					nodes {
						id
						body
						createdAt
						user {
							name
							email
						}
					}
				}
				history(first: $historyFirst) {
					nodes {
						id
						createdAt
						updatedAt
						actor {
							name
							email
						}
						fromAssignee {
							name
						}
						toAssignee {
							name
						}
						fromState {
							name
						}
						toState {
							name
						}
						fromPriority
						toPriority
						fromTitle
						toTitle
						fromCycle {
							name
						}
						toCycle {
							name
						}
						fromProject {
							name
						}
						toProject {
							name
						}
						addedLabelIds
						removedLabelIds
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":           issueID,
		"historyFirst": historyFirst,
	}

	var response struct {
		Issue Issue `json:"issue"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Issue, nil
}

// GetCycles returns cycles with optional filter
func (c *Client) GetCycles(ctx context.Context, filter map[string]interface{}, first int, after string) (*Cycles, error) {
	query := `
		query Cycles($filter: CycleFilter, $first: Int, $after: String) {
			cycles(filter: $filter, first: $first, after: $after, orderBy: createdAt) {
				nodes {
					id
					number
					name
					description
					startsAt
					endsAt
					progress
					completedAt
					createdAt
					updatedAt
					archivedAt
					autoArchivedAt
					isActive
					isFuture
					isPast
					isNext
					isPrevious
					scopeHistory
					completedScopeHistory
					inProgressScopeHistory
					issueCountHistory
					completedIssueCountHistory
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
	if len(filter) > 0 {
		variables["filter"] = filter
	}

	var response struct {
		Cycles Cycles `json:"cycles"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Cycles, nil
}

// GetCycle returns a single cycle by ID
func (c *Client) GetCycle(ctx context.Context, id string) (*Cycle, error) {
	query := `
		query Cycle($id: String!) {
			cycle(id: $id) {
				id
				number
				name
				description
				startsAt
				endsAt
				progress
				completedAt
				createdAt
				updatedAt
				archivedAt
				autoArchivedAt
				isActive
				isFuture
				isPast
				isNext
				isPrevious
				scopeHistory
				completedScopeHistory
				inProgressScopeHistory
				issueCountHistory
				completedIssueCountHistory
				team {
					id
					key
					name
				}
				issues {
					nodes {
						id
						identifier
						title
						priority
						state {
							id
							name
							type
							color
						}
						assignee {
							id
							name
							email
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Cycle Cycle `json:"cycle"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Cycle, nil
}

// CreateCycle creates a new cycle for a team
func (c *Client) CreateCycle(ctx context.Context, input map[string]interface{}) (*Cycle, error) {
	query := `
		mutation CreateCycle($input: CycleCreateInput!) {
			cycleCreate(input: $input) {
				cycle {
					id
					number
					name
					description
					startsAt
					endsAt
					progress
					completedAt
					createdAt
					updatedAt
					isActive
					isFuture
					isPast
					isNext
					isPrevious
					team {
						id
						key
						name
					}
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CycleCreate struct {
			Cycle   Cycle `json:"cycle"`
			Success bool  `json:"success"`
		} `json:"cycleCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CycleCreate.Cycle, nil
}

// UpdateTeam updates a team's settings
func (c *Client) UpdateTeam(ctx context.Context, id string, input map[string]interface{}) (*Team, error) {
	query := `
		mutation UpdateTeam($id: String!, $input: TeamUpdateInput!) {
			teamUpdate(id: $id, input: $input) {
				team {
					id
					key
					name
					displayName
					description
					icon
					color
					private
					timezone
					cyclesEnabled
					cycleStartDay
					cycleDuration
					cycleCooldownTime
					cycleIssueAutoAssignStarted
					cycleIssueAutoAssignCompleted
					cycleLockToActive
					upcomingCycleCount
					triageEnabled
					requirePriorityToLeaveTriage
					inheritIssueEstimation
					issueEstimationType
					issueEstimationAllowZero
					issueEstimationExtended
					defaultIssueEstimate
					setIssueSortOrderOnStateChange
					autoClosePeriod
					autoCloseStateId
					autoArchivePeriod
					autoCloseParentIssues
					autoCloseChildIssues
					inheritWorkflowStatuses
					groupIssueHistory
					joinByDefault
					allMembersCanJoin
					scimManaged
					aiThreadSummariesEnabled
					aiDiscussionSummariesEnabled
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		TeamUpdate struct {
			Team    Team `json:"team"`
			Success bool `json:"success"`
		} `json:"teamUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.TeamUpdate.Team, nil
}

// GetLabels returns labels with optional team filter
func (c *Client) GetLabels(ctx context.Context, filter map[string]interface{}, first int, after string) (*Labels, error) {
	query := `
		query Labels($filter: IssueLabelFilter, $first: Int, $after: String) {
			issueLabels(filter: $filter, first: $first, after: $after) {
				nodes {
					id
					name
					description
					color
					parent {
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
	if len(filter) > 0 {
		variables["filter"] = filter
	}

	var response struct {
		IssueLabels Labels `json:"issueLabels"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.IssueLabels, nil
}

// CreateLabel creates a new label
func (c *Client) CreateLabel(ctx context.Context, input map[string]interface{}) (*Label, error) {
	query := `
		mutation CreateLabel($input: IssueLabelCreateInput!) {
			issueLabelCreate(input: $input) {
				issueLabel {
					id
					name
					description
					color
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		IssueLabelCreate struct {
			IssueLabel Label `json:"issueLabel"`
			Success    bool  `json:"success"`
		} `json:"issueLabelCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.IssueLabelCreate.IssueLabel, nil
}

// GetInitiativeProjects returns projects for a specific initiative
func (c *Client) GetInitiativeProjects(ctx context.Context, initiativeID string, first int, after string) (*Projects, error) {
	query := `
		query InitiativeProjects($id: String!, $first: Int, $after: String) {
			initiative(id: $id) {
				projects(first: $first, after: $after) {
					nodes {
						id
						name
						state
						progress
						startDate
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
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    initiativeID,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		Initiative struct {
			Projects Projects `json:"projects"`
		} `json:"initiative"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Initiative.Projects, nil
}

// InitiativeToProject represents a link between an initiative and a project
type InitiativeToProject struct {
	ID         string   `json:"id"`
	Initiative *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"initiative"`
	Project   *Project `json:"project"`
	SortOrder string   `json:"sortOrder"`
}

// AddProjectToInitiative links a project to an initiative
func (c *Client) AddProjectToInitiative(ctx context.Context, initiativeID string, projectID string) (*InitiativeToProject, error) {
	query := `
		mutation InitiativeToProjectCreate($input: InitiativeToProjectCreateInput!) {
			initiativeToProjectCreate(input: $input) {
				initiativeToProject {
					id
					initiative {
						id
						name
					}
					project {
						id
						name
						state
						progress
					}
					sortOrder
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"initiativeId": initiativeID,
			"projectId":    projectID,
		},
	}

	var response struct {
		InitiativeToProjectCreate struct {
			InitiativeToProject InitiativeToProject `json:"initiativeToProject"`
			Success             bool                `json:"success"`
		} `json:"initiativeToProjectCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	if !response.InitiativeToProjectCreate.Success {
		return nil, fmt.Errorf("failed to add project to initiative")
	}

	return &response.InitiativeToProjectCreate.InitiativeToProject, nil
}

// GetInitiativeToProjectLink finds the link ID between an initiative and a project
func (c *Client) GetInitiativeToProjectLink(ctx context.Context, initiativeID string, projectID string) (string, error) {
	// Query all initiative-to-project links and filter for the matching pair
	// Linear API doesn't support direct filtering on this endpoint
	query := `
		query InitiativeToProjects($first: Int) {
			initiativeToProjects(first: $first) {
				nodes {
					id
					initiative {
						id
					}
					project {
						id
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"first": 250, // Reasonable limit
	}

	var response struct {
		InitiativeToProjects struct {
			Nodes []struct {
				ID         string `json:"id"`
				Initiative struct {
					ID string `json:"id"`
				} `json:"initiative"`
				Project struct {
					ID string `json:"id"`
				} `json:"project"`
			} `json:"nodes"`
		} `json:"initiativeToProjects"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return "", err
	}

	for _, link := range response.InitiativeToProjects.Nodes {
		if link.Initiative.ID == initiativeID && link.Project.ID == projectID {
			return link.ID, nil
		}
	}

	return "", fmt.Errorf("no link found between initiative %s and project %s", initiativeID, projectID)
}

// RemoveProjectFromInitiative unlinks a project from an initiative
func (c *Client) RemoveProjectFromInitiative(ctx context.Context, initiativeID string, projectID string) error {
	// First, find the link ID
	linkID, err := c.GetInitiativeToProjectLink(ctx, initiativeID, projectID)
	if err != nil {
		return err
	}

	// Delete by link ID
	query := `
		mutation InitiativeToProjectDelete($id: String!) {
			initiativeToProjectDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": linkID,
	}

	var response struct {
		InitiativeToProjectDelete struct {
			Success bool `json:"success"`
		} `json:"initiativeToProjectDelete"`
	}

	err = c.Execute(ctx, query, variables, &response)
	if err != nil {
		return err
	}

	if !response.InitiativeToProjectDelete.Success {
		return fmt.Errorf("failed to remove project from initiative")
	}

	return nil
}

// GetProjectIssues returns issues for a specific project
func (c *Client) GetProjectIssues(ctx context.Context, projectID string, first int, after string) (*Issues, error) {
	query := `
		query ProjectIssues($id: String!, $first: Int, $after: String) {
			project(id: $id) {
				issues(first: $first, after: $after) {
					nodes {
						id
						identifier
						title
						priority
						priorityLabel
						createdAt
						updatedAt
						state {
							id
							name
							type
							color
						}
						assignee {
							id
							name
							email
						}
					}
					pageInfo {
						hasNextPage
						endCursor
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    projectID,
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		Project struct {
			Issues Issues `json:"issues"`
		} `json:"project"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Project.Issues, nil
}

// UpdateComment updates an existing comment
func (c *Client) UpdateComment(ctx context.Context, commentID string, body string) (*Comment, error) {
	query := `
		mutation UpdateComment($id: String!, $input: CommentUpdateInput!) {
			commentUpdate(id: $id, input: $input) {
				comment {
					id
					body
					createdAt
					updatedAt
					editedAt
					user {
						id
						name
						email
					}
				}
				success
			}
		}
	`
	variables := map[string]interface{}{
		"id":    commentID,
		"input": map[string]interface{}{"body": body},
	}
	var response struct {
		CommentUpdate struct {
			Comment Comment `json:"comment"`
			Success bool    `json:"success"`
		} `json:"commentUpdate"`
	}
	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}
	return &response.CommentUpdate.Comment, nil
}

// DeleteComment deletes a comment
func (c *Client) DeleteComment(ctx context.Context, commentID string) error {
	query := `
		mutation DeleteComment($id: String!) {
			commentDelete(id: $id) {
				success
			}
		}
	`
	variables := map[string]interface{}{"id": commentID}
	var response struct {
		CommentDelete struct {
			Success bool `json:"success"`
		} `json:"commentDelete"`
	}
	err := c.Execute(ctx, query, variables, &response)
	return err
}

// UpdateLabel updates an existing label
func (c *Client) UpdateLabel(ctx context.Context, id string, input map[string]interface{}) (*Label, error) {
	query := `
		mutation UpdateLabel($id: String!, $input: IssueLabelUpdateInput!) {
			issueLabelUpdate(id: $id, input: $input) {
				issueLabel {
					id
					name
					description
					color
				}
				success
			}
		}
	`
	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}
	var response struct {
		IssueLabelUpdate struct {
			IssueLabel Label `json:"issueLabel"`
			Success    bool  `json:"success"`
		} `json:"issueLabelUpdate"`
	}
	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}
	return &response.IssueLabelUpdate.IssueLabel, nil
}

// DeleteLabel deletes a label
func (c *Client) DeleteLabel(ctx context.Context, id string) error {
	query := `
		mutation DeleteLabel($id: String!) {
			issueLabelDelete(id: $id) {
				success
			}
		}
	`
	variables := map[string]interface{}{"id": id}
	var response struct {
		IssueLabelDelete struct {
			Success bool `json:"success"`
		} `json:"issueLabelDelete"`
	}
	err := c.Execute(ctx, query, variables, &response)
	return err
}

// ArchiveIssue archives an issue (soft delete)
func (c *Client) ArchiveIssue(ctx context.Context, issueID string) (*Issue, error) {
	query := `
		mutation ArchiveIssue($id: String!) {
			issueArchive(id: $id) {
				success
			}
		}
	`
	variables := map[string]interface{}{"id": issueID}
	var response struct {
		IssueArchive struct {
			Success bool `json:"success"`
		} `json:"issueArchive"`
	}
	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(ctx context.Context, input map[string]interface{}) (*Project, error) {
	query := `
		mutation CreateProject($input: ProjectCreateInput!) {
			projectCreate(input: $input) {
				project {
					id
					name
					description
					state
					progress
					startDate
					targetDate
					url
					slugId
					icon
					color
					priority
					priorityLabel
					health
					scope
					createdAt
					lead {
						id
						name
						email
					}
					creator {
						id
						name
						email
					}
					members {
						nodes {
							id
							name
							email
						}
					}
					teams {
						nodes {
							id
							key
							name
						}
					}
				}
				success
			}
		}
	`
	variables := map[string]interface{}{
		"input": input,
	}
	var response struct {
		ProjectCreate struct {
			Project Project `json:"project"`
			Success bool    `json:"success"`
		} `json:"projectCreate"`
	}
	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}
	return &response.ProjectCreate.Project, nil
}

// ArchiveProject archives a project
func (c *Client) ArchiveProject(ctx context.Context, id string) error {
	query := `
		mutation ArchiveProject($id: String!) {
			projectArchive(id: $id) {
				success
			}
		}
	`
	variables := map[string]interface{}{"id": id}
	var response struct {
		ProjectArchive struct {
			Success bool `json:"success"`
		} `json:"projectArchive"`
	}
	return c.Execute(ctx, query, variables, &response)
}

// DeleteProject permanently deletes a project
func (c *Client) DeleteProject(ctx context.Context, id string) error {
	query := `
		mutation DeleteProject($id: String!) {
			projectDelete(id: $id) {
				success
			}
		}
	`
	variables := map[string]interface{}{"id": id}
	var response struct {
		ProjectDelete struct {
			Success bool `json:"success"`
		} `json:"projectDelete"`
	}
	return c.Execute(ctx, query, variables, &response)
}

// UpdateCycle updates an existing cycle
func (c *Client) UpdateCycle(ctx context.Context, id string, input map[string]interface{}) (*Cycle, error) {
	query := `
		mutation UpdateCycle($id: String!, $input: CycleUpdateInput!) {
			cycleUpdate(id: $id, input: $input) {
				cycle {
					id
					number
					name
					description
					startsAt
					endsAt
					progress
					completedAt
					createdAt
					updatedAt
					isActive
					isFuture
					isPast
					isNext
					isPrevious
					team {
						id
						key
						name
					}
				}
				success
			}
		}
	`
	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}
	var response struct {
		CycleUpdate struct {
			Cycle   Cycle `json:"cycle"`
			Success bool  `json:"success"`
		} `json:"cycleUpdate"`
	}
	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}
	return &response.CycleUpdate.Cycle, nil
}

// ArchiveCycle archives a cycle
func (c *Client) ArchiveCycle(ctx context.Context, id string) error {
	query := `
		mutation ArchiveCycle($id: String!) {
			cycleArchive(id: $id) {
				success
			}
		}
	`
	variables := map[string]interface{}{"id": id}
	var response struct {
		CycleArchive struct {
			Success bool `json:"success"`
		} `json:"cycleArchive"`
	}
	return c.Execute(ctx, query, variables, &response)
}

// GetNotifications returns the user's notifications (inbox items)
func (c *Client) GetNotifications(ctx context.Context, first int, after string, includeArchived bool) (*Notifications, error) {
	query := `
		query Notifications($first: Int, $after: String, $includeArchived: Boolean) {
			notifications(first: $first, after: $after, includeArchived: $includeArchived) {
				nodes {
					id
					type
					createdAt
					readAt
					snoozedUntilAt
					archivedAt
					actor {
						id
						name
						email
					}
					... on IssueNotification {
						issue {
							id
							identifier
							title
							state {
								id
								name
								type
								color
							}
							team {
								id
								key
								name
							}
						}
						commentId
						reactionEmoji
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
	var response struct {
		Notifications Notifications `json:"notifications"`
	}
	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}
	return &response.Notifications, nil
}

// Favorite represents a Linear favorite (sidebar shortcut)
type Favorite struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Title      string     `json:"title"`
	Detail     string     `json:"detail"`
	FolderName string     `json:"folderName"`
	SortOrder  float64    `json:"sortOrder"`
	Color      string     `json:"color"`
	Icon       string     `json:"icon"`
	URL        string     `json:"url"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
	// Parent folder reference
	Parent *Favorite `json:"parent"`
	// Children (for folders)
	Children *Favorites `json:"children"`
	// Owner
	Owner *User `json:"owner"`
	// Referenced entities (nullable)
	Issue       *Issue       `json:"issue"`
	Project     *Project     `json:"project"`
	Cycle       *Cycle       `json:"cycle"`
	CustomView  *CustomView  `json:"customView"`
	Document    *Document    `json:"document"`
	Initiative  *Initiative  `json:"initiative"`
	Label       *Label       `json:"label"`
	User        *User        `json:"user"`
	PullRequest *PullRequest `json:"pullRequest"`
}

// PullRequest represents a Linear pull request (simplified)
type PullRequest struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Number int    `json:"number"`
	URL    string `json:"url"`
}

// Favorites represents a paginated list of favorites
type Favorites struct {
	Nodes    []Favorite `json:"nodes"`
	PageInfo PageInfo   `json:"pageInfo"`
}

// GetFavorites returns the current user's favorites
func (c *Client) GetFavorites(ctx context.Context, first int, after string) (*Favorites, error) {
	query := `
		query Favorites($first: Int, $after: String) {
			favorites(first: $first, after: $after) {
				nodes {
					id
					type
					title
					detail
					folderName
					sortOrder
					color
					icon
					url
					createdAt
					updatedAt
					archivedAt
					parent {
						id
						title
						folderName
					}
					issue {
						id
						identifier
						title
					}
					project {
						id
						name
						state
					}
					cycle {
						id
						name
						number
					}
					customView {
						id
						name
						modelName
					}
					document {
						id
						title
					}
					initiative {
						id
						name
					}
					label {
						id
						name
						color
					}
					user {
						id
						name
						email
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

	var response struct {
		Favorites Favorites `json:"favorites"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Favorites, nil
}

// GetFavorite returns a single favorite by ID
func (c *Client) GetFavorite(ctx context.Context, id string) (*Favorite, error) {
	query := `
		query Favorite($id: String!) {
			favorite(id: $id) {
				id
				type
				title
				detail
				folderName
				sortOrder
				color
				icon
				url
				createdAt
				updatedAt
				archivedAt
				parent {
					id
					title
					folderName
				}
				children {
					nodes {
						id
						type
						title
						sortOrder
					}
				}
				issue {
					id
					identifier
					title
				}
				project {
					id
					name
					state
				}
				cycle {
					id
					name
					number
				}
				customView {
					id
					name
					modelName
				}
				document {
					id
					title
				}
				initiative {
					id
					name
				}
				label {
					id
					name
					color
				}
				user {
					id
					name
					email
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Favorite Favorite `json:"favorite"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Favorite, nil
}

// CreateFavorite creates a new favorite
func (c *Client) CreateFavorite(ctx context.Context, input map[string]interface{}) (*Favorite, error) {
	query := `
		mutation FavoriteCreate($input: FavoriteCreateInput!) {
			favoriteCreate(input: $input) {
				success
				favorite {
					id
					type
					title
					detail
					folderName
					sortOrder
					color
					icon
					url
					createdAt
					issue {
						id
						identifier
						title
					}
					project {
						id
						name
					}
					cycle {
						id
						name
						number
					}
					customView {
						id
						name
					}
					document {
						id
						title
					}
					initiative {
						id
						name
					}
					label {
						id
						name
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		FavoriteCreate struct {
			Success  bool     `json:"success"`
			Favorite Favorite `json:"favorite"`
		} `json:"favoriteCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.FavoriteCreate.Favorite, nil
}

// UpdateFavorite updates an existing favorite
func (c *Client) UpdateFavorite(ctx context.Context, id string, input map[string]interface{}) (*Favorite, error) {
	query := `
		mutation FavoriteUpdate($id: String!, $input: FavoriteUpdateInput!) {
			favoriteUpdate(id: $id, input: $input) {
				success
				favorite {
					id
					type
					title
					detail
					folderName
					sortOrder
					color
					icon
					url
					parent {
						id
						title
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id":    id,
		"input": input,
	}

	var response struct {
		FavoriteUpdate struct {
			Success  bool     `json:"success"`
			Favorite Favorite `json:"favorite"`
		} `json:"favoriteUpdate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.FavoriteUpdate.Favorite, nil
}

// DeleteFavorite removes a favorite
func (c *Client) DeleteFavorite(ctx context.Context, id string) error {
	query := `
		mutation FavoriteDelete($id: String!) {
			favoriteDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		FavoriteDelete struct {
			Success bool `json:"success"`
		} `json:"favoriteDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	return err
}

// CreateTeam creates a new team
func (c *Client) CreateTeam(ctx context.Context, input map[string]interface{}, copySettingsFromTeamID string) (*Team, error) {
	query := `
		mutation CreateTeam($input: TeamCreateInput!, $copySettingsFromTeamId: String) {
			teamCreate(input: $input, copySettingsFromTeamId: $copySettingsFromTeamId) {
				team {
					id
					key
					name
					displayName
					description
					icon
					color
					private
					timezone
					cyclesEnabled
					cycleStartDay
					cycleDuration
					cycleCooldownTime
					cycleIssueAutoAssignStarted
					cycleIssueAutoAssignCompleted
					cycleLockToActive
					upcomingCycleCount
					triageEnabled
					requirePriorityToLeaveTriage
					inheritIssueEstimation
					issueEstimationType
					issueEstimationAllowZero
					issueEstimationExtended
					defaultIssueEstimate
					setIssueSortOrderOnStateChange
					autoClosePeriod
					autoCloseStateId
					autoArchivePeriod
					inheritWorkflowStatuses
					groupIssueHistory
					joinByDefault
					allMembersCanJoin
					scimManaged
					aiThreadSummariesEnabled
					aiDiscussionSummariesEnabled
				}
				success
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}
	if copySettingsFromTeamID != "" {
		variables["copySettingsFromTeamId"] = copySettingsFromTeamID
	}

	var response struct {
		TeamCreate struct {
			Team    Team `json:"team"`
			Success bool `json:"success"`
		} `json:"teamCreate"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.TeamCreate.Team, nil
}

// DeleteTeam deletes (retires) a team with a 30-day grace period
func (c *Client) DeleteTeam(ctx context.Context, id string) error {
	query := `
		mutation DeleteTeam($id: String!) {
			teamDelete(id: $id) {
				success
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		TeamDelete struct {
			Success bool `json:"success"`
		} `json:"teamDelete"`
	}

	err := c.Execute(ctx, query, variables, &response)
	return err
}
