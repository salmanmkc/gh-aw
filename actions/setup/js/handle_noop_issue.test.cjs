// @ts-check
/// <reference types="@actions/github-script" />

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import fs from "fs";

describe("handle_noop_issue.cjs", () => {
  let main;
  let mockCore;
  let mockGithub;
  let mockContext;
  let originalEnv;

  beforeEach(async () => {
    // Save original environment
    originalEnv = { ...process.env };

    // Mock core
    mockCore = {
      info: vi.fn(),
      warning: vi.fn(),
      setFailed: vi.fn(),
      setOutput: vi.fn(),
      error: vi.fn(),
    };
    global.core = mockCore;

    // Mock github
    mockGithub = {
      rest: {
        search: {
          issuesAndPullRequests: vi.fn(),
        },
        issues: {
          create: vi.fn(),
          createComment: vi.fn(),
        },
      },
      graphql: vi.fn(),
    };
    global.github = mockGithub;

    // Mock context
    mockContext = {
      repo: {
        owner: "test-owner",
        repo: "test-repo",
      },
    };
    global.context = mockContext;

    // Set up environment
    process.env.GH_AW_WORKFLOW_NAME = "Test Workflow";
    process.env.GH_AW_AGENT_CONCLUSION = "success";
    process.env.GH_AW_RUN_URL = "https://github.com/test-owner/test-repo/actions/runs/123";
    process.env.GH_AW_WORKFLOW_SOURCE = "test-owner/test-repo/.github/workflows/test.md@main";
    process.env.GH_AW_WORKFLOW_SOURCE_URL = "https://github.com/test-owner/test-repo/blob/main/.github/workflows/test.md";
    process.env.GITHUB_REF_NAME = "copilot/test-branch";

    // Load the module
    const module = await import("./handle_noop_issue.cjs");
    main = module.main;
  });

  afterEach(() => {
    // Restore environment
    process.env = originalEnv;

    // Clean up globals
    delete global.core;
    delete global.github;
    delete global.context;

    // Clear all mocks
    vi.clearAllMocks();
  });

  describe("when agent succeeds with only noop messages", () => {
    it("should create a new no-op issue when none exists", async () => {
      // Set up agent output with only noop messages
      const tempFilePath = "/tmp/test_noop_issue_output.json";
      fs.writeFileSync(
        tempFilePath,
        JSON.stringify({
          items: [
            { type: "noop", message: "No changes needed - everything is up to date" },
            { type: "noop", message: "All tests passing, no action required" },
          ],
        })
      );
      process.env.GH_AW_AGENT_OUTPUT = tempFilePath;

      // Mock API responses
      mockGithub.rest.search.issuesAndPullRequests
        .mockResolvedValueOnce({ data: { total_count: 0, items: [] } }) // PR search
        .mockResolvedValueOnce({ data: { total_count: 1, items: [{ number: 1, html_url: "...", node_id: "I_1" }] } }) // Parent search
        .mockResolvedValueOnce({ data: { total_count: 0, items: [] } }); // No-op issue search

      mockGithub.rest.issues.create.mockResolvedValue({
        data: { number: 2, html_url: "https://example.com/2", node_id: "I_2" },
      });

      mockGithub.graphql.mockResolvedValue({
        addSubIssue: {
          issue: { id: "I_1", number: 1 },
          subIssue: { id: "I_2", number: 2 },
        },
      });

      try {
        await main();

        // Verify issue was created
        expect(mockGithub.rest.issues.create).toHaveBeenCalled();
        const createCall = mockGithub.rest.issues.create.mock.calls[0][0];

        // Check title format
        expect(createCall.title).toContain("Test Workflow - no action needed");

        // Check body contains no-op messages
        expect(createCall.body).toContain("No changes needed - everything is up to date");
        expect(createCall.body).toContain("All tests passing, no action required");

        // Check for proper formatting
        expect(createCall.body).toContain("**ℹ️ No-Op Messages**");
      } finally {
        if (fs.existsSync(tempFilePath)) {
          fs.unlinkSync(tempFilePath);
        }
      }
    });

    it("should add comment to existing no-op issue", async () => {
      // Set up agent output with noop message
      const tempFilePath = "/tmp/test_noop_comment.json";
      fs.writeFileSync(
        tempFilePath,
        JSON.stringify({
          items: [{ type: "noop", message: "Nothing to do this time" }],
        })
      );
      process.env.GH_AW_AGENT_OUTPUT = tempFilePath;

      // Mock API responses
      mockGithub.rest.search.issuesAndPullRequests
        .mockResolvedValueOnce({ data: { total_count: 0, items: [] } }) // PR search
        .mockResolvedValueOnce({ data: { total_count: 1, items: [{ number: 1, html_url: "...", node_id: "I_1" }] } }) // Parent search
        .mockResolvedValueOnce({
          data: {
            total_count: 1,
            items: [{ number: 5, html_url: "https://example.com/5" }],
          },
        }); // Existing no-op issue

      mockGithub.rest.issues.createComment.mockResolvedValue({
        data: { id: 123, html_url: "https://example.com/comment/123" },
      });

      try {
        await main();

        // Verify comment was created, not a new issue
        expect(mockGithub.rest.issues.createComment).toHaveBeenCalled();
        expect(mockGithub.rest.issues.create).not.toHaveBeenCalled();

        const commentCall = mockGithub.rest.issues.createComment.mock.calls[0][0];
        expect(commentCall.issue_number).toBe(5);
        expect(commentCall.body).toContain("Nothing to do this time");
      } finally {
        if (fs.existsSync(tempFilePath)) {
          fs.unlinkSync(tempFilePath);
        }
      }
    });
  });

  describe("when agent has non-noop outputs", () => {
    it("should skip processing", async () => {
      // Set up agent output with non-noop outputs
      const tempFilePath = "/tmp/test_skip_noop.json";
      fs.writeFileSync(
        tempFilePath,
        JSON.stringify({
          items: [
            { type: "noop", message: "Some message" },
            { type: "create_issue", title: "Test issue" },
          ],
        })
      );
      process.env.GH_AW_AGENT_OUTPUT = tempFilePath;

      try {
        await main();

        // Should skip without calling GitHub API
        expect(mockCore.info).toHaveBeenCalledWith(expect.stringContaining("non-noop outputs"));
        expect(mockGithub.rest.search.issuesAndPullRequests).not.toHaveBeenCalled();
        expect(mockGithub.rest.issues.create).not.toHaveBeenCalled();
      } finally {
        if (fs.existsSync(tempFilePath)) {
          fs.unlinkSync(tempFilePath);
        }
      }
    });
  });

  describe("when agent does not succeed", () => {
    it("should skip processing", async () => {
      process.env.GH_AW_AGENT_CONCLUSION = "failure";

      await main();

      expect(mockCore.info).toHaveBeenCalledWith(expect.stringContaining("Agent did not succeed"));
      expect(mockGithub.rest.search.issuesAndPullRequests).not.toHaveBeenCalled();
    });
  });

  describe("message sanitization", () => {
    it("should sanitize long noop messages", async () => {
      const longMessage = "A".repeat(6000);
      const tempFilePath = "/tmp/test_sanitize.json";
      fs.writeFileSync(
        tempFilePath,
        JSON.stringify({
          items: [{ type: "noop", message: longMessage }],
        })
      );
      process.env.GH_AW_AGENT_OUTPUT = tempFilePath;

      mockGithub.rest.search.issuesAndPullRequests
        .mockResolvedValueOnce({ data: { total_count: 0, items: [] } })
        .mockResolvedValueOnce({ data: { total_count: 1, items: [{ number: 1, html_url: "...", node_id: "I_1" }] } })
        .mockResolvedValueOnce({ data: { total_count: 0, items: [] } });

      mockGithub.rest.issues.create.mockResolvedValue({
        data: { number: 2, html_url: "https://example.com/2", node_id: "I_2" },
      });

      mockGithub.graphql.mockResolvedValue({
        addSubIssue: { issue: { id: "I_1", number: 1 }, subIssue: { id: "I_2", number: 2 } },
      });

      try {
        await main();

        const createCall = mockGithub.rest.issues.create.mock.calls[0][0];
        // Message should be truncated
        expect(createCall.body.length).toBeLessThan(longMessage.length + 1000);
      } finally {
        if (fs.existsSync(tempFilePath)) {
          fs.unlinkSync(tempFilePath);
        }
      }
    });
  });
});
