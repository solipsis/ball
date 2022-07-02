package main

import "github.com/hajimehoshi/ebiten/v2"

type state struct {
	frame   int
	ball    ball
	players []player
	//player1 player
	//player2 player
}

// TODO: ONLY IMPLEMENTED FOR CLIENTS
func readInput(s *state) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		s.players[0].input.dir = RIGHT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		s.players[0].input.dir = LEFT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		s.players[0].input.dir = UP
	} else {
		s.players[0].input.dir = NONE
	}
}

func step(s state, inputBuffer [][]input) state {

	for idx, player := range s.players {
		// if input is not authoritative, apply same input as previous frame
		if inputBuffer[idx][s.frame%len(inputBuffer)].predicted {
			prev := inputBuffer[idx][((s.frame-1)+len(inputBuffer))%len(inputBuffer)]
			inputBuffer[idx][s.frame%len(inputBuffer)] = input{
				dir:       prev.dir,
				predicted: true,
			}
		}

		s.players[idx].update()
		distBetweenBallPlayer := ((player.pos.x - s.ball.pos.x) * (player.pos.x - s.ball.pos.x)) + ((player.pos.y - s.ball.pos.y) * (player.pos.y - s.ball.pos.y))
		if distBetweenBallPlayer <= ((s.ball.radius + player.radius) * (s.ball.radius + player.radius)) {
			collidePlayer(&s.players[idx], &s.ball)
		}
	}
	s.ball.update()
	s.frame += 1

	return s
}
