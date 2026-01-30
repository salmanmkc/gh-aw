---
description: Test workflow to verify that skills from .github/skills are automatically loaded
on:
  workflow_dispatch:
  pull_request:
    types: [labeled]
    names: ["smoke"]
permissions:
  actions: read
  contents: read
  issues: read
  pull-requests: read
name: Skills Loading Test
engine: copilot
strict: true
network:
  allowed:
    - defaults
    - github
tools:
  agentic-workflows:
  bash:
    - "*"
  github:
    toolsets: [default]
safe-outputs:
  noop:
  add-comment:
    hide-older-comments: true
    max: 2
  messages:
    footer: "> 🔧 *Skills test by [{workflow_name}]({run_url})*"
    run-started: "🔬 Testing skills loading... [{workflow_name}]({run_url}) is verifying `.github/skills` discovery..."
    run-success: "✅ Skills loading test completed! [{workflow_name}]({run_url}) confirms all skills are accessible."
    run-failure: "❌ Skills loading test failed! [{workflow_name}]({run_url}) detected issues: {status}"
timeout-minutes: 10
---

{{#runtime-import agentics/test-skills-loading.md}}
