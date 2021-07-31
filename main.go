package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/unklearn/notebook-backend/channels"
	"github.com/unklearn/notebook-backend/connection"
	containerservices "github.com/unklearn/notebook-backend/container-services"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func CheckOrigin(r *http.Request) bool {
	return true
}

var dcs *containerservices.DockerContainerService

var upgrader = websocket.Upgrader{
	CheckOrigin: CheckOrigin,
} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	// Maps execId to a multiplexed connection
	mx := connection.NewMxedWebsocketConn(c)
	mx.RegisterChannel("root", channels.NewRootChannel("root"))

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	executor := NewCommandExecutor(dcs, mx)
	go executor.ExecuteIntents()
	for {
		d, err := mx.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		ch, e := mx.GetChannelById(d.ChannelId)
		if e != nil {
			// Respond with bad error-code
			break
		}
		intents, _ := ch.HandleMessage(d.EventName, d.Payload)
		go executor.DispatchIntents(intents)
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
	dcs = containerservices.NewDockerContainerService(cli)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
