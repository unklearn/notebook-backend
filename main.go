package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/unklearn/notebook-backend/connection"
	containerservices "github.com/unklearn/notebook-backend/container-services"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

var dcs = containerservices.DockerContainerService{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	// Maps execId to a multiplexed connection
	mx := connection.MxedWebsocketConn{Conn: c, Delimiter: "::"}

	var rootChannel = containerservices.RootChannel{RootConn: mx, Id: "root", ContainerService: dcs}

	// Register channels
	mx.RegisterChannel("root", &rootChannel)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// b := make([]byte, 1024)
	// ticker := time.Tick(time.Millisecond * 100)
	// go func() {
	// 	for {
	// 		n, err := respId.Reader.Read(b)
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		// Wait for next set
	// 		<-ticker
	// 		if len(b) > 0 {
	// 			mxed.WriteMessage(2, b[:n])
	// 		}
	// 	}
	// 	log.Println("Done reading")
	// }()

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
	http.HandleFunc("/", home)
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

// func contCreate(w http.ResponseWriter, r *http.Request) {
// 	queryParams := r.URL.Query()
// 	image := queryParams.Get("image")
// 	tag := queryParams.Get("tag")
// 	if len(tag) == 0 {
// 		tag = "latest"
// 	}
// 	id, err := dcs.CreateNew(r.Context(), image, tag)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	fmt.Fprintf(w, "{\"id\": \"%s\"}", id)
// }

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/websocket")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style type="text/css">
  
.xterm {
    position: relative;
    user-select: none;
    -ms-user-select: none;
    -webkit-user-select: none;
}

.xterm.focus,
.xterm:focus {
    outline: none;
}

.xterm .xterm-helpers {
    position: absolute;
    top: 0;
    /**
     * The z-index of the helpers must be higher than the canvases in order for
     * IMEs to appear on top.
     */
    z-index: 5;
}

.xterm .xterm-helper-textarea {
    padding: 0;
    border: 0;
    margin: 0;
    /* Move textarea out of the screen to the far left, so that the cursor is not visible */
    position: absolute;
    opacity: 0;
    left: -9999em;
    top: 0;
    width: 0;
    height: 0;
    z-index: -5;
    /** Prevent wrapping so the IME appears against the textarea at the correct position */
    white-space: nowrap;
    overflow: hidden;
    resize: none;
}

.xterm .composition-view {
    /* TODO: Composition position got messed up somewhere */
    background: #000;
    color: #FFF;
    display: none;
    position: absolute;
    white-space: nowrap;
    z-index: 1;
}

.xterm .composition-view.active {
    display: block;
}

.xterm .xterm-viewport {
    /* On OS X this is required in order for the scroll bar to appear fully opaque */
    background-color: #000;
    overflow-y: scroll;
    cursor: default;
    position: absolute;
    right: 0;
    left: 0;
    top: 0;
    bottom: 0;
}

.xterm .xterm-screen {
    position: relative;
}

.xterm .xterm-screen canvas {
    position: absolute;
    left: 0;
    top: 0;
}

.xterm .xterm-scroll-area {
    visibility: hidden;
}

.xterm-char-measure-element {
    display: inline-block;
    visibility: hidden;
    position: absolute;
    top: 0;
    left: -9999em;
    line-height: normal;
}

.xterm {
    cursor: text;
}

.xterm.enable-mouse-events {
    /* When mouse events are enabled (eg. tmux), revert to the standard pointer cursor */
    cursor: default;
}

.xterm.xterm-cursor-pointer {
    cursor: pointer;
}

.xterm.column-select.focus {
    /* Column selection mode */
    cursor: crosshair;
}

.xterm .xterm-accessibility,
.xterm .xterm-message {
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    right: 0;
    z-index: 10;
    color: transparent;
}

.xterm .live-region {
    position: absolute;
    left: -9999px;
    width: 1px;
    height: 1px;
    overflow: hidden;
}

.xterm-dim {
    opacity: 0.5;
}

.xterm-underline {
    text-decoration: underline;
}

.xterm-strikethrough {
    text-decoration: line-through;
}
</style>
<script type="text/javascript" src="https://cdn.jsdelivr.net/npm/xterm@4.13.0/lib/xterm.min.js"></script>
<script type="text/javascript" src="https://cdn.jsdelivr.net/npm/xterm-addon-attach@0.6.0/lib/xterm-addon-attach.min.js"></script>
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
//            const term = new Terminal({convertEol: true});
//const attachAddon = new AttachAddon.AttachAddon(ws, { bidirectional: true});

// Attach the socket to term
//term.loadAddon(attachAddon);
//term.open(document.getElementById('xterm-container'));
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<div id="xterm-container">
</div>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">

<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
