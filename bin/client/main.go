package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solipsis/ball"
)

func main() {

	c := &ball.Client{}
	go c.Run("http://localhost:8090")

	if err := ebiten.RunGame(c); err != nil {
		log.Fatal(err)
	}
	return
}
