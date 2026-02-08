// @ts-check
/// <reference types="@actions/github-script" />

const { getErrorMessage } = require("./error_helpers.cjs");
const { sanitizeContent } = require("./sanitize_content.cjs");
const { getFooterAgentFailureIssueMessage, getFooterAgentFailureCommentMessage, generateXMLMarker } = require("./messages.cjs");
const { renderTemplate } = require("./messages_core.cjs");
const { getCurrentBranch } = require("./get_current_branch.cjs");
const { generateFooterWithExpiration } = require("./ephemerals.cjs");
const { MAX_SUB_ISSUES, getSubIssueCount } = require("./sub_issue_helpers.cjs");
const fs = require("fs");

/**
 * Attempt to find a pull request for the current branch
 * @returns {Promise<{number: number, html_url: string} | null>} PR info or null if not found
 */
async function findPullRequestForCurrentBranch() {
  try {
    const { owner, repo } = context.repo;
    const currentBranch = getCurrentBranch();

    core.info(`Searching for pull request from branch: ${currentBranch}`);

    // Search for open PRs with the current branch as head
    const searchQuery = `repo:${owner}/${repo} is:pr is:open head:${currentBranch}`;

    const searchResult = await github.rest.search.issuesAndPullRequests({
      q: searchQuery,
      per_page: 1,
    });

    if (searchResult.data.total_count > 0) {
      const pr = searchResult.data.items[0];
      core.info(`Found pull request #${pr.number}: ${pr.html_url}`);
      return {
        number: pr.number,
        html_url: pr.html_url,
      };
    }

    core.info(`No pull request found for branch: ${currentBranch}`);
    return null;
  } catch (error) {
    core.warning(`Failed to find pull request for current branch: ${getErrorMessage(error)}`);
    return null;
  }
}

/**
 * Search for or create the parent issue for all agentic workflow no-op messages
 * @param {number|null} previousParentNumber - Previous parent issue number if creating due to limit
 * @returns {Promise<{number: number, node_id: string}>} Parent issue number and node ID
 */
async function ensureParentIssue(previousParentNumber = null) {
  const { owner, repo } = context.repo;
  const parentTitle = "[agentic-workflows] No-op runs";
  const parentLabel = "agentic-workflows";

  core.info(`Searching for parent issue: "${parentTitle}"`);

  // Search for existing parent issue
  const searchQuery = `repo:${owner}/${repo} is:issue is:open label:${parentLabel} in:title "${parentTitle}"`;

  try {
    const searchResult = await github.rest.search.issuesAndPullRequests({
      q: searchQuery,
      per_page: 1,
    });

    if (searchResult.data.total_count > 0) {
      const existingIssue = searchResult.data.items[0];
      core.info(`Found existing parent issue #${existingIssue.number}: ${existingIssue.html_url}`);

      // Check the sub-issue count
      const subIssueCount = await getSubIssueCount(owner, repo, existingIssue.number);

      if (subIssueCount !== null && subIssueCount >= MAX_SUB_ISSUES) {
        core.warning(`Parent issue #${existingIssue.number} has ${subIssueCount} sub-issues (max: ${MAX_SUB_ISSUES})`);
        core.info(`Creating a new parent issue (previous parent #${existingIssue.number} is full)`);

        // Fall through to create a new parent issue, passing the previous parent number
        previousParentNumber = existingIssue.number;
      } else {
        // Parent issue is within limits, return it
        if (subIssueCount !== null) {
          core.info(`Parent issue has ${subIssueCount} sub-issues (within limit of ${MAX_SUB_ISSUES})`);
        }
        return {
          number: existingIssue.number,
          node_id: existingIssue.node_id,
        };
      }
    }
  } catch (error) {
    core.warning(`Error searching for parent issue: ${getErrorMessage(error)}`);
  }

  // Create parent issue if it doesn't exist or if previous one is full
  const creationReason = previousParentNumber ? `creating new parent (previous #${previousParentNumber} reached limit)` : "creating first parent";
  core.info(`No suitable parent issue found, ${creationReason}`);

  let parentBodyContent = `This issue tracks all no-op runs from agentic workflows in this repository. Each workflow run where the agent determined no action was needed creates a sub-issue linked here for organization and easy filtering.`;

  // Add reference to previous parent if this is a continuation
  if (previousParentNumber) {
    parentBodyContent += `

> **Note:** This is a continuation parent issue. The previous parent issue #${previousParentNumber} reached the maximum of ${MAX_SUB_ISSUES} sub-issues.`;
  }

  parentBodyContent += `

### Purpose

This parent issue helps you:
- View all no-op runs in one place by checking the sub-issues below
- Filter out no-op issues from your main issue list using \`no:parent-issue\`
- Understand when the agent runs but determines no action is needed

### Sub-Issues

All individual no-op run issues are linked as sub-issues below. Click on any sub-issue to see details about a specific run.

---

> This issue is automatically managed by GitHub Agentic Workflows. Do not close this issue manually.`;

  // Add expiration marker (7 days from now) inside the quoted section using helper
  const footer = generateFooterWithExpiration({
    footerText: parentBodyContent,
    expiresHours: 24 * 7, // 7 days
  });
  const parentBody = footer;

  try {
    const newIssue = await github.rest.issues.create({
      owner,
      repo,
      title: parentTitle,
      body: parentBody,
      labels: [parentLabel],
    });

    core.info(`✓ Created parent issue #${newIssue.data.number}: ${newIssue.data.html_url}`);
    return {
      number: newIssue.data.number,
      node_id: newIssue.data.node_id,
    };
  } catch (error) {
    core.error(`Failed to create parent issue: ${getErrorMessage(error)}`);
    throw error;
  }
}

/**
 * Link an issue as a sub-issue to a parent issue
 * @param {string} parentNodeId - GraphQL node ID of the parent issue
 * @param {string} subIssueNodeId - GraphQL node ID of the sub-issue
 * @param {number} parentNumber - Parent issue number (for logging)
 * @param {number} subIssueNumber - Sub-issue number (for logging)
 */
async function linkSubIssue(parentNodeId, subIssueNodeId, parentNumber, subIssueNumber) {
  core.info(`Linking issue #${subIssueNumber} as sub-issue of #${parentNumber}`);

  try {
    // Use GraphQL to link the sub-issue
    await github.graphql(
      `mutation($parentId: ID!, $subIssueId: ID!) {
        addSubIssue(input: {issueId: $parentId, subIssueId: $subIssueId}) {
          issue {
            id
            number
          }
          subIssue {
            id
            number
          }
        }
      }`,
      {
        parentId: parentNodeId,
        subIssueId: subIssueNodeId,
      }
    );

    core.info(`✓ Successfully linked #${subIssueNumber} as sub-issue of #${parentNumber}`);
  } catch (error) {
    const errorMessage = getErrorMessage(error);
    if (errorMessage.includes("Field 'addSubIssue' doesn't exist") || errorMessage.includes("not yet available")) {
      core.warning(`Sub-issue API not available. Issue #${subIssueNumber} created but not linked to parent.`);
    } else {
      core.warning(`Failed to link sub-issue: ${errorMessage}`);
    }
  }
}

/**
 * Handle workflow runs where agent succeeded with only no-op messages
 * This script creates or updates an issue to track no-op runs separately from failures
 */
async function main() {
  try {
    // Get workflow context
    const workflowName = process.env.GH_AW_WORKFLOW_NAME || "unknown";
    const agentConclusion = process.env.GH_AW_AGENT_CONCLUSION || "";
    const runUrl = process.env.GH_AW_RUN_URL || "";
    const workflowSource = process.env.GH_AW_WORKFLOW_SOURCE || "";
    const workflowSourceURL = process.env.GH_AW_WORKFLOW_SOURCE_URL || "";

    core.info(`Agent conclusion: ${agentConclusion}`);
    core.info(`Workflow name: ${workflowName}`);

    // Only proceed if agent succeeded
    if (agentConclusion !== "success") {
      core.info(`Agent did not succeed (conclusion: ${agentConclusion}), skipping no-op issue handling`);
      return;
    }

    // Load agent output to check for noop messages
    const { loadAgentOutput } = require("./load_agent_output.cjs");
    const agentOutputResult = loadAgentOutput();

    if (!agentOutputResult.success || !agentOutputResult.items || agentOutputResult.items.length === 0) {
      core.info("No agent output found, skipping no-op issue handling");
      return;
    }

    // Check if all outputs are noop messages (no other safe-output types)
    const noopItems = agentOutputResult.items.filter(item => item.type === "noop");
    const nonNoopItems = agentOutputResult.items.filter(item => item.type !== "noop");

    if (noopItems.length === 0 || nonNoopItems.length > 0) {
      core.info(`Agent produced non-noop outputs (${nonNoopItems.length} non-noop, ${noopItems.length} noop), skipping no-op issue handling`);
      return;
    }

    // Extract noop messages
    const noopMessages = noopItems.map(item => item.message || "");
    core.info(`Agent succeeded with only noop messages (${noopMessages.length} message(s))`);

    const { owner, repo } = context.repo;

    // Try to find a pull request for the current branch
    const pullRequest = await findPullRequestForCurrentBranch();

    // Ensure parent issue exists first
    let parentIssue;
    try {
      parentIssue = await ensureParentIssue();
    } catch (error) {
      core.warning(`Could not create parent issue, proceeding without parent: ${getErrorMessage(error)}`);
      // Continue without parent issue
    }

    // Sanitize workflow name for title
    const sanitizedWorkflowName = sanitizeContent(workflowName, { maxLength: 100 });
    const issueTitle = `[agentic-workflows] ${sanitizedWorkflowName} - no action needed`;

    core.info(`Checking for existing issue with title: "${issueTitle}"`);

    // Search for existing open issue with this title and label
    const searchQuery = `repo:${owner}/${repo} is:issue is:open label:agentic-workflows in:title "${issueTitle}"`;

    try {
      const searchResult = await github.rest.search.issuesAndPullRequests({
        q: searchQuery,
        per_page: 1,
      });

      if (searchResult.data.total_count > 0) {
        // Issue exists, add a comment
        const existingIssue = searchResult.data.items[0];
        core.info(`Found existing issue #${existingIssue.number}: ${existingIssue.html_url}`);

        // Build noop messages context
        let noopMessagesText = "";
        noopMessages.forEach((msg, idx) => {
          const sanitizedMsg = sanitizeContent(msg, { maxLength: 5000 });
          if (noopMessages.length === 1) {
            noopMessagesText += `${sanitizedMsg}\n\n`;
          } else {
            noopMessagesText += `${idx + 1}. ${sanitizedMsg}\n\n`;
          }
        });

        const commentBody = `### No-Op Run

**Run URL:** ${runUrl}

**Messages:**

${noopMessagesText}

---

> Generated from [${workflowName}](${runUrl})`;

        const sanitizedCommentBody = sanitizeContent(commentBody, { maxLength: 65000 });

        await github.rest.issues.createComment({
          owner,
          repo,
          issue_number: existingIssue.number,
          body: sanitizedCommentBody,
        });

        core.info(`✓ Added comment to existing issue #${existingIssue.number}`);
      } else {
        // No existing issue, create a new one
        core.info("No existing issue found, creating a new one");

        // Get current branch information
        const currentBranch = getCurrentBranch();

        // Build noop messages context
        let noopMessagesText = "";
        noopMessages.forEach((msg, idx) => {
          const sanitizedMsg = sanitizeContent(msg, { maxLength: 5000 });
          if (noopMessages.length === 1) {
            noopMessagesText += `${sanitizedMsg}\n\n`;
          } else {
            noopMessagesText += `${idx + 1}. ${sanitizedMsg}\n\n`;
          }
        });

        // Build issue body
        let issueBodyContent = `### No Action Needed

**Workflow:** [${sanitizedWorkflowName}](${workflowSourceURL})  
**Branch:** ${currentBranch}  
**Run URL:** ${runUrl}`;

        if (pullRequest) {
          issueBodyContent += `  
**Pull Request:** [#${pullRequest.number}](${pullRequest.html_url})`;
        }

        issueBodyContent += `

**ℹ️ No-Op Messages**: The agent ran but determined no action was needed. The following message(s) were reported:

${noopMessagesText}

---

This issue is created for tracking purposes when an agentic workflow runs successfully but determines no action is needed.`;

        // Generate footer
        const ctx = {
          workflowName,
          runUrl,
          workflowSource,
          workflowSourceUrl: workflowSourceURL,
        };
        const footer = getFooterAgentFailureIssueMessage(ctx);

        // Add expiration marker (7 days from now) inside the quoted footer section using helper
        const footerWithExpires = generateFooterWithExpiration({
          footerText: footer,
          expiresHours: 24 * 7, // 7 days
          suffix: `\n\n${generateXMLMarker(workflowName, runUrl)}`,
        });

        // Combine issue body with footer
        const bodyLines = [issueBodyContent, "", footerWithExpires];
        const issueBody = bodyLines.join("\n");

        const newIssue = await github.rest.issues.create({
          owner,
          repo,
          title: issueTitle,
          body: issueBody,
          labels: ["agentic-workflows"],
        });

        core.info(`✓ Created new issue #${newIssue.data.number}: ${newIssue.data.html_url}`);

        // Link as sub-issue to parent if parent issue was created
        if (parentIssue) {
          try {
            await linkSubIssue(parentIssue.node_id, newIssue.data.node_id, parentIssue.number, newIssue.data.number);
          } catch (error) {
            core.warning(`Could not link issue as sub-issue: ${getErrorMessage(error)}`);
            // Continue even if linking fails
          }
        }
      }
    } catch (error) {
      core.warning(`Failed to create or update no-op tracking issue: ${getErrorMessage(error)}`);
      // Don't fail the workflow if we can't create the issue
    }
  } catch (error) {
    core.warning(`Error in handle_noop_issue: ${getErrorMessage(error)}`);
    // Don't fail the workflow
  }
}

module.exports = { main };
