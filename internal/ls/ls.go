package ls

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices" // For Go 1.21+
)

const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

// List scans the directory at the given path and prints its contents.
// This function will serve as the entry point for the ls logic and will be
// expanded throughout the project milestones.
func List(path string, all bool, long bool, sortByTime bool, sortBySize bool, reverse bool, classify bool, colorMode string, recursive bool) error {
	// Task 1.3: Core Directory Reading
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var fileEntries []*FileEntry
	var subDirs []*FileEntry

	for _, dirEntry := range dirEntries {
		// Task 2.2: Basic hidden file filtering logic
		if !all && len(dirEntry.Name()) > 0 && dirEntry.Name()[0] == '.' {
			continue
		}

		fullPath := filepath.Join(path, dirEntry.Name())
		fileEntry, err := NewFileEntry(fullPath, dirEntry)
		if err != nil {
			// fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", fullPath, err)
			continue // Skip files we can't get info for (e.g., broken symlinks)
		}
		fileEntries = append(fileEntries, fileEntry)

		if recursive && fileEntry.IsDir {
			// Standard ls -R ignores . and .. even with -a for the recursive part
			if fileEntry.Name != "." && fileEntry.Name != ".." {
				subDirs = append(subDirs, fileEntry)
			}
		}
	}

	// Task 2.3: Sorting Engine
	sortFiles(fileEntries, sortByTime, sortBySize, reverse)

	// For now, just print the names. Sorting and proper formatting will come later.
	for _, entry := range fileEntries {
		displayName := entry.Name
		if colorMode == "always" { // Simple check for now, Task 4.4 will handle 'auto'
			displayName = applyColor(entry.Name, entry.Mode)
		}

		if classify {
			displayName += getIndicator(entry.Mode)
		}

		// Milestone 3 will implement proper tabular formatting for 'long' mode.
		if long {
			fmt.Printf("Metadata placeholder for: %s\n", displayName)
		} else {
			fmt.Printf("%s  ", displayName)
		}
	}

	if !long {
		fmt.Println()
	}

	// Task 4.3: Recursive Listing
	if recursive {
		// Sort subDirs if they were collected (they might already be sorted if fileEntries was)
		// but it's safer to ensure they match the current sorting criteria.
		sortFiles(subDirs, sortByTime, sortBySize, reverse)

		for _, subDir := range subDirs {
			fmt.Printf("\n%s:\n", subDir.Path)
			err := List(subDir.Path, all, long, sortByTime, sortBySize, reverse, classify, colorMode, recursive)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ls: %s: %v\n", subDir.Path, err)
			}
		}
	}

	return nil
}

// applyColor returns the name wrapped in ANSI color codes based on file mode.
func applyColor(name string, mode os.FileMode) string {
	if mode.IsDir() {
		return colorBlue + name + colorReset
	}
	if mode&os.ModeSymlink != 0 {
		return colorCyan + name + colorReset
	}
	if mode&os.ModeSocket != 0 {
		return colorRed + name + colorReset
	}
	if mode&os.ModeNamedPipe != 0 {
		return colorYellow + name + colorReset
	}
	if mode&0111 != 0 {
		return colorGreen + name + colorReset
	}
	return name
}

// getIndicator returns the type-specific character based on file mode.
func getIndicator(mode os.FileMode) string {
	if mode.IsDir() {
		return "/"
	}
	if mode&os.ModeSymlink != 0 {
		return "@"
	}
	if mode&os.ModeSocket != 0 {
		return "="
	}
	if mode&os.ModeNamedPipe != 0 {
		return "|"
	}
	// Check for executable bits (owner, group, others)
	if mode&0111 != 0 {
		return "*"
	}
	return ""
}

// sortFiles sorts the slice of FileEntry based on the provided flags.
func sortFiles(entries []*FileEntry, sortByTime bool, sortBySize bool, reverse bool) {
	slices.SortFunc(entries, func(a, b *FileEntry) int {
		if sortByTime {
			if a.ModTime.Before(b.ModTime) {
				return -1
			}
			if a.ModTime.After(b.ModTime) {
				return 1
			}
		} else if sortBySize {
			if a.Size < b.Size {
				return -1
			}
			if a.Size > b.Size {
				return 1
			}
		}
		// Default: sort by name
		return cmp.Compare(a.Name, b.Name)
	})

	if reverse {
		slices.Reverse(entries)
	}
}
