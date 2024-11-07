package cross_test

import (
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	type test struct {
		platform platform.Platform
	}
	tests := []test{
		{platform: platform.Windows},
		{platform: platform.Linux},
		{platform: platform.Darwin},
	}
	for _, test := range tests {
		t.Run(test.platform.String(), func(t *testing.T) {
			a := arch.AMD64
			target := cross.NewTest(test.platform, a)
			require.Equal(t, a, target.Architecture())
			require.Equal(t, test.platform, target.Platform())
		})
	}
}

func TestInitializesDirectories(t *testing.T) {
	type test struct {
		platform platform.Platform
		paths    []string
	}
	tests := []test{
		{platform.Windows, []string{os.FakeWindowsHomeDirectory, os.FakeWindowsWorkingDirectory}},
		{platform.Linux, []string{os.FakeUnixHomeDirectory, os.FakeUnixWorkingDirectory}},
		{platform.Darwin, []string{os.FakeUnixHomeDirectory, os.FakeUnixWorkingDirectory}},
	}
	for _, test := range tests {
		t.Run(test.platform.String(), func(t *testing.T) {
			target := cross.NewTest(test.platform, arch.AMD64)
			fs := target.FS()
			for _, folder := range test.paths {
				ok, err := fs.Exists(folder)
				require.NoError(t, err)
				require.Truef(t, ok, "folder '%s' does not exist", folder)
			}
		})
	}
}
