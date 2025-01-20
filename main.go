package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/beklein/ccc/config"
	"github.com/beklein/ccc/fileutils"
)

func main() {
	configPath := flag.String("config", ".ccc", "Path to the .ccc configuration file.")
	outputToStdout := flag.Bool("o", false, "Print output to stdout instead of copying to clipboard.")
	flag.Parse()

	// Read the config file paths
	lines, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error reading ccc config file: %v\n", err)
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

	if *outputToStdout {
		fmt.Println(context.String())
	} else {
		// Copy to clipboard
		if err := clipboard.WriteAll(context.String()); err != nil {
			log.Fatalf("Failed to copy to clipboard: %v\n", err)
		}
	}
}
