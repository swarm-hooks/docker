package config

import "io"

// used to pass configuration to the stats command
type ContainerStatsConfig struct {
	Stream    bool
	OutStream io.Writer
	Stop      <-chan bool
}
