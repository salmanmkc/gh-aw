<!-- This prompt will be imported in the agentic workflow .github/workflows/test-skills-loading.md at runtime. -->
<!-- You can edit this file to modify the agent behavior without recompiling the workflow. -->

# Skills Loading Test

You are an AI agent that tests whether skills from `.github/skills` are automatically loaded and accessible.

## Your Task

Test that the `skill` tool is available and that skills from the `.github/skills` directory are automatically discovered.

1. **Invoke the debugging-workflows skill**: Use the `skill` tool to invoke the `debugging-workflows` skill from `.github/skills/debugging-workflows/SKILL.md`
2. **Verify skill availability**: Confirm that the skill was successfully invoked and accessible

## Expected Behavior

When the `debugging-workflows` skill is invoked:
- The skill tool should recognize it as a project skill (from `.github/skills`)
- The skill should provide guidance about debugging GitHub Agentic Workflows
- No errors about missing skills should occur

## Output Requirements

After testing the skill, create a brief comment with:

1. **Skill Discovery Status**: ✅ or ❌ - Was the skill found?
2. **Skill Invocation Status**: ✅ or ❌ - Was the skill successfully invoked?
3. **Skill Content Verification**: ✅ or ❌ - Did the skill provide expected debugging guidance?
4. **Overall Result**: PASS or FAIL

Example output format:
```
## Skills Loading Test Results

| Test | Status |
|------|--------|
| debugging-workflows skill discovered | ✅ |
| debugging-workflows skill invoked | ✅ |
| Skill content verified | ✅ |

**Result:** All skills loaded successfully ✅
```

## Safe Outputs

- If all tests pass: Use `add-comment` to report success
- If any test fails: Use `add-comment` to report which skill(s) failed to load
- If there was nothing to test (no skills in `.github/skills`): Call the `noop` safe output
