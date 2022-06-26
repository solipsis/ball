package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type input int

const (
	NONE input = iota
	LEFT
	RIGHT
	UP
	UP_LEFT
	UP_RIGHT
)

type player struct {
	pos    vec2
	vel    vec2
	radius float64
	mass   float64
	input  input
}

func (p *player) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.pos.x-p.radius, p.pos.y-p.radius)
	screen.DrawImage(playerImage, op)
}

//	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
func (p *player) update() {
	if p.input == LEFT {
		p.vel.x = -5.0
	} else if p.input == RIGHT {
		p.vel.x = 5.0
	} else {
		p.vel.x = 0.0
	}

	p.pos.x += p.vel.x

}

func collidePlayer(p *player, b *ball) {

	// 1. find unit normal and unit tangent vectors
	n := vec2{x: b.pos.x - p.pos.x, y: b.pos.y - p.pos.y}
	magnitude := math.Sqrt((n.x * n.x) + n.y*n.y)
	un := vec2{x: n.x / magnitude, y: n.y / magnitude}
	ut := vec2{x: -un.y, y: un.x}

	// 2. Create the initial(before the collision velocity vectors
	// already done

	// 3. Project velocity vectors onto unit normal and unit tangent vectors
	v1n := dot(un, p.vel)
	v1t := dot(ut, p.vel)
	v2n := dot(un, b.vel)
	v2t := dot(ut, b.vel)

	// 4. Find new tangential velocities after the collision
	// same as original because no force between circles in the tangential direction
	// v'1t = v1t    v'2t = v2t
	v1tPrime := v1t
	v2tPrime := v2t

	// 5. find new normal velocities
	v1nPrime := (v1n*(p.mass-b.mass) + (2 * b.mass * v2n)) / (p.mass + b.mass)
	v2nPrime := (v2n*(b.mass-p.mass) + (2 * p.mass * v1n)) / (p.mass + b.mass)

	// 6. convert scalar normal and tangential velocities into vectors
	v1nPrimeVec := vec2{x: v1nPrime * un.x, y: v1nPrime * un.y}
	v1tPrimeVec := vec2{x: v1tPrime * ut.x, y: v1tPrime * ut.y}
	v2nPrimeVec := vec2{x: v2nPrime * un.x, y: v2nPrime * un.y}
	v2tPrimeVec := vec2{x: v2tPrime * ut.x, y: v2tPrime * ut.y}

	// 7. Find final velocity vectors by adding normal and tangential components for each
	v1Prime := vec2{x: v1nPrimeVec.x + v1tPrimeVec.x, y: v1nPrimeVec.y + v1tPrimeVec.y}
	v2Prime := vec2{x: v2nPrimeVec.x + v2tPrimeVec.x, y: v2nPrimeVec.y + v2tPrimeVec.y}

	p.vel = v1Prime
	b.vel = v2Prime
}
