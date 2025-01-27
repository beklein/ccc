package ccc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/beklein/ccc/config"
	"github.com/beklein/ccc/fileutils"
)

// RunCCC implements reading your .ccc, optionally *not* using .gitignore, etc.
func RunCCC(configPath string, outputToStdout bool) error {
	// Read the config file paths
	lines, err := config.ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	var context strings.Builder

	for _, line := range lines {
		// Skip comments and empty lines
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches, err := filepath.Glob(line)
		// If no glob match or error, treat line as literal path
		if err != nil || matches == nil {
			matches = []string{line}
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				log.Printf("Skipping: %s (stat error: %v)\n", match, err)
				continue
			}
			if info.IsDir() {
				// Recursively gather files from the directory
				filepaths, err := fileutils.GatherFiles(match)
				if err != nil {
					log.Printf("Skipping directory: %s (error: %v)\n", match, err)
					continue
				}
				for _, f := range filepaths {
					fileutils.AppendFileContent(&context, f)
				}
			} else {
				fileutils.AppendFileContent(&context, match)
			}
		}
	}

	if outputToStdout {
		fmt.Println(context.String())
	} else {
		// Copy to clipboard
		if err := clipboard.WriteAll(context.String()); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %v", err)
		}
	}

	return nil
}
