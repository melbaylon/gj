package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/melbaylon/gj/internal/ls"
)

func main() {
	// Task 1.2 will expand these flags
	all := flag.Bool("a", false, "do not ignore entries starting with .")
	long := flag.Bool("l", false, "use a long listing format")
	sortByTime := flag.Bool("t", false, "sort by modification time")
	sortBySize := flag.Bool("S", false, "sort by file size") // Note: original ls uses -S for size
	reverse := flag.Bool("r", false, "reverse order while sorting")
	classify := flag.Bool("F", false, "append indicator (one of */=>@|) to entries")
	color := flag.String("color", "always", "colorize the output; WHEN can be 'always' (default), 'auto', or 'never'")
	recursive := flag.Bool("R", false, "list subdirectories recursively")

	flag.Parse()

	// Remaining arguments are the paths to list
	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	for i, path := range paths {
		if len(paths) > 1 {
			fmt.Printf("%s:\n", path)
		}
		// Basic orchestration for Milestone 1
		err := ls.List(path, *all, *long, *sortByTime, *sortBySize, *reverse, *classify, *color, *recursive)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", path, err)
		}
		if i < len(paths)-1 {
			fmt.Println()
		}
	}
}
