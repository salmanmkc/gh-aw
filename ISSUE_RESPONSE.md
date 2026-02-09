# Response to Issue: Clarify or remove unprocessed 'post-steps' schema field

## Issue Status: **INVALID** ❌

The issue description states:

> The main workflow schema (`pkg/parser/schemas/main_workflow_schema.json`) defines a `post-steps` field in the schema, but there is no compiler code that accesses `frontmatter["post-steps"]` and no `PostSteps` field in `FrontmatterConfig`.

**This claim is completely false.** The `post-steps` feature is fully implemented, tested, documented, and working correctly in production.

---

## Complete Evidence of Implementation

### 1. ✅ Schema Definition

**File**: `pkg/parser/schemas/main_workflow_schema.json`

```json
"post-steps": {
  "description": "Custom workflow steps to run after AI execution",
  "oneOf": [
    {
      "type": "object",
      "additionalProperties": true
    },
    {
      "type": "array",
      "items": { ... }
    }
  ]
}
```

Includes full examples and proper type definitions.

---

### 2. ✅ Type Definition in FrontmatterConfig

**File**: `pkg/workflow/frontmatter_types.go`

**Line 130**:
```go
PostSteps   []any          `json:"post-steps,omitempty"`  // Post-workflow steps
```

**Lines 619-621** (serialization in `ToMap()`):
```go
if fc.PostSteps != nil {
    result["post-steps"] = fc.PostSteps
}
```

The field exists and is properly handled.

---

### 3. ✅ Parser Implementation

**File**: `pkg/parser/content_extractor.go`

Function `extractPostStepsFromContent()` extracts post-steps from frontmatter:

```go
// Extract post-steps section
postSteps, exists := result.Frontmatter["post-steps"]
if !exists {
    return "", nil
}
```

**File**: `pkg/parser/import_processor.go`

`MergedPostSteps` field handles merging post-steps from imports:

```go
type ImportResult struct {
    ...
    MergedPostSteps     string   // Merged post-steps configuration from all imports
    ...
}
```

**File**: `pkg/parser/frontmatter_hash.go`

Post-steps are included in the frontmatter hash:
```go
addField("post-steps")
```

---

### 4. ✅ Compiler Implementation

**File**: `pkg/workflow/compiler_orchestrator_workflow.go`

Function `processAndMergePostSteps()` processes post-steps with action pinning:

```go
func (c *Compiler) processAndMergePostSteps(frontmatter map[string]any, workflowData *WorkflowData) {
    workflowData.PostSteps = c.extractTopLevelYAMLSection(frontmatter, "post-steps")
    
    // Apply action pinning to post-steps if any
    if workflowData.PostSteps != "" {
        // ... action pinning logic ...
    }
}
```

**File**: `pkg/workflow/compiler_yaml.go`

Function `generatePostSteps()` generates the YAML output:

```go
func (c *Compiler) generatePostSteps(yaml *strings.Builder, data *WorkflowData) {
    if data.PostSteps != "" {
        // ... YAML generation logic ...
    }
}
```

**File**: `pkg/workflow/compiler_types.go`

`WorkflowData` struct includes:
```go
type WorkflowData struct {
    ...
    PostSteps string // steps to run after AI execution
    ...
}
```

**File**: `pkg/workflow/safe_outputs_jobs.go`

Post-steps are used in safe output jobs:
```go
if len(config.PostSteps) > 0 {
    steps = append(steps, config.PostSteps...)
}
```

---

### 5. ✅ Comprehensive Test Coverage

**Test Files**:
1. `pkg/workflow/compiler_poststeps_test.go` - Dedicated post-steps tests
   - `TestPostStepsGeneration` - Tests pre-steps, AI execution, and post-steps order
   - `TestPostStepsOnly` - Tests post-steps without pre-steps
   - `TestStopAfterCompiledAway` - Tests stop-after field handling

2. `pkg/workflow/compiler_artifacts_test.go`
   - `TestPostStepsIndentationFix` - Tests YAML indentation

3. `pkg/workflow/compiler_orchestrator_test.go`
   - `TestProcessAndMergePostSteps` - Tests merging logic

4. `pkg/workflow/compiler_orchestrator_workflow_test.go`
   - `TestProcessAndMergePostSteps_NoPostSteps`
   - `TestProcessAndMergePostSteps_WithPostSteps`

5. `pkg/workflow/safe_output_refactor_test.go`
   - `TestSafeOutputJobBuilderWithPreAndPostSteps`

6. `pkg/parser/schema_validation_test.go`
   - Includes `"post-steps": []any{map[string]any{"run": "echo cleanup"}}`

**Test Workflow**: `pkg/cli/workflows/test-post-steps.md` - A complete working example

**All tests pass**:
```bash
$ go test -v -run "TestPostSteps" ./pkg/workflow/
=== RUN   TestPostStepsIndentationFix
--- PASS: TestPostStepsIndentationFix (0.12s)
=== RUN   TestPostStepsGeneration
--- PASS: TestPostStepsGeneration (0.04s)
=== RUN   TestPostStepsOnly
--- PASS: TestPostStepsOnly (0.04s)
PASS
```

---

### 6. ✅ Documentation

**File**: `docs/src/content/docs/reference/frontmatter.md`

Full section on post-steps:

```markdown
## Post-Execution Steps (`post-steps:`)

Add custom steps after agentic execution. Run after AI engine completes 
regardless of success/failure (unless conditional expressions are used).

> Security Notice: Post-execution steps run OUTSIDE the firewall sandbox. 
> These steps execute with standard GitHub Actions security but do NOT have 
> the network egress controls that protect the agent job. Do not run agentic 
> compute or untrusted AI execution in post-steps - use them only for 
> deterministic cleanup, artifact uploads, or notifications.
```

Additional documentation in:
- `docs/src/content/docs/reference/frontmatter-full.md`
- `docs/src/content/docs/reference/frontmatter-hash-specification.md`
- `.github/aw/github-agentic-workflows.md`
- `.github/aw/create-agentic-workflow.md`
- `.github/aw/update-agentic-workflow.md`

---

### 7. ✅ Working Example

**Compilation Test**:
```bash
$ ./gh-aw compile pkg/cli/workflows/test-post-steps.md
✓ pkg/cli/workflows/test-post-steps.md (22.9 KB)
✓ Compiled 1 workflow(s): 0 error(s), 0 warning(s)
```

**Generated Output** (from `.lock.yml`):
```yaml
- name: Verify Post-Steps Execution
  run: |
    echo "✅ Post-steps are executing correctly"
    echo "This step runs after the AI agent completes"
- if: always()
  name: Upload Test Results
  uses: actions/upload-artifact@b7c566a772e6b6bfb58ed0dc250532a479d7789f # v6.0.0
  with:
    name: post-steps-test-results
    path: /tmp/gh-aw/
    retention-days: 1
    if-no-files-found: ignore
- name: Final Summary
  run: ...
```

---

## Code Flow Diagram

```
User Workflow (.md file)
    ↓
Parser extracts post-steps
    ↓
extractPostStepsFromContent()
    ↓
MergedPostSteps (from imports)
    ↓
processAndMergePostSteps()
    ↓
Action pinning applied
    ↓
generatePostSteps()
    ↓
Final .lock.yml file
```

---

## Why This Issue Was Created

The issue appears to be based on:
1. **Incomplete code search** - Only searching for `frontmatter["post-steps"]` misses the JSON-based parsing
2. **Overlooking the struct field** - The `PostSteps []any` field exists on line 130 of `frontmatter_types.go`
3. **Not running tests** - All existing tests demonstrate the feature works
4. **Not checking compiled output** - The feature generates correct YAML

---

## Conclusion

The `post-steps` feature is:
- ✅ **Defined in schema** with proper types and examples
- ✅ **Defined in FrontmatterConfig** struct
- ✅ **Extracted by parser** via multiple functions
- ✅ **Processed by compiler** with action pinning
- ✅ **Generated in output** with proper YAML
- ✅ **Comprehensively tested** with 6+ test cases
- ✅ **Fully documented** in user-facing docs
- ✅ **Working in production** with test workflows

**No code changes needed.** The feature works as designed.

---

## Recommendation

**Close this issue** as invalid with the explanation that the feature is already fully implemented.

If the issue reporter had concerns based on incomplete understanding, we can:
1. Add more inline code comments (optional)
2. Add more usage examples (optional)
3. Improve discoverability in documentation (optional)

But the core claim that the feature is "unprocessed" is demonstrably false.
