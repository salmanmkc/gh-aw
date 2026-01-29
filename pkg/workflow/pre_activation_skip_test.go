//go:build !integration

package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/githubnext/gh-aw/pkg/constants"
	"github.com/githubnext/gh-aw/pkg/stringutil"
	"github.com/githubnext/gh-aw/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPreActivationSkipForScheduleEvents verifies that the pre_activation job is skipped
// for safe events (schedule, merge_group) and that downstream jobs handle this correctly.
func TestPreActivationSkipForScheduleEvents(t *testing.T) {
	tmpDir := testutil.TempDir(t, "pre-activation-skip-test")
	compiler := NewCompiler()

	t.Run("schedule_only_workflow_has_no_pre_activation", func(t *testing.T) {
		// Schedule-only workflows don't need a pre_activation job at all
		// because schedule is in SafeWorkflowEvents - no permission check needed
		workflowContent := `---
on:
  schedule:
    - cron: "0 8 * * *"
engine: codex
---

# Daily Schedule Workflow

This workflow runs on a daily schedule.
`
		workflowFile := filepath.Join(tmpDir, "schedule-only-workflow.md")
		require.NoError(t, os.WriteFile(workflowFile, []byte(workflowContent), 0644))

		err := compiler.CompileWorkflow(workflowFile)
		require.NoError(t, err, "Should compile workflow")

		lockFile := stringutil.MarkdownToLockFile(workflowFile)
		lockContent, err := os.ReadFile(lockFile)
		require.NoError(t, err, "Should read lock file")
		lockContentStr := string(lockContent)

		// Schedule-only workflows should NOT have a pre_activation job
		// because schedule is already considered a safe event
		preActivationSection := extractJobSectionForSkipTest(lockContentStr, string(constants.PreActivationJobName))
		assert.Empty(t, preActivationSection, "Schedule-only workflow should NOT have pre_activation job")
	})

	t.Run("schedule_and_issue_workflow_has_pre_activation_with_skip_condition", func(t *testing.T) {
		// Workflows with both schedule and issue triggers need pre_activation
		// but should skip it for schedule events
		workflowContent := `---
on:
  schedule:
    - cron: "0 8 * * *"
  issues:
    types: [opened]
engine: codex
---

# Mixed trigger workflow

This workflow runs on schedule and when issues are opened.
`
		workflowFile := filepath.Join(tmpDir, "schedule-and-issue-workflow.md")
		require.NoError(t, os.WriteFile(workflowFile, []byte(workflowContent), 0644))

		err := compiler.CompileWorkflow(workflowFile)
		require.NoError(t, err, "Should compile workflow")

		lockFile := stringutil.MarkdownToLockFile(workflowFile)
		lockContent, err := os.ReadFile(lockFile)
		require.NoError(t, err, "Should read lock file")
		lockContentStr := string(lockContent)

		// This workflow should have pre_activation because of the issues trigger
		preActivationSection := extractJobSectionForSkipTest(lockContentStr, string(constants.PreActivationJobName))
		require.NotEmpty(t, preActivationSection, "Should have pre_activation job")

		// The pre_activation job should have an if condition that skips for schedule
		assert.Contains(t, preActivationSection, "github.event_name != 'schedule'",
			"pre_activation should have condition to skip for schedule events")
		assert.Contains(t, preActivationSection, "github.event_name != 'merge_group'",
			"pre_activation should have condition to skip for merge_group events")

		// Verify activation job handles skipped pre_activation
		activationSection := extractJobSectionForSkipTest(lockContentStr, string(constants.ActivationJobName))
		require.NotEmpty(t, activationSection, "Should have activation job")
		assert.Contains(t, activationSection, "needs.pre_activation.result == 'skipped'",
			"activation should check if pre_activation was skipped")
		assert.Contains(t, activationSection, "needs.pre_activation.outputs.activated == 'true'",
			"activation should also check activated output for non-safe events")
	})

	t.Run("issue_workflow_does_not_skip_pre_activation", func(t *testing.T) {
		// Issue-triggered workflows should NOT skip pre_activation since
		// they require permission checks
		workflowContent := `---
on:
  issues:
    types: [opened]
engine: codex
---

# Issue Workflow

This workflow runs when issues are opened.
`
		workflowFile := filepath.Join(tmpDir, "issue-workflow.md")
		require.NoError(t, os.WriteFile(workflowFile, []byte(workflowContent), 0644))

		err := compiler.CompileWorkflow(workflowFile)
		require.NoError(t, err, "Should compile workflow")

		lockFile := stringutil.MarkdownToLockFile(workflowFile)
		lockContent, err := os.ReadFile(lockFile)
		require.NoError(t, err, "Should read lock file")
		lockContentStr := string(lockContent)

		// Verify pre_activation exists and has skip condition
		// (the skip condition is always present, but for issues it won't trigger)
		preActivationSection := extractJobSectionForSkipTest(lockContentStr, string(constants.PreActivationJobName))
		require.NotEmpty(t, preActivationSection, "Should have pre_activation job")

		// For issue-only workflows, the skip condition still exists but won't match
		// because github.event_name will be 'issues', not 'schedule' or 'merge_group'
		assert.Contains(t, preActivationSection, "github.event_name != 'schedule'",
			"pre_activation should have schedule skip condition")
	})

	t.Run("merge_group_workflow_skips_pre_activation", func(t *testing.T) {
		workflowContent := `---
on:
  merge_group:
    types: [checks_requested]
engine: codex
---

# Merge Queue Workflow

This workflow runs in the merge queue.
`
		workflowFile := filepath.Join(tmpDir, "merge-group-workflow.md")
		require.NoError(t, os.WriteFile(workflowFile, []byte(workflowContent), 0644))

		err := compiler.CompileWorkflow(workflowFile)
		require.NoError(t, err, "Should compile workflow")

		lockFile := stringutil.MarkdownToLockFile(workflowFile)
		lockContent, err := os.ReadFile(lockFile)
		require.NoError(t, err, "Should read lock file")
		lockContentStr := string(lockContent)

		// Verify pre_activation has skip condition for merge_group
		preActivationSection := extractJobSectionForSkipTest(lockContentStr, string(constants.PreActivationJobName))
		require.NotEmpty(t, preActivationSection, "Should have pre_activation job")
		assert.Contains(t, preActivationSection, "github.event_name != 'merge_group'",
			"pre_activation should skip for merge_group events")
	})
}

// TestSafePreActivationEventsListConsistency verifies that the safe events list
// in expression_builder.go matches what check_membership.cjs considers safe.
func TestSafePreActivationEventsListConsistency(t *testing.T) {
	// The safe events should match those defined in check_membership.cjs
	// This test ensures we don't accidentally diverge
	expectedSafeEvents := []string{"schedule", "merge_group"}

	assert.Equal(t, expectedSafeEvents, SafePreActivationEvents,
		"SafePreActivationEvents should match the safe events in check_membership.cjs")
}

// extractJobSectionForSkipTest extracts a job section from YAML content
// Named differently to avoid conflict with extractJobSection in compiler_test_helpers.go
func extractJobSectionForSkipTest(content string, jobName string) string {
	lines := strings.Split(content, "\n")
	var result strings.Builder
	inJob := false
	indent := 0

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		currentIndent := len(line) - len(trimmed)

		if !inJob {
			// Look for the job start
			if strings.HasPrefix(trimmed, jobName+":") {
				inJob = true
				indent = currentIndent
				result.WriteString(line + "\n")
			}
		} else {
			// We're inside the job
			if trimmed == "" {
				// Empty line, include it
				result.WriteString(line + "\n")
				continue
			}

			if currentIndent <= indent && !strings.HasPrefix(trimmed, "#") {
				// We've reached a new job at the same or lower indentation
				break
			}
			result.WriteString(line + "\n")
		}
	}

	return result.String()
}
