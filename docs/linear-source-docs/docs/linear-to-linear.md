# Linear

## Prerequisites

To import data from one Linear workspace to another, first ensure that you have an admin account in both workspaces using the same email login credential.

> [!NOTE]
> If there isn't already an admin using the same email credential in both source and destination workspaces, invite an admin from the destination workspace to the source workspace (or vice versa) and have them run the import. For best results, avoid creating multiple accounts on the same workspace under different emails.

## Running the import

1. Navigate to [Settings > Administration > Import/Export](https://linear.app/settings/import-export). You must be an admin in Linear to access this page.
2. Click on the "Linear to Linear import" button.
3. Select the workspace you would like to import. We will only show workspaces under the same email address where you are an admin member.
4. Select the teams you would like to import from this workspace.
5. Choose how to map members from the workspace you are importing from to your current workspace.
6. Review a summary of the import.
7. Click "Start import".

## Data included in the import

Linear (old) | Linear (new)
--- | ---
Title | Title
Description | Description
Estimate | Estimate
Labels | Labels
Due date | Due date
Comments | Comments
Workflow state | Workflow state
Sub-issues | Sub-issues
Relationships | Relationships
Projects | Projects
Initiatives | Initiatives
Team templates | Team templates
Dashboards | Dashboards
Documents | Documents

## Data excluded in the import

This import will not carry over all data associated with your source workspace. The data below will not transfer the following:

* Saved display preferences for custom, project, and initiative views
* Favorites / reminders / drafts / inbox notifications
* Integrations* / webhooks / OAuth clients / API keys
* Billing/plans**
* Workspace URL (you can update this manually after import in your [Settings > Workspace](https://linear.app/settings/workspace))
* Personal and Workspace settings
* Roles***

* Integrations will need to be installed on the new workspace and will not resume notifications or automations for imported issues in most cases.

** Contact support@linear.app for assistance with billing transfers.

*** Current admins will be imported as members. The person who initiated the import will initially be the only admin. Guests will carry over with the same permissions.