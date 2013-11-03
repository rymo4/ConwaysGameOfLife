package main

import (
	"encoding/json"
	"fmt"
	"github.com/rymo4/life/universe"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
  "code.google.com/p/go.net/websocket"
  "time"
)

type hub struct {
    // Registered connections.
    connections map[*connection]bool
    // Inbound messages from the connections.
    broadcast chan string
    // Register requests from the connections.
    register chan *connection
    // Unregister requests from connections.
    unregister chan *connection
}

var h = hub{
    broadcast:   make(chan string),
    register:    make(chan *connection),
    unregister:  make(chan *connection),
    connections: make(map[*connection]bool),
}

func (h *hub) run() {
    for {
        select {
        case c := <-h.register:
            h.connections[c] = true
            log.Println("New Connection. Total: ", len(h.connections))
        case c := <-h.unregister:
            delete(h.connections, c)
            close(c.send)
        case m := <-h.broadcast:
            for c := range h.connections {
                select {
                case c.send <- m:
                default:
                    delete(h.connections, c)
                    close(c.send)
                    go c.ws.Close()
                }
            }
        }
    }
}

type connection struct {
    // The websocket connection.
    ws *websocket.Conn
    // Buffered channel of outbound messages.
    send chan string
}

func (c *connection) reader() {
    for {
        var message string
        err := websocket.Message.Receive(c.ws, &message)
        if err != nil {
            break
        }
        h.broadcast <- message
    }
    c.ws.Close()
}

func (c *connection) writer() {
  ticker := time.NewTicker(100 * time.Millisecond)
  log.Println("starting to write to clients")
  u, _ := universe.LoadFromFile("maps/glider_gun.txt")
  for _ = range ticker.C {
    err := websocket.Message.Send(c.ws, u.CanonicalString())
    if err != nil {
      c.ws.Close()
      break
    }
    u.Next()
  }
  log.Println("stopped")
}

func wsHandler(ws *websocket.Conn) {
    log.Println("ws handler")
    c := &connection{send: make(chan string, 256), ws: ws}
    h.register <- c
    log.Println("Registered")
    defer func() { h.unregister <- c }()
    c.writer()
    //c.reader()
}

//// Protocol:
// width,height, [t/f]|i1,j1,i2,j2....
// where i = col # for a living cell
// and j = row # for living cell
// t/f for toroidal or not
func main() {
	log.Print("Starting webserver.")

	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Responding to ", r.URL.Path)
		val := r.FormValue("state")
		u := universe.LoadFromCanonicalString(val)
		u.Next()
		fmt.Fprintf(w, "%s\n", u.CanonicalString())
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/"))))

	http.HandleFunc("/maps", func(w http.ResponseWriter, r *http.Request) {
		log.Print("Responding to ", r.URL.Path)
		mapName := r.FormValue("mapName")
		u, _ := universe.LoadFromFile("./maps/" + mapName + ".txt")
		fmt.Fprintf(w, "%s\n", u.CanonicalString())
	})

    http.Handle("/ws", websocket.Handler(wsHandler))

	http.HandleFunc("/mapslist", func(w http.ResponseWriter, r *http.Request) {
		files, _ := ioutil.ReadDir("maps")
		filenames := make([]string, len(files))
		for i, file := range files {
			parts := strings.Split(file.Name(), ".")
			filenames[i] = parts[0]
		}
		b, _ := json.Marshal(filenames)
		fmt.Fprintf(w, string(b))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHTML(w, "./web/main.html")
	})

  go h.run()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHTML(w http.ResponseWriter, filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		// TODO: serve error
		log.Fatalf("Cannot find %s\n", filename)
		return
	}

	io.WriteString(w, string(content))
}
