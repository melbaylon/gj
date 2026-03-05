package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/gj/internal/ls"
)

func main() {
	// Task 1.2 will expand these flags
	all := flag.Bool("a", false, "do not ignore entries starting with .")
	long := flag.Bool("l", false, "use a long listing format")

	flag.Parse()

	// Remaining arguments are the paths to list
	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, path := range paths {
		// Basic orchestration for Milestone 1
		err := ls.List(path, *all, *long)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", path, err)
		}
	}
}
