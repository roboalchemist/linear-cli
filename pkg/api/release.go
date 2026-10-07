package api

import (
	"context"
	"time"
)

// NOTE: an existing minimal `Release` type (id/name/version, used for document
// linking) is declared in queries.go. The types below are intentionally
// distinctly named so they do not redeclare it.

// ReleaseDetail represents a full Linear release.
type ReleaseDetail struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Version         *string                `json:"version"`
	SlugID          string                 `json:"slugId"`
	URL             string                 `json:"url"`
	Description     *string                `json:"description"`
	CommitSha       *string                `json:"commitSha"`
	StartDate       *string                `json:"startDate"`
	TargetDate      *string                `json:"targetDate"`
	StartedAt       *time.Time             `json:"startedAt"`
	CompletedAt     *time.Time             `json:"completedAt"`
	CanceledAt      *time.Time             `json:"canceledAt"`
	ArchivedAt      *time.Time             `json:"archivedAt"`
	AutoArchivedAt  *time.Time             `json:"autoArchivedAt"`
	Trashed         *bool                  `json:"trashed"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	CurrentProgress map[string]interface{} `json:"currentProgress"`
	Pipeline        *ReleasePipeline       `json:"pipeline"`
	Stage           *ReleaseStage          `json:"stage"`
	Creator         *User                  `json:"creator"`
	ReleaseNotes    []ReleaseNote          `json:"releaseNotes"`
}

// ReleaseConnection is a paginated list of releases.
type ReleaseConnection struct {
	Nodes    []ReleaseDetail `json:"nodes"`
	PageInfo PageInfo        `json:"pageInfo"`
}

// ReleaseNote represents a Linear release note.
type ReleaseNote struct {
	ID               string           `json:"id"`
	Title            *string          `json:"title"`
	SlugID           string           `json:"slugId"`
	URL              string           `json:"url"`
	ReleaseCount     int              `json:"releaseCount"`
	GenerationStatus *string          `json:"generationStatus"`
	CreatedAt        time.Time        `json:"createdAt"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	ArchivedAt       *time.Time       `json:"archivedAt"`
	Pipeline         *ReleasePipeline `json:"pipeline"`
	FirstRelease     *ReleaseDetail   `json:"firstRelease"`
	LastRelease      *ReleaseDetail   `json:"lastRelease"`
	Releases         []ReleaseDetail  `json:"releases"`
}

// ReleaseNoteConnection is a paginated list of release notes.
type ReleaseNoteConnection struct {
	Nodes    []ReleaseNote `json:"nodes"`
	PageInfo PageInfo      `json:"pageInfo"`
}

// ReleasePipeline represents a Linear release pipeline.
type ReleasePipeline struct {
	ID                                   string     `json:"id"`
	Name                                 string     `json:"name"`
	Type                                 string     `json:"type"`
	SlugID                               string     `json:"slugId"`
	URL                                  string     `json:"url"`
	IsProduction                         bool       `json:"isProduction"`
	ApproximateReleaseCount              int        `json:"approximateReleaseCount"`
	AutoGenerateReleaseNotesOnCompletion bool       `json:"autoGenerateReleaseNotesOnCompletion"`
	RolloverIssuesOnCompletion           bool       `json:"rolloverIssuesOnCompletion"`
	IncludePathPatterns                  []string   `json:"includePathPatterns"`
	Trashed                              *bool      `json:"trashed"`
	CreatedAt                            time.Time  `json:"createdAt"`
	UpdatedAt                            time.Time  `json:"updatedAt"`
	ArchivedAt                           *time.Time `json:"archivedAt"`
}

// ReleasePipelineConnection is a paginated list of release pipelines.
type ReleasePipelineConnection struct {
	Nodes    []ReleasePipeline `json:"nodes"`
	PageInfo PageInfo          `json:"pageInfo"`
}

// ReleaseStage represents a stage within a release pipeline.
type ReleaseStage struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Type       string           `json:"type"`
	Color      string           `json:"color"`
	Position   float64          `json:"position"`
	Frozen     bool             `json:"frozen"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`
	ArchivedAt *time.Time       `json:"archivedAt"`
	Pipeline   *ReleasePipeline `json:"pipeline"`
}

// ReleaseStageConnection is a paginated list of release stages.
type ReleaseStageConnection struct {
	Nodes    []ReleaseStage `json:"nodes"`
	PageInfo PageInfo       `json:"pageInfo"`
}

const releaseFields = `
	id
	name
	version
	slugId
	url
	description
	commitSha
	startDate
	targetDate
	startedAt
	completedAt
	canceledAt
	archivedAt
	autoArchivedAt
	trashed
	createdAt
	updatedAt
	currentProgress
	pipeline {
		id
		name
		type
		slugId
		url
		isProduction
		approximateReleaseCount
		createdAt
		updatedAt
		archivedAt
	}
	stage {
		id
		name
		type
		color
		position
		frozen
		createdAt
		updatedAt
		archivedAt
	}
	creator {
		id
		name
		email
	}
`

// GetReleases returns a paginated list of releases.
func (c *Client) GetReleases(ctx context.Context, filter map[string]interface{}, first int, after string, includeArchived bool) (*ReleaseConnection, error) {
	query := `
		query Releases($filter: ReleaseFilter, $first: Int, $after: String, $includeArchived: Boolean) {
			releases(filter: $filter, first: $first, after: $after, includeArchived: $includeArchived) {
				nodes {` + releaseFields + `}
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
	if includeArchived {
		variables["includeArchived"] = includeArchived
	}

	var response struct {
		Releases ReleaseConnection `json:"releases"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Releases, nil
}

// GetRelease returns a single release by ID or slug.
func (c *Client) GetRelease(ctx context.Context, id string) (*ReleaseDetail, error) {
	query := `
		query Release($id: String!) {
			release(id: $id) {` + releaseFields + `
				releaseNotes {
					id
					title
					createdAt
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Release ReleaseDetail `json:"release"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.Release, nil
}

// GetReleaseNotes returns a paginated list of release notes.
func (c *Client) GetReleaseNotes(ctx context.Context, filter map[string]interface{}, first int, after string) (*ReleaseNoteConnection, error) {
	query := `
		query ReleaseNotes($filter: ReleaseNoteFilter, $first: Int, $after: String) {
			releaseNotes(filter: $filter, first: $first, after: $after) {
				nodes {
					id
					title
					slugId
					url
					releaseCount
					generationStatus
					createdAt
					updatedAt
					archivedAt
					pipeline {
						id
						name
						type
					}
					firstRelease {
						id
						name
						version
					}
					lastRelease {
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

	var response struct {
		ReleaseNotes ReleaseNoteConnection `json:"releaseNotes"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ReleaseNotes, nil
}

// GetReleasePipelines returns a paginated list of release pipelines.
func (c *Client) GetReleasePipelines(ctx context.Context, filter map[string]interface{}, first int, after string) (*ReleasePipelineConnection, error) {
	query := `
		query ReleasePipelines($filter: ReleasePipelineFilter, $first: Int, $after: String) {
			releasePipelines(filter: $filter, first: $first, after: $after) {
				nodes {
					id
					name
					type
					slugId
					url
					isProduction
					approximateReleaseCount
					autoGenerateReleaseNotesOnCompletion
					rolloverIssuesOnCompletion
					includePathPatterns
					trashed
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
	if filter != nil {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		ReleasePipelines ReleasePipelineConnection `json:"releasePipelines"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ReleasePipelines, nil
}

// GetReleaseStages returns a paginated list of release stages.
func (c *Client) GetReleaseStages(ctx context.Context, filter map[string]interface{}, first int, after string) (*ReleaseStageConnection, error) {
	query := `
		query ReleaseStages($filter: ReleaseStageFilter, $first: Int, $after: String) {
			releaseStages(filter: $filter, first: $first, after: $after) {
				nodes {
					id
					name
					type
					color
					position
					frozen
					createdAt
					updatedAt
					archivedAt
					pipeline {
						id
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
		"first": first,
	}
	if filter != nil {
		variables["filter"] = filter
	}
	if after != "" {
		variables["after"] = after
	}

	var response struct {
		ReleaseStages ReleaseStageConnection `json:"releaseStages"`
	}

	err := c.Execute(ctx, query, variables, &response)
	if err != nil {
		return nil, err
	}

	return &response.ReleaseStages, nil
}
