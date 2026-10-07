package api

import (
	"context"
	"time"
)

// Customer represents a single customer organization tracked in Linear's
// customer management system.
type Customer struct {
	ID                   string          `json:"id"`
	Name                 string          `json:"name"`
	SlugID               string          `json:"slugId"`
	ApproximateNeedCount float64         `json:"approximateNeedCount"`
	Domains              []string        `json:"domains"`
	ExternalIDs          []string        `json:"externalIds"`
	LogoURL              *string         `json:"logoUrl"`
	MainSourceID         *string         `json:"mainSourceId"`
	Revenue              *int            `json:"revenue"`
	Size                 *float64        `json:"size"`
	SlackChannelID       *string         `json:"slackChannelId"`
	Status               *CustomerStatus `json:"status"`
	Tier                 *CustomerTier   `json:"tier"`
	Owner                *User           `json:"owner"`
	Needs                []CustomerNeed  `json:"needs"`
	URL                  string          `json:"url"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`
	ArchivedAt           *time.Time      `json:"archivedAt"`
}

// CustomerNeed represents a product request or piece of feedback from a
// customer, bridged to an issue or project.
type CustomerNeed struct {
	ID         string     `json:"id"`
	Body       *string    `json:"body"`
	Content    *string    `json:"content"`
	Priority   float64    `json:"priority"`
	URL        *string    `json:"url"`
	Customer   *Customer  `json:"customer"`
	Issue      *Issue     `json:"issue"`
	Project    *Project   `json:"project"`
	Creator    *User      `json:"creator"`
	Comment    *Comment   `json:"comment"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
}

// CustomerStatus represents a workspace-defined lifecycle status for
// customers (e.g., Active, Churned, Trial).
type CustomerStatus struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	DisplayName string     `json:"displayName"`
	Color       string     `json:"color"`
	Description *string    `json:"description"`
	Type        *string    `json:"type"`
	Position    float64    `json:"position"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
}

// CustomerTier represents a tier or segment assigned to a customer for
// prioritization (e.g., Enterprise, Pro, Free).
type CustomerTier struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	DisplayName string     `json:"displayName"`
	Color       string     `json:"color"`
	Description *string    `json:"description"`
	Position    float64    `json:"position"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
}

// Customers is a paginated collection of customers.
type Customers struct {
	Nodes    []Customer `json:"nodes"`
	PageInfo PageInfo   `json:"pageInfo"`
}

// CustomerNeeds is a paginated collection of customer needs.
type CustomerNeeds struct {
	Nodes    []CustomerNeed `json:"nodes"`
	PageInfo PageInfo       `json:"pageInfo"`
}

// CustomerStatuses is a paginated collection of customer statuses.
type CustomerStatuses struct {
	Nodes    []CustomerStatus `json:"nodes"`
	PageInfo PageInfo         `json:"pageInfo"`
}

// CustomerTiers is a paginated collection of customer tiers.
type CustomerTiers struct {
	Nodes    []CustomerTier `json:"nodes"`
	PageInfo PageInfo       `json:"pageInfo"`
}

// customerFields is the shared field selection used when fetching customers.
const customerFields = `
	id
	name
	slugId
	approximateNeedCount
	domains
	externalIds
	logoUrl
	mainSourceId
	revenue
	size
	slackChannelId
	url
	createdAt
	updatedAt
	archivedAt
	status {
		id
		name
		displayName
		color
	}
	tier {
		id
		name
		displayName
		color
	}
	owner {
		id
		name
		displayName
		email
	}
`

// GetCustomers returns a list of customers.
//
// The optional filter map is passed through as the GraphQL CustomerFilter. A
// special "includeArchived" key may be set in filter to include archived
// customers; it is extracted and sent as the includeArchived argument rather
// than as part of the filter.
func (c *Client) GetCustomers(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string) (*Customers, error) {
	query := `
		query Customers($filter: CustomerFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			customers(filter: $filter, first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {` + customerFields + `}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	includeArchived := false
	if filter != nil {
		if v, ok := filter["includeArchived"].(bool); ok {
			includeArchived = v
			delete(filter, "includeArchived")
		}
	}

	variables := map[string]interface{}{
		"first": first,
	}
	if filter != nil && len(filter) > 0 {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}
	if includeArchived {
		variables["includeArchived"] = includeArchived
	}

	var response struct {
		Customers Customers `json:"customers"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Customers, nil
}

// GetCustomer returns a single customer by ID or slug.
func (c *Client) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	query := `
		query Customer($id: String!) {
			customer(id: $id) {` + customerFields + `
				needs {
					id
					body
					content
					priority
					url
					issue {
						id
						identifier
						title
					}
					project {
						id
						name
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Customer Customer `json:"customer"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Customer, nil
}

// GetCustomerNeeds returns customer needs, optionally filtered by customer.
func (c *Client) GetCustomerNeeds(ctx context.Context, customerID string, first int, after string) (*CustomerNeeds, error) {
	query := `
		query CustomerNeeds($filter: CustomerNeedFilter, $first: Int, $after: String) {
			customerNeeds(filter: $filter, first: $first, after: $after) {
				nodes {
					id
					body
					content
					priority
					url
					createdAt
					updatedAt
					archivedAt
					customer {
						id
						name
					}
					issue {
						id
						identifier
						title
						url
					}
					project {
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
	if customerID != "" {
		variables["filter"] = map[string]interface{}{
			"customer": map[string]interface{}{
				"id": map[string]interface{}{"eq": customerID},
			},
		}
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		CustomerNeeds CustomerNeeds `json:"customerNeeds"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CustomerNeeds, nil
}

// GetCustomerStatuses returns all customer statuses defined in the workspace.
func (c *Client) GetCustomerStatuses(ctx context.Context, first int, after string) (*CustomerStatuses, error) {
	query := `
		query CustomerStatuses($first: Int, $after: String) {
			customerStatuses(first: $first, after: $after) {
				nodes {
					id
					name
					displayName
					color
					description
					type
					position
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
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		CustomerStatuses CustomerStatuses `json:"customerStatuses"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CustomerStatuses, nil
}

// GetCustomerTiers returns all customer tiers defined in the workspace.
func (c *Client) GetCustomerTiers(ctx context.Context, first int, after string) (*CustomerTiers, error) {
	query := `
		query CustomerTiers($first: Int, $after: String) {
			customerTiers(first: $first, after: $after) {
				nodes {
					id
					name
					displayName
					color
					description
					position
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
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		CustomerTiers CustomerTiers `json:"customerTiers"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.CustomerTiers, nil
}

// CreateCustomer creates a new customer.
func (c *Client) CreateCustomer(ctx context.Context, input map[string]interface{}) (*Customer, error) {
	query := `
		mutation CustomerCreate($input: CustomerCreateInput!) {
			customerCreate(input: $input) {
				success
				customer {` + customerFields + `}
			}
		}
	`

	var response struct {
		CustomerCreate struct {
			Success  bool     `json:"success"`
			Customer Customer `json:"customer"`
		} `json:"customerCreate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"input": input}, &response); err != nil {
		return nil, err
	}
	return &response.CustomerCreate.Customer, nil
}

// UpdateCustomer updates an existing customer.
func (c *Client) UpdateCustomer(ctx context.Context, id string, input map[string]interface{}) (*Customer, error) {
	query := `
		mutation CustomerUpdate($id: String!, $input: CustomerUpdateInput!) {
			customerUpdate(id: $id, input: $input) {
				success
				customer {` + customerFields + `}
			}
		}
	`

	var response struct {
		CustomerUpdate struct {
			Success  bool     `json:"success"`
			Customer Customer `json:"customer"`
		} `json:"customerUpdate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"id": id, "input": input}, &response); err != nil {
		return nil, err
	}
	return &response.CustomerUpdate.Customer, nil
}

// DeleteCustomer deletes a customer.
func (c *Client) DeleteCustomer(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { customerDelete(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "customerDelete")
}

// MergeCustomer merges the source customer into the target customer, moving all
// needs to the target and archiving the source.
func (c *Client) MergeCustomer(ctx context.Context, sourceID, targetID string) (*Customer, error) {
	query := `
		mutation CustomerMerge($source: String!, $target: String!) {
			customerMerge(sourceCustomerId: $source, targetCustomerId: $target) {
				success
				customer {` + customerFields + `}
			}
		}
	`

	var response struct {
		CustomerMerge struct {
			Success  bool     `json:"success"`
			Customer Customer `json:"customer"`
		} `json:"customerMerge"`
	}

	variables := map[string]interface{}{"source": sourceID, "target": targetID}
	if err := c.Execute(ctx, query, variables, &response); err != nil {
		return nil, err
	}
	return &response.CustomerMerge.Customer, nil
}

// CreateCustomerStatus creates a new customer status.
func (c *Client) CreateCustomerStatus(ctx context.Context, input map[string]interface{}) (*CustomerStatus, error) {
	query := `
		mutation CustomerStatusCreate($input: CustomerStatusCreateInput!) {
			customerStatusCreate(input: $input) {
				success
				status {
					id
					name
					displayName
					color
					description
					type
					position
					createdAt
					updatedAt
					archivedAt
				}
			}
		}
	`

	var response struct {
		CustomerStatusCreate struct {
			Success bool           `json:"success"`
			Status  CustomerStatus `json:"status"`
		} `json:"customerStatusCreate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"input": input}, &response); err != nil {
		return nil, err
	}
	return &response.CustomerStatusCreate.Status, nil
}

// DeleteCustomerStatus deletes a customer status.
func (c *Client) DeleteCustomerStatus(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { customerStatusDelete(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "customerStatusDelete")
}

// CreateCustomerTier creates a new customer tier.
func (c *Client) CreateCustomerTier(ctx context.Context, input map[string]interface{}) (*CustomerTier, error) {
	query := `
		mutation CustomerTierCreate($input: CustomerTierCreateInput!) {
			customerTierCreate(input: $input) {
				success
				tier {
					id
					name
					displayName
					color
					description
					position
					createdAt
					updatedAt
					archivedAt
				}
			}
		}
	`

	var response struct {
		CustomerTierCreate struct {
			Success bool         `json:"success"`
			Tier    CustomerTier `json:"tier"`
		} `json:"customerTierCreate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"input": input}, &response); err != nil {
		return nil, err
	}
	return &response.CustomerTierCreate.Tier, nil
}

// DeleteCustomerTier deletes a customer tier.
func (c *Client) DeleteCustomerTier(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { customerTierDelete(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "customerTierDelete")
}
