package console

import (
	"bytes"
	"io"
)

const (
	MockWindowsExecutable = "c:\\ProgramData\\test\\fake.exe"
	MockLinuxExecutable   = "/opt/test/fake"
	MockDarwinExecutable  = MockLinuxExecutable
)

type memory struct {
	in         *bytes.Buffer
	out        *bytes.Buffer
	err        *bytes.Buffer
	args       []string
	executable string
}

type Memory interface {
	Console
	// OutBuffer exposes the output buffer for the memory console to enable testing
	OutBuffer() *bytes.Buffer
	// ErrBuffer exposes the error buffer for the memory console to enable testing
	ErrBuffer() *bytes.Buffer
	// InBuffer exposes the error buffer for the memory console to enable testing
	InBuffer() *bytes.Buffer
}

type MemoryOption func(*memory)

func WithExecutable(executable string) MemoryOption {
	return func(m *memory) {
		m.executable = executable
	}
}

func WithArgs(args []string) MemoryOption {
	return func(m *memory) {
		m.args = args
	}
}

func NewMemory(options ...MemoryOption) Memory {
	m := &memory{
		in:         &bytes.Buffer{},
		out:        &bytes.Buffer{},
		err:        &bytes.Buffer{},
		args:       []string{},
		executable: "string",
	}
	for _, option := range options {
		option(m)
	}
	return m
}

func (c *memory) In() io.Reader {
	return c.in
}

func (c *memory) Out() io.Writer {
	return c.out
}

func (c *memory) Error() io.Writer {
	return c.err
}

// ErrBuffer implements MemoryConsole
func (c *memory) ErrBuffer() *bytes.Buffer {
	return c.err
}

// InBuffer implements MemoryConsole
func (c *memory) InBuffer() *bytes.Buffer {
	return c.in
}

// OutBuffer implements MemoryConsole
func (c *memory) OutBuffer() *bytes.Buffer {
	return c.out
}

func (c *memory) Args() []string {
	return c.args
}

func (c *memory) Executable() (string, error) {
	return c.executable, nil
}
