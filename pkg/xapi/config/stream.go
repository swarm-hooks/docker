package config

import (
	"io"

	"github.com/docker/docker/pkg/broadcastwriter"
)

type StreamConfig struct {
	stdout    *broadcastwriter.BroadcastWriter
	stderr    *broadcastwriter.BroadcastWriter
	stdin     io.ReadCloser
	stdinPipe io.WriteCloser
}
