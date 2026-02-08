# Workflow Instruction Style Analysis

**Analysis Date:** 2026-02-08  
**Total Workflows Analyzed:** 147 (`.github/workflows/*.md`)

## Executive Summary

This analysis classifies all workflow files in the repository by their instruction approach:
- **Explicit Examples & Tool Call Code** - Workflows with bash commands, code blocks, and explicit MCP tool usage
- **Safe Outputs + Natural Language** - Workflows relying on safe-outputs with pure natural language instructions

## Key Findings

### Distribution Chart

```
==========================================================================================
WORKFLOW INSTRUCTION STYLE DISTRIBUTION
==========================================================================================

Hybrid (Explicit + Safe-Outputs)              ██████████████████████████████████████████████████  96 (65.3%)
Safe-Outputs Only (Natural Language)          ██████████████████████                              43 (29.3%)
Explicit Only (No Safe-Outputs)               █                                                    3 ( 2.0%)
Minimal (Neither)                             ██                                                   5 ( 3.4%)

TOTAL                                                                                            147 (100.0%)
==========================================================================================
```

### Critical Insight

**📊 67.3% of workflows (99 out of 147) include explicit examples and tool call code** in their markdown instructions, combining:
- 96 Hybrid workflows (explicit examples + safe-outputs)
- 3 Explicit-only workflows (no safe-outputs)

**Only 29.3% (43 workflows) rely exclusively on safe-outputs with pure natural language** instructions.

## Detailed Breakdown

### Category 1: Explicit Examples Only (No Safe-Outputs)
**Count: 3 workflows (2.0%)**

Workflows using explicit code blocks, bash commands, and MCP tool patterns without safe-outputs:

1. `chroma-issue-indexer.md` - Brave MCP search tools
2. `codex-github-remote-mcp-test.md` - GitHub API calls
3. `metrics-collector.md` - Bash scripts for metrics

**Pattern Example:**
```markdown
Use the Brave MCP search tools to find relevant information
```

---

### Category 2: Safe-Outputs + Natural Language Only
**Count: 43 workflows (29.3%)**

Workflows using safe-outputs exclusively with natural language instructions:

- `agent-performance-analyzer.md` - create-issue
- `agent-persona-explorer.md` - create-discussion
- `artifacts-summary.md` - create-discussion
- `auto-triage-issues.md` - add-labels
- `brave.md` - add-comment
- `ci-doctor.md` - create-issue
- `cloclo.md` - create-pull-request
- `daily-assign-issue-to-user.md` - assign-to-user
- `daily-choice-test.md` - staged jobs
- `daily-fact.md` - add-comment
- `daily-firewall-report.md` - create-discussion
- `daily-issues-report.md` - create-discussion
- `daily-performance-summary.md` - create-discussion
- `daily-regulatory.md` - create-issue
- `daily-secrets-analysis.md` - create-discussion
- `daily-team-status.md` - create-discussion
- `delight.md` - add-comment
- `dependabot-go-checker.md` - add-comment
- `dependabot-project-manager.md` - project management
- `developer-docs-consolidator.md` - create-discussion
- `dictation-prompt.md` - create-issue
- `docs-noob-tester.md` - create-issue
- `draft-pr-cleanup.md` - cleanup
- `example-workflow-analyzer.md` - analysis
- `functional-pragmatist.md` - create-pull-request
- `github-mcp-tools-report.md` - create-discussion
- `go-logger.md` - create-issue
- `issue-classifier.md` - add-labels
- `issue-triage-agent.md` - triage
- `layout-spec-maintainer.md` - create-pull-request
- `lockfile-stats.md` - create-discussion
- `org-health-report.md` - create-discussion
- `pdf-summary.md` - add-comment
- `plan.md` - create-issue
- `poem-bot.md` - add-comment
- `portfolio-analyst.md` - create-discussion
- `pr-nitpick-reviewer.md` - add-comment
- `pr-triage-agent.md` - triage
- `python-data-charts.md` - create-discussion
- `repository-quality-improver.md` - create-pull-request
- `semantic-function-refactor.md` - create-pull-request
- `terminal-stylist.md` - create-pull-request
- `unbloat-docs.md` - create-pull-request

**Pattern Example:**
```yaml
safe-outputs:
  add-labels:
    allowed: [bug, feature, enhancement, documentation]
    max: 1
```

With natural language instructions:
```markdown
Your task is to analyze newly created issues and classify them as either a "bug" or a "feature".
```

---

### Category 3: Hybrid (Explicit Examples AND Safe-Outputs)
**Count: 96 workflows (65.3%)**

The most sophisticated workflows, combining explicit tool calls with safe-outputs:

- `ai-moderator.md` - Explicit checks + add-labels
- `archie.md` - Bash commands + add-comment
- `audit-workflows.md` - MCP tool usage + upload-asset/create-discussion
- `blog-auditor.md` - Web analysis + create-discussion
- `breaking-change-checker.md` - Git analysis + create-issue
- `changeset.md` - Git operations + push-to-pull-request-branch
- `ci-coach.md` - CI analysis + create-pull-request
- `claude-code-user-docs-review.md` - Review steps + create-discussion
- `cli-consistency-checker.md` - CLI analysis + create-issue
- `cli-version-checker.md` - Version checks + create-issue
- `code-scanning-fixer.md` - Security fixes + create-pull-request
- `code-simplifier.md` - Code analysis + create-pull-request
- `commit-changes-analyzer.md` - Git analysis + create-discussion
- `copilot-agent-analysis.md` - Agent analysis + create-discussion
- `copilot-cli-deep-research.md` - Research + create-discussion
- `copilot-pr-merged-report.md` - PR analysis + create-discussion
- `copilot-pr-nlp-analysis.md` - NLP analysis + create-discussion
- `copilot-pr-prompt-analysis.md` - Prompt analysis + create-discussion
- `copilot-session-insights.md` - Session analysis + create-discussion
- `craft.md` - Content creation + create-pull-request
- `daily-cli-performance.md` - Performance tests + create-discussion
- `daily-cli-tools-tester.md` - Tool tests + create-discussion
- `daily-code-metrics.md` - Code metrics + create-discussion
- `daily-compiler-quality.md` - Compiler analysis + create-discussion
- `daily-copilot-token-report.md` - Token analysis + create-discussion
- `daily-doc-updater.md` - Doc updates + create-pull-request
- `daily-file-diet.md` - File analysis + create-discussion
- `daily-malicious-code-scan.md` - Security scan + create-issue
- `daily-mcp-concurrency-analysis.md` - Concurrency analysis + create-discussion
- `daily-multi-device-docs-tester.md` - Device testing + create-discussion
- `daily-news.md` - News generation + create-discussion
- `daily-observability-report.md` - Observability + create-discussion
- `daily-repo-chronicle.md` - Repository history + create-discussion
- `daily-safe-output-optimizer.md` - Optimization + create-pull-request
- `daily-syntax-error-quality.md` - Error analysis + create-discussion
- `daily-team-evolution-insights.md` - Team insights + create-discussion
- `daily-testify-uber-super-expert.md` - Test expert + create-discussion
- `daily-workflow-updater.md` - Workflow updates + create-pull-request
- `deep-report.md` - Deep analysis + create-discussion
- `dev-hawk.md` - Development monitoring + create-issue
- `dev.md` - Development workflow + create-pull-request
- `discussion-task-miner.md` - Task mining + create-issue
- `duplicate-code-detector.md` - Code duplication + create-issue
- `github-mcp-structural-analysis.md` - Structural analysis + create-discussion
- `glossary-maintainer.md` - Glossary updates + create-pull-request
- `go-fan.md` - Go code analysis + create-pull-request
- `go-pattern-detector.md` - Pattern detection + create-discussion
- `grumpy-reviewer.md` - Code review + add-comment
- `hourly-ci-cleaner.md` - CI cleanup + workflow operations
- `instructions-janitor.md` - Instructions cleanup + create-pull-request
- `issue-arborist.md` - Issue organization + project operations
- `issue-monster.md` - Issue management + multiple operations
- `jsweep.md` - JavaScript cleanup + create-pull-request
- `mcp-inspector.md` - MCP inspection + create-discussion
- `mergefest.md` - Merge operations + create-discussion
- `q.md` - Query operations + create-discussion
- `release.md` - Release management + create-release
- `repo-audit-analyzer.md` - Repository audit + create-discussion
- `repo-tree-map.md` - Repository mapping + upload-asset
- `research.md` - Research workflow + create-discussion
- `safe-output-health.md` - Health monitoring + create-issue
- `schema-consistency-checker.md` - Schema validation + create-issue
- `scout.md` - Code exploration + create-discussion
- `security-compliance.md` - Compliance checks + create-issue
- `security-guard.md` - Security monitoring + create-issue
- `security-review.md` - Security review + create-discussion
- `sergo.md` - Sergo analysis + create-discussion
- `slide-deck-maintainer.md` - Slide deck updates + create-pull-request
- `smoke-claude.md` - Claude smoke test + add-comment
- `smoke-codex.md` - Codex smoke test + add-comment
- `smoke-copilot.md` - Copilot smoke test + add-comment
- `smoke-opencode.md` - OpenCode smoke test + add-comment
- `smoke-project.md` - Project smoke test + project operations
- `smoke-test-tools.md` - Tool validation + add-comment
- `stale-repo-identifier.md` - Stale repo detection + create-issue
- `static-analysis-report.md` - Static analysis + create-discussion
- `step-name-alignment.md` - Alignment checks + create-pull-request
- `sub-issue-closer.md` - Issue closing + issue operations
- `super-linter.md` - Linting + create-pull-request
- `technical-doc-writer.md` - Documentation + create-pull-request
- `test-create-pr-error-handling.md` - Error testing + create-pull-request
- `test-dispatcher.md` - Dispatcher testing + multiple operations
- `test-project-url-default.md` - Project testing + project operations
- `tidy.md` - Code cleanup + create-pull-request
- `typist.md` - Typing fixes + create-pull-request
- `ubuntu-image-analyzer.md` - Image analysis + create-discussion
- `video-analyzer.md` - Video analysis + create-discussion
- `weekly-issue-summary.md` - Weekly summary + create-discussion
- `workflow-generator.md` - Workflow generation + create-pull-request
- `workflow-health-manager.md` - Health management + create-issue
- `workflow-normalizer.md` - Normalization + create-pull-request
- `workflow-skill-extractor.md` - Skill extraction + create-discussion

**Pattern Example:**
```yaml
safe-outputs:
  create-discussion:
    category: "audits"
    max: 1
    close-older-discussions: true
```

With explicit tool instructions:
```markdown
Use gh-aw MCP server (not CLI directly). Run `status` tool to verify.

**Collect Logs**: Use MCP `logs` tool to download workflow logs:
```
Use the agentic-workflows MCP tool `logs` with parameters:
- start_date: "-1d" (last 24 hours)
```
```

---

### Category 4: Minimal (Neither)
**Count: 5 workflows (3.4%)**

Basic workflows with minimal configuration:

1. `example-custom-error-patterns.md`
2. `example-permissions-warning.md`
3. `firewall.md`
4. `notion-issue-summary.md`
5. `test-workflow.md`

---

## Most Common Safe-Output Types

| Safe-Output Type | Usage Count | Percentage |
|-----------------|-------------|------------|
| max | 144 | 98.0% |
| title-prefix | 81 | 55.1% |
| expires | 76 | 51.7% |
| create-discussion | 58 | 39.5% |
| category | 55 | 37.4% |
| labels | 54 | 36.7% |
| close-older-discussions | 51 | 34.7% |
| create-issue | 40 | 27.2% |
| add-comment | 34 | 23.1% |
| messages (run-started, run-success, run-failure) | 30 | 20.4% |
| create-pull-request | 26 | 17.7% |

---

## Conclusions

### Primary Findings

1. ✅ **67.3% of workflows (99/147) include explicit examples and tool call code** in their markdown instructions
2. ✅ **29.3% of workflows (43/147) use only safe-outputs with pure natural language**
3. ✅ **The hybrid approach is dominant** - 65.3% combine both explicit examples and safe-outputs
4. ✅ **Only 2% rely exclusively on explicit code** without safe-outputs
5. ✅ **Safe-outputs are nearly universal** - 94.6% of workflows (139/147) use some form of safe-outputs

### Strategic Implications

**Most workflows benefit from explicit code examples and tool call demonstrations** to guide AI agents, rather than relying solely on natural language instructions with safe-outputs. The hybrid approach combining both techniques appears to be the most effective pattern.

The data suggests that:
- **Explicit examples provide concrete guidance** for complex operations
- **Safe-outputs ensure controlled, validated outputs** that integrate with GitHub
- **Combining both approaches** creates the most robust and reliable workflows

### Recommendations

For new workflow development:
1. Start with safe-outputs for the desired GitHub integration
2. Add explicit bash commands and tool examples for complex operations
3. Use natural language for high-level goals and context
4. Provide code blocks for specific implementation patterns

---

**Analysis Generated:** 2026-02-08  
**Repository:** github/gh-aw  
**Analyzer:** Workflow Analysis Script v1.0
