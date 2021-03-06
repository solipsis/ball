package ball

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var ballImage *ebiten.Image
var playerImage *ebiten.Image
var gravity = vec2{x: 0.0, y: 0.13}
var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.Black)
	f, err := os.Open("poke3.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	ballImage = ebiten.NewImageFromImage(img)

	f, err = os.Open("semi3.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err = png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	playerImage = ebiten.NewImageFromImage(img)
}

type vec2 struct {
	x float64
	y float64
}

type Game struct {
	state state
}

func dot(a, b vec2) float64 {
	return (a.x * b.x) + (a.y * b.y)
}

func (g *Game) Update() error {
	readInput(&g.state)
	//g.state = step(g.state)
	emptyImage.Fill(color.White)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.state.ball.pos.x-g.state.ball.radius, g.state.ball.pos.y-g.state.ball.radius)
	screen.DrawImage(ballImage, op)
	op = &ebiten.DrawImageOptions{}
	DrawRect(screen, 100, 200, 20, 50, color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0x00})
	g.state.players[0].draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Init() error {

	emptyImage.Fill(color.Black)
	f, err := os.Open("poke3.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	ballImage = ebiten.NewImageFromImage(img)

	f, err = os.Open("semi3.png")
	if err != nil {
		log.Fatal(err)
	}
	img, err = png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	playerImage = ebiten.NewImageFromImage(img)
	g.state = state{
		ball: ball{
			pos:    vec2{x: 50, y: 50},
			vel:    vec2{x: 2.5, y: 2.7},
			acc:    gravity,
			radius: 18,
			mass:   10,
		},
		players: []player{
			player{
				pos:    vec2{x: screenWidth / 2, y: screenHeight - 1},
				radius: 95.0 / 2.0,
				mass:   10000000,
			},
			player{
				pos:    vec2{x: -1000, y: -1000},
				radius: 95.0 / 2.0,
				mass:   10000000,
			},
		},
	}

	return nil
}

const BUFFER_SIZE = 60

func main() {
	fmt.Println("------------------------")

	l, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		log.Fatalf("server net.Listen(): %v", err)
	}
	log.Printf("listening on http://%v", l.Addr())

	srv := Server{
		clientInputs: make(chan playerUpdate),
		inputBuffer:  make([][]input, 2),
		stateBuffer:  make([]state, 60),
	}
	srv.inputBuffer[0] = make([]input, 60)
	srv.inputBuffer[1] = make([]input, 60)
	s := &http.Server{
		Handler:      &srv,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	//	go func() {
	fmt.Println("Serving")
	if err := s.Serve(l); err != nil {
		log.Fatalf("serve: %v", err)
	}
	//	}()

	fmt.Println("Sleeping...")
	//time.Sleep(1 * time.Second)
	fmt.Println("waking up")

	fmt.Println("pre run")
	// TODO: this never terminates
	/*
		go c.Run("http://localhost:8090")
		fmt.Println("post Run")

		if err := ebiten.RunGame(c); err != nil {
			log.Fatal(err)
		}
		return
	*/

	/*
		emptyImage.Fill(color.Black)
		f, err := os.Open("poke3.png")
		if err != nil {
			log.Fatal(err)
		}
		img, err := png.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		ballImage = ebiten.NewImageFromImage(img)

		f, err = os.Open("semi3.png")
		if err != nil {
			log.Fatal(err)
		}
		img, err = png.Decode(f)
		if err != nil {
			log.Fatal(err)
		}
		playerImage = ebiten.NewImageFromImage(img)

		ebiten.SetWindowSize(screenWidth, screenHeight)
	*/
	g := &Game{}
	g.Init()

	/*
		g := &Game{
			state: state{
				ball: ball{
					pos:    vec2{x: 50, y: 50},
					vel:    vec2{x: 2.5, y: 2.7},
					acc:    gravity,
					radius: 18,
					mass:   10,
				},
				player1: player{
					pos:    vec2{x: screenWidth / 2, y: screenHeight - 1},
					radius: 95.0 / 2.0,
					mass:   10000000,
				},
				player2: player{
					pos:    vec2{x: -1000, y: -1000},
					radius: 95.0 / 2.0,
					mass:   10000000,
				},
			},
		}
	*/

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}

// DrawRect draws a rectangle on the given destination dst.
//
// DrawRect is intended to be used mainly for debugging or prototyping purpose.
func DrawRect(dst *ebiten.Image, x, y, width, height float64, clr color.Color) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(width, height)
	op.GeoM.Translate(x, y)
	//emptyImage.Fill(color.White)	op.ColorM.ScaleWithColor(clr)
	// Filter must be 'nearest' filter (default).
	// Linear filtering would make edges blurred.
	dst.DrawImage(emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
}
