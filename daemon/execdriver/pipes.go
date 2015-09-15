package execdriver

import (
	"io"

	exec "github.com/docker/docker/pkg/xapi/types/exec"
)

// NewPipes returns a wrapper around a container's output
func NewPipes(stdin io.ReadCloser, stdout, stderr io.Writer, useStdin bool) *exec.Pipes {
	p := &exec.Pipes{
		Stdout: stdout,
		Stderr: stderr,
	}
	if useStdin {
		p.Stdin = stdin
	}
	return p
}
