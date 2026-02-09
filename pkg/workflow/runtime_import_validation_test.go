//go:build !integration

package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExtractRuntimeImportPaths tests the extractRuntimeImportPaths function
func TestExtractRuntimeImportPaths(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "no imports",
			content:  "# Simple markdown\n\nSome text here",
			expected: nil,
		},
		{
			name:     "single file import",
			content:  "{{#runtime-import ./shared.md}}",
			expected: []string{"./shared.md"},
		},
		{
			name:     "optional import",
			content:  "{{#runtime-import? ./optional.md}}",
			expected: []string{"./optional.md"},
		},
		{
			name:     "import with line range",
			content:  "{{#runtime-import ./file.md:10-20}}",
			expected: []string{"./file.md"},
		},
		{
			name:     "multiple imports",
			content:  "{{#runtime-import ./a.md}}\n{{#runtime-import ./b.md}}",
			expected: []string{"./a.md", "./b.md"},
		},
		{
			name:     "duplicate imports",
			content:  "{{#runtime-import ./shared.md}}\n{{#runtime-import ./shared.md}}",
			expected: []string{"./shared.md"}, // Deduplicated
		},
		{
			name:     "URL import (should be excluded)",
			content:  "{{#runtime-import https://example.com/file.md}}",
			expected: nil,
		},
		{
			name:     "mixed file and URL imports",
			content:  "{{#runtime-import ./local.md}}\n{{#runtime-import https://example.com/remote.md}}",
			expected: []string{"./local.md"},
		},
		{
			name:     ".github prefix in path",
			content:  "{{#runtime-import .github/shared/common.md}}",
			expected: []string{".github/shared/common.md"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractRuntimeImportPaths(tt.content)

			if tt.expected == nil {
				assert.Nil(t, result, "Expected nil result")
			} else {
				assert.Equal(t, tt.expected, result, "Extracted paths mismatch")
			}
		})
	}
}

// TestValidateRuntimeImportFiles tests the validateRuntimeImportFiles function
func TestValidateRuntimeImportFiles(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()
	githubDir := filepath.Join(tmpDir, ".github")
	sharedDir := filepath.Join(githubDir, "shared")
	require.NoError(t, os.MkdirAll(sharedDir, 0755))

	// Create test files with different content
	validFile := filepath.Join(sharedDir, "valid.md")
	validContent := `# Valid Content

This file has safe expressions:
- Actor: ${{ github.actor }}
- Repository: ${{ github.repository }}
- Issue number: ${{ github.event.issue.number }}
`
	require.NoError(t, os.WriteFile(validFile, []byte(validContent), 0644))

	invalidFile := filepath.Join(sharedDir, "invalid.md")
	invalidContent := `# Invalid Content

This file has unsafe expressions:
- Secret: ${{ secrets.MY_TOKEN }}
- Runner: ${{ runner.os }}
`
	require.NoError(t, os.WriteFile(invalidFile, []byte(invalidContent), 0644))

	multilineFile := filepath.Join(sharedDir, "multiline.md")
	multilineContent := `# Multiline Expression

This has a multiline expression:
${{ github.actor
    && github.run_id }}
`
	require.NoError(t, os.WriteFile(multilineFile, []byte(multilineContent), 0644))

	tests := []struct {
		name        string
		markdown    string
		expectError bool
		errorText   string
	}{
		{
			name:        "no runtime imports",
			markdown:    "# Simple workflow\n\nNo imports here",
			expectError: false,
		},
		{
			name:        "valid runtime import",
			markdown:    "{{#runtime-import ./shared/valid.md}}",
			expectError: false,
		},
		{
			name:        "invalid runtime import",
			markdown:    "{{#runtime-import ./shared/invalid.md}}",
			expectError: true,
			errorText:   "secrets.MY_TOKEN",
		},
		{
			name:        "multiline expression in import",
			markdown:    "{{#runtime-import ./shared/multiline.md}}",
			expectError: true,
			errorText:   "unauthorized expressions",
		},
		{
			name:        "multiple imports with one invalid",
			markdown:    "{{#runtime-import ./shared/valid.md}}\n{{#runtime-import ./shared/invalid.md}}",
			expectError: true,
			errorText:   "secrets.MY_TOKEN",
		},
		{
			name:        "non-existent file (should skip)",
			markdown:    "{{#runtime-import ./shared/nonexistent.md}}",
			expectError: false, // Should skip validation for non-existent files
		},
		{
			name:        "URL import (should skip)",
			markdown:    "{{#runtime-import https://example.com/remote.md}}",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRuntimeImportFiles(tt.markdown, tmpDir)

			if tt.expectError {
				require.Error(t, err, "Expected an error")
				if tt.errorText != "" {
					assert.Contains(t, err.Error(), tt.errorText, "Error should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}

// TestValidateRuntimeImportFiles_PathNormalization tests path normalization
func TestValidateRuntimeImportFiles_PathNormalization(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	githubDir := filepath.Join(tmpDir, ".github")
	sharedDir := filepath.Join(githubDir, "shared")
	require.NoError(t, os.MkdirAll(sharedDir, 0755))

	// Create a valid test file
	validFile := filepath.Join(sharedDir, "test.md")
	validContent := "# Test\n\nActor: ${{ github.actor }}"
	require.NoError(t, os.WriteFile(validFile, []byte(validContent), 0644))

	tests := []struct {
		name        string
		markdown    string
		expectError bool
	}{
		{
			name:        "path with ./",
			markdown:    "{{#runtime-import ./shared/test.md}}",
			expectError: false,
		},
		{
			name:        "path with .github/",
			markdown:    "{{#runtime-import .github/shared/test.md}}",
			expectError: false,
		},
		{
			name:        "path without prefix",
			markdown:    "{{#runtime-import shared/test.md}}",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRuntimeImportFiles(tt.markdown, tmpDir)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCompilerIntegration_RuntimeImportValidation tests the compiler integration
func TestCompilerIntegration_RuntimeImportValidation(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	githubDir := filepath.Join(tmpDir, ".github")
	workflowsDir := filepath.Join(githubDir, "workflows")
	sharedDir := filepath.Join(githubDir, "shared")
	require.NoError(t, os.MkdirAll(workflowsDir, 0755))
	require.NoError(t, os.MkdirAll(sharedDir, 0755))

	// Create a shared file with invalid expression
	sharedFile := filepath.Join(sharedDir, "instructions.md")
	sharedContent := `# Shared Instructions

Use this token: ${{ secrets.GITHUB_TOKEN }}
`
	require.NoError(t, os.WriteFile(sharedFile, []byte(sharedContent), 0644))

	// Create a workflow file that imports the shared file
	workflowFile := filepath.Join(workflowsDir, "test-workflow.md")
	workflowContent := `---
on:
  issues:
    types: [opened]
engine: copilot
---

# Test Workflow

{{#runtime-import ./shared/instructions.md}}

Please process the issue.
`
	require.NoError(t, os.WriteFile(workflowFile, []byte(workflowContent), 0644))

	// Create compiler and attempt to compile
	compiler := NewCompiler()

	err := compiler.CompileWorkflow(workflowFile)

	// Should fail due to invalid expression in runtime-import file
	require.Error(t, err, "Compilation should fail due to invalid expression in runtime-import file")
	assert.Contains(t, err.Error(), "runtime-import files contain expression errors", "Error should mention runtime-import files")
	assert.Contains(t, err.Error(), "secrets.GITHUB_TOKEN", "Error should mention the specific invalid expression")
}

// TestCompilerIntegration_RuntimeImportValidation_Valid tests successful compilation
func TestCompilerIntegration_RuntimeImportValidation_Valid(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	githubDir := filepath.Join(tmpDir, ".github")
	workflowsDir := filepath.Join(githubDir, "workflows")
	sharedDir := filepath.Join(githubDir, "shared")
	require.NoError(t, os.MkdirAll(workflowsDir, 0755))
	require.NoError(t, os.MkdirAll(sharedDir, 0755))

	// Create a shared file with valid expressions
	sharedFile := filepath.Join(sharedDir, "instructions.md")
	sharedContent := `# Shared Instructions

Actor: ${{ github.actor }}
Repository: ${{ github.repository }}
Issue: ${{ github.event.issue.number }}
`
	require.NoError(t, os.WriteFile(sharedFile, []byte(sharedContent), 0644))

	// Create a workflow file that imports the shared file
	workflowFile := filepath.Join(workflowsDir, "test-workflow.md")
	workflowContent := `---
on:
  issues:
    types: [opened]
engine: copilot
---

# Test Workflow

{{#runtime-import ./shared/instructions.md}}

Please process the issue.
`
	require.NoError(t, os.WriteFile(workflowFile, []byte(workflowContent), 0644))

	// Create compiler and compile
	compiler := NewCompiler()

	err := compiler.CompileWorkflow(workflowFile)

	// Should succeed - all expressions are valid
	require.NoError(t, err, "Compilation should succeed with valid expressions in runtime-import file")

	// Clean up lock file if it was created
	if err == nil {
		lockFile := strings.Replace(workflowFile, ".md", ".lock.yml", 1)
		os.Remove(lockFile)
	}
}

// TestValidateRuntimeImportFiles_AgentsFolder tests runtime imports from .agents folder
func TestValidateRuntimeImportFiles_AgentsFolder(t *testing.T) {
	// Create a temporary directory structure with both .github and .agents folders
	tmpDir := t.TempDir()
	githubDir := filepath.Join(tmpDir, ".github")
	agentsDir := filepath.Join(tmpDir, ".agents")
	githubSharedDir := filepath.Join(githubDir, "shared")
	agentsSharedDir := filepath.Join(agentsDir, "shared")
	require.NoError(t, os.MkdirAll(githubSharedDir, 0755))
	require.NoError(t, os.MkdirAll(agentsSharedDir, 0755))

	// Create test files in .github folder
	githubValidFile := filepath.Join(githubSharedDir, "valid.md")
	githubContent := `# GitHub Valid Content

This file has safe expressions:
- Actor: ${{ github.actor }}
- Repository: ${{ github.repository }}
`
	require.NoError(t, os.WriteFile(githubValidFile, []byte(githubContent), 0644))

	// Create test files in .agents folder
	agentsValidFile := filepath.Join(agentsSharedDir, "agent-instructions.md")
	agentsContent := `# Agents Valid Content

Agent instructions with safe expressions:
- Issue number: ${{ github.event.issue.number }}
- Run ID: ${{ github.run_id }}
`
	require.NoError(t, os.WriteFile(agentsValidFile, []byte(agentsContent), 0644))

	agentsInvalidFile := filepath.Join(agentsSharedDir, "agent-invalid.md")
	agentsInvalidContent := `# Agents Invalid Content

This has an unsafe expression:
- Secret: ${{ secrets.MY_TOKEN }}
`
	require.NoError(t, os.WriteFile(agentsInvalidFile, []byte(agentsInvalidContent), 0644))

	tests := []struct {
		name        string
		markdown    string
		expectError bool
		errorText   string
	}{
		{
			name:        "valid import from .github",
			markdown:    "{{#runtime-import .github/shared/valid.md}}",
			expectError: false,
		},
		{
			name:        "valid import from .agents",
			markdown:    "{{#runtime-import .agents/shared/agent-instructions.md}}",
			expectError: false,
		},
		{
			name:        "invalid import from .agents",
			markdown:    "{{#runtime-import .agents/shared/agent-invalid.md}}",
			expectError: true,
			errorText:   "secrets.MY_TOKEN",
		},
		{
			name:        "mixed imports from both folders",
			markdown:    "{{#runtime-import .github/shared/valid.md}}\n{{#runtime-import .agents/shared/agent-instructions.md}}",
			expectError: false,
		},
		{
			name:        "mixed imports with one invalid",
			markdown:    "{{#runtime-import .github/shared/valid.md}}\n{{#runtime-import .agents/shared/agent-invalid.md}}",
			expectError: true,
			errorText:   "secrets.MY_TOKEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRuntimeImportFiles(tt.markdown, tmpDir)

			if tt.expectError {
				require.Error(t, err, "Expected an error")
				if tt.errorText != "" {
					assert.Contains(t, err.Error(), tt.errorText, "Error should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}

// TestValidateRuntimeImportFiles_AgentsFolderSecurity tests security validation for .agents folder
func TestValidateRuntimeImportFiles_AgentsFolderSecurity(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	agentsDir := filepath.Join(tmpDir, ".agents")
	require.NoError(t, os.MkdirAll(agentsDir, 0755))

	// Create a file in .agents folder
	validFile := filepath.Join(agentsDir, "test.md")
	validContent := "# Test\n\nActor: ${{ github.actor }}"
	require.NoError(t, os.WriteFile(validFile, []byte(validContent), 0644))

	tests := []struct {
		name        string
		markdown    string
		expectError bool
		errorText   string
	}{
		{
			name:        "valid path within .agents",
			markdown:    "{{#runtime-import .agents/test.md}}",
			expectError: false,
		},
		{
			name:        "path with ./ prefix in .agents",
			markdown:    "{{#runtime-import .agents/./test.md}}",
			expectError: false,
		},
		{
			name:        "attempt to escape .agents folder",
			markdown:    "{{#runtime-import .agents/../.github/test.md}}",
			expectError: true,
			errorText:   "Security",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRuntimeImportFiles(tt.markdown, tmpDir)

			if tt.expectError {
				require.Error(t, err, "Expected an error")
				if tt.errorText != "" {
					assert.Contains(t, err.Error(), tt.errorText, "Error should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}
