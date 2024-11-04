package cross

import (
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
)

type Target interface {
	Platform() platform.Platform
	Architecture() arch.Arch
	OS() os.OS
	Path() filepath.Provider
	Env() env.Environment
	FS() fs.FS
	Console() console.Console
}

func NewTest(p platform.Platform, a arch.Arch, args ...string) Target {
	os := os.NewMemory(os.WithPlatform(p))
	path := filepath.NewProviderFromOS(os)
	return &target{
		os:           os,
		platform:     p,
		architecture: a,
		path:         path,
		env:          env.NewMemory(),
		fileSystem:   fs.NewMemory(path),
		console:      console.NewMemory(console.WithArgs(args)),
	}
}

func New() Target {
	os := os.New()
	return &target{
		os:           os,
		platform:     platform.Default(),
		architecture: arch.Default(),
		path:         filepath.NewProviderFromOS(os),
		env:          env.New(),
		fileSystem:   fs.New(),
		console:      console.New(),
	}
}

type target struct {
	platform     platform.Platform
	architecture arch.Arch
	path         filepath.Provider
	env          env.Environment
	fileSystem   fs.FS
	console      console.Console
	os           os.OS
}

func (t *target) Architecture() arch.Arch {
	return t.architecture
}

func (t *target) OS() os.OS {
	return t.os
}

func (t *target) Console() console.Console {
	return t.console
}

func (t *target) Platform() platform.Platform {
	return t.platform
}

func (t *target) Path() filepath.Provider {
	return t.path
}

func (t *target) Env() env.Environment {
	return t.env
}

func (t *target) FS() fs.FS {
	return t.fileSystem
}
