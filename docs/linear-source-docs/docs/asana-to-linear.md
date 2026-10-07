# Asana

Import Asana into Linear

## Prerequisites

We currently only support in app importing from Asana organizations. Make sure your Asana workspace is an organization before beginning an import, or [convert your Asana workspace to an organization](https://asana.com/guide/help/workspaces/create#gl-convert-workspace) if it is not. If you do not want to convert your Asana workspace to an organization, use the CLI importer.

## Running the import

1. Navigate to [Settings > Administration > Import/Export](https://linear.app/settings/import-export).
2. Click on the button for Asana.
3. Enter your Asana personal access token. You can obtain your personal access token by navigating to your Asana [developer console](https://app.asana.com/0/my-apps).
4. Enter your Asana team name.
5. Select which Linear team to import issues into.
6. Click **Next** to start importing.

## Data included in the import

Asana | Linear
--- | ---
Priority | Priority
Notes | Issue description (converted to Markdown)
Attached files | Added to description. All files except images and videos are converted to links.
Tags | Labels (Team-level)
Assignee | Assignee
Projects | Projects
Comments | Comments (Markdown respected)
Status | Imported status will be "Backlog" or "Done", only.
Sub-issue | Sub-issue
Blocked/blocking | Blocked/blocking