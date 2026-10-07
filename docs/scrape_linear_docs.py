#!/usr/bin/env python3
"""
scrape_linear_docs.py — recursively scrape Linear's official documentation
into a local, greppable corpus of Markdown source files.

Why this exists
---------------
Linear publishes an LLM-friendly index at https://linear.app/llms.txt whose
entries point at plain-Markdown (``.md``) renderings of the live docs. Those
``.md`` pages themselves link to further ``.md`` pages, so we crawl
*breadth-first* starting from the index and follow every Markdown link we
find on ``linear.app`` until the whole reachable Markdown graph is fetched.

That gives us a reproducible corpus we can diff over time to detect what
Linear has changed, and to measure feature parity for this CLI.

Usage
-----
    python3 docs/scrape_linear_docs.py                 # scrape into docs/linear-source-docs
    python3 docs/scrape_linear_docs.py --out DIR       # custom output directory
    python3 docs/scrape_linear_docs.py --dry-run       # show plan, fetch nothing
    python3 docs/scrape_linear_docs.py --jobs 8        # parallel fetches (default 8)
    python3 docs/scrape_linear_docs.py --max-depth 6   # recursion depth (default 6)
    python3 docs/scrape_linear_docs.py --prune         # delete stale corpus files

The script is idempotent: re-running it overwrites the corpus in place and
regenerates the manifest. That is the "regenerate from live Linear" workflow.

Stdlib-only: no third-party dependencies.
"""

from __future__ import annotations

import argparse
import concurrent.futures
import hashlib
import json
import os
import re
import sys
import time
from datetime import datetime, timezone
from urllib.error import HTTPError, URLError
from urllib.parse import urldefrag, urljoin, urlparse, urlunparse
from urllib.request import Request, urlopen

DEFAULT_INDEX = "https://linear.app/llms.txt"
DEFAULT_OUT = os.path.join(
    os.path.dirname(os.path.abspath(__file__)), "linear-source-docs"
)
USER_AGENT = (
    "linear-cli-docs-scraper/1.0 (+https://github.com/roboalchemist/linear-cli)"
)

# We follow Markdown links anywhere on linear.app. Assets (like the GraphQL
# SDL schema) are downloaded but not recursed into.
FOLLOW_EXT = (".md",)
ASSET_EXT = (".graphql", ".json", ".txt")
ALLOWED_HOSTS = ("linear.app",)

# [text](url) — url may be absolute or relative.
LINK_RE = re.compile(r"\[[^\]]*\]\(([^)\s]+)\)")
# ### Section headers in the llms.txt index.
SECTION_RE = re.compile(r"^###\s+(.*)$")


# ---------------------------------------------------------------------------
# HTTP
# ---------------------------------------------------------------------------
def fetch(url: str, timeout: int = 30, retries: int = 3) -> bytes:
    """GET a URL with retries. Raises on final failure.

    Client errors (4xx) are raised immediately — retrying a 404 is pointless
    and the crawler uses the status code to skip non-existent `.md` variants.
    """
    last_err: Exception | None = None
    for attempt in range(1, retries + 1):
        try:
            req = Request(
                url,
                headers={
                    "User-Agent": USER_AGENT,
                    "Accept": "text/markdown, text/plain, application/json, */*",
                },
            )
            with urlopen(req, timeout=timeout) as resp:
                return resp.read()
        except HTTPError as err:
            if 400 <= err.code < 500:
                raise
            last_err = err
            if attempt < retries:
                time.sleep(min(2 ** attempt, 8))
        except (URLError, TimeoutError, OSError) as err:
            last_err = err
            if attempt < retries:
                time.sleep(min(2 ** attempt, 8))
    raise last_err  # type: ignore[misc]


# ---------------------------------------------------------------------------
# URL helpers
# ---------------------------------------------------------------------------
def normalize(url: str, base: str) -> str | None:
    """Resolve ``url`` against ``base``, drop the fragment, return None if
    the result is not a linear.app http(s) URL."""
    abs_url = urljoin(base, url.strip())
    abs_url, _ = urldefrag(abs_url)
    parsed = urlparse(abs_url)
    if parsed.scheme not in ("http", "https"):
        return None
    if parsed.netloc not in ALLOWED_HOSTS:
        return None
    if parsed.scheme == "http":
        abs_url = "https://" + abs_url[len("http://"):]
    return abs_url


def is_asset(url: str) -> bool:
    path = urlparse(url).path.lower()
    return any(path.endswith(ext) for ext in ASSET_EXT)


def local_path_for(url: str, out_dir: str) -> str:
    """Map a URL to a file under out_dir, preserving source structure."""
    parsed = urlparse(url)
    path = parsed.path
    if not path or path.endswith("/"):
        path = path + "index"
    if parsed.netloc == "linear.app":
        rel = path.lstrip("/")
    else:
        rel = os.path.join("external", parsed.netloc, path.lstrip("/"))
    return os.path.join(out_dir, rel)


def candidate_links(text: str, base: str) -> list[tuple[str, bool]]:
    """Return (url, derived) candidates found in ``text``.

    ``derived`` is True when the link pointed at an extension-less doc page
    (e.g. ``/docs/projects``) and we synthesised its Markdown variant
    (``/docs/projects.md``). Derived candidates that 404 are skipped, since
    not every extension-less link has a Markdown twin.
    """
    found: dict[str, bool] = {}
    for raw in LINK_RE.findall(text):
        abs_url = urljoin(base, raw.strip())
        parsed = urlparse(abs_url)
        if parsed.scheme not in ("http", "https"):
            continue
        if parsed.netloc not in ALLOWED_HOSTS:
            continue
        if any(ch in abs_url for ch in "<>{}") or " " in abs_url:
            continue  # template placeholder, e.g. <preview-url>
        if parsed.scheme == "http":
            abs_url = "https://" + abs_url[len("http://"):]
            parsed = urlparse(abs_url)
        path = parsed.path
        if any(path.lower().endswith(ext) for ext in FOLLOW_EXT):
            canonical = urlunparse(
                (parsed.scheme, parsed.netloc, path, "", "", "")
            )
            found.setdefault(canonical, False)
        elif "." not in path.rsplit("/", 1)[-1] and (
            path.startswith("/docs/") or path.startswith("/developers/")
        ):
            canonical = urlunparse(
                (parsed.scheme, parsed.netloc, path + ".md", "", "", "")
            )
            found.setdefault(canonical, True)
    return list(found.items())


# ---------------------------------------------------------------------------
# Index parsing (for section + title metadata)
# ---------------------------------------------------------------------------
def parse_index(text: str, base: str = DEFAULT_INDEX) -> dict[str, dict[str, str]]:
    """Parse llms.txt -> {url: {section, title}}."""
    meta: dict[str, dict[str, str]] = {}
    section = ""
    for line in text.splitlines():
        m = SECTION_RE.match(line)
        if m:
            section = m.group(1).strip()
            continue
        for lg in re.finditer(r"\[([^\]]*)\]\(([^)\s]+)\)", line):
            title, raw = lg.group(1), lg.group(2)
            norm = normalize(raw, base)
            if norm:
                meta.setdefault(norm, {"section": section, "title": title})
    return meta


# ---------------------------------------------------------------------------
# Crawl
# ---------------------------------------------------------------------------
def crawl(
    index_url: str,
    out_dir: str,
    jobs: int,
    max_depth: int,
    dry_run: bool,
) -> dict:
    started = datetime.now(timezone.utc)
    index_bytes = fetch(index_url)
    index_text = index_bytes.decode("utf-8", "replace")
    meta = parse_index(index_text, index_url)

    # Assets referenced by the index (e.g. the GraphQL SDL) are downloaded
    # once, but never recursed into (and may live on a non-linear.app host).
    assets: list[str] = []
    for raw in LINK_RE.findall(index_text):
        abs_url, _ = urldefrag(urljoin(index_url, raw.strip()))
        if urlparse(abs_url).scheme in ("http", "https") and is_asset(abs_url):
            if abs_url not in assets:
                assets.append(abs_url)

    # --- Breadth-first crawl of the Markdown graph ----------------------
    visited: set[str] = set()
    records: dict[str, dict] = {}
    levels: list[list[str]] = []
    frontier: list[tuple[str, bool]] = candidate_links(index_text, index_url)

    if dry_run:
        planned = sorted({u for u, _ in frontier})
        print(f"==> Dry run: {len(planned)} seed documents from the index",
              file=sys.stderr)
        return {
            "index_url": index_url,
            "started": started.isoformat(),
            "dry_run": True,
            "planned": planned + [a for a in assets if a not in planned],
            "assets": assets,
        }

    os.makedirs(out_dir, exist_ok=True)

    depth = 0
    while frontier and depth <= max_depth:
        wave: list[tuple[str, bool]] = []
        seen: set[str] = set()
        for url, derived in frontier:
            if url in visited or url in seen:
                continue
            seen.add(url)
            wave.append((url, derived))
        if not wave:
            break
        visited.update(url for url, _ in wave)
        levels.append([u for u, _ in wave])
        print(f"    depth {depth}: {len(wave)} page(s)", file=sys.stderr)

        next_frontier: list[tuple[str, bool]] = []
        with concurrent.futures.ThreadPoolExecutor(max_workers=jobs) as pool:
            futs = {pool.submit(fetch, u): (u, derived) for u, derived in wave}
            for fut in concurrent.futures.as_completed(futs):
                url, derived = futs[fut]
                try:
                    body = fut.result()
                except HTTPError as err:
                    if derived and err.code == 404:
                        continue  # no Markdown twin for this link
                    records[url] = {
                        "status": "fail",
                        "error": f"HTTP {err.code}",
                        "path": local_path_for(url, out_dir),
                    }
                    continue
                except Exception as err:  # noqa: BLE001
                    records[url] = {
                        "status": "fail",
                        "error": str(err),
                        "path": local_path_for(url, out_dir),
                    }
                    continue
                dest = local_path_for(url, out_dir)
                os.makedirs(os.path.dirname(dest), exist_ok=True)
                tmp = dest + ".tmp"
                with open(tmp, "wb") as fh:
                    fh.write(body)
                os.replace(tmp, dest)
                records[url] = {
                    "status": "ok",
                    "path": dest,
                    "bytes": len(body),
                    "sha256": hashlib.sha256(body).hexdigest(),
                }
                for cand, is_derived in candidate_links(
                    body.decode("utf-8", "replace"), url
                ):
                    if cand not in visited:
                        next_frontier.append((cand, is_derived))
        frontier = sorted(set(next_frontier))
        depth += 1

    # --- Download index-referenced assets (non-recursive) ---------------
    if assets:
        print(f"    assets: {len(assets)} file(s)", file=sys.stderr)
        with concurrent.futures.ThreadPoolExecutor(max_workers=jobs) as pool:
            for url, info in pool.map(_grab_asset(out_dir), assets):
                records[url] = info

    # --- Persist the raw index as part of the corpus --------------------
    with open(os.path.join(out_dir, "llms.txt"), "wb") as fh:
        fh.write(index_bytes)

    finished = datetime.now(timezone.utc)
    ok = sum(1 for r in records.values() if r["status"] == "ok")
    failed = len(records) - ok
    return {
        "index_url": index_url,
        "started": started.isoformat(),
        "finished": finished.isoformat(),
        "dry_run": False,
        "counts": {"ok": ok, "failed": failed, "total": len(records)},
        "results": records,
        "meta": meta,
        "levels": levels,
        "assets": assets,
        "out_dir": out_dir,
    }


def _grab_asset(out_dir: str):
    """Return a function that downloads a non-recursive asset URL."""

    def _grab(url: str) -> tuple[str, dict]:
        dest = local_path_for(url, out_dir)
        try:
            body = fetch(url)
        except Exception as err:  # noqa: BLE001
            return url, {"status": "fail", "error": str(err), "path": dest}
        os.makedirs(os.path.dirname(dest), exist_ok=True)
        tmp = dest + ".tmp"
        with open(tmp, "wb") as fh:
            fh.write(body)
        os.replace(tmp, dest)
        return url, {
            "status": "ok",
            "path": dest,
            "bytes": len(body),
            "sha256": hashlib.sha256(body).hexdigest(),
        }

    return _grab


# ---------------------------------------------------------------------------
# Manifest rendering
# ---------------------------------------------------------------------------
def write_manifests(report: dict, out_dir: str) -> None:
    meta = report.get("meta", {})
    results = report["results"]
    rel = lambda p: os.path.relpath(p, out_dir)  # noqa: E731

    # machine-readable JSON
    manifest = []
    for url in sorted(results):
        info = results[url]
        m = meta.get(url, {})
        manifest.append(
            {
                "url": url,
                "title": m.get("title", ""),
                "section": m.get("section", "Discovered"),
                "local_path": rel(info["path"]),
                "status": info["status"],
                "bytes": info.get("bytes", 0),
                "sha256": info.get("sha256", ""),
                "error": info.get("error", ""),
            }
        )
    payload = {
        "index_url": report["index_url"],
        "generated_at": report["finished"],
        "counts": report["counts"],
        "documents": manifest,
    }
    with open(os.path.join(out_dir, "manifest.json"), "w") as fh:
        json.dump(payload, fh, indent=2, sort_keys=False)
        fh.write("\n")

    # TSV for shell tooling
    with open(os.path.join(out_dir, "manifest.tsv"), "w") as fh:
        fh.write("section\ttitle\turl\tlocal_path\tbytes\tsha256\tstatus\n")
        for d in manifest:
            row = [
                d["section"],
                d["title"],
                d["url"],
                d["local_path"],
                str(d["bytes"]),
                d["sha256"],
                d["status"],
            ]
            fh.write("\t".join(c.replace("\t", " ") for c in row) + "\n")

    # Human-readable Markdown
    counts = report["counts"]
    lines = [
        "# Linear Source Docs — Manifest",
        "",
        f"Auto-generated by `docs/scrape_linear_docs.py` on {report['finished']}.",
        "",
        f"- Index: <{report['index_url']}>",
        f"- Documents: **{counts['ok']}** fetched, {counts['failed']} failed "
        f"({counts['total']} total)",
        "- Machine-readable: [`manifest.json`](manifest.json) · [`manifest.tsv`](manifest.tsv)",
        "",
        "| Section | Title | Source | Local file | Bytes | Status |",
        "|---------|-------|--------|-----------|------:|--------|",
    ]
    for d in manifest:
        title = d["title"] or "—"
        lines.append(
            f"| {d['section']} | {title} | [link]({d['url']}) | "
            f"[`{d['local_path']}`]({d['local_path']}) | {d['bytes']} | {d['status']} |"
        )
    lines.append("")
    with open(os.path.join(out_dir, "MANIFEST.md"), "w") as fh:
        fh.write("\n".join(lines))


def prune_stale(out_dir: str, keep: set[str]) -> list[str]:
    """Delete files under out_dir that are not part of the current corpus.

    Returns the list of removed paths. Empty directories left behind are
    removed too.
    """
    removed: list[str] = []
    for root, _dirs, files in os.walk(out_dir):
        for name in files:
            full = os.path.join(root, name)
            if full not in keep:
                os.remove(full)
                removed.append(full)
    # clean up now-empty directories, deepest first
    for root, _dirs, _files in sorted(
        os.walk(out_dir), key=lambda t: -len(t[0])
    ):
        if root != out_dir and not os.listdir(root):
            os.rmdir(root)
    return removed


# ---------------------------------------------------------------------------
# Entry point
# ---------------------------------------------------------------------------
def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        description="Recursively scrape Linear's Markdown documentation."
    )
    parser.add_argument("--index", default=DEFAULT_INDEX, help="index URL")
    parser.add_argument("--out", default=DEFAULT_OUT, help="output directory")
    parser.add_argument("--jobs", type=int, default=8, help="parallel fetches")
    parser.add_argument(
        "--max-depth", type=int, default=6, help="max recursion depth"
    )
    parser.add_argument(
        "--dry-run", action="store_true", help="show plan, fetch nothing"
    )
    parser.add_argument(
        "--prune",
        action="store_true",
        help="delete corpus files not present in this run (stale pages)",
    )
    args = parser.parse_args(argv)

    out_dir = os.path.abspath(os.path.expanduser(args.out))
    print(f"==> Index:  {args.index}", file=sys.stderr)
    print(f"==> Output: {out_dir}", file=sys.stderr)

    report = crawl(
        index_url=args.index,
        out_dir=out_dir,
        jobs=max(1, args.jobs),
        max_depth=max(0, args.max_depth),
        dry_run=args.dry_run,
    )

    if report.get("dry_run"):
        print(f"==> Dry run: {len(report['planned'])} documents would be fetched",
              file=sys.stderr)
        for url in report["planned"]:
            print(f"    {url}", file=sys.stderr)
        return 0

    write_manifests(report, out_dir)
    counts = report["counts"]

    if args.prune:
        keep = {info["path"] for info in report["results"].values()}
        keep.update(
            os.path.join(out_dir, name)
            for name in ("llms.txt", "manifest.json", "manifest.tsv", "MANIFEST.md")
        )
        removed = prune_stale(out_dir, keep)
        if removed:
            print(f"==> Pruned {len(removed)} stale file(s)", file=sys.stderr)

    print("", file=sys.stderr)
    print("==> Done.", file=sys.stderr)
    print(f"    corpus:    {out_dir}", file=sys.stderr)
    print(f"    index:     {os.path.join(out_dir, 'llms.txt')}", file=sys.stderr)
    print(f"    manifest:  {os.path.join(out_dir, 'manifest.json')}", file=sys.stderr)
    print(f"    fetched:   {counts['ok']} ok, {counts['failed']} failed", file=sys.stderr)
    if counts["failed"]:
        for url, info in sorted(report["results"].items()):
            if info["status"] != "ok":
                print(f"      ! {url}: {info.get('error', '')}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
