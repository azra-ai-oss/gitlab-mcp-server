package main

import (
	"context"
	"fmt"
	"os"

	"github.com/azra-ai-oss/gitlab-mcp-server/gitlab"
	"github.com/azra-ai-oss/gitlab-mcp-server/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const version = "0.1.0"

func main() {
	token := os.Getenv("GITLAB_PERSONAL_ACCESS_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: GITLAB_PERSONAL_ACCESS_TOKEN environment variable is not set")
		os.Exit(1)
	}

	apiURL := os.Getenv("GITLAB_API_URL")
	client := gitlab.NewClient(apiURL, token)
	h := tools.NewHandlers(client)

	server := mcp.NewServer(
		&mcp.Implementation{Name: "gitlab-mcp-server", Version: version},
		nil,
	)

	// Register all 9 tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_repositories",
		Description: "Search for GitLab projects",
	}, h.SearchRepositories)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_file_contents",
		Description: "Get the contents of a file or directory from a GitLab project",
	}, h.GetFileContents)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_or_update_file",
		Description: "Create or update a single file in a GitLab project",
	}, h.CreateOrUpdateFile)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "push_files",
		Description: "Push multiple files to a GitLab project in a single commit",
	}, h.PushFiles)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_repository",
		Description: "Create a new GitLab project",
	}, h.CreateRepository)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_issue",
		Description: "Create a new issue in a GitLab project",
	}, h.CreateIssue)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_merge_request",
		Description: "Create a new merge request in a GitLab project",
	}, h.CreateMergeRequest)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "fork_repository",
		Description: "Fork a GitLab project to your account or specified namespace",
	}, h.ForkRepository)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_branch",
		Description: "Create a new branch in a GitLab project",
	}, h.CreateBranch)

	// Run the server over stdio
	fmt.Fprintln(os.Stderr, "GitLab MCP Server running on stdio")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}
