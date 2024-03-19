package filepath_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"
)

func TestCanParse(t *testing.T) {
	type test struct {
		path string
		fp   filepath.FilePath
	}
	windowsparse := []test{
		{path: "c:", fp: filepath.FilePath{
			Volume:   filepath.Volume{Drive: filepath.Nullable[string]{Value: `c:`, HasValue: true}},
			Absolute: false,
		}},
		{path: "c:/", fp: filepath.FilePath{
			Volume:   filepath.Volume{Drive: filepath.Nullable[string]{Value: `c:`, HasValue: true}},
			Segments: []string{""},
			Absolute: true,
		}},
		{path: "c:/foo", fp: filepath.FilePath{
			Volume:   filepath.Volume{Drive: filepath.Nullable[string]{Value: `c:`, HasValue: true}},
			Segments: []string{"foo"},
			Absolute: true,
		}},
		{path: "c:/foo/bar", fp: filepath.FilePath{
			Volume:   filepath.Volume{Drive: filepath.Nullable[string]{Value: `c:`, HasValue: true}},
			Segments: []string{"foo", "bar"},
			Absolute: true,
		}},
		{path: "//host/share", fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "host", HasValue: true},
				Share: filepath.Nullable[string]{Value: "share", HasValue: true},
			},
			Absolute: true,
		}},
		{path: "//host/share/", fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "host", HasValue: true},
				Share: filepath.Nullable[string]{Value: "share", HasValue: true},
			},
			Segments: []string{""},
			Absolute: true,
		}},
		{path: "//host/share/foo", fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "host", HasValue: true},
				Share: filepath.Nullable[string]{Value: "share", HasValue: true},
			},
			Segments: []string{"foo"},
			Absolute: true,
		}},
		{path: `\\host\share`, fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "host", HasValue: true},
				Share: filepath.Nullable[string]{Value: "share", HasValue: true},
			},
			Absolute: true,
		}},
		{path: `\\host\share\`, fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "host", HasValue: true},
				Share: filepath.Nullable[string]{Value: "share", HasValue: true},
			},
			Segments: []string{""}, // trailing slash
			Absolute: true,
		}},
		{path: `\\host\share\foo`, fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "host", HasValue: true},
				Share: filepath.Nullable[string]{Value: "share", HasValue: true},
			},
			Segments: []string{"foo"},
			Absolute: true,
		}},
		{path: `//./NUL`, fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: ".", HasValue: true},
				Share: filepath.Nullable[string]{Value: "NUL", HasValue: true},
			},
			Absolute: true,
		}},
		{path: `//?/NUL`, fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "?", HasValue: true},
				Share: filepath.Nullable[string]{Value: "NUL", HasValue: true},
			},
			Absolute: true,
		}},
		{path: "//abc//", fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "abc", HasValue: true},
				Share: filepath.Nullable[string]{Value: "", HasValue: true},
			},
			Segments: []string{""},
			Absolute: true,
		}},
		{path: "//abc", fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host: filepath.Nullable[string]{Value: "abc", HasValue: true},
			},
			Absolute: true,
		}},
		{path: "///abc", fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "", HasValue: true},
				Share: filepath.Nullable[string]{Value: "abc", HasValue: true},
			},
			Absolute: true,
		}},
		{path: "//abc//", fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "abc", HasValue: true},
				Share: filepath.Nullable[string]{Value: "", HasValue: true},
			},
			Segments: []string{""},
			Absolute: true,
		}},
		{path: "///abc/", fp: filepath.FilePath{
			Volume: filepath.Volume{
				Host:  filepath.Nullable[string]{Value: "", HasValue: true},
				Share: filepath.Nullable[string]{Value: "abc", HasValue: true},
			},
			Absolute: true,
			Segments: []string{""},
		}},
		{path: "a/b", fp: filepath.FilePath{
			Segments: []string{"a", "b"},
			Absolute: false,
		}},
		{path: "a/b/", fp: filepath.FilePath{
			Segments: []string{"a", "b", ""},
			Absolute: false,
		}},
		{path: "a/", fp: filepath.FilePath{
			Segments: []string{"a", ""},
			Absolute: false,
		}},
		{path: "a", fp: filepath.FilePath{
			Segments: []string{"a"},
			Absolute: false,
		}},
		{path: "", fp: filepath.FilePath{
			Absolute: false,
		}},
		// {path: `\a`, fp: filepath.FilePath{
		// 	Absolute: false,
		// 	Segments: []string{"a"},
		// }},
	}
	linuxparse := []test{
		{path: "/", fp: filepath.FilePath{
			Absolute: true,
		}},
	}
	run := func(tests []test, name string, plat platform.Platform) {
		for _, test := range tests {
			t.Run(name, func(t *testing.T) {
				parser := filepath.NewParserFromPlatform(plat)
				actual, err := parser.Parse(test.path)
				require.NoError(t, err)
				require.Equal(t, test.fp, actual, "unable to parse path '%s'", test.path)
			})
		}
	}
	run(windowsparse, "windowsparse", platform.Windows)
	run(linuxparse, "linuxparse", platform.Linux)

}
