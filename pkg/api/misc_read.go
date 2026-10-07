package api

import (
	"context"
	"fmt"
)

// ---------------------------------------------------------------------------
// OAuth applications
// ---------------------------------------------------------------------------

type OAuthApplication struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	ClientID       string   `json:"clientId"`
	Description    string   `json:"description"`
	DeveloperURL   string   `json:"developerUrl"`
	Distribution   string   `json:"distribution"`
	WebhookURL     string   `json:"webhookUrl"`
	WebhookEnabled bool     `json:"webhookEnabled"`
	RedirectURIs   []string `json:"redirectUris"`
	GrantTypes     []string `json:"grantTypes"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

// GetOAuthApplications lists the workspace OAuth applications.
func (c *Client) GetOAuthApplications(ctx context.Context) ([]OAuthApplication, error) {
	query := `
		query {
			oauthApplications {
				id
				name
				clientId
				description
				developerUrl
				distribution
				webhookUrl
				webhookEnabled
				redirectUris
				grantTypes
				createdAt
				updatedAt
			}
		}
	`
	var response struct {
		OAuthApplications []OAuthApplication `json:"oauthApplications"`
	}
	if err := c.Execute(ctx, query, nil, &response); err != nil {
		return nil, err
	}
	return response.OAuthApplications, nil
}

// ---------------------------------------------------------------------------
// Organization invites
// ---------------------------------------------------------------------------

type OrganizationInvite struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	External   bool   `json:"external"`
	AcceptedAt string `json:"acceptedAt"`
	ExpiresAt  string `json:"expiresAt"`
	CreatedAt  string `json:"createdAt"`
	ArchivedAt string `json:"archivedAt"`
}

type OrganizationInvites struct {
	Nodes    []OrganizationInvite `json:"nodes"`
	PageInfo PageInfo             `json:"pageInfo"`
}

// GetOrganizationInvites lists pending and past organization invites.
func (c *Client) GetOrganizationInvites(ctx context.Context, first int, after string, includeArchived bool) (*OrganizationInvites, error) {
	query := `
		query Invites($first: Int, $after: String, $includeArchived: Boolean) {
			organizationInvites(first: $first, after: $after, includeArchived: $includeArchived) {
				nodes { id email role external acceptedAt expiresAt createdAt archivedAt }
				pageInfo { hasNextPage endCursor }
			}
		}
	`
	vars := map[string]interface{}{"first": first}
	if after != "" {
		vars["after"] = after
	}
	if includeArchived {
		vars["includeArchived"] = true
	}
	var response struct {
		OrganizationInvites OrganizationInvites `json:"organizationInvites"`
	}
	if err := c.Execute(ctx, query, vars, &response); err != nil {
		return nil, err
	}
	return &response.OrganizationInvites, nil
}

// ---------------------------------------------------------------------------
// Audit entry types
// ---------------------------------------------------------------------------

type AuditEntryType struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// GetAuditEntryTypes lists the audit log event types.
func (c *Client) GetAuditEntryTypes(ctx context.Context) ([]AuditEntryType, error) {
	var response struct {
		AuditEntryTypes []AuditEntryType `json:"auditEntryTypes"`
	}
	if err := c.Execute(ctx, `query { auditEntryTypes { type description } }`, nil, &response); err != nil {
		return nil, err
	}
	return response.AuditEntryTypes, nil
}

// ---------------------------------------------------------------------------
// Team memberships
// ---------------------------------------------------------------------------

type TeamMembership struct {
	ID        string  `json:"id"`
	Owner     bool    `json:"owner"`
	SortOrder float64 `json:"sortOrder"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Team      *Team   `json:"team"`
	User      *User   `json:"user"`
}

type TeamMemberships struct {
	Nodes    []TeamMembership `json:"nodes"`
	PageInfo PageInfo         `json:"pageInfo"`
}

// GetTeamMemberships lists team membership records.
func (c *Client) GetTeamMemberships(ctx context.Context, first int, after string, includeArchived bool) (*TeamMemberships, error) {
	query := `
		query Memberships($first: Int, $after: String, $includeArchived: Boolean) {
			teamMemberships(first: $first, after: $after, includeArchived: $includeArchived) {
				nodes {
					id owner sortOrder createdAt updatedAt
					team { id key name }
					user { id name email }
				}
				pageInfo { hasNextPage endCursor }
			}
		}
	`
	vars := map[string]interface{}{"first": first}
	if after != "" {
		vars["after"] = after
	}
	if includeArchived {
		vars["includeArchived"] = true
	}
	var response struct {
		TeamMemberships TeamMemberships `json:"teamMemberships"`
	}
	if err := c.Execute(ctx, query, vars, &response); err != nil {
		return nil, err
	}
	return &response.TeamMemberships, nil
}

// ---------------------------------------------------------------------------
// Webhook secret rotation
// ---------------------------------------------------------------------------

// RotateWebhookSecret rotates a webhook's signing secret.
func (c *Client) RotateWebhookSecret(ctx context.Context, id string) (*Webhook, error) {
	query := `
		mutation($id: String!) {
			webhookRotateSecret(id: $id) {
				success
				webhook { id url enabled secret }
			}
		}
	`
	var response struct {
		WebhookRotateSecret struct {
			Success bool     `json:"success"`
			Webhook *Webhook `json:"webhook"`
		} `json:"webhookRotateSecret"`
	}
	if err := c.Execute(ctx, query, map[string]interface{}{"id": id}, &response); err != nil {
		return nil, err
	}
	if !response.WebhookRotateSecret.Success {
		return nil, fmt.Errorf("webhookRotateSecret reported failure")
	}
	return response.WebhookRotateSecret.Webhook, nil
}

// ---------------------------------------------------------------------------
// Emojis
// ---------------------------------------------------------------------------

// CreateEmoji creates a custom emoji.
func (c *Client) CreateEmoji(ctx context.Context, name, url string) error {
	input := map[string]interface{}{"name": name}
	if url != "" {
		input["url"] = url
	}
	return c.mutationOK(ctx,
		`mutation($input: EmojiCreateInput!) { emojiCreate(input: $input) { success } }`,
		map[string]interface{}{"input": input}, "emojiCreate")
}

// DeleteEmoji deletes a custom emoji.
func (c *Client) DeleteEmoji(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { emojiDelete(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "emojiDelete")
}
