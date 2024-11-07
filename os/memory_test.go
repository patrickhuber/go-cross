package os_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"
)

func TestPlatform(t *testing.T) {
	type test struct {
		name     string
		expected platform.Platform
		o        os.OS
	}
	platforms := []platform.Platform{
		platform.Darwin,
		platform.Linux,
		platform.Windows,
		platform.Plan9,
		platform.FreeBSD,
		platform.AIX,
		platform.Android,
		platform.Dragonfly,
		platform.Hurd,
		platform.IOS,
		platform.Illumos,
		platform.JS,
		platform.Wasip1,
		platform.NACL,
		platform.ZOS,
		platform.Solaris,
	}
	var tests []test
	for _, p := range platforms {
		tests = append(tests, test{
			name:     p.String(),
			expected: p,
			o:        os.NewMemory(os.WithPlatform(p)),
		})
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, test.o.Platform())
		})
	}
}

func TestArchitecture(t *testing.T) {
	type test struct {
		name     string
		expected arch.Arch
		o        os.OS
	}
	architectures := []arch.Arch{
		arch.AMD64,
		arch.AMD64p32,
		arch.ARM,
		arch.ARM64be,
		arch.I386,
		arch.Loong64,
		arch.MIPS,
		arch.Mips64,
		arch.Mips64le,
		arch.Mips64p32,
		arch.Mips64p32le,
		arch.PPC,
		arch.PPC64le,
		arch.RISCV,
		arch.RISCV64,
		arch.S390,
		arch.S390x,
		arch.SPARC,
		arch.SPARC64,
		arch.WASM,
	}
	var tests []test
	for _, a := range architectures {
		tests = append(tests, test{
			name:     a.String(),
			expected: a,
			o:        os.NewMemory(os.WithArchitecture(a)),
		})
	}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, test.o.Architecture(), "test [%d] failed", i)
		})
	}
}

func TestHome(t *testing.T) {
	type test struct {
		expected string
		o        os.OS
	}
	const (
		OtherHome = "/home/other"
	)
	tests := []test{
		{expected: os.FakeUnixHomeDirectory, o: os.NewMemory(os.WithPlatform(platform.Darwin))},
		{expected: os.FakeUnixHomeDirectory, o: os.NewMemory(os.WithPlatform(platform.Linux))},
		{expected: os.FakeWindowsHomeDirectory, o: os.NewMemory(os.WithPlatform(platform.Windows))},
		{expected: OtherHome, o: os.NewMemory(os.WithHomeDirectory(OtherHome))},
	}
	for i, test := range tests {
		home, err := test.o.Home()
		require.NoError(t, err)
		require.Equal(t, test.expected, home, "test [%d] failed", i)
	}
}

func TestWorkingDirectory(t *testing.T) {
	type test struct {
		expected string
		o        os.OS
	}
	const (
		OtherWorkingDirectory = "/home/other/wd"
	)
	tests := []test{
		{expected: os.FakeUnixWorkingDirectory, o: os.NewMemory(os.WithPlatform(platform.Darwin))},
		{expected: os.FakeUnixWorkingDirectory, o: os.NewMemory(os.WithPlatform(platform.Linux))},
		{expected: os.FakeWindowsWorkingDirectory, o: os.NewMemory(os.WithPlatform(platform.Windows))},
		{expected: OtherWorkingDirectory, o: os.NewMemory(os.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: os.NewMemory(os.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: os.NewMemory(os.WithWorkingDirectory(OtherWorkingDirectory))},
		{expected: OtherWorkingDirectory, o: os.NewMemory(os.WithWorkingDirectory(OtherWorkingDirectory))},
	}
	for i, test := range tests {
		workingDirectory, err := test.o.WorkingDirectory()
		require.NoError(t, err, "test [%d] o.WorkingDirectory() returned error")
		require.Equal(t, test.expected, workingDirectory, "test [%d] failed", i)
	}
}
