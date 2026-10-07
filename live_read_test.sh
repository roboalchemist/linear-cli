#!/usr/bin/env bash
#
# live_read_test.sh — automated LIVE read-only test suite for linear-cli.
#
# Exercises every read/list/get/search command against the real Linear API and
# validates the output. WRITE commands are never invoked.
#
# Safety & etiquette:
#   * Only read operations are run (list/get/search/…).
#   * Small --limit values and a delay between calls keep us well under Linear's
#     rate limit (~1,500 requests/hour). The suite issues ~40 requests.
#   * Transient failures (rate limit / timeout) are retried with backoff.
#
# Usage:
#   ./live_read_test.sh                 # build + run everything
#   BIN=/path/to/linear-cli ./live_read_test.sh
#   ./live_read_test.sh --report out.json --fast
#
# Exit code is 0 only when there are no unexpected failures. Commands that are
# legitimately unavailable to this workspace/key (feature-gated customers,
# admin-only webhooks/audit) are reported as SKIP, not FAIL.
set -uo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="${BIN:-$REPO_ROOT/linear-cli-livetest}"
DELAY="${DELAY:-0.35}"
RETRIES="${RETRIES:-3}"
REPORT=""
FAST=0

while [ $# -gt 0 ]; do
  case "$1" in
    --report) REPORT="$2"; shift 2 ;;
    --fast) FAST=1; DELAY=0.1; shift ;;
    --bin) BIN="$2"; shift 2 ;;
    -h|--help) sed -n '2,20p' "${BASH_SOURCE[0]}" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) echo "unknown option: $1" >&2; exit 2 ;;
  esac
done

# Colors
GREEN=$'\033[0;32m'; RED=$'\033[0;31m'; YELLOW=$'\033[0;33m'; CYAN=$'\033[0;36m'; NC=$'\033[0m'

PASS=0; FAIL=0; SKIP=0
declare -a RESULTS

command -v jq >/dev/null 2>&1 && HAVE_JQ=1 || HAVE_JQ=0

log() { printf '%s\n' "$*"; }

# ---------------------------------------------------------------------------
# Build (unless a binary was provided and exists)
# ---------------------------------------------------------------------------
if [ ! -x "$BIN" ]; then
  log "Building linear-cli..."
  ( cd "$REPO_ROOT" && go build -o "$BIN" . ) || { echo "build failed" >&2; exit 1; }
fi

if [ -z "${LINEAR_API_KEY:-}" ]; then
  log "${YELLOW}Warning:${NC} LINEAR_API_KEY is not set; relying on stored auth."
fi

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

# json_valid <string> -> 0 if parses as JSON
json_valid() {
  if [ "$HAVE_JQ" -eq 1 ]; then printf '%s' "$1" | jq -e . >/dev/null 2>&1
  else printf '%s' "$1" | python3 -c 'import sys,json; json.load(sys.stdin)' >/dev/null 2>&1; fi
}

# run_capture <cmd...>; sets RC, OUT
run_capture() {
  OUT="$("$@" 2>&1)"; RC=$?
}

# is_transient <output>
is_transient() {
  echo "$1" | grep -qiE 'rate ?limit|too many requests|context deadline|timeout|connection reset|502|503|429'
}

# is_unavailable <output>  (legitimately not usable with this workspace/key)
is_unavailable() {
  echo "$1" | grep -qiE 'not available in this workspace|owner or admin required|access denied|Invalid role'
}

# record <label> <status> <detail>
record() {
  RESULTS+=("$1|$2|$3")
}

# check <label> <command...>
check() {
  local label="$1"; shift
  local attempt=1 out rc
  while :; do
    run_capture "$@"
    out="$OUT"; rc="$RC"
    if [ $rc -eq 0 ]; then break; fi
    if is_transient "$out" && [ $attempt -lt "$RETRIES" ]; then
      sleep $((attempt * 2)); attempt=$((attempt + 1)); continue
    fi
    break
  done

  if [ $rc -eq 0 ]; then
    PASS=$((PASS + 1)); printf '  %-6s %s\n' "${GREEN}PASS${NC}" "$label"; record "$label" PASS ""
  elif is_unavailable "$out"; then
    SKIP=$((SKIP + 1)); printf '  %-6s %s\n' "${YELLOW}SKIP${NC}" "$label  ($(echo "$out" | head -1 | cut -c1-60))"; record "$label" SKIP "$(echo "$out" | head -1)"
  else
    FAIL=$((FAIL + 1)); printf '  %-6s %s\n' "${RED}FAIL${NC}" "$label"; echo "        $out" | head -3; record "$label" FAIL "$(echo "$out" | head -1)"
  fi
  sleep "$DELAY"
}

# check_json <label> <command...>  -> must exit 0 AND emit valid JSON
check_json() {
  local label="$1"; shift
  local attempt=1 out rc
  while :; do
    run_capture "$@"
    out="$OUT"; rc="$RC"
    if [ $rc -eq 0 ]; then break; fi
    if is_transient "$out" && [ $attempt -lt "$RETRIES" ]; then
      sleep $((attempt * 2)); attempt=$((attempt + 1)); continue
    fi
    break
  done

  if [ $rc -ne 0 ]; then
    if is_unavailable "$out"; then
      SKIP=$((SKIP + 1)); printf '  %-6s %s\n' "${YELLOW}SKIP${NC}" "$label"; record "$label" SKIP "$(echo "$out" | head -1)"
    else
      FAIL=$((FAIL + 1)); printf '  %-6s %s\n' "${RED}FAIL${NC}" "$label"; echo "        $out" | head -3; record "$label" FAIL "$(echo "$out" | head -1)"
    fi
  elif json_valid "$out"; then
    PASS=$((PASS + 1)); printf '  %-6s %s\n' "${GREEN}PASS${NC}" "$label (json)"; record "$label" PASS ""
  else
    FAIL=$((FAIL + 1)); printf '  %-6s %s\n' "${RED}FAIL${NC}" "$label (invalid JSON)"; record "$label" FAIL "invalid json"
  fi
  sleep "$DELAY"
}

# first_id <list-command...> -> prints first "id" value found in JSON output
first_id() {
  local out
  out="$("$@" -j -l 1 2>/dev/null)"
  if [ "$HAVE_JQ" -eq 1 ]; then
    printf '%s' "$out" | jq -r '.. | objects | .id? // empty' 2>/dev/null | head -1
  else
    printf '%s' "$out" | python3 -c '
import sys,json
try: data=json.load(sys.stdin)
except Exception: sys.exit(0)
def walk(o):
    if isinstance(o,dict):
        if isinstance(o.get("id"),str): print(o["id"]); raise SystemExit
        for v in o.values(): walk(v)
    elif isinstance(o,list):
        for v in o: walk(v)
walk(data)' 2>/dev/null | head -1
  fi
}

section() { printf '\n%s== %s ==%s\n' "$CYAN" "$1" "$NC"; }

log "linear-cli live read-only test suite"
log "binary: $BIN"

# ---------------------------------------------------------------------------
# Read-only commands. Format: "<label>|<args>"
# ---------------------------------------------------------------------------
section "Auth / user"
check "auth status"      "$BIN" auth status
check "whoami"           "$BIN" whoami
check "user me"          "$BIN" user me
check "user list"        "$BIN" user list
check_json "user list json" "$BIN" user list -j -l 3
check "auth rate-limit"  "$BIN" auth rate-limit

section "Teams"
check "team list"        "$BIN" team list -l 3
check_json "team list json" "$BIN" team list -j -l 3
TEAM_KEY="$("$BIN" team list -j -l 1 2>/dev/null | { [ "$HAVE_JQ" -eq 1 ] && jq -r '.nodes[0].key // .[0].key // empty' || python3 -c 'import sys,json;d=json.load(sys.stdin);print((d.get("nodes") or d)[0].get("key",""))' 2>/dev/null; })"
[ -z "$TEAM_KEY" ] && TEAM_KEY="$("$BIN" team list 2>/dev/null | awk 'NR>1{print $1;exit}')"
if [ -n "$TEAM_KEY" ]; then
  check "team get $TEAM_KEY"      "$BIN" team get "$TEAM_KEY"
  check "team members $TEAM_KEY"  "$BIN" team members "$TEAM_KEY"
  check "team states $TEAM_KEY"   "$BIN" team states "$TEAM_KEY"
fi

section "Issues"
check "issue list"       "$BIN" issue list -l 3
check_json "issue list json" "$BIN" issue list -j -l 3
check "issue list filter" "$BIN" issue list -l 3 --assignee me
if [ -n "$TEAM_KEY" ]; then check "issue triage $TEAM_KEY" "$BIN" issue triage "$TEAM_KEY" -l 3; fi
ISSUE_ID="$("$BIN" issue list 2>/dev/null | grep -oE '[A-Z]+-[0-9]+' | head -1)"
if [ -n "$ISSUE_ID" ]; then
  check "issue get $ISSUE_ID"        "$BIN" issue get "$ISSUE_ID"
  check "issue search $ISSUE_ID"     "$BIN" issue search "$ISSUE_ID" -l 3
  check "issue activity $ISSUE_ID"   "$BIN" issue activity "$ISSUE_ID"
  check "issue comment list"         "$BIN" issue comment list "$ISSUE_ID" -l 3
  check "issue attachment list"      "$BIN" issue attachment list "$ISSUE_ID" -l 3
  check "issue relation list"        "$BIN" issue relation list "$ISSUE_ID"
fi

section "Projects"
check "project list"     "$BIN" project list -l 3
PID="$(first_id "$BIN" project list)"
if [ -n "$PID" ]; then
  check "project get"           "$BIN" project get "$PID"
  check "project issues"        "$BIN" project issues "$PID" -l 3
  check "project milestone list" "$BIN" project milestone list "$PID" -l 3
  check "project status list"   "$BIN" project status list "$PID" -l 3
fi
check "project label list"   "$BIN" project label list -l 3
check "project relation list" "$BIN" project relation list -l 3

section "Initiatives"
check "initiative list"  "$BIN" initiative list -l 3
IID="$(first_id "$BIN" initiative list)"
if [ -n "$IID" ]; then check "initiative projects" "$BIN" initiative projects "$IID" -l 3; fi
check "initiative label list"       "$BIN" initiative label list -l 3
check "initiative relation list"    "$BIN" initiative relation list -l 3
check "initiative update list"      "$BIN" initiative update list -l 3

section "Cycles / labels / views / favorites / documents"
check "cycle list"       "$BIN" cycle list -l 3
check "label list"       "$BIN" label list -l 3
check "view list"        "$BIN" view list -l 3
VID="$(first_id "$BIN" view list)"
if [ -n "$VID" ]; then check "view run" "$BIN" view run "$VID"; fi
check "favorite list"    "$BIN" favorite list
check "document list"    "$BIN" document list -l 3
check "document search"  "$BIN" document search test

section "Releases"
check "release list"          "$BIN" release list -l 3
check "release note list"     "$BIN" release note list -l 3
check "release pipeline list" "$BIN" release pipeline list -l 3
check "release stage list"    "$BIN" release stage list -l 3
RID="$(first_id "$BIN" release list)"
if [ -n "$RID" ]; then check "release get" "$BIN" release get "$RID"; fi

section "Customers (feature-gated in some workspaces)"
check "customer list"        "$BIN" customer list -l 3
check "customer status list" "$BIN" customer status list -l 3
check "customer tier list"   "$BIN" customer tier list -l 3
check "customer need list"   "$BIN" customer need list -l 3

section "Planning / admin"
check "roadmap list"                "$BIN" roadmap list -l 3
check "schedule list"               "$BIN" schedule list -l 3
check "triage responsibility list"  "$BIN" triage responsibility list -l 3
check "emoji list"                  "$BIN" emoji list -l 3
check "template list"               "$BIN" template list -l 3
check "template search"             "$BIN" template search bug -l 3
check "webhook list"                "$BIN" webhook list -l 3
check "audit log list"              "$BIN" audit log list -l 3

section "Platform"
check "agent session list"  "$BIN" agent session list -l 3
check "agent skill list"    "$BIN" agent skill list -l 3
check "integration list"    "$BIN" integration list -l 3
check "external user list"  "$BIN" external user list -l 3
check "organization get"    "$BIN" organization get
check "subscription list"   "$BIN" subscription list -l 3

section "Search"
check "search issues"    "$BIN" search issues login -l 3
check "search projects"  "$BIN" search projects command -l 3
check "search semantic"  "$BIN" search semantic onboarding -l 3

section "Inbox"
check "inbox list"       "$BIN" inbox

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
log ""
log "================================"
log "Live read-only test summary:"
log "  Passed: ${GREEN}$PASS${NC}"
log "  Skipped (unavailable): ${YELLOW}$SKIP${NC}"
log "  Failed: ${RED}$FAIL${NC}"
log "  Total:  $((PASS + SKIP + FAIL))"

if [ -n "$REPORT" ]; then
  {
    echo "{"
    echo "  \"generatedAt\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\","
    echo "  \"passed\": $PASS, \"skipped\": $SKIP, \"failed\": $FAIL,"
    echo "  \"results\": ["
    first=1
    for r in "${RESULTS[@]}"; do
      label="${r%%|*}"; rest="${r#*|}"; status="${rest%%|*}"; detail="${rest#*|}"
      [ $first -eq 0 ] && echo ","
      first=0
      printf '    {"name": %s, "status": %s, "detail": %s}' \
        "$(printf '%s' "$label" | python3 -c 'import sys,json;print(json.dumps(sys.stdin.read()))')" \
        "\"$status\"" \
        "$(printf '%s' "$detail" | python3 -c 'import sys,json;print(json.dumps(sys.stdin.read()))')"
    done
    echo ""
    echo "  ]"
    echo "}"
  } > "$REPORT"
  log "  Report: $REPORT"
fi

if [ "$FAIL" -eq 0 ]; then
  log ""
  log "${GREEN}✅ All live read-only tests passed${NC}"
  exit 0
fi
log ""
log "${RED}❌ $FAIL live read-only test(s) failed${NC}"
exit 1
