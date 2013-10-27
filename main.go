package main

import (
	"fmt"
	"github.com/rymo4/life/universe"
  "github.com/realistschuckle/gohaml"
  "net/http"
  //"net/url"
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
    vals := r.URL.Query
    fmt.Printf("%s", vals.Get("gen"))
    fmt.Fprintf(w, "%s", vals.Get("gen"))
    u.Next()
  })

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    serveHaml(w, "web/views/index.haml")
  })

  log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHaml(w http.ResponseWriter, filename string) {
  var scope = make(map[string]interface{})
  scope["lang"] = "HAML"

  content, err := universe.readLines(filename)
  if err != nil {
    // TODO: serve error
    return
  }

  // TODO: serve error
  engine, _ := gohaml.NewEngine(content)
  output := engine.Render(scope)
  fmt.Println(w, output)
}
