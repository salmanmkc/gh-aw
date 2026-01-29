//go:build !integration

package constants

import (
	"path/filepath"
	"testing"
	"time"
)

func TestGetWorkflowDir(t *testing.T) {
	expected := filepath.Join(".github", "workflows")
	result := GetWorkflowDir()

	if result != expected {
		t.Errorf("GetWorkflowDir() = %q, want %q", result, expected)
	}
}

func TestDefaultAllowedDomains(t *testing.T) {
	if len(DefaultAllowedDomains) == 0 {
		t.Error("DefaultAllowedDomains should not be empty")
	}

	expectedDomains := []string{"localhost", "localhost:*", "127.0.0.1", "127.0.0.1:*"}
	if len(DefaultAllowedDomains) != len(expectedDomains) {
		t.Errorf("DefaultAllowedDomains length = %d, want %d", len(DefaultAllowedDomains), len(expectedDomains))
	}

	for i, domain := range expectedDomains {
		if DefaultAllowedDomains[i] != domain {
			t.Errorf("DefaultAllowedDomains[%d] = %q, want %q", i, DefaultAllowedDomains[i], domain)
		}
	}
}

func TestSafeWorkflowEvents(t *testing.T) {
	if len(SafeWorkflowEvents) == 0 {
		t.Error("SafeWorkflowEvents should not be empty")
	}

	// workflow_run is intentionally excluded due to HIGH security risks
	expectedEvents := []string{"workflow_dispatch", "schedule"}
	if len(SafeWorkflowEvents) != len(expectedEvents) {
		t.Errorf("SafeWorkflowEvents length = %d, want %d", len(SafeWorkflowEvents), len(expectedEvents))
	}

	for i, event := range expectedEvents {
		if SafeWorkflowEvents[i] != event {
			t.Errorf("SafeWorkflowEvents[%d] = %q, want %q", i, SafeWorkflowEvents[i], event)
		}
	}
}

func TestAllowedExpressions(t *testing.T) {
	if len(AllowedExpressions) == 0 {
		t.Error("AllowedExpressions should not be empty")
	}

	// Test a few key expressions are present
	requiredExpressions := []string{
		"github.event.issue.number",
		"github.event.pull_request.number",
		"github.repository",
		"github.run_id",
		"github.workspace",
	}

	expressionsMap := make(map[string]bool)
	for _, expr := range AllowedExpressions {
		expressionsMap[expr] = true
	}

	for _, required := range requiredExpressions {
		if !expressionsMap[required] {
			t.Errorf("AllowedExpressions missing required expression: %q", required)
		}
	}
}

func TestAgenticEngines(t *testing.T) {
	if len(AgenticEngines) == 0 {
		t.Error("AgenticEngines should not be empty")
	}

	expectedEngines := []string{"claude", "codex", "copilot"}
	if len(AgenticEngines) != len(expectedEngines) {
		t.Errorf("AgenticEngines length = %d, want %d", len(AgenticEngines), len(expectedEngines))
	}

	for i, engine := range expectedEngines {
		if AgenticEngines[i] != engine {
			t.Errorf("AgenticEngines[%d] = %q, want %q", i, AgenticEngines[i], engine)
		}
	}

	// Verify that engine constants can be converted to strings for AgenticEngines
	if string(ClaudeEngine) != "claude" {
		t.Errorf("ClaudeEngine constant = %q, want %q", ClaudeEngine, "claude")
	}
	if string(CodexEngine) != "codex" {
		t.Errorf("CodexEngine constant = %q, want %q", CodexEngine, "codex")
	}
	if string(CopilotEngine) != "copilot" {
		t.Errorf("CopilotEngine constant = %q, want %q", CopilotEngine, "copilot")
	}
	if string(CustomEngine) != "custom" {
		t.Errorf("CustomEngine constant = %q, want %q", CustomEngine, "custom")
	}
}

func TestDefaultGitHubTools(t *testing.T) {
	if len(DefaultGitHubToolsLocal) == 0 {
		t.Error("DefaultGitHubToolsLocal should not be empty")
	}

	if len(DefaultGitHubToolsRemote) == 0 {
		t.Error("DefaultGitHubToolsRemote should not be empty")
	}

	if len(DefaultReadOnlyGitHubTools) == 0 {
		t.Error("DefaultReadOnlyGitHubTools should not be empty")
	}

	// Test that DefaultGitHubTools defaults to local mode
	if len(DefaultGitHubTools) != len(DefaultGitHubToolsLocal) {
		t.Errorf("DefaultGitHubTools should default to DefaultGitHubToolsLocal")
	}

	// Test that Local and Remote tools reference the same shared list
	if len(DefaultGitHubToolsLocal) != len(DefaultReadOnlyGitHubTools) {
		t.Errorf("DefaultGitHubToolsLocal should have same length as DefaultReadOnlyGitHubTools, got %d vs %d",
			len(DefaultGitHubToolsLocal), len(DefaultReadOnlyGitHubTools))
	}

	if len(DefaultGitHubToolsRemote) != len(DefaultReadOnlyGitHubTools) {
		t.Errorf("DefaultGitHubToolsRemote should have same length as DefaultReadOnlyGitHubTools, got %d vs %d",
			len(DefaultGitHubToolsRemote), len(DefaultReadOnlyGitHubTools))
	}

	// Test a few key tools are present in all lists
	requiredTools := []string{
		"get_me",
		"list_issues",
		"pull_request_read",
		"get_file_contents",
		"search_code",
	}

	for name, tools := range map[string][]string{
		"DefaultGitHubToolsLocal":    DefaultGitHubToolsLocal,
		"DefaultGitHubToolsRemote":   DefaultGitHubToolsRemote,
		"DefaultReadOnlyGitHubTools": DefaultReadOnlyGitHubTools,
	} {
		toolsMap := make(map[string]bool)
		for _, tool := range tools {
			toolsMap[tool] = true
		}

		for _, required := range requiredTools {
			if !toolsMap[required] {
				t.Errorf("%s missing required tool: %q", name, required)
			}
		}
	}
}

func TestDefaultBashTools(t *testing.T) {
	if len(DefaultBashTools) == 0 {
		t.Error("DefaultBashTools should not be empty")
	}

	// Test a few key bash tools are present
	requiredTools := []string{
		"echo",
		"ls",
		"cat",
		"grep",
	}

	toolsMap := make(map[string]bool)
	for _, tool := range DefaultBashTools {
		toolsMap[tool] = true
	}

	for _, required := range requiredTools {
		if !toolsMap[required] {
			t.Errorf("DefaultBashTools missing required tool: %q", required)
		}
	}
}

func TestPriorityFields(t *testing.T) {
	if len(PriorityStepFields) == 0 {
		t.Error("PriorityStepFields should not be empty")
	}

	if len(PriorityJobFields) == 0 {
		t.Error("PriorityJobFields should not be empty")
	}

	if len(PriorityWorkflowFields) == 0 {
		t.Error("PriorityWorkflowFields should not be empty")
	}

	// Test that "name" is first in step fields
	if PriorityStepFields[0] != "name" {
		t.Errorf("PriorityStepFields[0] = %q, want %q", PriorityStepFields[0], "name")
	}

	// Test that "name" is first in job fields
	if PriorityJobFields[0] != "name" {
		t.Errorf("PriorityJobFields[0] = %q, want %q", PriorityJobFields[0], "name")
	}

	// Test that "on" is first in workflow fields
	if PriorityWorkflowFields[0] != "on" {
		t.Errorf("PriorityWorkflowFields[0] = %q, want %q", PriorityWorkflowFields[0], "on")
	}
}

func TestConstantValues(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"CLIExtensionPrefix", string(CLIExtensionPrefix), "gh aw"},
		{"DefaultMCPRegistryURL", string(DefaultMCPRegistryURL), "https://api.mcp.github.com/v0.1"},
		{"AgentJobName", string(AgentJobName), "agent"},
		{"ActivationJobName", string(ActivationJobName), "activation"},
		{"PreActivationJobName", string(PreActivationJobName), "pre_activation"},
		{"DetectionJobName", string(DetectionJobName), "detection"},
		{"SafeOutputArtifactName", string(SafeOutputArtifactName), "safe-output"},
		{"AgentOutputArtifactName", string(AgentOutputArtifactName), "agent-output"},
		{"SafeOutputsMCPServerID", string(SafeOutputsMCPServerID), "safeoutputs"},
		{"SafeInputsMCPServerID", string(SafeInputsMCPServerID), "safeinputs"},
		{"CheckMembershipStepID", string(CheckMembershipStepID), "check_membership"},
		{"CheckStopTimeStepID", string(CheckStopTimeStepID), "check_stop_time"},
		{"CheckSkipIfMatchStepID", string(CheckSkipIfMatchStepID), "check_skip_if_match"},
		{"CheckSkipIfNoMatchStepID", string(CheckSkipIfNoMatchStepID), "check_skip_if_no_match"},
		{"CheckCommandPositionStepID", string(CheckCommandPositionStepID), "check_command_position"},
		{"IsTeamMemberOutput", string(IsTeamMemberOutput), "is_team_member"},
		{"StopTimeOkOutput", string(StopTimeOkOutput), "stop_time_ok"},
		{"SkipCheckOkOutput", string(SkipCheckOkOutput), "skip_check_ok"},
		{"SkipNoMatchCheckOkOutput", string(SkipNoMatchCheckOkOutput), "skip_no_match_check_ok"},
		{"CommandPositionOkOutput", string(CommandPositionOkOutput), "command_position_ok"},
		{"ActivatedOutput", string(ActivatedOutput), "activated"},
		{"DefaultActivationJobRunnerImage", string(DefaultActivationJobRunnerImage), "ubuntu-slim"},
		{"GitHubCopilotMCPDomain", string(GitHubCopilotMCPDomain), "api.githubcopilot.com"},
		{"AgenticCampaignLabel", string(AgenticCampaignLabel), "agentic-campaign"},
		{"CampaignLabelPrefix", string(CampaignLabelPrefix), "z_campaign_"},
		{"DefaultMCPGatewayContainer", string(DefaultMCPGatewayContainer), "ghcr.io/githubnext/gh-aw-mcpg"},
		{"DefaultSerenaMCPServerContainer", string(DefaultSerenaMCPServerContainer), "ghcr.io/githubnext/serena-mcp-server"},
		{"OraiosSerenaContainer", string(OraiosSerenaContainer), "ghcr.io/oraios/serena"},
		{"DefaultNodeAlpineLTSImage", string(DefaultNodeAlpineLTSImage), "node:lts-alpine"},
		{"DefaultPythonAlpineLTSImage", string(DefaultPythonAlpineLTSImage), "python:alpine"},
		{"DefaultAlpineImage", string(DefaultAlpineImage), "alpine:latest"},
		{"DefaultGhAwMount", string(DefaultGhAwMount), "/opt/gh-aw:/opt/gh-aw:ro"},
		{"DefaultTmpGhAwMount", string(DefaultTmpGhAwMount), "/tmp/gh-aw:/tmp/gh-aw:rw"},
		{"DefaultWorkspaceMount", string(DefaultWorkspaceMount), "${{ github.workspace }}:${{ github.workspace }}:rw"},
		{"AgentOutputFilename", string(AgentOutputFilename), "agent_output.json"},
		{"MatchedCommandOutput", string(MatchedCommandOutput), "matched_command"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestModelNameConstants(t *testing.T) {
	// Test that DefaultCopilotDetectionModel has the correct type and value
	tests := []struct {
		name     string
		value    ModelName
		expected string
	}{
		{"DefaultCopilotDetectionModel", DefaultCopilotDetectionModel, "gpt-5.1-codex-mini"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestVersionConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    Version
		expected Version
	}{
		{"DefaultClaudeCodeVersion", DefaultClaudeCodeVersion, "2.1.22"},
		{"DefaultCopilotVersion", DefaultCopilotVersion, "0.0.397"},
		{"DefaultCodexVersion", DefaultCodexVersion, "0.92.0"},
		{"DefaultGitHubMCPServerVersion", DefaultGitHubMCPServerVersion, "v0.30.2"},
		{"DefaultMCPGatewayVersion", DefaultMCPGatewayVersion, "v0.0.84"},
		{"DefaultSandboxRuntimeVersion", DefaultSandboxRuntimeVersion, "0.0.32"},
		{"DefaultFirewallVersion", DefaultFirewallVersion, "v0.11.2"},
		{"DefaultPlaywrightMCPVersion", DefaultPlaywrightMCPVersion, "0.0.61"},
		{"DefaultPlaywrightBrowserVersion", DefaultPlaywrightBrowserVersion, "v1.58.0"},
		{"DefaultBunVersion", DefaultBunVersion, "1.1"},
		{"DefaultNodeVersion", DefaultNodeVersion, "24"},
		{"DefaultPythonVersion", DefaultPythonVersion, "3.12"},
		{"DefaultRubyVersion", DefaultRubyVersion, "3.3"},
		{"DefaultDotNetVersion", DefaultDotNetVersion, "8.0"},
		{"DefaultJavaVersion", DefaultJavaVersion, "21"},
		{"DefaultElixirVersion", DefaultElixirVersion, "1.17"},
		{"DefaultHaskellVersion", DefaultHaskellVersion, "9.10"},
		{"DefaultDenoVersion", DefaultDenoVersion, "2.x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestNumericConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    LineLength
		minValue LineLength
	}{
		{"MaxExpressionLineLength", MaxExpressionLineLength, 1},
		{"ExpressionBreakThreshold", ExpressionBreakThreshold, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value < tt.minValue {
				t.Errorf("%s = %d, should be >= %d", tt.name, tt.value, tt.minValue)
			}
		})
	}
}

func TestTimeoutConstants(t *testing.T) {
	// Test new time.Duration-based constants
	tests := []struct {
		name     string
		value    time.Duration
		minValue time.Duration
	}{
		{"DefaultAgenticWorkflowTimeout", DefaultAgenticWorkflowTimeout, 1 * time.Minute},
		{"DefaultToolTimeout", DefaultToolTimeout, 1 * time.Second},
		{"DefaultMCPStartupTimeout", DefaultMCPStartupTimeout, 1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value < tt.minValue {
				t.Errorf("%s = %v, should be >= %v", tt.name, tt.value, tt.minValue)
			}
		})
	}
}

func TestFeatureFlagConstants(t *testing.T) {
	// Test that feature flag constants have the correct type and values
	tests := []struct {
		name     string
		value    FeatureFlag
		expected string
	}{
		{"SafeInputsFeatureFlag", SafeInputsFeatureFlag, "safe-inputs"},
		{"MCPGatewayFeatureFlag", MCPGatewayFeatureFlag, "mcp-gateway"},
		{"SandboxRuntimeFeatureFlag", SandboxRuntimeFeatureFlag, "sandbox-runtime"},
		{"DangerousPermissionsWriteFeatureFlag", DangerousPermissionsWriteFeatureFlag, "dangerous-permissions-write"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.value) != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestFeatureFlagType(t *testing.T) {
	// Test that FeatureFlag type can be used as expected
	var flag FeatureFlag = "test-flag"
	if string(flag) != "test-flag" {
		t.Errorf("FeatureFlag conversion failed: got %q, want %q", flag, "test-flag")
	}

	// Test that constants can be assigned to FeatureFlag variables
	safeInputsFlag := SafeInputsFeatureFlag
	if safeInputsFlag != "safe-inputs" {
		t.Errorf("SafeInputsFeatureFlag assignment failed: got %q, want %q", safeInputsFlag, "safe-inputs")
	}
}

func TestSemanticTypeAliases(t *testing.T) {
	// Test URL type
	t.Run("URL type", func(t *testing.T) {
		var testURL URL = "https://example.com"
		if string(testURL) != "https://example.com" {
			t.Errorf("URL conversion failed: got %q, want %q", testURL, "https://example.com")
		}

		// Test DefaultMCPRegistryURL has the correct type
		registryURL := DefaultMCPRegistryURL
		if string(registryURL) != "https://api.mcp.github.com/v0.1" {
			t.Errorf("DefaultMCPRegistryURL = %q, want %q", registryURL, "https://api.mcp.github.com/v0.1")
		}
	})

	// Test ModelName type
	t.Run("ModelName type", func(t *testing.T) {
		var testModel ModelName = "test-model"
		if string(testModel) != "test-model" {
			t.Errorf("ModelName conversion failed: got %q, want %q", testModel, "test-model")
		}

		// Test DefaultCopilotDetectionModel has the correct type
		detectionModel := DefaultCopilotDetectionModel
		if string(detectionModel) != "gpt-5.1-codex-mini" {
			t.Errorf("DefaultCopilotDetectionModel = %q, want %q", detectionModel, "gpt-5.1-codex-mini")
		}
	})

	// Test JobName type
	t.Run("JobName type", func(t *testing.T) {
		var testJob JobName = "test-job"
		if string(testJob) != "test-job" {
			t.Errorf("JobName conversion failed: got %q, want %q", testJob, "test-job")
		}

		// Test job name constants have the correct type
		agentJob := AgentJobName
		if string(agentJob) != "agent" {
			t.Errorf("AgentJobName = %q, want %q", agentJob, "agent")
		}

		activationJob := ActivationJobName
		if string(activationJob) != "activation" {
			t.Errorf("ActivationJobName = %q, want %q", activationJob, "activation")
		}

		preActivationJob := PreActivationJobName
		if string(preActivationJob) != "pre_activation" {
			t.Errorf("PreActivationJobName = %q, want %q", preActivationJob, "pre_activation")
		}

		detectionJob := DetectionJobName
		if string(detectionJob) != "detection" {
			t.Errorf("DetectionJobName = %q, want %q", detectionJob, "detection")
		}
	})

	// Test StepID type
	t.Run("StepID type", func(t *testing.T) {
		var testStep StepID = "test-step"
		if string(testStep) != "test-step" {
			t.Errorf("StepID conversion failed: got %q, want %q", testStep, "test-step")
		}

		// Test step ID constants have the correct type
		membershipStep := CheckMembershipStepID
		if string(membershipStep) != "check_membership" {
			t.Errorf("CheckMembershipStepID = %q, want %q", membershipStep, "check_membership")
		}

		stopTimeStep := CheckStopTimeStepID
		if string(stopTimeStep) != "check_stop_time" {
			t.Errorf("CheckStopTimeStepID = %q, want %q", stopTimeStep, "check_stop_time")
		}

		skipMatchStep := CheckSkipIfMatchStepID
		if string(skipMatchStep) != "check_skip_if_match" {
			t.Errorf("CheckSkipIfMatchStepID = %q, want %q", skipMatchStep, "check_skip_if_match")
		}

		commandPosStep := CheckCommandPositionStepID
		if string(commandPosStep) != "check_command_position" {
			t.Errorf("CheckCommandPositionStepID = %q, want %q", commandPosStep, "check_command_position")
		}
	})

	// Test CommandPrefix type
	t.Run("CommandPrefix type", func(t *testing.T) {
		var testPrefix CommandPrefix = "test-prefix"
		if string(testPrefix) != "test-prefix" {
			t.Errorf("CommandPrefix conversion failed: got %q, want %q", testPrefix, "test-prefix")
		}

		// Test CLIExtensionPrefix has the correct type
		cliPrefix := CLIExtensionPrefix
		if string(cliPrefix) != "gh aw" {
			t.Errorf("CLIExtensionPrefix = %q, want %q", cliPrefix, "gh aw")
		}
	})

	// Test WorkflowID type
	t.Run("WorkflowID type", func(t *testing.T) {
		var testWorkflow WorkflowID = "ci-doctor"
		if string(testWorkflow) != "ci-doctor" {
			t.Errorf("WorkflowID conversion failed: got %q, want %q", testWorkflow, "ci-doctor")
		}

		// Test that WorkflowID can hold typical workflow identifiers
		workflows := []WorkflowID{"ci-doctor", "deploy-prod", "test-workflow"}
		for i, wf := range workflows {
			if !wf.IsValid() {
				t.Errorf("WorkflowID[%d] should be valid: %q", i, wf)
			}
		}
	})

	// Test EngineName type
	t.Run("EngineName type", func(t *testing.T) {
		var testEngine EngineName = "copilot"
		if string(testEngine) != "copilot" {
			t.Errorf("EngineName conversion failed: got %q, want %q", testEngine, "copilot")
		}

		// Test engine constants have the correct type
		copilot := CopilotEngine
		if string(copilot) != "copilot" {
			t.Errorf("CopilotEngine = %q, want %q", copilot, "copilot")
		}

		claude := ClaudeEngine
		if string(claude) != "claude" {
			t.Errorf("ClaudeEngine = %q, want %q", claude, "claude")
		}

		codex := CodexEngine
		if string(codex) != "codex" {
			t.Errorf("CodexEngine = %q, want %q", codex, "codex")
		}

		custom := CustomEngine
		if string(custom) != "custom" {
			t.Errorf("CustomEngine = %q, want %q", custom, "custom")
		}
	})

	// Test FilePath type
	t.Run("FilePath type", func(t *testing.T) {
		var testPath FilePath = "/opt/gh-aw/safe-inputs"
		if string(testPath) != "/opt/gh-aw/safe-inputs" {
			t.Errorf("FilePath conversion failed: got %q, want %q", testPath, "/opt/gh-aw/safe-inputs")
		}
	})

	// Test ContainerImage type
	t.Run("ContainerImage type", func(t *testing.T) {
		var testImage ContainerImage = "ghcr.io/test/image"
		if string(testImage) != "ghcr.io/test/image" {
			t.Errorf("ContainerImage conversion failed: got %q, want %q", testImage, "ghcr.io/test/image")
		}

		// Test container image constants have the correct type
		gateway := DefaultMCPGatewayContainer
		if string(gateway) != "ghcr.io/githubnext/gh-aw-mcpg" {
			t.Errorf("DefaultMCPGatewayContainer = %q, want %q", gateway, "ghcr.io/githubnext/gh-aw-mcpg")
		}

		serena := DefaultSerenaMCPServerContainer
		if string(serena) != "ghcr.io/githubnext/serena-mcp-server" {
			t.Errorf("DefaultSerenaMCPServerContainer = %q, want %q", serena, "ghcr.io/githubnext/serena-mcp-server")
		}

		oraios := OraiosSerenaContainer
		if string(oraios) != "ghcr.io/oraios/serena" {
			t.Errorf("OraiosSerenaContainer = %q, want %q", oraios, "ghcr.io/oraios/serena")
		}

		node := DefaultNodeAlpineLTSImage
		if string(node) != "node:lts-alpine" {
			t.Errorf("DefaultNodeAlpineLTSImage = %q, want %q", node, "node:lts-alpine")
		}

		python := DefaultPythonAlpineLTSImage
		if string(python) != "python:alpine" {
			t.Errorf("DefaultPythonAlpineLTSImage = %q, want %q", python, "python:alpine")
		}

		alpine := DefaultAlpineImage
		if string(alpine) != "alpine:latest" {
			t.Errorf("DefaultAlpineImage = %q, want %q", alpine, "alpine:latest")
		}
	})

	// Test Domain type
	t.Run("Domain type", func(t *testing.T) {
		var testDomain Domain = "example.com"
		if string(testDomain) != "example.com" {
			t.Errorf("Domain conversion failed: got %q, want %q", testDomain, "example.com")
		}

		// Test domain constant has the correct type
		copilotDomain := GitHubCopilotMCPDomain
		if string(copilotDomain) != "api.githubcopilot.com" {
			t.Errorf("GitHubCopilotMCPDomain = %q, want %q", copilotDomain, "api.githubcopilot.com")
		}
	})

	// Test Label type
	t.Run("Label type", func(t *testing.T) {
		var testLabel Label = "test-label"
		if string(testLabel) != "test-label" {
			t.Errorf("Label conversion failed: got %q, want %q", testLabel, "test-label")
		}

		// Test label constants have the correct type
		campaignLabel := AgenticCampaignLabel
		if string(campaignLabel) != "agentic-campaign" {
			t.Errorf("AgenticCampaignLabel = %q, want %q", campaignLabel, "agentic-campaign")
		}

		labelPrefix := CampaignLabelPrefix
		if string(labelPrefix) != "z_campaign_" {
			t.Errorf("CampaignLabelPrefix = %q, want %q", labelPrefix, "z_campaign_")
		}
	})

	// Test ArtifactName type
	t.Run("ArtifactName type", func(t *testing.T) {
		var testArtifact ArtifactName = "test-artifact"
		if string(testArtifact) != "test-artifact" {
			t.Errorf("ArtifactName conversion failed: got %q, want %q", testArtifact, "test-artifact")
		}

		// Test artifact name constants have the correct type
		safeOutput := SafeOutputArtifactName
		if string(safeOutput) != "safe-output" {
			t.Errorf("SafeOutputArtifactName = %q, want %q", safeOutput, "safe-output")
		}

		agentOutput := AgentOutputArtifactName
		if string(agentOutput) != "agent-output" {
			t.Errorf("AgentOutputArtifactName = %q, want %q", agentOutput, "agent-output")
		}
	})

	// Test FileName type
	t.Run("FileName type", func(t *testing.T) {
		var testFile FileName = "test.json"
		if string(testFile) != "test.json" {
			t.Errorf("FileName conversion failed: got %q, want %q", testFile, "test.json")
		}

		// Test filename constant has the correct type
		agentFile := AgentOutputFilename
		if string(agentFile) != "agent_output.json" {
			t.Errorf("AgentOutputFilename = %q, want %q", agentFile, "agent_output.json")
		}
	})

	// Test ServerID type
	t.Run("ServerID type", func(t *testing.T) {
		var testServer ServerID = "testserver"
		if string(testServer) != "testserver" {
			t.Errorf("ServerID conversion failed: got %q, want %q", testServer, "testserver")
		}

		// Test server ID constants have the correct type
		safeOutputsServer := SafeOutputsMCPServerID
		if string(safeOutputsServer) != "safeoutputs" {
			t.Errorf("SafeOutputsMCPServerID = %q, want %q", safeOutputsServer, "safeoutputs")
		}

		safeInputsServer := SafeInputsMCPServerID
		if string(safeInputsServer) != "safeinputs" {
			t.Errorf("SafeInputsMCPServerID = %q, want %q", safeInputsServer, "safeinputs")
		}
	})

	// Test OutputName type
	t.Run("OutputName type", func(t *testing.T) {
		var testOutput OutputName = "test_output"
		if string(testOutput) != "test_output" {
			t.Errorf("OutputName conversion failed: got %q, want %q", testOutput, "test_output")
		}

		// Test output name constants have the correct type
		teamMember := IsTeamMemberOutput
		if string(teamMember) != "is_team_member" {
			t.Errorf("IsTeamMemberOutput = %q, want %q", teamMember, "is_team_member")
		}

		stopTime := StopTimeOkOutput
		if string(stopTime) != "stop_time_ok" {
			t.Errorf("StopTimeOkOutput = %q, want %q", stopTime, "stop_time_ok")
		}

		skipCheck := SkipCheckOkOutput
		if string(skipCheck) != "skip_check_ok" {
			t.Errorf("SkipCheckOkOutput = %q, want %q", skipCheck, "skip_check_ok")
		}

		skipNoMatch := SkipNoMatchCheckOkOutput
		if string(skipNoMatch) != "skip_no_match_check_ok" {
			t.Errorf("SkipNoMatchCheckOkOutput = %q, want %q", skipNoMatch, "skip_no_match_check_ok")
		}

		commandPos := CommandPositionOkOutput
		if string(commandPos) != "command_position_ok" {
			t.Errorf("CommandPositionOkOutput = %q, want %q", commandPos, "command_position_ok")
		}

		matchedCmd := MatchedCommandOutput
		if string(matchedCmd) != "matched_command" {
			t.Errorf("MatchedCommandOutput = %q, want %q", matchedCmd, "matched_command")
		}

		activated := ActivatedOutput
		if string(activated) != "activated" {
			t.Errorf("ActivatedOutput = %q, want %q", activated, "activated")
		}
	})

	// Test RunnerImage type
	t.Run("RunnerImage type", func(t *testing.T) {
		var testRunner RunnerImage = "ubuntu-latest"
		if string(testRunner) != "ubuntu-latest" {
			t.Errorf("RunnerImage conversion failed: got %q, want %q", testRunner, "ubuntu-latest")
		}

		// Test runner image constant has the correct type
		activationRunner := DefaultActivationJobRunnerImage
		if string(activationRunner) != "ubuntu-slim" {
			t.Errorf("DefaultActivationJobRunnerImage = %q, want %q", activationRunner, "ubuntu-slim")
		}
	})

	// Test MountPath type
	t.Run("MountPath type", func(t *testing.T) {
		var testMount MountPath = "/src:/dst:ro"
		if string(testMount) != "/src:/dst:ro" {
			t.Errorf("MountPath conversion failed: got %q, want %q", testMount, "/src:/dst:ro")
		}

		// Test mount path constants have the correct type
		ghAwMount := DefaultGhAwMount
		if string(ghAwMount) != "/opt/gh-aw:/opt/gh-aw:ro" {
			t.Errorf("DefaultGhAwMount = %q, want %q", ghAwMount, "/opt/gh-aw:/opt/gh-aw:ro")
		}

		tmpMount := DefaultTmpGhAwMount
		if string(tmpMount) != "/tmp/gh-aw:/tmp/gh-aw:rw" {
			t.Errorf("DefaultTmpGhAwMount = %q, want %q", tmpMount, "/tmp/gh-aw:/tmp/gh-aw:rw")
		}

		workspaceMount := DefaultWorkspaceMount
		if string(workspaceMount) != "${{ github.workspace }}:${{ github.workspace }}:rw" {
			t.Errorf("DefaultWorkspaceMount = %q, want %q", workspaceMount, "${{ github.workspace }}:${{ github.workspace }}:rw")
		}
	})
}

func TestTypeSafetyBetweenSemanticTypes(t *testing.T) {
	// This test demonstrates that semantic types provide type safety
	// by preventing accidental mixing of different string types

	// These assignments should work (same types)
	job1 := AgentJobName
	job2 := ActivationJobName
	if job1 == job2 {
		t.Error("AgentJobName should not equal ActivationJobName")
	}

	step1 := CheckMembershipStepID
	step2 := CheckStopTimeStepID
	if step1 == step2 {
		t.Error("CheckMembershipStepID should not equal CheckStopTimeStepID")
	}

	// Verify that we can still convert to string when needed
	if string(job1) != "agent" {
		t.Errorf("JobName string conversion failed: got %q, want %q", job1, "agent")
	}

	if string(step1) != "check_membership" {
		t.Errorf("StepID string conversion failed: got %q, want %q", step1, "check_membership")
	}

	// Verify that different semantic types have different underlying types
	// (this is a compile-time check, but we verify the values are correct)
	jobStr := string(AgentJobName)
	stepStr := string(CheckMembershipStepID)
	_ = jobStr  // Used for demonstration
	_ = stepStr // Used for demonstration
	// Different semantic types prevent accidental mixing even if string values match
}

// TestHelperMethods tests the helper methods on semantic types
func TestHelperMethods(t *testing.T) {
	t.Run("LineLength", func(t *testing.T) {
		length := LineLength(120)
		if length.String() != "120" {
			t.Errorf("LineLength.String() = %q, want %q", length.String(), "120")
		}
		if !length.IsValid() {
			t.Error("LineLength.IsValid() = false, want true for positive value")
		}

		invalidLength := LineLength(0)
		if invalidLength.IsValid() {
			t.Error("LineLength.IsValid() = true, want false for zero value")
		}

		negativeLength := LineLength(-1)
		if negativeLength.IsValid() {
			t.Error("LineLength.IsValid() = true, want false for negative value")
		}
	})

	t.Run("Version", func(t *testing.T) {
		version := Version("1.0.0")
		if version.String() != "1.0.0" {
			t.Errorf("Version.String() = %q, want %q", version.String(), "1.0.0")
		}
		if !version.IsValid() {
			t.Error("Version.IsValid() = false, want true for non-empty value")
		}

		emptyVersion := Version("")
		if emptyVersion.IsValid() {
			t.Error("Version.IsValid() = true, want false for empty value")
		}
	})

	t.Run("FeatureFlag", func(t *testing.T) {
		flag := FeatureFlag("test-flag")
		if flag.String() != "test-flag" {
			t.Errorf("FeatureFlag.String() = %q, want %q", flag.String(), "test-flag")
		}
		if !flag.IsValid() {
			t.Error("FeatureFlag.IsValid() = false, want true for non-empty value")
		}

		emptyFlag := FeatureFlag("")
		if emptyFlag.IsValid() {
			t.Error("FeatureFlag.IsValid() = true, want false for empty value")
		}
	})

	t.Run("URL", func(t *testing.T) {
		url := URL("https://example.com")
		if url.String() != "https://example.com" {
			t.Errorf("URL.String() = %q, want %q", url.String(), "https://example.com")
		}
		if !url.IsValid() {
			t.Error("URL.IsValid() = false, want true for non-empty value")
		}

		emptyURL := URL("")
		if emptyURL.IsValid() {
			t.Error("URL.IsValid() = true, want false for empty value")
		}
	})

	t.Run("ModelName", func(t *testing.T) {
		model := ModelName("gpt-5-mini")
		if model.String() != "gpt-5-mini" {
			t.Errorf("ModelName.String() = %q, want %q", model.String(), "gpt-5-mini")
		}
		if !model.IsValid() {
			t.Error("ModelName.IsValid() = false, want true for non-empty value")
		}

		emptyModel := ModelName("")
		if emptyModel.IsValid() {
			t.Error("ModelName.IsValid() = true, want false for empty value")
		}
	})

	t.Run("JobName", func(t *testing.T) {
		job := JobName("agent")
		if job.String() != "agent" {
			t.Errorf("JobName.String() = %q, want %q", job.String(), "agent")
		}
		if !job.IsValid() {
			t.Error("JobName.IsValid() = false, want true for non-empty value")
		}

		emptyJob := JobName("")
		if emptyJob.IsValid() {
			t.Error("JobName.IsValid() = true, want false for empty value")
		}
	})

	t.Run("StepID", func(t *testing.T) {
		step := StepID("check_membership")
		if step.String() != "check_membership" {
			t.Errorf("StepID.String() = %q, want %q", step.String(), "check_membership")
		}
		if !step.IsValid() {
			t.Error("StepID.IsValid() = false, want true for non-empty value")
		}

		emptyStep := StepID("")
		if emptyStep.IsValid() {
			t.Error("StepID.IsValid() = true, want false for empty value")
		}
	})

	t.Run("CommandPrefix", func(t *testing.T) {
		prefix := CommandPrefix("gh aw")
		if prefix.String() != "gh aw" {
			t.Errorf("CommandPrefix.String() = %q, want %q", prefix.String(), "gh aw")
		}
		if !prefix.IsValid() {
			t.Error("CommandPrefix.IsValid() = false, want true for non-empty value")
		}

		emptyPrefix := CommandPrefix("")
		if emptyPrefix.IsValid() {
			t.Error("CommandPrefix.IsValid() = true, want false for empty value")
		}
	})

	t.Run("WorkflowID", func(t *testing.T) {
		workflow := WorkflowID("ci-doctor")
		if workflow.String() != "ci-doctor" {
			t.Errorf("WorkflowID.String() = %q, want %q", workflow.String(), "ci-doctor")
		}
		if !workflow.IsValid() {
			t.Error("WorkflowID.IsValid() = false, want true for non-empty value")
		}

		emptyWorkflow := WorkflowID("")
		if emptyWorkflow.IsValid() {
			t.Error("WorkflowID.IsValid() = true, want false for empty value")
		}
	})

	t.Run("EngineName", func(t *testing.T) {
		engine := EngineName("copilot")
		if engine.String() != "copilot" {
			t.Errorf("EngineName.String() = %q, want %q", engine.String(), "copilot")
		}
		if !engine.IsValid() {
			t.Error("EngineName.IsValid() = false, want true for non-empty value")
		}

		emptyEngine := EngineName("")
		if emptyEngine.IsValid() {
			t.Error("EngineName.IsValid() = true, want false for empty value")
		}
	})

	t.Run("FilePath", func(t *testing.T) {
		path := FilePath("/opt/gh-aw")
		if path.String() != "/opt/gh-aw" {
			t.Errorf("FilePath.String() = %q, want %q", path.String(), "/opt/gh-aw")
		}
		if !path.IsValid() {
			t.Error("FilePath.IsValid() = false, want true for non-empty value")
		}

		emptyPath := FilePath("")
		if emptyPath.IsValid() {
			t.Error("FilePath.IsValid() = true, want false for empty value")
		}
	})

	t.Run("ContainerImage", func(t *testing.T) {
		image := ContainerImage("ghcr.io/test/image")
		if image.String() != "ghcr.io/test/image" {
			t.Errorf("ContainerImage.String() = %q, want %q", image.String(), "ghcr.io/test/image")
		}
		if !image.IsValid() {
			t.Error("ContainerImage.IsValid() = false, want true for non-empty value")
		}

		emptyImage := ContainerImage("")
		if emptyImage.IsValid() {
			t.Error("ContainerImage.IsValid() = true, want false for empty value")
		}
	})

	t.Run("Domain", func(t *testing.T) {
		domain := Domain("example.com")
		if domain.String() != "example.com" {
			t.Errorf("Domain.String() = %q, want %q", domain.String(), "example.com")
		}
		if !domain.IsValid() {
			t.Error("Domain.IsValid() = false, want true for non-empty value")
		}

		emptyDomain := Domain("")
		if emptyDomain.IsValid() {
			t.Error("Domain.IsValid() = true, want false for empty value")
		}
	})

	t.Run("Label", func(t *testing.T) {
		label := Label("test-label")
		if label.String() != "test-label" {
			t.Errorf("Label.String() = %q, want %q", label.String(), "test-label")
		}
		if !label.IsValid() {
			t.Error("Label.IsValid() = false, want true for non-empty value")
		}

		emptyLabel := Label("")
		if emptyLabel.IsValid() {
			t.Error("Label.IsValid() = true, want false for empty value")
		}
	})

	t.Run("ArtifactName", func(t *testing.T) {
		artifact := ArtifactName("test-artifact")
		if artifact.String() != "test-artifact" {
			t.Errorf("ArtifactName.String() = %q, want %q", artifact.String(), "test-artifact")
		}
		if !artifact.IsValid() {
			t.Error("ArtifactName.IsValid() = false, want true for non-empty value")
		}

		emptyArtifact := ArtifactName("")
		if emptyArtifact.IsValid() {
			t.Error("ArtifactName.IsValid() = true, want false for empty value")
		}
	})

	t.Run("FileName", func(t *testing.T) {
		filename := FileName("test.json")
		if filename.String() != "test.json" {
			t.Errorf("FileName.String() = %q, want %q", filename.String(), "test.json")
		}
		if !filename.IsValid() {
			t.Error("FileName.IsValid() = false, want true for non-empty value")
		}

		emptyFilename := FileName("")
		if emptyFilename.IsValid() {
			t.Error("FileName.IsValid() = true, want false for empty value")
		}
	})

	t.Run("ServerID", func(t *testing.T) {
		server := ServerID("testserver")
		if server.String() != "testserver" {
			t.Errorf("ServerID.String() = %q, want %q", server.String(), "testserver")
		}
		if !server.IsValid() {
			t.Error("ServerID.IsValid() = false, want true for non-empty value")
		}

		emptyServer := ServerID("")
		if emptyServer.IsValid() {
			t.Error("ServerID.IsValid() = true, want false for empty value")
		}
	})

	t.Run("OutputName", func(t *testing.T) {
		output := OutputName("test_output")
		if output.String() != "test_output" {
			t.Errorf("OutputName.String() = %q, want %q", output.String(), "test_output")
		}
		if !output.IsValid() {
			t.Error("OutputName.IsValid() = false, want true for non-empty value")
		}

		emptyOutput := OutputName("")
		if emptyOutput.IsValid() {
			t.Error("OutputName.IsValid() = true, want false for empty value")
		}
	})

	t.Run("RunnerImage", func(t *testing.T) {
		runner := RunnerImage("ubuntu-latest")
		if runner.String() != "ubuntu-latest" {
			t.Errorf("RunnerImage.String() = %q, want %q", runner.String(), "ubuntu-latest")
		}
		if !runner.IsValid() {
			t.Error("RunnerImage.IsValid() = false, want true for non-empty value")
		}

		emptyRunner := RunnerImage("")
		if emptyRunner.IsValid() {
			t.Error("RunnerImage.IsValid() = true, want false for empty value")
		}
	})

	t.Run("MountPath", func(t *testing.T) {
		mount := MountPath("/src:/dst:ro")
		if mount.String() != "/src:/dst:ro" {
			t.Errorf("MountPath.String() = %q, want %q", mount.String(), "/src:/dst:ro")
		}
		if !mount.IsValid() {
			t.Error("MountPath.IsValid() = false, want true for non-empty value")
		}

		emptyMount := MountPath("")
		if emptyMount.IsValid() {
			t.Error("MountPath.IsValid() = true, want false for empty value")
		}
	})
}
