package config

import (
	"sync"

	// get rid of this somehow
	execdriver "github.com/docker/docker/pkg/xapi/types/exec"
)

type ExecConfig struct {
	sync.Mutex
	ID            string
	Running       bool
	ExitCode      int
	ProcessConfig *execdriver.ProcessConfig
	StreamConfig
	OpenStdin  bool
	OpenStderr bool
	OpenStdout bool
	Container  *Container
	canRemove  bool

	// waitStart will be closed immediately after the exec is really started.
	waitStart chan struct{}
}

type ExecStore struct {
	s map[string]*ExecConfig
	sync.RWMutex
}
