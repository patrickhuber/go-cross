package fs

import (
	iofs "io/fs"
)

type OpenFileFS interface {
	OpenFile(name string, flag int, perm iofs.FileMode) (File, error)
}

type RenameFS interface {
	iofs.FS
	Rename(oldPath, newPath string) error
}

type RemoveFS interface {
	iofs.FS
	Remove(name string) error
	RemoveAll(name string) error
}

type WriteFileFS interface {
	WriteFile(name string, data []byte, perm iofs.FileMode) error
}

type ExistsFS interface {
	Exists(path string) (bool, error)
}

type MakeDirFS interface {
	Mkdir(path string, perm iofs.FileMode) error
	MkdirAll(path string, perm iofs.FileMode) error
}

type CreateFS interface {
	Create(path string) (File, error)
}

type FS interface {
	iofs.FS
	OpenFileFS
	CreateFS
	RenameFS
	RemoveFS
	WriteFileFS
	ExistsFS
	iofs.GlobFS
	iofs.ReadFileFS
	iofs.ReadFileFS
	iofs.StatFS
	iofs.SubFS
	iofs.ReadDirFS
	MakeDirFS
}
