# CLI Importer

Use the CLI Importer when you’re moving from a tool without a dedicated Linear import assistant.

> [!NOTE]
> While CLI imports can also be used to customize imports for services with an importer tool, we don't recommend this as the CLI importer will not import data supported by the import assistant like comments or projects.

## Preparing your CSV file for import

Go to Linear [Settings → Administration → Import/Export → Export CSV](https://linear.app/settings/import-export) to download a template. Keep its column headers, replace the sample rows with your data, and save the file.

The essential fields are:

* `Title` - Issue title
* `Description` - Issue description
* `Priority` - Issue priority
* `Status` - Issue state
* `Assignee` - Issue assignee (user’s full name)
* `Created` - Issue created date
* `Completed` - Issue completed timestamp. This will not appear in issue activity log, but will be stored on the issue
* `Labels` - Separate each label by commas for multiple labels
* `Estimate` - Please note that these values only appear on the issue after enabling estimates on the imported issue’s team.

## Running the import

You’ll need to be a workspace admin and a CSV export from your source tool.

1. Navigate to [Settings > Workspace > Import/Export](https://linear.app/settings/import-export)
2. Under “CLI import”, click **Open** to navigate to our [Command Line interface (CLI) tool](https://github.com/linear/linear/tree/master/packages/import).
3. Follow instructions to run the CLI Importer
4. Choose the "Linear" CSV option
5. Follow prompts to select your file and destination team

Re-imports to the same team skip issues that were already imported, so delete a previous import first if you need a clean retry.