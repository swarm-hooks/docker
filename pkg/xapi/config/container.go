package config

import (
	"sync"
	"time"

	"github.com/docker/docker/runconfig"

	"github.com/docker/docker/pkg/xapi/types/exec"
	"github.com/docker/docker/volume"
)

// CommonContainer holds the settings for a container which are applicable
// across all platforms supported by the daemon.
type CommonContainer struct {
	StreamConfig

	*State `json:"State"` // Needed for remote api version <= 1.11
	root   string         // Path to the "home" of the container, including metadata.
	basefs string         // Path to the graphdriver mountpoint

	ID                       string
	Created                  time.Time
	Path                     string
	Args                     []string
	Config                   *runconfig.Config
	ImageID                  string `json:"Image"`
	NetworkSettings          *Settings
	LogPath                  string
	Name                     string
	Driver                   string
	ExecDriver               string
	MountLabel, ProcessLabel string
	RestartCount             int
	HasBeenStartedBefore     bool
	hostConfig               *runconfig.HostConfig
	command                  *types.Command
	monitor                  *containerMonitor
	execCommands             *execStore
	//daemon                   *Backend
	// logDriver for closing
	logDriver Logger
	logCopier *Copier
}

type Container struct {
	CommonContainer

	// Fields below here are platform specific.
	activeLinks     map[string]*Link
	AppArmorProfile string
	HostnamePath    string
	HostsPath       string
	MountPoints     map[string]*mountPoint
	ResolvConfPath  string
	UpdateDns       bool
	Volumes         map[string]string // Deprecated since 1.7, kept for backwards compatibility
	VolumesRW       map[string]bool   // Deprecated since 1.7, kept for backwards compatibility
}

type State struct {
	sync.Mutex
	Running           bool
	Paused            bool
	Restarting        bool
	OOMKilled         bool
	removalInProgress bool // Not need for this to be persistent on disk.
	Dead              bool
	Pid               int
	ExitCode          int
	Error             string // contains last known error when starting the container
	StartedAt         time.Time
	FinishedAt        time.Time
	waitChan          chan struct{}
}

// containerMonitor monitors the execution of a container's main process.
// If a restart policy is specified for the container the monitor will ensure that the
// process is restarted based on the rules of the policy.  When the container is finally stopped
// the monitor will reset and cleanup any of the container resources such as networking allocations
// and the rootfs
type containerMonitor struct {
	mux sync.Mutex

	// container is the container being monitored
	container *Container

	// restartPolicy is the current policy being applied to the container monitor
	restartPolicy runconfig.RestartPolicy

	// failureCount is the number of times the container has failed to
	// start in a row
	failureCount int

	// shouldStop signals the monitor that the next time the container exits it is
	// either because docker or the user asked for the container to be stopped
	shouldStop bool

	// startSignal is a channel that is closes after the container initially starts
	startSignal chan struct{}

	// stopChan is used to signal to the monitor whenever there is a wait for the
	// next restart so that the timeIncrement is not honored and the user is not
	// left waiting for nothing to happen during this time
	stopChan chan struct{}

	// timeIncrement is the amount of time to wait between restarts
	// this is in milliseconds
	timeIncrement int

	// lastStartTime is the time which the monitor last exec'd the container's process
	lastStartTime time.Time
}

//// From volumes

// non json and json, UGGHHGHGH

// mountPoint is the intersection point between a volume and a container. It
// specifies which volume is to be used and where inside a container it should
// be mounted.
type mountPoint struct {
	Name        string
	Destination string
	Driver      string
	RW          bool
	Volume      volume.Volume `json:"-"`
	Source      string
	Mode        string `json:"Relabel"` // Originally field was `Relabel`"
}

//// From daemon

// Message is datastructure that represents record from some container.
type Message struct {
	ContainerID string
	Line        []byte
	Source      string
	Timestamp   time.Time
}
