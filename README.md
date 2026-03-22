# GitLab MCP Server

A Go implementation of the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for GitLab, providing AI assistants with seamless access to GitLab repositories, issues, merge requests, and more.

## Why This Exists

The [official GitLab MCP server](https://github.com/modelcontextprotocol/servers/tree/main/src/gitlab) is written in TypeScript and requires a Node.js runtime. This creates friction for teams that:

- **Don't want to install Node.js** — This server compiles to a **single static binary** (~10-15 MB) with zero runtime dependencies. Download it, set your token, and run it.
- **Need better resource efficiency** — The Node.js GitLab MCP server uses ~30-50 MB of RAM and ~200ms startup time. This Go version uses ~5-10 MB of RAM and starts in <5ms — meaningful when running multiple MCP servers simultaneously.
- **Want reliable schema validation** — The community [custom-gitlab-mcp-server](https://github.com/chris-miaskowski/custom-gitlab-mcp-server) was created specifically to fix schema validation bugs in the official TypeScript server. This Go implementation gets correct schemas for free via Go struct tags and the official Go MCP SDK.
- **Prefer a single binary for distribution** — `go install`, copy the binary, or use the Nix flake. No `npm install`, no `node_modules`, no version conflicts.

## Quick Start

### Prerequisites

- Go 1.25+ (or use the included Nix flake: `nix develop`)
- A GitLab Personal Access Token with appropriate scopes (`api`, `read_api`, `read_repository`, `write_repository`)

### Build

```bash
go build -o gitlab-mcp-server .
```

### Configure Your AI Client

Add to your MCP client configuration (e.g., Claude Desktop, Gemini CLI):

```json
{
  "mcpServers": {
    "gitlab": {
      "command": "/path/to/gitlab-mcp-server",
      "env": {
        "GITLAB_PERSONAL_ACCESS_TOKEN": "your-token-here",
        "GITLAB_API_URL": "https://gitlab.com/api/v4"
      }
    }
  }
}
```

> `GITLAB_API_URL` defaults to `https://gitlab.com/api/v4` if not set. Point it to your self-hosted instance if needed.

## Available Tools

| Tool | Description |
|------|-------------|
| `search_repositories` | Search for GitLab projects |
| `get_file_contents` | Get the contents of a file from a project |
| `create_or_update_file` | Create or update a single file in a project |
| `push_files` | Push multiple files in a single commit |
| `create_repository` | Create a new GitLab project |
| `create_issue` | Create a new issue in a project |
| `create_merge_request` | Create a new merge request in a project |
| `fork_repository` | Fork a project to your account or namespace |
| `create_branch` | Create a new branch in a project |

## Development

### Using Nix (recommended)

```bash
# Enter dev shell with Go 1.25, gopls, etc.
nix develop

# Or if you have direnv:
direnv allow
```

### Build & Run

```bash
export GITLAB_PERSONAL_ACCESS_TOKEN="your-token"
go build -o gitlab-mcp-server .
./gitlab-mcp-server
```

## License

MIT
