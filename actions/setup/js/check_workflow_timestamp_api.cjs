// @ts-check
/// <reference types="@actions/github-script" />

/**
 * Check workflow frontmatter hashes to detect outdated lock files
 * This script compares the frontmatter hash from the lock file with
 * a recomputed hash from the source .md file to detect configuration drift
 */

const { getErrorMessage } = require("./error_helpers.cjs");
const { computeFrontmatterHash, extractHashFromLockFile } = require("./frontmatter_hash.cjs");
const fs = require("fs");
const path = require("path");

async function main() {
  const workflowFile = process.env.GH_AW_WORKFLOW_FILE;

  if (!workflowFile) {
    core.setFailed("Configuration error: GH_AW_WORKFLOW_FILE not available.");
    return;
  }

  // Construct file paths
  const workflowBasename = workflowFile.replace(".lock.yml", "");
  const workflowMdPath = `.github/workflows/${workflowBasename}.md`;
  const lockFilePath = `.github/workflows/${workflowFile}`;

  core.info(`Checking workflow frontmatter hashes:`);
  core.info(`  Source: ${workflowMdPath}`);
  core.info(`  Lock file: ${lockFilePath}`);

  const { owner, repo } = context.repo;

  // Check frontmatter hashes and fail if they don't match
  await checkFrontmatterHashes(workflowMdPath, lockFilePath, owner, repo);
}

/**
 * Check frontmatter hashes from lock file and recomputed from source
 * Fails the step if hashes don't match
 * @param {string} workflowMdPath - Path to the source .md file
 * @param {string} lockFilePath - Path to the compiled .lock.yml file
 * @param {string} owner - Repository owner
 * @param {string} repo - Repository name
 */
async function checkFrontmatterHashes(workflowMdPath, lockFilePath, owner, repo) {
  // Extract hash from lock file
  let lockFileHash = "";
  try {
    const lockContent = fs.readFileSync(lockFilePath, "utf8");
    lockFileHash = extractHashFromLockFile(lockContent);
  } catch (error) {
    const errorMessage = getErrorMessage(error);
    core.setFailed(`Could not read lock file: ${errorMessage}`);
    return;
  }

  if (!lockFileHash) {
    core.setFailed(`Lock file '${lockFilePath}' does not contain a frontmatter hash comment.`);
    return;
  }

  // Recompute hash from .md file
  let recomputedHash = "";
  try {
    // Get the absolute path to the workflow file
    const workspacePath = process.env.GITHUB_WORKSPACE || process.cwd();
    const absoluteMdPath = path.join(workspacePath, workflowMdPath);

    if (!fs.existsSync(absoluteMdPath)) {
      core.setFailed(`Source file not found: ${workflowMdPath}`);
      return;
    }

    recomputedHash = await computeFrontmatterHash(absoluteMdPath);
  } catch (error) {
    const errorMessage = getErrorMessage(error);
    core.setFailed(`Could not compute hash from source: ${errorMessage}`);
    return;
  }

  // Display both hashes
  core.info(`  Lock file hash:  ${lockFileHash}`);
  core.info(`  Recomputed hash: ${recomputedHash}`);

  // Check if hashes match
  if (lockFileHash === recomputedHash) {
    core.info("✅ Frontmatter hashes match - lock file is up to date");
  } else {
    const errorMessage = `Lock file '${lockFilePath}' is outdated! The frontmatter hash does not match the source file '${workflowMdPath}'. Run 'gh aw compile' to regenerate the lock file.`;

    // Add summary to GitHub Step Summary
    let summary = core.summary
      .addRaw("### ⚠️ Workflow Lock File Outdated\n\n")
      .addRaw("**ERROR**: Lock file frontmatter hash does not match the source file.\n\n")
      .addRaw("**Files:**\n")
      .addRaw(`- Source: \`${workflowMdPath}\`\n`)
      .addRaw(`- Lock: \`${lockFilePath}\`\n\n`)
      .addRaw("**Hashes:**\n")
      .addRaw(`- Lock file hash: \`${lockFileHash}\`\n`)
      .addRaw(`- Recomputed hash: \`${recomputedHash}\`\n\n`)
      .addRaw("**Action Required:** Run `gh aw compile` to regenerate the lock file.\n\n");

    await summary.write();

    // Fail the step to prevent workflow from running with outdated configuration
    core.setFailed(errorMessage);
  }
}

module.exports = { main };
