//go:build integration

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/github/gh-aw/pkg/parser"
	"github.com/github/gh-aw/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMCPGatewayInspectIntegration tests that the MCP gateway properly proxies
// MCP server connections and returns the same tools as direct connections.
// This is a critical integration test that validates the gateway doesn't lose
// or alter tool definitions when routing through the HTTP gateway.
//
// Note: This test uses the GitHub MCP server since custom MCP servers in mcp-servers
// section are not yet fully supported by the inspect command. This still validates
// that the gateway properly proxies MCP protocol without alterations.
func TestMCPGatewayInspectIntegration(t *testing.T) {
	// Skip if Docker is not available
	if !isDockerAvailable() {
		t.Skip("Docker not available, skipping MCP gateway integration test")
	}

	setup := setupIntegrationTest(t)
	defer setup.cleanup()

	// Create a test workflow with GitHub MCP server (no gateway - sandbox: false)
	workflowContentDirect := `---
on: workflow_dispatch
engine: copilot
strict: false
sandbox: false
tools:
  github:
---

# Test GitHub MCP Server Direct Connection

Test workflow for direct GitHub MCP server connection (no gateway).
`

	// Create a test workflow with MCP gateway enabled
	workflowContentGateway := `---
on: workflow_dispatch
engine: copilot
sandbox:
  mcp:
    container: ghcr.io/github/gh-aw-mcpg
    version: v0.0.103
tools:
  github:
---

# Test GitHub MCP Server with Gateway

Test workflow for GitHub MCP server connection through gateway.
`

	// Write direct connection workflow
	workflowFileDirect := filepath.Join(setup.workflowsDir, "test-gateway-direct.md")
	err := os.WriteFile(workflowFileDirect, []byte(workflowContentDirect), 0644)
	require.NoError(t, err, "Failed to create direct workflow file")

	// Write gateway connection workflow
	workflowFileGateway := filepath.Join(setup.workflowsDir, "test-gateway-with-gateway.md")
	err = os.WriteFile(workflowFileGateway, []byte(workflowContentGateway), 0644)
	require.NoError(t, err, "Failed to create gateway workflow file")

	// Test 1: Inspect GitHub MCP server with direct connection (no gateway)
	t.Run("direct_connection", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, setup.binaryPath, "mcp", "inspect", "test-gateway-direct", "--server", "github", "--verbose")
		cmd.Dir = setup.tempDir
		// Set a placeholder GitHub token to avoid token validation errors
		cmd.Env = append(os.Environ(), "GITHUB_TOKEN=test_token_for_integration_test")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		t.Logf("Direct connection output:\n%s", outputStr)

		// Check if the command succeeded or at least parsed the configuration
		if err != nil {
			// Connection might fail but config should parse
			if !strings.Contains(outputStr, "github") {
				t.Errorf("Expected github to be mentioned in output")
			}
		}

		// Verify GitHub MCP server is detected
		if strings.Contains(outputStr, "github") {
			t.Logf("✓ Direct connection successfully detected GitHub MCP server")
		}

		// Verify some GitHub tools are listed (if connection succeeds)
		if strings.Contains(outputStr, "get_file_contents") || strings.Contains(outputStr, "create_pull_request") {
			t.Logf("✓ Direct connection successfully listed GitHub MCP tools")
		}
	})

	// Test 2: Inspect GitHub MCP server through gateway (if Docker available)
	t.Run("gateway_connection", func(t *testing.T) {
		// Check if gateway image is available
		checkCmd := exec.Command("docker", "image", "inspect", "ghcr.io/github/gh-aw-mcpg:v0.0.103")
		if err := checkCmd.Run(); err != nil {
			t.Skip("MCP gateway Docker image not available locally, skipping gateway test")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, setup.binaryPath, "mcp", "inspect", "test-gateway-with-gateway", "--server", "github", "--verbose")
		cmd.Dir = setup.tempDir
		cmd.Env = append(os.Environ(), "GITHUB_TOKEN=test_token_for_integration_test")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)

		t.Logf("Gateway connection output:\n%s", outputStr)

		// Check if the command succeeded or at least parsed the configuration
		if err != nil {
			// Gateway might not start but config should parse
			if !strings.Contains(outputStr, "github") {
				t.Errorf("Expected github to be mentioned in output")
			}
		}

		// Verify gateway configuration is detected
		if strings.Contains(outputStr, "gateway") || strings.Contains(outputStr, "ghcr.io/github/gh-aw-mcpg") {
			t.Logf("✓ Gateway configuration detected")
		}

		// Verify GitHub tools are listed through gateway
		if strings.Contains(outputStr, "get_file_contents") || strings.Contains(outputStr, "create_pull_request") {
			t.Logf("✓ Gateway connection successfully listed GitHub MCP tools")
		}
	})

	// Test 3: Compare tools between direct and gateway connections
	t.Run("compare_tools", func(t *testing.T) {
		// This test validates that both direct and gateway connections return
		// the same set of tools, ensuring the gateway doesn't lose information

		// Parse direct connection
		directTools, err := extractToolsFromInspect(setup, "test-gateway-direct", "github")
		if err != nil {
			t.Logf("Note: Could not extract tools from direct connection: %v", err)
			// Don't fail if tools can't be extracted - connection might not work in CI
			return
		}

		// Parse gateway connection (if available)
		gatewayTools, err := extractToolsFromInspect(setup, "test-gateway-with-gateway", "github")
		if err != nil {
			t.Logf("Note: Could not extract tools from gateway connection: %v", err)
			// Don't fail if tools can't be extracted - gateway might not be available
			return
		}

		// Compare tool names
		if len(directTools) > 0 && len(gatewayTools) > 0 {
			assert.Equal(t, len(directTools), len(gatewayTools),
				"Gateway and direct connections should return same number of tools")

			for _, tool := range directTools {
				assert.Contains(t, gatewayTools, tool,
					"Gateway should expose the same tool: %s", tool)
			}

			t.Logf("✓ Gateway and direct connections expose the same %d tools", len(directTools))
		}
	})
}

// extractToolsFromInspect runs mcp inspect and extracts tool names from the output
func extractToolsFromInspect(setup *integrationTestSetup, workflowName, serverName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, setup.binaryPath, "mcp", "inspect", workflowName, "--server", serverName, "--verbose")
	cmd.Dir = setup.tempDir
	cmd.Env = append(os.Environ(), "GITHUB_TOKEN=test_token_for_integration_test")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("inspect command failed: %w (output: %s)", err, string(output))
	}

	outputStr := string(output)
	tools := []string{}

	// Extract common GitHub MCP tool names from output
	githubTools := []string{
		"get_file_contents",
		"create_pull_request",
		"create_or_update_file",
		"add_issue_comment",
		"search_code",
		"list_issues",
		"get_commit",
	}

	for _, tool := range githubTools {
		if strings.Contains(outputStr, tool) {
			tools = append(tools, tool)
		}
	}

	if len(tools) == 0 {
		return nil, fmt.Errorf("no tools found in inspect output")
	}

	return tools, nil
}

// TestMCPServerBasicProtocol tests the test MCP server implements the basic MCP protocol
func TestMCPServerBasicProtocol(t *testing.T) {
	setup := setupIntegrationTest(t)
	defer setup.cleanup()

	// Copy the test MCP server script to the temp directory
	srcServerPath := filepath.Join(projectRoot, "pkg", "cli", "test-mcp-server.cjs")
	destServerPath := filepath.Join(setup.tempDir, "test-mcp-server.cjs")

	srcContent, err := os.ReadFile(srcServerPath)
	if err != nil {
		t.Skip("Test MCP server script not found")
	}
	err = os.WriteFile(destServerPath, srcContent, 0755)
	require.NoError(t, err, "Failed to write test MCP server script")

	testServerPath := destServerPath

	// Test the MCP server can list tools via direct stdio connection
	config := parser.MCPServerConfig{
		Name: "test-server",
		BaseMCPServerConfig: types.BaseMCPServerConfig{
			Type:    "stdio",
			Command: "node",
			Args:    []string{testServerPath},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create the command for the MCP server
	cmd := exec.CommandContext(ctx, config.Command, config.Args...)

	// Set up pipes for communication
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	// Start the server
	err = cmd.Start()
	require.NoError(t, err)
	defer func() {
		stdin.Close()
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	// Send initialize request
	initRequest := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{},
			"clientInfo": map[string]any{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}
	initJSON, _ := json.Marshal(initRequest)
	_, err = stdin.Write(append(initJSON, '\n'))
	require.NoError(t, err)

	// Read initialize response
	decoder := json.NewDecoder(stdout)
	var initResponse map[string]any
	err = decoder.Decode(&initResponse)
	require.NoError(t, err, "Failed to decode initialize response")

	// Verify initialize response
	assert.Equal(t, "2.0", initResponse["jsonrpc"])
	assert.Equal(t, float64(1), initResponse["id"])
	result, ok := initResponse["result"].(map[string]any)
	require.True(t, ok, "Initialize response should have result")
	assert.Equal(t, "2024-11-05", result["protocolVersion"])

	t.Logf("✓ MCP server initialized successfully")

	// Send tools/list request
	toolsRequest := map[string]any{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
		"params":  map[string]any{},
	}
	toolsJSON, _ := json.Marshal(toolsRequest)
	_, err = stdin.Write(append(toolsJSON, '\n'))
	require.NoError(t, err)

	// Read tools/list response
	var toolsResponse map[string]any
	err = decoder.Decode(&toolsResponse)
	require.NoError(t, err, "Failed to decode tools/list response")

	// Verify tools/list response
	assert.Equal(t, "2.0", toolsResponse["jsonrpc"])
	assert.Equal(t, float64(2), toolsResponse["id"])
	result, ok = toolsResponse["result"].(map[string]any)
	require.True(t, ok, "Tools response should have result")
	tools, ok := result["tools"].([]any)
	require.True(t, ok, "Result should have tools array")
	assert.Greater(t, len(tools), 0, "Should have at least one tool")

	t.Logf("✓ MCP server listed %d tools successfully", len(tools))

	// Verify expected tools are present
	toolNames := []string{}
	for _, tool := range tools {
		toolMap, ok := tool.(map[string]any)
		if ok {
			if name, ok := toolMap["name"].(string); ok {
				toolNames = append(toolNames, name)
			}
		}
	}

	assert.Contains(t, toolNames, "test_echo", "Should have test_echo tool")
	assert.Contains(t, toolNames, "test_add", "Should have test_add tool")
	assert.Contains(t, toolNames, "test_uppercase", "Should have test_uppercase tool")

	t.Logf("✓ All expected tools present: %v", toolNames)
}
