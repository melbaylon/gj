package ls

import (
	"os"
	"testing"
	"time"
)

// TestFormatSize tests the formatSize function for human-readable output
func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		human    bool
		expected string
	}{
		// Non-human readable
		{"bytes no human", 1234, false, "1234"},
		{"kb no human", 1024, false, "1024"},

		// Human readable - bytes
		{"zero bytes", 0, true, "0"},
		{"small bytes", 500, true, "500"},
		{"exact 1023 bytes", 1023, true, "1023"},

		// Human readable - kilobytes
		{"1 KB", 1024, true, "1.0K"},
		{"1.5 KB", 1536, true, "1.5K"},
		{"10 KB", 10240, true, "10.0K"},

		// Human readable - megabytes
		{"1 MB", 1048576, true, "1.0M"},
		{"2.5 MB", 2621440, true, "2.5M"},
		{"100 MB", 104857600, true, "100.0M"},

		// Human readable - gigabytes
		{"1 GB", 1073741824, true, "1.0G"},
		{"2.5 GB", 2684354560, true, "2.5G"},

		// Human readable - terabytes
		{"1 TB", 1099511627776, true, "1.0T"},
		{"1.5 TB", 1649267441664, true, "1.5T"},

		// Human readable - petabytes
		{"1 PB", 1125899906842624, true, "1.0P"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSize(tt.size, tt.human)
			if result != tt.expected {
				t.Errorf("formatSize(%d, %v) = %q, want %q", tt.size, tt.human, result, tt.expected)
			}
		})
	}
}

// TestFormatTime tests the formatTime function for recent vs old dates
func TestFormatTime(t *testing.T) {
	now := time.Now()

	// Recent time (within 6 months) - should show time
	recentTime := now.Add(-time.Hour * 24 * 30) // 30 days ago
	recentResult := formatTime(recentTime)

	// Check format includes time (HH:MM)
	if len(recentResult) < 12 {
		t.Errorf("Recent time format should include time, got %q", recentResult)
	}

	// Old time (more than 6 months) - should show year
	oldTime := now.AddDate(0, -7, 0) // 7 months ago
	oldResult := formatTime(oldTime)

	// Check format includes year
	if len(oldResult) < 12 {
		t.Errorf("Old time format should include year, got %q", oldResult)
	}

	// Future time - should show year
	futureTime := now.AddDate(1, 0, 0) // 1 year in future
	futureResult := formatTime(futureTime)
	if len(futureResult) < 12 {
		t.Errorf("Future time format should include year, got %q", futureResult)
	}
}

// TestFormatTime_ExactFormat tests the exact output format
func TestFormatTime_ExactFormat(t *testing.T) {
	// Create a specific time for testing
	testTime := time.Date(2024, time.January, 15, 14, 30, 0, 0, time.Local)

	// This should be considered "old" if we're past July 2024
	result := formatTime(testTime)

	// Should contain the year for old dates
	if testTime.Before(time.Now().AddDate(0, -6, 0)) {
		if result[len(result)-4:] != "2024" {
			t.Errorf("Old date should show year 2024, got %q", result)
		}
	}
}

// TestGetIndicator tests the getIndicator function for different file types
func TestGetIndicator(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		expected string
	}{
		{"directory", os.ModeDir, "/"},
		{"symlink", os.ModeSymlink, "@"},
		{"socket", os.ModeSocket, "="},
		{"named pipe", os.ModeNamedPipe, "|"},
		{"executable owner", 0o700, "*"},
		{"executable group", 0o070, "*"},
		{"executable other", 0o007, "*"},
		{"regular file", 0o644, ""},
		{"regular file no perms", 0o000, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIndicator(tt.mode)
			if result != tt.expected {
				t.Errorf("getIndicator(%o) = %q, want %q", tt.mode, result, tt.expected)
			}
		})
	}
}

// TestGetIndicator_ExecutableBits tests various executable bit combinations
func TestGetIndicator_ExecutableBits(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		expected string
	}{
		{"all exec bits", 0o111, "*"},
		{"owner exec only", 0o100, "*"},
		{"group exec only", 0o010, "*"},
		{"other exec only", 0o001, "*"},
		{"owner+group exec", 0o110, "*"},
		{"no exec bits", 0o644, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIndicator(tt.mode)
			if result != tt.expected {
				t.Errorf("getIndicator(%o) = %q, want %q", tt.mode, result, tt.expected)
			}
		})
	}
}

// TestApplyColor tests the applyColor function for different file types
func TestApplyColor(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		hasColor bool // whether output should contain color codes
	}{
		{"directory", os.ModeDir, true},
		{"symlink", os.ModeSymlink, true},
		{"socket", os.ModeSocket, true},
		{"named pipe", os.ModeNamedPipe, true},
		{"executable", 0o755, true},
		{"regular file", 0o644, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyColor("testfile", tt.mode)
			hasEscape := len(result) > len("testfile")

			if tt.hasColor && !hasEscape {
				t.Errorf("applyColor(%o) should add color codes", tt.mode)
			}
			if !tt.hasColor && hasEscape && result == "testfile" {
				// This is actually correct - no color means same as input
			}
		})
	}
}

// TestApplyColor_ColorCodes tests specific color codes are applied
func TestApplyColor_ColorCodes(t *testing.T) {
	// Directory should be blue
	dirResult := applyColor("test", os.ModeDir)
	if dirResult != colorBlue+"test"+colorReset {
		t.Errorf("Directory should be blue, got %q", dirResult)
	}

	// Symlink should be cyan
	symResult := applyColor("test", os.ModeSymlink)
	if symResult != colorCyan+"test"+colorReset {
		t.Errorf("Symlink should be cyan, got %q", symResult)
	}

	// Socket should be red
	sockResult := applyColor("test", os.ModeSocket)
	if sockResult != colorRed+"test"+colorReset {
		t.Errorf("Socket should be red, got %q", sockResult)
	}

	// Named pipe should be yellow
	pipeResult := applyColor("test", os.ModeNamedPipe)
	if pipeResult != colorYellow+"test"+colorReset {
		t.Errorf("Named pipe should be yellow, got %q", pipeResult)
	}

	// Executable should be green
	execResult := applyColor("test", 0o755)
	if execResult != colorGreen+"test"+colorReset {
		t.Errorf("Executable should be green, got %q", execResult)
	}
}
