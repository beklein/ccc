package fileutils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GatherFiles recursively collects all files under a directory,
func GatherFiles(dir string) ([]string, error) {
	var collected []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			collected = append(collected, path)
		}
		return nil
	})
	return collected, err
}

// AppendFileContent opens a file, reads it, and appends its filename + content.
func AppendFileContent(builder *strings.Builder, path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Skipping file: %s (read error: %v)\n", path, err)
		return
	}
	builder.WriteString(fmt.Sprintf("// %s\n", path))
	builder.Write(content)
	builder.WriteString("\n")
}
