# Loops

![Loops](https://webassets.linear.app/images/ornj730p/production/d50a83fccc38ccc78ad8c0092dbd8796d092e759-3600x2080.png?q=95&auto=format&dpr=2)

## Overview

> [!NOTE]
> Available on paid plans. Loops that use Linear Agent draw from your workspace’s [AI credits](https://linear.app/docs/ai-credits).

Loops let Linear drive work forward automatically. They can run on a set schedule or when supported triggering events occur on [issues](https://linear.app/docs/creating-issues), [projects](https://linear.app/docs/projects), [initiatives](https://linear.app/docs/initiatives), [cycles](https://linear.app/docs/use-cycles) or [releases](https://linear.app/docs/releases). 

Loops follow your written instructions, along with any triggering changes to an issue’s status, priority, assignee, cycle, project, or labels. 

They can do work in your Slack workspace through Linear’s [Slack integration](https://linear.app/docs/loops#use-slack-with-loops), and with other optional [MCP connectors](https://linear.app/docs/loops#mcp-connectors).

## Example loops

Here are some basic loops that are useful for most workspaces. Create custom loops for your team by prompting Linear, or get started with one of these examples.

### Investigate and delegate bug reports

When new bugs arrive in Triage, check the codebase for a root cause. Delegate the issue to Linear to fix if the cause is found. ([Build this loop ↗](https://linear.app/agent?skill=loop-bug-triage))

![Video](https://webassets.linear.app/files/ornj730p/production/aaafc7577b0e892bef05a313a18d0692bba838af.mp4)

### Create follow-up issues from resolved incidents

When an issue with an incident label is marked done, do a root cause analysis and create follow-up issues. ([Build this loop ↗](https://linear.app/agent?skill=loop-incident-follow-up))

![Video](https://webassets.linear.app/files/ornj730p/production/c4c37995fa812bf22c56261c7c11e97eae861dc0.mp4)

### Create platform specific versions of incoming requests

When a new incoming issue is created, determine if new issues also need to be created for multiple platforms. ([Build this loop ↗](https://linear.app/agent?skill=loop-cross-platform-handoff))

![Video](https://webassets.linear.app/files/ornj730p/production/94fd2b994d955887413cbbbbfe0c2776816cb632.mp4)

### Create user-facing messaging after issues close

After an issue closes, check the PR that closed it and generate user-facing messaging that Support can share with customers. ([Build this loop ↗](https://linear.app/agent?skill=loop-customer-follow-up))

![Video](https://webassets.linear.app/files/ornj730p/production/7eea0fd4bd8d2f019cb4f846cdb4b57ac444e1a4.mp4)



## How Linear uses loops

Teams at Linear use loops to automate work. Learn more about some of our own use cases:

* [Slow PostgresSQL queries are surfaced for review](https://linear.app/learn/database-query)
* [Technical sales questions get automatically answered in Slack ](https://linear.app/learn/docs-maintenance)
* [Customer feedback is routinely collated for Product](https://linear.app/learn/customer-feedback-loop)
* [When help documentation fall out of date, a loop suggests changes](https://linear.app/learn/docs-maintenance)

## Create loops with Linear Agent

The best approach is to first achieve the outcome you want in a Linear Agent chat, then refine it into a loop. For example, you might be leading a project that has check-in meetings every Thursday morning. You're about to write the agenda, but realize you could automate this instead.

* Press ⌘/Ctrl + J to open a Linear Agent chat
* Give Linear Agent a prompt such as:

> Review the project’s Slack channel, the open questions in the latest project update, and any relevant issue threads. Add ,**Open questions**,, ,**Relevant issue threads**,, and ,**Next steps**, sections to the Meetings project document. Summarize your findings under each heading and include links to the relevant messages or issues.

* Review the result to make sure it looks right. If it doesn't, continue prompting Linear until the agenda looks complete, and in the desired format.
* Then, ask Linear to create the loop. You’ll be prompted to add any missing details.

## Create loops manually

You can also create a new loop directly in the interface.

1. In your sidebar, click **Loops** under the workspace group to start a new workspace-level loop. To create a new team loop, click the team name in your sidebar, then choose the **Loops** tab on that page.
2. Click **New loop**
3. Choose a trigger: a supported event on an issue, project, initiative, or cycle, or a schedule
4. Add instructions describing the outcome you want Linear to achieve
5. Optionally, add [MCP connectors](https://linear.app/docs/loops#mcp-connectors)
6. Review the loop's scope and permissions
7. Click **Create loop**

![Create a Loop](https://webassets.linear.app/images/ornj730p/production/a74f69f266bec89619ddeec9d71d5571e6fe30e7-3088x2088.png?q=95&auto=format&dpr=2)

## Update loops

Prompt Linear to update an existing loop, or make changes directly in the interface:

1. Open the list of **Loops**
2. Select the loop you want to update
3. Click **Edit**
4. Update the loop fields as needed. These changes are saved as a draft, and not yet live.
5. Click on **Publish** to apply all changes at once to the loop

![Edit a Loop](https://webassets.linear.app/images/ornj730p/production/de5836e670cc38d0b5b8818e266dd8b94cc184c1-1800x1200.png?q=95&auto=format&dpr=2)

## Review a loop's runs

Each loop has a run history to let you audit its behavior.

Open a loop’s **Run history** to ask Linear Agent follow-up questions or request additional actions using the run’s context.

1. Open **Loops** from the workspace or team that contains the loop
2. Open the loop you want to review
3. Select **Run history**
4. Review recent runs to see when the loop executed and what actions it took

![Run History Menu](https://webassets.linear.app/images/ornj730p/production/831089d29114df09dbc0dd742b4fd6448121d0aa-1988x882.png?q=95&auto=format&dpr=2)

![Run History](https://webassets.linear.app/images/ornj730p/production/9abb797a2ce8eb11c7072cbc780468b81144a63f-3088x2088.png?q=95&auto=format&dpr=2)

## Restoring a previous version

All published versions of a loop are saved and can be restored. Please note that removed connectors cannot be restored as they need to be authenticated again manually.

1. Start editing a loop as usual
2. Click on **Published versions**
3. Select a past version you want to restore
4. Click on **Restore version**

## Use Slack with Loops

Connect a Loop to your workspace’s Slack integration to let it work in Slack as Linear. A Loop can use one connected Slack workspace. You must link your Slack account to Linear and be a full member of the workspace to connect.

To connect Slack to a Loop:

1. A Linear admin must first turn on **Enable Linear Agent**, and **Allow Loops access** via [Settings > Integrations > Slack > Workspace](https://linear.app/settings/integrations/slack)
2. In the Loop draft, open the connector menu and select **Slack**
3. Select the appropriate Slack workspace

![Slack Integration Connector](https://webassets.linear.app/images/ornj730p/production/34f6d699f4080f76d2789955b00efc93158df8ec-989x389.png?q=95&auto=format&dpr=2)

The Slack integration uses your workspace’s shared Slack connection, rather than a personal Slack MCP connection. It is intended for shared automation with defined public-channel access. You should use a personal Slack MCP connector when a Loop needs access based on an individual’s Slack authorization or broader Slack tools.

### What Loops can do in Slack

* Loops can find public channels by name, including Slack Connect channels, and read their message history. They can also read messages from Slack links, post new messages, and add or remove emoji reactions. Messages posted to Slack by a loop are identifiable by a direct link to its [run](https://linear.app/docs/loops#review-a-loop's-runs).
* Loops can send direct messages to full members in Slack, and those messages are one-way
* In a public-channel thread started by a Loop, anyone with access to the loop run can reply from there to continue the conversation
* Loops can send broadcast mentions such as `@here` and `@channel` in Slack messages

### Slack access and limitations

* Loops cannot access private channels, direct-message history, or group direct messages
* Loops cannot search messages across the workspace, create channels, edit or delete messages, upload files, or manage Slack canvas documents
* Replying to a loop message in Slack Connect channels, direct messages and unrelated threads do not continue a conversation
* When automatic channel joining is disabled in Slack, Loops can only use public channels that Linear has already joined. If message history is unavailable, Loops can send messages into those channels but cannot read their history, or continue conversations through replies.

## MCP Connectors

MCP connectors extend what Linear Agent can do during a loop run. When a connected MCP is available to a loop, Linear can use it to gather context or take supported actions in other services. For example:

* Search or retrieve content from connected sources such as GitHub, Notion, or Sentry
* Post comments or updates to external services like Slack
* Fetch documentation, pull request details, or error reports to enrich an issue

An MCP connector can only access data and perform actions that are permitted by the available MCP’s configuration and the loop's data scope. When a user connects an MCP to a loop, they authorize that loop to act in the other service, so connector access should be granted thoughtfully.

![MCP Connectors available during a loop run](https://webassets.linear.app/images/ornj730p/production/2978ce6c0e93bc3e51f484ec013f72c0f7acbe47-996x432.png?q=95&auto=format&dpr=2)

### Managing available Connectors

Owners (Enterprise) or workspace Admins (Business) can elect to let users connect any MCP connector, or only allowlisted servers.

![security settings MCP server section -- All severs or only specific servers](https://webassets.linear.app/images/ornj730p/production/56dc4e78b28dc26e05409d93c0f7ef8e1ad2a3d7-950x634.png?q=95&auto=format&dpr=2)

If _Only specific connectors_ are selected, each MCP must first be approved in your workspace before it can be used in a loop.

Workspace admins and owners can configure the full set of MCP connectors allowed in a workspace in Security settings. Once an MCP server has been added, users can add it when creating or updating a loop. Adding one directly to a loop requires the configuring user to authenticate with the connector, and it will use their personal account's connection.

## Setup

* Workspace owners can manage who can create and manage workspace loops in [**Settings → Security → Workspace management → Manage loops**](https://linear.app/settings/security)
* Team owners can control who can create and manage loops in **Settings → Teams → [Team name] → Access and permissions → Loop management**
* To change who can edit a specific loop, open it and then select **Who can edit**. Choose between all members, workspace admins or team owners, or only the loop owner.

![Manage Loops Access](https://webassets.linear.app/images/ornj730p/production/358fe38b142196a3259b9bbb79290df811462fb5-1924x850.png?q=95&auto=format&dpr=2)

## Permissions

Loops support a range of permissions to control what they can access and do. Enable only the permissions your loop needs to fulfil its intended purpose, and consider the impact of each one before enabling it.

![Loop Permissions](https://webassets.linear.app/images/ornj730p/production/b7815db70390ffeb5a9a60a075cf6f27d255e9e5-1988x1704.png?q=95&auto=format&dpr=2)

### Team access

Controls which teams the loop can access. It can read and write data only in those teams. By default, a loop at the workspace level or in a public team will have access to all public teams, while a loop inside a private team will have access only to that team.

Loops that can access all public teams can also access workspace-level objects like Initiatives and Customers.

### Web access

If enabled, the loop can query any website. You can use this to set up loops that can take actions like researching competitor announcements, or checking the documentation of services your product integrates with.

Use caution when enabling web access. **Web access can send workspace content to external services.** Do not enable it for loops that may process sensitive data unless you are comfortable with that data leaving Linear.

### Code Intelligence

If enabled, the loop can use [Code Intelligence](https://linear.app/docs/code-intelligence) to browse and analyze repositories configured in your workspace. Enable Code Intelligence for loops designed to investigate bugs or answer questions about your code.

### Coding sessions

If enabled, the loop will be able to start a [coding session](https://linear.app/docs/coding-sessions) to open a draft pull request. Allow coding session access for loops that delegate implementation work to Linear.

### Externally synced issues and comments

If enabled, the loop will be able to write data on issues or comment threads that are synced with external applications, like issues created from Slack with bi-directionally synced threads.

Enable this with care, as some synced threads may be visible outside your workspace - if you use [GitHub sync](https://linear.app/docs/github#configure-github-issues-sync) in a public repo for instance, the loop would be able to post comments to that repo's GitHub Issues.

### External sources

Issues can be created from sources outside Linear, like Slack messages or emails received to a particular address. For security reasons, loops will only run by default on issues created from within Linear.

If you want a loop to run on issues created from  external sources, you can enable specific external sources on a per-loop basis. Workspace owners can configure the list of allowed external sources in [Security](https://linear.app/settings/security) settings. After that, configured external sources can be enabled for individual loops.

### Allow changes outside of triggering entity

If enabled, the loop can write data to other entities included in the **Team access** scope. If disabled, the loop can only write data on the entity that triggered the loop for a given run. This setting does not apply to scheduled loops.

---

## AI credits

Loops use [AI credits](https://linear.app/docs/ai-credits), which workspace admins can purchase and manage in [Settings](https://linear.app/settings/billing). Purchased credits expire 12 months after the purchase date. If a workspace runs out of credits, Loops pause unless an admin purchases more credits or enables auto top-ups. Workspaces are charged automatically only when auto top-ups are enabled.

**Promotional Loops AI credits**

Some workspaces may receive promotional Loops AI credits with separate expiration terms. Eligible workspaces will be notified by email.

---

## FAQ

<details>
<summary>What are the best practices for writing instructions?</summary>
Clearly describe the intended outcome and any actions Linear Agent should avoid.

For example:

> Investigate the issue using the available context. Add a comment summarizing the likely cause and recommend the next action. Do not change the issue's assignee or status.

For best results:

* Use Linear Agent chat to help build your instructions. See [Create loops with Linear Agent](https://linear.app/docs/loops#create-loops-with-linear-agent) for guidance.
* Describe the outcome, not only an action
* Specify which context or connected tools Linear Agent should use
* State which changes are permitted
</details>

<details>
<summary>What MCP connectors can loops use?</summary>
A loop can only use MCP connectors that are connected to your workspace, enabled for Linear Agent, and allowed by your workspace admins. Access also depends on the loop's scope and permissions.
</details>

<details>
<summary>When should web access be enabled?</summary>
Enable web access only when a loop needs information from external websites.

A good fit is research on public information, like competitor announcements, documentation, or API references, where the loop does not need to send sensitive workspace content to external services.

Enable this intentionally, since allowing external web access can increase the risk of exposing workspace data outside Linear.
</details>

<details>
<summary>How do I use Code Intelligence?</summary>
First, a workspace admin must enable Code Intelligence and configure code access for the workspace. Then, enable the [Code Intelligence permission](https://linear.app/settings/ai/code-intelligence) on the loop itself.
</details>

<details>
<summary>How do I use Coding Sessions?</summary>
First, a workspace admin must enable coding sessions and configure code access for the workspace. Then, enable the [Coding Sessions permission](https://linear.app/settings/ai/coding-sessions) on the loop itself.
</details>

<details>
<summary>How can I review loop costs and AI credit usage?</summary>
Open **Settings → Usage & limits** to review AI credit usage. Loop spend includes the loop’s own work and any coding sessions it starts. Learn more about [usage and spend limits](https://linear.app/docs/usage-and-limits).

If your workspace runs out of AI credits, loops that use Linear Agent stop running until more credits are added. Loops that only apply configured issue updates can continue to run without AI credits.
</details>

<details>
<summary>Why did I get a loop run failed notification?</summary>
To troubleshoot, open the loop's **Run history** to check recent runs and any visible failures.

Common causes a loop run might fail include missing permissions, untrusted external sources, or your workspace running out of AI credits.

To get a more detailed description of what went wrong, you can also ask Linear. Press ⌘/Ctrl + J while looking at the run history and ask for a detailed description of why the loop run failed.
</details>

<details>
<summary>Can I run a loop manually?</summary>
You can use Run now for enabled loops that use Linear Agent and have a schedule, issue, project, or initiative trigger.

For issue, project, and initiative triggers, choose the entity you want to use for the run. Run now is not available for cycle triggers or loops that only apply configured issue updates.
</details>

<details>
<summary>What schedules can loops use?</summary>
Scheduled loops can run hourly, daily, weekly, monthly, or yearly. Choose the timezone and when the loop should first run.

Weekly schedules can run on one or more days of the week. The first run date cannot be in the past.
</details>

<details>
<summary>How do I enable or disable a loop?</summary>
Right-click on the loop from either the workspace or team loop view and select **Enable** or **Disable.**
</details>

<details>
<summary>How do triage loops work?</summary>
Triage loops are created for a specific team and require triage to be enabled for that team. They run when an issue enters the team's triage queue and begins matching the configured conditions.

Triage loops do not support schedules or non-issue triggers.
</details>

<details>
<summary>Are there limits on loops?</summary>
Teams can create up to 100 loops, and workspaces can create up to 500 workspace-level loops. Each person can keep up to 100 unpublished loop drafts.

Loop names can be up to 64 characters, and descriptions can be up to 255 characters. A loop can include up to eight conditions and eight actions.
</details>

<details>
<summary>What happens when I delete a loop?</summary>
Deleting a loop is permanent and cannot be undone. Consider disabling the loop instead if you may need to restore it in the future.
</details>

<details>
<summary>When does a Loop with configured conditions run?</summary>
For Loops with a set of defined conditions, those conditions must be met before the Loop will runs. It also won’t run again on subsequent updates, if those same conditions remain matched.
</details>