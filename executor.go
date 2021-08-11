package main

import (
	"context"
	"encoding/json"
	"log"
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
	ExecuteContainerCommand(ctx context.Context, intent commands.ContainerExecuteCommandIntent) (*channels.BidirectionalContainerConduit, error)
}

func (ce CommandExecutor) createNewContainerSaga(intent commands.ContainerCreateCommandIntent) {
	// Business logic is encapsulated in this saga
	containerId, err := ce.IContainerCommandService.CreateNew(context.Background(), intent)
	// Let conn know that new channel has been registered
	failed, _ := json.Marshal(commands.ContainerStatusResponse{Id: containerId, Hash: intent.Hash, Status: "failed"})
	conn := ce.conn

	if err != nil {
		// Write a message stating that container has failed
		conn.WriteMessage(intent.ChannelId, string(channels.ContainerStatusEventName), failed)
		return
	}
	// Create new container channel
	conn.RegisterChannel(containerId, channels.NewContainerChannel(containerId))

	// Let conn know that new channel has been registered
	response, _ := json.Marshal(commands.ContainerStatusResponse{Id: containerId, Hash: intent.Hash, Status: "pending"})
	// Write a message stating that container has started
	conn.WriteMessage(intent.ChannelId, string(channels.ContainerStatusEventName), response)

	// Wait for container status
	go ce.waitForContainerSaga(intent.ChannelId, commands.ContainerWaitCommandIntent{ContainerId: containerId})
}

func (ce CommandExecutor) waitForContainerSaga(channelId string, intent commands.ContainerWaitCommandIntent) {
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
			conn.WriteMessage(channelId, string(channels.ContainerStatusEventName), out)
			break
		}
		times += 1
		if (times * sleepTime) > timeout {
			statusResponse.Status = "timed-out"
			out, _ := json.Marshal(statusResponse)
			conn.WriteMessage(channelId, string(channels.ContainerStatusEventName), out)
			break
		}
		if status == "running" {
			statusResponse.Status = "running"
			out, _ := json.Marshal(statusResponse)
			conn.WriteMessage(channelId, string(channels.ContainerStatusEventName), out)
			break
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

func (ce CommandExecutor) executeContainerCommandSaga(intent commands.ContainerExecuteCommandIntent) {
	conduit, err := ce.ExecuteContainerCommand(context.Background(), intent)
	conn := ce.conn
	if err != nil {
		failed, _ := json.Marshal(commands.ContainerCommandStatusResponse{ExecId: conduit.ExecId, CellId: intent.CellId, Status: "failed", Reason: err.Error()})
		// Write a message stating that container command execution has failed
		conn.WriteMessage(intent.ContainerId, string(channels.ContainerCommandStatusEventName), failed)
		return
	}
	// Create new container command channel
	ch := channels.NewContainerCommandChannel(intent.CellId, conduit)
	conn.RegisterChannel(intent.CellId, ch)
	// No wait here, some commands never send output
	success, _ := json.Marshal(commands.ContainerCommandStatusResponse{ExecId: conduit.ExecId, CellId: intent.CellId, Status: "success"})
	conn.WriteMessage(intent.ContainerId, string(channels.ContainerCommandStatusEventName), success)
	// Run go-routines to wait on command output, and command input if interactive is true
	if intent.Interactive {
		go ce.listenForComandOutput(conduit, intent.ContainerId, intent.CellId)
	}
}

func (ce CommandExecutor) listenForComandOutput(conduit *channels.BidirectionalContainerConduit, containerId string, cellId string) {
	go func() {
	L:
		for {
		S:
			select {
			case read := <-conduit.ReadChan:
				ce.conn.WriteMessage(cellId, string(channels.ContainerCommandOutputEventName), read)
				break S
			case cmd := <-conduit.CommChan:
				// Parse command. If it is a close op, exit the loop and update status
				// Other ops are pending
				if cmd == "quit" {
					stopped, _ := json.Marshal(commands.ContainerCommandStatusResponse{ExecId: conduit.ExecId, CellId: cellId, Status: "stopped"})
					ce.conn.WriteMessage(containerId, string(channels.ContainerCommandStatusEventName), stopped)
					break L
				}
				break S
			}
		}
	}()
}

// Executor channel <- receive intent and run it

func (ce CommandExecutor) ExecuteIntents() {
	// Create a container channel and register it
	for intent := range ce.dispatch {
		log.Printf("Handling intent %s\n", intent.ToString())
		switch i := intent.(type) {
		case commands.ContainerCreateCommandIntent:
			ce.createNewContainerSaga(i)
			continue
		// case commands.ContainerWaitCommandIntent:
		// 	ce.waitForContainerSaga(i)
		// 	continue
		case commands.ContainerExecuteCommandIntent:
			ce.executeContainerCommandSaga(i)
			continue
		default:
			log.Printf("Got typo %T\n", intent)
			continue
		}
	}
}

func (ce CommandExecutor) DispatchIntents(intents []commands.ActionIntent) {
	for _, intent := range intents {
		ce.dispatch <- intent
	}
}

func (ce CommandExecutor) ConnectionHandler() {
	mx := ce.conn
	go ce.ExecuteIntents()
	for {
		d, err := mx.ReadMessage()
		if err != nil {
			log.Println("Error while reading from connection:", err)
			break
		}
		ch, e := mx.GetChannelById(d.ChannelId)
		if e != nil {
			// Respond with bad error-code
			log.Printf("Error while retrieving channel %s", d.ChannelId)
			continue
		}
		intents, e := ch.HandleMessage(d.EventName, d.Payload)
		if e != nil {
			// Write the error to the end user
			mx.WriteMessage(d.ChannelId, d.EventName, []byte(e.Error()))
		}
		// Dispatch intents
		ce.DispatchIntents(intents)
	}
}
