package fs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	fstest "testing/fstest"

	"github.com/patrickhuber/go-cross/filepath"
)

type memory struct {
	fs   fstest.MapFS
	path filepath.Provider
}

func NewMemory(path filepath.Provider) FS {
	m := &memory{
		fs:   fstest.MapFS{},
		path: path,
	}
	return m
}

func (m *memory) Create(name string) (File, error) {
	original := name
	name, err := m.path.Normalize(name)
	if err != nil {
		return nil, err
	}

	file, ok := m.fs[name]
	if !ok {
		file = &fstest.MapFile{}
		m.fs[name] = file
	}
	file.Data = nil
	file.Mode = 0666
	return &openFile{
		path: original,
		infoFile: infoFile{
			name: m.path.Base(original),
			file: file,
		},
	}, nil
}

// Open implements FS
func (m *memory) Open(name string) (fs.File, error) {
	op := "open"
	original := name
	name, err := m.path.Normalize(name)
	if err != nil {
		return nil, err
	}

	f, ok := m.fs[name]
	if !ok {
		return nil, &fs.PathError{
			Op:   op,
			Path: original,
			Err:  fs.ErrNotExist,
		}
	}
	return &openFile{
		path: name,
		infoFile: infoFile{
			name: m.path.Base(name),
			file: f,
		},
	}, nil
}

func isReadOnly(mode int) bool {
	switch {
	case mode&os.O_APPEND == os.O_APPEND:
		return false
	case mode&os.O_CREATE == os.O_CREATE:
		return false
	case mode&os.O_TRUNC == os.O_CREATE:
		return false
	case mode&os.O_WRONLY == os.O_WRONLY:
		return false
	case mode&os.O_RDWR == os.O_RDWR:
		return false
	}
	return true
}

// OpenFile implements OpenFS
func (m *memory) OpenFile(name string, mode int, perm fs.FileMode) (File, error) {
	op := "openFile"
	original := name

	name, err := m.path.Normalize(name)
	if err != nil {
		return nil, err
	}

	f, ok := m.fs[name]
	if !ok {
		// for readonly files, if the file doesn't exist return an error
		if isReadOnly(mode) {
			return nil, &fs.PathError{
				Op:   op,
				Path: original,
				Err:  fs.ErrNotExist,
			}
		}

		f = &fstest.MapFile{}
		m.fs[name] = f
	}

	// truncate if O_TRUNC specified
	if mode&os.O_TRUNC != 0 {
		f.Data = nil
	}

	// seek pos
	offset := 0
	if mode&os.O_APPEND != 0 {
		offset = len(f.Data)
	}

	return &openFile{
		path:   name,
		offset: int64(offset),
		infoFile: infoFile{
			name: m.path.Base(name),
			file: f,
		},
	}, nil
}

// Rename implements FS
func (m *memory) Rename(oldPath string, newPath string) error {

	var err error

	oldPath, err = m.path.Normalize(oldPath)
	if err != nil {
		return err
	}

	newPath, err = m.path.Normalize(newPath)
	if err != nil {
		return err
	}

	file, ok := m.fs[oldPath]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.fs, oldPath)
	m.fs[newPath] = file
	return nil
}

// Remove implements FS
func (m *memory) Remove(path string) error {
	path, err := m.path.Normalize(path)
	if err != nil {
		return err
	}
	_, ok := m.fs[path]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.fs, path)
	return nil
}

// RemoveAll implements FS
func (m *memory) RemoveAll(path string) error {
	paths := []string{}
	for p := range m.fs {
		if strings.HasPrefix(p, path) {
			paths = append(paths, p)
		}
	}
	for _, p := range paths {
		delete(m.fs, p)
	}
	return nil
}

// Glob implements FS
func (m *memory) Glob(pattern string) ([]string, error) {

	return m.fs.Glob(pattern)
}

// ReadDir implements FS
func (m *memory) ReadDir(name string) ([]fs.DirEntry, error) {
	// check that we can open the file
	d, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	// create the list of entries
	var entries []fs.DirEntry
	for path, file := range m.fs {
		originalPath := path

		path, err = m.path.Normalize(path)
		if err != nil {
			return nil, err
		}

		name, err = m.path.Normalize(name)
		if err != nil {
			return nil, err
		}

		// same dir
		if path == name {
			continue
		}

		// is the file's the directory the same as the
		if m.path.Dir(path) == name {

			// get the file name
			fileName := m.path.Base(originalPath)

			// append
			entries = append(entries, &infoFile{name: fileName, file: file})
		}
	}
	return entries, nil
}

// ReadFile implements FS
func (m *memory) ReadFile(name string) ([]byte, error) {

	f, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, stat.Size())
	if len(buf) == 0 {
		return buf, nil
	}

	_, err = f.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// WriteFile implements FS
func (m *memory) WriteFile(name string, data []byte, perm os.FileMode) error {
	name, err := m.path.Normalize(name)
	if err != nil {
		return err
	}

	file, ok := m.fs[name]
	if !ok {
		file = &fstest.MapFile{}
		m.fs[name] = file
	}

	file.Data = data
	file.Mode = perm

	return nil
}

// Exists implements FS
func (m *memory) Exists(path string) (bool, error) {
	name, err := m.path.Normalize(path)
	if err != nil {
		return false, err
	}
	_, ok := m.fs[name]
	return ok, nil
}

// Stat implements FS
func (m *memory) Stat(name string) (fs.FileInfo, error) {
	f, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

// Sub implements FS
func (m *memory) Sub(dir string) (fs.FS, error) {
	return m.fs.Sub(dir)
}

// Mkdir implements MakeDirFS
func (m *memory) Mkdir(path string, perm fs.FileMode) error {

	fp, err := m.path.Parse(path)
	if err != nil {
		return err
	}
	accumulator := fp.Root()

	// check each ancestor path
	for i := 0; i < len(fp.Segments); i++ {
		currentPath := m.path.String(accumulator)
		currentPath, err = m.path.Normalize(currentPath)
		if err != nil {
			return err
		}
		_, ok := m.fs[currentPath]
		if !ok {
			return errNotExist(currentPath)
		}
		seg := fp.Segments[i]
		fpseg, err := m.path.Parse(seg)
		if err != nil {
			return err
		}
		accumulator = accumulator.Join(fpseg)
	}

	// write the segment
	m.fs[path] = &fstest.MapFile{
		Mode: perm | fs.ModeDir,
	}

	return nil
}

// MkdirAll implements MakeDirFS
func (m *memory) MkdirAll(path string, perm fs.FileMode) error {

	// create all child paths of the current path from the root
	// so first, grab the root
	fp, err := m.path.Parse(path)
	if err != nil {
		return err
	}
	accumulator := fp.Root()

	// create each ancestor path
	for i := 0; i <= len(fp.Segments); i++ {
		currentPath := m.path.String(accumulator)
		currentPath, err = m.path.Normalize(currentPath)
		if err != nil {
			return err
		}
		_, ok := m.fs[currentPath]

		if !ok {
			m.fs[currentPath] = &fstest.MapFile{
				Mode: perm | fs.ModeDir,
			}
		}
		if i == len(fp.Segments) {
			break
		}

		seg := fp.Segments[i]
		fpseg, err := m.path.Parse(seg)
		if err != nil {
			return err
		}
		accumulator = accumulator.Join(fpseg)
	}
	return nil
}

func errNotExist(path string) error {
	return fmt.Errorf("'%s' %w", path, fs.ErrNotExist)
}
