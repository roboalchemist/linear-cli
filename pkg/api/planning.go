package api

import (
	"context"
	"time"
)

// This file adds read/list support for Linear "planning" entities:
// roadmaps, project labels, initiative labels, project relations and
// initiative relations. Read-only; no mutations live here.

// ---------------------------------------------------------------------------
// Roadmaps (deprecated in Linear in favor of initiatives, but still readable)
// ---------------------------------------------------------------------------

// RoadmapDetails is a paginated list of roadmaps.
type RoadmapDetails struct {
	Nodes    []RoadmapDetail `json:"nodes"`
	PageInfo PageInfo        `json:"pageInfo"`
}

// RoadmapDetail is a roadmap with the owner/slug/projects fields that the
// minimal api.Roadmap type does not carry.
type RoadmapDetail struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SlugId      string     `json:"slugId"`
	Color       string     `json:"color"`
	SortOrder   float64    `json:"sortOrder"`
	URL         string     `json:"url"`
	Owner       *User      `json:"owner"`
	Creator     *User      `json:"creator"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
	Projects    *Projects  `json:"projects"`
}

// roadmapDetailFields is the shared field selection for a single roadmap.
const roadmapDetailFields = `
	id
	name
	description
	slugId
	color
	sortOrder
	url
	createdAt
	updatedAt
	archivedAt
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
`

// GetRoadmaps returns all roadmaps in the workspace.
func (c *Client) GetRoadmaps(ctx context.Context, first int, after string, orderBy string, includeArchived bool) (*RoadmapDetails, error) {
	query := `
		query Roadmaps($first: Int, $after: String, $orderBy: PaginationOrderBy) {
			roadmaps(first: $first, after: $after, orderBy: $orderBy) {
				nodes {
					id
					name
					description
					slugId
					color
					sortOrder
					url
					createdAt
					updatedAt
					archivedAt
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
		"first": first,
	}
	if after != "" {
		variables["after"] = after
	}
	if orderBy != "" {
		variables["orderBy"] = orderBy
	}

	var response struct {
		Roadmaps RoadmapDetails `json:"roadmaps"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Roadmaps, nil
}

// GetRoadmap returns a single roadmap by ID, including its linked projects.
func (c *Client) GetRoadmap(ctx context.Context, id string) (*RoadmapDetail, error) {
	query := `
		query Roadmap($id: String!) {
			roadmap(id: $id) {
				id
				name
				description
				slugId
				color
				sortOrder
				url
				createdAt
				updatedAt
				archivedAt
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
				projects {
					nodes {
						id
						name
						state
						progress
						url
						lead {
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
		Roadmap RoadmapDetail `json:"roadmap"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Roadmap, nil
}

// ---------------------------------------------------------------------------
// Project labels
// ---------------------------------------------------------------------------

// ProjectLabelDetail is a project label including its team scope, which the
// pre-existing api.ProjectLabel type does not carry.
type ProjectLabelDetail struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Color       string     `json:"color"`
	Description *string    `json:"description"`
	IsGroup     bool       `json:"isGroup"`
	Team        *Team      `json:"team"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
}

// ProjectLabels is a paginated list of project labels.
type ProjectLabels struct {
	Nodes    []ProjectLabelDetail `json:"nodes"`
	PageInfo PageInfo             `json:"pageInfo"`
}

// GetProjectLabels returns all project labels in the workspace.
func (c *Client) GetProjectLabels(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string, includeArchived bool) (*ProjectLabels, error) {
	query := `
		query ProjectLabels($filter: ProjectLabelFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			projectLabels(filter: $filter, first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {
					id
					name
					color
					description
					isGroup
					createdAt
					updatedAt
					archivedAt
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
		ProjectLabels ProjectLabels `json:"projectLabels"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ProjectLabels, nil
}

// ---------------------------------------------------------------------------
// Initiative labels
// ---------------------------------------------------------------------------

// InitiativeLabel represents a Linear initiative label.
type InitiativeLabel struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Color       string     `json:"color"`
	Description *string    `json:"description"`
	IsGroup     bool       `json:"isGroup"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ArchivedAt  *time.Time `json:"archivedAt"`
}

// InitiativeLabels is a paginated list of initiative labels.
type InitiativeLabels struct {
	Nodes    []InitiativeLabel `json:"nodes"`
	PageInfo PageInfo          `json:"pageInfo"`
}

// GetInitiativeLabels returns all initiative labels in the workspace.
func (c *Client) GetInitiativeLabels(ctx context.Context, filter map[string]interface{}, first int, after string, orderBy string, includeArchived bool) (*InitiativeLabels, error) {
	query := `
		query InitiativeLabels($filter: InitiativeLabelFilter, $first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			initiativeLabels(filter: $filter, first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {
					id
					name
					color
					description
					isGroup
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
		InitiativeLabels InitiativeLabels `json:"initiativeLabels"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.InitiativeLabels, nil
}

// ---------------------------------------------------------------------------
// Project relations
// ---------------------------------------------------------------------------

// CreateRoadmap creates a new roadmap.
func (c *Client) CreateRoadmap(ctx context.Context, input map[string]interface{}) (*RoadmapDetail, error) {
	query := `
		mutation RoadmapCreate($input: RoadmapCreateInput!) {
			roadmapCreate(input: $input) {
				success
				roadmap {` + roadmapDetailFields + `}
			}
		}
	`

	var response struct {
		RoadmapCreate struct {
			Success bool          `json:"success"`
			Roadmap RoadmapDetail `json:"roadmap"`
		} `json:"roadmapCreate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"input": input}, &response); err != nil {
		return nil, err
	}
	return &response.RoadmapCreate.Roadmap, nil
}

// UpdateRoadmap updates an existing roadmap.
func (c *Client) UpdateRoadmap(ctx context.Context, id string, input map[string]interface{}) (*RoadmapDetail, error) {
	query := `
		mutation RoadmapUpdate($id: String!, $input: RoadmapUpdateInput!) {
			roadmapUpdate(id: $id, input: $input) {
				success
				roadmap {` + roadmapDetailFields + `}
			}
		}
	`

	var response struct {
		RoadmapUpdate struct {
			Success bool          `json:"success"`
			Roadmap RoadmapDetail `json:"roadmap"`
		} `json:"roadmapUpdate"`
	}

	if err := c.Execute(ctx, query, map[string]interface{}{"id": id, "input": input}, &response); err != nil {
		return nil, err
	}
	return &response.RoadmapUpdate.Roadmap, nil
}

// DeleteRoadmap permanently deletes a roadmap.
func (c *Client) DeleteRoadmap(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { roadmapDelete(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "roadmapDelete")
}

// ArchiveRoadmap archives a roadmap.
func (c *Client) ArchiveRoadmap(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { roadmapArchive(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "roadmapArchive")
}

// ProjectRelation represents a dependency between two projects.
type ProjectRelation struct {
	ID                      string            `json:"id"`
	Type                    string            `json:"type"`
	AnchorType              string            `json:"anchorType"`
	RelatedAnchorType       string            `json:"relatedAnchorType"`
	Project                 *Project          `json:"project"`
	RelatedProject          *Project          `json:"relatedProject"`
	ProjectMilestone        *ProjectMilestone `json:"projectMilestone"`
	RelatedProjectMilestone *ProjectMilestone `json:"relatedProjectMilestone"`
	CreatedAt               time.Time         `json:"createdAt"`
	UpdatedAt               time.Time         `json:"updatedAt"`
	ArchivedAt              *time.Time        `json:"archivedAt"`
}

// ProjectRelations is a paginated list of project relations.
type ProjectRelations struct {
	Nodes    []ProjectRelation `json:"nodes"`
	PageInfo PageInfo          `json:"pageInfo"`
}

// GetProjectRelations returns all project dependency relations in the workspace.
func (c *Client) GetProjectRelations(ctx context.Context, first int, after string, orderBy string, includeArchived bool) (*ProjectRelations, error) {
	query := `
		query ProjectRelations($first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			projectRelations(first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {
					id
					type
					anchorType
					relatedAnchorType
					createdAt
					updatedAt
					archivedAt
					project {
						id
						name
						state
					}
					relatedProject {
						id
						name
						state
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
		ProjectRelations ProjectRelations `json:"projectRelations"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ProjectRelations, nil
}

// ---------------------------------------------------------------------------
// Initiative relations
// ---------------------------------------------------------------------------

// InitiativeRelation represents a parent-child relation between two initiatives.
type InitiativeRelation struct {
	ID                string      `json:"id"`
	Initiative        *Initiative `json:"initiative"`
	RelatedInitiative *Initiative `json:"relatedInitiative"`
	SortOrder         float64     `json:"sortOrder"`
	CreatedAt         time.Time   `json:"createdAt"`
	UpdatedAt         time.Time   `json:"updatedAt"`
	ArchivedAt        *time.Time  `json:"archivedAt"`
}

// InitiativeRelations is a paginated list of initiative relations.
type InitiativeRelations struct {
	Nodes    []InitiativeRelation `json:"nodes"`
	PageInfo PageInfo             `json:"pageInfo"`
}

// GetInitiativeRelations returns all initiative parent-child relations in the workspace.
func (c *Client) GetInitiativeRelations(ctx context.Context, first int, after string, orderBy string, includeArchived bool) (*InitiativeRelations, error) {
	query := `
		query InitiativeRelations($first: Int, $after: String, $orderBy: PaginationOrderBy, $includeArchived: Boolean) {
			initiativeRelations(first: $first, after: $after, orderBy: $orderBy, includeArchived: $includeArchived) {
				nodes {
					id
					sortOrder
					createdAt
					updatedAt
					archivedAt
					initiative {
						id
						name
						status
					}
					relatedInitiative {
						id
						name
						status
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
		InitiativeRelations InitiativeRelations `json:"initiativeRelations"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.InitiativeRelations, nil
}
