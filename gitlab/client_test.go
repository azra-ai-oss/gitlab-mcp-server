package gitlab

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient creates a Client pointing at a test HTTP server.
func newTestClient(handler http.Handler) (*Client, *httptest.Server) {
	ts := httptest.NewServer(handler)
	client := NewClient(ts.URL, "test-token", nil)
	return client, ts
}

func TestSearchProjects(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("search") != "myproject" {
			t.Errorf("unexpected search: %s", r.URL.Query().Get("search"))
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}
		w.Header().Set("X-Total", "1")
		_ = json.NewEncoder(w).Encode([]Project{
			{ID: 1, Name: "myproject", DefaultBranch: "main"},
		})
	})

	client, ts := newTestClient(handler)
	defer ts.Close()

	result, err := client.SearchProjects("myproject", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	if result.Items[0].Name != "myproject" {
		t.Errorf("expected name myproject, got %s", result.Items[0].Name)
	}
}

func TestGetFileContents(t *testing.T) {
	content := base64.StdEncoding.EncodeToString([]byte("Hello, World!"))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(FileContent{
			FileName: "README.md",
			FilePath: "README.md",
			Size:     13,
			Encoding: "base64",
			Content:  content,
		})
	})

	client, ts := newTestClient(handler)
	defer ts.Close()

	result, err := client.GetFileContents("mygroup/myproject", "README.md", "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Content != "Hello, World!" {
		t.Errorf("expected decoded content 'Hello, World!', got '%s'", result.Content)
	}
}

func TestGetDefaultBranch(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(Project{
			ID:            42,
			Name:          "test",
			DefaultBranch: "develop",
		})
	})

	client, ts := newTestClient(handler)
	defer ts.Close()

	branch, err := client.GetDefaultBranch("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch != "develop" {
		t.Errorf("expected 'develop', got '%s'", branch)
	}
}

func TestCreateIssue(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "Bug report" {
			t.Errorf("expected title 'Bug report', got '%v'", body["title"])
		}

		_ = json.NewEncoder(w).Encode(Issue{
			ID:    1,
			IID:   1,
			Title: "Bug report",
		})
	})

	client, ts := newTestClient(handler)
	defer ts.Close()

	issue, err := client.CreateIssue("42", "Bug report", "Something broke", nil, 0, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.Title != "Bug report" {
		t.Errorf("expected title 'Bug report', got '%s'", issue.Title)
	}
}

func TestListIssues(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != "opened" {
			t.Errorf("expected state=opened, got %s", r.URL.Query().Get("state"))
		}
		w.Header().Set("X-Total", "2")
		_ = json.NewEncoder(w).Encode([]Issue{
			{ID: 1, IID: 1, Title: "First"},
			{ID: 2, IID: 2, Title: "Second"},
		})
	})

	client, ts := newTestClient(handler)
	defer ts.Close()

	issues, total, err := client.ListIssues("42", "opened", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(issues))
	}
}

func TestGetMergeRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(MergeRequest{
			ID:           10,
			IID:          5,
			Title:        "Add feature",
			SourceBranch: "feature",
			TargetBranch: "main",
		})
	})

	client, ts := newTestClient(handler)
	defer ts.Close()

	mr, err := client.GetMergeRequest("42", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mr.Title != "Add feature" {
		t.Errorf("expected title 'Add feature', got '%s'", mr.Title)
	}
	if mr.SourceBranch != "feature" {
		t.Errorf("expected source branch 'feature', got '%s'", mr.SourceBranch)
	}
}

func TestListPipelines(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Total", "1")
		_ = json.NewEncoder(w).Encode([]Pipeline{
			{ID: 100, Status: "success", Ref: "main"},
		})
	})

	client, ts := newTestClient(handler)
	defer ts.Close()

	pipelines, total, err := client.ListPipelines("42", "", "", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if pipelines[0].Status != "success" {
		t.Errorf("expected status 'success', got '%s'", pipelines[0].Status)
	}
}

func TestAPIError(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "404 Project Not Found"}`))
	})

	client, ts := newTestClient(handler)
	defer ts.Close()

	_, err := client.GetDefaultBranch("nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAddNote(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["body"] != "LGTM!" {
			t.Errorf("expected body 'LGTM!', got '%v'", body["body"])
		}

		_ = json.NewEncoder(w).Encode(Note{
			ID:   1,
			Body: "LGTM!",
			Author: Author{
				ID:       1,
				Username: "testuser",
				Name:     "Test User",
			},
		})
	})

	client, ts := newTestClient(handler)
	defer ts.Close()

	note, err := client.AddNote("42", "merge_requests", 5, "LGTM!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Body != "LGTM!" {
		t.Errorf("expected body 'LGTM!', got '%s'", note.Body)
	}
}
