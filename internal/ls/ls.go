package ls

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices" // For Go 1.21+
	"strings"
	"text/tabwriter"
	"time"

	"golang.org/x/term"
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
func List(path string, all bool, long bool, sortByTime bool, sortBySize bool, reverse bool, classify bool, colorMode string, recursive bool, humanReadable bool) error {
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

	// Task 4.4: Terminal Detection
	isTTY := term.IsTerminal(int(os.Stdout.Fd()))
	shouldColor := false
	if colorMode == "always" || (colorMode == "auto" && isTTY) {
		shouldColor = true
	}

	if long {
		// Task 3.4: Tabular Long Listing
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		for _, entry := range fileEntries {
			displayName := entry.Name
			if shouldColor {
				displayName = applyColor(entry.Name, entry.Mode)
			}
			if classify {
				displayName += getIndicator(entry.Mode)
			}

			modeStr := FormatMode(entry.Mode)
			sizeStr := formatSize(entry.Size, humanReadable)
			timeStr := formatTime(entry.ModTime)

			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%s\t%s\n",
				modeStr, entry.Nlink, entry.Owner, entry.Group, sizeStr, timeStr, displayName)
		}
		w.Flush()
	} else if !isTTY {
		// Task 4.5: If not a TTY, print one per line or simple space-separated
		// Default behavior of ls without a TTY is often one-per-line (e.g., when piped)
		for _, entry := range fileEntries {
			displayName := entry.Name
			if shouldColor {
				displayName = applyColor(entry.Name, entry.Mode)
			}
			if classify {
				displayName += getIndicator(entry.Mode)
			}
			fmt.Println(displayName)
		}
	} else {
		// Task 4.5: Multi-column Formatting for TTY
		printColumns(fileEntries, shouldColor, classify)
	}

	// Task 4.3: Recursive Listing
	if recursive {
		// Sort subDirs if they were collected (they might already be sorted if fileEntries was)
		// but it's safer to ensure they match the current sorting criteria.
		sortFiles(subDirs, sortByTime, sortBySize, reverse)

		for _, subDir := range subDirs {
			fmt.Printf("\n%s:\n", subDir.Path)
			err := List(subDir.Path, all, long, sortByTime, sortBySize, reverse, classify, colorMode, recursive, humanReadable)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ls: %s: %v\n", subDir.Path, err)
			}
		}
	}

	return nil
}

// formatSize converts bytes into human-readable strings if enabled.
func formatSize(size int64, human bool) string {
	if !human {
		return fmt.Sprintf("%d", size)
	}
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(size)/float64(div), "KMGTP"[exp])
}

// formatTime formats the modification time according to ls standards.
func formatTime(t time.Time) string {
	now := time.Now()
	// If the time is older than 6 months or in the future, show year instead of time
	sixMonthsAgo := now.AddDate(0, -6, 0)
	if t.Before(sixMonthsAgo) || t.After(now) {
		return t.Format("Jan _2  2006")
	}
	return t.Format("Jan _2 15:04")
}

// printColumns implements a grid-based layout for TTY output.
func printColumns(entries []*FileEntry, shouldColor bool, classify bool) {
	if len(entries) == 0 {
		return
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 80 // Fallback
	}

	// Prepare list of display names and calculate maximum length
	maxLen := 0
	type item struct {
		display string
		length  int
	}
	items := make([]item, len(entries))

	for i, entry := range entries {
		display := entry.Name
		length := len(display)
		if classify {
			ind := getIndicator(entry.Mode)
			display += ind
			length += len(ind)
		}
		if shouldColor {
			display = applyColor(display, entry.Mode)
		}
		items[i] = item{display, length}
		if length > maxLen {
			maxLen = length
		}
	}

	// Task 4.5: Layout logic
	columnWidth := maxLen + 2
	cols := width / columnWidth
	if cols <= 0 {
		cols = 1
	}
	rows := (len(items) + cols - 1) / cols

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			idx := c*rows + r
			if idx < len(items) {
				padding := ""
				if c < cols-1 {
					padding = strings.Repeat(" ", columnWidth-items[idx].length)
				}
				fmt.Printf("%s%s", items[idx].display, padding)
			}
		}
		fmt.Println()
	}
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
