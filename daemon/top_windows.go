package daemon

import (
	"fmt"

	"github.com/docker/docker/context"
	derr "github.com/docker/docker/errors"
	"github.com/docker/docker/pkg/xapi/types"
)

// ContainerTop is not supported on Windows and returns an error.
func (daemon *Daemon) ContainerTop(ctx context.Context, name string, psArgs string) (*types.ContainerProcessList, error) {
	return nil, derr.ErrorCodeNoTop
}
