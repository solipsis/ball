package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solipsis/ball"
)

func main() {

	srv := ball.NewServer()
	go srv.Listen()

	//g := &Game{}
	//g.Init()

	if err := ebiten.RunGame(srv); err != nil {
		log.Fatal(err)
	}
}
