package channels

import (
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
	ContainerStartEventName  RootChannelEventNames = "root/container-start"
	ContainerStopEventName   RootChannelEventNames = "root/container-stop"
	ContainerStatusEventName RootChannelEventNames = "root/container-status"
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
	ContainerExecuteCommandEventName ContainerChannelEventNames = "container/execute-command"
	ContainerCommandStatusEventName  ContainerChannelEventNames = "container/command-status"
	ContainerSyncFileEventName       ContainerChannelEventNames = "container/sync-file"
	ContainerSyncFileOutputEventName ContainerChannelEventNames = "container/file-output"
	ContainerWriteToFile             ContainerChannelEventNames = "container/write-file"
)

// Return id for external callers
func (cc ContainerChannel) GetId() string {
	return cc.id
}

// HandleMessage takes care of a given event and payload. If payload cannot be handled, error
// is returned
func (cc ContainerChannel) HandleMessage(eventName string, payload []byte) ([]commands.ActionIntent, error) {
	switch eventName {
	case string(ContainerExecuteCommandEventName):
		// Parse body into StartContainer
		c, e := commands.NewContainerExecuteCommandIntent(cc.id, payload)
		if e != nil {
			return []commands.ActionIntent{}, e
		}
		return []commands.ActionIntent{c}, nil
	case string(ContainerSyncFileEventName):
		// Sync file to/from the container.
		c, e := commands.NewSyncFileIntent(cc.id, payload)
		if e != nil {
			return []commands.ActionIntent{}, e
		}
		return []commands.ActionIntent{c}, nil
	default:
		break
	}
	return []commands.ActionIntent{}, fmt.Errorf("unknown event name %s", eventName)
}

// A container conduit acts as a bidirectional communication channel between
// the container and the webserver
type BidirectionalContainerConduit struct {
	// Used to refer to the running command, if required
	ExecId string
	// Used for stdout and stderr streams
	ReadChan chan []byte
	// Used for stdin
	WriteChan chan []byte
	// Used for communicating error codes etc
	CommChan chan string
}

type ContainerCommandChannelEventNames string

const (
	ContainerCommandOutputEventName ContainerCommandChannelEventNames = "command/output"
	ContainerCommandInputEventname  ContainerCommandChannelEventNames = "command/input"
)

type ContainerCommandChannel struct {
	id          string
	conduit     *BidirectionalContainerConduit
	emptyIntent []commands.ActionIntent
}

// Constructor function for new command channel
func NewContainerCommandChannel(id string, conduit *BidirectionalContainerConduit) *ContainerCommandChannel {
	return &ContainerCommandChannel{id: id, conduit: conduit, emptyIntent: []commands.ActionIntent{}}
}

// Return id for external callers
func (cce ContainerCommandChannel) GetId() string {
	return cce.id
}

// HandleMessage takes care of a given event and payload. If payload cannot be handled, error
// is returned
func (cce ContainerCommandChannel) HandleMessage(eventName string, payload []byte) ([]commands.ActionIntent, error) {
	switch eventName {
	case string(ContainerExecuteCommandEventName):
		// Parse body into StartContainer
		c, e := commands.NewContainerExecuteCommandIntent(cce.id, payload)
		if e != nil {
			return cce.emptyIntent, e
		}
		return []commands.ActionIntent{c}, nil
	case string(ContainerCommandInputEventname):
		cce.conduit.WriteChan <- payload
		return cce.emptyIntent, nil
	default:
		break
	}
	return cce.emptyIntent, fmt.Errorf("unknown event name %s", eventName)
}
