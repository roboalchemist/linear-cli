package api

import (
	"context"
	"time"
)

// This file implements admin/misc entities: webhooks, time schedules,
// triage responsibilities, custom emojis, templates, and audit entries.
//
// Note: pkg/api/queries.go already declares a minimal `Template` type used by
// team/project payloads. The richer admin-facing template entity is declared
// here as `AdminTemplate` to avoid a redeclaration collision.

// Webhook represents a Linear webhook subscription.
type Webhook struct {
	ID             string     `json:"id"`
	Label          *string    `json:"label"`
	URL            *string    `json:"url"`
	Enabled        bool       `json:"enabled"`
	AllPublicTeams bool       `json:"allPublicTeams"`
	ResourceTypes  []string   `json:"resourceTypes"`
	Team           *Team      `json:"team"`
	TeamIds        []string   `json:"teamIds"`
	Creator        *User      `json:"creator"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	ArchivedAt     *time.Time `json:"archivedAt"`
}

// Webhooks is a paginated list of webhooks.
type Webhooks struct {
	Nodes    []Webhook `json:"nodes"`
	PageInfo PageInfo  `json:"pageInfo"`
}

// TimeSchedule represents a Linear time schedule (used by triage shifts).
type TimeSchedule struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	ExternalId  *string    `json:"externalId"`
	ExternalUrl *string    `json:"externalUrl"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
}

// TimeSchedules is a paginated list of time schedules.
type TimeSchedules struct {
	Nodes    []TimeSchedule `json:"nodes"`
	PageInfo PageInfo       `json:"pageInfo"`
}

// TriageResponsibility represents a team's triage responsibility assignment.
type TriageResponsibility struct {
	ID           string        `json:"id"`
	Action       string        `json:"action"`
	Team         *Team         `json:"team"`
	TimeSchedule *TimeSchedule `json:"timeSchedule"`
	CurrentUser  *User         `json:"currentUser"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	ArchivedAt   *time.Time    `json:"archivedAt"`
}

// TriageResponsibilities is a paginated list of triage responsibilities.
type TriageResponsibilities struct {
	Nodes    []TriageResponsibility `json:"nodes"`
	PageInfo PageInfo               `json:"pageInfo"`
}

// Emoji represents a custom workspace emoji.
type Emoji struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	Source     string     `json:"source"`
	Creator    *User      `json:"creator"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

// Emojis is a paginated list of custom emojis.
type Emojis struct {
	Nodes    []Emoji  `json:"nodes"`
	PageInfo PageInfo `json:"pageInfo"`
}

// AdminTemplate represents a Linear template (issue/project/document/etc.).
type AdminTemplate struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Description   *string    `json:"description"`
	Content       *string    `json:"content"`
	Color         *string    `json:"color"`
	Icon          *string    `json:"icon"`
	SortOrder     float64    `json:"sortOrder"`
	HasFormFields bool       `json:"hasFormFields"`
	Team          *Team      `json:"team"`
	Creator       *User      `json:"creator"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	ArchivedAt    *time.Time `json:"archivedAt"`
}

// AuditEntry represents a single audit log entry. Reading audit entries
// requires workspace admin privileges.
type AuditEntry struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Actor       *User      `json:"actor"`
	ActorId     *string    `json:"actorId"`
	IP          *string    `json:"ip"`
	CountryCode *string    `json:"countryCode"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
}

// AuditEntries is a paginated list of audit entries.
type AuditEntries struct {
	Nodes    []AuditEntry `json:"nodes"`
	PageInfo PageInfo     `json:"pageInfo"`
}

const webhookFields = `
	id
	label
	url
	enabled
	allPublicTeams
	resourceTypes
	team { id key name }
	teamIds
	creator { id name email }
	createdAt
	updatedAt
	archivedAt
`

// GetWebhooks returns a list of webhooks for the current workspace.
func (c *Client) GetWebhooks(ctx context.Context, first int, after string) (*Webhooks, error) {
	query := `
		query Webhooks($first: Int, $after: String) {
			webhooks(first: $first, after: $after) {
				nodes {` + webhookFields + `}
				pageInfo { hasNextPage endCursor }
			}
		}
	`

	variables := map[string]interface{}{"first": first}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		Webhooks Webhooks `json:"webhooks"`
	}
	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.Webhooks, nil
}

// GetWebhook returns a single webhook by ID.
func (c *Client) GetWebhook(ctx context.Context, id string) (*Webhook, error) {
	query := `
		query Webhook($id: String!) {
			webhook(id: $id) {` + webhookFields + `}
		}
	`

	var response struct {
		Webhook Webhook `json:"webhook"`
	}
	if err := c.Execute(ctx, query, map[string]interface{}{"id": id}, &response); err != nil {
		return nil, err
	}
	return &response.Webhook, nil
}

const timeScheduleFields = `
	id
	name
	externalId
	externalUrl
	createdAt
	updatedAt
	archivedAt
`

// GetTimeSchedules returns a list of time schedules.
func (c *Client) GetTimeSchedules(ctx context.Context, first int, after string) (*TimeSchedules, error) {
	query := `
		query TimeSchedules($first: Int, $after: String) {
			timeSchedules(first: $first, after: $after) {
				nodes {` + timeScheduleFields + `}
				pageInfo { hasNextPage endCursor }
			}
		}
	`

	variables := map[string]interface{}{"first": first}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		TimeSchedules TimeSchedules `json:"timeSchedules"`
	}
	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.TimeSchedules, nil
}

// GetTimeSchedule returns a single time schedule by ID.
func (c *Client) GetTimeSchedule(ctx context.Context, id string) (*TimeSchedule, error) {
	query := `
		query TimeSchedule($id: String!) {
			timeSchedule(id: $id) {` + timeScheduleFields + `}
		}
	`

	var response struct {
		TimeSchedule TimeSchedule `json:"timeSchedule"`
	}
	if err := c.Execute(ctx, query, map[string]interface{}{"id": id}, &response); err != nil {
		return nil, err
	}
	return &response.TimeSchedule, nil
}

const triageResponsibilityFields = `
	id
	action
	team { id key name }
	timeSchedule { id name }
	currentUser { id name email }
	createdAt
	updatedAt
	archivedAt
`

// GetTriageResponsibilities returns a list of triage responsibilities.
func (c *Client) GetTriageResponsibilities(ctx context.Context, first int, after string) (*TriageResponsibilities, error) {
	query := `
		query TriageResponsibilities($first: Int, $after: String) {
			triageResponsibilities(first: $first, after: $after) {
				nodes {` + triageResponsibilityFields + `}
				pageInfo { hasNextPage endCursor }
			}
		}
	`

	variables := map[string]interface{}{"first": first}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		TriageResponsibilities TriageResponsibilities `json:"triageResponsibilities"`
	}
	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.TriageResponsibilities, nil
}

// GetTriageResponsibility returns a single triage responsibility by ID.
func (c *Client) GetTriageResponsibility(ctx context.Context, id string) (*TriageResponsibility, error) {
	query := `
		query TriageResponsibility($id: String!) {
			triageResponsibility(id: $id) {` + triageResponsibilityFields + `}
		}
	`

	var response struct {
		TriageResponsibility TriageResponsibility `json:"triageResponsibility"`
	}
	if err := c.Execute(ctx, query, map[string]interface{}{"id": id}, &response); err != nil {
		return nil, err
	}
	return &response.TriageResponsibility, nil
}

const emojiFields = `
	id
	name
	url
	source
	creator { id name email }
	createdAt
	updatedAt
	archivedAt
`

// GetEmojis returns a list of custom emojis in the workspace.
func (c *Client) GetEmojis(ctx context.Context, first int, after string) (*Emojis, error) {
	query := `
		query Emojis($first: Int, $after: String) {
			emojis(first: $first, after: $after) {
				nodes {` + emojiFields + `}
				pageInfo { hasNextPage endCursor }
			}
		}
	`

	variables := map[string]interface{}{"first": first}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		Emojis Emojis `json:"emojis"`
	}
	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.Emojis, nil
}

// GetEmoji returns a single custom emoji by ID or name.
func (c *Client) GetEmoji(ctx context.Context, idOrName string) (*Emoji, error) {
	query := `
		query Emoji($id: String!) {
			emoji(id: $id) {` + emojiFields + `}
		}
	`

	var response struct {
		Emoji Emoji `json:"emoji"`
	}
	if err := c.Execute(ctx, query, map[string]interface{}{"id": idOrName}, &response); err != nil {
		return nil, err
	}
	return &response.Emoji, nil
}

const adminTemplateFields = `
	id
	name
	type
	description
	content
	color
	icon
	sortOrder
	hasFormFields
	team { id key name }
	creator { id name email }
	createdAt
	updatedAt
	archivedAt
`

// GetTemplates returns all templates in the workspace (team-scoped and
// workspace-level). The Linear API returns this as an unpaginated list.
func (c *Client) GetTemplates(ctx context.Context) ([]AdminTemplate, error) {
	query := `
		query Templates {
			templates {` + adminTemplateFields + `}
		}
	`

	var response struct {
		Templates []AdminTemplate `json:"templates"`
	}
	if err := c.Execute(ctx, query, nil, &response); err != nil {
		return nil, err
	}
	return response.Templates, nil
}

// SearchTemplates returns templates matching a name filter.
func (c *Client) SearchTemplates(ctx context.Context, term string, first int) ([]AdminTemplate, error) {
	query := `
		query TemplateSearch($filter: TemplateFilter, $first: Int) {
			templateSearch(filter: $filter, first: $first) {` + adminTemplateFields + `}
		}
	`

	filter := map[string]interface{}{}
	if term != "" {
		filter["name"] = map[string]interface{}{"containsIgnoreCase": term}
	}

	variables := map[string]interface{}{"first": first}
	if len(filter) > 0 {
		variables["filter"] = filter
	}

	var response struct {
		TemplateSearch []AdminTemplate `json:"templateSearch"`
	}
	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return response.TemplateSearch, nil
}

// GetTemplate returns a single template by ID.
func (c *Client) GetTemplate(ctx context.Context, id string) (*AdminTemplate, error) {
	query := `
		query Template($id: String!) {
			template(id: $id) {` + adminTemplateFields + `}
		}
	`

	var response struct {
		Template AdminTemplate `json:"template"`
	}
	if err := c.Execute(ctx, query, map[string]interface{}{"id": id}, &response); err != nil {
		return nil, err
	}
	return &response.Template, nil
}

const auditEntryFields = `
	id
	type
	actor { id name email }
	actorId
	ip
	countryCode
	createdAt
	updatedAt
	archivedAt
`

// GetAuditEntries returns audit log entries. This requires workspace admin
// privileges; non-admin API keys receive a permission error.
func (c *Client) GetAuditEntries(ctx context.Context, filter map[string]interface{}, first int, after string) (*AuditEntries, error) {
	query := `
		query AuditEntries($filter: AuditEntryFilter, $first: Int, $after: String) {
			auditEntries(filter: $filter, first: $first, after: $after) {
				nodes {` + auditEntryFields + `}
				pageInfo { hasNextPage endCursor }
			}
		}
	`

	variables := map[string]interface{}{"first": first}
	if filter != nil {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		AuditEntries AuditEntries `json:"auditEntries"`
	}
	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.AuditEntries, nil
}

// ---------------------------------------------------------------------------
// Templates (mutations)
// ---------------------------------------------------------------------------

// CreateTemplate creates a new template.
func (c *Client) CreateTemplate(ctx context.Context, input map[string]interface{}) (*AdminTemplate, error) {
	query := `
		mutation TemplateCreate($input: TemplateCreateInput!) {
			templateCreate(input: $input) {
				success
				template {` + adminTemplateFields + `}
			}
		}
	`

	var response struct {
		TemplateCreate struct {
			Success  bool          `json:"success"`
			Template AdminTemplate `json:"template"`
		} `json:"templateCreate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"input": input}, &response); err != nil {
		return nil, err
	}
	return &response.TemplateCreate.Template, nil
}

// UpdateTemplate updates an existing template.
func (c *Client) UpdateTemplate(ctx context.Context, id string, input map[string]interface{}) (*AdminTemplate, error) {
	query := `
		mutation TemplateUpdate($id: String!, $input: TemplateUpdateInput!) {
			templateUpdate(id: $id, input: $input) {
				success
				template {` + adminTemplateFields + `}
			}
		}
	`

	var response struct {
		TemplateUpdate struct {
			Success  bool          `json:"success"`
			Template AdminTemplate `json:"template"`
		} `json:"templateUpdate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"id": id, "input": input}, &response); err != nil {
		return nil, err
	}
	return &response.TemplateUpdate.Template, nil
}

// DeleteTemplate deletes a template.
func (c *Client) DeleteTemplate(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { templateDelete(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "templateDelete")
}

// ---------------------------------------------------------------------------
// Webhooks (mutations)
// ---------------------------------------------------------------------------

// CreateWebhook creates a new webhook subscription.
func (c *Client) CreateWebhook(ctx context.Context, input map[string]interface{}) (*Webhook, error) {
	query := `
		mutation WebhookCreate($input: WebhookCreateInput!) {
			webhookCreate(input: $input) {
				success
				webhook {` + webhookFields + `}
			}
		}
	`

	var response struct {
		WebhookCreate struct {
			Success bool    `json:"success"`
			Webhook Webhook `json:"webhook"`
		} `json:"webhookCreate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"input": input}, &response); err != nil {
		return nil, err
	}
	return &response.WebhookCreate.Webhook, nil
}

// UpdateWebhook updates an existing webhook.
func (c *Client) UpdateWebhook(ctx context.Context, id string, input map[string]interface{}) (*Webhook, error) {
	query := `
		mutation WebhookUpdate($id: String!, $input: WebhookUpdateInput!) {
			webhookUpdate(id: $id, input: $input) {
				success
				webhook {` + webhookFields + `}
			}
		}
	`

	var response struct {
		WebhookUpdate struct {
			Success bool    `json:"success"`
			Webhook Webhook `json:"webhook"`
		} `json:"webhookUpdate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"id": id, "input": input}, &response); err != nil {
		return nil, err
	}
	return &response.WebhookUpdate.Webhook, nil
}

// DeleteWebhook deletes a webhook.
func (c *Client) DeleteWebhook(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { webhookDelete(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "webhookDelete")
}

// ---------------------------------------------------------------------------
// Time schedules (mutations)
// ---------------------------------------------------------------------------

// CreateTimeSchedule creates a new time schedule.
func (c *Client) CreateTimeSchedule(ctx context.Context, input map[string]interface{}) (*TimeSchedule, error) {
	query := `
		mutation TimeScheduleCreate($input: TimeScheduleCreateInput!) {
			timeScheduleCreate(input: $input) {
				success
				timeSchedule {` + timeScheduleFields + `}
			}
		}
	`

	var response struct {
		TimeScheduleCreate struct {
			Success      bool         `json:"success"`
			TimeSchedule TimeSchedule `json:"timeSchedule"`
		} `json:"timeScheduleCreate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"input": input}, &response); err != nil {
		return nil, err
	}
	return &response.TimeScheduleCreate.TimeSchedule, nil
}

// UpdateTimeSchedule updates an existing time schedule.
func (c *Client) UpdateTimeSchedule(ctx context.Context, id string, input map[string]interface{}) (*TimeSchedule, error) {
	query := `
		mutation TimeScheduleUpdate($id: String!, $input: TimeScheduleUpdateInput!) {
			timeScheduleUpdate(id: $id, input: $input) {
				success
				timeSchedule {` + timeScheduleFields + `}
			}
		}
	`

	var response struct {
		TimeScheduleUpdate struct {
			Success      bool         `json:"success"`
			TimeSchedule TimeSchedule `json:"timeSchedule"`
		} `json:"timeScheduleUpdate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"id": id, "input": input}, &response); err != nil {
		return nil, err
	}
	return &response.TimeScheduleUpdate.TimeSchedule, nil
}

// DeleteTimeSchedule deletes a time schedule.
func (c *Client) DeleteTimeSchedule(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { timeScheduleDelete(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "timeScheduleDelete")
}
