# Coding settings

Review GitHub pull requests in Linear and set your preferences for how Linear works with your git workflow. Open [**Settings → Account → Code & reviews**](https://linear.app/settings/account/code-and-reviews) to configure these preferences.

![Code & reviews settings](https://webassets.linear.app/images/ornj730p/production/123dd417ae00aa099c201f9b934b926fe501de02-1556x1772.png?q=95&auto=format&dpr=2)

## Setup

### Enable code reviews

Code reviews require a workspace GitHub integration with repository code access, a connected personal GitHub account, and Reviews enabled in _Settings_ → _Account_ → _Code & reviews_. 

Workspace integration and code-access changes require workspace permissions, but each user manages their own GitHub connection individually. 

See the [Reviews setup guide](https://linear.app/docs/diffs#setup) for detailed setup instructions.

## Pull request preferences

### Auto-convert draft pull requests

Automatically mark your draft pull requests as ready when a review is requested or the pull request is approved. It applies only to your own draft pull requests.

### Merge strategy

Choose how you merge pull requests from Linear:

* **Squash and merge** (default)
* **Rebase and merge**
* **Create merge commit**
* **Choose for each pull request**

This is a personal preference. Repository configuration still determines which merge methods are available, so a method disabled in GitHub, or restricted by a merge queue, cannot be used even when it is selected here.

## Code theme

Select the syntax highlighting theme and typography used in code diffs and viewers.

Make code easier to read in the way that works best for you. Choose separate syntax themes for light and dark mode, including _Linear Colorblind_. This lets you tailor code review views to your display mode and accessibility preferences.

Some themes support additional display options. _Use theme background_ applies the theme’s own background colors, while _Show diagonal lines_ adds patterns to additions and removals so changes remain distinguishable without relying on color alone. These options only appear for compatible themes.

![Code theme settings in Linear](https://webassets.linear.app/images/ornj730p/production/6c33487a9b1952711dca93e0dd3bbcef6f987333-1312x1114.png?q=95&auto=format&dpr=2)

## Review notifications

These controls determine which code review events reach your Linear inbox and your push notifications. Each one is set separately.

* **Comments & reviews**: choose **All**, **Exclude Bots** to skip comments and reviews left by bots, or **None**.
* **Review requests**: get notified when your review is requested.
* **GitHub team review requests**: get notified when a review is requested from a GitHub team you belong to. This requires **Review requests** to be enabled and applies only to GitHub teams with 10 or fewer members.
* **Checks & merge queue**: get notified when checks fail or when a pull request is removed from the merge queue.

![Review notifications](https://webassets.linear.app/images/ornj730p/production/c8d1a0befa2f53ec70b73a94ba1d944cd59583e3-1378x812.png?q=95&auto=format&dpr=2)

## Signed commits

When _Require signed commits_ is enabled, commits made by coding sessions must be signed. This is a workspace-wide setting and applies to commits created by coding sessions, not to every commit pushed to GitHub. Changing it requires workspace admin (or owner if on Enterprise) permissions.

_Commit signing key_ is a personal setting on your own account. Add a signing key so that commits Linear makes on your behalf can be signed. It is relevant when you use coding sessions and code reviews in Linear: without a key, coding sessions cannot commit in a workspace that requires signed commits. Follow the instructions shown in the setting to generate a key and add it to your account.

## External tools

### Configure coding tools

Choose the coding tools you can open Linear issues in, such as Claude Code, Codex, Cursor, and more. Terminal-based tools require the Linear desktop app on a supported platform, because Linear needs to launch a local terminal. On macOS you can also pick your default terminal application; that selector is macOS only.

You can also add custom links to open issues in any tool that accepts a URL, and custom scripts to run a local command with issue context. See [coding tools](https://linear.app/docs/assigning-issues#open-issues-in-coding-tools) and [open issues with custom scripts](https://linear.app/docs/open-issues-with-custom-scripts) for the full workflows.

## Git attachment format

Choose how GitHub and GitLab attachments are displayed on issues: Title, or Title+Repository.

## Move issue to started status

Linear can move an issue to your team's first started status when you begin work on it. There are separate settings for copying the git branch name and for opening the issue in a coding tool, which includes **Copy as prompt**.

In both cases:

* Issues that have already started, and closed issues, are not moved.
* If the issue is unassigned, it may also be assigned to you when your separate [self-assignment preference](https://linear.app/docs/account-preferences#auto-assign-to-self) is enabled.
* Hold **Option** on macOS, or **Alt** on Windows, while performing the action to skip the status change for that one action.

### Open GitHub pull requests in

Choose whether pull requests open in GitHub or Graphite. This control appears once GitHub is connected and is separate from enabling reviews in Linear: it sets your external destination whether or not you review pull requests in Linear.