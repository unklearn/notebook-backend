package containerservices

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

/// Implementation of Docker container service
type DockerContainerService struct {
	Client *client.Client
}

/// Create a new docker container with given image and tag
func (dcs DockerContainerService) CreateNew(ctx context.Context, image string, tag string) (ContainerChannel, error) {
	ctrConfig := container.Config{
		Image: fmt.Sprintf("%s:%s", image, tag),
		Cmd:   []string{"sleep", "infinity"},
	}
	hostConfig := container.HostConfig{
		AutoRemove: true,
	}
	netConfig := network.NetworkingConfig{}
	resp, err := dcs.Client.ContainerCreate(ctx, &ctrConfig, &hostConfig, &netConfig, "")
	if err != nil {
		return ContainerChannel{}, err
	}
	if err := dcs.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	cc := ContainerChannel{Id: resp.ID[:4]}
	return cc, nil
}

func (dcs DockerContainerService) ExecuteCommand(ctx context.Context, containerId string, command []string) (ContainerCommandChannel, error) {
	execConfig := types.ExecConfig{Tty: true, AttachStdout: true, AttachStderr: true, AttachStdin: true, Cmd: command}
	respIdExecCreate, err := dcs.Client.ContainerExecCreate(ctx, containerId, execConfig)
	if err != nil {
		fmt.Println(err)
	}
	respId, err := dcs.Client.ContainerExecAttach(ctx, respIdExecCreate.ID, execConfig)

	if err != nil {
		return ContainerCommandChannel{}, err
	}
	return ContainerCommandChannel{Id: containerId + "-" + respIdExecCreate.ID, ContainerId: containerId, ExecConn: respId.Conn, ExecReader: respId.Reader}, nil
}
