# Contributing to linear-cli

Thanks for contributing! This repo aims to keep changes simple, focused, and tested.

## Development

- Requirements: Go 1.22+, `gh` CLI (optional), `jq` (for examples), `golangci-lint` (optional).
- Useful targets:
  - `make deps` — install/tidy deps
  - `make build` — local build
  - `make test` — smoke tests (read-only commands)
  - `make lint` — lint if you have golangci-lint
  - `make fmt` — go fmt

## Release Checklist

Follow this checklist to cut a new release and update Homebrew:

1) Prepare
- Ensure README and help text match behavior.
- Run `make test` to verify smoke tests pass.
- Optionally draft release notes (highlights, fixes, breaking changes).

2) Tag and Release (vX.Y.Z)
- Create tag and push:
  ```bash
  git tag vX.Y.Z -a -m "vX.Y.Z: short summary"
  git push origin vX.Y.Z
  ```
- Create GitHub release (with notes):
  ```bash
  gh release create vX.Y.Z \
    --title "linear-cli vX.Y.Z" \
    --notes "<highlights/fixes>"
  ```

3) Homebrew Tap Bump (auto)
- This repo has a GitHub Action that auto-opens a PR to the tap on release publish.
- Required secret: `HOMEBREW_TAP_TOKEN` (fine‑grained PAT with contents:write on `roboalchemist/homebrew-linear-cli`).
  - Add in GitHub: repo Settings → Secrets and variables → Actions → New repository secret.

4) Homebrew Tap Bump (manual fallback)
If the action is disabled or no secret is configured:
```bash
TAG=vX.Y.Z
TARBALL=https://github.com/roboalchemist/linear-cli/archive/refs/tags/${TAG}.tar.gz
curl -sL "$TARBALL" -o /tmp/linear-cli.tgz
SHA=$(shasum -a 256 /tmp/linear-cli.tgz | awk '{print $1}')

git clone https://github.com/roboalchemist/homebrew-linear-cli.git
cd homebrew-linear-cli
git checkout -b bump-linear-cli-${TAG#v}
sed -i.bak -E "s|url \"[^\"]+\"|url \"$TARBALL\"|g" Formula/linear-cli.rb
sed -i.bak -E "s|sha256 \"[0-9a-f]+\"|sha256 \"$SHA\"|g" Formula/linear-cli.rb
rm -f Formula/linear-cli.rb.bak
git commit -am "linear-cli: bump to ${TAG}"
git push -u origin HEAD
gh pr create --title "linear-cli: bump to ${TAG}" --body "Update formula to ${TAG}." --base master --head bump-linear-cli-${TAG#v}
```

5) Validate
- After the tap PR merges:
  ```bash
  brew update && brew upgrade linear-cli
  linear-cli --version
  linear-cli docs | head -n 5
  ```
- Run a quick smoke test against your Linear workspace if possible.

6) Housekeeping
- Close any issues tied to the release.
- Start a new milestone if applicable.

