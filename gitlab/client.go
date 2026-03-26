// Package gitlab provides an HTTP client for the GitLab REST API v4.
package gitlab

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const defaultTimeout = 30 * time.Second

// Client is a thin HTTP wrapper for GitLab API v4.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	logger     *slog.Logger
}

// NewClient creates a new GitLab API client with a 30s timeout.
func NewClient(baseURL, token string, logger *slog.Logger) *Client {
	if baseURL == "" {
		baseURL = "https://gitlab.com/api/v4"
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		logger: logger,
	}
}

// --- Response types ---

type Project struct {
	ID            any     `json:"id"`
	Name          string  `json:"name"`
	Description   *string `json:"description"`
	DefaultBranch string  `json:"default_branch"`
	Visibility    string  `json:"visibility,omitempty"`
	WebURL        string  `json:"web_url,omitempty"`
	HTTPURLToRepo string  `json:"http_url_to_repo,omitempty"`
}

type Branch struct {
	Name   string       `json:"name"`
	Commit BranchCommit `json:"commit"`
}

type BranchCommit struct {
	ID      string `json:"id"`
	ShortID string `json:"short_id"`
	Title   string `json:"title"`
}

type FileContent struct {
	FileName      string `json:"file_name"`
	FilePath      string `json:"file_path"`
	Size          int    `json:"size"`
	Encoding      string `json:"encoding"`
	Content       string `json:"content"`
	ContentSHA256 string `json:"content_sha256,omitempty"`
	Ref           string `json:"ref,omitempty"`
	BlobID        string `json:"blob_id,omitempty"`
	CommitID      string `json:"commit_id,omitempty"`
	LastCommitID  string `json:"last_commit_id,omitempty"`
}

type FileResponse struct {
	FilePath string `json:"file_path"`
	Branch   string `json:"branch"`
}

type Commit struct {
	ID      string `json:"id"`
	ShortID string `json:"short_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Issue struct {
	ID          int      `json:"id"`
	IID         int      `json:"iid"`
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	State       string   `json:"state,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	WebURL      string   `json:"web_url,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	UpdatedAt   string   `json:"updated_at,omitempty"`
}

type MergeRequest struct {
	ID           int     `json:"id"`
	IID          int     `json:"iid"`
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	State        string  `json:"state,omitempty"`
	SourceBranch string  `json:"source_branch"`
	TargetBranch string  `json:"target_branch"`
	WebURL       string  `json:"web_url,omitempty"`
	MergeStatus  string  `json:"merge_status,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
}

type Note struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Author    Author `json:"author"`
	CreatedAt string `json:"created_at,omitempty"`
}

type Author struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type Pipeline struct {
	ID        int    `json:"id"`
	IID       int    `json:"iid"`
	Status    string `json:"status"`
	Ref       string `json:"ref"`
	SHA       string `json:"sha"`
	WebURL    string `json:"web_url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Source    string `json:"source,omitempty"`
}

type SearchResult struct {
	Count int       `json:"count"`
	Items []Project `json:"items"`
}

// CommitAction represents a single file action in a commit.
type CommitAction struct {
	Action   string `json:"action"`
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

// --- Internal HTTP helpers ---

func (c *Client) do(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	c.logger.Debug("GitLab API request", "method", method, "path", path)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("GitLab API request failed", "method", method, "path", path, "error", err)
		return nil, err
	}

	c.logger.Debug("GitLab API response", "method", method, "path", path, "status", resp.StatusCode)
	return resp, nil
}

func (c *Client) doJSON(method, path string, body io.Reader, result any) error {
	resp, err := c.do(method, path, body)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		c.logger.Warn("GitLab API error", "status", resp.StatusCode, "body", string(respBody))
		return fmt.Errorf("GitLab API error (%d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}

func (c *Client) doJSONList(method, path string, result any) (int, error) {
	resp, err := c.do(method, path, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("GitLab API error (%d): %s", resp.StatusCode, string(respBody))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return 0, err
	}

	total := 0
	if v := resp.Header.Get("X-Total"); v != "" {
		total, _ = strconv.Atoi(v)
	}
	return total, nil
}

func encodeProjectID(projectID string) string {
	return url.PathEscape(projectID)
}

func jsonBody(v any) *strings.Reader {
	b, _ := json.Marshal(v)
	return strings.NewReader(string(b))
}

// --- Original 9 API methods ---

// GetDefaultBranch returns the default branch of a project.
func (c *Client) GetDefaultBranch(projectID string) (string, error) {
	var project Project
	path := fmt.Sprintf("/projects/%s", encodeProjectID(projectID))
	if err := c.doJSON("GET", path, nil, &project); err != nil {
		return "", err
	}
	return project.DefaultBranch, nil
}

// SearchProjects searches for GitLab projects.
func (c *Client) SearchProjects(query string, page, perPage int) (*SearchResult, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	params := url.Values{}
	params.Set("search", query)
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))

	var projects []Project
	total, err := c.doJSONList("GET", "/projects?"+params.Encode(), &projects)
	if err != nil {
		return nil, err
	}
	return &SearchResult{Count: total, Items: projects}, nil
}

// GetFileContents retrieves the contents of a file from a project.
func (c *Client) GetFileContents(projectID, filePath, ref string) (*FileContent, error) {
	encodedPath := url.PathEscape(filePath)
	path := fmt.Sprintf("/projects/%s/repository/files/%s", encodeProjectID(projectID), encodedPath)
	if ref != "" {
		path += "?ref=" + url.QueryEscape(ref)
	}

	var content FileContent
	if err := c.doJSON("GET", path, nil, &content); err != nil {
		return nil, err
	}

	// Decode base64 content
	if content.Content != "" {
		decoded, err := base64.StdEncoding.DecodeString(content.Content)
		if err == nil {
			content.Content = string(decoded)
		}
	}

	return &content, nil
}

// CreateOrUpdateFile creates or updates a file in a project.
func (c *Client) CreateOrUpdateFile(projectID, filePath, content, commitMessage, branch, previousPath string) (*FileResponse, error) {
	encodedPath := url.PathEscape(filePath)
	apiPath := fmt.Sprintf("/projects/%s/repository/files/%s", encodeProjectID(projectID), encodedPath)

	body := map[string]string{
		"branch":         branch,
		"content":        content,
		"commit_message": commitMessage,
	}
	if previousPath != "" {
		body["previous_path"] = previousPath
	}

	// Check if file exists to decide POST vs PUT
	method := "POST"
	_, getErr := c.GetFileContents(projectID, filePath, branch)
	if getErr == nil {
		method = "PUT"
	}

	var result FileResponse
	if err := c.doJSON(method, apiPath, jsonBody(body), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateCommit pushes multiple files in a single commit.
func (c *Client) CreateCommit(projectID, message, branch string, files []CommitAction) (*Commit, error) {
	path := fmt.Sprintf("/projects/%s/repository/commits", encodeProjectID(projectID))

	var commit Commit
	if err := c.doJSON("POST", path, jsonBody(map[string]any{
		"branch":         branch,
		"commit_message": message,
		"actions":        files,
	}), &commit); err != nil {
		return nil, err
	}
	return &commit, nil
}

// CreateProject creates a new GitLab project.
func (c *Client) CreateProject(name, description, visibility string, initReadme bool) (*Project, error) {
	body := map[string]any{"name": name}
	if description != "" {
		body["description"] = description
	}
	if visibility != "" {
		body["visibility"] = visibility
	}
	if initReadme {
		body["initialize_with_readme"] = true
	}

	var project Project
	if err := c.doJSON("POST", "/projects", jsonBody(body), &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateIssue creates a new issue in a project.
func (c *Client) CreateIssue(projectID, title, description string, assigneeIDs []int, milestoneID int, labels []string) (*Issue, error) {
	path := fmt.Sprintf("/projects/%s/issues", encodeProjectID(projectID))

	body := map[string]any{"title": title}
	if description != "" {
		body["description"] = description
	}
	if len(assigneeIDs) > 0 {
		body["assignee_ids"] = assigneeIDs
	}
	if milestoneID > 0 {
		body["milestone_id"] = milestoneID
	}
	if len(labels) > 0 {
		body["labels"] = strings.Join(labels, ",")
	}

	var issue Issue
	if err := c.doJSON("POST", path, jsonBody(body), &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

// CreateMergeRequest creates a new merge request in a project.
func (c *Client) CreateMergeRequest(projectID, title, description, sourceBranch, targetBranch string, draft, allowCollaboration bool) (*MergeRequest, error) {
	path := fmt.Sprintf("/projects/%s/merge_requests", encodeProjectID(projectID))

	body := map[string]any{
		"title":         title,
		"source_branch": sourceBranch,
		"target_branch": targetBranch,
	}
	if description != "" {
		body["description"] = description
	}
	if draft {
		body["draft"] = true
	}
	if allowCollaboration {
		body["allow_collaboration"] = true
	}

	var mr MergeRequest
	if err := c.doJSON("POST", path, jsonBody(body), &mr); err != nil {
		return nil, err
	}
	return &mr, nil
}

// ForkProject forks a project to your account or a specified namespace.
func (c *Client) ForkProject(projectID, namespace string) (*Project, error) {
	path := fmt.Sprintf("/projects/%s/fork", encodeProjectID(projectID))
	if namespace != "" {
		path += "?namespace=" + url.QueryEscape(namespace)
	}

	var project Project
	if err := c.doJSON("POST", path, nil, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

// CreateBranch creates a new branch in a project.
func (c *Client) CreateBranch(projectID, branchName, ref string) (*Branch, error) {
	path := fmt.Sprintf("/projects/%s/repository/branches", encodeProjectID(projectID))

	var branch Branch
	if err := c.doJSON("POST", path, jsonBody(map[string]string{
		"branch": branchName,
		"ref":    ref,
	}), &branch); err != nil {
		return nil, err
	}
	return &branch, nil
}

// --- New API methods ---

// ListIssues lists issues for a project with optional filtering.
func (c *Client) ListIssues(projectID, state string, page, perPage int) ([]Issue, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	if state != "" {
		params.Set("state", state)
	}

	path := fmt.Sprintf("/projects/%s/issues?%s", encodeProjectID(projectID), params.Encode())
	var issues []Issue
	total, err := c.doJSONList("GET", path, &issues)
	if err != nil {
		return nil, 0, err
	}
	return issues, total, nil
}

// GetIssue retrieves a single issue by IID.
func (c *Client) GetIssue(projectID string, issueIID int) (*Issue, error) {
	path := fmt.Sprintf("/projects/%s/issues/%d", encodeProjectID(projectID), issueIID)
	var issue Issue
	if err := c.doJSON("GET", path, nil, &issue); err != nil {
		return nil, err
	}
	return &issue, nil
}

// ListMergeRequests lists merge requests for a project with optional filtering.
func (c *Client) ListMergeRequests(projectID, state string, page, perPage int) ([]MergeRequest, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	if state != "" {
		params.Set("state", state)
	}

	path := fmt.Sprintf("/projects/%s/merge_requests?%s", encodeProjectID(projectID), params.Encode())
	var mrs []MergeRequest
	total, err := c.doJSONList("GET", path, &mrs)
	if err != nil {
		return nil, 0, err
	}
	return mrs, total, nil
}

// GetMergeRequest retrieves a single merge request by IID.
func (c *Client) GetMergeRequest(projectID string, mrIID int) (*MergeRequest, error) {
	path := fmt.Sprintf("/projects/%s/merge_requests/%d", encodeProjectID(projectID), mrIID)
	var mr MergeRequest
	if err := c.doJSON("GET", path, nil, &mr); err != nil {
		return nil, err
	}
	return &mr, nil
}

// AddNote adds a comment to an issue or merge request.
// noteableType should be "issues" or "merge_requests".
func (c *Client) AddNote(projectID, noteableType string, noteableIID int, body string) (*Note, error) {
	path := fmt.Sprintf("/projects/%s/%s/%d/notes", encodeProjectID(projectID), noteableType, noteableIID)
	var note Note
	if err := c.doJSON("POST", path, jsonBody(map[string]string{"body": body}), &note); err != nil {
		return nil, err
	}
	return &note, nil
}

// ListPipelines lists pipelines for a project.
func (c *Client) ListPipelines(projectID, ref, status string, page, perPage int) ([]Pipeline, int, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("per_page", strconv.Itoa(perPage))
	if ref != "" {
		params.Set("ref", ref)
	}
	if status != "" {
		params.Set("status", status)
	}

	path := fmt.Sprintf("/projects/%s/pipelines?%s", encodeProjectID(projectID), params.Encode())
	var pipelines []Pipeline
	total, err := c.doJSONList("GET", path, &pipelines)
	if err != nil {
		return nil, 0, err
	}
	return pipelines, total, nil
}

// GetPipeline retrieves a single pipeline by ID.
func (c *Client) GetPipeline(projectID string, pipelineID int) (*Pipeline, error) {
	path := fmt.Sprintf("/projects/%s/pipelines/%d", encodeProjectID(projectID), pipelineID)
	var pipeline Pipeline
	if err := c.doJSON("GET", path, nil, &pipeline); err != nil {
		return nil, err
	}
	return &pipeline, nil
}
