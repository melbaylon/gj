package ls

import (
	"fmt"
	"os"
	"path/filepath"
	"slices" // For Go 1.21+
	"sort"
)

// List scans the directory at the given path and prints its contents.
// This function will serve as the entry point for the ls logic and will be
// expanded throughout the project milestones.
func List(path string, all bool, long bool, sortByTime bool, sortBySize bool, reverse bool) error {
	// Task 1.3: Core Directory Reading
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var fileEntries []*FileEntry
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
	}

	// Task 2.3: Sorting Engine
	sortFiles(fileEntries, sortByTime, sortBySize, reverse)

	// For now, just print the names. Sorting and proper formatting will come later.
	for _, entry := range fileEntries {
		// Milestone 3 will implement proper tabular formatting for 'long' mode.
		if long {
			fmt.Printf("Metadata placeholder for: %s\n", entry.Name)
		} else {
			fmt.Printf("%s  ", entry.Name)
		}
	}

	if !long {
		fmt.Println()
	}

	return nil
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
		return sort.String(a.Name).Compare(sort.String(b.Name))
	})

	if reverse {
		slices.Reverse(entries)
	}
}
