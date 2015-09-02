package swapi

// ContainersService provides the containers management aspects of a Docker
// server.
type ContainersService interface {
	List(*ListContainersParams) ([]*Container, error)
	Create(string, []string) (*ContainerCreateResponse, error)
	Start(string) (int, error)
}

// ExtendedContainerService provides the containers management aspects of a
// Docker server, including those who cannot be implemented over a strictly
// HTTP/1.x compliant connection.
type ExtendedContainerService interface {
	ContainersService
}

type ListContainersParams struct {
	All     bool
	Limit   int
	Since   int
	Before  int
	Size    bool
	Filters map[string][]string
}

// ContainerCreateResponse contains the information returned to a client on the
// creation of a new container.
type ContainerCreateResponse struct {
	// ID is the ID of the created container.
	ID string `json:"Id"`

	// Warnings are any warnings encountered during the creation of the container.
	Warnings []string `json:"Warnings"`
}

// GET  "/containers/json"
type Port struct {
	IP          string `json:",omitempty"`
	PrivatePort int
	PublicPort  int `json:",omitempty"`
	Type        string
}

// Container holds data for an existing container.
type Container struct {
	ID         string `json:"Id"`
	Names      []string
	Image      string
	Command    string
	Created    int
	Ports      []Port
	SizeRw     int `json:",omitempty"`
	SizeRootFs int `json:",omitempty"`
	Labels     map[string]string
	Status     string
}
