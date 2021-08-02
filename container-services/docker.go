package containerservices

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/unklearn/notebook-backend/commands"
)

/// Implementation of Docker container service
type DockerContainerService struct {
	// The underlying docker client
	client *client.Client
	// Stores a map of networks associated with docker daemon.
	// Allows notebooks/channels to create user defined networks
	networkMap map[string]types.NetworkResource
}

func (dcs DockerContainerService) GetClient() *client.Client {
	return dcs.client
}

func NewDockerContainerService(c *client.Client) *DockerContainerService {
	return &DockerContainerService{client: c}
}

const NETWORK_NAME = "unk_default_network"

// Create a new network for a channel if it does not exist. If it does exist, then
// simply return it. This function is idempotent
func (dcs DockerContainerService) createNetworkForChannel(ctx context.Context) (types.NetworkResource, error) {
	// Check if there exists a network with channelId
	dummyResp := types.NetworkResource{}
	netLsResponse, err := dcs.client.NetworkList(ctx, types.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: NETWORK_NAME}, filters.KeyValuePair{Key: "type", Value: "custom"}),
	})
	if err != nil {
		return dummyResp, err
	}
	if len(netLsResponse) == 1 && err == nil {
		return netLsResponse[0], nil
	} else {
		networkResponse, e := dcs.client.NetworkCreate(ctx, NETWORK_NAME, types.NetworkCreate{})
		if e != nil {
			return dummyResp, e
		} else {
			return dcs.client.NetworkInspect(ctx, networkResponse.ID, types.NetworkInspectOptions{})
		}
	}
}

func (dcs DockerContainerService) EnsureImage(ctx context.Context, image string, tag string, repoUrl string) error {
	// TODO: Use repoUrl
	_, e := dcs.client.ImageList(ctx, types.ImageListOptions{All: true})
	if e != nil {
		return e
	}
	return nil
}

// Return the status of a running container. If container is missing
// error is returned
func (dcs DockerContainerService) GetContainerStatus(ctx context.Context, containerId string) (string, error) {
	ctr, e := dcs.client.ContainerInspect(ctx, containerId)
	if e != nil {
		return "", e
	}
	return ctr.State.Status, nil
}

// Create a new docker container with given image and tag
// Returns containerId and err if any
func (dcs DockerContainerService) CreateNew(ctx context.Context, intent commands.ContainerCreateCommandIntent) (string, error) {

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
	channelNetwork, e := dcs.createNetworkForChannel(ctx)
	if e != nil {
		return "", e
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
	}, intent.Name)
	if err != nil {
		return "", err
	}
	err = dcs.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	return resp.ID, err
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
