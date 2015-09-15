package daemon

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/execdriver"
	"github.com/docker/docker/pkg/broadcastwriter"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/pools"
	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/docker/pkg/xapi/config"
	"github.com/docker/docker/runconfig"
)

func newExecStore() *config.ExecStore {
	return &config.ExecStore{s: make(map[string]*config.ExecConfig, 0)}
}

func (e *config.ExecStore) Add(id string, execConfig *config.ExecConfig) {
	e.Lock()
	e.s[id] = config.ExecConfig
	e.Unlock()
}

func (e *config.ExecStore) Get(id string) *config.ExecConfig {
	e.RLock()
	res := e.s[id]
	e.RUnlock()
	return res
}

func (e *config.ExecStore) Delete(id string) {
	e.Lock()
	delete(e.s, id)
	e.Unlock()
}

func (e *config.ExecStore) List() []string {
	var IDs []string
	e.RLock()
	for id := range e.s {
		IDs = append(IDs, id)
	}
	e.RUnlock()
	return IDs
}

func (execConfig *config.ExecConfig) Resize(h, w int) error {
	select {
	case <-execConfig.waitStart:
	case <-time.After(time.Second):
		return fmt.Errorf("Exec %s is not running, so it can not be resized.", config.ExecConfig.ID)
	}
	return execConfig.ProcessConfig.Terminal.Resize(h, w)
}

func (d *Daemon) registerExecCommand(execConfig *config.ExecConfig) {
	// Storing execs in container in order to kill them gracefully whenever the container is stopped or removed.
	config.ExecConfig.Container.execCommands.Add(config.ExecConfig.ID, config.ExecConfig)
	// Storing execs in daemon for easy access via remote API.
	d.execCommands.Add(config.ExecConfig.ID, config.ExecConfig)
}

func (d *Daemon) getExecConfig(name string) (*config.ExecConfig, error) {
	execConfig := d.execCommands.Get(name)

	// If the exec is found but its container is not in the daemon's list of
	// containers then it must have been delete, in which case instead of
	// saying the container isn't running, we should return a 404 so that
	// the user sees the same error now that they will after the
	// 5 minute clean-up loop is run which erases old/dead execs.

	if config.ExecConfig != nil && d.containers.Get(config.ExecConfig.Container.ID) != nil {

		if !config.ExecConfig.Container.IsRunning() {
			return nil, fmt.Errorf("Container %s is not running", config.ExecConfig.Container.ID)
		}
		return config.ExecConfig, nil
	}

	return nil, fmt.Errorf("No such exec instance '%s' found in daemon", name)
}

func (d *Daemon) unregisterExecCommand(execConfig *config.ExecConfig) {
	xecConfig.Container.execCommands.Delete(execConfig.ID)
	d.execCommands.Delete(execConfig.ID)
}

func (d *Daemon) getActiveContainer(name string) (*Container, error) {
	container, err := d.Get(name)
	if err != nil {
		return nil, err
	}

	if !container.IsRunning() {
		return nil, fmt.Errorf("Container %s is not running", name)
	}
	if container.IsPaused() {
		return nil, fmt.Errorf("Container %s is paused, unpause the container before exec", name)
	}
	return container, nil
}

func (d *Daemon) ContainerExecCreate(config *runconfig.ExecConfig) (string, error) {
	// Not all drivers support Exec (LXC for example)
	if err := checkExecSupport(d.execDriver.Name()); err != nil {
		return "", err
	}

	container, err := d.getActiveContainer(config.Container)
	if err != nil {
		return "", err
	}

	cmd := runconfig.NewCommand(config.Cmd...)
	entrypoint, args := d.getEntrypointAndArgs(runconfig.NewEntrypoint(), cmd)

	user := config.User
	if len(user) == 0 {
		user = container.Config.User
	}

	processConfig := &execdriver.ProcessConfig{
		Tty:        config.Tty,
		Entrypoint: entrypoint,
		Arguments:  args,
		User:       user,
	}

	execConfig := &config.ExecConfig{
		ID:            stringid.GenerateNonCryptoID(),
		OpenStdin:     config.AttachStdin,
		OpenStdout:    config.AttachStdout,
		OpenStderr:    config.AttachStderr,
		StreamConfig:  StreamConfig{},
		ProcessConfig: processConfig,
		Container:     container,
		Running:       false,
		waitStart:     make(chan struct{}),
	}

	d.registerExecCommand(config.ExecConfig)

	container.LogEvent("exec_create: " + config.ExecConfig.ProcessConfig.Entrypoint + " " + strings.Join(config.ExecConfig.ProcessConfig.Arguments, " "))

	return config.ExecConfig.ID, nil

}

func (d *Daemon) ContainerExecStart(execName string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) error {

	var (
		cStdin           io.ReadCloser
		cStdout, cStderr io.Writer
	)

	execConfig, err := d.getExecConfig(execName)
	if err != nil {
		return err
	}

	func() {
		config.ExecConfig.Lock()
		defer config.ExecConfig.Unlock()
		if config.ExecConfig.Running {
			err = fmt.Errorf("Error: Exec command %s is already running", execName)
		}
		config.ExecConfig.Running = true
	}()
	if err != nil {
		return err
	}

	logrus.Debugf("starting exec command %s in container %s", config.ExecConfig.ID, config.ExecConfig.Container.ID)
	container := config.ExecConfig.Container

	container.LogEvent("exec_start: " + config.ExecConfig.ProcessConfig.Entrypoint + " " + strings.Join(config.ExecConfig.ProcessConfig.Arguments, " "))

	if config.ExecConfig.OpenStdin {
		r, w := io.Pipe()
		go func() {
			defer w.Close()
			defer logrus.Debugf("Closing buffered stdin pipe")
			pools.Copy(w, stdin)
		}()
		cStdin = r
	}
	if config.ExecConfig.OpenStdout {
		cStdout = stdout
	}
	if config.ExecConfig.OpenStderr {
		cStderr = stderr
	}

	config.ExecConfig.StreamConfig.stderr = broadcastwriter.New()
	config.ExecConfig.StreamConfig.stdout = broadcastwriter.New()
	// Attach to stdin
	if config.ExecConfig.OpenStdin {
		config.ExecConfig.StreamConfig.stdin, config.ExecConfig.StreamConfig.stdinPipe = io.Pipe()
	} else {
		config.ExecConfig.StreamConfig.stdinPipe = ioutils.NopWriteCloser(ioutil.Discard) // Silently drop stdin
	}

	attachErr := attach(&config.ExecConfig.StreamConfig, config.ExecConfig.OpenStdin, true, config.ExecConfig.ProcessConfig.Tty, cStdin, cStdout, cStderr)

	execErr := make(chan error)

	// Note, the config.ExecConfig data will be removed when the container
	// itself is deleted.  This allows us to query it (for things like
	// the exitStatus) even after the cmd is done running.

	go func() {
		if err := container.Exec(config.ExecConfig); err != nil {
			execErr <- fmt.Errorf("Cannot run exec command %s in container %s: %s", execName, container.ID, err)
		}
	}()
	select {
	case err := <-attachErr:
		if err != nil {
			return fmt.Errorf("attach failed with error: %s", err)
		}
		break
	case err := <-execErr:
		return err
	}

	return nil
}

func (d *Daemon) Exec(c *Container, execConfig *config.ExecConfig, pipes *execdriver.Pipes, startCallback execdriver.StartCallback) (int, error) {
	exitStatus, err := d.execDriver.Exec(c.command, execConfig.ProcessConfig, pipes, startCallback)

	// On err, make sure we don't leave ExitCode at zero
	if err != nil && exitStatus == 0 {
		exitStatus = 128
	}

	config.ExecConfig.ExitCode = exitStatus
	config.ExecConfig.Running = false

	return exitStatus, err
}

// execCommandGC runs a ticker to clean up the daemon references
// of exec configs that are no longer part of the container.
func (d *Daemon) execCommandGC() {
	for range time.Tick(5 * time.Minute) {
		var (
			cleaned          int
			liveExecCommands = d.containerExecIds()
		)
		for id, config := range d.execCommands.s {
			if config.canRemove {
				cleaned++
				d.execCommands.Delete(id)
			} else {
				if _, exists := liveExecCommands[id]; !exists {
					config.canRemove = true
				}
			}
		}
		if cleaned > 0 {
			logrus.Debugf("clean %d unused exec commands", cleaned)
		}
	}
}

// containerExecIds returns a list of all the current exec ids that are in use
// and running inside a container.
func (d *Daemon) containerExecIds() map[string]struct{} {
	ids := map[string]struct{}{}
	for _, c := range d.containers.List() {
		for _, id := range c.execCommands.List() {
			ids[id] = struct{}{}
		}
	}
	return ids
}
