// package containerservices

// import (
// 	"context"
// 	"errors"
// 	"net"

// 	"github.com/unklearn/notebook-backend/connection"
// )

// // Interface for container service
// type IContainerService interface {
// 	// Create a new container with image and tag as parameters. Returns the created container id
// 	CreateNew(ctx context.Context, image string, tag string, containerCommand []string, containerEnv []string, ports []string) (ContainerChannel, error)
// 	ExecuteCommand(ctx context.Context, containerId string, command []string, useTty bool) (ContainerCommandChannel, error)
// }

// type IChannelRegistry interface {
// 	// Register channel stores the channel id > channel mapping for a given connection id
// 	RegisterChannel(connId string, channelId string, channel IChannel)
// 	// Retrieve a channel if it exists, otherwise return error
// 	GetChannel(connId string, channelId string) (IChannel, error)
// }

// // Proxy interface to net.Conn
// type IChannel interface {
// 	net.Conn
// 	GetId() string
// }

// type RootChannel struct {
// 	Id string
// 	IChannel
// 	RootConn         connection.MxedWebsocketConn
// 	ContainerService IContainerService
// }

// type ContainerChannel struct {
// 	IChannel
// 	Id               string
// 	RootConn         connection.MxedWebsocketConn
// 	ContainerService IContainerService
// }

// type IContainerIntent interface {
// 	Execute(*connection.MxedWebsocketConn, IContainerService) error
// }

// // A container execute command intent stores information about
// // command execution in a container. The command will use interactive
// // mode and TTY mode.
// type ContainerExecuteCommandIntent struct {
// 	IContainerIntent
// 	containerID string
// 	command     []string
// }

// func (ceci ContainerExecuteCommandIntent) Execute(channelRegistry IChannelRegistry, conn *connection.MxedWebsocketConn, cs IContainerService) error {
// 	// Execute command
// 	ccc, _ := cs.ExecuteCommand(context.Background(), ceci.containerID, ceci.command, true)
// 	// Register container command channel
// 	channelRegistry.RegisterChannel("conn", ccc.GetId(), ccc)
// 	return errors.New("omg")
// }
