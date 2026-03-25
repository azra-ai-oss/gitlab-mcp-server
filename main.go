package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/azra-ai-oss/gitlab-mcp-server/gitlab"
	"github.com/azra-ai-oss/gitlab-mcp-server/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.2.0"

func main() {
	// Structured logging to stderr (MCP uses stdout for JSON-RPC)
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	token := os.Getenv("GITLAB_PERSONAL_ACCESS_TOKEN")
	if token == "" {
		logger.Error("GITLAB_PERSONAL_ACCESS_TOKEN environment variable is not set")
		os.Exit(1)
	}

	apiURL := os.Getenv("GITLAB_API_URL")
	client := gitlab.NewClient(apiURL, token, logger)
	h := tools.NewHandlers(client)

	server := mcp.NewServer(
		&mcp.Implementation{Name: "gitlab-mcp-server", Version: version},
		nil,
	)

	// --- Read-only tools only (defense-in-depth: no write operations) ---

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_repositories",
		Description: "Search for GitLab projects",
	}, h.SearchRepositories)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_file_contents",
		Description: "Get the contents of a file or directory from a GitLab project",
	}, h.GetFileContents)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_issues",
		Description: "List issues in a GitLab project with optional state filtering",
	}, h.ListIssues)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_issue",
		Description: "Get details of a specific issue by its IID",
	}, h.GetIssue)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_merge_requests",
		Description: "List merge requests in a GitLab project with optional state filtering",
	}, h.ListMergeRequests)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_merge_request",
		Description: "Get details of a specific merge request by its IID",
	}, h.GetMergeRequest)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_pipelines",
		Description: "List CI/CD pipelines in a GitLab project with optional filtering",
	}, h.ListPipelines)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_pipeline",
		Description: "Get details of a specific CI/CD pipeline",
	}, h.GetPipeline)

	// Run the server over stdio
	logger.Info("GitLab MCP Server starting", "version", version, "tools", 8)
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		logger.Error("Fatal error", "error", err)
		os.Exit(1)
	}
}
