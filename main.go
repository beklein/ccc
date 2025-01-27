package main

import (
	"flag"
	"log"

	"github.com/beklein/ccc/ccc"
)

func main() {
	configPath := flag.String("c", ".ccc", "Path to the .ccc configuration file.")
	outputToStdout := flag.Bool("o", false, "Print output to stdout instead of copying to clipboard.")
	flag.Parse()

	if err := ccc.RunCCC(*configPath, *outputToStdout); err != nil {
		log.Fatal(err)
	}
}
