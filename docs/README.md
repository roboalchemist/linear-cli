# Documentation

This folder holds the tooling and the generated source material used to keep
`linear-cli` aligned with the live Linear product and API.

```
docs/
├── README.md                       # this file
├── FEATURE-GAP-ANALYSIS.md         # gap between linear-cli and live Linear + plan
├── scrape_linear_docs.py           # the auto-scraper
└── linear-source-docs/             # generated corpus (committed)
    ├── llms.txt                    # raw Linear LLM index
    ├── manifest.json               # machine-readable manifest
    ├── manifest.tsv                # same, for shell tooling
    ├── MANIFEST.md                 # human-readable manifest
    ├── docs/                       # linear.app/docs/*.md
    ├── developers/                 # linear.app/developers/*.md
    └── external/                   # e.g. the GraphQL SDL schema
```

## Auto-scraper: `scrape_linear_docs.py`

Linear publishes an LLM-friendly index at
<https://linear.app/llms.txt>. Its entries point at plain-Markdown (`.md`)
renderings of the live docs pages, and those pages link on to further `.md`
pages. The scraper crawls that Markdown graph breadth-first and writes every
page into `linear-source-docs/`, alongside a manifest.

It is **stdlib-only** (no third-party dependencies).

```bash
# Regenerate the corpus from live Linear (the normal workflow)
python3 docs/scrape_linear_docs.py

# Show what would be fetched, write nothing
python3 docs/scrape_linear_docs.py --dry-run

# Options
python3 docs/scrape_linear_docs.py --out DIR      # output directory
python3 docs/scrape_linear_docs.py --jobs 12      # parallel fetches (default 8)
python3 docs/scrape_linear_docs.py --max-depth 6  # recursion depth (default 6)
python3 docs/scrape_linear_docs.py --prune        # delete stale corpus files
python3 docs/scrape_linear_docs.py --index URL    # alternate index
```

Or via Make:

```bash
make scrape-docs
```

### What it crawls

1. Fetches `https://linear.app/llms.txt`.
2. Extracts every `.md` link (also `.graphql`/`.json`/`.txt` assets, e.g. the
   GraphQL SDL schema).
3. Fetches each page and **recurses**: any extension-less internal link
   (`/docs/<slug>`, `/developers/<slug>`) is retried as its `.md` twin, and any
   `.md` link is followed. Non-existent Markdown twins (HTTP 404) are skipped.
4. Writes each page under `linear-source-docs/` preserving the URL path, and
   emits `manifest.json`, `manifest.tsv`, and `MANIFEST.md`.

Re-running is idempotent: the corpus is overwritten in place and the manifest
regenerated. Diff `manifest.json` between runs to see what Linear changed.

### Why committed

The corpus is committed so the docs are greppable offline and so the
feature-gap analysis has a stable, reviewable source of truth. Regenerate it
periodically (or in CI) to detect drift.

## Feature gap

See [`FEATURE-GAP-ANALYSIS.md`](FEATURE-GAP-ANALYSIS.md) for the comparison of
`linear-cli`'s implemented surface against the live docs and GraphQL schema,
plus the prioritized plan to close the gap.
