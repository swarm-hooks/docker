package xapi

import (
	//"time"
	"github.com/docker/docker/api/types"
)

// Backend is all the methods that need to be implemented to provide
// all the stuff a server wants
type Backend interface {
	SystemInfo() (*types.Info, error)
	// NetworkApiRouter()

	// Get(string) // get container
	// ContainerCopy(string, string)
	// ContainerStatPath(string, string)
	// ContainerArchivePath(string, string)
	// ContainerExtractToDir(string, string, bool, io.Reader)
	// ContainerInspect(string)
	// Containers(config)
	// ContainerStats(string, config)
	// ContainerLogs(c, logsConfig)
	// ContainerExport(string, w)
	// ContainerStart(string, hostConfig)
	// ContainerStop(string, seconds)
	// ContainerKill(name, sig)
	// ContainerRestart(string, timeout)
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

	// ContainerExecInspect(string)
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
	// ContainerInspectPre120(namevar)
}
