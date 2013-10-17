package main

import (
	"fmt"
	"github.com/rymo4/life/universe"
	"time"
)

const (
	framerate = 10
)

func main() {
	u, err := universe.LoadFromFile("maps/glider_gun.txt")

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
