package ls

import (
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	userCache  = make(map[string]string)
	groupCache = make(map[string]string)
	cacheMu    sync.RWMutex
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
	}

	// Extract Unix-specific metadata
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		entry.Nlink = uint64(stat.Nlink)
		entry.Blocks = int64(stat.Blocks)
		entry.Owner = resolveUser(strconv.Itoa(int(stat.Uid)))
		entry.Group = resolveGroup(strconv.Itoa(int(stat.Gid)))
	}

	return entry, nil
}

func resolveUser(uid string) string {
	cacheMu.RLock()
	if name, ok := userCache[uid]; ok {
		cacheMu.RUnlock()
		return name
	}
	cacheMu.RUnlock()

	cacheMu.Lock()
	defer cacheMu.Unlock()
	u, err := user.LookupId(uid)
	if err != nil {
		userCache[uid] = uid
		return uid
	}
	userCache[uid] = u.Username
	return u.Username
}

func resolveGroup(gid string) string {
	cacheMu.RLock()
	if name, ok := groupCache[gid]; ok {
		cacheMu.RUnlock()
		return name
	}
	cacheMu.RUnlock()

	cacheMu.Lock()
	defer cacheMu.Unlock()
	g, err := user.LookupGroupId(gid)
	if err != nil {
		groupCache[gid] = gid
		return gid
	}
	groupCache[gid] = g.Name
	return g.Name
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
