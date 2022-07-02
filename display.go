package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type EbitenWrapper struct {
	State        *state
	ScreenHeight int
	ScreenWidth  int

	UpdateFunc func() error
	DrawFunc   func(screen *ebiten.Image)
}

func (w *EbitenWrapper) Update() error {
	w.UpdateFunc()

	return nil
}

func (w *EbitenWrapper) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.State.ball.pos.x-w.State.ball.radius, w.State.ball.pos.y-w.State.ball.radius)
	screen.DrawImage(ballImage, op)
	op = &ebiten.DrawImageOptions{}
	DrawRect(screen, 100, 200, 20, 50, color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0x00})
	w.State.players[0].draw(screen)
}

func (w *EbitenWrapper) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
