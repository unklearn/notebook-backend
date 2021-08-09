package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/unklearn/notebook-backend/channels"
	"github.com/unklearn/notebook-backend/connection"
	containerservices "github.com/unklearn/notebook-backend/container-services"
	"github.com/unklearn/notebook-backend/notebooks"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func CheckOrigin(r *http.Request) bool {
	return true
}

var dcs *containerservices.DockerContainerService

var upgrader = websocket.Upgrader{
	CheckOrigin: CheckOrigin,
} // use default options

func HandleWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	// Maps execId to a multiplexed connection
	mx := connection.NewMxedWebsocketConn(c)
	// Use the path for registering channels to conn
	vars := mux.Vars(r)
	rootChName := vars["notebookId"]
	mx.RegisterChannel(rootChName, channels.NewRootChannel(rootChName))

	if err != nil {
		log.Print("upgrade error:", err)
		return
	}
	defer c.Close()
	executor := NewCommandExecutor(dcs, mx)
	// Run connector handler
	executor.ConnectionHandler()
}

// Main serve function that runs the HTTP handler routes as well as the websocker handler.
// The routes will handle notebook related API calls, and the websocket will relay container
// outputs and execution status of a cell.
func main() {
	// Create new docker client
	cli, err := client.NewClientWithOpts()
	router := mux.NewRouter()
	if err != nil {
		panic(err)
	}
	dcs = containerservices.NewDockerContainerService(cli)

	// Register websocket handler
	router.HandleFunc("/websocket/{notebookId}", HandleWS)
	router.HandleFunc("/api/v1/notebooks", notebooks.NotebooksHandler)
	router.HandleFunc("/api/v1/notebooks/{notebookId}", notebooks.NotebookHandler)
	// Root router route
	http.Handle("/", router)
	log.Printf("Listening on %v\n", *addr)
	// Listen and serve
	log.Fatal(http.ListenAndServe(*addr, nil))
}
