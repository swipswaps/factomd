package entryfs

import (
	"path/filepath"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type Storage struct {
	DB   interfaces.DBOverlaySimple
	Path string
}

func (s *Storage) fetchEntry(key string) (interfaces.IEBEntry, error) {
	h, err := primitives.HexToHash(key)
	if err != nil {
		return nil, err
	}
	return s.DB.FetchEntry(h)
}

func (s *Storage) FetchEntryData(key string) ([]byte, error) {
	entry, err := s.fetchEntry(key)
	if err != nil {
		return nil, err
	}
	return entry.MarshalBinary()
}

func Mount(s *Storage) (*fuse.Server, error) {
	mnt, err := filepath.Abs(s.Path)

	if err != nil {
		return nil, err
	}

	fs := &EntryFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		Dirs:       make(map[string][]string),
		Sep:        ":", // KEY separator - REVIEW: a holdout from redis
		Storage:    s,
	}

	nfs := pathfs.NewPathNodeFs(fs, nil)
	server, _, err := nodefs.MountRoot(mnt, nfs.Root(), nil)

	if err != nil {
		return nil, err
	}

	return server, nil
}
