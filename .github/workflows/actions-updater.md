---
name: Actions Version Updater
description: Scans .md workflow source files for outdated GitHub Actions and creates issues for each update, similar to Dependabot
on:
  schedule:
    - cron: '0 6 * * 1'  # Weekly on Monday at 6am UTC
  workflow_dispatch:
    inputs:
      repository:
        description: "Repository to scan (format: owner/repo). Defaults to current repo."
        required: false
        type: string
      organization:
        description: "Scan all repos in an organization (overrides repository input)"
        required: false
        type: string
      dry-run:
        description: "Only report outdated actions, don't create issues"
        required: false
        type: boolean
        default: false

permissions:
  contents: read
  issues: write
  pull-requests: read

engine: copilot
timeout-minutes: 15
strict: true

network:
  allowed:
    - defaults
    - github

tools:
  github:
    read-only: true
    lockdown: true
    toolsets:
      - repos
      - issues
  bash:
    - "find .github/workflows -name '*.md' -type f"
    - "grep -rn 'uses:' .github/workflows/*.md"
    - "cat .github/workflows/*.md"
    - "head -n * .github/workflows/*.md"
    - "gh api repos/*/releases/latest"
    - "gh api repos/*/tags"
    - "gh api repos/*/contents/*"
    - "gh repo list *"

safe-outputs:
  create-issue:
    expires: 7d
    title-prefix: "[actions-update] "
    labels: [dependencies, automation, github-actions]
    max: 20
    close-older-issues: true
  missing-tool:

if: needs.check_outdated.outputs.has_outdated == 'true'

jobs:
  check_outdated:
    runs-on: ubuntu-latest
    permissions:
      contents: read
    outputs:
      has_outdated: ${{ steps.check.outputs.has_outdated }}
      outdated_list: ${{ steps.check.outputs.outdated_list }}
      outdated_count: ${{ steps.check.outputs.outdated_count }}
    steps:
      - name: Checkout repository
        uses: actions/checkout@v6
        with:
          sparse-checkout: |
            .github/workflows
          persist-credentials: false

      - name: Extract and check action versions
        id: check
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          # Extract unique action@version pairs from .md workflow files
          grep -rhn 'uses:' .github/workflows/*.md 2>/dev/null \
            | sed -n 's/.*uses:[[:space:]]*\([^@#"'"'"'[:space:]]*\)@\([^#"'"'"'[:space:]]*\).*/\1@\2/p' \
            | sort -u > /tmp/current-actions.txt

          echo "Found $(wc -l < /tmp/current-actions.txt) unique action references"

          OUTDATED=0
          OUTDATED_LIST=""

          # Extract unique action names (without version)
          cut -d@ -f1 /tmp/current-actions.txt | sort -u > /tmp/action-names.txt

          while IFS= read -r action; do
            # Get current version(s) used
            current=$(grep "^${action}@" /tmp/current-actions.txt | cut -d@ -f2 | sort -V | tail -1)

            # For SHA-pinned refs, resolve to a version tag via the comment or API
            IS_SHA="false"
            if echo "$current" | grep -qE '^[0-9a-f]{40}$'; then
              IS_SHA="true"
              # Try to find the version from a comment in the .md file (e.g., # v6.1.0)
              comment_ver=$(grep -rh "${current}" .github/workflows/*.md 2>/dev/null | sed -n 's/.*#[[:space:]]*\(v[0-9][0-9.]*\).*/\1/p' | head -1)
              if [ -n "$comment_ver" ]; then
                current="$comment_ver"
              else
                # Resolve SHA to tag via API
                resolved=$(gh api "repos/${action}/git/matching-refs/tags" --jq '.[].ref' 2>/dev/null | while read ref; do
                  tag_sha=$(gh api "repos/${action}/git/ref/${ref#refs/}" --jq '.object.sha' 2>/dev/null)
                  if [ "$tag_sha" = "$current" ]; then
                    echo "${ref#refs/tags/}"
                    break
                  fi
                done | head -1)
                if [ -n "$resolved" ]; then
                  current="$resolved"
                else
                  echo "SKIP: ${action}@${current} — SHA could not be resolved to a tag"
                  continue
                fi
              fi
            fi

            # Get latest release tag
            latest=$(gh api "repos/${action}/releases/latest" --jq '.tag_name' 2>/dev/null || echo "")
            if [ -z "$latest" ]; then
              latest=$(gh api "repos/${action}/tags" --jq '.[0].name' 2>/dev/null || echo "")
            fi

            if [ -z "$latest" ]; then
              echo "SKIP: ${action} — could not determine latest version"
              continue
            fi

            # Normalize versions for comparison (strip leading 'v')
            current_num=$(echo "$current" | sed 's/^v//')
            latest_num=$(echo "$latest" | sed 's/^v//')
            current_major=$(echo "$current_num" | cut -d. -f1)
            latest_major=$(echo "$latest_num" | cut -d. -f1)

            if [ "$current_major" -lt "$latest_major" ] 2>/dev/null; then
              echo "OUTDATED: ${action} ${current} -> ${latest} (major, sha=${IS_SHA})"
              OUTDATED_LIST="${OUTDATED_LIST}${action}|${current}|${latest}|${IS_SHA}\n"
              OUTDATED=$((OUTDATED + 1))
            elif [ "$current_major" = "$latest_major" ] && [ "$current_num" != "$latest_num" ]; then
              if [ "$(printf '%s\n%s' "$current_num" "$latest_num" | sort -V | head -1)" = "$current_num" ] && [ "$current_num" != "$latest_num" ]; then
                echo "OUTDATED: ${action} ${current} -> ${latest} (minor/patch, sha=${IS_SHA})"
                OUTDATED_LIST="${OUTDATED_LIST}${action}|${current}|${latest}|${IS_SHA}\n"
                OUTDATED=$((OUTDATED + 1))
              fi
            else
              echo "OK: ${action} ${current} (latest: ${latest})"
            fi
          done < /tmp/action-names.txt

          echo ""
          echo "=== SUMMARY: ${OUTDATED} outdated actions ==="

          if [ "$OUTDATED" -gt 0 ]; then
            echo "has_outdated=true" >> "$GITHUB_OUTPUT"
            # Truncate to fit output limits
            echo "outdated_list=$(printf '%b' "$OUTDATED_LIST" | head -50)" >> "$GITHUB_OUTPUT"
          else
            echo "has_outdated=false" >> "$GITHUB_OUTPUT"
          fi
          echo "outdated_count=${OUTDATED}" >> "$GITHUB_OUTPUT"
---

# GitHub Actions Version Updater Agent 🔄

You are an expert CI/CD maintenance agent. Your job is to scan `.md` workflow source files for outdated GitHub Actions and **create issues** for each outdated action — similar to how Dependabot creates alerts for outdated dependencies.

> **You do NOT modify files or create PRs.** You only create issues. A developer will assign the issue to CCA (Claude Code Agent) which will handle the actual file changes and PR creation.

## Important Context

- This repository compiles `.md` files into `.lock.yml` workflow files
- The compiler automatically handles SHA pinning — `.md` files use **version tags** (e.g., `v6`, `v5.1.0`)
- Only `.md` files are the source of truth — `.lock.yml` files are generated

## Target Repository

- **Repository input**: `${{ github.event.inputs.repository || github.repository }}`
- **Organization input**: `${{ github.event.inputs.organization }}`
- **Dry run**: `${{ github.event.inputs.dry-run }}`

If an **organization** is provided, list all repos in the org and scan each one. Otherwise, scan only the target repository.

If **dry-run** is `true`, report findings but do **not** create issues.

## Phase 1: Determine Target Repos

If an **organization** input is provided, list all repos:

```bash
gh repo list {organization} --limit 500 --json nameWithOwner --jq '.[].nameWithOwner'
```

Otherwise, use the single target repository: `${{ github.event.inputs.repository || github.repository }}`

## Phase 2: Discover Current Action Versions

For each target repo, fetch the `.md` workflow files and scan for `uses:` directives.

If scanning the **current repo**:

```bash
find .github/workflows -name '*.md' -type f -exec grep -Hn 'uses:' {} \;
```

If scanning a **remote repo**:

```bash
gh api repos/{owner}/{repo}/contents/.github/workflows --jq '.[] | select(.name | endswith(".md")) | .name' | while read f; do
  echo "=== $f ==="
  gh api repos/{owner}/{repo}/contents/.github/workflows/$f -H "Accept: application/vnd.github.raw" | grep -n 'uses:'
done
```

Parse each match to extract:
- **Action name** (e.g., `actions/checkout`)
- **Current version** (e.g., `v5`, `v6.1.0`)
- **File** and **line number**

Build a deduplicated map of: `action → { current_versions, files_using_it }`

## Phase 3: Check for Updates

For each unique action found, check the latest release:

```bash
gh api repos/{owner}/{action}/releases/latest --jq '.tag_name'
```

Compare the current version to the latest. An action needs updating if:
- Its major version is behind (e.g., `v5` when `v6` is available)
- Its minor/patch version is behind within the same major (e.g., `v6.1.0` when `v6.3.0` is available)

Skip actions that are already on the latest version.

## Phase 4: Create Issues

> If **dry-run** is `true`, skip issue creation and just output the summary.

For each outdated action (or group of related actions), create an issue **in the target repo** with:

### Issue Title Format
`Upgrade {action} from {old_version} to {new_version}`

Example: `Upgrade actions/checkout from v5 to v6`

### Issue Body Format

```markdown
## Action Update Available

| Field | Value |
|-------|-------|
| **Action** | `{action}` |
| **Current Version** | `{old_version}` |
| **Latest Version** | `{new_version}` |
| **Release Notes** | [{new_version}](https://github.com/{action}/releases/tag/{new_version}) |

### Affected Files

The following `.md` workflow source files use this action:

- `.github/workflows/{file1}.md` (line {N})
- `.github/workflows/{file2}.md` (line {N})

### Update Instructions

Update `uses:` directives in the listed `.md` files:

\`\`\`yaml
# Before
uses: {action}@{old_version}

# After
uses: {action}@{new_version}
\`\`\`

Then run `gh aw compile` to regenerate `.lock.yml` files.

### Breaking Changes

{If major version bump: note that this is a major version upgrade and link to release notes for breaking changes}
{If Node 24 related: note Node 20 EOL April 2026 and the need for Node 24 compatible actions}

> 💡 **Assigned to CCA** — this issue will be automatically picked up to create a PR with the update.
```

### Grouping Rules

- **Group by action**: One issue per unique action that needs updating (not per file)
- **Node 24 label**: If the action is a core `actions/*` action being upgraded for Node 24 compatibility, add the `node24` label
- **Priority**: Create issues for `actions/*` (official) actions first, then third-party actions
- **Auto-assign to CCA**: After creating each issue, assign it to CCA so it automatically creates the PR with the `.md` file updates and runs `gh aw compile`

## Phase 5: Summary

After creating all issues, output a summary:

```
Actions Update Summary:
- N repositories scanned
- X actions checked
- Y actions need updating
- Z issues created

Per-repo breakdown:
  owner/repo-1: 3 outdated actions, 3 issues created
  owner/repo-2: all up to date
  ...

Issues created:
- owner/repo-1#123: Upgrade actions/checkout from v5 to v6
- owner/repo-1#124: Upgrade actions/setup-go from v5 to v6
...
```

If all actions are up to date, create a single informational issue noting that all actions are current.

## Rules

1. **Read-only on files** — do NOT modify any files, only create and assign issues
2. **One issue per action** — group all files using the same action into one issue
3. **Skip current actions** — don’t create issues for actions already on the latest version
4. **Check for existing issues** — before creating, check if an open issue for the same action/version already exists and skip if so
5. **Auto-assign to CCA** — after creating each issue, assign it to CCA so it automatically creates the PR
