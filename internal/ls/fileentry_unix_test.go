//go:build unix

package ls

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// TestNewFileEntry tests the Unix-specific NewFileEntry function
func TestNewFileEntry(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a regular file
	filePath := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(filePath, []byte("test content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	dirEntry, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	entry, err := NewFileEntry(filePath, dirEntry[0])
	if err != nil {
		t.Fatalf("NewFileEntry() returned error: %v", err)
	}

	if entry.Name != "test.txt" {
		t.Errorf("Expected name 'test.txt', got %q", entry.Name)
	}
	if entry.Size != 12 {
		t.Errorf("Expected size 12, got %d", entry.Size)
	}
	if entry.IsDir {
		t.Error("Expected IsDir to be false")
	}
	// Owner and Group should be resolved (or fallback to numeric)
	if entry.Owner == "" {
		t.Error("Expected Owner to be set")
	}
	if entry.Group == "" {
		t.Error("Expected Group to be set")
	}
	if entry.Nlink == 0 {
		t.Error("Expected Nlink to be set")
	}
}

// TestNewFileEntry_Directory tests NewFileEntry with a directory
func TestNewFileEntry_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	allEntries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(allEntries) == 0 {
		t.Fatal("Expected at least one entry in directory")
	}

	for _, de := range allEntries {
		if de.Name() == "subdir" {
			entry, err := NewFileEntry(subDir, de)
			if err != nil {
				t.Fatalf("NewFileEntry() for directory returned error: %v", err)
			}
			if !entry.IsDir {
				t.Error("Expected directory entry to have IsDir=true")
			}
			return
		}
	}

	t.Error("Could not find subdir entry")
}

// TestNewFileEntry_Symlink tests NewFileEntry with a symlink
func TestNewFileEntry_Symlink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a target file
	targetPath := filepath.Join(tmpDir, "target.txt")
	err := os.WriteFile(targetPath, []byte("target"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create target file: %v", err)
	}

	// Create a symlink
	linkPath := filepath.Join(tmpDir, "link.txt")
	err = os.Symlink(targetPath, linkPath)
	if err != nil {
		t.Skipf("Cannot create symlink: %v (may not be supported on this filesystem)", err)
	}

	allEntries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	for _, de := range allEntries {
		if de.Name() == "link.txt" {
			entry, err := NewFileEntry(linkPath, de)
			if err != nil {
				t.Logf("NewFileEntry() for symlink returned error: %v", err)
				return // Symlinks may have issues in some environments
			}
			if entry.Mode&os.ModeSymlink == 0 {
				t.Error("Expected symlink to have ModeSymlink set")
			}
			break
		}
	}
}

// TestNewFileEntry_BrokenSymlink tests NewFileEntry with a broken symlink
func TestNewFileEntry_BrokenSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a broken symlink
	linkPath := filepath.Join(tmpDir, "broken_link.txt")
	err := os.Symlink("/nonexistent/target", linkPath)
	if err != nil {
		t.Skipf("Cannot create symlink: %v", err)
	}

	allEntries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	for _, de := range allEntries {
		if de.Name() == "broken_link.txt" {
			// This should handle the broken symlink gracefully
			entry, err := NewFileEntry(linkPath, de)
			if err != nil {
				// Broken symlinks may return an error, which is acceptable
				t.Logf("Broken symlink handled with error: %v", err)
				return
			}
			if entry != nil {
				t.Log("Broken symlink handled gracefully")
			}
			break
		}
	}
}

// TestResolveUser_Cache tests the user resolution caching
func TestResolveUser_Cache(t *testing.T) {
	// First call should populate cache
	uid := "0" // root user exists on most Unix systems
	user1 := resolveUser(uid)

	// Second call should use cache
	user2 := resolveUser(uid)

	if user1 != user2 {
		t.Errorf("Cached user resolution mismatch: %q != %q", user1, user2)
	}

	if user1 == "" {
		t.Error("Expected user resolution for uid 0")
	}
}

// TestResolveGroup_Cache tests the group resolution caching
func TestResolveGroup_Cache(t *testing.T) {
	// First call should populate cache
	gid := "0" // root group exists on most Unix systems
	group1 := resolveGroup(gid)

	// Second call should use cache
	group2 := resolveGroup(gid)

	if group1 != group2 {
		t.Errorf("Cached group resolution mismatch: %q != %q", group1, group2)
	}

	if group1 == "" {
		t.Error("Expected group resolution for gid 0")
	}
}

// TestResolveUser_InvalidUID tests resolution of invalid UID
func TestResolveUser_InvalidUID(t *testing.T) {
	// Use a very large UID that shouldn't exist
	invalidUID := "99999999"
	result := resolveUser(invalidUID)

	// Should fallback to numeric UID
	if result != invalidUID {
		t.Errorf("Expected fallback to numeric UID %q, got %q", invalidUID, result)
	}
}

// TestResolveGroup_InvalidGID tests resolution of invalid GID
func TestResolveGroup_InvalidGID(t *testing.T) {
	// Use a very large GID that shouldn't exist
	invalidGID := "99999999"
	result := resolveGroup(invalidGID)

	// Should fallback to numeric GID
	if result != invalidGID {
		t.Errorf("Expected fallback to numeric GID %q, got %q", invalidGID, result)
	}
}

// TestNewFileEntry_InfoError tests the error path when dirEntry.Info() fails
func TestNewFileEntry_InfoError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file
	filePath := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(filePath, []byte("test"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a mock DirEntry that will fail on Info()
	mockDirEntry := &mockDirEntry{
		name:    "test.txt",
		infoErr: os.ErrNotExist, // Force Info() to return error
	}

	// NewFileEntry should fall back to os.Lstat
	entry, err := NewFileEntry(filePath, mockDirEntry)
	if err != nil {
		t.Fatalf("NewFileEntry() with failing Info() returned error: %v", err)
	}

	if entry.Name != "test.txt" {
		t.Errorf("Expected name 'test.txt', got %q", entry.Name)
	}
}

// mockDirEntry implements fs.DirEntry for testing error paths
type mockDirEntry struct {
	name    string
	mode    fs.FileMode
	infoErr error
}

func (m *mockDirEntry) Name() string {
	return m.name
}

func (m *mockDirEntry) IsDir() bool {
	return false
}

func (m *mockDirEntry) Type() fs.FileMode {
	return m.mode
}

func (m *mockDirEntry) Info() (fs.FileInfo, error) {
	if m.infoErr != nil {
		return nil, m.infoErr
	}
	return nil, nil
}

// TestNewFileEntry_BothInfoAndLstatFail tests when both Info() and Lstat() fail
func TestNewFileEntry_BothInfoAndLstatFail(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mock DirEntry that will fail on Info()
	// and point to a non-existent path for Lstat() fallback
	mockDirEntry := &mockDirEntry{
		name:    "nonexistent.txt",
		infoErr: os.ErrNotExist,
	}

	nonExistentPath := filepath.Join(tmpDir, "nonexistent.txt")

	// This should return an error since both paths fail
	_, err := NewFileEntry(nonExistentPath, mockDirEntry)
	if err == nil {
		t.Error("Expected error when both Info() and Lstat() fail")
	}
}
