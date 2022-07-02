package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const clientBuffer = 5

type Client struct {
	ID             int
	state          state
	lastInputFrame int
	inputBuffer    [][]input
	conn           *websocket.Conn

	serverUpdates chan state
}

/*
type inputUpdate struct {
	Frame int
	Input input
}
*/

func (c *Client) readInput() input {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		return input{dir: RIGHT, predicted: false}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		return input{dir: LEFT, predicted: false}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		return input{dir: UP, predicted: false}
	} else {
		return input{dir: NONE, predicted: false}
	}
}

func (c *Client) update() {
	in := c.readInput()
	c.inputBuffer[c.ID][c.state.frame%len(c.inputBuffer)] = in

	// send input
	go func() {
		if err := wsjson.Write(context.TODO(), c.conn,
			playerUpdate{
				Input: in,
				Frame: c.state.frame,
			}); err != nil {
			log.Fatalf("writing input: %v", err)
		}
	}()

	// read updates
	for {
		select {
		case up := <-c.serverUpdates:
			fmt.Println("received update")
			c.handleRollback(up)
		default:
			break
		}
	}

	// advance frame
	c.state.players[c.ID].input = in
	c.state = step(c.state, c.inputBuffer)
}

func (c *Client) handleRollback(rollback state) {
	currentFrame := c.state.frame
	fmt.Printf("Rolled back from: %d to %d\n", currentFrame, rollback.frame)

	for rollback.frame < currentFrame {
		// override input
		rollback.players[c.ID].input = c.inputBuffer[c.ID][rollback.frame%len(c.inputBuffer)]

		rollback = step(rollback, c.inputBuffer)
	}
	c.state = rollback
}

// go function that receives all input dumps to channel
// start of update loop reads all from that channel
//  while select { <- , default break}

func (c *Client) Run(url string) {
	c.ID = 0 // TODO: Get value from server
	c.serverUpdates = make(chan state)
	c.inputBuffer = make([][]input, 2)
	c.inputBuffer[0] = make([]input, 60)
	c.inputBuffer[1] = make([]input, 60)

	conn, _, err := websocket.Dial(context.TODO(), url, nil)
	if err != nil {
		log.Printf("client websocket close: %v", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "client socket closing")
	c.conn = conn

	go func() {
		// read frame
		_, stateBuf, err := c.conn.Read(context.TODO())
		if err != nil {
			log.Printf("client ws read: %v", err)
			return
		}
		var s state
		if err := json.Unmarshal(stateBuf, &s); err != nil {
			log.Printf("client unmarshal: %v", err)
			return
		}

		c.serverUpdates <- s
	}()
}

// implements ebiten.Game
func (c *Client) Update() error {
	c.update()
	return nil
}

// implements ebiten.Game
func (c *Client) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.state.ball.pos.x-c.state.ball.radius, c.state.ball.pos.y-c.state.ball.radius)
	screen.DrawImage(ballImage, op)
	op = &ebiten.DrawImageOptions{}
	DrawRect(screen, 100, 200, 20, 50, color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0x00})
	c.state.players[0].draw(screen)
}

// implements ebiten.Game
func (c *Client) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
