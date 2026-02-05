package bookmarks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBookmarkStore(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "stui-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Override config dir
	store := &Store{
		path:      filepath.Join(tmpDir, "bookmarks.json"),
		bookmarks: []Bookmark{},
	}

	// Test Add
	bm, err := store.Add("test-bookmark", "my-bucket", "some/prefix/")
	if err != nil {
		t.Fatalf("failed to add bookmark: %v", err)
	}
	if bm.Name != "test-bookmark" {
		t.Errorf("expected name 'test-bookmark', got '%s'", bm.Name)
	}
	if bm.Bucket != "my-bucket" {
		t.Errorf("expected bucket 'my-bucket', got '%s'", bm.Bucket)
	}

	// Test List
	list := store.List()
	if len(list) != 1 {
		t.Errorf("expected 1 bookmark, got %d", len(list))
	}

	// Test Get
	found, ok := store.Get(bm.ID)
	if !ok {
		t.Error("bookmark not found")
	}
	if found.Name != bm.Name {
		t.Errorf("expected name '%s', got '%s'", bm.Name, found.Name)
	}

	// Test DisplayName
	if bm.DisplayName() != "test-bookmark" {
		t.Errorf("expected display name 'test-bookmark', got '%s'", bm.DisplayName())
	}

	// Test Path
	expectedPath := "s3://my-bucket/some/prefix/"
	if bm.Path() != expectedPath {
		t.Errorf("expected path '%s', got '%s'", expectedPath, bm.Path())
	}

	// Test Remove
	err = store.Remove(bm.ID)
	if err != nil {
		t.Fatalf("failed to remove bookmark: %v", err)
	}
	list = store.List()
	if len(list) != 0 {
		t.Errorf("expected 0 bookmarks after remove, got %d", len(list))
	}
}

func TestBookmarkDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		bookmark Bookmark
		expected string
	}{
		{
			name:     "with name",
			bookmark: Bookmark{Name: "My Bookmark", Bucket: "bucket", Prefix: "prefix/"},
			expected: "My Bookmark",
		},
		{
			name:     "without name, with prefix",
			bookmark: Bookmark{Name: "", Bucket: "bucket", Prefix: "prefix/"},
			expected: "s3://bucket/prefix/",
		},
		{
			name:     "without name, without prefix",
			bookmark: Bookmark{Name: "", Bucket: "bucket", Prefix: ""},
			expected: "s3://bucket",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bookmark.DisplayName(); got != tt.expected {
				t.Errorf("DisplayName() = %s, want %s", got, tt.expected)
			}
		})
	}
}
