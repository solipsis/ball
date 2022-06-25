package main

import (
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 600
	screenHeight = 600
)

var ballImage *ebiten.Image
var playerImage *ebiten.Image
var gravity = vec2{x: 0.0, y: 0.13}

type vec2 struct {
	x float64
	y float64
}

type Game struct {
	ball    *ball
	ball2   *ball
	player1 *player
}

type ball struct {
	pos    vec2
	vel    vec2
	acc    vec2
	radius float64
	mass   float64
}

func dot(a, b vec2) float64 {
	return (a.x * b.x) + (a.y * b.y)
}

func collide(b1 *ball, b2 *ball) {

	/*
		b1 = &ball{
			pos:  vec2{x: 50, y: 0},
			vel:  vec2{x: 10, y: 10.0},
			mass: 10.0,
		}
		b2 = &ball{
			pos:  vec2{x: 70, y: 10},
			vel:  vec2{x: -5.0, y: -3.0},
			mass: 10.0,
		}
	*/
	//fmt.Println("------------------------------")

	// 1. find unit normal and unit tangent vectors
	n := vec2{x: b2.pos.x - b1.pos.x, y: b2.pos.y - b1.pos.y}
	magnitude := math.Sqrt((n.x * n.x) + n.y*n.y)
	un := vec2{x: n.x / magnitude, y: n.y / magnitude}
	ut := vec2{x: -un.y, y: un.x}

	//fmt.Println(n, magnitude, un, ut)

	// 2. Create the initial(before the collision velocity vectors
	// already done

	// 3. Project velocity vectors onto unit normal and unit tangent vectors
	v1n := dot(un, b1.vel)
	v1t := dot(ut, b1.vel)
	v2n := dot(un, b2.vel)
	v2t := dot(ut, b2.vel)

	// 4. Find new tangential velocities after the collision
	// same as original because no force between circles in the tangential direction
	// v'1t = v1t    v'2t = v2t
	v1tPrime := v1t
	v2tPrime := v2t

	// 5. find new normal velocities
	v1nPrime := (v1n*(b1.mass-b2.mass) + (2 * b2.mass * v2n)) / (b1.mass + b2.mass)
	v2nPrime := (v2n*(b2.mass-b1.mass) + (2 * b1.mass * v1n)) / (b1.mass + b2.mass)

	// 6. convert scalar normal and tangential velocities into vectors
	v1nPrimeVec := vec2{x: v1nPrime * un.x, y: v1nPrime * un.y}
	v1tPrimeVec := vec2{x: v1tPrime * ut.x, y: v1tPrime * ut.y}
	v2nPrimeVec := vec2{x: v2nPrime * un.x, y: v2nPrime * un.y}
	v2tPrimeVec := vec2{x: v2tPrime * ut.x, y: v2tPrime * ut.y}

	// 7. Find final velocity vectors by adding normal and tangential components for each
	v1Prime := vec2{x: v1nPrimeVec.x + v1tPrimeVec.x, y: v1nPrimeVec.y + v1tPrimeVec.y}
	v2Prime := vec2{x: v2nPrimeVec.x + v2tPrimeVec.x, y: v2nPrimeVec.y + v2tPrimeVec.y}

	//fmt.Println(b1.vel, b2.vel)
	b1.vel = v1Prime
	b2.vel = v2Prime
	//	fmt.Println("final")
	//	fmt.Println(b1.vel, b2.vel)
}

func (b *ball) update() {
	b.vel.x += b.acc.x
	b.vel.y += b.acc.y

	b.pos.x += b.vel.x
	b.pos.y += b.vel.y

	mag := (b.vel.x * b.vel.x) + (b.vel.y * b.vel.y)
	if mag > 12*12 {
		b.vel.x *= 0.9
		b.vel.y *= 0.9
	}

	if b.pos.x-b.radius < 0 {
		b.vel.x = b.vel.x * -1.0
		b.pos.x = b.radius
	}
	if b.pos.x+b.radius > screenWidth {
		b.vel.x = b.vel.x * -1.0
		b.pos.x = screenWidth - b.radius
	}
	if b.pos.y+b.radius > screenHeight {
		b.vel.y = b.vel.y * -1.0
		b.pos.y = screenHeight - b.radius
	}
	if b.pos.y-b.radius < 0 {
		b.vel.y = b.vel.y * -1.0
		b.pos.y = b.radius
	}
}

func (g *Game) Update() error {
	/*
		g.ball.vel.x += gravity.x
		g.ball.vel.y += gravity.y

		g.ball.pos.x += g.ball.vel.x
		g.ball.pos.y += g.ball.vel.y

		if g.ball.pos.y+g.ball.radius > screenHeight {
			g.ball.vel.y = g.ball.vel.y * -1.0
			g.ball.pos.y = screenHeight - g.ball.radius
		}
	*/
	g.ball.update()
	//g.ball2.update()
	g.player1.update()

	/*
		distBetweenBallsSquared := ((g.ball2.pos.x - g.ball.pos.x) * (g.ball2.pos.x - g.ball.pos.x)) + ((g.ball2.pos.y - g.ball.pos.y) * (g.ball2.pos.y - g.ball.pos.y))
		if distBetweenBallsSquared < ((g.ball.radius + g.ball2.radius) * (g.ball.radius + g.ball2.radius)) {
			collide(g.ball, g.ball2)
		}
	*/

	distBetweenBallPlayer := ((g.player1.pos.x - g.ball.pos.x) * (g.player1.pos.x - g.ball.pos.x)) + ((g.player1.pos.y - g.ball.pos.y) * (g.player1.pos.y - g.ball.pos.y))
	if distBetweenBallPlayer < ((g.ball.radius + g.player1.radius) * (g.ball.radius + g.player1.radius)) {
		//	fmt.Println("b1p")
		collidePlayer(g.player1, g.ball)
	}
	/*
		distBetweenBallPlayer = ((g.player1.pos.x - g.ball2.pos.x) * (g.player1.pos.x - g.ball2.pos.x)) + ((g.player1.pos.y - g.ball2.pos.y) * (g.player1.pos.y - g.ball2.pos.y))
		if distBetweenBallPlayer < ((g.ball2.radius + g.player1.radius) * (g.ball2.radius + g.player1.radius)) {
			//	fmt.Println("b2p")
			collidePlayer(g.player1, g.ball2)
		}
	*/

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.ball.pos.x-g.ball.radius, g.ball.pos.y-g.ball.radius)
	screen.DrawImage(ballImage, op)
	op = &ebiten.DrawImageOptions{}
	/*
		op.GeoM.Translate(g.ball2.pos.x-g.ball2.radius, g.ball2.pos.y-g.ball2.radius)
		screen.DrawImage(ballImage, op)
	*/
	g.player1.draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
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

	//collide(&ball{}, &ball{})
	//	return

	g := &Game{
		ball: &ball{
			pos:    vec2{x: 50, y: 50},
			vel:    vec2{x: 2.5, y: 2.7},
			acc:    gravity,
			radius: 18,
			mass:   10,
		},
		ball2: &ball{
			pos: vec2{x: 200, y: 200},
			vel: vec2{x: 4.5, y: 3.7},
			//		acc:    gravity,
			radius: 18,
			mass:   10,
		},
		player1: &player{
			pos:    vec2{x: screenWidth / 2, y: screenHeight - 1},
			radius: 95.0 / 2.0,
			mass:   10000000,
		},
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
