package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/unklearn/notebook-backend/connection"
	containerservices "github.com/unklearn/notebook-backend/container-services"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func CheckOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: CheckOrigin,
} // use default options

var dcs = containerservices.DockerContainerService{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	// Maps execId to a multiplexed connection
	mx := connection.MxedWebsocketConn{Conn: c, Delimiter: "::"}
	if mx.ChannelMap == nil {
		mx.ChannelMap = make(map[string]connection.Channel)
	}
	var rootChannel = containerservices.RootChannel{RootConn: mx, Id: "root", ContainerService: dcs}

	// Register channels
	mx.RegisterChannel("root", &rootChannel)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		err := mx.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
	}
}

// Main serve function that runs the HTTP handler routes as well as the websocker handler.
// The routes will handle notebook related API calls, and the websocket will relay container
// outputs and execution status of a cell.
func main() {
	http.HandleFunc("/websocket", echo)
	// http.HandleFunc("/container-create", contCreate)
	// Create new docker client
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	// Set client
	dcs.Client = cli
	log.Fatal(http.ListenAndServe(*addr, nil))
}
