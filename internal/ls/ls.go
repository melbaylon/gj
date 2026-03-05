package ls

import (
	"fmt"
	"os"
)

// List scans the directory at the given path and prints its contents.
// This function will serve as the entry point for the ls logic and will be
// expanded throughout the project milestones.
func List(path string, all bool, long bool) error {
	// Task 1.3: Core Directory Reading
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// Task 2.2: Basic hidden file filtering logic
		if !all && len(entry.Name()) > 0 && entry.Name()[0] == '.' {
			continue
		}

		// Initial display logic
		// Milestone 3 will implement proper tabular formatting for 'long' mode.
		if long {
			fmt.Printf("Metadata placeholder for: %s\n", entry.Name())
		} else {
			fmt.Printf("%s  ", entry.Name())
		}
	}

	if !long {
		fmt.Println()
	}

	return nil
}
