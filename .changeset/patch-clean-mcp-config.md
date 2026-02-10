---
"gh-aw": patch
---

Gateway startup now removes `/home/runner/.copilot/mcp-config.json` after configuring the MCP server so the temporary config isn't left on disk.
