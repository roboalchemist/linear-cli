# Issues

Issues are the core unit of work and can represent anything a team needs to track, like bugs, feature work, follow-up tasks, or internal requests.

![issue creation dialogue box](https://webassets.linear.app/images/ornj730p/production/25ae979891503f4780b6e837aaa71b5027c67281-1224x534.png?q=95&auto=format&dpr=2)

## Overview

Issues must belong to a single team and are given a unique issue ID number that is assigned in consecutive order of when it was created. They are required to have a title and a status.

## Creating issues

* Use the keyboard shortcut `C` to open up an issue creation modal
* Click the **Create new issue** icon in the upper left of the app
* To create an issue from a template, use `Option/Alt + C`.
* Enter [https://linear.new](https://linear.new) into your browser URL bar to create a new issue. It will redirect you to the new issue creation page as long as you are logged into your Linear account.
* Highlight text and click the **Create issue from selection** icon in the formatting bar

> [!NOTE]
> Changes made to an issue's properties in the first 3 minutes are considered part of the issue creation process, and won't be added to the activity log as changes to the issue.

## Create an issue via email

To allow issues to be created without requiring an account in Linear, you can set up email intake by creating a custom email address for a team or for a team’s template.

When someone sends an email to one of the intake email addresses, they _will not receive a confirmation or be notified_ when the issue is updated or resolved.

The submitted Linear issue will be included as an attachment on the Linear issue. Attachment files will be synced over, though email attachments are limited to 25 MB.

> [!NOTE]
> Refer to the [Linear Asks](https://linear.app/docs/linear-asks) feature for the ability to:
> 
> * Customize the email address you want requesters to send to
> * Reply to customer emails from Linear
> * Customize response emails

### Team intake email

Navigate to **Team settings > General > Create issues by email** and enable the toggle.

### Templated intake email

To create an email address for a team template:

1. Navigate to **Team Settings > Templates**
2. Click the three dots on the right of an issue template
3. Select **Configure email address**
4. Enable the toggle in the pop-up

When a team template is used, the issue's title and description will be overwritten by the email contents, but the properties of the template will be applied to the new issue. Replies sent on the original email to the forwarding address will not create additional issues. 

## Create recurring issues

To automate your repeated tasks on a cadence of your choosing, you can choose to make it recurring:

1. Open the issue composer
2. Apply a template (optional)
3. Click the `…` menu
4. Select _"Make recurring…"_

Recurring issues can be found from your **Team settings > Recurring issues**.

![Creating a recurring issue from the issue composer](https://webassets.linear.app/images/ornj730p/production/7d4cb7c61fa50aad11ca1bca4e9844ec6bd98489-1091x944.png?q=95&auto=format&dpr=2)

Once you create a recurring issue, future issues in the cadence are expected to be created once the due date passes (00:01 the following day in your team’s timezone.)

> [!NOTE]
> Changes to a template in future will not affect recurring issues that were created from this template. You will need to edit the recurring issue directly or recreate it from your updated template.

### Convert an existing issue into a recurring issue

To convert any issue into a recurring issue, open the issue and in the `…`  menu choose` Convert into > Recurring issue…,  
,`You can also use the `Cmd/Ctrl` + `K `menu by typing _"Convert into recurring issue"_.  
You can then choose your first due date, and the cadence at which it repeats.

### Create recurring issues from team settings

1. Navigate to **Team settings > Recurring issues.**
2. Click the `+` icon and set your chosen schedule of recurrence

## Create a new issue URL

The following links trigger the creation of a new Linear issue in any browser and you can add query parameters after any of them to pre-fill issue fields.

* [http://linear.app/team/<team ID>/new](https://linear.app/team/%3Cteam%20ID%3E/new)
* [http://linear.app/new](https://linear.app/new)
* [http://linear.new](https://linear.new/)

To pre-fill issue fields and/or properties:

1. Add a `?` at the end of the link
2. Include the field or issue property you want to pre-set
3. Add `=`
4. Add the parameter you are setting
5. Use `&` between each field or issue property when creating a string of pre-settings

### Apply pre-set properties

We support the following query parameters:

Property | How to apply | Example
--- | --- | ---
Title and description | Use `+` to indicate empty space in the keyword, or fully url encode content if more complex as description can be a markdown document | https://linear.new?title=My+issue+title&description=This+is+my+issue+description
Status | Can be set by `UUID` or name of the workflow status | https://linear.new?status=Todo
Team | Can be set by `UUID` or the slug/identifier of the team | [https://linear.new?team=ENG](https://linear.new?team=ENG)
Priority | Can be set by Urgent, High, Medium or Low | https://linear.new?priority=Urgent
Assignee | Can be set by `UUID`, display name/name of the user, or `assignee=me` to assign the creator | https://linear.new?assignee=john

https://linear.new?assignee=John+Smith

https://linear.new?assignee=me
Estimate | Can be set by their point number

T-shirt sizes have the following point values: No priority (0), XS (1), S(2), M (3), L (5), XL (8), XXL (13), XXXL (21) | https://linear.app/team/LIN/new?estimate=2
Cycle | Can be set by `UUID`, cycle number of a name of a cycle | https://linear.app/team/MOB?cycle=36

https://linear.new/team/EU/new?cycle=focus+on+bugs
Label | Use a comma between each label you want to apply | https://linear.new/team/LIN/new?label=bug

https://linear.new?labels=bug,android,comments
Project | Can be set by `UUID` or the name of the project | https://linear.new/team/LIN/new?project=Project+page+improvements
Milestone | Can be set by `UUID` or the name of the project milestone

Project milestone can be read only if `project` is also passed in the URL | https://linear.app/team/LIN/new?project=Project+page+improvements&projectMilestone=Beta
Links | URL encoded comma delimited urls with optional title, in format `url|title`. These will be attached to the issue as link attachments | https://linear.new/team/LIN/new?links=https%3A%2F%2Flinear.app%2Chttp%3A%2F%2Fgoogle.com%7CGoogle
Template | Can be set by `UUID` or the name of the template | [https://linear.app/team/LIN/new?template=30cd8534-6b24-40bb-8ee7-bb0d83df7d4f](https://linear.app/team/LIN/new?template=30cd8534-6b24-40bb-8ee7-bb0d83df7d4f)

### Create a URL from a template

1. Go to **Settings > Team > Templates**
2. Click the three dots to the right of the template
3. Select **Copy URL to create issue from template**

### Copy a URL from an issue

While on an issue's view, or an issue is highlighted or selected, open command bar by typing `Cmd/Ctrl` + `K `and select **Copy pre-filled create issue URL to clipboard**.

## Drafts

When writing an issue and navigating away to another part of the app, Linear will hide the issue modal and keep a temporary draft. The next time you go to create an issue, the editor will re-open with the previous content draft. This type of draft is saved locally and only available on the client used to create it. Logging out, restarting, or resetting Linear will clear this type of draft. 

If you use `Esc` or click on the close button, a pop-up modal will appear giving you the option to save the issue as a draft. This draft type persists across clients and will not clear on logout or reset. To access your saved drafts, open the Drafts page in your sidebar.

> [!NOTE]
> Drafts are stored for 6 months before being deleted automatically

## Editing issues

Edit an issue title or description by clicking directly on the title or description and editing inline. 

## Revert/Restore issue description

Use `Cmd/Ctrl K` and search for **Issue description history**. Or open the issue menu and select **Show description history**.

Open issue description history to view and restore earlier versions of the description.

## Move an issue to another team

When work needs to be passed over to another team, or when you are consolidating teams, issues can be moved to the appropriate Linear team within the same workspace.

For a single issue, simply use `Cmd/Ctrl Shift M` to move an issue to a new team. To move issues in bulk while retaining as much data as possible, select issues manually or with filters before moving them. Use `Cmd/Ctrl A` to select all issues on the list or board.

You can undo the move with `Cmd/Ctrl Z`. Most fields are restored, but changes to labels, subscribers, estimates, or access-related assignments may remain.

### Old Issue IDs and URLS

When you move an issue to a new team, we generate a new issue ID and unique URL for the issue. Old URLs will still work and redirect to the new issue URL. Searching for old issue IDs will also bring up the current issue (unfortunately, this doesn't work for old issue titles). Inline references to issues (like #ENG-123) will redirect when clicked, but won't update visually from the original issue ID they're associated with.

Some fields may be remapped or cleared based on the destination team's configuration, as detailed below.

Issue property | Effect | Workaround
--- | --- | ---
Cycle | May be cleared | The cycle may be cleared if there isn't a corresponding cycle in the destination team.
Team Labels | Removed | Create a label in the new team with the same name.
Projects | Removed | Add the new team to the current team's project before moving the issue.
Relations | Remain
Priority | Remain
Issue ID | Changed | The issue receives a new identifier for the destination team. Previous identifiers remain searchable and continue to resolve to the issue.
Status | Changed | When you move an issue, Linear maps its status to the closest corresponding status in the destination team's workflow. If the destination team uses triage, open issues moved by someone outside that team will move to triage. Closed issues remain closed.



## FAQ

<details>
<summary>Can you create issues and use keyboard shortcuts in Safari?</summary>
If Safari is stealing your focus when hitting tab during issue creation, enable this Safari advanced preference:   
  
Safari > Preferences > Advanced > Enable "Press tab to highlight each item on a webpage".

![Safari accessibility preferences modal showing press tab](https://webassets.linear.app/images/ornj730p/production/d1b238d1514c7004f9e3b3d9f4431d4bb27a8a38-992x622.png?q=95&auto=format&dpr=2)
*Safari Accessibility Preferences*
</details>

<details>
<summary>Is there an attachment file size limit when creating issues through email?</summary>
Email routing does not support sending messages with file sizes greater than 25MB. If you have attachments exceeding that, it will fail to deliver. Message body must also be less than 250,000 characters.
</details>