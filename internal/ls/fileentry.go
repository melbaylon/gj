package ls

import (
	"io/fs"
	"time"
)

// FileEntry represents a file or directory entry with its relevant metadata.
// This struct will be enriched in later milestones for the long listing format.
type FileEntry struct {
	Name    string
	Path    string // Full path to the file/directory
	Mode    fs.FileMode
	Size    int64
	ModTime time.Time
	IsDir   bool
	// Add more fields here as needed for future milestones (e.g., Owner, Group, Nlink for -l)
}

// NewFileEntry creates a FileEntry from an os.DirEntry and its full path.
// It performs an os.Lstat to get detailed file info.
func NewFileEntry(path string, dirEntry fs.DirEntry) (*FileEntry, error) {
	info, err := dirEntry.Info()
	if err != nil {
		// If dirEntry.Info() returns an error, try os.Lstat directly on the path.
		// This can happen for broken symlinks where Info() might fail,
		// but Lstat would still give information about the link itself.
		info, err = fs.Lstat(path)
		if err != nil {
			return nil, err
		}
	}

	return &FileEntry{
		Name:    dirEntry.Name(),
		Path:    path,
		Mode:    info.Mode(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
	}, nil
}
