package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/impzero/creeder/pkg/file"
)

func TestReadFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	content := "hello, world!"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Errorf("failed to write content to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Errorf("failed to close temp file: %v", err)
	}

	readContent, err := file.ReadFile(tmpfile.Name())
	if err != nil {
		t.Errorf("ReadFile(%s) returned error: %v", tmpfile.Name(), err)
	}
	if readContent != content {
		t.Errorf("ReadFile(%s) = %q, want %q", tmpfile.Name(), readContent, content)
	}
}

func TestShouldIgnorePath(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		ignore  string
		wantRes bool
	}{
		{
			name:    "empty ignore",
			path:    "/path/to/file",
			ignore:  "",
			wantRes: false,
		},
		{
			name:    "exact match",
			path:    "/path/to/file",
			ignore:  "file",
			wantRes: true,
		},
		{
			name:    "prefix match",
			path:    "/path/to/file",
			ignore:  "/path/to",
			wantRes: true,
		},
		{
			name:    "no match",
			path:    "/path/to/file",
			ignore:  "otherfile",
			wantRes: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotRes := file.ShouldIgnorePath(tc.path, tc.ignore)
			if gotRes != tc.wantRes {
				t.Errorf("ShouldIgnorePath(%q, %q) = %v, want %v", tc.path, tc.ignore, gotRes, tc.wantRes)
			}
		})
	}
}

func TestShouldIncludeFile(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		filter  string
		wantRes bool
	}{
		{
			name:    "empty filter",
			path:    "/path/to/file",
			filter:  "",
			wantRes: true,
		},
		{
			name:    "matching extension",
			path:    "/path/to/file.go",
			filter:  "go",
			wantRes: true,
		},
		{
			name:    "non-matching extension",
			path:    "/path/to/file.txt",
			filter:  "go",
			wantRes: false,
		},
		{
			name:    "multiple extensions",
			path:    "/path/to/file.txt",
			filter:  "go,txt",
			wantRes: true,
		},
		{
			name:    "no match",
			path:    "/path/to/file",
			filter:  "go",
			wantRes: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotRes := file.ShouldIncludeFile(tc.path, tc.filter)
			if gotRes != tc.wantRes {
				t.Errorf("ShouldIncludeFile(%q, %q) = %v, want %v", tc.path, tc.filter, gotRes, tc.wantRes)
			}
		})
	}
}

func TestReadFilesFromPath(t *testing.T) {
	root := "test"
	err := os.Mkdir(root, os.ModePerm)
	if err != nil {
		t.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(root)

	// create test files
	file1 := filepath.Join(root, "file1.go")
	file2 := filepath.Join(root, "file2.txt")
	file3 := filepath.Join(root, "dir/file3.go")
	file4 := filepath.Join(root, "dir/file4.txt")

	if err := os.Mkdir(filepath.Dir(file3), os.ModePerm); err != nil {
		t.Errorf("failed to create dir for test files: %v", err)
	}

	files := []struct {
		path    string
		content string
	}{
		{path: file1, content: "package main\n\nfunc main() {\n}\n"},
		{path: file2, content: "this is a text file"},
		{path: file3, content: "package main\n\nfunc test() {\n}\n"},
		{path: file4, content: "another text file"},
	}

	for _, f := range files {
		if err := os.WriteFile(f.path, []byte(f.content), 0644); err != nil {
			t.Errorf("failed to write content to test file %q: %v", f.path, err)
		}
	}

	// test list files with no filters
	gotFiles := []string{}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if file.ShouldIncludeFile(path, "") {
			if info.IsDir() {
				return nil
			}
			gotFiles = append(gotFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Errorf("filePath.Walk returned error: %v", err)
	}

	wantFiles := []string{file1, file2, file3, file4}
	if len(gotFiles) != len(wantFiles) {
		t.Errorf("ShouldIncludeFile(%q, \"\") = %v, want %v", root, gotFiles, wantFiles)
	}

	for _, f := range wantFiles {
		if !contains(gotFiles, f) {
			t.Errorf("ShouldIncludeFiles(%q, \"\") = %v, does not contain %q", root, gotFiles, f)
		}
	}

	// test list files with go filter
	gotFiles = []string{}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if file.ShouldIncludeFile(path, "go") {
			gotFiles = append(gotFiles, path)
		}
		return nil
	})
	if err != nil {
		t.Errorf("filePath.Walk returned error: %v", err)
	}

	wantFiles = []string{file1, file3}
	if len(gotFiles) != len(wantFiles) {
		t.Errorf("ListFiles(%q, \"go\") = %v, want %v", root, gotFiles, wantFiles)
	}

	for _, f := range wantFiles {
		if !contains(gotFiles, f) {
			t.Errorf("ListFiles(%q, \"go\") = %v, does not contain %q", root, gotFiles, f)
		}
	}
}

func contains(files []string, path string) bool {
	for _, f := range files {
		if f == path {
			return true
		}
	}
	return false
}
