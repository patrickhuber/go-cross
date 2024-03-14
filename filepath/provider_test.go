package filepath_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"
)

func TestJoin(t *testing.T) {
	type test struct {
		elements []string
		sep      filepath.PathSeparator
		result   string
	}

	tests := []test{
		{
			[]string{"a", "b", "c"},
			filepath.ForwardSlash,
			"a/b/c",
		},
		{
			[]string{"a", "b/c"},
			filepath.ForwardSlash,
			"a/b/c",
		},
		{
			[]string{"a/b", "c"},
			filepath.ForwardSlash,
			"a/b/c",
		},
		{
			[]string{"a/b", "/c"},
			filepath.ForwardSlash,
			"a/b/c",
		},
		{
			[]string{"/a/b", "/c"},
			filepath.ForwardSlash,
			"/a/b/c",
		},
		{
			[]string{`c:\`, `a\b`, `c`},
			filepath.BackwardSlash,
			`c:\a\b\c`,
		},
	}

	for _, test := range tests {
		provider := filepath.NewProvider(filepath.WithSeparator(test.sep))
		actual := provider.Join(test.elements...)
		require.Equal(t, test.result, actual)
	}
}

func TestRoot(t *testing.T) {
	type test struct {
		path     string
		sep      filepath.PathSeparator
		expected string
	}

	tests := []test{
		{
			// UNC forward slash
			"//host/share/gran/parent/child",
			filepath.ForwardSlash,
			"//host/share",
		},
		{
			// UNC backward slash
			`\\host\share\gran\parent\child`,
			filepath.BackwardSlash,
			`\\host\share`,
		},
		{
			// Unix Path
			"/gran/parent/child",
			filepath.ForwardSlash,
			"/",
		},
		{
			// Windows Path
			`c:\gran\parent\child`,
			filepath.BackwardSlash,
			`c:\`,
		},
	}

	for _, test := range tests {
		provider := filepath.NewProvider(filepath.WithSeparator(test.sep))
		actual := provider.Root(test.path)
		require.Equal(t, test.expected, actual)
	}
}

func TestRel(t *testing.T) {
	type test struct {
		source   string
		target   string
		expected string
	}
	reltests := []test{
		{"a/b", "a/b", "."},
		{"a/b/.", "a/b", "."},
		{"a/b", "a/b/.", "."},
		{"./a/b", "a/b", "."},
		{"a/b", "./a/b", "."},
		{"ab/cd", "ab/cde", "../cde"},
		{"ab/cd", "ab/c", "../c"},
		{"a/b", "a/b/c/d", "c/d"},
		{"a/b", "a/b/../c", "../c"},
		{"a/b/../c", "a/b", "../b"},
		{"a/b/c", "a/c/d", "../../c/d"},
		{"a/b", "c/d", "../../c/d"},
		{"a/b/c/d", "a/b", "../.."},
		{"a/b/c/d", "a/b/", "../.."},
		{"a/b/c/d/", "a/b", "../.."},
		{"a/b/c/d/", "a/b/", "../.."},
		{"../../a/b", "../../a/b/c/d", "c/d"},
		{"/a/b", "/a/b", "."},
		{"/a/b/.", "/a/b", "."},
		{"/a/b", "/a/b/.", "."},
		{"/ab/cd", "/ab/cde", "../cde"},
		{"/ab/cd", "/ab/c", "../c"},
		{"/a/b", "/a/b/c/d", "c/d"},
		{"/a/b", "/a/b/../c", "../c"},
		{"/a/b/../c", "/a/b", "../b"},
		{"/a/b/c", "/a/c/d", "../../c/d"},
		{"/a/b", "/c/d", "../../c/d"},
		{"/a/b/c/d", "/a/b", "../.."},
		{"/a/b/c/d", "/a/b/", "../.."},
		{"/a/b/c/d/", "/a/b", "../.."},
		{"/a/b/c/d/", "/a/b/", "../.."},
		{"/../../a/b", "/../../a/b/c/d", "c/d"},
		{".", "a/b", "a/b"},
		{".", "..", ".."},

		// can't do purely lexically
		{"..", ".", "err"},
		{"..", "a", "err"},
		{"../..", "..", "err"},
		{"a", "/a", "err"},
		{"/a", "a", "err"},
	}

	winreltests := []test{
		{`C:a\b\c`, `C:a/b/d`, `..\d`},
		{`C:\`, `D:\`, `err`},
		{`C:`, `D:`, `err`},
		{`C:\Projects`, `c:\projects\src`, `src`},
		{`C:\Projects`, `c:\projects`, `.`},
		{`C:\Projects\a\..`, `c:\projects`, `.`},
		{`\\host\share`, `\\host\share\file.txt`, `file.txt`},
		{`\\host\share\folder`, `\\other\test\share`, `err`},
	}

	run := func(tests []test, name string, o os.OS) {
		provider := filepath.NewProvider(filepath.WithOS(o))
		for i, test := range tests {

			actual, err := provider.Rel(test.source, test.target)
			if err != nil {
				actual = "err"
			}
			require.Equal(t, test.expected, actual,
				"test %s[%d] failed. source '%s' target '%s' expected '%s' actual '%s'",
				name, i, test.source, test.target, test.expected, actual)
		}
	}
	run(reltests, "reltests", os.NewMemory(os.WithPlatform(platform.Linux)))
	run(winreltests, "winreltests", os.NewMemory(os.WithPlatform(platform.Windows)))
}

func TestClean(t *testing.T) {
	type test struct {
		path     string
		expected string
	}

	cleantests := []test{
		// Already clean
		{"abc", "abc"},
		{"abc/def", "abc/def"},
		{"a/b/c", "a/b/c"},
		{".", "."},
		{"..", ".."},
		{"../..", "../.."},
		{"../../abc", "../../abc"},
		{"/abc", "/abc"},
		{"/", "/"},

		// Empty is current dir
		{"", "."},

		// Remove trailing slash
		{"abc/", "abc"},
		{"abc/def/", "abc/def"},
		{"a/b/c/", "a/b/c"},
		{"./", "."},
		{"../", ".."},
		{"../../", "../.."},
		{"/abc/", "/abc"},

		// Remove doubled slash
		{"abc//def//ghi", "abc/def/ghi"},
		{"abc//", "abc"},

		// Remove . elements
		{"abc/./def", "abc/def"},
		{"/./abc/def", "/abc/def"},
		{"abc/.", "abc"},

		// Remove .. elements
		{"abc/def/ghi/../jkl", "abc/def/jkl"},
		{"abc/def/../ghi/../jkl", "abc/jkl"},
		{"abc/def/..", "abc"},
		{"abc/def/../..", "."},
		{"/abc/def/../..", "/"},
		{"abc/def/../../..", ".."},
		{"/abc/def/../../..", "/"},
		{"abc/def/../../../ghi/jkl/../../../mno", "../../mno"},
		{"/../abc", "/abc"},

		// Combinations
		{"abc/./../def", "def"},
		{"abc//./../def", "def"},
		{"abc/../../././../def", "../../def"},
	}

	nonwincleantests := []test{
		// Remove leading doubled slash
		{"//abc", "/abc"},
		{"///abc", "/abc"},
		{"//abc//", "/abc"},
	}

	wincleantests := []test{
		{`c:`, `c:.`},
		{`c:\`, `c:\`},
		{`c:\abc`, `c:\abc`},
		{`c:abc\..\..\.\.\..\def`, `c:..\..\def`},
		{`c:\abc\def\..\..`, `c:\`},
		{`c:\..\abc`, `c:\abc`},
		{`c:..\abc`, `c:..\abc`},
		{`\`, `\`},
		{`/`, `\`},
		{`\\i\..\c$`, `\\i\..\c$`},
		{`\\i\..\i\c$`, `\\i\..\i\c$`},
		{`\\i\..\I\c$`, `\\i\..\I\c$`},
		{`\\host\share\foo\..\bar`, `\\host\share\bar`},
		{`//host/share/foo/../baz`, `\\host\share\baz`},
		{`\\host\share\foo\..\..\..\..\bar`, `\\host\share\bar`},
		{`\\.\C:\a\..\..\..\..\bar`, `\\.\C:\bar`},
		{`\\.\C:\\\\a`, `\\.\C:\a`},
		{`\\a\b\..\c`, `\\a\b\c`},
		{`\\a\b`, `\\a\b`},
		{`.\c:`, `.\c:`},
		{`.\c:\foo`, `.\c:\foo`},
		{`.\c:foo`, `.\c:foo`},
		{`//abc`, `\\abc`},
		{`//abc/`, `\\abc\`},
		{`///abc`, `\\\abc`},
		{`//abc//`, `\\abc\\`},
		{`///abc/`, `\\\abc\`},

		// Don't allow cleaning to move an element with a colon to the start of the path.
		{`a/../c:`, `.\c:`},
		{`a\..\c:`, `.\c:`},
		{`a/../c:/a`, `.\c:\a`},
		{`a/../../c:`, `..\c:`},
		{`foo:bar`, `foo:bar`},
	}

	run := func(tests []test, name string, o os.OS) {
		provider := filepath.NewProvider(filepath.WithOS(o))
		for i, test := range tests {
			actual := provider.Clean(test.path)
			require.Equal(t, test.expected, actual,
				"%s test failed: %d given '%s' expected '%s' actual '%s'", name, i, test.path, test.expected, actual)
		}
	}

	run(cleantests, "cleantests", os.NewMemory(os.WithPlatform(platform.Linux)))
	run(nonwincleantests, "nonwincleantests", os.NewMemory(os.WithPlatform(platform.Linux)))
	run(wincleantests, "wincleantests", os.NewMemory(os.WithPlatform(platform.Windows)))
}

func TestDir(t *testing.T) {

	type test struct {
		path     string
		expected string
	}

	var dirtests = []test{
		{"", "."},
		{".", "."},
		{"/.", "/"},
		{"/", "/"},
		{"/foo", "/"},
		{"x/", "x"},
		{"abc", "."},
		{"abc/def", "abc"},
		{"a/b/.x", "a/b"},
		{"a/b/c.", "a/b"},
		{"a/b/c.x", "a/b"},
	}

	var nonwindirtests = []test{
		{"////", "/"},
	}

	var windirtests = []test{
		{`c:\`, `c:\`},
		{`c:.`, `c:.`},
		{`c:\a\b`, `c:\a`},
		{`c:a\b`, `c:a`},
		{`c:a\b\c`, `c:a\b`},
		{`\\host\share`, `\\host\share`},
		{`\\host\share\`, `\\host\share\`},
		{`\\host\share\a`, `\\host\share\`},
		{`\\host\share\a\b`, `\\host\share\a`},
		{`\\\\`, `\\\\`},
	}

	run := func(tests []test, name string, o os.OS) {
		provider := filepath.NewProvider(filepath.WithOS(o))
		for i, test := range tests {
			actual := provider.Dir(test.path)
			require.Equal(t, test.expected, actual,
				"%s[%d] given '%s' expected '%s' actual '%s'",
				name, i, test.path, test.expected, actual)
		}
	}

	run(dirtests, "dirtests", os.NewMemory(os.WithPlatform(platform.Linux)))
	run(nonwindirtests, "nonwindirtests", os.NewMemory(os.WithPlatform(platform.Linux)))
	run(windirtests, "windirtests", os.NewMemory(os.WithPlatform(platform.Windows)))
}

func TestExt(t *testing.T) {
	type test struct {
		path string
		ext  string
	}
	var exttests = []test{
		{"path.go", ".go"},
		{"path.pb.go", ".go"},
		{"a.dir/b", ""},
		{"a.dir/b.go", ".go"},
		{"a.dir/", ""},
	}

	provider := filepath.NewProvider()
	for _, test := range exttests {
		actual := provider.Ext(test.path)
		require.Equal(t, test.ext, actual)
	}
}

func TestBase(t *testing.T) {
	type test struct {
		path     string
		expected string
	}
	var basetests = []test{
		{"", "."},
		{".", "."},
		{"/.", "."},
		{"/", "/"},
		{"////", "/"},
		{"x/", "x"},
		{"abc", "abc"},
		{"abc/def", "def"},
		{"a/b/.x", ".x"},
		{"a/b/c.", "c."},
		{"a/b/c.x", "c.x"},
	}

	var winbasetests = []test{
		{`c:\`, `\`},
		{`c:.`, `.`},
		{`c:\a\b`, `b`},
		{`c:a\b`, `b`},
		{`c:a\b\c`, `c`},
		{`\\host\share\`, `\`},
		{`\\host\share\a`, `a`},
		{`\\host\share\a\b`, `b`},
	}
	run := func(tests []test, name string, o os.OS) {
		provider := filepath.NewProvider(filepath.WithOS(o))
		for i, test := range tests {
			actual := provider.Base(test.path)
			require.Equal(t, test.expected, actual,
				"%s[%d] given: '%s' expected: '%s' actual: '%s'", name, i, test.path, test.expected, actual)
		}
	}
	run(basetests, "basetests", os.NewMemory(os.WithPlatform(platform.Linux)))
	run(winbasetests, "winbasetests", os.NewMemory(os.WithPlatform(platform.Windows)))
}

func TestAbs(t *testing.T) {
	var absDirs = []string{
		"a",
		"a/b",
		"a/b/c",
	}
	var relPaths = []string{
		".",
		"b",
		"b/",
		"../a",
		"../a/b",
		"../a/b/./c/../../.././a",
		"../a/b/./c/../../.././a/",
		//"$",
		//"$/.",
		//"$/a/../a/b",
		//"$/a/b/c/../../.././a",
		//"$/a/b/c/../../.././a/",
	}
	run := func(name string, abs, rel []string, o os.OS) {
		t.Run(name, func(t *testing.T) {
			for _, a := range abs {
				o.ChangeDirectory(a)
				for _, r := range rel {
					path := filepath.NewProvider(filepath.WithOS(o))
					result, err := path.Abs(r)
					require.NoError(t, err)
					require.NotEmpty(t, result)
				}
			}
		})
	}
	run("abs_linux", absDirs, relPaths, os.NewMemory(os.WithPlatform(platform.Linux)))
	run("abs_darwin", absDirs, relPaths, os.NewMemory(os.WithPlatform(platform.Darwin)))
	run("abs_windows", absDirs, relPaths, os.NewMemory(os.WithPlatform(platform.Windows)))
}
