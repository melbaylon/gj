package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/melbaylon/gj/internal/ls"
)

const version = "1.0.0"

func main() {
	all := flag.Bool("a", false, "do not ignore entries starting with .")
	long := flag.Bool("l", false, "use a long listing format")
	sortByTime := flag.Bool("t", false, "sort by modification time")
	sortBySize := flag.Bool("S", false, "sort by file size")
	reverse := flag.Bool("r", false, "reverse order while sorting")
	classify := flag.Bool("F", false, "append indicator (one of */=>@|) to entries")
	color := flag.String("color", "auto", "colorize the output; WHEN can be 'always', 'auto' (default), or 'never'")
	recursive := flag.Bool("R", false, "list subdirectories recursively")
	humanReadable := flag.Bool("h", false, "with -l, print sizes like 1K 234M 2G etc.")
	showVersion := flag.Bool("v", false, "display version information and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gs [OPTION]... [FILE]...\n")
		fmt.Fprintf(os.Stderr, "List information about the FILEs (the current directory by default).\n\n")
		fmt.Fprintf(os.Stderr, "Sort entries alphabetically if none of -cftuSUX nor --sort is specified.\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  gs -l -h          List files in long format with human-readable sizes.\n")
		fmt.Fprintf(os.Stderr, "  gs -R             Recursively list all files and subdirectories.\n")
		fmt.Fprintf(os.Stderr, "  gs --color=always Force colorized output.\n")
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("gs version %s\n", version)
		return
	}

	// Remaining arguments are the paths to list
	paths := flag.Args()

	// Task 5.1: Custom 'help' command
	if len(paths) > 0 && paths[0] == "help" {
		flag.Usage()
		return
	}

	if len(paths) == 0 {
		paths = []string{"."}
	}

	for i, path := range paths {
		if len(paths) > 1 {
			fmt.Printf("%s:\n", path)
		}
		// Basic orchestration for Milestone 1
		err := ls.List(path, *all, *long, *sortByTime, *sortBySize, *reverse, *classify, *color, *recursive, *humanReadable)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: %s: %v\n", path, err)
		}
		if i < len(paths)-1 {
			fmt.Println()
		}
	}
}
