package channels

import (
	"encoding/json"
	"fmt"

	"github.com/unklearn/notebook-backend/commands"
)

type IChannel interface {
	HandleMessage(eventName string, payload []byte) ([]commands.ActionIntent, error)
	GetId() string
}

// RootChannel is used by the notebook to communicate any top level socket events
// Examples include stop-container, start-container etc. Since these execute outside
// of a given container, we use the root channel, the root channel is registered
// against the connection id, which is the id of the RootChannel
type RootChannel struct {
	// The id of the channel
	id string
}

func NewRootChannel(id string) *RootChannel {
	return &RootChannel{id: id}
}

type RootChannelEventNames string

const (
	ContainerStartEventName   RootChannelEventNames = "container:start"
	ContainerStopEventName    RootChannelEventNames = "container:stop"
	ContainerStartedEventName RootChannelEventNames = "container:started"
	ContainerStatusEventName  RootChannelEventNames = "container:status"
)

// Return id for external callers
func (rc RootChannel) GetId() string {
	return rc.id
}

// HandleMessage takes care of a given event and payload. If payload cannot be handled, error
// is returned
func (rc RootChannel) HandleMessage(eventName string, payload []byte) ([]commands.ActionIntent, error) {
	switch eventName {
	case string(ContainerStartEventName):
		// Parse body into StartContainer
		c, e := commands.NewContainerCreateCommandIntent(rc.id, payload)
		if e != nil {
			return []commands.ActionIntent{}, e
		}
		return []commands.ActionIntent{c}, nil
	default:
		break
	}
	return []commands.ActionIntent{}, fmt.Errorf("unknown event name %s", eventName)
}

// A container channel wraps a single container
// and can be used to perform actions inside the container
type ContainerChannel struct {
	// id of the channel
	id string
}

func NewContainerChannel(id string) *ContainerChannel {
	return &ContainerChannel{id: id}
}

type ContainerChannelEventNames string

const (
	ExecuteCommand ContainerChannelEventNames = "execute:command"
	CommandStatus  ContainerChannelEventNames = "command:status"
	ReadFile       ContainerChannelEventNames = "read:file"
	WriteToFile    ContainerChannelEventNames = "write:file"
)

// Return id for external callers
func (cc ContainerChannel) GetId() string {
	return cc.id
}

// HandleMessage takes care of a given event and payload. If payload cannot be handled, error
// is returned
func (cc ContainerChannel) HandleMessage(eventName string, payload []byte) ([]commands.ActionIntent, error) {
	switch eventName {
	case string(ExecuteCommand):
		// Parse body into StartContainer
		var cmd []string
		json.Unmarshal(payload, &cmd)
		c := commands.NewContainerExecuteCommandIntent(cc.id, true, true, -1, cmd)
		return []commands.ActionIntent{*c}, nil
	default:
		break
	}
	return []commands.ActionIntent{}, fmt.Errorf("unknown event name %s", eventName)
}
