package ccc

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureOutput helps us see what RunCCC prints to stdout.
func captureOutput(runFunc func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runFunc()

	w.Close()
	os.Stdout = old

	outBytes, _ := io.ReadAll(r)
	return string(outBytes), err
}

func TestReadCCCFile(t *testing.T) {
	// Make a temp dir and chdir into it
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	tempDir, err := os.MkdirTemp("", "ccc-test-read-ccc-file")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to chdir to temp dir: %v", err)
	}

	// Create .ccc with comments + actual entries
	cccPath := filepath.Join(tempDir, ".ccc")
	cccContent := `
# This is a comment
some_file.xyz
src/*.go
README.md
`
	if err := os.WriteFile(cccPath, []byte(cccContent), 0644); err != nil {
		t.Fatalf("Failed to write .ccc file: %v", err)
	}

	// Create a file matching "some_file.xyz"
	filePath := filepath.Join(tempDir, "some_file.xyz")
	if err := os.WriteFile(filePath, []byte("Contents of some_file.xyz"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create a folder named "src" + Go files
	srcFolderPath := filepath.Join(tempDir, "src")
	if err := os.MkdirAll(srcFolderPath, 0755); err != nil {
		t.Fatalf("Failed to create src folder: %v", err)
	}
	srcMainPath := filepath.Join(srcFolderPath, "main.go")
	if err := os.WriteFile(srcMainPath, []byte("package main\nfunc main() {}"), 0644); err != nil {
		t.Fatalf("Failed to write main.go: %v", err)
	}
	srcUtilsPath := filepath.Join(srcFolderPath, "utils.go")
	if err := os.WriteFile(srcUtilsPath, []byte("package main\nfunc Util() {}"), 0644); err != nil {
		t.Fatalf("Failed to write utils.go: %v", err)
	}

	// Create README.md
	readmePath := filepath.Join(tempDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Project Documentation"), 0644); err != nil {
		t.Fatalf("Failed to write README.md: %v", err)
	}

	output, runErr := captureOutput(func() error {
		return RunCCC(cccPath, true)
	})
	if runErr != nil {
		t.Fatalf("RunCCC returned an error: %v\nOutput: %s", runErr, output)
	}

	if !strings.Contains(output, "some_file.xyz") {
		t.Errorf("Expected 'some_file.xyz' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "Contents of some_file.xyz") {
		t.Errorf("Expected 'Contents of some_file.xyz' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "src/main.go") {
		t.Errorf("Expected 'src/main.go' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "package main") {
		t.Errorf("Expected 'package main' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "src/utils.go") {
		t.Errorf("Expected 'src/utils.go' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "func Util() {}") {
		t.Errorf("Expected 'func Util() {}' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "README.md") {
		t.Errorf("Expected 'README.md' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "# Project Documentation") {
		t.Errorf("Expected '# Project Documentation' in output, got:\n%s", output)
	}
}
