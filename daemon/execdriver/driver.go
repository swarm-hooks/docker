package execdriver

import (
	"errors"

	exec "github.com/docker/docker/pkg/xapi/types/exec"
)

// Context is a generic key value pair that allows
// arbatrary data to be sent
type Context map[string]string

// Define error messages
var (
	ErrNotRunning              = errors.New("Container is not running")
	ErrWaitTimeoutReached      = errors.New("Wait timeout reached")
	ErrDriverAlreadyRegistered = errors.New("A driver already registered this docker init function")
	ErrDriverNotFound          = errors.New("The requested docker init has not been found")
)

// StartCallback defines a callback function.
// It's used by 'Run' and 'Exec', does some work in parent process
// after child process is started.
type StartCallback func(*exec.ProcessConfig, int)

// Info is driver specific information based on
// processes registered with the driver
type Info interface {
	IsRunning() bool
}

// Driver is an interface for drivers to implement
// including all basic functions a driver should have
type Driver interface {
	// Run executes the process, blocks until the process exits and returns
	// the exit code. It's the last stage on Docker side for running a container.
	Run(c *exec.Command, pipes *exec.Pipes, startCallback StartCallback) (exec.ExitStatus, error)

	// Exec executes the process in an existing container, blocks until the
	// process exits and returns the exit code.
	Exec(c *exec.Command, processConfig *exec.ProcessConfig, pipes *exec.Pipes, startCallback StartCallback) (int, error)

	// Kill sends signals to process in container.
	Kill(c *exec.Command, sig int) error

	// Pause pauses a container.
	Pause(c *exec.Command) error

	// Unpause unpauses a container.
	Unpause(c *exec.Command) error

	// Name returns the name of the driver.
	Name() string

	// Info returns the configuration stored in the driver struct,
	// "temporary" hack (until we move state from core to plugins).
	Info(id string) Info

	// GetPidsForContainer returns a list of pid for the processes running in a container.
	GetPidsForContainer(id string) ([]int, error)

	// Terminate kills a container by sending signal SIGKILL.
	Terminate(c *exec.Command) error

	// Clean removes all traces of container exec.
	Clean(id string) error

	// Stats returns resource stats for a running container
	Stats(id string) (*exec.ResourceStats, error)
}
