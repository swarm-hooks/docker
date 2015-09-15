package config

// ContainersConfig is used to pass the configuration to the list
// containers command supported by the backend (daemon)
type ContainersConfig struct {
	All     bool
	Since   string
	Before  string
	Limit   int
	Size    bool
	Filters string
}
