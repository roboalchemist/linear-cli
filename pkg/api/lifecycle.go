package api

import (
	"context"
	"encoding/json"
	"fmt"
)

// mutationOK runs a mutation that returns a payload with a `success` boolean and
// returns an error when the mutation does not report success.
func (c *Client) mutationOK(ctx context.Context, query string, variables map[string]interface{}, field string) error {
	var resp map[string]json.RawMessage
	if err := c.Execute(ctx, query, variables, &resp); err != nil {
		return err
	}
	raw, ok := resp[field]
	if !ok {
		return fmt.Errorf("unexpected response: missing %q", field)
	}
	var payload struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return err
	}
	if !payload.Success {
		return fmt.Errorf("%s reported failure", field)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Issues
// ---------------------------------------------------------------------------

// UnarchiveIssue restores an archived issue.
func (c *Client) UnarchiveIssue(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { issueUnarchive(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "issueUnarchive")
}

// DeleteIssue deletes an issue. When permanent is false the issue is trashed.
func (c *Client) DeleteIssue(ctx context.Context, id string, permanent bool) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $permanent: Boolean) { issueDelete(id: $id, permanentlyDelete: $permanent) { success } }`,
		map[string]interface{}{"id": id, "permanent": permanent}, "issueDelete")
}

// AddIssueLabel attaches a label to an issue.
func (c *Client) AddIssueLabel(ctx context.Context, issueID, labelID string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $labelId: String!) { issueAddLabel(id: $id, labelId: $labelId) { success } }`,
		map[string]interface{}{"id": issueID, "labelId": labelID}, "issueAddLabel")
}

// RemoveIssueLabel detaches a label from an issue.
func (c *Client) RemoveIssueLabel(ctx context.Context, issueID, labelID string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $labelId: String!) { issueRemoveLabel(id: $id, labelId: $labelId) { success } }`,
		map[string]interface{}{"id": issueID, "labelId": labelID}, "issueRemoveLabel")
}

// SubscribeIssue subscribes a user (default: the viewer) to an issue.
func (c *Client) SubscribeIssue(ctx context.Context, issueID, email string) error {
	vars := map[string]interface{}{"id": issueID}
	q := `mutation($id: String!) { issueSubscribe(id: $id) { success } }`
	if email != "" {
		q = `mutation($id: String!, $email: String!) { issueSubscribe(id: $id, userEmail: $email) { success } }`
		vars["email"] = email
	}
	return c.mutationOK(ctx, q, vars, "issueSubscribe")
}

// UnsubscribeIssue unsubscribes a user (default: the viewer) from an issue.
func (c *Client) UnsubscribeIssue(ctx context.Context, issueID, email string) error {
	vars := map[string]interface{}{"id": issueID}
	q := `mutation($id: String!) { issueUnsubscribe(id: $id) { success } }`
	if email != "" {
		q = `mutation($id: String!, $email: String!) { issueUnsubscribe(id: $id, userEmail: $email) { success } }`
		vars["email"] = email
	}
	return c.mutationOK(ctx, q, vars, "issueUnsubscribe")
}

// ShareIssue shares an issue with a user.
func (c *Client) ShareIssue(ctx context.Context, issueID, userID string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $userId: String!) { issueShare(id: $id, userId: $userId) { success } }`,
		map[string]interface{}{"id": issueID, "userId": userID}, "issueShare")
}

// UnshareIssue removes a user's access to a shared issue.
func (c *Client) UnshareIssue(ctx context.Context, issueID, userID string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $userId: String!) { issueUnshare(id: $id, userId: $userId) { success } }`,
		map[string]interface{}{"id": issueID, "userId": userID}, "issueUnshare")
}

// AddIssueReminder sets a reminder for the viewer on an issue.
func (c *Client) AddIssueReminder(ctx context.Context, issueID, reminderAt string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $reminderAt: DateTime!) { issueReminder(id: $id, reminderAt: $reminderAt) { success } }`,
		map[string]interface{}{"id": issueID, "reminderAt": reminderAt}, "issueReminder")
}

// RemoveIssueReminder clears the viewer's reminder on an issue.
func (c *Client) RemoveIssueReminder(ctx context.Context, issueID string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { issueRemoveReminder(id: $id) { success } }`,
		map[string]interface{}{"id": issueID}, "issueRemoveReminder")
}

// BatchCreateIssues creates many issues at once.
func (c *Client) BatchCreateIssues(ctx context.Context, issues []map[string]interface{}) error {
	return c.mutationOK(ctx,
		`mutation($input: IssueBatchCreateInput!) { issueBatchCreate(input: $input) { success } }`,
		map[string]interface{}{"input": map[string]interface{}{"issues": issues}}, "issueBatchCreate")
}

// BatchUpdateIssues applies the same update to many issues at once.
func (c *Client) BatchUpdateIssues(ctx context.Context, ids []string, input map[string]interface{}) error {
	return c.mutationOK(ctx,
		`mutation($ids: [UUID!]!, $input: IssueUpdateInput!) { issueBatchUpdate(ids: $ids, input: $input) { success } }`,
		map[string]interface{}{"ids": ids, "input": input}, "issueBatchUpdate")
}

// ---------------------------------------------------------------------------
// Comments
// ---------------------------------------------------------------------------

// ResolveComment resolves a comment thread.
func (c *Client) ResolveComment(ctx context.Context, id, resolvingCommentID string) error {
	vars := map[string]interface{}{"id": id}
	q := `mutation($id: String!) { commentResolve(id: $id) { success } }`
	if resolvingCommentID != "" {
		q = `mutation($id: String!, $cid: String!) { commentResolve(id: $id, resolvingCommentId: $cid) { success } }`
		vars["cid"] = resolvingCommentID
	}
	return c.mutationOK(ctx, q, vars, "commentResolve")
}

// UnresolveComment reopens a resolved comment thread.
func (c *Client) UnresolveComment(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { commentUnresolve(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "commentUnresolve")
}

// ---------------------------------------------------------------------------
// Projects / project updates
// ---------------------------------------------------------------------------

// UnarchiveProject restores an archived project.
func (c *Client) UnarchiveProject(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { projectUnarchive(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "projectUnarchive")
}

// AddProjectLabel attaches a label to a project.
func (c *Client) AddProjectLabel(ctx context.Context, projectID, labelID string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $labelId: String!) { projectAddLabel(id: $id, labelId: $labelId) { success } }`,
		map[string]interface{}{"id": projectID, "labelId": labelID}, "projectAddLabel")
}

// RemoveProjectLabel detaches a label from a project.
func (c *Client) RemoveProjectLabel(ctx context.Context, projectID, labelID string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $labelId: String!) { projectRemoveLabel(id: $id, labelId: $labelId) { success } }`,
		map[string]interface{}{"id": projectID, "labelId": labelID}, "projectRemoveLabel")
}

// ReassignProjectStatus moves projects from one status to another.
func (c *Client) ReassignProjectStatus(ctx context.Context, newStatusID, originalStatusID string) error {
	return c.mutationOK(ctx,
		`mutation($newId: String!, $origId: String!) { projectReassignStatus(newProjectStatusId: $newId, originalProjectStatusId: $origId) { success } }`,
		map[string]interface{}{"newId": newStatusID, "origId": originalStatusID}, "projectReassignStatus")
}

// UnarchiveProjectUpdate restores an archived project status update.
func (c *Client) UnarchiveProjectUpdate(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { projectUpdateUnarchive(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "projectUpdateUnarchive")
}

// DeleteProjectUpdate permanently deletes a project status update.
func (c *Client) DeleteProjectUpdate(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { projectUpdateDelete(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "projectUpdateDelete")
}

// ---------------------------------------------------------------------------
// Initiatives
// ---------------------------------------------------------------------------

// ArchiveInitiative archives an initiative.
func (c *Client) ArchiveInitiative(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { initiativeArchive(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "initiativeArchive")
}

// UnarchiveInitiative restores an archived initiative.
func (c *Client) UnarchiveInitiative(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { initiativeUnarchive(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "initiativeUnarchive")
}

// AddInitiativeLabel attaches a label to an initiative.
func (c *Client) AddInitiativeLabel(ctx context.Context, initiativeID, labelID string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $labelId: String!) { initiativeAddLabel(id: $id, labelId: $labelId) { success } }`,
		map[string]interface{}{"id": initiativeID, "labelId": labelID}, "initiativeAddLabel")
}

// RemoveInitiativeLabel detaches a label from an initiative.
func (c *Client) RemoveInitiativeLabel(ctx context.Context, initiativeID, labelID string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $labelId: String!) { initiativeRemoveLabel(id: $id, labelId: $labelId) { success } }`,
		map[string]interface{}{"id": initiativeID, "labelId": labelID}, "initiativeRemoveLabel")
}

// ---------------------------------------------------------------------------
// Documents
// ---------------------------------------------------------------------------

// UnarchiveDocument restores an archived document.
func (c *Client) UnarchiveDocument(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { documentUnarchive(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "documentUnarchive")
}

// ---------------------------------------------------------------------------
// Cycles
// ---------------------------------------------------------------------------

// StartUpcomingCycleToday starts the given upcoming cycle as of today.
func (c *Client) StartUpcomingCycleToday(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { cycleStartUpcomingCycleToday(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "cycleStartUpcomingCycleToday")
}

// ShiftAllCycles shifts all cycles after the given one by a number of days.
func (c *Client) ShiftAllCycles(ctx context.Context, id string, days float64) error {
	return c.mutationOK(ctx,
		`mutation($input: CycleShiftAllInput!) { cycleShiftAll(input: $input) { success } }`,
		map[string]interface{}{"input": map[string]interface{}{"id": id, "daysToShift": days}}, "cycleShiftAll")
}

// ---------------------------------------------------------------------------
// Workflow states
// ---------------------------------------------------------------------------

// CreateWorkflowState creates a workflow state for a team.
func (c *Client) CreateWorkflowState(ctx context.Context, input map[string]interface{}) (string, error) {
	var resp struct {
		WorkflowStateCreate struct {
			Success       bool `json:"success"`
			WorkflowState struct {
				ID string `json:"id"`
			} `json:"workflowState"`
		} `json:"workflowStateCreate"`
	}
	err := c.Execute(ctx,
		`mutation($input: WorkflowStateCreateInput!) { workflowStateCreate(input: $input) { success workflowState { id name } } }`,
		map[string]interface{}{"input": input}, &resp)
	if err != nil {
		return "", err
	}
	if !resp.WorkflowStateCreate.Success {
		return "", fmt.Errorf("workflowStateCreate reported failure")
	}
	return resp.WorkflowStateCreate.WorkflowState.ID, nil
}

// UpdateWorkflowState updates a workflow state.
func (c *Client) UpdateWorkflowState(ctx context.Context, id string, input map[string]interface{}) error {
	return c.mutationOK(ctx,
		`mutation($id: String!, $input: WorkflowStateUpdateInput!) { workflowStateUpdate(id: $id, input: $input) { success } }`,
		map[string]interface{}{"id": id, "input": input}, "workflowStateUpdate")
}

// ArchiveWorkflowState archives a workflow state.
func (c *Client) ArchiveWorkflowState(ctx context.Context, id string) error {
	return c.mutationOK(ctx,
		`mutation($id: String!) { workflowStateArchive(id: $id) { success } }`,
		map[string]interface{}{"id": id}, "workflowStateArchive")
}

// ---------------------------------------------------------------------------
// Notifications (bulk)
// ---------------------------------------------------------------------------

func (c *Client) notificationBulk(ctx context.Context, mutation, arg, field string, input map[string]interface{}) error {
	q := fmt.Sprintf(
		`mutation($input: NotificationEntityInput!) { %s(input: $input) { success } }`,
		mutation,
	)
	if arg != "" {
		q = fmt.Sprintf(
			`mutation($input: NotificationEntityInput!) { %s(input: $input, %s) { success } }`,
			mutation, arg,
		)
	}
	return c.mutationOK(ctx, q, map[string]interface{}{"input": input}, field)
}

// ArchiveAllNotifications archives all notifications for the given entity (or everything).
func (c *Client) ArchiveAllNotifications(ctx context.Context, input map[string]interface{}) error {
	return c.notificationBulk(ctx, "notificationArchiveAll", "", "notificationArchiveAll", input)
}

// MarkAllNotificationsUnread marks all notifications for the given entity as unread.
func (c *Client) MarkAllNotificationsUnread(ctx context.Context, input map[string]interface{}) error {
	return c.notificationBulk(ctx, "notificationMarkUnreadAll", "", "notificationMarkUnreadAll", input)
}

// SnoozeAllNotifications snoozes all notifications for the given entity until the given time.
func (c *Client) SnoozeAllNotifications(ctx context.Context, input map[string]interface{}, snoozedUntilAt string) error {
	return c.mutationOK(ctx,
		`mutation($input: NotificationEntityInput!, $at: DateTime!) { notificationSnoozeAll(input: $input, snoozedUntilAt: $at) { success } }`,
		map[string]interface{}{"input": input, "at": snoozedUntilAt}, "notificationSnoozeAll")
}

// UnsnoozeAllNotifications clears snoozes for all notifications for the given entity.
func (c *Client) UnsnoozeAllNotifications(ctx context.Context, input map[string]interface{}, unsnoozedAt string) error {
	return c.mutationOK(ctx,
		`mutation($input: NotificationEntityInput!, $at: DateTime!) { notificationUnsnoozeAll(input: $input, unsnoozedAt: $at) { success } }`,
		map[string]interface{}{"input": input, "at": unsnoozedAt}, "notificationUnsnoozeAll")
}
