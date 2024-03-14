package platform

import "runtime"

type Platform interface {
	String() string
	platform()
}

type platform string

func (p platform) platform() {}

const (
	AIX       platform = "aix"
	Android   platform = "android"
	Darwin    platform = "darwin"
	Dragonfly platform = "dragonfly"
	FreeBSD   platform = "freebsd"
	Hurd      platform = "hurd"
	Illumos   platform = "illumos"
	IOS       platform = "ios"
	JS        platform = "js"
	Linux     platform = "linux"
	NACL      platform = "nacl"
	NetBSD    platform = "netbsd"
	OpenBSD   platform = "openbsd"
	Plan9     platform = "plan9"
	Solaris   platform = "solaris"
	Wasip1    platform = "wasip1"
	Windows   platform = "windows"
	ZOS       platform = "zos"
	Unknown   platform = ""
)

var lookup map[string]platform

func init() {
	all := []platform{
		AIX,
		Android,
		Darwin,
		Dragonfly,
		FreeBSD,
		Hurd,
		Illumos,
		IOS,
		JS,
		Linux,
		NACL,
		NetBSD,
		OpenBSD,
		Plan9,
		Solaris,
		Wasip1,
		Windows,
		ZOS}
	lookup = make(map[string]platform)
	for _, item := range all {
		lookup[item.String()] = item
	}
}

// String returns the string representation of the Platform
func (p platform) String() string {
	return string(p)
}

// IsPosix returns true if the platform is a unix platform
func IsPosix(p Platform) bool {
	switch p {
	case AIX:
		return true
	case Android:
		return true
	case Darwin:
		return true
	case Dragonfly:
		return true
	case FreeBSD:
		return true
	case Hurd:
		return true
	case Illumos:
		return true
	case IOS:
		return true
	case Linux:
		return true
	case NetBSD:
		return true
	case OpenBSD:
		return true
	case Solaris:
		return true
	}
	return false
}

// IsWindows returns true if the platform is windows
func IsWindows(p Platform) bool {
	return p == Windows
}

func Default() Platform {
	return Parse(runtime.GOOS)
}

func Parse(s string) Platform {
	p, ok := lookup[s]
	if !ok {
		return Unknown
	}
	return p
}
