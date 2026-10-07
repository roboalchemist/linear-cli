# Reviews

![Review your delegated agent's code changes](https://webassets.linear.app/images/ornj730p/production/e6108a650d6fd66011fbde5cbb8446ec2b46cd9a-3330x2194.png?q=95&auto=format&dpr=2)

## Overview

Linear Diffs let you review code changes directly in Linear—so you can read diffs, follow conversations, and complete reviews without leaving your workflow. When enabled, Linear shows pull request details, changed files, checks, and comments, and keeps them in sync with GitHub.

A Reviews section appears in the sidebar where you can see pull requests that need your attention, track pull requests you authored, and return to reviews you've participated in. Opening a PR in Linear shows its details, activity, CI checks, and associated comments, and updates sync bidirectionally between GitHub and Linear automatically.

Navigate to Linear Code quickly with the shortcut `G` + `R`, or open a specific review with `O` + `R`.

## Setup

### Organization setup

To display pull request diffs and file changes inside Linear, you’ll need to grant Linear access to your repository code through the GitHub integration.

If your workspace hasn’t connected GitHub yet, a Linear owner or admin can [configure the GitHub integration](https://linear.app/docs/github#enable-the-github-integration) to enable code access for the repositories you selected.

If your workspace has already connected GitHub, a GitHub organization owner can enable code access for the selected repositories by updating the integration in your [GitHub integration settings](https://linear.app/settings/integrations/github).

Your current GitHub integration settings and any existing pull request links and history in Linear are preserved—enabling code access simply adds permissions so Linear can display diffs and file contents for the repositories you select.

<details>
<summary>GitHub organizations using IP allow lists</summary>
GitHub enforces your organization’s IP allow list for review actions performed as the authenticated GitHub user, even when the Linear GitHub App is installed.

To allow Linear to load pull request diffs and perform review actions, add Linear’s IP addresses listed under **GitHub Enterprise options → GitHub Enterprise Cloud** in the [GitHub integration documentation](https://linear.app/docs/github).
</details>

![Enable your code reviews for your Github organization](https://webassets.linear.app/images/ornj730p/production/072dab8e645b2c0c53ba1f8e665658828090ddea-1694x1212.png?q=95&auto=format&dpr=2)

### Personal account setup

Personal GitHub connections are required for accessing pull requests, repository code, and review information associated with your specific user. If your account isn’t already connected, you can do so [here](https://linear.app/enable-reviews).

![Personal account connection](https://webassets.linear.app/images/ornj730p/production/6cdb4e1d2f0be76e9d4e5e2658520b300df4acc1-1526x1080.png?q=95&auto=format&dpr=2)

> [!NOTE]
> If you don't see the Reviews tab in your sidebar, navigate to [Settings → Code & reviews](https://linear.app/settings/account/code-and-reviews) and confirm the _Enable code reviews_ toggle is switched on.
> 
> If you see the Reviews tab but no diffs, confirm that [code access is granted](https://linear.app/docs/diffs#organization-setup) to the repository you expect.

## Display and grouping

At the top of the Reviews tab, you can switch between two tabs: **For me**, which shows the pull requests you're involved in or responsible for, and **Created**, which displays every pull request you've authored. This makes it easy to focus on work that needs your attention while still being able to browse or reference other pull requests when needed.

You can group by status, author, or repository, and sort or filter to hide draft and closed PRs.

You can choose whether to display draft or closed pull requests and whether to show additional fields such as repository, failed checks, or preview links if applicable.

![Display options for your Reviews](https://webassets.linear.app/images/ornj730p/production/7f5086a8ef00e02de7cd8b6eb1ec9ae06306af39-760x882.png?q=95&auto=format&dpr=2)

If your workspace uses GitHub team reviewers, display options  let you control whether team-level review requests appear in your Reviews list and whether they count toward the Reviews badge, which can help reduce noise when you’re part of larger review groups.

![Team review Display options](https://webassets.linear.app/images/ornj730p/production/e6cedd7eb3f05124eb6d18952266420c8ef17fe8-760x894.png?q=95&auto=format&dpr=2)

## Code review

### Review a pull request in Linear

Linear displays all pull request activity in one place, including comments, reviews, and discussions. When you open a pull request, you can view the files changed and read through diffs directly in Linear.

Line counts in pull request diffs show implementation changes by default, excluding things like tests and documentation. When this implementation count differs from the total changed lines, Linear shows a `[*]` indicator next to the count. Hover over it to see the full total line count.

Comments appear inline alongside the relevant sections of code and you can start new threads, reply to existing ones, or react with emoji.

You can complete your review directly from Linear. Approving a pull request, requesting changes, or submitting a review comment updates the review state and syncs with GitHub.

Once a pull request is ready and you have permission, you can merge it directly from Linear.



![Linear Code](https://webassets.linear.app/images/ornj730p/production/3d63a6ce70f4556fdac04c678575c37f0774b958-2450x1024.png?q=95&auto=format&dpr=2)

### Guides

Guided reviews help you understand large pull requests faster by organizing related changes into structured sections with explanations of their purpose and impact.

> [!NOTE]
> Guides is available on Basic, Business and Enterprise plans.

Instead of working through files one-by-one to piece together what changed, Guides surface the core parts of an implementation first while grouping supporting or lower-signal changes separately. Each section pairs a high-level explanation of _why_ a part of the change exists alongside the relevant diffs, making it easier to understand the intent behind the implementation before diving into the code itself.

Guided reviews appear in a dedicated **Guide** tab alongside the diff view, with direct links into the relevant parts of the pull request so you can move naturally between the guide and the underlying code.

![Guided review to help you understand large pull requests faster](https://webassets.linear.app/images/ornj730p/production/c6edafc1dde5f8863154ef98bb879ec5492bbc4e-2648x1914.png?q=95&auto=format&dpr=2)

If you’d prefer to disable this feature, you can disable the _Generate Pull Request guides_ toggle within your GitHub integration settings' _Pull Requests_ section. 

![Disable guided reviews](https://webassets.linear.app/images/ornj730p/production/c46e1c258d3a708c9c29b1643147755c7a8a956e-1456x418.png?q=95&auto=format&dpr=2)

### Diff view options

#### Unified and split view

Linear supports both _Unified_ and _Split_ diff views when reviewing a pull request.

Unified shows changes in a single column, while Split shows the before and after side-by-side to make comparisons easier. You can switch between Unified and Split from the pull request display options, and you can also toggle between Split and Unified view with `Ctrl` or `⌘` + `B`. 

Note that on smaller screens, Split view may be unavailable if there isn't enough horizontal space.

#### Stacked diffs

Stacked diffs let you split a larger change into smaller, dependent pull requests, with each building on the one before it. Linear displays the stack’s order and your current position, so you can move between pull requests while reviewing. Each member’s review, check, and merge status helps you understand what’s ready and what still needs attention.  
  
Linear supports [GitHub-native stacks](https://docs.github.com/en/pull-requests/get-started/about-stacked-prs) and linked pull requests inferred from their branch relationships. For native stacks, merging a pull request also merges every pull request below it, leaving those above it open. Linear shows how many pull requests the action affects and tracks GitHub’s asynchronous merge progress. Linked branch chains merge individually, not as a single stack.

  
Stacks are created and managed outside Linear; Linear doesn’t create, reorder, rebase, or restack them, however you can ask your agent to use stacks.

#### Structural highlighting

When structural highlight is disabled, the diff uses a standard line-based comparison, so added and removed lines are emphasized. When it’s on, Linear uses syntax-aware highlighting to better understand the structure of the code and highlight the specific parts of a line that changed, such as renamed variables, edited expressions, or moved code blocks.

![Unified and Split diff view options](https://webassets.linear.app/images/ornj730p/production/7d0e5b91472e0fb3962333a783d41a951429298f-788x900.png?q=95&auto=format&dpr=2)

### Open GitHub PR URLs in Linear

If you have an existing GitHub pull request URL, you can also open that same pull request in Linear by replacing `github.com` in your URL with `linear.review`. 

**Original Github URL**: `github.com/owner/repo/pull/123`

**New Linear URL**: `linear.review/owner/repo/pull/123`

With this pattern, Linear will redirect you to the matching pull request review page in your Linear workspace automatically.

### Organize files with ,`.gitattributes`

Use `.gitattributes` to organize files in pull request diffs by review category. This helps reviewers focus on implementation changes while keeping tests, documentation, generated files, and other supporting changes easy to find. 

Add rules to your repository’s `.gitattributes` file:

```sh
tests/** review-test
docs/** review-documentation
generated/** linguist-generated=true
*.stories.tsx -review-test
```

Use `.gitattributes` to group pull request files by type. Linear supports `review-implementation`, `review-test`, `review-documentation`, `review-generated`, `review-agent-guidance`, `review-localization`, and `review-assets`. GitHub-compatible `linguist-generated` and `linguist-documentation` work too.

Add a category with `attribute` or `attribute=true`; remove one with `-attribute` or `attribute=false`. In a pull request, choose **File type** grouping to browse the resulting categories.

### Pull request preview links

If your PR contains one or more preview links, this will add a preview link shortcut to the Linear issue. More details on this feature [here](https://linear.app/docs/github#pull-request-preview-links).

![Preview links on a Linear issue](https://webassets.linear.app/images/ornj730p/production/a5b099a909766d6a2718826182155ba372563cf2-2140x1105.png?q=95&auto=format&dpr=2)

## Notifications

With Diffs enabled, Linear can notify you about pull request activity in your inbox and push notifications. 

In  [Code & reviews settings](https://linear.app/settings/account/code-and-reviews), you can choose which kinds of review activity reach you. 

* **Comments & reviews** covers comments, mentions, and submitted reviews, and you can choose whether to include all activity or exclude activity from bot actors. **Review requests** controls requests for your personal review. 
* **Team review requests** controls review requests that come through your GitHub teams, which apply when a pull request is assigned to a GitHub team you belong to rather than directly to you. Note that team review requests are only sent for GitHub teams with 10 or fewer members.
* **Checks & merge queue** controls notifications about failed checks and merge queue updates. 

![Review notification settings in Code & reviews](https://webassets.linear.app/images/ornj730p/production/d68f78471d3d578b171db1f3a00d41c4d5f0178b-1524x872.png?q=95&auto=format&dpr=2)

## Additional settings

Navigate into [Code & reviews settings](https://linear.app/settings/account/code-and-reviews) for additional preferences to the Reviews experience within Linear. These settings are personal, so they only affect your own workflow.

### Auto-convert draft PRs

When enabling the **Auto-convert draft PRs** toggle, this setting automatically moves a draft pull request into a ready-for-review status once a review is requested or when the pull request is approved. 

![Auto-convert draft PRs setting](https://webassets.linear.app/images/ornj730p/production/9be5ef8a7f2cb2ced47e157fbf5e285894b0e5fb-1128x216.png?q=95&auto=format&dpr=2)

### Code theme

You can control the appearance of code while manually reviewing the diff or following a guided review. This includes separate themes for light and dark mode, along with display settings that provide control over the font size and line height.

![Code theme settings](https://webassets.linear.app/images/ornj730p/production/d66cca37e485c7d374d410f394bf50fcf452c439-1184x744.png?q=95&auto=format&dpr=2)

### Move issues to start

Save yourself a few steps by toggling on the automation that moves a Linear issue into a started status when you copy the git branch name. To set up this automation, refer to the _On git branch copy, move issues to a started status_ toggle in your [Code & reviews settings](https://linear.app/settings/account/code-and-reviews).

![Setting to move issues to a started status on git branch copy](https://webassets.linear.app/images/ornj730p/production/fe2e9b02774759562da7a6c1b70e56ed2591b010-1156x238.png?q=95&auto=format&dpr=2)

Extend this further by enabling automation within your [Preferences](https://linear.app/docs/account-preferences#auto-assign-to-self) that assigns the issues you move to a started status to yourself. 

## Review Platform

### On behalf of attribution

Pull request activity can show when a comment was made by an AI coding agent on behalf of a person. This attribution makes it clear which agent produced the comment and, when available, the model it used.

You’ll see the agent’s name and avatar alongside the comment, helping reviewers understand the source of the feedback or update. Supported agents include Claude Code, Codex, Linear, Pi, and OpenCode.

<details>
<summary>On behalf of attribution example</summary>
![on behalf of example](https://webassets.linear.app/images/ornj730p/production/4fbc8b2de3968efe9e9c8201d923b051d3d70b21-1778x548.png?q=95&auto=format&dpr=2)
</details>

### Risk score

Risk scores help reviewers quickly assess the potential impact of a pull request. When a connected provider supplies a score, Linear displays one of four levels:

Low risk, Typical risk, High risk, Very high risk

The score may also include an explanation, along with who reported it and when it was last updated. Use this context to prioritize review attention, it complements, rather than replaces, your team’s normal review process.

<details>
<summary>Risk score example</summary>
![typical risk showing in linear on a review](https://webassets.linear.app/images/ornj730p/production/0871528769766bed8ad410e95b71d9d17da36afe-887x419.png?q=95&auto=format&dpr=2)
</details>

### Add data from GitHub

Connected tools can add this metadata in an HTML comment on a GitHub pull request, review, or review comment. Use a `linear:extension` block with version `1` and one or more plugins:

```json
<!-- linear:extension {"version":1,"plugins":[...]} -->
```

### Risk score

```json
<!-- linear:extension {
  "version": 1,
  "plugins": [
    {
      "plugin": "riskScore",
      "sha": "0123456789abcdef0123456789abcdef01234567",
      "level": 3,
      "explanations": [
        "Touches authentication",
        "Changes database migrations"
      ]
    }
  ]
} -->
```

* `level` is required and must be an integer from `1` (low risk) to `4` (very high risk).
* `sha` is optional and must be a 40-character commit SHA. Omit it to apply the score to the pull request’s current revision.
* `explanations` is optional: up to 8 non-empty strings, each up to 200 characters.

### On behalf of

```json
<!-- linear:extension {
  "version": 1,
  "plugins": [
    {
      "plugin": "onBehalfOf",
      "agent": "codex",
      "model": "gpt-5"
    }
  ]
} -->
```

* `agent` is required. Supported values are `claude`, `codex`, `linear`, `pi`, and `opencode`; other values are shown as a generic agent.
* `model` is optional and can be up to 200 characters.

Multiple plugins can be included in the same block. Extension-only comments are hidden from Linear by default. Add `"visible":true` to preserve the GitHub comment in the activity feed. On-behalf-of blocks are visible by default unless `"visible":false` is set.

## FAQ

<details>
<summary>AI Privacy</summary>
Linear does not utilize your data to train its own AI models. Any data processed to enable Linear's AI features is shared with our trusted partners (AI subprocessors, see our DPA) exclusively to deliver those AI functionalities to you without permission to train on provided data.

To provide features powered by AI and large language models (LLMs), Linear utilizes voluntary data provided by the user in terms of labeling feature outputs (thumbs up/down) or in other opt-in ways. If you have any questions or concerns, please let us know at security@linear.app.

For further information, please see AI Security FAQ in our [Trust Center](https://trustcenter.linear.app/).
</details>

<details>
<summary>How do I update a PR that is showing the wrong state?</summary>
This can happen sometimes if a webhook is missed between Linear and GitHub and we fail to capture the change in PR state automatically.   
  
You can trigger an update to the PR state shown in Linear by making a small edit to the PR description in GitHub — for example adding and removing a space.
</details>

<details>
<summary>Why didn't someone get notified on GitHub when I mentioned them from Linear?</summary>
GitHub only recognizes mentions for users it can map to a GitHub identity. If you mention someone in Linear who hasn't connected their personal GitHub account (or otherwise isn't mapped), GitHub won't treat it as a real @mention — so they won't get a GitHub notification. The fix is for that person to connect GitHub in Linear so mentions can resolve correctly.
</details>

<details>
<summary>Can I view and review a pull request commit by commit in Linear?</summary>
Yes. You can view a pull request’s changes commit by commit by selecting individual commits to filter on in the diff. Reviews are still submitted at the pull request level rather than as separate per-commit reviews, and you can comment on specific lines or changes.

![View a PR's changes filtered by individual commits.](https://webassets.linear.app/images/ornj730p/production/8a8c5d79abb2c12f2da54d91ee2e952fb2f5dfc6-1332x734.png?q=95&auto=format&dpr=2)
</details>

<details>
<summary>Why can’t I always create or view inline comments in a diff?</summary>
With GitHub’s current API, some inline comments can’t always be shown in Linear or created from Linear in the same way they appear in GitHub. 

We’re tracking improvements to make this experience more consistent.
</details>

<details>
<summary>Do you show check annotations or rich external status details (like inline CI annotations)?</summary>
Not yet. Linear shows overall check status and basic details, but doesn't currently sync or display rich check-run annotations (for example, inline failure locations or detailed external tool output).
</details>

<details>
<summary>Why don't draft/unfinished reviews sync between GitHub and Linear?</summary>
Draft reviews that are started in GitHub but not submitted don't sync into Linear today. Linear syncs review state when a review is actually submitted, and it doesn't mirror GitHub's draft review state due to API limitations.
</details>