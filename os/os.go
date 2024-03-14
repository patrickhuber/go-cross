// Package os provides abstraction of the os package. Any os function that doesn't have a better fit lands in this package
package os

import (
	"os"
	"runtime"

	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
)

type OS interface {
	WorkingDirectory() (string, error)
	ChangeDirectory(dir string) error
	Platform() platform.Platform
	Architecture() arch.Arch
	Home() string
}

type realOS struct {
}

func New() OS {
	return &realOS{}
}

func (o *realOS) WorkingDirectory() (string, error) {
	return os.Getwd()
}

func (o *realOS) Executable() (string, error) {
	return os.Executable()
}

func (o *realOS) Platform() platform.Platform {
	return platform.Parse(runtime.GOOS)
}

func (o *realOS) Architecture() arch.Arch {
	return arch.Parse(runtime.GOARCH)
}

func (o *realOS) Home() string {
	dir, _ := os.UserHomeDir()
	return dir
}

func (o *realOS) ChangeDirectory(dir string) error {
	return os.Chdir(dir)
}
