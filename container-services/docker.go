// package containerservices

// import (
// 	"context"
// 	"fmt"

// 	"github.com/docker/docker/api/types"
// 	"github.com/docker/docker/api/types/container"
// 	"github.com/docker/docker/api/types/network"
// 	"github.com/docker/docker/client"
// 	"github.com/docker/go-connections/nat"
// 	v1 "github.com/opencontainers/image-spec/specs-go/v1"
// )

// /// Implementation of Docker container service
// type DockerContainerService struct {
// 	Client *client.Client
// }

// /// Create a new docker container with given image and tag
// func (dcs DockerContainerService) CreateNew(ctx context.Context, image string, tag string, containerCommand []string, containerEnv []string, ports []string) (ContainerChannel, error) {

// 	exposedPorts := make(map[nat.Port]struct{})
// 	portMap := make(map[nat.Port][]nat.PortBinding)
// 	for _, port := range ports {
// 		p := nat.Port(port + "/tcp")
// 		exposedPorts[p] = struct{}{}
// 		portMap[p] = []nat.PortBinding{{
// 			HostPort: port,
// 		}}
// 	}
// 	ctrConfig := container.Config{
// 		Image:        fmt.Sprintf("%s:%s", image, tag),
// 		Cmd:          containerCommand,
// 		Env:          containerEnv,
// 		ExposedPorts: exposedPorts,
// 	}
// 	hostConfig := container.HostConfig{
// 		AutoRemove:   true,
// 		PortBindings: portMap,
// 		NetworkMode:  "host",
// 	}
// 	netConfig := network.NetworkingConfig{}
// 	resp, err := dcs.Client.ContainerCreate(ctx, &ctrConfig, &hostConfig, &netConfig, &v1.Platform{
// 		Architecture: "amd64",
// 		OS:           "linux",
// 	}, "")
// 	if err != nil {
// 		return ContainerChannel{}, err
// 	}
// 	if err := dcs.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
// 		panic(err)
// 	}
// 	cc := ContainerChannel{Id: resp.ID[:4]}
// 	return cc, nil
// }

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
