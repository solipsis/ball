package ball

import "github.com/hajimehoshi/ebiten/v2"

type state struct {
	frame   int
	ball    ball
	players []player
	//player1 player
	//player2 player
}

func newState() state {
	return state{
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
}

// TODO: ONLY IMPLEMENTED FOR CLIENTS
func readInput(s *state) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		s.players[0].input.Dir = RIGHT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		s.players[0].input.Dir = LEFT
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		s.players[0].input.Dir = UP
	} else {
		s.players[0].input.Dir = NONE
	}
}

func step(s state, inputBuffer [][]input) state {

	for idx, player := range s.players {
		// if input is not authoritative, apply same input as previous frame
		if inputBuffer[idx][s.frame%len(inputBuffer)].Predicted {
			prev := inputBuffer[idx][((s.frame-1)+len(inputBuffer))%len(inputBuffer)]
			inputBuffer[idx][s.frame%len(inputBuffer)] = input{
				Dir:       prev.Dir,
				Predicted: true,
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
