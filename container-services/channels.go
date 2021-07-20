package containerservices

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net"

	"github.com/unklearn/notebook-backend/connection"
)

type RootChannel struct {
	Id string
	connection.Channel
	RootConn         connection.MxedWebsocketConn
	ContainerService IContainerService
}

func (rc RootChannel) GetId() string {
	return rc.Id
}

func (rc RootChannel) Write(message []byte) (int, error) {
	// Parse message intent
	var f interface{}
	err := json.Unmarshal(message, &f)
	if err != nil {
		return len(message), err
	}
	m := f.(map[string]interface{})
	// Look for exec handlers
	action := m["Action"]
	log.Print(action)
	switch action {
	case "start":
		cc, err := rc.ContainerService.CreateNew(context.Background(), m["Image"].(string), m["Tag"].(string))
		if err != nil {
			return len(message), err
		} else {
			// Add new container channel
			id := cc.GetId()
			cc.RootConn = rc.RootConn
			rc.RootConn.RegisterChannel(id, cc)
			// Write back on root channel the id of the container
			rc.RootConn.WriteMessage(2, rc.GetId(), []byte(id))
		}
	}
	return len(message), nil
}

type ContainerChannel struct {
	connection.Channel
	Id               string
	RootConn         connection.MxedWebsocketConn
	ContainerService IContainerService
}

func (cc ContainerChannel) GetId() string {
	return cc.Id
}

func (cc ContainerChannel) Write(message []byte) (int, error) {
	// Parse message intent
	var f interface{}
	err := json.Unmarshal(message, &f)
	if err != nil {
		return len(message), err
	}
	m := f.(map[string]interface{})
	// Look for exec handlers
	action := m["Action"]
	switch action {
	case "exec-command":
		ccc, err := cc.ContainerService.ExecuteCommand(context.Background(), cc.Id, m["Command"].([]string))

		if err != nil {
			return len(message), err
		} else {
			ccc.RootConn = cc.RootConn
			id := ccc.GetId()
			ccc.RootConn.RegisterChannel(id, ccc)
			ccc.RootConn.WriteMessage(2, id, []byte(id))
		}
	}
	return len(message), nil
}

type ContainerCommandChannel struct {
	connection.Channel
	Id               string
	ContainerId      string
	RootConn         connection.MxedWebsocketConn
	ExecConn         net.Conn
	ExecReader       *bufio.Reader
	ContainerService IContainerService
}

func (ccc ContainerCommandChannel) GetId() string {
	return ccc.Id
}

func (ccc ContainerCommandChannel) Write(message []byte) (int, error) {
	// Write it to underlying handler
	return ccc.ExecConn.Write(message)
}
