package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/unklearn/notebook-backend/channels"
	"github.com/unklearn/notebook-backend/commands"
	"github.com/unklearn/notebook-backend/connection"
)

type CommandExecutor struct {
	// A channel to receive action commands.
	dispatch chan commands.ActionIntent
	// Container service
	IContainerCommandService
	// The multiplexed connection
	conn *connection.MxedWebsocketConn
}

func NewCommandExecutor(cs IContainerCommandService, conn *connection.MxedWebsocketConn) *CommandExecutor {
	ce := &CommandExecutor{
		dispatch:                 make(chan commands.ActionIntent, 1),
		IContainerCommandService: cs,
		conn:                     conn,
	}
	// Start a go routine that listens and executes ExecuteIntents
	return ce
}

type IContainerCommandService interface {
	CreateNew(ctx context.Context, intent commands.ContainerCreateCommandIntent) (containerId string, err error)
	GetContainerStatus(ctx context.Context, containerId string) (status string, err error)
}

func (ce CommandExecutor) CreateNewContainerSaga(intent commands.ContainerCreateCommandIntent) {
	// Business logic is encapsulated in this saga
	containerId, err := ce.IContainerCommandService.CreateNew(context.Background(), intent)
	// Let conn know that new channel has been registered
	failed, _ := json.Marshal(commands.ContainerStatusResponse{Id: containerId, Status: "failed"})
	conn := ce.conn

	if err != nil {
		// Write a message stating that container has failed
		conn.WriteMessage(intent.ChannelId, string(channels.ContainerStatusEventName), failed)
		return
	}
	// Create new container channel
	conn.RegisterChannel(containerId, channels.NewContainerChannel(containerId))

	// Let conn know that new channel has been registered
	response, _ := json.Marshal(commands.ContainerStatusResponse{Id: containerId, Status: "pending"})
	// Write a message stating that container has started
	conn.WriteMessage(intent.ChannelId, string(channels.ContainerStartedEventName), response)

	// Wait for container status
	go ce.WaitForContainerSaga(commands.ContainerWaitCommandIntent{ContainerId: containerId})
}

func (ce CommandExecutor) WaitForContainerSaga(intent commands.ContainerWaitCommandIntent) {
	times := 0
	timeout := intent.Timeout
	if timeout == 0 {
		timeout = 15
	}
	conn := ce.conn
	sleepTime := 3
	statusResponse := commands.ContainerStatusResponse{Id: intent.ContainerId, Status: "failed"}
	for {
		// Inspect the container
		status, e := ce.IContainerCommandService.GetContainerStatus(context.Background(), intent.ContainerId)
		if e != nil {
			statusResponse.Status = "error"
			out, _ := json.Marshal(statusResponse)
			conn.WriteMessage(intent.ContainerId, string(channels.ContainerStatusEventName), out)
			break
		}
		times += 1
		if (times * sleepTime) > timeout {
			statusResponse.Status = "timed-out"
			out, _ := json.Marshal(statusResponse)
			conn.WriteMessage(intent.ContainerId, string(channels.ContainerStatusEventName), out)
			break
		}
		if status == "running" {
			statusResponse.Status = "running"
			out, _ := json.Marshal(statusResponse)
			conn.WriteMessage(intent.ContainerId, string(channels.ContainerStatusEventName), out)
			break
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

// Executor channel <- receive intent and run it

func (ce CommandExecutor) ExecuteIntents() {
	// Create a container channel and register it
	for intent := range ce.dispatch {
		switch i := intent.(type) {
		case commands.ContainerCreateCommandIntent:
			ce.CreateNewContainerSaga(i)
		case commands.ContainerWaitCommandIntent:
			ce.WaitForContainerSaga(i)
		default:
			continue
		}
	}
}

func (ce CommandExecutor) DispatchIntents(intents []commands.ActionIntent) {
	for _, intent := range intents {
		ce.dispatch <- intent
	}
}
