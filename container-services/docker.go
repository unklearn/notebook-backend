package containerservices

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/unklearn/notebook-backend/channels"
	"github.com/unklearn/notebook-backend/commands"
	"github.com/unklearn/notebook-backend/connection"
)

/// Implementation of Docker container service
type DockerContainerService struct {
	// The underlying docker client
	client *client.Client
	// Stores a map of networks associated with docker daemon.
	// Allows notebooks/channels to create user defined networks
	networkMap map[string]types.NetworkResource
}

// Execute intents
func (dcs DockerContainerService) ExecuteIntents(intents []commands.ActionIntent, channelRegistry *channels.Registry, conn *connection.MxedWebsocketConn) error {
	for _, intent := range intents {
		switch intent.(type) {
		case commands.ContainerCreateCommandIntent:
			return dcs.CreateNew(context.Background(), channelRegistry, conn, intent.(commands.ContainerCreateCommandIntent))
		default:
			return nil
		}
	}
	return nil
}

// Create a new network for a channel if it does not exist. If it does exist, then
// simply return it. This function is idempotent
func (dcs DockerContainerService) createNetworkForChannel(ctx context.Context, channelId string) (types.NetworkResource, error) {
	// Check if there exists a network with channelId
	dummyResp := types.NetworkResource{}
	netinspectResponse, err := dcs.client.NetworkInspect(ctx, channelId, types.NetworkInspectOptions{})
	if err != nil {
		if client.IsErrNotFound(err) {
			// Create network and return
			_, e := dcs.client.NetworkCreate(ctx, channelId, types.NetworkCreate{})
			if e != nil {
				return dummyResp, e
			} else {
				return dcs.client.NetworkInspect(ctx, channelId, types.NetworkInspectOptions{})
			}
		}
		return dummyResp, err
	}
	return netinspectResponse, nil
}

/// Create a new docker container with given image and tag
func (dcs DockerContainerService) CreateNew(ctx context.Context, channelRegistry *channels.Registry, conn *connection.MxedWebsocketConn, intent commands.ContainerCreateCommandIntent) error {

	exposedPorts := make(map[nat.Port]struct{})
	portMap := make(map[nat.Port][]nat.PortBinding)
	for _, port := range intent.NetworkOptions.Ports {
		p := nat.Port(port + "/tcp")
		exposedPorts[p] = struct{}{}
		portMap[p] = []nat.PortBinding{{
			HostPort: port,
		}}
	}
	ctrConfig := container.Config{
		Image:        fmt.Sprintf("%s:%s", intent.Image, intent.ImageTag),
		Cmd:          intent.Command,
		Env:          intent.EnvVars,
		ExposedPorts: exposedPorts,
	}
	hostConfig := container.HostConfig{
		AutoRemove:   true,
		PortBindings: portMap,
	}

	// Network for channel
	channelNetwork, e := dcs.createNetworkForChannel(ctx, intent.ChannelId)
	if e != nil {
		return e
	}
	endpointsConfig := make(map[string]*network.EndpointSettings)
	endpointsConfig[channelNetwork.Name] = &network.EndpointSettings{
		NetworkID: channelNetwork.ID,
	}
	netConfig := network.NetworkingConfig{
		EndpointsConfig: endpointsConfig,
	}
	resp, err := dcs.client.ContainerCreate(ctx, &ctrConfig, &hostConfig, &netConfig, &v1.Platform{
		Architecture: "amd64",
		OS:           "linux",
	}, "")
	if err != nil {
		return err
	}
	err = dcs.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	return err
}

// func (dcs DockerContainerService) ExecuteCommand(ctx context.Context, containerId string, command []string, useTty bool) (ContainerCommandChannel, error) {
// 	execConfig := types.ExecConfig{Tty: useTty, AttachStdout: true, AttachStderr: true, AttachStdin: true, Cmd: command}
// 	respIdExecCreate, err := dcs.Client.ContainerExecCreate(ctx, containerId, execConfig)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	respId, err := dcs.Client.ContainerExecAttach(ctx, respIdExecCreate.ID, types.ExecStartCheck{
// 		Tty: useTty,
// 	})

// 	if err != nil {
// 		fmt.Println(err)
// 		return ContainerCommandChannel{}, err
// 	}
// 	return ContainerCommandChannel{Id: containerId + "-" + respIdExecCreate.ID[:4], ContainerId: containerId, ExecConn: respId.Conn, ExecReader: respId.Reader}, nil
// }
