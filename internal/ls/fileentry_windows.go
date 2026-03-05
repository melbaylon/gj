//go:build windows

package ls

import (
	"io/fs"
	"os"
)

// NewFileEntry creates a FileEntry from an os.DirEntry and its full path.
func NewFileEntry(path string, dirEntry fs.DirEntry) (*FileEntry, error) {
	info, err := dirEntry.Info()
	if err != nil {
		info, err = os.Lstat(path)
		if err != nil {
			return nil, err
		}
	}

	entry := &FileEntry{
		Name:    dirEntry.Name(),
		Path:    path,
		Mode:    info.Mode(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
		Owner:   "N/A", // Windows identity resolution is different
		Group:   "N/A",
		Nlink:   1,
	}

	return entry, nil
}
