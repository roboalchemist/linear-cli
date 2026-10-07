# Projects

Create Projects for features or large units of work.

![Project overview in Linear](https://webassets.linear.app/images/ornj730p/production/5622bb57d625b2dbbe27639fe6f685e1c43589bd-2105x1332.png?q=95&auto=format&dpr=2)

## Overview

A project is a grouping of issues that have a clear outcome or planned completion date, such as a new feature launch. They can be shared across multiple teams and come with their own unique features, progress graph, and notification options.

Projects have a project ID, formatted as `P-TEAM-123`, where “P” is for project, “TEAM” is the team identifier, and the number is the sequence in which the project was created.

![Watch on YouTube](https://www.youtube.com/watch?v=JLpDL7x50hA)

## Create a project

Create new projects directly from your workspace or team's project view by clicking the `+` button. This opens a prompt for project properties. The only required fields are the project name and status.

We recommend assigning a project lead and updating the project icon for better visibility and organization, as well as filling out the overview to describe what’s known about the project. You can also select a [project status](https://linear.app/docs/project-status) and apply [project labels](https://linear.app/docs/project-labels) when creating a project.

To edit these properties later, right-click a project in an [initiative](https://linear.app/docs/initiatives) or [custom project view](https://linear.app/docs/custom-views#create-views) and select **Edit project**, or click the properties on the project overview page.

## Delete a project

To delete a project that you have in focus, click the three dots next to the project name beside the Overview and Issues tabs and select **Delete**.

From a view containing multiple projects, you can also right click on any project and select **Delete** from the bottom of the menu.

Once you've confirmed this action, your deleted projects can be found in the team's archive under the "Recently deleted projects" tab for 30 days before they are permanently deleted.

## View your projects

Each team has a _Projects_ page which organizes the team's projects into a list, board, or timeline. There is also a _Projects_ page at the workspace level so you can view all projects within your workspace in one location.

You can organize project views by properties such as status, lead, target date, and activity. For example, group projects by target date to see what is planned for each quarter, or use activity to see which projects have recent momentum relative to other active projects in your workspace.

To explore a specific project, click on it to access the project overview. Here, you'll find a brief project summary, project properties, any associated documents and links, a detailed project description, and a list of project milestones.

## Add issues to a project

Add existing issues to projects or move issues between projects by selecting the issues and using the shortcut `Shift` `P`. To create new issues in a project, use `C` from the project view or add the project property manually if creating the issue from another page.

> [!NOTE]
> **How we work**: The project lead is in charge of writing the spec and general execution. Other team members collaborate on the brief, split up areas of work, and then write their own issues.

Issues can only be associated with one project at a time. A workaround would be to create sub-issues for the task, then assign each sub-issue to a different project.

## Project documents

To create a document in a Linear project, go to the `Resources` section. Press the  `+` button and select Create Document. Documents use the same Markdown as issues and have similar editing capabilities including tagging others and linking to issues.

### External links

To add an external link to a Linear project, go to the `Resources` section. Press the `+` button, then add the link and press enter. Add a link label to give the link a title, otherwise we'll show the link source (e.g. docs.google.com). To edit an existing link label, hover over the relevant link, click on the **`⋮ `** icon and choose the edit option here.

### Milestones

You can view a list of [project milestones](https://linear.app/docs/project-milestones) in the project overview. You can reorder the milestones from here, set their target dates and edit their names. You can also include a detailed description below each milestone.

## Attach issue views to projects

Next to the "Issues" tab in each project, you'll find the `new view` icon. This feature enables the creation of custom views of the project's issues. By clicking this icon, you can filter a subset of the issues and save this filtered perspective under a specific view name.

Once a view is created within the project, it can be repositioned by clicking and dragging it to the left or right to adjust the order in which these views appear. These views are accessible to anyone who has access to the project and are displayed as tabs, similar to the `Overview` and `Issues` tabs, allowing users to easily share and switch between different perspectives within the project.

Additionally, right-clicking on a view opens options to copy the view link, add it to your favorites, edit its details, or delete it. 

> [!NOTE]
> **Example views we've attached to projects:** A view we use often is a view that filters for all issues assigned to the `current user`. This ensures everyone viewing a project can quickly view all their issues. We also use a `bug` view that filters for all issues with the "bug" label. A "standup" view filtering for `In Progress` helps us have more efficient standups.

## Attach project views to workspace projects page 

Similar to attaching issue views to projects, you can also attach project views within the workspace projects page. This allows you to create custom views of projects in your workspace. These views can be repositioned just like attached issue views within projects and are accessible to all members of your workspace.

> [!NOTE]
> **Example project views we've attached**: In our workspace, we have attached a project view that filters for all projects with an “In Progress” status. Projects are automatically added to this view based on the selected filters, eliminating the need for manual curation. This makes it easy to reference and see all in-progress projects at a glance. We also have an attached project view that filters for all projects involving a specific team, allowing us to quickly review that team’s activity.

## Project details sidebar

Press the sidebar icon or use `Cmd/Ctrl` `I` to toggle the project details sidebar on or off. You'll be able to update any properties from here as well as add resources such as external files, create project documents, and view the project progress [graph](https://linear.app/docs/project-graph). 

## Multi-team projects

Projects can be shared across multiple teams. Add more teams when creating a project or from the project details page by clicking on the team name, then marking off more teams in the dropdown. Once you add more than one team, the project view will create tabs so you can toggle between seeing all issues or issues on specific teams. 

One of a project's teams is its lead team, which clarifies ownership and determines its available project statuses. The other teams contribute to the project. On Business and Enterprise plans, a project can use project labels from its teams alongside workspace-level labels.

![Adding team to a project in Linear](https://webassets.linear.app/images/ornj730p/production/97431e2be69c1959a7759239cc709a965473c386-1573x1035.png?q=95&auto=format&dpr=2)

## Timeframes

Rarely will a project's precise end date be known in its early stages. Select start and target dates that match your level of certainty. Options are available to choose a year, half-year, quarter, month or precise day. Choose start and target dates using these timeframes by selecting the appropriate fields in the project's right sidebar.

You can also use these timeframes in project views in [Initiatives](https://linear.app/docs/initiatives) when grouping by start or target date, or in start/target date filters as well as created at under the "custom date or timeframe…" option.

## FAQ

<details>
<summary>Can I make a project for ongoing work?</summary>
Technically yes, you can create a project without a completion date -- but we encourage projects to have a clear start and end. We're looking into better ways to manage this type of work but in the meantime there are a few workarounds:

* Use projects without target dates (you may want to create an Initiative to group these projects)
* Break up the work into smaller time-based projects (e.g. bugweek, Q1 infrastructure)
* In cases where a project should remain unclosed for some time but not forever, consider creating a custom Maintenance project status.
</details>

<details>
<summary>Can I add multiple project leads?</summary>
We have a single lead field to keep ownership of the project clear. If more people are involved with the project, consider adding them as members, however members have to opt-in to receive notifications.
</details>

<details>
<summary>How is it decided what projects are listed in the project picker menu?</summary>
It will default to show you the following options when adding a project to an issue:   
  
Projects you lead, are a member of, created by user recently, has overlapping teams with user, active, recently created, then cancelled and completed.
</details>

<details>
<summary>Where is the roadmap timeline view?</summary>
Roadmaps are renamed to Initiatives and can be enabled in your workspace settings.   
Timeline views are available on any view with projects in your team projects view, inside an initiative or the workspace projects view. You can choose a timeline view under "Display".
</details>