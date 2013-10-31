package main

import (
	"fmt"
	"github.com/rymo4/life/universe"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveHTML(w, "./web/main.html")
	})

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
