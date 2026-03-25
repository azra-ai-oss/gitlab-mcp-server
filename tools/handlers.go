package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/azra-ai-oss/gitlab-mcp-server/gitlab"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Handlers holds the GitLab client and provides tool handler methods.
type Handlers struct {
	Client *gitlab.Client
}

// NewHandlers creates a new Handlers instance with the given GitLab client.
func NewHandlers(client *gitlab.Client) *Handlers {
	return &Handlers{Client: client}
}

func toJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// textResult is a helper to create a text-only CallToolResult.
func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

// SearchRepositories handles the search_repositories tool.
func (h *Handlers) SearchRepositories(ctx context.Context, req *mcp.CallToolRequest, input SearchRepositoriesInput) (*mcp.CallToolResult, any, error) {
	result, err := h.Client.SearchProjects(input.Search, input.Page, input.PerPage)
	if err != nil {
		return nil, nil, fmt.Errorf("searching repositories: %w", err)
	}
	return textResult(toJSON(result)), nil, nil
}

// GetFileContents handles the get_file_contents tool.
func (h *Handlers) GetFileContents(ctx context.Context, req *mcp.CallToolRequest, input GetFileContentsInput) (*mcp.CallToolResult, any, error) {
	content, err := h.Client.GetFileContents(input.ProjectID, input.FilePath, input.Ref)
	if err != nil {
		return nil, nil, fmt.Errorf("getting file contents: %w", err)
	}
	return textResult(toJSON(content)), nil, nil
}

// CreateOrUpdateFile handles the create_or_update_file tool.
func (h *Handlers) CreateOrUpdateFile(ctx context.Context, req *mcp.CallToolRequest, input CreateOrUpdateFileInput) (*mcp.CallToolResult, any, error) {
	result, err := h.Client.CreateOrUpdateFile(
		input.ProjectID, input.FilePath, input.Content, input.CommitMessage, input.Branch, input.PreviousPath,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating/updating file: %w", err)
	}
	return textResult(toJSON(result)), nil, nil
}

// PushFiles handles the push_files tool.
func (h *Handlers) PushFiles(ctx context.Context, req *mcp.CallToolRequest, input PushFilesInput) (*mcp.CallToolResult, any, error) {
	actions := make([]gitlab.CommitAction, len(input.Files))
	for i, f := range input.Files {
		actions[i] = gitlab.CommitAction{
			Action:   "create",
			FilePath: f.FilePath,
			Content:  f.Content,
		}
	}

	commit, err := h.Client.CreateCommit(input.ProjectID, input.CommitMessage, input.Branch, actions)
	if err != nil {
		return nil, nil, fmt.Errorf("pushing files: %w", err)
	}
	return textResult(toJSON(commit)), nil, nil
}

// CreateRepository handles the create_repository tool.
func (h *Handlers) CreateRepository(ctx context.Context, req *mcp.CallToolRequest, input CreateRepositoryInput) (*mcp.CallToolResult, any, error) {
	project, err := h.Client.CreateProject(input.Name, input.Description, input.Visibility, input.InitializeWithReadme)
	if err != nil {
		return nil, nil, fmt.Errorf("creating repository: %w", err)
	}
	return textResult(toJSON(project)), nil, nil
}

// CreateIssue handles the create_issue tool.
func (h *Handlers) CreateIssue(ctx context.Context, req *mcp.CallToolRequest, input CreateIssueInput) (*mcp.CallToolResult, any, error) {
	issue, err := h.Client.CreateIssue(
		input.ProjectID, input.Title, input.Description, input.AssigneeIDs, input.MilestoneID, input.Labels,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating issue: %w", err)
	}
	return textResult(toJSON(issue)), nil, nil
}

// CreateMergeRequest handles the create_merge_request tool.
func (h *Handlers) CreateMergeRequest(ctx context.Context, req *mcp.CallToolRequest, input CreateMergeRequestInput) (*mcp.CallToolResult, any, error) {
	mr, err := h.Client.CreateMergeRequest(
		input.ProjectID, input.Title, input.Description, input.SourceBranch, input.TargetBranch, input.Draft, input.AllowCollaboration,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating merge request: %w", err)
	}
	return textResult(toJSON(mr)), nil, nil
}

// ForkRepository handles the fork_repository tool.
func (h *Handlers) ForkRepository(ctx context.Context, req *mcp.CallToolRequest, input ForkRepositoryInput) (*mcp.CallToolResult, any, error) {
	fork, err := h.Client.ForkProject(input.ProjectID, input.Namespace)
	if err != nil {
		return nil, nil, fmt.Errorf("forking repository: %w", err)
	}
	return textResult(toJSON(fork)), nil, nil
}

// CreateBranch handles the create_branch tool.
func (h *Handlers) CreateBranch(ctx context.Context, req *mcp.CallToolRequest, input CreateBranchInput) (*mcp.CallToolResult, any, error) {
	ref := input.Ref
	if ref == "" {
		defaultBranch, err := h.Client.GetDefaultBranch(input.ProjectID)
		if err != nil {
			return nil, nil, fmt.Errorf("getting default branch: %w", err)
		}
		ref = defaultBranch
	}

	branch, err := h.Client.CreateBranch(input.ProjectID, input.Branch, ref)
	if err != nil {
		return nil, nil, fmt.Errorf("creating branch: %w", err)
	}
	return textResult(toJSON(branch)), nil, nil
}

// --- New handlers ---

// ListIssues handles the list_issues tool.
func (h *Handlers) ListIssues(ctx context.Context, req *mcp.CallToolRequest, input ListIssuesInput) (*mcp.CallToolResult, any, error) {
	issues, total, err := h.Client.ListIssues(input.ProjectID, input.State, input.Page, input.PerPage)
	if err != nil {
		return nil, nil, fmt.Errorf("listing issues: %w", err)
	}
	return textResult(toJSON(map[string]any{"count": total, "items": issues})), nil, nil
}

// GetIssue handles the get_issue tool.
func (h *Handlers) GetIssue(ctx context.Context, req *mcp.CallToolRequest, input GetIssueInput) (*mcp.CallToolResult, any, error) {
	issue, err := h.Client.GetIssue(input.ProjectID, input.IssueIID)
	if err != nil {
		return nil, nil, fmt.Errorf("getting issue: %w", err)
	}
	return textResult(toJSON(issue)), nil, nil
}

// ListMergeRequests handles the list_merge_requests tool.
func (h *Handlers) ListMergeRequests(ctx context.Context, req *mcp.CallToolRequest, input ListMergeRequestsInput) (*mcp.CallToolResult, any, error) {
	mrs, total, err := h.Client.ListMergeRequests(input.ProjectID, input.State, input.Page, input.PerPage)
	if err != nil {
		return nil, nil, fmt.Errorf("listing merge requests: %w", err)
	}
	return textResult(toJSON(map[string]any{"count": total, "items": mrs})), nil, nil
}

// GetMergeRequest handles the get_merge_request tool.
func (h *Handlers) GetMergeRequest(ctx context.Context, req *mcp.CallToolRequest, input GetMergeRequestInput) (*mcp.CallToolResult, any, error) {
	mr, err := h.Client.GetMergeRequest(input.ProjectID, input.MRIID)
	if err != nil {
		return nil, nil, fmt.Errorf("getting merge request: %w", err)
	}
	return textResult(toJSON(mr)), nil, nil
}

// AddNote handles the add_note tool.
func (h *Handlers) AddNote(ctx context.Context, req *mcp.CallToolRequest, input AddNoteInput) (*mcp.CallToolResult, any, error) {
	// Map user-friendly type to GitLab API path segment
	noteableType := input.NoteableType
	switch noteableType {
	case "issue":
		noteableType = "issues"
	case "merge_request":
		noteableType = "merge_requests"
	}

	note, err := h.Client.AddNote(input.ProjectID, noteableType, input.NoteableIID, input.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("adding note: %w", err)
	}
	return textResult(toJSON(note)), nil, nil
}

// ListPipelines handles the list_pipelines tool.
func (h *Handlers) ListPipelines(ctx context.Context, req *mcp.CallToolRequest, input ListPipelinesInput) (*mcp.CallToolResult, any, error) {
	pipelines, total, err := h.Client.ListPipelines(input.ProjectID, input.Ref, input.Status, input.Page, input.PerPage)
	if err != nil {
		return nil, nil, fmt.Errorf("listing pipelines: %w", err)
	}
	return textResult(toJSON(map[string]any{"count": total, "items": pipelines})), nil, nil
}

// GetPipeline handles the get_pipeline tool.
func (h *Handlers) GetPipeline(ctx context.Context, req *mcp.CallToolRequest, input GetPipelineInput) (*mcp.CallToolResult, any, error) {
	pipeline, err := h.Client.GetPipeline(input.ProjectID, input.PipelineID)
	if err != nil {
		return nil, nil, fmt.Errorf("getting pipeline: %w", err)
	}
	return textResult(toJSON(pipeline)), nil, nil
}
