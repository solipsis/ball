package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const clientBuffer = 5

type client struct {
	ID             int
	state          state
	lastInputFrame int
	inputBuffer    []input
	conn           *websocket.Conn
}

type inputUpdate struct {
	//Frame int
	Input input
}

func (c *client) readInput() input {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		return RIGHT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		return LEFT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		return UP
	} else {
		return NONE
	}
}

func (c *client) run() {

	input := c.readInput()

	// send input
	go func() {
		if err := wsjson.Write(context.TODO(), c.conn,
			inputUpdate{
				Input: input,
			}); err != nil {

			log.Fatalf("writing input: %v", err)
		}
	}()

	// read frame
	_, state, err := c.conn.Read(context.TODO())
	if err != nil {
		log.Fatalf("client ws read: %v", err)
	}
	fmt.Println(state)

	//	if c.lastInputFrame < state.Frame+clientBuffer {
	// send input

	//	}

}
