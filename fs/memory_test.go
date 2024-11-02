package fs_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
)

func TestMemoryMkdirCreatesRootUnix(t *testing.T) {
	newConformance(platform.Linux).
		TestMkdirCreatesRoot(t, "/")
}

func TestMemoryMkdirFailsWhenRootNotExists(t *testing.T) {
	newConformance(platform.Linux).
		TestMkdirFailsWhenRootNotExists(t, "/test")
}

func TestMemoryMkdirAllCreatesAllDirectories(t *testing.T) {
	newConformance(platform.Linux).
		TestMkdirAllCreatesAllDirectories(t, "/gran/parent/child", []string{
			"/",
			"/gran",
			"/gran/parent",
			"/gran/parent/child",
		})
}

func TestMemoryWriteFile(t *testing.T) {
	newConformance(platform.Linux).
		TestWriteFile(t, "/gran/parent/child", "file.txt", "file")
}

func TestMemoryWriteCanGrow(t *testing.T) {
	newConformance(platform.Linux).
		TestWrite(t,
			"/gran/parent/child",
			"grow.txt",
			[]byte("this is test data"),
			7,
			[]byte(" more data than expected"),
			[]byte("this is more data than expected"))
}

func TestMemoryWriteCanOverwriteMiddle(t *testing.T) {
	newConformance(platform.Linux).
		TestWrite(t,
			"/gran/parent/child",
			"less.txt",
			[]byte("this is test data"),
			8,
			[]byte("also"),
			[]byte("this is also data"))
}

func TestMemoryWriteCanOverwriteEnd(t *testing.T) {
	newConformance(platform.Linux).
		TestWrite(t,
			"/gran/parent/child",
			"end.txt",
			[]byte("this is test data"),
			13,
			[]byte("info"),
			[]byte("this is test info"))

}

func TestMemoryReadDir(t *testing.T) {
	newConformance(platform.Linux).
		TestReadDir(t,
			"/gran/parent/child", []file{
				{"one.txt", []byte("one")},
				{"two.txt", []byte("two")},
				{"three.txt", []byte("three")},
				{"sub/one.txt", []byte("one")},
			}, []file{
				{"one.txt", []byte("one")},
				{"two.txt", []byte("two")},
				{"three.txt", []byte("three")},
				{"sub", []byte{}}})
}

func TestMemoryCanCreateFile(t *testing.T) {
	newConformance(platform.Linux).
		TestCanCreateFile(t, "/gran/parent/child", []file{
			{"one.txt", []byte("one")},
			{"two.txt", []byte("two")},
			{"three.txt", []byte("three")},
		})
}

func TestMemoryCanWriteFile(t *testing.T) {
	newConformance(platform.Linux).
		TestCanWriteFile(t, "/gran/parent/child", []file{
			{"one.txt", []byte("one")},
			{"two.txt", []byte("two")},
			{"three.txt", []byte("three")},
		})
}

func TestMemoryOpenFileFailsWhenReadOnlyAndNotExists(t *testing.T) {
	newConformance(platform.Linux).
		TestOpenFileFailsWhenNotExists(t, "/gran/parent/child", "/gran/parent/child/one.txt")

}

func TestWindowsWillNormalizePath(t *testing.T) {
	newConformance(platform.Windows).
		TestWindowsWillNormalizePath(t, `c:/ProgramData/fake/folder`, `test.txt`)
}

func TestWindowsFileExists(t *testing.T) {
	newConformance(platform.Windows).
		TestWindowsFileForwardAndBackwardSlash(t, "c:/ProgramData/fake/folder/test.txt")
}

func TestCanChmod(t *testing.T) {
	newConformance(platform.Windows).
		TestCanChangePermission(t, "/opt/fake/folder/test.txt")
}

func newMemory(o os.OS) (fs.FS, filepath.Provider) {
	path := filepath.NewProviderFromOS(o)
	fs := fs.NewMemory(path)
	return fs, path
}

func newOS(plat platform.Platform) os.OS {
	return os.NewMemory(os.WithPlatform(plat))
}

func newConformance(plat platform.Platform) *conformance {
	return NewConformanceWithProvider(newMemory(newOS(plat)))
}
