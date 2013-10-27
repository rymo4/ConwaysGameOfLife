package main

import (
	"fmt"
	"github.com/rymo4/life/universe"
  "net/http"
  //"net/url"
  "io/ioutil"
  "io"
  "log"
)

const (
	framerate = 10
)

func main() {
  log.Print("Starting webserver.")

  http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
    log.Print("Responding to ", r.URL.Path)
    u, err := universe.LoadFromFile("maps/glider_gun.txt")
    if err != nil {
      fmt.Printf("Please provide a valid file")
      return
    }
    //vals := r.URL.Query
    //fmt.Printf("%s", vals.Get("gen"))
    //fmt.Fprintf(w, "%s", vals.Get("gen"))
    u.Next()
  })

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    serveHTML(w, "./web/main.html")
  })

  http.Handle("/static/",  http.StripPrefix("/static/", http.FileServer(http.Dir("./web/"))))

  log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHTML(w http.ResponseWriter, filename string) {

  content, err := ioutil.ReadFile(filename)
  if err != nil {
    // TODO: serve error
    fmt.Printf("Cannot find %s\n", filename);
    return
  }

  io.WriteString(w, string(content))
}
