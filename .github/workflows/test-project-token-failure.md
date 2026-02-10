---
name: Test Project Token Failure
description: Test workflow to verify token failure path for project-related safe outputs
on:
  workflow_dispatch:

permissions:
  contents: read
  issues: read
  pull-requests: read

engine: copilot
strict: true
timeout-minutes: 5

safe-outputs:
  update-project:
    max: 1
    project: "https://github.com/orgs/github/projects/24068"
  create-project-status-update:
    max: 1
    project: "https://github.com/orgs/github/projects/24068"
  create-project:
    max: 1

tools:
  cache-memory: true

imports:
  - shared/mood.md
---

# Test Project Token Failure Path

This workflow specifically tests the **token failure path** for project-related safe outputs.

## Purpose

Project-related operations (update-project, create-project, create-project-status-update) require a Personal Access Token (PAT) with Projects access. The default `GITHUB_TOKEN` used by GitHub Actions (`github-actions[bot]`) does not have Projects v2 access.

This workflow intentionally:
- Uses the default GITHUB_TOKEN (no custom token configured)
- Attempts project-related operations
- Triggers token failure detection and error messages

## Expected Behavior

When you attempt to use project-related safe outputs WITHOUT a proper token:

1. **Token Detection**: The system should detect that authentication is `github-actions[bot]`
2. **Early Failure**: Should fail fast with clear error message BEFORE attempting GraphQL queries
3. **Actionable Guidance**: Error message should explain:
   - Why the operation failed (default token doesn't have Projects access)
   - What token is needed (PAT with Projects scope)
   - How to fix it (set `secrets.GH_AW_PROJECT_GITHUB_TOKEN` or configure `safe-outputs.*.github-token`)

## Test Cases

### Test 1: Update Project with Default Token

Try to update a project item using the default token:

```json
{
  "type": "update_project",
  "project": "https://github.com/orgs/example-org/projects/1",
  "content_type": "issue",
  "content_number": 1,
  "fields": {
    "status": "In Progress"
  }
}
```

Expected: Should fail with clear error about github-actions[bot] token.

### Test 2: Create Project Status Update with Default Token

Try to create a project status update:

```json
{
  "type": "create_project_status_update",
  "project": "https://github.com/orgs/example-org/projects/1",
  "body": "Test status update to trigger token failure",
  "status": "ON_TRACK"
}
```

Expected: Should fail with INSUFFICIENT_SCOPES or authentication error.

### Test 3: Create Project with Default Token

Try to create a new project:

```json
{
  "type": "create_project",
  "owner_type": "org",
  "owner_login": "example-org",
  "title": "Test Project Creation with Default Token"
}
```

Expected: Should fail indicating token lacks Projects creation permissions.

## Task

Attempt each of the three test cases above. For each attempt:

1. **Call the safe output tool** with the example JSON
2. **Document the error message** you receive
3. **Verify the error is actionable** - does it clearly explain:
   - What went wrong?
   - Why it failed?
   - How to fix it?

## Success Criteria

The token failure path is working correctly if:

- ✅ Errors are detected early (before making API calls)
- ✅ Error messages mention `github-actions[bot]` or `GITHUB_TOKEN`
- ✅ Error messages explain Projects v2 requires PAT
- ✅ Error messages provide fix instructions with specific secret names
- ✅ No confusing fallback behavior or silent failures

Report your findings for each test case.
