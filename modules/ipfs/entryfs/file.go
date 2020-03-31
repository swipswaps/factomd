package entryfs

import "log"
import "time"
import "github.com/hanwen/go-fuse/fuse"
import "github.com/hanwen/go-fuse/fuse/nodefs"

type entryFile struct {
	*Storage
	key string
}

func NewEntryFile(s *Storage, key string) nodefs.File {
	file := new(entryFile)
	file.Storage = s
	file.key = key
	return file
}

func (f *entryFile) SetInode(*nodefs.Inode) {
}

func (f *entryFile) InnerFile() nodefs.File {
	return nil
}

func (f *entryFile) GetLk(uint64, *fuse.FileLock, uint32, *fuse.FileLock) fuse.Status {
	return fuse.OK
}

func (f *entryFile) SetLk(uint64, *fuse.FileLock, uint32) fuse.Status {
	return fuse.OK
}

func (f *entryFile) SetLkw(uint64, *fuse.FileLock, uint32) fuse.Status {
	return fuse.OK
}

func (f *entryFile) String() string {
	return "entryFile"
}

func (f *entryFile) Read(buf []byte, off int64) (fuse.ReadResult, fuse.Status) {
	data, err := f.Storage.FetchEntryData(f.key)

	if err != nil {
		log.Println("ERROR:", err)
		return nil, fuse.EIO
	}

	end := int(off) + int(len(buf))
	dataLen := len(data)

	if end > dataLen {
		end = dataLen
	}

	return fuse.ReadResultData(data[off:end]), fuse.OK
}

// no support for writes
func (f *entryFile) Write(data []byte, off int64) (uint32, fuse.Status) {
	return 0, fuse.EIO
}

func (f *entryFile) Flush() fuse.Status {
	return fuse.OK
}

func (f *entryFile) Release() {
}

func (f *entryFile) GetAttr(out *fuse.Attr) fuse.Status {
	content, err := f.Storage.FetchEntryData(f.key)

	if err != nil {
		log.Println("Error:", err)
		return fuse.EIO
	}

	out.Mode = fuse.S_IFREG | 0644
	out.Size = uint64(len(content))

	return fuse.OK
}

func (f *entryFile) Fsync(flags int) (code fuse.Status) {
	return fuse.OK
}

func (f *entryFile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	return fuse.ENOSYS
}

func (f *entryFile) Truncate(size uint64) fuse.Status {
	return fuse.OK
}

func (f *entryFile) Chown(uid uint32, gid uint32) fuse.Status {
	return fuse.ENOSYS
}

func (f *entryFile) Chmod(perms uint32) fuse.Status {
	return fuse.ENOSYS
}

func (f *entryFile) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) {
	return fuse.OK
}
