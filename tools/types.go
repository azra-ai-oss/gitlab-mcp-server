// Package tools defines the MCP tool input types and handler functions for the GitLab MCP server.
package tools

// --- Tool input types ---
// The `jsonschema` tag provides the description shown to AI clients.
// Required vs optional is inferred from the `json` tag: omitempty = optional.

// SearchRepositoriesInput is the input for the search_repositories tool.
type SearchRepositoriesInput struct {
	Search  string `json:"search" jsonschema:"Search query string"`
	Page    int    `json:"page,omitempty" jsonschema:"Page number (default: 1)"`
	PerPage int    `json:"per_page,omitempty" jsonschema:"Results per page (default: 20)"`
}

// GetFileContentsInput is the input for the get_file_contents tool.
type GetFileContentsInput struct {
	ProjectID string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	FilePath  string `json:"file_path" jsonschema:"Path to the file within the repository"`
	Ref       string `json:"ref,omitempty" jsonschema:"Branch name or commit SHA (default: default branch)"`
}

// CreateOrUpdateFileInput is the input for the create_or_update_file tool.
type CreateOrUpdateFileInput struct {
	ProjectID     string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	FilePath      string `json:"file_path" jsonschema:"Path where the file should be created/updated"`
	Content       string `json:"content" jsonschema:"File content"`
	CommitMessage string `json:"commit_message" jsonschema:"Commit message"`
	Branch        string `json:"branch" jsonschema:"Branch to create/update the file in"`
	PreviousPath  string `json:"previous_path,omitempty" jsonschema:"Previous file path if renaming"`
}

// FileEntry represents a single file in a push_files operation.
type FileEntry struct {
	FilePath string `json:"file_path" jsonschema:"Path for the file"`
	Content  string `json:"content" jsonschema:"File content"`
}

// PushFilesInput is the input for the push_files tool.
type PushFilesInput struct {
	ProjectID     string      `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	Branch        string      `json:"branch" jsonschema:"Branch to push to"`
	CommitMessage string      `json:"commit_message" jsonschema:"Commit message"`
	Files         []FileEntry `json:"files" jsonschema:"Array of files to push"`
}

// CreateRepositoryInput is the input for the create_repository tool.
type CreateRepositoryInput struct {
	Name                 string `json:"name" jsonschema:"Project name"`
	Description          string `json:"description,omitempty" jsonschema:"Project description"`
	Visibility           string `json:"visibility,omitempty" jsonschema:"Project visibility: private, internal, or public"`
	InitializeWithReadme bool   `json:"initialize_with_readme,omitempty" jsonschema:"Initialize with a README file"`
}

// CreateIssueInput is the input for the create_issue tool.
type CreateIssueInput struct {
	ProjectID   string   `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	Title       string   `json:"title" jsonschema:"Issue title"`
	Description string   `json:"description,omitempty" jsonschema:"Issue description"`
	AssigneeIDs []int    `json:"assignee_ids,omitempty" jsonschema:"Array of user IDs to assign"`
	MilestoneID int      `json:"milestone_id,omitempty" jsonschema:"Milestone ID"`
	Labels      []string `json:"labels,omitempty" jsonschema:"Array of label names"`
}

// CreateMergeRequestInput is the input for the create_merge_request tool.
type CreateMergeRequestInput struct {
	ProjectID          string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	Title              string `json:"title" jsonschema:"Merge request title"`
	Description        string `json:"description,omitempty" jsonschema:"Merge request description"`
	SourceBranch       string `json:"source_branch" jsonschema:"Source branch"`
	TargetBranch       string `json:"target_branch" jsonschema:"Target branch"`
	Draft              bool   `json:"draft,omitempty" jsonschema:"Create as draft merge request"`
	AllowCollaboration bool   `json:"allow_collaboration,omitempty" jsonschema:"Allow commits from upstream members"`
}

// ForkRepositoryInput is the input for the fork_repository tool.
type ForkRepositoryInput struct {
	ProjectID string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	Namespace string `json:"namespace,omitempty" jsonschema:"Namespace to fork to (default: your account)"`
}

// CreateBranchInput is the input for the create_branch tool.
type CreateBranchInput struct {
	ProjectID string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	Branch    string `json:"branch" jsonschema:"Name for the new branch"`
	Ref       string `json:"ref,omitempty" jsonschema:"Source ref to create from (default: default branch)"`
}

// --- New tool input types ---

// ListIssuesInput is the input for the list_issues tool.
type ListIssuesInput struct {
	ProjectID string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	State     string `json:"state,omitempty" jsonschema:"Filter by state: opened, closed, or all (default: all)"`
	Page      int    `json:"page,omitempty" jsonschema:"Page number (default: 1)"`
	PerPage   int    `json:"per_page,omitempty" jsonschema:"Results per page (default: 20)"`
}

// GetIssueInput is the input for the get_issue tool.
type GetIssueInput struct {
	ProjectID string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	IssueIID  int    `json:"issue_iid" jsonschema:"Issue internal ID (IID) within the project"`
}

// ListMergeRequestsInput is the input for the list_merge_requests tool.
type ListMergeRequestsInput struct {
	ProjectID string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	State     string `json:"state,omitempty" jsonschema:"Filter by state: opened, merged, closed, or all (default: all)"`
	Page      int    `json:"page,omitempty" jsonschema:"Page number (default: 1)"`
	PerPage   int    `json:"per_page,omitempty" jsonschema:"Results per page (default: 20)"`
}

// GetMergeRequestInput is the input for the get_merge_request tool.
type GetMergeRequestInput struct {
	ProjectID string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	MRIID     int    `json:"mr_iid" jsonschema:"Merge request internal ID (IID) within the project"`
}

// AddNoteInput is the input for the add_note tool.
type AddNoteInput struct {
	ProjectID   string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	NotableType string `json:"notable_type" jsonschema:"Type of object to comment on: issue or merge_request"`
	NotableIID  int    `json:"notable_iid" jsonschema:"IID of the issue or merge request"`
	Body        string `json:"body" jsonschema:"Comment body text (supports Markdown)"`
}

// ListPipelinesInput is the input for the list_pipelines tool.
type ListPipelinesInput struct {
	ProjectID string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	Ref       string `json:"ref,omitempty" jsonschema:"Filter by branch or tag name"`
	Status    string `json:"status,omitempty" jsonschema:"Filter by status: running, pending, success, failed, canceled, skipped, manual"`
	Page      int    `json:"page,omitempty" jsonschema:"Page number (default: 1)"`
	PerPage   int    `json:"per_page,omitempty" jsonschema:"Results per page (default: 20)"`
}

// GetPipelineInput is the input for the get_pipeline tool.
type GetPipelineInput struct {
	ProjectID  string `json:"project_id" jsonschema:"GitLab project ID or URL-encoded path"`
	PipelineID int    `json:"pipeline_id" jsonschema:"Pipeline ID"`
}
