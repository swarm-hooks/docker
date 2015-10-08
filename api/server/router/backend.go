package router

import (
	"io"
	"time"

	// everything from here needs to move to types
	// consists mainly of XYZConfig structs to pass information
	"github.com/docker/docker/api/types"                // bunch of return types
	"github.com/docker/docker/api/types/versions/v1p19" // container json
	"github.com/docker/docker/api/types/versions/v1p20" // container json
	"github.com/docker/docker/cliconfig"                // configs, duh
	"github.com/docker/docker/daemon"                   // many configs
	"github.com/docker/docker/daemon/events"            // event format
	"github.com/docker/docker/daemon/network"           // network config
	"github.com/docker/docker/graph"                    // image pull config
	"github.com/docker/docker/registry"                 // registry.searchresults
	"github.com/docker/docker/runconfig"                // configs, duh

	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/idtools"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/parsers/filters"
	"github.com/docker/libnetwork"
)

// Backend is all the methods that need to be implemented to provide
// all the stuff a server wants
type Backend interface {
	SystemInfo() (*types.Info, error)

	Exists(id string) bool

	ContainerCopy(name string, res string) (io.ReadCloser, error)
	ContainerStatPath(name string,
		path string) (stat *types.ContainerPathStat, err error)
	ContainerArchivePath(name string,
		path string) (content io.ReadCloser,
		stat *types.ContainerPathStat, err error)
	ContainerExtractToDir(name, path string, noOverwriteDirNonDir bool,
		content io.Reader) error

	ContainerInspect(name string, size bool) (*types.ContainerJSON, error)
	ContainerInspect120(name string) (*v1p20.ContainerJSON, error)
	ContainerInspectPre120(name string) (*v1p19.ContainerJSON, error)

	Containers(config *daemon.ContainersConfig) ([]*types.Container, error)
	ContainerStats(prefixOrName string,
		config *daemon.ContainerStatsConfig) error
	ContainerLogs(containerName string,
		logsConfig *daemon.ContainerLogsConfig) error
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

	ContainerCreate(name string, config *runconfig.Config,
		hostConfig *runconfig.HostConfig,
		adjustCPUShares bool) (types.ContainerCreateResponse, error)

	ContainerRm(name string, config *daemon.ContainerRmConfig) error
	ContainerResize(name string, height, width int) error
	ContainerExecResize(name string, height, width int) error

	ContainerAttachWithLogs(prefixOrName string,
		c *daemon.ContainerAttachWithLogsConfig) error
	ContainerWsAttachWithLogs(prefixOrName string,
		c *daemon.ContainerWsAttachWithLogsConfig) error

	ContainerExecStart(execName string, stdin io.ReadCloser,
		stdout io.Writer, stderr io.Writer) error
	// two different versions of ExecConfig, oi vey!
	ContainerExecCreate(config *runconfig.ExecConfig) (string, error)
	ContainerExecInspect(id string) (*daemon.ExecConfig, error)

	ExecExists(name string) (bool, error)

	GetUIDGIDMaps() ([]idtools.IDMap, []idtools.IDMap)
	GetRemappedUIDGID() (int, int)

	ImageDelete(imageRef string, force, prune bool) ([]types.ImageDelete, error)
	TagImage(repoName, tag, imageName string, force bool) error
	PullImage(image string, tag string,
		imagePullConfig *graph.ImagePullConfig) error
	LoadImage(inTar io.ReadCloser, outStream io.Writer) error
	LookupImage(name string) (*types.ImageInspect, error)

	ImportImage(src, repo, tag, msg string, inConfig io.ReadCloser,
		outStream io.Writer, containerConfig *runconfig.Config) error
	ExportImage(names []string, outStream io.Writer) error
	PushImage(localName string, imagePushConfig *graph.ImagePushConfig) error
	ListImages(filterArgs, filter string, all bool) ([]*types.Image, error)
	ImageHistory(name string) ([]*types.ImageHistory, error)
	AuthenticateToRegistry(authConfig *cliconfig.AuthConfig) (string, error)
	SearchRegistryForImages(term string,
		authConfig *cliconfig.AuthConfig,
		headers map[string][]string) (*registry.SearchResults, error)
	Volumes(filter string) ([]*types.Volume, error)
	VolumeInspect(name string) (*types.Volume, error)
	VolumeCreate(name, driverName string,
		opts map[string]string) (*types.Volume, error)
	VolumeRm(name string) error

	FindNetwork(idName string) (libnetwork.Network, error)
	GetNetwork(idName string, by int) (libnetwork.Network, error)
	GetNetworksByID(partialID string) []libnetwork.Network
	CreateNetwork(name, driver string, ipam network.IPAM,
		options map[string]string) (libnetwork.Network, error)
	ConnectContainerToNetwork(containerName, networkName string) error
	DisconnectContainerFromNetwork(containerName string,
		network libnetwork.Network) error

	GetEventFilter(filter filters.Args) *events.Filter
	SubscribeToEvents() ([]*jsonmessage.JSONMessage,
		chan interface{}, func())
}
