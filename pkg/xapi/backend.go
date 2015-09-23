package xapi

import (
	"io"
	"time"

	// everything from here needs to move to types
	// consists mainly of XYZConfig structs to pass information
	"github.com/docker/docker/daemon"
	"github.com/docker/docker/runconfig"

	"github.com/docker/docker/pkg/xapi/types"

	"github.com/docker/docker/pkg/archive"
)

// Backend is all the methods that need to be implemented to provide
// all the stuff a server wants
type Backend interface {
	SystemInfo() (*types.Info, error)
	//ContainerStart(name string, hostConfig *runconfig.HostConfig) error
	// NetworkApiRouter()

	Exists(id string) bool

	ContainerCopy(name string, res string) (io.ReadCloser, error)
	ContainerStatPath(name string, path string) (stat *types.ContainerPathStat, err error)
	ContainerArchivePath(name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error)
	ContainerExtractToDir(name, path string, noOverwriteDirNonDir bool, content io.Reader) error

	// ContainerInspect(string)
	ContainerInspect(name string) (*types.ContainerJSON, error)
	ContainerInspect120(name string) (*types.ContainerJSON120, error)
	ContainerInspectPre120(name string) (*types.ContainerJSONPre120, error)

	// Containers(config)
	ContainerStats(prefixOrName string, config *daemon.ContainerStatsConfig) error
	// ContainerLogs(c, logsConfig)
	ContainerExport(name string, out io.Writer) error

	ContainerStart(name string, hostConfig *runconfig.HostConfig) error
	ContainerStop(name string, seconds int) error
	ContainerKill(name string, sig uint64) error
	ContainerRestart(name string, seconds int) error
	ContainerPause(name string) error
	ContainerUnpause(name string) error
	ContainerWait(name string, timeout time.Duration) (int, error)

	ContainerChanges(name string) ([]archive.Change, error)

	ContainerTop(name string, psArgs string) (*types.ContainerProcessList, error)
	ContainerRename(oldName, newName string) error

	// ContainerCreate(name, config, hostConfig)
	// ContainerRm(name, config)
	ContainerResize(name string, height, width int) error
	ContainerExecResize(name string, height, width int) error

	// ContainerAttachWithLogs(cont, attachWithLogsConfig)
	// ContainerWsAttachWithLogs(cont, wsAttachWithLogsConfig)

	ContainerExecStart(execName string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) error
	// two different versions of ExecConfig, oi vey!
	ContainerExecCreate(config *runconfig.ExecConfig) (string, error)
	ContainerExecInspect(id string) (*daemon.ExecConfig, error)
	// Repositories()

	// ContainerAttachWithLogs(cont, attachWithLogsConfig)
	// ContainerWsAttachWithLogs(cont, wsAttachWithLogsConfig)
	// ContainerExecInspect(string)
	// ContainerExecCreate(execConfig)
	// ContainerExecStart(execName, stdin, stdout, stderr)
	// ContainerExecResize(string, height, width)

	// NetworkApiRouter()
	ImageDelete(imageRef string, force, prune bool) ([]types.ImageDelete, error)
	// EventsService
	// RegistryService

	Volumes(filter string) ([]*types.Volume, error)
	VolumeInspect(name string) (*types.Volume, error)
	VolumeCreate(name, driverName string, opts map[string]string) (*types.Volume, error)
	VolumeRm(name string) error
}
