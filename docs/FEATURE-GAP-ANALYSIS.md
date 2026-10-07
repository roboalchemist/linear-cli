# Feature Gap Analysis — `linear-cli` vs. live Linear

Generated from the scraped corpus in [`linear-source-docs/`](linear-source-docs/)
and the live GraphQL SDL schema
(`linear-source-docs/external/raw.githubusercontent.com/linear/linear/master/packages/sdk/src/schema.graphql`).

**Baseline:** `linear-cli` was forked from `linctl`; first commit `2025-07-13`,
last code commit `2026-03-04` (94 commits, through `v0.5.0`). Linear keeps
shipping, so the CLI is roughly **seven months behind** the product and the API.

---

## Implementation status (update)

**Closed — read/list parity for every first-class entity.** Added read/list/get/
search commands for: customers (+ needs/statuses/tiers), releases (+ notes/
pipelines/stages), roadmaps, project & initiative labels, project & initiative
relations, templates, webhooks, time schedules, triage responsibilities, emojis,
audit log, agent sessions/skills, integrations, external users, notification
subscriptions, organization (+ invites), OAuth applications, team memberships,
audit entry types, initiative updates, and unified search
(`search issues|projects|semantic`).

**Closed — lifecycle operations.** Unarchive (issue/project/initiative/document),
issue delete (trash/permanent), issue label add/remove, comment resolve/
unresolve, issue subscribe/unsubscribe/share/unshare/reminder, issue batch
create/update, project label add/remove, project status archive/unarchive/purge,
project reassign-status, initiative archive/unarchive + label add/remove,
workflow-state CRUD (`team state create/update/archive`), cycle start-today/
shift-all, and bulk inbox actions (mark-all-read/unread, snooze-all,
unsnooze-all, archive-all).

**Closed — entity CRUD (write).** create/update/delete/archive for roadmaps,
releases, templates, webhooks, time schedules, and customers (+ status/tier,
plus customer merge); emoji create/delete; webhook secret rotation.

**Verified live (read-only).** `live_read_test.sh` runs every read/list command
against the real API: **65 pass, 6 skip (workspace feature/scope gates), 0 fail**.
Automated via `make test-live-read` and the *Live read-only tests* workflow.

**Still open (write-side, not live-testable here).** Integration connect flows
(~68 mutations), OAuth applications, email intake addresses, organization/user
admin (invites, domains, role changes, session revocation), agent session
creation, project milestone move, provider-specific attachment links, push
subscriptions, view preferences, git automation, package/webhook rotation
extras, and import (`issueImport*`) flows. The `graphql` command remains the
escape hatch for these.

---

## 1. Method

1. Scrape Linear's live docs + GraphQL SDL into `docs/linear-source-docs/`
   (`docs/scrape_linear_docs.py`).
2. Extract the API surface from the SDL: root `Query` and `Mutation` fields and
   all object types.
3. Extract the operations `linear-cli` actually issues from `pkg/api/*.go`
   (GraphQL mutation/query names) and its user-facing command tree from `cmd/`.
4. Diff: schema minus implemented = gap. Group by product domain and rank by
   value to a CLI user.

Regenerate the underlying corpus at any time with `make scrape-docs`.

---

## 2. API-surface coverage

From the live SDL:

| Surface | Count |
|---|---:|
| GraphQL object types | 729 |
| Root `Query` fields | 177 |
| Root `Mutation` fields | 383 |
| Mutations exercised by `linear-cli` | **51** |
| **Unimplemented mutations** | **332** (~87%) |

`linear-cli`'s `graphql`/`gql` command is an escape hatch, so nothing is
strictly impossible today — but only ~13% of the write API has a first-class,
typed command.

### Unimplemented mutations, by domain

| Domain | Missing | Notes |
|---|---:|---|
| `integration*` | 68 | Mostly OAuth/connect flows — low CLI value |
| `issue*` | 27 | Label ops, batch, subscribe/share, reminders, imports, releases |
| `release*` | 23 | Entire Releases domain (new product area) |
| `project*` | 22 | Labels, relations, unarchive, move milestone, status archive |
| `customer*` | 18 | Entire Customers / Customer Requests domain |
| `initiative*` | 18 | Labels, relations, archive, initiative updates |
| `organization*` | 14 | Admin: invites, domains, update, trial |
| `agent*` | 13 | Agent sessions / activities / skills |
| `attachment*` | 11 | Provider-specific smart links (Slack, GitHub, Jira, …) |
| `user*` | 11 | Role change, suspend, revoke sessions, settings |
| `email*` | 8 | Email intake addresses |
| `notification*` | 8 | Bulk read/unread/snooze, subscriptions |
| `roadmap*` | 8 | Entire Roadmaps domain |
| `gitAutomation*`, `team*`, `oauth*`, `timeSchedule*`, `webhook*`, `template*`, `triageResponsibility*`, `viewPreferences*`, `workflowState*`, others | ~70 | |

---

## 3. Missing first-class entities

These are product entities Linear exposes with full CRUD that `linear-cli` has
**no command for at all**:

| Entity | Docs | Key mutations | Priority |
|---|---|---|---|
| **Customers / Customer Requests** | `docs/customer-requests.md` | `customerCreate/Update/Delete/Merge`, `customerNeed*`, `customerStatus*`, `customerTier*`, `customerUpsert/Sync` | **P1** |
| **Releases** | `docs/releases.md` | `releaseCreate/Update/Delete/Archive/Complete`, `releaseNote*`, `releasePipeline*`, `releaseStage*`, `issueToRelease*` | **P1** |
| **Roadmaps** | (query `roadmaps`, schema) | `roadmapCreate/Update/Delete/Archive`, `roadmapToProject*` | **P2** |
| **Project labels** | `docs/project-labels.md` | `projectLabelCreate/Update/Delete/Restore/Retire`, `projectAddLabel/RemoveLabel` | **P1** |
| **Initiative labels** | (schema) | `initiativeLabelCreate/Update/Delete/Restore/Retire`, `initiativeAddLabel/RemoveLabel` | **P1** |
| **Project relations** | (schema) | `projectRelationCreate/Update/Delete`, `projectRelation(s)` query | **P1** |
| **Initiative relations** | (schema) | `initiativeRelationCreate/Update/Delete` | **P2** |
| **Project status updates (delete/unarchive)** | `docs/initiative-and-project-updates.md` | `projectUpdateDelete/Unarchive` | **P2** |
| **Templates** | (query `templates`, `templateSearch`) | `templateCreate/Update/Delete` | **P2** |
| **Webhooks** | `docs/api-and-webhooks.md` | `webhookCreate/Update/Delete/RotateSecret`, `webhooks` query | **P2** |
| **Time schedules** | `docs/saml-and-access-control.md` | `timeScheduleCreate/Update/Delete/UpsertExternal/Refresh` | **P3** |
| **Triage responsibilities** | `docs/triage.md` | `triageResponsibilityCreate/Update/Delete` | **P3** |
| **Workflow state CRUD** | `docs/teams.md` | `workflowStateCreate/Update/Archive` (today: read-only `team states`) | **P1** |
| **Emojis** | (schema) | `emojiCreate/Delete`, `emojis` query | **P3** |
| **Reactions** | `docs/comment-on-issues.md` | `reactionCreate/Delete` | **P3** |
| **Entity external links** | (schema) | `entityExternalLinkCreate/Update/Delete` | **P3** |
| **Agent sessions / skills / activities** | `docs/coding-sessions.md`, `docs/agents-in-linear.md`, `docs/linear-agent.md` | `agentSessionCreateOnIssue/OnComment`, `agentSessionUpdate/Restart`, `agentActivity*`, `agentSkill*` | **P2** |
| **Push subscriptions** | (schema) | `pushSubscriptionCreate/Delete` | **P3** |
| **View preferences** | (schema) | `viewPreferencesCreate/Update/Delete` | **P3** |
| **Git automation** | `docs/github-integration.md` | `gitAutomationState*`, `gitAutomationTargetBranch*` | **P3** |
| **Integrations (connect)** | `docs/*` integrations | `integration*Connect` (68 ops) | **P4** |

---

## 4. Partial entities — gaps on things already supported

These entities have commands already; the listed operations are missing.

**Issues**
- Label attach/detach on an existing issue: `issueAddLabel`, `issueRemoveLabel`
  (create/update accept labels, but there is no standalone command).
- Bulk: `issueBatchCreate`, `issueBatchUpdate`.
- Collaboration: `issueSubscribe`, `issueUnsubscribe`, `issueShare`,
  `issueUnshare`.
- Reminders: `issueReminder`, `issueRemoveReminder`.
- Lifecycle: `issueUnarchive` (only `archive` exists today),
  `issueDelete(permanentlyDelete: true)`.
- Label lifecycle: `issueLabelRestore`, `issueLabelRetire`.
- Imports: `issueImportCreateAsana/CSVJira/Clubhouse/Github/Jira/LinearV2`,
  `issueImportProcess/Delete/Update` (docs: `import-issues.md`, `*-to-linear.md`).
- Releases link: `issueToReleaseCreate/Delete`.
- Queries: `searchIssues`, `searchProjects`, `semanticSearch`,
  `issueVcsBranchSearch`, `issueFilterSuggestion`,
  `issueTitleSuggestionFromCustomerRequest`.

**Projects**
- `projectUnarchive` (only `archive` exists), `projectDelete(trash)`.
- `projectMilestoneMove`, `projectRelation*`, `projectAddLabel/RemoveLabel`,
  `projectReassignStatus`.
- `projectStatusArchive`, `projectStatusUnarchive` (create/update/delete exist).
- `projectUpdateDelete`, `projectUpdateUnarchive`.
- `projectCreateSlackChannel`, `projectExternalSyncDisable`.

**Initiatives**
- `initiativeArchive` / `initiativeUnarchive` (only `delete` exists).
- `initiativeAddLabel` / `initiativeRemoveLabel`, `initiativeLabel*`.
- `initiativeRelation*`, `initiativeToProjectUpdate`, `initiativeLeadTeamUpdate`.
- Initiative status updates: `initiativeUpdateCreate/Update/Archive/Unarchive`.

**Cycles**
- `cycleShiftAll`, `cycleStartUpcomingCycleToday`, `teamCyclesDelete`.

**Documents**
- `documentUnarchive`; `documentContentHistory*` queries.

**Comments**
- `commentResolve`, `commentUnresolve` (threaded comment resolution).

**Teams**
- `teamMembershipCreate/Update/Delete`, `teamUnarchive`, `teamKeyDelete`,
  `teamCyclesDelete`.

**Users**
- `userChangeRole`, `userSuspend`, `userUnsuspend`, `userRevokeSession`,
  `userRevokeAllSessions`, `userSettingsUpdate`, `userFlagUpdate`,
  `authenticationSessions` / `userSessions` queries.

**Notifications / Inbox**
- `notificationArchiveAll`, `notificationMarkReadAll` (bulk input form),
  `notificationMarkUnreadAll`, `notificationSnoozeAll`,
  `notificationUnsnoozeAll`, `notificationSubscription{Create,Update,Delete}`,
  `notificationsUnreadCount` query.

**Attachments**
- Provider-specific smart links: `attachmentLinkSlack`, `attachmentLinkGitHubPR`,
  `attachmentLinkGitHubIssue`, `attachmentLinkGitLabMR`,
  `attachmentLinkJiraIssue`, `attachmentLinkZendesk`, `attachmentLinkDiscord`,
  `attachmentLinkFront`, `attachmentLinkIntercom`, `attachmentLinkSalesforce`,
  `attachmentSyncToSlack`, `fileUploadDangerouslyDelete`.

---

## 5. What changed in Linear since the CLI's baseline

Features that ship in today's docs/schema but postdate the CLI's last commit
(`2026-03-04`) — or were never covered:

| New area | Docs | Why it matters |
|---|---|---|
| **Coding sessions & Reviews** | `docs/coding-sessions.md`, `docs/diffs.md`, `docs/releases.md` | Linear's agentic dev workflow; new entities + `agentSession*` API |
| **Releases / pipelines / stages** | `docs/releases.md` | 23 mutations, zero CLI coverage |
| **Customers & Customer Requests** | `docs/customer-requests.md` | 18 mutations; intake workflow |
| **Triage Intelligence** | `docs/triage-intelligence.md` | AI triage on top of `triage.md` |
| **Linear Agent / AI Agents / Loops** | `docs/linear-agent.md`, `docs/agents-in-linear.md`, `docs/loops.md` | Agent platform |
| **MCP server** | `docs/mcp.md`, `docs/connect-mcp-servers.md` | MCP integration surface |
| **Code Intelligence** | `docs/code-intelligence.md` | New |
| **Dashboards / Insights / AI Credits** | `docs/dashboards.md`, `docs/insights.md`, `docs/ai-credits.md` | Reporting + billing |
| **Sub-initiatives / sub-teams** | `docs/sub-initiatives.md`, `docs/sub-teams.md` | Hierarchy the CLI doesn't model |
| **Project & initiative labels, relations** | `docs/project-labels.md`, schema | New label/relation types |
| **Project milestones move / status archive** | schema | Lifecycle ops |

The CLI's own `master_api_ref.md` predates most of the above and is now stale.

---

## 6. Plan to close the gap

Principle: extend breadth in **thin, consistent slices** — one entity per PR,
CRUD + a couple of query subcommands, `--json/-p` output, and integration
coverage in `crud_test.go`. Keep the `graphql` escape hatch for the long tail
(especially integrations).

### Phase 0 — high-leverage lifecycle fixes (small, ~days)

Mostly one-liners on existing entities; immediately useful.

1. Unarchive everywhere: `issue unarchive`, `project unarchive`,
   `initiative unarchive`, `document unarchive`.
2. `comment resolve` / `comment unresolve`.
3. Issue label attach/detach: `issue label add REMOVE` →
   `issueAddLabel`/`issueRemoveLabel`; mirror `project label add/remove`,
   `initiative label add/remove`.
4. Notification bulk: `inbox mark-all-read`, `mark-all-unread`, `snooze-all`,
   `unsnooze-all`; add `--all` archive.
5. Cycle ops: `cycle shift-all`, `cycle start-today`.
6. Project: `project milestone move`, `project status archive/unarchive`.

### Phase 1 — labels, relations, workflow states, project updates (medium)

7. `project-label` CRUD (mirror existing `label` command) →
   `projectLabelCreate/Update/Delete/Restore/Retire`.
8. `initiative-label` CRUD → `initiativeLabel*`.
9. `project relation` CRUD → `projectRelation*` (+ `project relations` query).
10. `team state create/update/archive` → `workflowState*` (today read-only).
11. `project status delete/unarchive`; `project update delete/unarchive`.
12. `projectUpdateDelete` + `projectUpdateUnarchive`.

### Phase 2 — new first-class domains (larger, one PR each)

13. **Customers**: `customer list/get/create/update/delete/merge`,
    `customer need …`, `customer status/tier …` →
    `customer*`, `customerNeed*`, `customerStatus*`, `customerTier*`.
14. **Releases**: `release list/get/create/update/archive/complete`, `release note …`,
    `release pipeline …`, `release stage …`, `release add-issue` →
    `release*`, `releaseNote*`, `releasePipeline*`, `releaseStage*`,
    `issueToRelease*`.
15. **Roadmaps**: `roadmap list/get/create/update/archive` + link/unlink project.
16. **Templates**: `template list/get/create/update/delete` + `issue create --template`.
17. **Webhooks**: `webhook list/get/create/update/delete/rotate-secret`.
18. **Imports**: `import jira|github|asana|shortcut|csv` → `issueImportCreate*`
    + `import process/status`; `export csv` → `createCsvExportReport`.
19. **Triage responsibilities**: `triage responsibility …`.
20. **Time schedules**: `schedule …` → `timeSchedule*`.

### Phase 3 — agent platform & admin (as demand dictates)

21. **Agent sessions**: `agent session create/update/restart`, `agent activity …`,
    `agent skill …` → `agentSession*`, `agentActivity*`, `agentSkill*`.
22. **User admin**: `user role`, `user suspend/unsuspend`, `user revoke-sessions`,
    `user settings`.
23. **Organization admin**: `org invite …`, `org domain …`, `org update`.
24. **Reactions / emojis / entity external links**.
25. **Push subscriptions / view preferences**.
26. **Integrations**: leave to `graphql` (OAuth flows don't fit a CLI), revisit
    only for read-only settings.

### Cross-cutting work

- **Search parity**: add `searchIssues`, `searchProjects`, `semanticSearch`
  commands (the `search` query family is richer than today's `issue search`).
- **Schema/command consistency**: add an `--include-archived` and pagination
  flag audit across list commands; standardize on the existing `-j/-p` outputs.
- **Generated coverage check**: script the §2 numbers (schema vs. `pkg/api`
  operations) and fail CI when new mutations appear, so drift is visible.
  Source: re-run `make scrape-docs`, then diff `linear-source-docs/manifest.json`.

### Suggested sequencing

| Phase | Scope | Rough effort |
|---|---|---|
| 0 | Lifecycle fixes on existing entities | 3–5 PRs, days |
| 1 | Labels / relations / workflow states | 5–7 PRs, ~1 week |
| 2 | Customers, Releases, Roadmaps, Templates, Webhooks, Imports | 6–8 PRs, ~2–3 weeks |
| 3 | Agents + admin + misc | open-ended, as needed |

Every slice should land with: command wiring in `cmd/`, client method in
`pkg/api/queries.go`, README section, and a `crud_test.go` case.

---

## 7. Verifying progress

```bash
# Refresh the source corpus (also refreshes the SDL schema)
make scrape-docs

# Compare against the previous corpus to see what Linear changed
git diff docs/linear-source-docs/manifest.json

# Re-measure surface coverage
python3 - <<'PY'
import re, pathlib
s = pathlib.Path("docs/linear-source-docs/external/raw.githubusercontent.com/"
                 "linear/linear/master/packages/sdk/src/schema.graphql").read_text()
block = s[re.search(r'^type Mutation [^{]*\{', s, re.M).end():]
depth, i = 1, 0
while depth and i < len(block):
    depth += block[i] == '{'
    depth -= block[i] == '}'
    i += 1
fields = re.findall(r'^  ([a-z][A-Za-z0-9_]*)\s*[(:]', block[:i], re.M)
print("mutations in schema:", len(set(fields)))
PY
```

Then grep `pkg/api/` for the operations you added:

```bash
grep -ohE '\b(issue|project|customer|release|roadmap)[A-Z][A-Za-z]*\b' pkg/api/*.go | sort -u
```
