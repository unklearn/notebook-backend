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
			rc.RootConn.WriteMessage(2, rc.GetId(), append([]byte("container-started::"), []byte(id)...))
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
	case "exec-terminal":
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
			ccc.RootConn.WriteMessage(2, "root", append([]byte("terminal-started::"), []byte(id)...))
			time.Sleep(1 * time.Second)
			go ccc.Listen()
		}
	case "sync-file":
		filePath := m["Path"].(string)
		fileContents := m["Content"]
		var ccc ContainerCommandChannel
		switch fileContents := fileContents.(type) {
		case nil:
			// Do not write file, simply sync
			ccc, _ = cc.ContainerService.ExecuteCommand(context.Background(), cc.Id, []string{"cat", filePath})
		case string:
			ccc, _ = cc.ContainerService.ExecuteCommand(context.Background(), cc.Id, []string{"echo", fileContents, ">", filePath, "&&", "cat", filePath})
		}
		ccc.RootConn = cc.RootConn
		// Immediately grab the output and send
		b := make([]byte, 1024)
		var contents []byte
		go func() {
			for {
				n, err := ccc.ExecReader.Read(b)
				if err == io.EOF {
					break
				}

				if len(b) > 0 {
					if contents == nil {
						contents = b[:n]
					} else {

						contents = append(contents, b[:n]...)
					}
					contents = append(contents, []byte("\r\n")...)
				}
			}
			ccc.RootConn.WriteMessage(2, "root", append([]byte("file-contents::"), contents...))
		}()

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
	return ccc.ExecConn.Write(message)
}

func (ccc ContainerCommandChannel) Listen() {
	// Listen for messages on exec-reader
	b := make([]byte, 1024)
	//ticker := time.NewTicker(time.Millisecond * 100)
	id := ccc.Id
	go func() {
		for {
			n, err := ccc.ExecReader.Read(b)
			if err == io.EOF {
				break
			}
			// Wait for next set
			// <-ticker.C
			if len(b) > 0 {
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
/// af74-cd3d70cc51934f5d49a9871f2c461e996c90ec1273d78627ba343ef1c27086e3::ls
