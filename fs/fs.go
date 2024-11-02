package fs

import (
	iofs "io/fs"
)

type OpenFileFS interface {
	OpenFile(name string, flag int, perm iofs.FileMode) (File, error)
}

type RenameFS interface {
	iofs.FS
	Rename(oldName, newName string) error
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
	Exists(name string) (bool, error)
}

type MakeDirFS interface {
	Mkdir(name string, perm iofs.FileMode) error
	MkdirAll(name string, perm iofs.FileMode) error
}

type CreateFS interface {
	Create(name string) (File, error)
}

type ChmodFS interface {
	Chmod(name string, mode iofs.FileMode) error
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
	ChmodFS
}
