package daemon

import (
	"fmt"
	"io"

	"github.com/docker/docker/pkg/stdcopy"
)

// ContainerAttachWithLogsConfig holds the streams to use when connecting to a container to view logs.
type ContainerAttachWithLogsConfig struct {
	InStream                       io.ReadCloser
	OutStream                      io.Writer
	ErrStream                      io.Writer
	SetupStreams                   func() (io.ReadCloser, io.Writer, error)
	UseStdin, UseStdout, UseStderr bool
	Logs, Stream                   bool
}

// ContainerAttachWithLogs attaches to logs according to the config passed in. See ContainerAttachWithLogsConfig.
func (daemon *Daemon) ContainerAttachWithLogs(prefixOrName string, c *ContainerAttachWithLogsConfig) error {
	container, err := daemon.Get(prefixOrName)
	if err != nil {
		return err
	}

	var errStream io.Writer

	if nil == c.SetupStreams {
		return fmt.Errorf("no streams to set up connection")
	}

	InStream, OutStream, err := c.SetupStreams()
	// overwrite streams if we changed them to pass back to the calling context
	c.InStream = InStream
	c.OutStream = OutStream

	if nil != err {
		return err
	}

	if !container.Config.Tty {
		errStream = stdcopy.NewStdWriter(OutStream, stdcopy.Stderr)
		OutStream = stdcopy.NewStdWriter(OutStream, stdcopy.Stdout)
	} else {
		errStream = OutStream
	}

	var stdin io.ReadCloser
	var stdout, stderr io.Writer

	if c.UseStdin {
		stdin = InStream
	}
	if c.UseStdout {
		stdout = OutStream
	}
	if c.UseStderr {
		stderr = errStream
	}

	return container.attachWithLogs(stdin, stdout, stderr, c.Logs, c.Stream)
}

// ContainerWsAttachWithLogsConfig attach with websockets, since all
// stream data is delegated to the websocket to handle, there
type ContainerWsAttachWithLogsConfig struct {
	InStream             io.ReadCloser
	OutStream, ErrStream io.Writer
	Logs, Stream         bool
}

// ContainerWsAttachWithLogs websocket connection
func (daemon *Daemon) ContainerWsAttachWithLogs(prefixOrName string, c *ContainerWsAttachWithLogsConfig) error {
	container, err := daemon.Get(prefixOrName)
	if err != nil {
		return err
	}
	return container.attachWithLogs(c.InStream, c.OutStream, c.ErrStream, c.Logs, c.Stream)
}
