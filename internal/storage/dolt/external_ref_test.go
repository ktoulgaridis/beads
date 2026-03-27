package dolt

import (
	"testing"

	"github.com/steveyegge/beads/internal/types"
)

func TestGetIssueByExternalRef_FindsWisp(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx, cancel := testContext(t)
	defer cancel()

	extRef := "https://gitlab.example.com/group/project/-/issues/42"
	wisp := &types.Issue{
		Title:       "GitLab issue stored as wisp",
		ExternalRef: &extRef,
		Status:      types.StatusOpen,
		Priority:    3,
		IssueType:   types.TypeTask,
		Ephemeral:   true,
	}
	if err := store.CreateIssue(ctx, wisp, "tester"); err != nil {
		t.Fatalf("CreateIssue (wisp) failed: %v", err)
	}

	got, err := store.GetIssueByExternalRef(ctx, extRef)
	if err != nil {
		t.Fatalf("GetIssueByExternalRef failed for wisp: %v", err)
	}
	if got == nil {
		t.Fatal("GetIssueByExternalRef returned nil for wisp with matching external_ref")
	}
	if got.ID != wisp.ID {
		t.Errorf("expected ID %s, got %s", wisp.ID, got.ID)
	}
}

func TestGetIssueByExternalRef_FindsRegularIssue(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx, cancel := testContext(t)
	defer cancel()

	extRef := "https://gitlab.example.com/group/project/-/issues/99"
	issue := &types.Issue{
		ID:          "ext-regular-abc",
		Title:       "Regular issue with external ref",
		ExternalRef: &extRef,
		Status:      types.StatusOpen,
		Priority:    2,
		IssueType:   types.TypeBug,
	}
	if err := store.CreateIssue(ctx, issue, "tester"); err != nil {
		t.Fatalf("CreateIssue failed: %v", err)
	}

	got, err := store.GetIssueByExternalRef(ctx, extRef)
	if err != nil {
		t.Fatalf("GetIssueByExternalRef failed: %v", err)
	}
	if got == nil {
		t.Fatal("GetIssueByExternalRef returned nil for regular issue with matching external_ref")
	}
	if got.ID != issue.ID {
		t.Errorf("expected ID %s, got %s", issue.ID, got.ID)
	}
}

func TestGetIssueByExternalRef_NotFound(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	ctx, cancel := testContext(t)
	defer cancel()

	got, err := store.GetIssueByExternalRef(ctx, "https://nonexistent.example.com/issue/1")
	if err == nil {
		t.Fatal("expected error for nonexistent external ref")
	}
	if got != nil {
		t.Errorf("expected nil issue, got %+v", got)
	}
}
