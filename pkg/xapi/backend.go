package xapi

import (
	"io"
	"time"

	// everything from here needs to move to types
	// consists mainly of XYZConfig structs to pass information
	"github.com/docker/docker/context"
	"github.com/docker/docker/daemon"
	"github.com/docker/docker/runconfig"

	"github.com/docker/docker/pkg/xapi/types"

	"github.com/docker/docker/pkg/archive"
)

// Backend is all the methods that need to be implemented to provide
// all the stuff a server wants
type Backend interface {
	SystemInfo(ctx context.Context) (*types.Info, error)
	//ContainerStart(name string, hostConfig *runconfig.HostConfig) error

	Exists(ctx context.Context, id string) bool

	ContainerCopy(ctx context.Context, name string, res string) (io.ReadCloser, error)
	ContainerStatPath(ctx context.Context, name string, path string) (stat *types.ContainerPathStat, err error)
	ContainerArchivePath(ctx context.Context, name string, path string) (content io.ReadCloser, stat *types.ContainerPathStat, err error)
	ContainerExtractToDir(ctx context.Context, name, path string, noOverwriteDirNonDir bool, content io.Reader) error

	// ContainerInspect(string)
	ContainerInspect(ctx context.Context, name string) (*types.ContainerJSON, error)
	ContainerInspect120(ctx context.Context, name string) (*types.ContainerJSON120, error)
	ContainerInspectPre120(ctx context.Context, name string) (*types.ContainerJSONPre120, error)

	// Containers(config)
	Containers(ctx context.Context, config *daemon.ContainersConfig) ([]*types.Container, error)
	ContainerStats(ctx context.Context, prefixOrName string, config *daemon.ContainerStatsConfig) error
	// ContainerLogs(c, logsConfig)
	ContainerExport(ctx context.Context, name string, out io.Writer) error

	ContainerStart(ctx context.Context, name string, hostConfig *runconfig.HostConfig) error
	ContainerStop(ctx context.Context, name string, seconds int) error
	ContainerKill(ctx context.Context, name string, sig uint64) error
	ContainerRestart(ctx context.Context, name string, seconds int) error
	ContainerPause(ctx context.Context, name string) error
	ContainerUnpause(ctx context.Context, name string) error
	ContainerWait(ctx context.Context, name string, timeout time.Duration) (int, error)

	ContainerChanges(ctx context.Context, name string) ([]archive.Change, error)

	ContainerTop(ctx context.Context, name string, psArgs string) (*types.ContainerProcessList, error)
	ContainerRename(ctx context.Context, oldName, newName string) error

	// ContainerCreate(name, config, hostConfig)
	ContainerRm(ctx context.Context, name string, config *daemon.ContainerRmConfig) error
	ContainerResize(ctx context.Context, name string, height, width int) error
	ContainerExecResize(ctx context.Context, name string, height, width int) error

	// ContainerAttachWithLogs(cont, attachWithLogsConfig)
	// ContainerWsAttachWithLogs(cont, wsAttachWithLogsConfig)

	ContainerExecStart(ctx context.Context, execName string, stdin io.ReadCloser, stdout io.Writer, stderr io.Writer) error
	// two different versions of ExecConfig, oi vey!
	ContainerExecCreate(ctx context.Context, config *runconfig.ExecConfig) (string, error)
	ContainerExecInspect(ctx context.Context, id string) (*daemon.ExecConfig, error)
	// Repositories()

	// ContainerAttachWithLogs(cont, attachWithLogsConfig)
	// ContainerWsAttachWithLogs(cont, wsAttachWithLogsConfig)

	// NetworkApiRouter()
	ImageDelete(ctx context.Context, imageRef string, force, prune bool) ([]types.ImageDelete, error)
	// EventsService
	// RegistryService

	Volumes(ctx context.Context, filter string) ([]*types.Volume, error)
	VolumeInspect(ctx context.Context, name string) (*types.Volume, error)
	VolumeCreate(ctx context.Context, name, driverName string, opts map[string]string) (*types.Volume, error)
	VolumeRm(ctx context.Context, name string) error
}
