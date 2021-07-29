package channels

import (
	"fmt"

	"github.com/unklearn/notebook-backend/commands"
)

type IChannel interface {
	HandleMessage(eventName string, payload []byte) (commands.ActionIntent, error)
	GetId()
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
	ContainerStartEventName RootChannelEventNames = "container:start"
	ContainerStopEventName  RootChannelEventNames = "container:stop"
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
		c := commands.ContainerCreateCommandIntent{}
		c.Parse(rc.id, payload)
		return []commands.ActionIntent{&c}, nil
	default:
		break
	}
	return []commands.ActionIntent{}, fmt.Errorf("unknown event name %s", eventName)
}
