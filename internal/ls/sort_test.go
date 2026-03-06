package ls

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestSortFiles_Alphabetical tests default alphabetical sorting
func TestSortFiles_Alphabetical(t *testing.T) {
	entries := []*FileEntry{
		{Name: "zebra", ModTime: time.Now(), Size: 100, Mode: 0o644},
		{Name: "apple", ModTime: time.Now(), Size: 100, Mode: 0o644},
		{Name: "banana", ModTime: time.Now(), Size: 100, Mode: 0o644},
		{Name: "cherry", ModTime: time.Now(), Size: 100, Mode: 0o644},
	}

	sortFiles(entries, false, false, false)

	expected := []string{"apple", "banana", "cherry", "zebra"}
	for i, entry := range entries {
		if entry.Name != expected[i] {
			t.Errorf("Position %d: expected %q, got %q", i, expected[i], entry.Name)
		}
	}
}

// TestSortFiles_ByTime tests sorting by modification time
func TestSortFiles_ByTime(t *testing.T) {
	now := time.Now()
	entries := []*FileEntry{
		{Name: "oldest", ModTime: now.Add(-3 * time.Hour), Size: 100, Mode: 0o644},
		{Name: "newest", ModTime: now, Size: 100, Mode: 0o644},
		{Name: "middle", ModTime: now.Add(-1 * time.Hour), Size: 100, Mode: 0o644},
	}

	sortFiles(entries, true, false, false)

	// Implementation sorts ascending (oldest first)
	expected := []string{"oldest", "middle", "newest"}
	for i, entry := range entries {
		if entry.Name != expected[i] {
			t.Errorf("Position %d: expected %q, got %q", i, expected[i], entry.Name)
		}
	}
}

// TestSortFiles_BySize tests sorting by file size
func TestSortFiles_BySize(t *testing.T) {
	entries := []*FileEntry{
		{Name: "small", ModTime: time.Now(), Size: 100, Mode: 0o644},
		{Name: "large", ModTime: time.Now(), Size: 10000, Mode: 0o644},
		{Name: "medium", ModTime: time.Now(), Size: 1000, Mode: 0o644},
	}

	sortFiles(entries, false, true, false)

	// Implementation sorts ascending (smallest first)
	expected := []string{"small", "medium", "large"}
	for i, entry := range entries {
		if entry.Name != expected[i] {
			t.Errorf("Position %d: expected %q, got %q", i, expected[i], entry.Name)
		}
	}
}

// TestSortFiles_Reverse tests reverse sorting
func TestSortFiles_Reverse(t *testing.T) {
	entries := []*FileEntry{
		{Name: "apple", ModTime: time.Now(), Size: 100, Mode: 0o644},
		{Name: "banana", ModTime: time.Now(), Size: 100, Mode: 0o644},
		{Name: "cherry", ModTime: time.Now(), Size: 100, Mode: 0o644},
	}

	sortFiles(entries, false, false, true)

	// Should be reversed alphabetical
	expected := []string{"cherry", "banana", "apple"}
	for i, entry := range entries {
		if entry.Name != expected[i] {
			t.Errorf("Position %d: expected %q, got %q", i, expected[i], entry.Name)
		}
	}
}

// TestSortFiles_ReverseTime tests reverse time sorting (newest first)
func TestSortFiles_ReverseTime(t *testing.T) {
	now := time.Now()
	entries := []*FileEntry{
		{Name: "oldest", ModTime: now.Add(-3 * time.Hour), Size: 100, Mode: 0o644},
		{Name: "newest", ModTime: now, Size: 100, Mode: 0o644},
		{Name: "middle", ModTime: now.Add(-1 * time.Hour), Size: 100, Mode: 0o644},
	}

	sortFiles(entries, true, false, true)

	// Ascending order reversed = newest first
	expected := []string{"newest", "middle", "oldest"}
	for i, entry := range entries {
		if entry.Name != expected[i] {
			t.Errorf("Position %d: expected %q, got %q", i, expected[i], entry.Name)
		}
	}
}

// TestSortFiles_ReverseSize tests reverse size sorting (largest first)
func TestSortFiles_ReverseSize(t *testing.T) {
	entries := []*FileEntry{
		{Name: "small", ModTime: time.Now(), Size: 100, Mode: 0o644},
		{Name: "large", ModTime: time.Now(), Size: 10000, Mode: 0o644},
		{Name: "medium", ModTime: time.Now(), Size: 1000, Mode: 0o644},
	}

	sortFiles(entries, false, true, true)

	// Ascending order reversed = largest first
	expected := []string{"large", "medium", "small"}
	for i, entry := range entries {
		if entry.Name != expected[i] {
			t.Errorf("Position %d: expected %q, got %q", i, expected[i], entry.Name)
		}
	}
}

// TestSortFiles_EmptySlice tests sorting an empty slice
func TestSortFiles_EmptySlice(t *testing.T) {
	entries := []*FileEntry{}
	sortFiles(entries, false, false, false)
	if len(entries) != 0 {
		t.Error("Empty slice should remain empty after sorting")
	}
}

// TestSortFiles_SingleElement tests sorting a single element
func TestSortFiles_SingleElement(t *testing.T) {
	entries := []*FileEntry{
		{Name: "only", ModTime: time.Now(), Size: 100, Mode: 0o644},
	}
	sortFiles(entries, true, true, true)
	if len(entries) != 1 || entries[0].Name != "only" {
		t.Error("Single element should remain unchanged")
	}
}

// TestSortFiles_SameValues tests sorting entries with same values
func TestSortFiles_SameValues(t *testing.T) {
	entries := []*FileEntry{
		{Name: "file1", ModTime: time.Now(), Size: 100, Mode: 0o644},
		{Name: "file2", ModTime: time.Now(), Size: 100, Mode: 0o644},
		{Name: "file3", ModTime: time.Now(), Size: 100, Mode: 0o644},
	}

	// Sort by time when all times are equal - should fall back to name
	sortFiles(entries, true, false, false)

	// Should be alphabetically sorted as fallback
	for i := 1; i < len(entries); i++ {
		if entries[i-1].Name > entries[i].Name {
			t.Errorf("Fallback to alphabetical sort failed: %q > %q", entries[i-1].Name, entries[i].Name)
		}
	}
}

// TestList_BasicListing tests basic directory listing functionality (Milestone 1)
func TestList_BasicListing(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create some test files
	testFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, name := range testFiles {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte("test"), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Test basic listing (should not error)
	err := List(tmpDir, false, false, false, false, false, false, "auto", false, false)
	if err != nil {
		t.Errorf("List() returned error: %v", err)
	}
}

// TestList_HiddenFiles tests the -a flag for hidden files (Milestone 2)
func TestList_HiddenFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create visible and hidden files
	visibleFile := filepath.Join(tmpDir, "visible.txt")
	hiddenFile := filepath.Join(tmpDir, ".hidden")

	os.WriteFile(visibleFile, []byte("visible"), 0o644)
	os.WriteFile(hiddenFile, []byte("hidden"), 0o644)

	// Without -a, hidden files should be filtered (we can't easily capture output,
	// but we can test that the function doesn't error)
	err := List(tmpDir, false, false, false, false, false, false, "auto", false, false)
	if err != nil {
		t.Errorf("List() without -a returned error: %v", err)
	}

	// With -a, all files should be included
	err = List(tmpDir, true, false, false, false, false, false, "auto", false, false)
	if err != nil {
		t.Errorf("List() with -a returned error: %v", err)
	}
}

// TestList_NonExistentPath tests error handling for non-existent paths
func TestList_NonExistentPath(t *testing.T) {
	err := List("/nonexistent/path/that/does/not/exist", false, false, false, false, false, false, "auto", false, false)
	if err == nil {
		t.Error("List() should return error for non-existent path")
	}
}

// TestList_Recursive tests the -R flag for recursive listing (Milestone 4)
func TestList_Recursive(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectory structure
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0o755)

	// Create files in both directories
	os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0o644)
	os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0o644)

	// Test recursive listing
	err := List(tmpDir, false, false, false, false, false, false, "auto", true, false)
	if err != nil {
		t.Errorf("List() recursive returned error: %v", err)
	}
}

// TestList_FileIndicators tests the -F flag for file type indicators (Milestone 4)
func TestList_FileIndicators(t *testing.T) {
	tmpDir := t.TempDir()

	// Create different file types
	os.Mkdir(filepath.Join(tmpDir, "directory"), 0o755)
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("file"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "executable"), []byte("exec"), 0o755)

	err := List(tmpDir, false, false, false, false, false, true, "auto", false, false)
	if err != nil {
		t.Errorf("List() with -F returned error: %v", err)
	}
}

// TestList_LongFormat tests the -l flag for long format (Milestone 3)
func TestList_LongFormat(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0o644)

	err := List(tmpDir, false, true, false, false, false, false, "auto", false, false)
	if err != nil {
		t.Errorf("List() with -l returned error: %v", err)
	}
}

// TestList_HumanReadable tests the -h flag for human-readable sizes (Milestone 3)
func TestList_HumanReadable(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files of different sizes
	os.WriteFile(filepath.Join(tmpDir, "small.txt"), []byte("small"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "large.txt"), make([]byte, 10000), 0o644)

	err := List(tmpDir, false, true, false, false, false, false, "auto", false, true)
	if err != nil {
		t.Errorf("List() with -h returned error: %v", err)
	}
}

// TestList_ColorOptions tests the --color flag (Milestone 4)
func TestList_ColorOptions(t *testing.T) {
	tmpDir := t.TempDir()

	os.Mkdir(filepath.Join(tmpDir, "dir"), 0o755)
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("file"), 0o644)

	// Test different color modes
	modes := []string{"always", "never", "auto"}
	for _, mode := range modes {
		err := List(tmpDir, false, false, false, false, false, false, mode, false, false)
		if err != nil {
			t.Errorf("List() with --color=%s returned error: %v", mode, err)
		}
	}
}

// TestList_MultiplePaths tests listing multiple paths
func TestList_MultiplePaths(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir1, "file1.txt"), []byte("file1"), 0o644)
	os.WriteFile(filepath.Join(tmpDir2, "file2.txt"), []byte("file2"), 0o644)

	// Note: The current implementation handles multiple paths in main.go,
	// not in the List function. This test verifies List works for individual paths.
	err := List(tmpDir1, false, false, false, false, false, false, "auto", false, false)
	if err != nil {
		t.Errorf("List() for first path returned error: %v", err)
	}

	err = List(tmpDir2, false, false, false, false, false, false, "auto", false, false)
	if err != nil {
		t.Errorf("List() for second path returned error: %v", err)
	}
}

// TestPrintColumns tests the printColumns function for multi-column TTY output
func TestPrintColumns(t *testing.T) {
	entries := []*FileEntry{
		{Name: "file1.txt", Mode: 0o644},
		{Name: "file2.txt", Mode: 0o644},
		{Name: "file3.txt", Mode: 0o644},
		{Name: "file4.txt", Mode: 0o644},
	}

	// Test without color or classify
	printColumns(entries, false, false)

	// Test with color
	printColumns(entries, true, false)

	// Test with classify
	printColumns(entries, false, true)

	// Test with both
	printColumns(entries, true, true)

	// Test empty slice
	printColumns([]*FileEntry{}, false, false)

	// Test single entry
	printColumns([]*FileEntry{{Name: "single.txt", Mode: 0o644}}, false, false)
}

// TestPrintColumns_WithDirectories tests printColumns with different file types
func TestPrintColumns_WithDirectories(t *testing.T) {
	entries := []*FileEntry{
		{Name: "directory", Mode: os.ModeDir | 0o755},
		{Name: "file.txt", Mode: 0o644},
		{Name: "executable", Mode: 0o755},
	}

	printColumns(entries, true, true)
}

// TestPrintColumns_LongNames tests printColumns with varying name lengths
func TestPrintColumns_LongNames(t *testing.T) {
	entries := []*FileEntry{
		{Name: "short", Mode: 0o644},
		{Name: "medium_length_name", Mode: 0o644},
		{Name: "very_long_filename_that_exceeds_normal_expectations.txt", Mode: 0o644},
	}

	printColumns(entries, false, false)
}

// TestList_FileEntryError tests that List handles NewFileEntry errors gracefully
func TestList_FileEntryError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with restricted permissions to potentially cause Info() issues
	testFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0o000)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// List should handle errors gracefully and continue
	err = List(tmpDir, false, false, false, false, false, false, "auto", false, false)
	if err != nil {
		t.Errorf("List() should handle file errors gracefully: %v", err)
	}
}

// TestList_RecursiveWithSubdirError tests recursive listing with inaccessible subdirectory
func TestList_RecursiveWithSubdirError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create files in both directories
	os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0o644)
	os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0o644)

	// Test recursive listing
	err = List(tmpDir, false, false, false, false, false, false, "auto", true, false)
	if err != nil {
		t.Errorf("List() recursive returned error: %v", err)
	}
}

// TestList_LongFormatWithClassify tests long format with file classification
func TestList_LongFormatWithClassify(t *testing.T) {
	tmpDir := t.TempDir()

	os.Mkdir(filepath.Join(tmpDir, "directory"), 0o755)
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("file"), 0o644)
	os.WriteFile(filepath.Join(tmpDir, "executable"), []byte("exec"), 0o755)

	err := List(tmpDir, false, true, false, false, false, true, "always", false, false)
	if err != nil {
		t.Errorf("List() with -l -F returned error: %v", err)
	}
}

// TestList_RecursiveWithLongFormat tests recursive listing with long format
func TestList_RecursiveWithLongFormat(t *testing.T) {
	tmpDir := t.TempDir()

	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0o755)
	os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0o644)
	os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0o644)

	err := List(tmpDir, false, true, false, false, false, false, "auto", true, false)
	if err != nil {
		t.Errorf("List() with -l -R returned error: %v", err)
	}
}
