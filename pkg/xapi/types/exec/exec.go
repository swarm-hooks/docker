package types

// I don't know if this is in the right place or not, it seems like
// it's more for the http interface rather than a binary interface...
import (
	"io"
	"os/exec"
	"time"

	//"github.com/docker/docker/pkg/xapi/config"

	// TODO Windows: Factor out ulimit
	"github.com/docker/docker/pkg/ulimit"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/configs"
)

// Network settings of the container
type Network struct {
	Interface      *NetworkInterface `json:"interface"` // if interface is nil then networking is disabled
	Mtu            int               `json:"mtu"`
	ContainerID    string            `json:"container_id"` // id of the container to join network.
	NamespacePath  string            `json:"namespace_path"`
	HostNetworking bool              `json:"host_networking"`
}

// Ipc settings of the container
// It is for IPC namespace setting. Usually different containers
// have their own IPC namespace, however this specifies to use
// an existing IPC namespace.
// You can join the host's or a container's IPC namespace.
type Ipc struct {
	ContainerID string `json:"container_id"` // id of the container to join ipc.
	HostIpc     bool   `json:"host_ipc"`
}

// Pid settings of the container
// It is for PID namespace setting. Usually different containers
// have their own PID namespace, however this specifies to use
// an existing PID namespace.
// Joining the host's PID namespace is currently the only supported
// option.
type Pid struct {
	HostPid bool `json:"host_pid"`
}

// UTS settings of the container
// It is for UTS namespace setting. Usually different containers
// have their own UTS namespace, however this specifies to use
// an existing UTS namespace.
// Joining the host's UTS namespace is currently the only supported
// option.
type UTS struct {
	HostUTS bool `json:"host_uts"`
}

// NetworkInterface contains all network configs for a driver
type NetworkInterface struct {
	Gateway              string `json:"gateway"`
	IPAddress            string `json:"ip"`
	IPPrefixLen          int    `json:"ip_prefix_len"`
	MacAddress           string `json:"mac"`
	Bridge               string `json:"bridge"`
	GlobalIPv6Address    string `json:"global_ipv6"`
	LinkLocalIPv6Address string `json:"link_local_ipv6"`
	GlobalIPv6PrefixLen  int    `json:"global_ipv6_prefix_len"`
	IPv6Gateway          string `json:"ipv6_gateway"`
	HairpinMode          bool   `json:"hairpin_mode"`
}

// Resources contains all resource configs for a driver.
// Currently these are all for cgroup configs.
// TODO Windows: Factor out ulimit.Rlimit
type Resources struct {
	Memory           int64            `json:"memory"`
	MemorySwap       int64            `json:"memory_swap"`
	CPUShares        int64            `json:"cpu_shares"`
	CpusetCpus       string           `json:"cpuset_cpus"`
	CpusetMems       string           `json:"cpuset_mems"`
	CPUPeriod        int64            `json:"cpu_period"`
	CPUQuota         int64            `json:"cpu_quota"`
	BlkioWeight      int64            `json:"blkio_weight"`
	Rlimits          []*ulimit.Rlimit `json:"rlimits"`
	OomKillDisable   bool             `json:"oom_kill_disable"`
	MemorySwappiness int64            `json:"memory_swappiness"`
}

// ResourceStats contains information about resource usage by a container.
type ResourceStats struct {
	*libcontainer.Stats
	Read        time.Time `json:"read"`
	MemoryLimit int64     `json:"memory_limit"`
	SystemUsage uint64    `json:"system_usage"`
}

// Mount contains information for a mount operation.
type Mount struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Writable    bool   `json:"writable"`
	Private     bool   `json:"private"`
	Slave       bool   `json:"slave"`
}

// ProcessConfig describes a process that will be run inside a container.
type ProcessConfig struct {
	exec.Cmd `json:"-"`

	Privileged  bool     `json:"privileged"`
	User        string   `json:"user"`
	Tty         bool     `json:"tty"`
	Entrypoint  string   `json:"entrypoint"`
	Arguments   []string `json:"arguments"`
	Terminal    Terminal `json:"-"` // standard or tty terminal
	Console     string   `json:"-"` // dev/console path
	ConsoleSize [2]int   `json:"-"` // h,w of initial console size
}

// Command wrapps an os/exec.Cmd to add more metadata
//
// TODO Windows: Factor out unused fields such as LxcConfig, AppArmorProfile,
// and CgroupParent.
type Command struct {
	ID                 string            `json:"id"`
	Rootfs             string            `json:"rootfs"` // root fs of the container
	ReadonlyRootfs     bool              `json:"readonly_rootfs"`
	InitPath           string            `json:"initpath"` // dockerinit
	WorkingDir         string            `json:"working_dir"`
	ConfigPath         string            `json:"config_path"` // this should be able to be removed when the lxc template is moved into the driver
	Network            *Network          `json:"network"`
	Ipc                *Ipc              `json:"ipc"`
	Pid                *Pid              `json:"pid"`
	UTS                *UTS              `json:"uts"`
	Resources          *Resources        `json:"resources"`
	Mounts             []Mount           `json:"mounts"`
	AllowedDevices     []*configs.Device `json:"allowed_devices"`
	AutoCreatedDevices []*configs.Device `json:"autocreated_devices"`
	CapAdd             []string          `json:"cap_add"`
	CapDrop            []string          `json:"cap_drop"`
	GroupAdd           []string          `json:"group_add"`
	ContainerPid       int               `json:"container_pid"`  // the pid for the process inside a container
	ProcessConfig      ProcessConfig     `json:"process_config"` // Describes the init process of the container.
	ProcessLabel       string            `json:"process_label"`
	MountLabel         string            `json:"mount_label"`
	LxcConfig          []string          `json:"lxc_config"`
	AppArmorProfile    string            `json:"apparmor_profile"`
	CgroupParent       string            `json:"cgroup_parent"` // The parent cgroup for this command.
	FirstStart         bool              `json:"first_start"`
	LayerPaths         []string          `json:"layer_paths"` // Windows needs to know the layer paths and folder for a command
	LayerFolder        string            `json:"layer_folder"`
}

// Terminal represents a pseudo TTY, it is for when
// using a container interactively.
type Terminal interface {
	io.Closer
	Resize(height, width int) error
}

// ExitStatus provides exit reasons for a container.
type ExitStatus struct {
	// The exit code with which the container exited.
	ExitCode int

	// Whether the container encountered an OOM.
	OOMKilled bool
}

// Pipes is a wrapper around a container's output for
// stdin, stdout, stderr
type Pipes struct {
	Stdin          io.ReadCloser
	Stdout, Stderr io.Writer
}
