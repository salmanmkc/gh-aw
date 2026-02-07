---
name: Dependabot Project Manager
description: Automatically bundles Dependabot PRs by runtime and manifest, creates project items, and assigns them to Copilot for remediation with a "Review Required" status column
on:
  #schedule: daily
  workflow_dispatch:

timeout-minutes: 30

permissions:
  contents: read
  issues: read
  pull-requests: read
  security-events: read

network:
  allowed:
    - defaults
    - github

tools:
  github:
    toolsets:
      - default
      - dependabot
      - projects
  bash:
    - "jq *"

safe-outputs:
  update-project:
    project: "https://github.com/orgs/github/projects/24060"
    max: 50
    github-token: ${{ secrets.GH_AW_PROJECT_GITHUB_TOKEN }}
    views:
      - name: "Dependabot Alerts Board"
        layout: board
        filter: "is:open"
      - name: "Review Required"
        layout: board
        filter: 'is:open status:"Review Required"'
      - name: "All Alerts Table"
        layout: table

  create-project-status-update:
    project: "https://github.com/orgs/github/projects/24060"
    max: 1
    github-token: ${{ secrets.GH_AW_PROJECT_GITHUB_TOKEN }}

  create-issue:
    expires: 7d
    title-prefix: "[Dependabot Bundle] "
    labels: [dependencies, dependabot]
    assignees: copilot  # Automatically assigns Copilot when creating issues
    max: 20
    group: false

  add-comment:
    max: 10
---

# Dependabot Project Manager

- Find all open Dependabot PRs and add them to the project (up to max 50).
- Create bundle issues, each for exactly **one runtime + one manifest file** (up to max 20).
- Prioritize security updates first, then oldest PRs. If limits exceeded, create a summary issue noting overflow.
- Add bundle issues to the project, and assign them to Copilot.
