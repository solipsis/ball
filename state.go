package main

import "github.com/hajimehoshi/ebiten/v2"

type state struct {
	frame   int
	ball    ball
	player1 player
	player2 player
}

// TODO: ONLY IMPLEMENTED FOR CLIENTS
func readInput(s *state) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		s.player1.input = RIGHT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		s.player1.input = LEFT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		s.player1.input = UP
	} else {
		s.player1.input = NONE
	}
}

func step(s state) state {
	s.frame += 1
	s.player1.update()
	s.player2.update()
	s.ball.update()

	distBetweenBallPlayer1 := ((s.player1.pos.x - s.ball.pos.x) * (s.player1.pos.x - s.ball.pos.x)) + ((s.player1.pos.y - s.ball.pos.y) * (s.player1.pos.y - s.ball.pos.y))
	if distBetweenBallPlayer1 <= ((s.ball.radius + s.player1.radius) * (s.ball.radius + s.player1.radius)) {
		collidePlayer(&s.player1, &s.ball)
	}
	distBetweenBallPlayer2 := ((s.player2.pos.x - s.ball.pos.x) * (s.player2.pos.x - s.ball.pos.x)) + ((s.player2.pos.y - s.ball.pos.y) * (s.player2.pos.y - s.ball.pos.y))
	if distBetweenBallPlayer2 <= ((s.ball.radius + s.player2.radius) * (s.ball.radius + s.player2.radius)) {
		collidePlayer(&s.player2, &s.ball)
	}

	return s
}
