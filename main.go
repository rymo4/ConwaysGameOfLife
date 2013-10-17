package main

import (
  "fmt"
  "time"
  "github.com/rymo4/life/universe"
)

const (
  framerate = 10
)

func main() {
  u, err := universe.LoadFromFile("maps/glider_gun.txt")
  u.AtGeneration(100);

  if err != nil {
    fmt.Printf("Please provide a valid file")
    return
  }

  for {
    time.Sleep(time.Second / framerate)
    u.Next()
    u.Show()
  }
}
