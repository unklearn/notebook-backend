package containerservices

import (
	"context"
)

// Interface for container service
type IContainerService interface {
	// Create a new container with image and tag as parameters. Returns the created container id
	CreateNew(ctx context.Context, image string, tag string, containerCommand []string, containerEnv []string, ports []string) (ContainerChannel, error)
	ExecuteCommand(ctx context.Context, containerId string, command []string, useTty bool) (ContainerCommandChannel, error)
}
