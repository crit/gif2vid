package inputs

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestGetFilesFromDir(t *testing.T) {
	tmp := t.TempDir()

	// Create some files
	files := []string{"a.gif", "b.webp", "c.txt", "d.GIF"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmp, f), []byte("dummy"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	// Create a subdir
	if err := os.Mkdir(filepath.Join(tmp, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}

	got, err := GetFilesFromDir(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 files, got %d: %v", len(got), got)
	}

	// Paths are absolute and cleaned, let's just check basenames for simplicity
	var names []string
	for _, p := range got {
		names = append(names, filepath.Base(p))
	}
	sort.Strings(names)

	want := []string{"a.gif", "b.webp", "d.GIF"}
	for i, v := range want {
		if names[i] != v {
			t.Errorf("got %q, want %q", names[i], v)
		}
	}
}

func TestGetFilesFromDir_Empty(t *testing.T) {
	tmp := t.TempDir()
	_, err := GetFilesFromDir(tmp)
	if err == nil {
		t.Error("expected error for empty directory")
	}
}

func TestGetFilesFromDir_NotADir(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "file.txt")
	os.WriteFile(f, []byte("hi"), 0644)

	_, err := GetFilesFromDir(f)
	if err == nil {
		t.Error("expected error for non-directory")
	}
}
