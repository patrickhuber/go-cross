package arch

type Arch interface {
	arch()
	String() string
}

type arch string

func (arch) arch() {}

const (
	I386        arch = "386"
	AMD64       arch = "amd64"
	AMD64p32    arch = "amd64p32"
	ARM         arch = "arm"
	ARM64be     arch = "arm64be"
	Loong64     arch = "loong64"
	MIPS        arch = "mips"
	MIPSLE      arch = "mipsle"
	Mips64      arch = "mips64"
	Mips64le    arch = "mips64le"
	Mips64p32   arch = "mips64p32"
	Mips64p32le arch = "mips64p32le"
	PPC         arch = "ppc"
	PPC64le     arch = "ppc64le"
	RISCV       arch = "riscv"
	RISCV64     arch = "riscv64"
	S390        arch = "s390"
	S390x       arch = "s390x"
	SPARC       arch = "sparc"
	SPARC64     arch = "sparc64"
	WASM        arch = "wasm"
	Unknown     arch = ""
)

var lookup map[string]arch

func init() {
	all := []arch{
		I386,
		AMD64,
		AMD64p32,
		ARM,
		ARM64be,
		Loong64,
		MIPS,
		MIPSLE,
		Mips64,
		Mips64le,
		Mips64p32,
		Mips64p32le,
		PPC,
		PPC64le,
		RISCV,
		RISCV64,
		S390,
		S390x,
		SPARC,
		SPARC64,
		WASM,
	}
	lookup = make(map[string]arch)
	for _, a := range all {
		lookup[a.String()] = a
	}
}

func (a arch) String() string {
	return string(a)
}

func Parse(s string) Arch {
	a, ok := lookup[s]
	if !ok {
		return Unknown
	}
	return a
}
