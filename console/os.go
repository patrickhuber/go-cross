package console

import (
	"io"
	"os"
)

type console struct {
}

func New() Console {
	return &console{}
}

func (c *console) In() io.Reader {
	return os.Stdin
}

func (c *console) Out() io.Writer {
	return os.Stdout
}

func (c *console) Error() io.Writer {
	return os.Stderr
}

func (c *console) Args() []string {
	return os.Args
}

func (c *console) Executable() (string, error) {
	return os.Executable()
}
