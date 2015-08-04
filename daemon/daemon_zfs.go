// +build !exclude_graphdriver_zfs,linux !exclude_graphdriver_zfs,freebsd

package daemon

import (
	// FIXME: I have no idea what the justification is for this import. I
	// imagine some interaction with the build tags and something
	// else. golint wants a comment, and is even satisfied by a blank comment.
	_ "github.com/docker/docker/daemon/graphdriver/zfs"
)
