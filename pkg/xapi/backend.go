package xapi

import (
	"github.com/docker/docker/pkg/xapi/config"
	"github.com/docker/docker/pkg/xapi/types"
	"github.com/docker/docker/registry"
	"github.com/docker/docker/runconfig"
)

// Backend is all the methods that need to be implemented to provide
// all the stuff a server wants
type Backend interface {
	SystemInfo() (*types.Info, error)
	// NetworkApiRouter()

	// interesting inconsistency, Container, Container, ContainerJSON
	Get(prefixOrName string) (*config.Container, error)  // I think maybe this shouldn't be something exported and used.
	Containers(*config.ContainersConfig) ([]*types.Container, error)
	ContainerInspect(name string) (*types.ContainerJSON, error)

	ContainerStats(name string, config *config.ContainerStatsConfig) error
	ContainerLogs(container *types.Container, config *config.ContainerLogsConfig) error
	ContainerStart(name string, hostConfig *runconfig.HostConfig) error
	ContainerStop(name string, seconds int) error
	ContainerKill(name string, sig uint64) error

	ContainerRestart(name string, seconds int) error
	ContainerPause(name string) error

	// ContainerCopy(string, string)
	// ContainerStatPath(string, string)
	// ContainerArchivePath(string, string)
	// ContainerExtractToDir(string, string, bool, io.Reader)
	// ContainerExport(string, w)

	// ContainerPause(string)
	// ContainerUnpause(string)
	// ContainerWait(string, time.duration)
	// ContainerChanges(string)
	// ContainerTop(string, string)
	// ContainerRename(name, newName)
	// ContainerCreate(name, config, hostConfig)
	// ContainerRm(name, config)
	// ContainerResize(string, height, width)
	// ContainerAttachWithLogs(cont, attachWithLogsConfig)
	// ContainerWsAttachWithLogs(cont, wsAttachWithLogsConfig)

	// ContainerExecInspect(string)o
	// ContainerExecCreate(execConfig)
	// ContainerExecStart(execName, stdin, stdout, stderr)
	// ContainerExecResize(string, height, width)

	// Repositories()

	// ContainerAttachWithLogs(cont, attachWithLogsConfig)
	// ContainerWsAttachWithLogs(cont, wsAttachWithLogsConfig)
	// ContainerExecInspect(string)
	// ContainerExecCreate(execConfig)
	// ContainerExecStart(execName, stdin, stdout, stderr)
	// ContainerExecResize(string, height, width)

	// NetworkApiRouter()
	// ImageDelete(name, force, noprune)
	// EventsService
	// RegistryService
	RegistryService() *registry.Service
	// ContainerInspectPre120(namevar)
}
