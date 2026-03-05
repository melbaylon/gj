package ls

import (
	"os"
	"testing"
	"time"
)

// TestFormatMode tests the FormatMode function for various file types and permissions
func TestFormatMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		expected string
	}{
		// Regular files
		{"regular file no perms", 0o644, "-rw-r--r--"},
		{"regular file all perms", 0o777, "-rwxrwxrwx"},
		{"regular file no perms", 0o000, "----------"},
		{"regular file owner only", 0o700, "-rwx------"},

		// Directories
		{"directory", os.ModeDir | 0o755, "drwxr-xr-x"},
		{"directory full perms", os.ModeDir | 0o777, "drwxrwxrwx"},
		{"directory no perms", os.ModeDir | 0o000, "d---------"},

		// Symlinks
		{"symlink", os.ModeSymlink | 0o777, "lrwxrwxrwx"},
		{"symlink typical", os.ModeSymlink | 0o755, "lrwxr-xr-x"},

		// Sockets
		{"socket", os.ModeSocket | 0o755, "srwxr-xr-x"},

		// Named pipes (FIFO)
		{"named pipe", os.ModeNamedPipe | 0o644, "prw-r--r--"},

		// Character devices
		{"char device", os.ModeDevice | os.ModeCharDevice | 0o666, "crw-rw-rw-"},

		// Block devices
		{"block device", os.ModeDevice | 0o660, "brw-rw----"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatMode(tt.mode)
			if result != tt.expected {
				t.Errorf("FormatMode(%o) = %q, want %q", tt.mode, result, tt.expected)
			}
		})
	}
}

// TestFormatMode_FileTypeDetection tests that file type detection works correctly
func TestFormatMode_FileTypeDetection(t *testing.T) {
	// Ensure directory is detected first
	dirMode := os.ModeDir | 0o755
	result := FormatMode(dirMode)
	if result[0] != 'd' {
		t.Errorf("Directory mode should start with 'd', got %q", result[0])
	}

	// Ensure symlink is detected
	symMode := os.ModeSymlink | 0o777
	result = FormatMode(symMode)
	if result[0] != 'l' {
		t.Errorf("Symlink mode should start with 'l', got %q", result[0])
	}
}

// TestFormatMode_PermissionBits tests individual permission bits
func TestFormatMode_PermissionBits(t *testing.T) {
	tests := []struct {
		bit      os.FileMode
		position int
		char     byte
	}{
		{0o400, 1, 'r'}, // owner read
		{0o200, 2, 'w'}, // owner write
		{0o100, 3, 'x'}, // owner execute
		{0o040, 4, 'r'}, // group read
		{0o020, 5, 'w'}, // group write
		{0o010, 6, 'x'}, // group execute
		{0o004, 7, 'r'}, // other read
		{0o002, 8, 'w'}, // other write
		{0o001, 9, 'x'}, // other execute
	}

	for _, tt := range tests {
		t.Run(tt.char, func(t *testing.T) {
			mode := tt.bit
			result := FormatMode(mode)
			if len(result) != 10 {
				t.Fatalf("FormatMode result should be 10 chars, got %d", len(result))
			}
			if result[tt.position] != tt.char {
				t.Errorf("Position %d: expected %c, got %c (mode: %o)", tt.position, tt.char, result[tt.position], mode)
			}
		})
	}
}

// TestFormatMode_MissingPermissions tests that missing permissions show as '-'
func TestFormatMode_MissingPermissions(t *testing.T) {
	mode := os.FileMode(0o000)
	result := FormatMode(mode)

	// Check that all permission positions show '-'
	expectedPerms := "---------"
	actualPerms := result[1:]
	if actualPerms != expectedPerms {
		t.Errorf("Missing permissions should show as '-', got %q", actualPerms)
	}
}

// TestFileEntry_StructFields tests that FileEntry struct has all required fields
func TestFileEntry_StructFields(t *testing.T) {
	now := time.Now()
	entry := FileEntry{
		Name:    "test.txt",
		Path:    "/path/to/test.txt",
		Mode:    0o644,
		Size:    1024,
		ModTime: now,
		IsDir:   false,
		Owner:   "testuser",
		Group:   "testgroup",
		Nlink:   1,
		Blocks:  8,
	}

	if entry.Name != "test.txt" {
		t.Errorf("Name field mismatch")
	}
	if entry.Path != "/path/to/test.txt" {
		t.Errorf("Path field mismatch")
	}
	if entry.Size != 1024 {
		t.Errorf("Size field mismatch")
	}
	if entry.Owner != "testuser" {
		t.Errorf("Owner field mismatch")
	}
	if entry.Group != "testgroup" {
		t.Errorf("Group field mismatch")
	}
}
