# Connect MCP Servers to Linear Agent

![DataDog and Sentry MCP servers connected with Linear Agent to help analyze a bug report.](https://webassets.linear.app/images/ornj730p/production/74b8e0c4616b9593aad2293afb62477c096652a4-758x1134.png?q=95&auto=format&dpr=2)

## Overview

Linear Agent can connect to MCP servers to access tools and data outside of Linear.

This lets the agent bring external context into your workflows to investigate issues, plan projects, write specs, and draft updates grounded in your full context.

MCP works across Linear Agent chat, comments, and automations, and is enabled at the workspace level by Admins. 

### Available MCP servers

* Amplitude
* Attio
* Better Stack
* Datadog
* GitHub
* Glean
* Granola
* HubSpot
* incident.io
* Intercom
* Jam
* Notion
* PostHog
* Sentry
* Slack
* Stripe

## Setup

A workspace owner or admin can enable or disable MCP support for the workspace, restrict usage to approved MCP servers, and manage workspace-connected MCP servers.

Inherited from workspace-level settings, users can then configure their preferred MCP servers to leverage when collaborating with Linear Agent.

### Workspace configuration

1. Open [Workspace settings > Security](https://linear.app/settings/security)
2. Beside **Enable MCP servers** toggle it to **“On”**
3. Under **Allowed MCP servers**, choose **All servers** or **Only specific servers**
4. Scroll down to MCP servers and select any that users and automations can connect to via Linear Agent

![Enable MCP Servers for the Workspace](https://webassets.linear.app/images/ornj730p/production/1e9a947797b964167c0416c08563e85610f25eea-859x522.png?q=95&auto=format&dpr=2)

![Select available servers for the workspace including any custom servers using a URL](https://webassets.linear.app/images/ornj730p/production/ebb2f790ef72631c8dbeed19f49f43a371fabdc0-462x592.png?q=95&auto=format&dpr=2)
*Select available servers for the workspace including any custom servers using a URL*

### Individual server connection

Each user in your Linear workspace can configure their preferred MCP servers, based on restrictions defined at the workspace-level. 

1. Open [Settings > Agent personalization](https://linear.app/settings/account/agents)
2. Under MCP servers click on the **“+” (Add server)** button
3. Choose a server and start the connection flow
4. Complete authentication with the provider
5. Return to Linear and confirm the server shows as connected

![Connect Agent to MCP Server from Settings](https://webassets.linear.app/images/ornj730p/production/c1ead295c6b114f9369f7d603e03d254ee8d477d-1974x1018.png?q=95&auto=format&dpr=2)

After a server is connected, users can reference it in natural language when using Linear Agent. The agent can use the connected MCP server to gather information or perform supported actions. 

![Prompt example for MCP server](https://webassets.linear.app/images/ornj730p/production/461e38065d0a92659e96488a59c4f0b46e74d0dc-503x636.png?q=95&auto=format&dpr=2)

### Custom server connection

1. Click **+ Custom URL…**
2. Select **Custom URL…**
3. Enter your server URL (must start with `http://` or `https://`) and complete the connection flow:

* No authentication → connects immediately
* OAuth supported → sign in and authorize
* Authentication headers → enter your custom headers

Custom servers can authenticate with OAuth or with Authentication headers, depending on how the server is set up. Once connected, the custom server will be available for users to connect to through their settings.

![Custom server enablement](https://webassets.linear.app/images/ornj730p/production/4fb96d1ef5c9173eaf92f61910481671047a2815-1580x700.png?q=95&auto=format&dpr=2)

## Example prompts

> [Investigate this issue using Sentry and summarize the likely cause](https://linear.app/agent?prompt=Investigate%20this%20issue%20using%20Sentry%20and%20summarize%20the%20likely%20cause)

> [Check Datadog for recent errors related to this workflow](https://linear.app/agent?prompt=Check%20Datadog%20for%20recent%20errors%20related%20to%20this%20workflow)

> [Look up the latest notes about this customer in Intercom](https://linear.app/agent?prompt=Look%20up%20the%20latest%20notes%20about%20this%20customer%20in%20Intercom)

> [Search Notion for relevant documents and generate a guide ahead of launch](https://linear.app/agent?prompt=Search%20Notion%20for%20relevant%20documents%20and%20generate%20a%20guide%20ahead%20of%20launch)

> [Use GitHub to find the pull request that likely introduced this change](https://linear.app/agent?prompt=Use%20GitHub%20to%20find%20the%20pull%20request%20that%20likely%20introduced%20this%20change)