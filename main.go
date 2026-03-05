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
	sortByTime := flag.Bool("t", false, "sort by modification time")
	sortBySize := flag.Bool("S", false, "sort by file size") // Note: original ls uses -S for size
	reverse := flag.Bool("r", false, "reverse order while sorting")
	classify := flag.Bool("F", false, "append indicator (one of */=>@|) to entries")

	flag.Parse()

	// Remaining arguments are the paths to list
	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for _, path := range paths {
		// Basic orchestration for Milestone 1
		err := ls.List(path, *all, *long, *sortByTime, *sortBySize, *reverse, *classify)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", path, err)
		}
	}
}
