# Import Guide

Learn best practices for importing to Linear. Imports can only be started by workspace admins.

![Image of migration assistants](https://webassets.linear.app/images/ornj730p/production/d2f61b4eeb982988a066a5416eaa81bb547178b3-2880x1769.png?q=95&auto=format&dpr=2)

## Linear Concepts

A Linear [workspace](https://linear.app/docs/workspaces) is our top level concept, which contains one or many [teams](https://linear.app/docs/teams). We recommend that each company using Linear uses only one workspace.

[Issues](https://linear.app/docs/creating-issues) and [projects](https://linear.app/docs/projects) are the core entities used to manage work in Linear. Each issue belongs to a single team, while projects can belong to one or many teams. Other concepts in Linear can be scoped to a team or a workspace like [Views](https://linear.app/docs/custom-views) (filter based groups of issues or projects), [Initiatives](https://linear.app/docs/initiatives) (hand-picked groups of projects) and more.

When you import from another service to Linear, we'll attempt to match data from the source tool with the closest concept in Linear. If a concept from the source tool does not translate well to Linear, it may not be imported; please see details on your individual service for further details.

## Pre-import best practices

 Consider the below before starting your first import for a smoother experience.

### Which data is worth importing?

 If there is data that is no longer relevant to your organization's day to day work, consider whether it needs to be in Linear at all (perhaps a CSV of exported data from long-resolved issues is sufficient). Some organizations choose to use Linear as a "clean break" from their legacy tool and import only where absolutely necessary so they can start fresh with minimal clutter. Others prefer to maintain as full a historical record as possible in one place.

If you're importing when evaluating Linear instead of transitioning all at once, you may wish to run a pilot by importing just a few teams. We also have a resource for switching tools in our [switch instruction manual. ](https://linear.app/switch)

### Choose an import method

We offer two main methods of importing to Linear; our dedicated import assistants in-product and a CLI import tool. We recommend using a dedicated importer whenever possible as they retain much more data from the original source, are easier to use, and provide the option to delete an import in bulk if desired.

Jump to our importers:

* [Jira](https://linear.app/docs/jira-to-linear)
* [GitHub Issues](https://linear.app/docs/github-to-linear)
* [Asana](https://linear.app/docs/asana-to-linear)
* [Shortcut](https://linear.app/docs/shortcut-to-linear)
* [Linear](https://linear.app/docs/linear-to-linear)
* [CLI](https://linear.app/docs/cli-importer)

### Understand the import process

You'll first need to have an Admin role in Linear in order to access our import tools. You may also need high permissions in the tool you're importing from in order to access the data. In general, the import assistants follow this path:

Step | Detail
--- | ---
Setup  | Provide an access token, or sign in to your source tool from Linear. Choose an existing team to import to, or create a new team. You can import issues into a top-level team or a sub-team with its own workflow statuses. Imports aren't available for sub-teams that inherit workflow statuses from their parent team.
Review | We'll display the issues, projects, labels, users and other data fetched from the information provided. 
Choose what to import | Choose which issues to import, including active-only or broader sets such as stale, completed, or all issues depending on the source. The archive can be accessed through the overflow menu on each team in your sidebar.
Map users | For each user fetched, choose not to import them, to create a new user from an email address, or to map them to an existing user in your Linear workspace.
Confirm | At this point, all data available to import will be summarized. If it looks as you expect, hit Finish to start the import. 

Some importer behavior varies by source. For details not covered in this article, check the guide for your source tool in the section below, especially for Jira and GitHub.

### Troubleshooting

If something did not import as you expected, please check the section for your specific service to confirm whether we support importing that property. If the property is supported but didn't import as documented, please let us know at [support@linear.app](mailto:support@linear.app).

If you need to delete an import in order to re-import once more, you can do so through Import/Export settings, on the overflow menu on a specific import.  Imports can only be deleted for a limited time after creation, so if the delete option is no longer available, the import may be outside that window.

Reimporting from the same external source to the same Linear team without deleting the initial import first will skip any already-imported issues.