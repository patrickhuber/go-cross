package os

import (
	"runtime"

	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
)

const (
	FakeWindowsWorkingDirectory = "c:\\working"
	FakeWindowsHomeDirectory    = "c:\\users\\fake"

	FakeUnixWorkingDirectory = "/working"
	FakeUnixHomeDirectory    = "/home/fake"
)

type memory struct {
	workingDirectory string
	platform         platform.Platform

	architecture  arch.Arch
	homeDirectory string
}

type MemoryOption func(*memory)

func WithHomeDirectory(homeDirectory string) MemoryOption {
	return func(o *memory) {
		o.homeDirectory = homeDirectory
	}
}

func WithArchitecture(architecture arch.Arch) MemoryOption {
	return func(o *memory) {
		o.architecture = architecture
	}
}

func WithWorkingDirectory(workingDirectory string) MemoryOption {
	return func(o *memory) {
		o.workingDirectory = workingDirectory
	}
}

func WithPlatform(platform platform.Platform) MemoryOption {
	return func(o *memory) {
		o.platform = platform
	}
}

// NewMemory creates a new OS from the mock OS request
func NewMemory(options ...MemoryOption) OS {
	o := &memory{}
	for _, option := range options {
		option(o)
	}
	if o.architecture == nil {
		o.architecture = arch.Parse(runtime.GOARCH)
	}
	if o.platform == nil {
		o.platform = platform.Parse(runtime.GOOS)
	}
	if o.workingDirectory == "" {
		if platform.IsWindows(o.platform) {
			o.workingDirectory = FakeWindowsWorkingDirectory
		} else {
			o.workingDirectory = FakeUnixWorkingDirectory
		}
	}
	if o.homeDirectory == "" {
		if platform.IsWindows(o.platform) {
			o.homeDirectory = FakeWindowsHomeDirectory
		} else {
			o.homeDirectory = FakeUnixHomeDirectory
		}
	}
	return o
}

func (o *memory) WorkingDirectory() (string, error) {
	return o.workingDirectory, nil
}

func (o *memory) Platform() platform.Platform {
	return platform.Platform(o.platform)
}

func (o *memory) Architecture() arch.Arch {
	return o.architecture
}

func (o *memory) Home() string {
	return o.homeDirectory
}

func (o *memory) ChangeDirectory(dir string) error {
	o.workingDirectory = dir
	return nil
}
