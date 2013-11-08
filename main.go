package main

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/rymo4/life/universe"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type hub struct {
	connections map[*connection]bool
	register    chan *connection
	unregister  chan *connection
}

var h = hub{
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = false
			log.Println("New Connection. Total: ", len(h.connections))
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		}
	}
}

type connection struct {
	ws      *websocket.Conn
	Initial string
	send    chan string
}

func (c *connection) reader() {
	for {
		var canonical string
		err := websocket.Message.Receive(c.ws, &canonical)
		if err != nil {
			break
		}
		c.Initial = canonical
		h.connections[c] = true
	}
	c.ws.Close()
}

func (c *connection) writer() {
	ticker := time.NewTicker(100 * time.Millisecond)
	log.Println("starting to write to clients")
	var u *universe.Universe
	for _ = range ticker.C {
		if c.Initial == "" {
			// Wait for a start state before sending the data at it
			log.Print("No start state received")
			continue
		}
		if h.connections[c] {
			// use start state, then mark it to not be used again
			h.connections[c] = false
			u = universe.LoadFromCanonicalString(c.Initial)
		}
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
	log.Print("wsHandler")
	c := &connection{send: make(chan string, 256), ws: ws}
	h.register <- c
	log.Println("Registered new connection")
	defer func() { h.unregister <- c }()
	go c.reader()
	c.writer()
}

//// Protocol:
// width,height,[true/false]|i1,j1,i2,j2....
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

	http.Handle("/stream", websocket.Handler(wsHandler))

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
