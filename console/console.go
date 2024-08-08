package console

import "io"

type Console interface {
	In() io.Reader
	Out() io.Writer
	Error() io.Writer
	Executable() (string, error)
	Args() []string
}
