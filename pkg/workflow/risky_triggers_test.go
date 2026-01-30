//go:build !integration

package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasRiskyTriggers(t *testing.T) {
	compiler := NewCompiler()

	tests := []struct {
		name       string
		onSection  string
		expectRisky bool
	}{
		{
			name: "issue_comment trigger",
			onSection: `on:
  issue_comment:
    types: [created]`,
			expectRisky: true,
		},
		{
			name: "pull_request_target trigger",
			onSection: `on:
  pull_request_target:
    types: [opened]`,
			expectRisky: true,
		},
		{
			name: "workflow_run trigger",
			onSection: `on:
  workflow_run:
    workflows: ["CI"]`,
			expectRisky: true,
		},
		{
			name: "pull_request_review_comment trigger",
			onSection: `on:
  pull_request_review_comment:
    types: [created]`,
			expectRisky: true,
		},
		{
			name: "multiple triggers including risky",
			onSection: `on:
  push:
    branches: [main]
  issue_comment:
    types: [created]`,
			expectRisky: true,
		},
		{
			name: "pull_request trigger (not risky)",
			onSection: `on:
  pull_request:
    types: [opened]`,
			expectRisky: false,
		},
		{
			name: "push trigger (not risky)",
			onSection: `on:
  push:
    branches: [main]`,
			expectRisky: false,
		},
		{
			name: "issues trigger (not risky)",
			onSection: `on:
  issues:
    types: [opened]`,
			expectRisky: false,
		},
		{
			name: "schedule trigger (not risky)",
			onSection: `on:
  schedule:
    - cron: '0 0 * * *'`,
			expectRisky: false,
		},
		{
			name: "workflow_dispatch trigger (not risky)",
			onSection: `on:
  workflow_dispatch:`,
			expectRisky: false,
		},
		{
			name: "risky trigger in comment should not match",
			onSection: `on:
  push:
    branches: [main]
  # issue_comment: trigger would be risky`,
			expectRisky: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compiler.hasRiskyTriggers(tt.onSection)
			assert.Equal(t, tt.expectRisky, result,
				"Expected hasRiskyTriggers to return %v for %s", tt.expectRisky, tt.name)
		})
	}
}
