package containerservices

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"

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
			cc.ContainerService = rc.ContainerService
			rc.RootConn.RegisterChannel(id, cc)
			// Write back on root channel the id of the container
			rc.RootConn.WriteMessage(2, rc.GetId(), []byte(id))
		}
	case "sync":
		containerId := m["ContainerId"].(string)
		cc := ContainerChannel{Id: containerId, RootConn: rc.RootConn, ContainerService: rc.ContainerService}
		rc.RootConn.RegisterChannel(containerId, cc)
		rc.RootConn.WriteMessage(2, rc.GetId(), []byte("ok"))
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
		log.Printf("%T\n", m["Command"])
		// Convert to string
		var strCommand []string
		for _, el := range m["Command"].([]interface{}) {
			strCommand = append(strCommand, el.(string))
		}
		ccc, err := cc.ContainerService.ExecuteCommand(context.Background(), cc.Id, strCommand)

		if err != nil {
			return len(message), err
		} else {
			ccc.RootConn = cc.RootConn
			ccc.ContainerService = cc.ContainerService
			id := ccc.GetId()
			ccc.RootConn.RegisterChannel(id, ccc)
			ccc.RootConn.WriteMessage(2, id, []byte(id))
			go ccc.Listen()
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
	log.Print("Writing message to connection handler")
	return ccc.ExecConn.Write(message)
}

func (ccc ContainerCommandChannel) Listen() {
	// Listen for messages on exec-reader
	b := make([]byte, 1024)
	ticker := time.Tick(time.Millisecond * 100)
	id := ccc.Id
	go func() {
		for {
			n, err := ccc.ExecReader.Read(b)
			if err == io.EOF {
				break
			}
			// Wait for next set
			<-ticker
			if len(b) > 0 {
				log.Print(string(b[:n]))
				ccc.RootConn.WriteMessage(2, id, b[:n])
			}
		}
		log.Println("Done reading from command")
	}()
}

///
/// root::{"Action":"start","Image":"python","Tag":"3.6"}
/// root::{"Action": "sync", "ContainerId": "af74"}
/// af74::{"Action":"exec-command","Command": ["bash"]}
/// af74-de5aa317c774859b50c6e8865fc298deeeb9b7fd69f26b0fa17961444655d3d9::ls
