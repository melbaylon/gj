//go:build unix

package ls

import (
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"sync"
	"syscall"
)

var (
	userCache  = make(map[string]string)
	groupCache = make(map[string]string)
	cacheMu    sync.RWMutex
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
