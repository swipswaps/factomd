package entryfs

import "os"
import "log"
import "regexp"
import "strings"
import "github.com/hanwen/go-fuse/fuse"
import "github.com/hanwen/go-fuse/fuse/nodefs"
import "github.com/hanwen/go-fuse/fuse/pathfs"

type EntryFs struct {
	pathfs.FileSystem
	Dirs map[string][]string
	Sep  string
	*Storage
}

func (fs *EntryFs) GetAttr(name string, ctx *fuse.Context) (*fuse.Attr, fuse.Status) {
	if name == "" {
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	// ignore hidden files
	if string(name[0]) == "." {
		return nil, fuse.ENOENT
	}

	content, err := fs.Storage.FetchEntryData(name)

	if err == nil {
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644,
			Size: uint64(len(content)),
		}, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (fs *EntryFs) OpenDir(name string, ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	// REVIEW: may be nice to be able to treat directory blocks as directories of entries
	return nil, fuse.OK
}

func (fs *EntryFs) Open(name string, flags uint32, ctx *fuse.Context) (nodefs.File, fuse.Status) {
	key := fs.nameToKey(name)
	return NewEntryFile(fs.Storage, key), fuse.OK
}

func (fs *EntryFs) Create(name string, flags uint32, mode uint32, ctx *fuse.Context) (nodefs.File, fuse.Status) {
	key := fs.nameToKey(name)
	return NewEntryFile(fs.Storage, key), fuse.OK
}

func (fs *EntryFs) Rename(oldName string, newName string, ctx *fuse.Context) fuse.Status {
	return fuse.ENOENT
}

func (fs *EntryFs) Unlink(name string, ctx *fuse.Context) fuse.Status {
	return fuse.ENOENT
}

func (fs *EntryFs) Rmdir(name string, ctx *fuse.Context) fuse.Status {
	return fuse.ENOENT
}

func (fs *EntryFs) Mkdir(name string, mode uint32, ctx *fuse.Context) fuse.Status {
	return fuse.ENOENT
}

func (fs *EntryFs) nameToPattern(name string) string {
	pattern := fs.nameToKey(name)

	if name == "" {
		pattern += "*"
	} else {
		pattern += fs.Sep + "*"
	}

	return pattern
}

func (fs *EntryFs) dirsToEntries(dir string, m map[string]bool) []fuse.DirEntry {
	entries := make([]fuse.DirEntry, 0, 2)

	if dir == "" {
		dir = "."
	}

	if list, ok := fs.Dirs[dir]; ok {
		for _, key := range list {
			m[key] = true
			entries = append(entries, fuse.DirEntry{
				Name: key,
				Mode: fuse.S_IFDIR,
			})
		}
	}

	return entries
}

func (fs *EntryFs) resToEntries(dir string, list []string, m map[string]bool) []fuse.DirEntry {
	entries := make([]fuse.DirEntry, 0, 2)
	offset := len(dir)
	sepCount := strings.Count(dir, string(os.PathSeparator)) + 1

	if offset != 0 {
		offset += 1
	}

	for _, el := range list {
		key := el[offset:]
		keySepCount := strings.Count(key, fs.Sep)

		switch true {
		case keySepCount == 0:
			entries = append(entries, fuse.DirEntry{
				Name: fs.keyToName(key),
				Mode: fuse.S_IFREG,
			})
			break
		case keySepCount >= sepCount:
			tmp := strings.SplitN(key, fs.Sep, 2)
			key = tmp[0]

			if _, ok := m[key]; !ok {
				m[key] = true
				entries = append(entries, fuse.DirEntry{
					Name: key,
					Mode: fuse.S_IFDIR,
				})
			}
		}
	}

	return entries
}

func (fs *EntryFs) nameToKey(name string) string {
	re := regexp.MustCompile(string(os.PathSeparator))
	key := re.ReplaceAllLiteralString(name, fs.Sep)
	key = fs.decodePathSeparator(key)
	return key
}

func (fs *EntryFs) keyToName(key string) string {
	name := fs.encodePathSeparator(key)
	re := regexp.MustCompile(fs.Sep)
	name = re.ReplaceAllLiteralString(name, string(os.PathSeparator))
	return name
}

func (fs *EntryFs) encodePathSeparator(str string) string {
	re := regexp.MustCompile(string(os.PathSeparator))
	str = re.ReplaceAllLiteralString(str, "\uffff")
	return str
}

func (fs *EntryFs) decodePathSeparator(str string) string {
	re := regexp.MustCompile("\uffff")
	str = re.ReplaceAllLiteralString(str, string(os.PathSeparator))
	return str
}

func (fs *EntryFs) printError(err error) {
	log.Println("Error:", err)
}

func (fs *EntryFs) stringInSlice(target string, list []string) (bool, int) {
	for i, str := range list {
		if str == target {
			return true, i
		}
	}
	return false, -1
}
