package filepath_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	type test struct {
		fp       filepath.FilePath
		platform platform.Platform
		expected string
	}

	tests := []test{
		{
			// UNC share forward slash
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host:  filepath.Nullable[string]{Value: `host`, HasValue: true},
					Share: filepath.Nullable[string]{Value: `share`, HasValue: true},
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			platform: platform.Linux,
			expected: "//host/share/gran/parent/child",
		},
		{
			// UNC share backward slash
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host:  filepath.Nullable[string]{Value: `host`, HasValue: true},
					Share: filepath.Nullable[string]{Value: `share`, HasValue: true},
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			platform: platform.Windows,
			expected: `\\host\share\gran\parent\child`,
		},
		{
			// UNC share only root
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host:  filepath.Nullable[string]{Value: `host`, HasValue: true},
					Share: filepath.Nullable[string]{Value: `share`, HasValue: true},
				},
				Segments: []string{},
				Absolute: true,
			},
			platform: platform.Windows,
			expected: `\\host\share`,
		},
		{
			// UNC share without share name and trailing slash
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host:  filepath.Nullable[string]{Value: `abc`, HasValue: true},
					Share: filepath.Nullable[string]{Value: ``, HasValue: true},
				},
				Segments: []string{""},
				Absolute: true,
			},
			platform: platform.Windows,
			expected: `\\abc\\`,
		},
		{
			// UNC share without share name and trailing slash
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Host: filepath.Nullable[string]{Value: `abc`, HasValue: true},
				},
				Absolute: true,
			},
			platform: platform.Windows,
			expected: `\\abc`,
		},
		{
			// Unix Path
			fp: filepath.FilePath{
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			platform: platform.Linux,
			expected: "/gran/parent/child",
		},
		{
			// Windows Path
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Drive: filepath.Nullable[string]{Value: `c:`, HasValue: true},
				},
				Segments: []string{"gran", "parent", "child"},
				Absolute: true,
			},
			platform: platform.Windows,
			expected: `c:\gran\parent\child`,
		},
		{
			// relative
			fp: filepath.FilePath{
				Segments: []string{"gran", "parent", "child"},
				Absolute: false,
			},
			platform: platform.Linux,
			expected: "gran/parent/child",
		},
		{
			// root unix path
			fp: filepath.FilePath{
				Absolute: true,
			},
			platform: platform.Linux,
			expected: "/",
		},
		{
			// root windows path
			fp: filepath.FilePath{
				Volume: filepath.Volume{
					Drive: filepath.Nullable[string]{Value: `c:`, HasValue: true},
				},
				Absolute: true,
			},
			platform: platform.Windows,
			expected: `c:\`,
		},
	}

	for i, test := range tests {
		provider := filepath.NewProviderFromOS(os.NewMemory(os.WithPlatform(test.platform)))
		actual := test.fp.String(provider.Separator())
		require.Equal(t, test.expected, actual, "failed test at [%d]", i)
	}
}

func TestVolumeName(t *testing.T) {
	type test struct {
		path     string
		expected string
		platform platform.Platform
	}

	tests := []test{
		{
			"//host/share/gran/parent/child",
			`\\host\share`,
			platform.Windows,
		},
		{
			`\\host\share\gran\parent\child`,
			`\\host\share`,
			platform.Windows,
		},
		{
			"/gran/parent/child",
			"",
			platform.Linux,
		},
		{
			// Windows Path
			`c:\gran\parent\child`,
			`c:`,
			platform.Windows,
		},
	}

	for _, test := range tests {
		provider := filepath.NewProviderFromOS(os.NewMemory(os.WithPlatform(test.platform)))
		actual := provider.VolumeName(test.path)
		require.Equal(t, test.expected, actual)
	}
}
