package ls

import (
	"io/fs"
	"os"
	"time"
)

// FileEntry represents a file or directory entry with its relevant metadata.
type FileEntry struct {
	Name    string
	Path    string // Full path to the file/directory
	Mode    fs.FileMode
	Size    int64
	ModTime time.Time
	IsDir   bool
	Owner   string
	Group   string
	Nlink   uint64
	Blocks  int64
}

// FormatMode converts os.FileMode to a Unix-style string (e.g., drwxr-xr-x).
func FormatMode(mode fs.FileMode) string {
	res := make([]byte, 10)

	// File type
	if mode.IsDir() {
		res[0] = 'd'
	} else if mode&os.ModeSymlink != 0 {
		res[0] = 'l'
	} else if mode&os.ModeSocket != 0 {
		res[0] = 's'
	} else if mode&os.ModeNamedPipe != 0 {
		res[0] = 'p'
	} else if mode&os.ModeDevice != 0 {
		if mode&os.ModeCharDevice != 0 {
			res[0] = 'c'
		} else {
			res[0] = 'b'
		}
	} else {
		res[0] = '-'
	}

	// Owner
	res[1] = formatBit(mode&0400 != 0, 'r')
	res[2] = formatBit(mode&0200 != 0, 'w')
	res[3] = formatBit(mode&0100 != 0, 'x')

	// Group
	res[4] = formatBit(mode&0040 != 0, 'r')
	res[5] = formatBit(mode&0020 != 0, 'w')
	res[6] = formatBit(mode&0010 != 0, 'x')

	// Other
	res[7] = formatBit(mode&0004 != 0, 'r')
	res[8] = formatBit(mode&0002 != 0, 'w')
	res[9] = formatBit(mode&0001 != 0, 'x')

	return string(res)
}

func formatBit(condition bool, char byte) byte {
	if condition {
		return char
	}
	return '-'
}
